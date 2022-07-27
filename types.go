package unisenderclient

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type AttachmentType string

const (
	Octet AttachmentType = "application/octet-stream"
)

type Config struct {
	APIKey      string
	APIURI      string
	SenderEmail string
	Timeout     time.Duration
}

var (
	ErrAPIKeyEmpty        = errors.New("api key is empty")
	ErrAPIUrlEmpty        = errors.New("api url is empty")
	ErrSenderEmailEmpty   = errors.New("sender email is empty")
	ErrSenderEmailInvalid = errors.New("sender email is invalid")
)

func (c Config) Validate() error {
	if strings.Trim(c.APIKey, " ") == "" {
		return ErrAPIKeyEmpty
	}

	if strings.Trim(c.APIURI, " ") == "" {
		return ErrAPIUrlEmpty
	}

	if strings.Trim(c.SenderEmail, " ") == "" {
		return ErrSenderEmailEmpty
	}

	if _, err := mail.ParseAddress(c.SenderEmail); err != nil {
		return ErrSenderEmailInvalid
	}

	return nil
}

type Service struct {
	log    *zerolog.Logger
	config Config
	client http.Client
}

type MessageDecorator struct {
	Message Message `json:"message"`
}

type To struct {
	Email string `json:"email"`
}

type Recipient []To

type Attachments []Attachment

type Attachment struct {
	Type    AttachmentType `json:"type"`
	Name    string         `json:"name"`
	Content string         `json:"content"`
}

type Message struct {
	Recipients     Recipient   `json:"recipients"`
	Body           Body        `json:"body"`
	SenderEmail    string      `json:"from_email"`
	SenderName     string      `json:"from_name"`
	Subject        string      `json:"subject"`
	GlobalLanguage string      `json:"global_language"`
	Attachments    Attachments `json:"attachments"`
}

type Body struct {
	HTML string `json:"html"`
}

func (m *Message) AddTo(to ...mail.Address) {
	for i := range to {
		m.Recipients = append(m.Recipients, To{Email: to[i].Address})
	}
}

func (m *Message) AddFile(name string, content []byte) {
	m.Attachments = append(m.Attachments,
		Attachment{
			Type:    Octet,
			Name:    name,
			Content: base64.StdEncoding.EncodeToString(content),
		})
}
