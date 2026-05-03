package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"project/internal/service"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	svc    *service.UserService
	jwtKey string
}

func NewUserHandler(svc *service.UserService, jwtKey string) *UserHandler {
	return &UserHandler{svc: svc, jwtKey: jwtKey}
}

func (h *UserHandler) Registr(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Некорректный json", http.StatusBadRequest)
		return
	}
	user, err := h.svc.Registration(req.Name, req.Email, req.Password)
	if err != nil {
		http.Error(w, "Ошибка валидации", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("ошибка при отправке ответа: %v", err)
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	user, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.Id,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})
	TokenString, err := token.SignedString([]byte(h.jwtKey))
	if err != nil {
		http.Error(w, "ошибка создания токена", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": TokenString,
		"user":  user,
	})
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}

	user, err := h.svc.Create(r.Context(), req.Name, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userId").(int)
	id, _ := strconv.Atoi(r.PathValue("id"))

	if userID != id {
		http.Error(w, "доступ запрещён", http.StatusForbidden)
		return
	}

	user, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "не найден", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	id, _ := strconv.Atoi(r.PathValue("id"))
	if userID != id {
		http.Error(w, "доступ запрещён", http.StatusForbidden)
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}

	user, err := h.svc.Update(r.Context(), id, req.Name, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(r.PathValue("id"))
	userID := r.Context().Value("userID").(int)
	if userID != id {
		http.Error(w, "доступ запрещён", http.StatusForbidden)
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, "не найден", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
