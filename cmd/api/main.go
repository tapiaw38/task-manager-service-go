package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources"
	firestoreclient "github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources/firestore"
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources/jsonstore"
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/web"
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/web/middlewares"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/appcontext"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/config"
	"github.com/tapiaw38/task-manager-service-go/internal/usecases"
)

const shutdownTimeout = 10 * time.Second

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := initConfig(); err != nil {
		slog.Error("configuration initialization failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := run(); err != nil {
		slog.Error("server stopped with error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run() error {
	configService := config.GetConfigService()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	store, firestoreClient, err := newDatasources(ctx, configService)
	if err != nil {
		return err
	}

	if firestoreClient != nil {
		defer func() {
			if closeErr := firestoreClient.Close(); closeErr != nil {
				slog.Error("closing firestore client", slog.String("error", closeErr.Error()))
			}
		}()
	}

	if configService.ServerConfig.GinMode == config.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()
	app.Use(gin.Recovery(), middlewares.Logging(), newCORS(configService))

	if err := bootstrap(app, store, firestoreClient, &configService); err != nil {
		return err
	}

	server := &http.Server{
		Addr:    ":" + configService.ServerConfig.Port,
		Handler: app,
	}

	serverErrors := make(chan error, 1)

	go func() {
		slog.Info("server listening",
			slog.String("application", configService.AppName),
			slog.String("version", configService.AppVersion),
			slog.String("port", configService.ServerConfig.Port),
			slog.String("store", configService.StoreConfig.FilePath),
		)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		return err
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}

func newDatasources(
	ctx context.Context,
	configService config.ConfigurationService,
) (*jsonstore.Store, *firestoreclient.Client, error) {
	if configService.StoreConfig.Driver == config.FirestoreDriver {
		client, err := firestoreclient.New(
			ctx,
			configService.StoreConfig.FirestoreProjectID,
			configService.StoreConfig.FirestoreCollection,
		)
		if err != nil {
			return nil, nil, err
		}

		return nil, client, nil
	}

	store, err := jsonstore.New(configService.StoreConfig.FilePath)
	if err != nil {
		return nil, nil, err
	}

	return store, nil, nil
}

func bootstrap(
	app *gin.Engine,
	store *jsonstore.Store,
	firestoreClient *firestoreclient.Client,
	configService *config.ConfigurationService,
) error {
	datasources := datasources.CreateDatasources(store, firestoreClient)
	contextFactory := appcontext.NewFactory(datasources, configService)
	useCases := usecases.CreateUsecases(contextFactory)

	web.RegisterApplicationRoutes(app, useCases, *configService)

	return nil
}

func newCORS(configService config.ConfigurationService) gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = configService.ServerConfig.AllowedOrigins
	corsConfig.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}

	return cors.New(corsConfig)
}
