package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"
	"trigger-queue/sensors"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var (
	// Flags
	walkFlg    bool // Walk and print the routes - default: false
	verboseFlg bool // Verbose logging - default: false

	// Sensors
	Scale       *sensors.Scale
	Thermometer *sensors.Thermometer
	Stove       *sensors.Stove
)

func init() {
	// Handle cli flags
	flag.BoolVar(&walkFlg, "walk", false, "Walk the routes")
	flag.BoolVar(&verboseFlg, "verbose", false, "Display verbose logs")

	flag.Parse()

	// Load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Setup logger
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)

	// Print env values
	if verboseFlg {
		var env map[string]string
		env, err = godotenv.Read()

		if err != nil {
			log.Warn("Error printing .env files")
		} else {
			log.Info("Environment Variables")
			for k, v := range env {
				log.Printf("%s=%s\n", k, v)
			}
		}
	}

	// Setup sensors
	Scale, err = new(sensors.Scale).GetInstance()
	if err != nil {
		log.Error("Error setting up Scale")
		log.Fatal(err.Error())
	}

	Thermometer, err = new(sensors.Thermometer).GetInstance()
	if err != nil {
		log.Error("Error setting up Thermometer")
		log.Fatal(err.Error())
	}

	Stove, err = new(sensors.Stove).GetInstance()
	if err != nil {
		log.Error("Error setting up Stove")
		log.Fatal(err.Error())
	}
}

func main() {
	// Defer cleaning up the Stove
	defer Stove.Cleanup()

	// Setup Mux
	r := mux.NewRouter()
	var wait time.Duration

	// Log endpoints Method + URI
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// https://golang.org/pkg/net/http/#Request
			log.Info(r.Method + " " + r.RequestURI)
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		})
	})

	// Route handles & endpoints
	CreateRoutes(r, AllRoutes[:], "/api")

	// Walk all the routes
	if walkFlg {
		r.Walk(RouteWalker)
	}

	// Start server
	srv := &http.Server{
		Handler: r,
		Addr:    os.Getenv("TRIGGER_QUEUE_ADDR"),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("Starting trigger-queue server at %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Run the trigger queue consumer thread
	go func() {
		if verboseFlg {
			log.Info("Started trigger queue consumer thread")
		}
		for {
			for _, q := range TriggerQueue {
				triggersTrue, elem := q.EvaluateFront(evaluateTriggers, triggerError)
				if triggersTrue {
					executeAction(*elem)
				}
			}
			time.Sleep(2 * time.Second)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	os.Exit(0)
}
