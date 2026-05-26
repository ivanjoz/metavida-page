package security

import (
	"strings"

	"metavida/backend/config"
	"metavida/backend/core"
	"metavida/backend/db"
)

type User struct {
	db.TableStruct[UserTable, User]
	ID           int32  `json:",omitempty"`
	Name         string `json:",omitempty"`
	Username     string `json:",omitempty"`
	Email        string `json:",omitempty"`
	Phone        string `json:",omitempty"`
	Message      string `json:",omitempty"`
	PasswordHash string `json:"-"`
	Role         string `json:",omitempty"`
	Type         int8   `json:",omitempty"` // 1 = client, 2 = admin
	Status       int8   `json:"ss,omitempty"`
	Created      int32  `json:",omitempty"`
	Updated      int32  `json:"upd,omitempty"`
}

type UserTable struct {
	db.TableStruct[UserTable, User]
	ID           db.Col[UserTable, int32]
	Name         db.Col[UserTable, string]
	Username     db.Col[UserTable, string]
	Email        db.Col[UserTable, string]
	Phone        db.Col[UserTable, string]
	Message      db.Col[UserTable, string]
	PasswordHash db.Col[UserTable, string]
	Role         db.Col[UserTable, string]
	Type         db.Col[UserTable, int8]
	Status       db.Col[UserTable, int8]
	Created      db.Col[UserTable, int32]
	Updated      db.Col[UserTable, int32]
}

func (UserTable) GetSchema() db.TableSchema {
	table := db.Table[User]()
	return db.TableSchema{
		Name: "users",
		Keys: []db.Coln{table.ID},
		Indexes: []db.Index{
			{Keys: []db.Coln{table.Email}},
			{Keys: []db.Coln{table.Type}},
			{Keys: []db.Coln{table.Status}},
		},
	}
}

func EnsureBootstrapAdmin(cfg config.Config) error {
	email := NormalizeEmail(cfg.AdminEmail)
	if email == "" || strings.TrimSpace(cfg.AdminPassword) == "" {
		return nil
	}

	users := []User{}
	query := db.Query(&users)
	query.Email.Equals(email).Limit(1)
	if err := query.Exec(); err != nil {
		return err
	}
	if len(users) > 0 {
		return nil
	}

	passwordHash, err := HashPassword(cfg.AdminPassword)
	if err != nil {
		return err
	}

	now := core.SUnixTime()
	user := User{
		ID:           1,
		Name:         "admin",
		Username:     "admin",
		Email:        email,
		PasswordHash: passwordHash,
		Role:         "admin",
		Type:         2, // 2 = admin
		Status:       1,
		Created:      now,
		Updated:      now,
	}
	return db.InsertOne(user)
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ValidateTokenUser(userToken *core.UsuarioToken) error {
	users := []User{}
	query := db.Query(&users)
	query.ID.Equals(userToken.ID).Limit(1)
	if err := query.Exec(); err != nil {
		return err
	}
	if len(users) == 0 {
		return core.Err("No se encontró el user.")
	}
	user := users[0]
	if user.Status != 0 && user.Status != 1 {
		return core.Err("El user está inactivo.")
	}
	if user.Email != userToken.User || user.Role != userToken.Role {
		return core.Err("El token de inicio de sesión no es válido.")
	}
	return nil
}
