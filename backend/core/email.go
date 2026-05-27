package core

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type EmailConfig struct {
	AccessKey string
	SecretKey string
	Region    string
	From      string
	To        string
}

type EmailMessage struct {
	Subject string
	Body    string
	ReplyTo string
}

func SendEmail(emailConfig EmailConfig, message EmailMessage) error {
	emailConfig.Region = strings.TrimSpace(emailConfig.Region)
	if emailConfig.Region == "" {
		emailConfig.Region = "us-east-1"
	}
	if strings.TrimSpace(emailConfig.AccessKey) == "" || strings.TrimSpace(emailConfig.SecretKey) == "" {
		return fmt.Errorf("AWS SES credentials are missing")
	}
	if strings.TrimSpace(emailConfig.From) == "" || strings.TrimSpace(emailConfig.To) == "" {
		return fmt.Errorf("AWS SES sender and recipient are required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Static credentials come from credentials.json so SES can run outside an AWS-hosted environment.
	awsConfig := aws.Config{
		Region: emailConfig.Region,
		Credentials: aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(emailConfig.AccessKey, emailConfig.SecretKey, ""),
		),
	}

	input := &ses.SendEmailInput{
		Source: aws.String(emailConfig.From),
		Destination: &types.Destination{
			ToAddresses: []string{emailConfig.To},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(message.Subject),
			},
			Body: &types.Body{
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(message.Body),
				},
			},
		},
	}
	if strings.TrimSpace(message.ReplyTo) != "" {
		input.ReplyToAddresses = []string{strings.TrimSpace(message.ReplyTo)}
	}

	log.Printf("email: sending SES message to %s in %s", emailConfig.To, emailConfig.Region)
	output, err := ses.NewFromConfig(awsConfig).SendEmail(ctx, input)
	if err != nil {
		return err
	}
	log.Printf("email: SES accepted message %s", aws.ToString(output.MessageId))
	return nil
}
