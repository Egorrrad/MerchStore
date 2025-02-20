package e2e

import (
	"MerchStore/src/cmd"
	"MerchStore/src/test"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"log"
	"os"
	"testing"
)

var baseURL string

func TestMain(m *testing.M) {
	var (
		network *dockertest.Network
	)
	// Загружаем переменные окружения
	if err := test.LoadEnv(); err != nil {
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
	dbResource, err1 := test.DeployPostgres(cfg, pool, network)
	if err1 != nil {
		test.TearDown(pool, resources, network)
		log.Fatalf("Could not start PostgreSQL: %v", err1)
	}
	resources = append(resources, dbResource)

	redisResource, err1 := test.DeployRedis(cfg, pool, network)
	if err1 != nil {
		test.TearDown(pool, resources, network)
		log.Fatalf("Could not start Redis: %v", err1)
	}
	resources = append(resources, redisResource)

	apiResource, err1 := test.DeployAPIContainer(cfg, pool, network)
	if err1 != nil {
		test.TearDown(pool, resources, network)
		log.Fatalf("Could not start API container: %v", err1)
	}
	resources = append(resources, apiResource)

	baseURL = fmt.Sprintf("http://localhost:%s", cfg.Service.ServerPort)
	// Запуск тестов
	exitCode := m.Run()

	// Удаляем контейнеры
	if err := test.TearDown(pool, resources, network); err != nil {
		log.Fatalf("Could not purge resource: %v", err)
	}
	os.Exit(exitCode)
}
