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

func init() {
	redisConn := asynq.RedisClientOpt{Addr: "localhost:6379"}

    srv := asynq.NewServer(redisConn, asynq.Config{
        Concurrency: 10, 
    })

    mux := asynq.NewServeMux()
    mux.HandleFunc(tasks.AdderTask, tasks.HandleLogUserTask)
	
    fmt.Println("Worker started...")
    if err := srv.Run(mux); err != nil {
        log.Fatal(err)
    }
}

func main() {
	r := mux.NewRouter()
	routes.RegisterRoutes(r)
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
        path, err := route.GetPathTemplate()
        if err == nil {
            fmt.Println("Registered route:", path)
        }else {
			fmt.Print(err)
		}
        return nil
    	})

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":9010", r))
}