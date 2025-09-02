package main

import (
	"context"
	"crud/internal/config"
	"crud/internal/db"
	"crud/internal/handler"
	"crud/internal/models"
	"crud/internal/repository"
	"crud/internal/router"
	"crud/internal/service"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	port = flag.Int("port", 8080, "server endpoint")
)

func main() {
	flag.Parse()
	ctx := context.Background()
	cfg := config.Load()

	// connect to postgreSQL
	database := db.ConnnectPostgreSQL(cfg)
	defer database.Close()

	// connect to redis
	rdb := db.ConnectRedis(cfg)

	// test
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(pong)

	// migrate
	_, err = database.NewCreateTable().Model((*models.User)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		panic(err)
	}
	// Run migrations
	// if err := migrations.InitMigrations(ctx, database); err != nil {
	// 	log.Fatalf("Failed to run migrations: %v", err)
	// }

	userRepo := repository.NewUserRepository(database)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	server := http.Server{
		Addr:         fmt.Sprintf(":%v", *port),
		Handler:      router.NewRouter(userHandler),
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}

	log.Fatal(server.ListenAndServe())
}
