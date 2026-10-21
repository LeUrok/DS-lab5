package api

import (
	"encoding/json"
	"net/http"

	"github.com/LeUrok/DS-lab2/gateway-service/internal/auth"
)

type AuthHandler struct {
	client *auth.Client
}

func NewAuthHandler(client *auth.Client) *AuthHandler {
	return &AuthHandler{client: client}
}

type authorizeRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	var req authorizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "username and password are required")
		return
	}

	token, err := h.client.GetToken(r.Context(), req.Username, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"access_token": token.AccessToken,
		"token_type":   token.TokenType,
		"expires_in":   token.ExpiresIn,
	})
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	writeJSON(w, http.StatusOK, map[string]string{"message": "callback ok"})
}