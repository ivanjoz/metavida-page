package operations

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"metavida/backend/config"
	"metavida/backend/core"
	"metavida/backend/db"
)

type UserSendedEmails struct {
	db.TableStruct[UserSendedEmailsTable, UserSendedEmails]
	ID         int32  `json:",omitempty"`
	Name       string `json:",omitempty"`
	Email      string `json:",omitempty"`
	Phone      string `json:",omitempty"`
	Subject    string `json:",omitempty"`
	Message    string `json:",omitempty"`
	IP         string `json:",omitempty"`
	Metadata   string `json:",omitempty"`
	EmailError string `json:",omitempty"`
	Status     int8   `json:"ss,omitempty"` // 1 = saved, 2 = sent, 3 = send failed
	Created    int32  `json:",omitempty"`
	Updated    int32  `json:"upd,omitempty"`
}

type UserSendedEmailsTable struct {
	db.TableStruct[UserSendedEmailsTable, UserSendedEmails]
	ID         db.Col[UserSendedEmailsTable, int32]
	Name       db.Col[UserSendedEmailsTable, string]
	Email      db.Col[UserSendedEmailsTable, string]
	Phone      db.Col[UserSendedEmailsTable, string]
	Subject    db.Col[UserSendedEmailsTable, string]
	Message    db.Col[UserSendedEmailsTable, string]
	IP         db.Col[UserSendedEmailsTable, string]
	Metadata   db.Col[UserSendedEmailsTable, string]
	EmailError db.Col[UserSendedEmailsTable, string]
	Status     db.Col[UserSendedEmailsTable, int8]
	Created    db.Col[UserSendedEmailsTable, int32]
	Updated    db.Col[UserSendedEmailsTable, int32]
}

func (UserSendedEmailsTable) GetSchema() db.TableSchema {
	table := db.Table[UserSendedEmails]()
	return db.TableSchema{
		Name: "user_sended_emails",
		Keys: []db.Coln{table.ID},
		Indexes: []db.Index{
			{Keys: []db.Coln{table.Email}},
			{Keys: []db.Coln{table.Status}},
			{Keys: []db.Coln{table.Created}},
			{Keys: []db.Coln{table.IP, table.Created}},
		},
	}
}

type contactEmailInput struct {
	Name    string `json:",omitempty"`
	Email   string `json:",omitempty"`
	Phone   string `json:",omitempty"`
	Subject string `json:",omitempty"`
	Message string `json:",omitempty"`
}

