package config

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type App struct {
	Env   *Env
	DB    *mongo.Database
	Redis *redis.Client
}

func HealthChecker(app App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		dbErr := app.DB.Client().Ping(r.Context(), readpref.Primary())
		if dbErr != nil {
			log.Fatal(dbErr.Error())
			http.Error(w, "MongoDB not connected", http.StatusInternalServerError)
			return
		}

		_, clientErr := app.Redis.Ping(r.Context()).Result()
		if clientErr != nil {
			log.Fatal(clientErr.Error())
			http.Error(w, "Redis not connected", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Server Starter Successfully!")
	}
}
