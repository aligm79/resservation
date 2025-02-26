package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aligm79/reservation/pkg/routes"
	"github.com/gorilla/mux"
)

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