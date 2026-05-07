package api

import (
	"encoding/json"
	"net/http"
)

type authRequest struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Invalid JSON body.", requestID(r))
		return
	}
	token, user, err := s.auth.Register(req.ID, req.Email, req.Password)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), requestID(r))
		return
	}
	writeJSON(w, http.StatusCreated, authResponse{UserID: user.ID, Email: user.Email, Token: token})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Invalid JSON body.", requestID(r))
		return
	}
	token, user, err := s.auth.Login(req.Email, req.Password)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "validation_error", "Invalid email or password.", requestID(r))
		return
	}
	writeJSON(w, http.StatusOK, authResponse{UserID: user.ID, Email: user.Email, Token: token})
}
