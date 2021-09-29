package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gitlab.linoproject.lab/myteam/myapp/mycomponent/server"
)

var endpointsAccessed = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "mycomponent_endpoints_accessed_total",
		Help: "Total number of accessed to a given endpoint",
	},
	[]string{"accessed_endpoint"},
)

// Root returns the release configured by the user
func Root(w http.ResponseWriter, r *http.Request) {
	version := getEnv("VERSION", "NotSet")
	env := getEnv("ENV", "NotSet")
	output := fmt.Sprintf("Hello VMworld! Version: %s - Environment: %s\n", version, env)
	_, err := w.Write([]byte(output))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	endpointsAccessed.WithLabelValues("root").Inc()
}

//Health returns healthy string, can be used for monitoring pourposes
func Health(w http.ResponseWriter, r *http.Request) {
	health := "Healthy"
	_, err := w.Write([]byte(health + "\n"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	endpointsAccessed.WithLabelValues("health").Inc()
}

//getEnv returns the value for a given Env Var
func getEnv(varName, defaultValue string) string {
	if varValue, ok := os.LookupEnv(varName); ok {
		return varValue
	}
	return defaultValue
}

func main() {
	version := getEnv("VERSION", "NotSet")
	log.Println("🧩  Version", version)

	prometheus.MustRegister(endpointsAccessed)

	router := mux.NewRouter()
	router.Use(func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, next)
	})
	router.HandleFunc("/", Root).Methods("GET")
	router.HandleFunc("/health", Health).Methods("GET")
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	srv, err := server.Listen(router, "")
	if err != nil {
		panic(err)
	}
	log.Printf("👋  Listening on %v", srv.Listener().Addr())

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	log.Printf("☠️  Shutting down")
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Stop(c)
}
