package leads

import "time"

type Lead struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Email     string `json:"email" db:"email"`
	Message   string `json:"message" db:"message"`
	Status    string `json:"status" db:"status"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

func New(name string, email string, message string, id string) Lead {
	return Lead{
		ID:        id,
		Name:      name,
		Email:     email,
		Message:   message,
		Status:    "new",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}
