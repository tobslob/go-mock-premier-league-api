package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"tobslob.com/go-mock-premier-league-api/pkg/config"
)

func main() {
	var err error
	var client *mongo.Client
	var db *mongo.Database
	var redisClient *redis.Client

	var env config.Env
	err = config.LoadEnv(&env)
	if err != nil {
		panic(err.Error())
	}

	ctx, cancel := config.WithCancel(context.Background())
	defer cancel()

	startupCtx, cancel := context.WithTimeout(ctx, time.Minute)

	client, db, err = config.ConnectMongo(startupCtx, env.MongodbURL, env.MongodbName)
	if err != nil {
		panic(err.Error())
	}

	defer func() {
		err = client.Disconnect(context.TODO())
		if err != nil {
			log.Fatal(err.Error())
		}
	}()
	log.Print("Successfully connected to mongodb")

	redisClient, err = config.ConnectRedis(startupCtx, env)
	if err != nil {
		panic(err.Error())
	}

	defer func() {
		err := redisClient.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()
	log.Print("Successfully connected to redis")

	app := config.App{
		DB:    db,
		Env:   &env,
		Redis: redisClient,
	}

	// API router
	router := chi.NewRouter()

	// setup routes

	// mount API on app router
	appRouter := chi.NewRouter()
	appRouter.Mount("/api/v1", router)
	appRouter.Get("/", config.HealthChecker(app))
	appRouter.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Route Not Found", http.StatusNotFound)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", env.Port),
		Handler: appRouter,
	}

	go func() {
		<-ctx.Done()

		shutCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err := server.Shutdown(shutCtx)
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	log.Printf("serving api at http://127.0.0.1:%d", env.Port)
	serverErr := server.ListenAndServe()
	if serverErr != http.ErrServerClosed {
		log.Print("could not start the server")
	}
	<-ctx.Done()
}
