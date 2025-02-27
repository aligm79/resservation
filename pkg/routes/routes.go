package routes

import (
	"net/http"
	"github.com/aligm79/reservation/pkg/controllers"
	"github.com/aligm79/reservation/pkg/utils"
	"github.com/gorilla/mux"
)


var RegisterRoutes = func(router *mux.Router) {
	router.Handle("/tickets", utils.JWTMiddleware(http.HandlerFunc(controllers.TicketsList))).Methods("GET")
	router.Handle("/my_tickets", utils.JWTMiddleware(http.HandlerFunc(controllers.MyTicketsList))).Methods("GET")
	router.Handle("/ticket/{id}/", utils.JWTMiddleware(http.HandlerFunc(controllers.GetOrReserveTicket))).Methods("GET", "POST")
	router.HandleFunc("/login/", controllers.LoginHandler).Methods("POST")
}