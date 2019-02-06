package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cedric-parisi/payment-api/pkg/auth"

	"github.com/jinzhu/gorm"

	"github.com/cedric-parisi/payment-api/docs"
	"github.com/cedric-parisi/payment-api/internal/config"

	"github.com/cedric-parisi/payment-api/internal/payments"
	"github.com/cedric-parisi/payment-api/internal/repository"

	opentracing "github.com/opentracing/opentracing-go"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	jwt "github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	kitlog "github.com/go-kit/kit/log"

	jaegercfg "github.com/uber/jaeger-client-go/config"

	// import postgres drivers
	_ "github.com/lib/pq"
)

const (
	appName          = "payment-api"
	gracefulPeriod   = 5 * time.Second
	httpReadTimeout  = 5 * time.Second
	httpWriteTimeout = 60 * time.Second
)

var (
	errorLogger      = kitlog.NewJSONLogger(os.Stderr)
	connectionString = "host=%s port=%s user=%s dbname=%s password=%s sslmode=disable"
)

func main() {
	// Init default context
	ctx, cancel := context.WithTimeout(context.Background(), gracefulPeriod)
	defer cancel()

	// Setup configuration
	cfg := config.SetConfiguration()

	// Init DB connection
	db, err := gorm.Open("postgres", fmt.Sprintf(connectionString, cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbName, cfg.DbPassword))
	if err != nil {
		log.Fatalf("could not open db connection: %s", err.Error())
	}
	defer db.Close()
	db.LogMode(cfg.DbLogMode)

	// Init tracing
	jgCfg, err := jaegercfg.FromEnv()
	if err != nil {
		log.Fatal("could not set jaeger config")
	}

	tracer, closer, err := jgCfg.NewTracer()
	if err != nil {
		errorLogger.Log(err, "could not set jaeger")
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// Listen to interruption signal from the system
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	// Init HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		ReadTimeout:  httpReadTimeout,
		WriteTimeout: httpWriteTimeout,
	}

	// Dummy authentication service
	key := []byte(cfg.JwtSigningKey)
	keyFunc := func(*jwt.Token) (interface{}, error) {
		return key, nil
	}
	authSvc := auth.NewService(500, key, keyFunc)

	// Middleware that will check jwt validity
	JWTMiddleware := kitjwt.NewParser(keyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory)

	// payment resource endpoints
	var paymentEndpoints payments.Endpoints
	{
		repository := repository.NewPaymentRepository(db)
		service := payments.NewService(repository)
		paymentEndpoints = payments.MakeEndpoints(service, tracer, JWTMiddleware)
	}

	go func() {
		var mux *http.ServeMux
		{
			mux = http.NewServeMux()

			// TODO V1 == no version, documentation about the v1 alias for fixing the version
			mux.Handle("/payments/", payments.MakePaymentHTTPHandler(
				errorLogger,
				tracer,
				paymentEndpoints))

			// authentication endpoint to receive a JWT
			mux.Handle("/auth/", auth.MakeAuthHandler(authSvc, errorLogger, tracer))
			// For liveness probe
			mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
			// Expose metrics endpoint
			mux.Handle("/metrics", promhttp.Handler())
			// Expose documentation endpoint
			mux.Handle("/swaggerui/", docs.Handler())

			srv.Handler = mux

			log.Printf("payment-api listening on port %s", cfg.AppPort)
			if err := srv.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					log.Fatal(err)
				}
			}
		}
	}()

	// Block here until stop signal received
	<-stopChan

	// Graceful shutdown
	log.Print("shutting down...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Print("gracefully stopped")
}
