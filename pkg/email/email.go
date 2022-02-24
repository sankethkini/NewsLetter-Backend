package email

import (
	"context"
	"fmt"
	"net/smtp"

	newsletter "github.com/sankethkini/NewsLetter-Backend/internal/service/news_letter"
	subscriptionpb "github.com/sankethkini/NewsLetter-Backend/proto/subscriptionpb/v1"
	userpb "github.com/sankethkini/NewsLetter-Backend/proto/userpb/v1"
	"google.golang.org/grpc"
)

// nolint:revive
type EmailConfig struct {
	From     string `yaml:"from"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	GrpcPort string `yaml:"server"`
}

type Email struct {
	email EmailConfig
}

func NewEmailServer(email EmailConfig) *Email {
	return &Email{email: email}
}

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

func NewMail(to []string, subject string, body string) Mail {
	return Mail{
		To:      to,
		Subject: subject,
		Body:    body,
	}
}

// get emails of the users to whom newsletter have to be send.
// nolint: govet,staticcheck
func (em Email) GetEmails(ctx context.Context, n newsletter.EmailData) ([]string, error) {
	serverAddress := em.email.GrpcPort
	sid := n.Scheme.SchemeId
	conn, e := grpc.Dial(serverAddress, grpc.WithInsecure())

	if e != nil {
		return nil, e
	}
	defer conn.Close()
	cl := subscriptionpb.NewSubscriptionServiceClient(conn)
	data := subscriptionpb.GetUsersRequest{
		SchemeId: sid,
	}

	resp, err := cl.GetUsers(ctx, &data)
	if err != nil {
		return nil, err
	}
	emails := make([]string, 0, len(resp.UserIds))

	cl1 := userpb.NewUserServiceClient(conn)
	for _, val := range resp.UserIds {
		data1 := userpb.GetEmailRequest{
			Name: val,
		}
		res, err := cl1.GetEmail(ctx, &data1)
		if err != nil {
			return nil, err
		}
		emails = append(emails, res.Email)
	}
	return emails, nil
}

// send emails to respective emails.
func (em Email) SendEmail(m Mail) error {
	from := em.email.From
	password := em.email.Password

	host := em.email.Host
	port := em.email.Port

	address := host + ":" + port

	m.Sender = from

	message := BuildMessage(m)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, m.To, []byte(message))
	if err != nil {
		return err
	}
	return nil
}

// build message, add title and body.
func BuildMessage(mail Mail) string {
	msg := ""
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)

	if len(mail.To) > 0 {
		msg += fmt.Sprintf("To: %s\r\n", mail.To[0])
	}

	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}
