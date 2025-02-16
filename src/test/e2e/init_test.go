package e2e

import (
	"MerchStore/src/cmd"
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Declare a global variable to hold the Docker pool and resource.
var (
	network *dockertest.Network
)

func TestMain(m *testing.M) {
	// Загружаем переменные окружения
	if err := loadEnv(); err != nil {
		log.Fatalf("Could not load env: %v", err)
	}
	cfg, err := cmd.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	var resources []*dockertest.Resource
	pool, err1 := dockertest.NewPool("")
	if err1 != nil {
		log.Fatalf("Could not connect to Docker: %v", err1)
	}

	// Создаём тестовую сеть
	network, err1 = pool.CreateNetwork("test-network")
	if err1 != nil {
		log.Fatalf("Could not create network: %v", err1)
	}

	// Запускаем контейнеры
	dbResource, err1 := deployPostgres(pool)
	if err1 != nil {
		TearDown(pool, resources)
		log.Fatalf("Could not start PostgreSQL: %v", err1)
	}
	resources = append(resources, dbResource)

	redisResource, err1 := deployRedis(pool)
	if err1 != nil {
		TearDown(pool, resources)
		log.Fatalf("Could not start Redis: %v", err1)
	}
	resources = append(resources, redisResource)

	apiResource, err1 := deployAPIContainer(cfg, pool)
	if err1 != nil {
		TearDown(pool, resources)
		log.Fatalf("Could not start API container: %v", err1)
	}
	resources = append(resources, apiResource)

	// Запуск тестов
	exitCode := m.Run()

	// Удаляем контейнеры
	if err := TearDown(pool, resources); err != nil {
		log.Fatalf("Could not purge resource: %v", err)
	}
	os.Exit(exitCode)
}

// Функция загрузки переменных окружения из .env файла
func loadEnv() error {
	// Определяем путь к .env файлу
	envFilePath := filepath.Join("..", ".env") // Путь относительно текущей директории

	// Загружаем переменные окружения из файла .env
	err := godotenv.Load(envFilePath)
	if err != nil {
		return fmt.Errorf("Error loading .env file")
	}
	return nil
}

// Redis на порту 6389
func deployRedis(pool *dockertest.Pool) (*dockertest.Resource, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Hostname:     "redis-container",
		Repository:   "redis",
		Tag:          "latest",
		ExposedPorts: []string{"6379/tcp"}, // Внутри контейнера порт 6379
		PortBindings: map[docker.Port][]docker.PortBinding{
			"6379/tcp": {{HostIP: "", HostPort: "6389"}}, // Пробрасываем на хост-порт 6389
		},
		Networks: []*dockertest.Network{network},
	})
	if err != nil {
		return nil, fmt.Errorf("could not start Redis: %v", err)
	}

	if err := pool.Retry(func() error {
		fmt.Println("Checking Redis connection...")
		client := redis.NewClient(&redis.Options{
			Addr:     "localhost:6389", // Подключаемся к хостовому порту 6389
			Password: "",
			DB:       0,
		})
		defer client.Close()

		_, err := client.Ping(context.Background()).Result()
		return err
	}); err != nil {
		return nil, fmt.Errorf("could not connect to Redis: %v", err)
	}

	return resource, nil
}

// PostgreSQL на порту 5440
func deployPostgres(pool *dockertest.Pool) (*dockertest.Resource, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Hostname:     "postgres",
		Repository:   "postgres",
		Tag:          "13",
		ExposedPorts: []string{"5432/tcp"}, // Внутри контейнера порт 5432
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {{HostIP: "", HostPort: "5440"}}, // Пробрасываем на хост-порт 5440
		},
		Networks: []*dockertest.Network{network},
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=password",
			"POSTGRES_DB=shop",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("could not start PostgreSQL: %v", err)
	}

	if err := pool.Retry(func() error {
		fmt.Println("Checking PostgreSQL connection...")
		dsn := "postgres://postgres:password@localhost:5440/shop?sslmode=disable" // Подключаемся к хостовому порту 5440
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return err
		}
		defer db.Close()

		time.Sleep(2 * time.Second) // Даем время на инициализацию
		return db.Ping()
	}); err != nil {
		return nil, fmt.Errorf("could not connect to PostgreSQL: %v", err)
	}

	return resource, nil
}

// TearDown purges the resources and removes the network.
func TearDown(pool *dockertest.Pool, resources []*dockertest.Resource) error {
	for _, resource := range resources {
		if err := pool.Purge(resource); err != nil {
			return fmt.Errorf("could not purge resource: %v", err)
		}
	}
	if network == nil {
		return nil
	}
	if err := pool.RemoveNetwork(network); err != nil {
		return fmt.Errorf("could not remove network: %v", err)
	}

	return nil
}

// deployAPIContainer builds and runs the API container.
func deployAPIContainer(cfg *cmd.Config, pool *dockertest.Pool) (*dockertest.Resource, error) {
	// Получаем переменные окружения
	envVars := []string{
		"SERVER_PORT=" + "8090",
		"LOG_LEVEL=" + cfg.Service.LogLevel,
		"SECRET_KEY=" + cfg.Service.SecretKey,
		"DATABASE_USER=" + cfg.Database.User,
		"DATABASE_PASSWORD=" + cfg.Database.Password,
		"DATABASE_HOST=" + cfg.Database.Host,
		"DATABASE_PORT=" + cfg.Database.Port,
		"DATABASE_NAME=" + cfg.Database.Name,
		"CACHE_HOST=" + cfg.Cache.Host,
		"CACHE_PORT=" + cfg.Cache.Port,
	}
	// build and run the API container
	resource, err := pool.BuildAndRunWithBuildOptions(&dockertest.BuildOptions{
		ContextDir: "../../../",
		Dockerfile: "DockerfileTest",
	}, &dockertest.RunOptions{
		Name:         "api-container",
		ExposedPorts: []string{"8090/tcp"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"8090/tcp": {{HostIP: "0.0.0.0", HostPort: "8090"}},
		},
		Networks: []*dockertest.Network{
			network,
		},
		Env: envVars,
	})

	if err != nil {
		return nil, fmt.Errorf("could not start resource: %v", err)
	}

	// check if the API container is ready to accept connections
	if err = pool.Retry(func() error {
		fmt.Println("Checking API connection...")
		_, err := http.Get("http://localhost:8090/api/info")
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("could not start resource: %v", err)
	}

	return resource, nil
}
