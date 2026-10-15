package auth

import (
	"encoding/json"
	"errors"
	"jwt-auth/internal/user"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

type registerRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	u := &user.User{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	createdUser, err := h.service.Register(r.Context(), u)
	if err != nil {
		if errors.Is(err, user.ErrAllFieldsRequired) {
			http.Error(w, "all fields required", http.StatusBadRequest)
			return
		}

		if errors.Is(err, user.ErrAlreadyExist) {
			http.Error(w, "email already exist", http.StatusConflict)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	u := &user.User{
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	token, loggedUser, err := h.service.Login(r.Context(), u)
	if err != nil {
		if errors.Is(err, user.ErrAllFieldsRequired) {
			http.Error(w, "all fields required", http.StatusBadRequest)
			return
		}

		if errors.Is(err, user.ErrInvalidCredential) {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{
		"data":  loggedUser,
		"token": token,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
