package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/aligm79/reservation/pkg/routes"
	"github.com/aligm79/reservation/pkg/tasks"
	"github.com/gorilla/mux"
	"github.com/hibiken/asynq"
)

func startWorker() {
	redisConn := asynq.RedisClientOpt{Addr: "localhost:6379"}
	scheduler := asynq.NewScheduler(redisConn, nil)

	srv := asynq.NewServer(redisConn, asynq.Config{
		Concurrency: 10,
	})

	mux := asynq.NewServeMux()
	_, err := scheduler.Register("*/1 * * * *", asynq.NewTask(tasks.TenMinuteCheck, nil))
	if err != nil {
		log.Fatalf("Could not schedule task: %v", err)
	}

	log.Println("Starting periodic task scheduler...")
	if err := scheduler.Start(); err != nil {
		log.Fatalf("Scheduler failed: %v", err)
	}
	mux.HandleFunc(tasks.TenMinuteCheck, tasks.HandleTenMinuteCheck)

	fmt.Println("Worker started...")

	if err := srv.Run(mux); err != nil {
		log.Fatalf("Worker error: %v", err)
	}
}

func startHTTPServer() {
	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err == nil {
			log.Println("Registered route:", path)
		} else {
			log.Println("Error registering route:", err)
		}
		return nil
	})

	log.Println("HTTP server started on :9010")
	log.Fatal(http.ListenAndServe(":9010", r))
}

func main() {
	//go startWorker()
	startHTTPServer()
}