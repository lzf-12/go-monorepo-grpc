package http

import (
	"encoding/json"
	"net/http"
	"ops-monorepo/services/svc-user/internal/model"
	"ops-monorepo/services/svc-user/internal/usecase"
	"ops-monorepo/shared-libs/logger"

	"github.com/gorilla/mux"
)

type UserHTTPHandler struct {
	usecase usecase.IUserUsecase
	log     logger.Logger
}

func NewUserHTTPHandler(usecase usecase.IUserUsecase, log logger.Logger) *UserHTTPHandler {
	return &UserHTTPHandler{
		usecase: usecase,
		log:     log,
	}
}

func (h *UserHTTPHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/auth/register", h.Register).Methods("POST")
	router.HandleFunc("/api/v1/auth/login", h.Login).Methods("POST")
	router.HandleFunc("/api/v1/auth/refresh", h.RefreshToken).Methods("POST")
}

func (h *UserHTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorf("failed to decode request: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.usecase.Register(r.Context(), &req)
	if err != nil {
		h.log.Errorf("failed to register user: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "user registered successfully",
		"data":    user,
	})
}

func (h *UserHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorf("failed to decode request: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	loginResponse, err := h.usecase.Login(r.Context(), &req)
	if err != nil {
		h.log.Errorf("failed to login user: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "login successful",
		"data":    loginResponse,
	})
}

func (h *UserHTTPHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorf("failed to decode request: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.RefreshToken == "" {
		http.Error(w, "refresh token is required", http.StatusBadRequest)
		return
	}

	refreshResponse, err := h.usecase.RefreshToken(r.Context(), &req)
	if err != nil {
		h.log.Errorf("failed to refresh token: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "token refreshed successfully",
		"data":    refreshResponse,
	})
}