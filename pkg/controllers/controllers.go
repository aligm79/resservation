package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aligm79/reservation/pkg/models"
	"github.com/aligm79/reservation/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)


func TicketsList(w http.ResponseWriter, r *http.Request) {
	tickets := models.GetTickets()
	res, _ := json.Marshal(tickets)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
type ContextKey string
const UserContextKey ContextKey = "user"

func GetTicket(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserContextKey).(*models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}
	fmt.Print(user, "this is the user")
	switch r.Method{
	case http.MethodGet:
		params := mux.Vars(r)
		ticketId, err := uuid.Parse(params["id"])
		if err != nil {
			http.Error(w, "Bad Id", http.StatusBadRequest)
			return
		}
		ticket, err := models.GetTicket(ticketId)
		if err != nil {
			http.Error(w, "Not found", http.StatusBadRequest)
			return
		}
		result, _ := json.Marshal(ticket)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write(result)
	case http.MethodPost:
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

	user , err := models.GetUserForLogin(loginReq.Username, loginReq.Password)
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
