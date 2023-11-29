package main

import (
	"context"
	"fmt"
	"log"

	"github.com/samber/lo"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()

	const (
		user   = "user"
		pass   = "pass"
		dbname = "dbname"
	)

	// Стартуем Postgres контейнер
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForLog("PostgreSQL init process complete"),
		Env: map[string]string{
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": pass,
			"POSTGRES_DB":       dbname,
		},
		SkipReaper: true,
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("could not start postgres container: %s", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("could not get mapped port: %s", err)
	}
	// Подключаемся к базе данных
	dsn := fmt.Sprintf(
		"host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable",
		port.Port(), user, pass, dbname,
	)
	db := lo.Must(gorm.Open(postgres.Open(dsn)))

	type customer struct {
		ID   int
		Name string
	}
	err = db.AutoMigrate(customer{})
	if err != nil {
		log.Fatalln(err)
	}

	lo.Must0(db.Create(&customer{Name: "Vasya"}).Error)
	lo.Must0(db.Create(&customer{Name: "Petya"}).Error)
	lo.Must0(db.Create(&customer{Name: "Kolya"}).Error)

	var users []customer
	lo.Must0(db.Find(&users).Error)

	fmt.Println(users)

	log.Println("Test finished!")
}
