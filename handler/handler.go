package handler

import (
	"encoding/json"
	"github.com/sillamilla/user_microservice/internal/users/model"
	"github.com/sillamilla/user_microservice/internal/users/service"
	"io/ioutil"
	"net/http"
)

type Handler struct {
	srv service.Service
}

func NewHandler(service service.Service) Handler {
	return Handler{
		srv: service,
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var inputUser model.Input
	err = json.Unmarshal(body, &inputUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.srv.SignUp(r.Context(), inputUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make(map[string]interface{})
	response["message"] = "Sign up successful"
	response["user"] = user

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var inputUser model.Input
	err = json.Unmarshal(body, &inputUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.srv.SignIn(r.Context(), inputUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make(map[string]interface{})
	response["message"] = "Sign in successful"
	response["user"] = user

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) EditProfile(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updateUser model.UpdateUser
	err = json.Unmarshal(body, &updateUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.srv.EditProfile(r.Context(), r.Header.Get("ID"), updateUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make(map[string]string)
	response["message"] = "Edit profile successful"

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var changePassword model.ChangePassword
	err = json.Unmarshal(body, &changePassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.srv.EditPassword(r.Context(), r.Header.Get("ID"), changePassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make(map[string]string)
	response["message"] = "Edit password successful"

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpsertSessions(w http.ResponseWriter, r *http.Request) {
	session := r.Header.Get("ID")

	session, err := h.srv.UpsertSessions(r.Context(), session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Set session successful"}
	response["sessions"] = session

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	err := h.srv.Logout(r.Context(), r.Header.Get("ID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Logout successful"}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetBySession(w http.ResponseWriter, r *http.Request) {
	user, err := h.srv.GetBySession(r.Context(), r.Header.Get("Session"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make(map[string]interface{})
	response["message"] = "Get user by session successful"
	response["user"] = user

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	session, err := h.srv.GetSession(r.Context(), r.Header.Get("ID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Get session successful"}
	response["session"] = session

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SearchByUsername(w http.ResponseWriter, r *http.Request) {
	user, err := h.srv.SearchByUsername(r.Context(), r.Header.Get("Username"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make(map[string]interface{})
	response["message"] = "Get username successful"
	response["user"] = user

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetByUsername(w http.ResponseWriter, r *http.Request) {
	user, err := h.srv.GetByUsername(r.Context(), r.Header.Get("Username"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make(map[string]interface{})
	response["message"] = "Get username successful"
	response["user"] = user

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	user, err := h.srv.GetByID(r.Context(), r.Header.Get("ID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make(map[string]interface{})
	response["message"] = "Get user by id successful"
	response["user"] = user

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
