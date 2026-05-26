package security

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"metavida/backend/core"
	"metavida/backend/db"
)

type loginInput struct {
	Email     string `json:"email"`
	User      string `json:"user"`
	Password  string `json:"password"`
	CipherKey string `json:"cipherKey"`
}

func PostLogin(req *core.HandlerArgs) core.HandlerResponse {
	input := loginInput{}
	if err := core.DecodeJSONBody(req.Body, &input); err != nil {
		return core.HandlerResponse{Error: "invalid JSON body", StatusCode: http.StatusBadRequest}
	}

	email := NormalizeEmail(input.Email)
	if email == "" {
		email = NormalizeEmail(input.User)
	}
	if email == "" || strings.TrimSpace(input.Password) == "" {
		return core.HandlerResponse{Error: "email and password are required", StatusCode: http.StatusBadRequest}
	}

	users := []User{}
	query := db.Query(&users)
	query.Email.Equals(email).Limit(1)
	if err := query.Exec(); err != nil {
		return core.HandlerResponse{Error: "error querying user: " + err.Error(), StatusCode: http.StatusInternalServerError}
	}
	if len(users) == 0 || !VerifyPassword(input.Password, users[0].PasswordHash) {
		return core.HandlerResponse{Error: "invalid email or password", StatusCode: http.StatusUnauthorized}
	}

	user := users[0]
	if user.Status != 0 && user.Status != 1 { // 1 = active
		return core.HandlerResponse{Error: "user is inactive", StatusCode: http.StatusForbidden}
	}

	response, err := MakeUsuarioResponse(user, input.CipherKey)
	if err != nil {
		return core.HandlerResponse{Error: err.Error(), StatusCode: http.StatusInternalServerError}
	}
	return core.HandlerResponse{Body: response}
}

func MakeUsuarioResponse(user User, _ string) (map[string]any, error) {
	usuarioToken := core.UsuarioToken{
		ID:      user.ID,
		Created: core.SUnixTime(),
		User:    user.Email,
		Role:    user.Role,
	}

	userToken, err := core.EncodeUsuarioToken(usuarioToken)
	if err != nil {
		return nil, err
	}

	userInfo := map[string]any{
		"ID":     user.ID,
		"Name":   user.Name,
		"Email":  user.Email,
		"Role":   user.Role,
		"Status": user.Status,
	}
	userInfoBytes, err := json.Marshal(userInfo)
	if err != nil {
		return nil, err
	}

	tokenExpiresAt := time.Now().UTC().Add(4 * time.Hour)
	response := map[string]any{
		"UserID":       user.ID,
		"UserToken":    userToken,
		"TokenExpTime": tokenExpiresAt.Unix(),
		"UserInfo":     core.BytesToBase64(userInfoBytes),
		"user": map[string]any{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"role":   user.Role,
			"status": user.Status,
		},
		"token":            userToken,
		"token_expires_at": tokenExpiresAt.Format(time.RFC3339),
	}
	return response, nil
}
