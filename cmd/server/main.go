package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/lisvindanuu/indonesiaql/graph"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/lisvindanuu/indonesiaql/internal/cache"
	"github.com/lisvindanuu/indonesiaql/internal/config"
	"github.com/lisvindanuu/indonesiaql/internal/database"
	"github.com/lisvindanuu/indonesiaql/internal/middleware"
	"github.com/lisvindanuu/indonesiaql/internal/repository"
	"github.com/lisvindanuu/indonesiaql/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := database.NewPool(ctx, cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	c := cache.New(5*time.Minute, 10*time.Minute)

	wilayahRepo := repository.NewWilayahRepository(pool)
	hariLiburRepo := repository.NewHariLiburRepository(pool)

	wilayahSvc := service.NewWilayahService(wilayahRepo)
	hariLiburSvc := service.NewHariLiburService(hariLiburRepo)
	cuacaSvc := service.NewCuacaService(c)
	kursSvc := service.NewKursService(c)
	nikSvc := service.NewNIKService(wilayahRepo)
	kodeBankSvc := service.NewKodeBankService()
	platNomorSvc := service.NewPlatNomorService()
	waktuSholatSvc := service.NewWaktuSholatService(c, cuacaSvc)
	gempaSvc := service.NewGempaService()
	kodePosSvc := service.NewKodePosService(pool)
	kalenderHijriyahSvc := service.NewKalenderHijriyahService()
	hargaBBMSvc := service.NewHargaBBMService()
	ihsgSvc := service.NewIHSGService()
	krlSvc := service.NewKRLService()
	bpjsSvc := service.NewBPJSService()
	rekeningSvc := service.NewRekeningService()
	inflasiSvc := service.NewInflasiService()

	resolver := graph.NewResolver(wilayahSvc, hariLiburSvc, cuacaSvc, kursSvc, nikSvc, kodeBankSvc, platNomorSvc, waktuSholatSvc, gempaSvc, kodePosSvc, kalenderHijriyahSvc, hargaBBMSvc, ihsgSvc, krlSvc, bpjsSvc, rekeningSvc, inflasiSvc)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.GET{})
	srv.Use(extension.Introspection{})
	srv.Use(extension.FixedComplexityLimit(200))

	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		gqlErr := graphql.DefaultErrorPresenter(ctx, err)
		if !cfg.IsDev() {
			// Jangan bocorkan detail internal error ke client di production
			if gqlErr.Extensions == nil {
				return &gqlerror.Error{Message: gqlErr.Message}
			}
		}
		return gqlErr
	})

	srv.SetRecoverFunc(func(ctx context.Context, err any) error {
		slog.Error("graphql resolver panic", "error", err)
		return &gqlerror.Error{Message: "internal server error"}
	})

	mux := http.NewServeMux()

	if cfg.IsDev() {
		mux.Handle("/", playground.Handler("IndonesiaQL Playground", "/query"))
	}

	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}))

	mux.Handle("/query", srv)

	var h http.Handler = mux
	h = middleware.Recovery(h)
	h = middleware.Logging(h)
	h = middleware.RateLimit(cfg.RateLimitRPM)(h)
	h = middleware.CORS()(h)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      h,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server started", "port", cfg.Port, "env", cfg.Environment)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
	slog.Info("server stopped")
}
