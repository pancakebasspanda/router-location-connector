package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/rs/zerolog"

	"router-location-connecter/api"
	"router-location-connecter/app"
	"router-location-connecter/storage"
)

const (
	_errRedisClient = "redis client initialization error"
	_appName        = "router-location-connector"
)

var (
	maxRetries         int
	baseURL            string
	timeout            int64
	redisURL, redisPWD string
	persistData        bool
)

func init() {
	flag.StringVar(&baseURL, "base-url", "https://my-json-server.typicode.com/marcuzh/router_location_test_api/db", "base url to get router location data")
	flag.IntVar(&maxRetries, "retries", 3, "max retries")
	flag.Int64Var(&timeout, "timeout", 20, "time in seconds")
	flag.BoolVar(&persistData, "persist-data", false, "keep router location data between runs")
}

func main() {
	ctx := context.Background()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app_name", _appName).
		Logger()

	flag.Parse()

	apiClient := api.New(api.WithMaxRetries(maxRetries),
		api.WithBaseURL(baseURL),
		api.WithTimeout(time.Duration(timeout)*time.Second))

	// we don't pass these in as flags as we ideally would want to create a Kubernetes secret,
	// mount this secret into your Pods where the application
	//can read them as environmental variables
	redisURL = getEnv("REDIS_URL", "localhost:6379")

	redisPWD = getEnv("REDIS_PASSWORD", "")

	redisClient, err := storage.New(ctx, redisURL, redisPWD)
	if err != nil {
		log.Panic().Err(err).Msg(_errRedisClient)
	}

	runner := app.NewApp(apiClient, redisClient, log)

	runner.Process(ctx)

	// Close the Redis client after finishing
	if !persistData {
		if err := redisClient.FlushAll(ctx); err != nil {
			log.Error().Err(err).Msg("error flushing redis store")
		}
	}

	if err := redisClient.Close(); err != nil {
		log.Error().Err(err).Msg("error closing redis conn")
	}

}

// getEnv gets any environment variables that are set
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
