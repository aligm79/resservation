package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/aligm79/reservation/pkg/models"
	"github.com/aligm79/reservation/pkg/services"
	"github.com/aligm79/reservation/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)


func TicketsList(w http.ResponseWriter, r *http.Request) {
	tickets := services.GetTickets()
	res, _ := json.Marshal(tickets)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func MyTicketsList(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(utils.UserContextKey).(*models.User)
	tickets := services.MyTickets(user.ID)
	result, _ := json.Marshal(tickets)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func GetOrReserveTicket(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(utils.UserContextKey).(*models.User)
	params := mux.Vars(r)
	ticketId, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Bad Id", http.StatusBadRequest)
		return
	}
	switch r.Method{
	case http.MethodGet:	
		ticket, err := services.GetTicket(ticketId)
		if err != nil {
			http.Error(w, "Not found", http.StatusBadRequest)
			return
		}
		result, _ := json.Marshal(ticket)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write(result)
	case http.MethodPost:
		newReservation := models.Reserved{
			ID: 			uuid.New(),
			UserId: 		user.ID,
			TicketId: 		ticketId,
			CreatedDate: 	time.Now(),
		}
		if !services.ReserveTicket(&newReservation) {
			http.Error(w, "an error occured", http.StatusBadRequest)
		}
		result, _ := json.Marshal(newReservation)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	}
}

type LoginRequest struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user , err := services.GetUserForLogin(loginReq.Username, loginReq.Password)
	if err != nil {
		http.Error(w, "user not found", http.StatusForbidden)
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "token could not be generated", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, token)))
}
