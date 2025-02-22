package test

import (
	"MerchStore/src/cmd"
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/redis/go-redis/v9"
	"net/http"
	"path/filepath"
	"runtime"
)

func LoadEnv() error {
	envFilePath := filepath.Join("..", ".env")
	if err := godotenv.Load(envFilePath); err != nil {
		return fmt.Errorf("ошибка загрузки .env: %v", err)
	}
	return nil
}

func DeployRedis(cfg *cmd.Config, pool *dockertest.Pool, network *dockertest.Network) (*dockertest.Resource, error) {
	port := docker.Port(cfg.Cache.Port)
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Hostname:     cfg.Cache.Host,
		Repository:   "redis",
		Tag:          "latest",
		ExposedPorts: []string{string(port)},
		PortBindings: map[docker.Port][]docker.PortBinding{
			port: {{HostIP: "", HostPort: cfg.Cache.HostPort}},
		},
		Networks: []*dockertest.Network{network},
	})
	if err != nil {
		return nil, fmt.Errorf("could not start Redis: %v", err)
	}

	if err := pool.Retry(func() error {
		fmt.Println("Checking Redis connection...")
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("localhost:%s", cfg.Cache.HostPort),
			Password: "",
			DB:       0,
		})
		defer func(client *redis.Client) {
			err := client.Close()
			if err != nil {
			}
		}(client)
		_, err := client.Ping(context.Background()).Result()
		return err
	}); err != nil {
		return nil, fmt.Errorf("could not connect to Redis: %v", err)
	}

	return resource, nil
}

func getInitSQLPath() string {
	_, filename, _, _ := runtime.Caller(0)
	rootDir := filepath.Join(filepath.Dir(filename), "..", "..")
	return filepath.Join(rootDir, "migrations", "init.sql")
}

func DeployPostgres(cfg *cmd.Config, pool *dockertest.Pool, network *dockertest.Network) (*dockertest.Resource, error) {
	port := docker.Port(cfg.Database.Port)
	migrationPath := getInitSQLPath()
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Hostname:     cfg.Database.Host,
		Repository:   "postgres",
		Tag:          "13",
		ExposedPorts: []string{string(port)},
		PortBindings: map[docker.Port][]docker.PortBinding{
			port: {{HostIP: "", HostPort: cfg.Database.HostPort}},
		},
		Networks: []*dockertest.Network{network},
		Env: []string{
			"POSTGRES_USER=" + cfg.Database.User,
			"POSTGRES_PASSWORD=" + cfg.Database.Password,
			"POSTGRES_DB=" + cfg.Database.Name,
		},
		Mounts: []string{
			fmt.Sprintf("%s:/docker-entrypoint-initdb.d/init.sql", migrationPath),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("could not start PostgreSQL: %v", err)
	}

	if err := pool.Retry(func() error {
		fmt.Println("Checking PostgreSQL connection...")
		dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
			cfg.Database.User, cfg.Database.Password, cfg.Database.HostPort, cfg.Database.Name)
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return err
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
			}
		}(db)
		return db.Ping()
	}); err != nil {
		return nil, fmt.Errorf("could not connect to PostgreSQL: %v", err)
	}

	return resource, nil
}

func DeployAPIContainer(cfg *cmd.Config, pool *dockertest.Pool, network *dockertest.Network) (*dockertest.Resource, error) {
	port := docker.Port(cfg.Service.ServerPort)
	envVars := []string{
		"SERVER_PORT=" + cfg.Service.ServerPort,
		"LOG_LEVEL=" + cfg.Service.LogLevel,
		"SECRET_KEY=" + cfg.Service.SecretKey,
		"DATABASE_USER=" + cfg.Database.User,
		"DATABASE_PASSWORD=" + cfg.Database.Password,
		"DATABASE_HOST=" + cfg.Database.Host,
		"DATABASE_PORT=" + cfg.Database.Port,
		"DATABASE_NAME=" + cfg.Database.Name,
		"CACHE_HOST=" + cfg.Cache.Host,
		"CACHE_PORT=" + cfg.Cache.Port,
		"HOST_CACHE_PORT=" + cfg.Cache.HostPort,
		"HOST_DATABASE_PORT=" + cfg.Database.HostPort,
	}

	resource, err := pool.BuildAndRunWithBuildOptions(&dockertest.BuildOptions{
		ContextDir: "../../../",
		Dockerfile: "DockerfileTest",
	}, &dockertest.RunOptions{
		Name:         "api-container",
		ExposedPorts: []string{string(port)},
		PortBindings: map[docker.Port][]docker.PortBinding{
			port: {{HostIP: "0.0.0.0", HostPort: cfg.Service.ServerPort}},
		},
		Networks: []*dockertest.Network{network},
		Env:      envVars,
	})

	if err != nil {
		return nil, fmt.Errorf("could not start API container: %v", err)
	}

	if err = pool.Retry(func() error {
		fmt.Println("Checking API connection...")
		url := fmt.Sprintf("http://localhost:%s/api/info", cfg.Service.ServerPort)
		_, err := http.Get(url)
		return err
	}); err != nil {
		return nil, fmt.Errorf("API container is not responding: %v", err)
	}

	return resource, nil
}

func TearDown(pool *dockertest.Pool, resources []*dockertest.Resource, network *dockertest.Network) error {
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
