package routes

import (
	//	"github.com/aligm79/reservation/pkg/config"
	"net/http"

	"github.com/aligm79/reservation/pkg/controllers"
	"github.com/aligm79/reservation/pkg/utils"
	"github.com/gorilla/mux"
	// "gorm.io/gorm"
)

//var db *gorm.DB = config.GetDB()

var RegisterRoutes = func(router *mux.Router) {
	router.Handle("/tickets", utils.JWTMiddleware(http.HandlerFunc(controllers.TicketsList))).Methods("GET")
	router.Handle("/ticket/{id}/", utils.JWTMiddleware(http.HandlerFunc(controllers.GetTicket))).Methods("GET", "POST")
	router.HandleFunc("/login", controllers.LoginHandler).Methods("POST")
}