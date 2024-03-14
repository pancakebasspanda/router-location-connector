package integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"router-location-connecter/api"
	"router-location-connecter/app"
	"router-location-connecter/storage"
)

func Test_app_Process(t *testing.T) {
	ctx := context.Background()

	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	redisHandler, err := storage.New(nil, _redisAddress, _redisPassword)
	if err != nil {
		assert.NoError(t, err)
	}

	// serverURL running in dockerfile
	apiClient := api.New(api.WithMaxRetries(2),
		api.WithBaseURL("http://localhost:1090/test/router-location-data"),
		api.WithTimeout(time.Duration(10)*time.Second))

	tests := []struct {
		name     string
		api      api.API
		redis    storage.Storage
		log      zerolog.Logger
		expected string
	}{
		{
			name:     "process the API data by calling a server, store router and location data and print the linked router locations",
			api:      apiClient,
			redis:    redisHandler,
			log:      logger,
			expected: "[Williamson Park] <-> [Birmingham Hippodrome]\n[Loughborough University] <-> [Lancaster Brewery]\n[Lancaster Brewery] <-> [Lancaster University]\n[Loughborough University] <-> [Lancaster Castle]\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Redirect stdout to buffer using pipe
			r, w, err := os.Pipe()
			if err != nil {
				assert.NoError(t, err)
			}
			origStdout := os.Stdout
			os.Stdout = w

			a := app.NewApp(tt.api, tt.redis, tt.log)

			a.Process(ctx)

			buf := make([]byte, 1024)
			n, err := r.Read(buf)
			if err != nil {
				assert.NoError(t, err)
			}

			os.Stdout = origStdout

			assert.Equal(t, string(buf[:n]), tt.expected)

		})
	}

}