func PostContactEmail(args *core.HandlerArgs) core.HandlerResponse {
	var input contactEmailInput
	if err := core.DecodeJSONBody(args.Body, &input); err != nil {
		log.Printf("contact-email: invalid JSON body: %v", err)
		return core.HandlerResponse{Error: "invalid JSON body", StatusCode: http.StatusBadRequest}
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Subject = strings.TrimSpace(input.Subject)
	input.Message = strings.TrimSpace(input.Message)
	if input.Subject == "" {
		input.Subject = "Contacto desde MetaVida"
	}
	if input.Name == "" || input.Email == "" || input.Phone == "" || input.Message == "" {
		return core.HandlerResponse{Error: "name, email, phone and message are required", StatusCode: http.StatusBadRequest}
	}
	parsedEmail, err := mail.ParseAddress(input.Email)
	if err != nil {
		return core.HandlerResponse{Error: "email is invalid", StatusCode: http.StatusBadRequest}
	}
	input.Email = parsedEmail.Address

	requestIP := contactRequestIP(args)
	if waitSeconds, err := contactEmailCooldownSeconds(requestIP); err != nil {
		log.Printf("contact-email: cooldown lookup failed for ip %s: %v", requestIP, err)
		return core.HandlerResponse{Error: err.Error(), StatusCode: http.StatusInternalServerError}
	} else if waitSeconds > 0 {
		log.Printf("contact-email: rejected cooldown for ip %s, wait %d seconds", requestIP, waitSeconds)
		return core.HandlerResponse{Error: fmt.Sprintf("Debes esperar %d segundos antes de enviar otro mensaje.", waitSeconds), StatusCode: http.StatusTooManyRequests}
	}

	now := core.SUnixTime()
	record := UserSendedEmails{
		ID:      contactEmailID(),
		Name:    input.Name,
		Email:   input.Email,
		Phone:   input.Phone,
		Subject: input.Subject,
		Message: input.Message,
		IP:      requestIP,
		Status:  1,
		Created: now,
		Updated: now,
	}
	log.Printf("contact-email: saving contact request from %s", record.Email)
	if err := db.InsertOne(record); err != nil {
		log.Printf("contact-email: save failed: %v", err)
		return core.HandlerResponse{Error: err.Error(), StatusCode: http.StatusInternalServerError}
	}

	emailConfig, err := loadContactEmailConfig()
	if err != nil {
		markContactEmailFailed(record, err)
		log.Printf("contact-email: config failed after save: %v", err)
		return core.HandlerResponse{Error: err.Error(), StatusCode: http.StatusInternalServerError}
	}

	err = core.SendEmail(emailConfig, core.EmailMessage{
		Subject: record.Subject,
		ReplyTo: record.Email,
		Body: strings.Join([]string{
			"Nuevo contacto desde MetaVida",
			"",
			"Nombre: " + record.Name,
			"Email: " + record.Email,
			"Telefono: " + record.Phone,
			"Asunto: " + record.Subject,
			"IP: " + record.IP,
			"",
			record.Message,
		}, "\n"),
	})
	if err != nil {
		markContactEmailFailed(record, err)
		log.Printf("contact-email: send failed after save: %v", err)
		return core.HandlerResponse{Error: "contact request was saved but email could not be sent", StatusCode: http.StatusInternalServerError}
	}

	record.Status = 2
	record.Updated = core.SUnixTime()
	if err := db.UpdateOne(record, db.Table[UserSendedEmails]().Status, db.Table[UserSendedEmails]().Updated); err != nil {
		log.Printf("contact-email: sent email but status update failed: %v", err)
	}
	log.Printf("contact-email: saved and sent contact request %d", record.ID)
	return core.HandlerResponse{Body: map[string]string{"status": "sent"}, StatusCode: http.StatusCreated}
}

func contactRequestIP(args *core.HandlerArgs) string {
	// Prefer proxy headers because production requests may arrive behind Cloudflare or a reverse proxy.
	for _, headerName := range []string{"CF-Connecting-IP", "X-Real-IP", "X-Forwarded-For"} {
		value := strings.TrimSpace(args.Headers[headerName])
		if value == "" {
			value = strings.TrimSpace(args.Headers[strings.ToLower(headerName)])
		}
		if value == "" {
			continue
		}
		firstValue := strings.TrimSpace(strings.Split(value, ",")[0])
		if parsedIP := net.ParseIP(firstValue); parsedIP != nil {
			return parsedIP.String()
		}
	}
	if args.Request == nil {
		return ""
	}
	host, _, err := net.SplitHostPort(args.Request.RemoteAddr)
	if err != nil {
		host = args.Request.RemoteAddr
	}
	if parsedIP := net.ParseIP(strings.TrimSpace(host)); parsedIP != nil {
		return parsedIP.String()
	}
	return strings.TrimSpace(host)
}

func contactEmailCooldownSeconds(ip string) (int, error) {
	if strings.TrimSpace(ip) == "" {
		return 0, nil
	}
	now := core.SUnixTime()
	cooldownStart := now - 15 // SUnixTime advances once every two real seconds; 15 units = 30 seconds.
	rows := []UserSendedEmails{}
	query := db.Query(&rows)
	query.IP.Equals(ip).Created.GreaterEqual(cooldownStart).Limit(20)
	if err := query.Exec(); err != nil {
		return 0, err
	}
	newestCreated := int32(0)
	for _, row := range rows {
		if row.Created > newestCreated {
			newestCreated = row.Created
		}
	}
	if newestCreated == 0 {
		return 0, nil
	}
	elapsedSeconds := int((now - newestCreated) * 2)
	if elapsedSeconds >= 30 {
		return 0, nil
	}
	return 30 - elapsedSeconds, nil
}

func contactEmailID() int32 {
	// Keep the primary key compact for D1 while avoiding SUnixTime collisions on repeated submissions.
	return int32(time.Now().UnixNano() % 2147483647)
}

func markContactEmailFailed(record UserSendedEmails, err error) {
	record.Status = 3
	record.EmailError = err.Error()
	record.Updated = core.SUnixTime()
	table := db.Table[UserSendedEmails]()
	if updateErr := db.UpdateOne(record, table.Status, table.EmailError, table.Updated); updateErr != nil {
		log.Printf("contact-email: failed status update failed: %v", updateErr)
	}
}

func loadContactEmailConfig() (core.EmailConfig, error) {
	cfg, err := config.Load()
	if err != nil {
		return core.EmailConfig{}, err
	}
	sender := strings.TrimSpace(cfg.SenderEmail)
	if sender == "" {
		sender = strings.TrimSpace(cfg.AWSSESSender)
	}
	if sender == "" {
		sender = "informes@metavida.life"
	}
	recipient := strings.TrimSpace(cfg.AWSSESRecipient)
	if recipient == "" {
		recipient = strings.TrimSpace(cfg.AdminEmail)
	}
	if recipient == "" {
		recipient = sender
	}
	return core.EmailConfig{
		AccessKey: cfg.AWSUserKey,
		SecretKey: cfg.AWSSecretKey,
		Region:    cfg.AWSRegion,
		From:      sender,
		To:        recipient,
	}, nil
}
