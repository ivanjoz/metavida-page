package operations

import (
	"net/http"
	"strings"

	"metavida/backend/core"
	"metavida/backend/db"
	"metavida/backend/security"
)

type clientInput struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

func PostClients(args *core.HandlerArgs) core.HandlerResponse {
	var input clientInput
	if err := core.DecodeJSONBody(args.Body, &input); err != nil {
		return core.HandlerResponse{Error: "invalid JSON body", StatusCode: http.StatusBadRequest}
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Message = strings.TrimSpace(input.Message)
	if input.Name == "" {
		return core.HandlerResponse{Error: "name is required", StatusCode: http.StatusBadRequest}
	}

	now := core.SUnixTime()
	user := security.User{
		ID:      core.SUnixTime(),
		Name:    input.Name,
		Email:   input.Email,
		Phone:   input.Phone,
		Message: input.Message,
		Type:    1, // 1 = client
		Status:  1,
		Created: now,
		Updated: now,
	}
	if err := db.InsertOne(user); err != nil {
		return core.HandlerResponse{Error: err.Error(), StatusCode: http.StatusInternalServerError}
	}
	return core.HandlerResponse{Body: user, StatusCode: http.StatusCreated}
}

func GetClients(_ *core.HandlerArgs) core.HandlerResponse {
	rows := []security.User{}
	query := db.Query(&rows)
	query.Type.Equals(int8(1)).Limit(50)
	if err := query.Exec(); err != nil {
		return core.HandlerResponse{Error: err.Error(), StatusCode: http.StatusInternalServerError}
	}
	return core.HandlerResponse{Body: rows}
}
