package unisenderclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/mail"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func New(config Config, log *zerolog.Logger) (*Service, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Service{
		log:    log,
		config: config,
		client: http.Client{Timeout: config.Timeout},
	}, nil
}

// NewMessage
// lang may be  “be”, “de”, “en”, “es”, “fr”, “it”, “pl”, “pt”, “ru”, “ua”.
func (s *Service) NewMessage(subject, body, senderName, lang string, to ...mail.Address) (message Message) {
	message.Subject = subject
	message.Body.HTML = body
	message.SenderName = senderName
	message.AddTo(to...)
	message.SenderEmail = s.config.SenderEmail
	message.GlobalLanguage = lang

	return message
}

func (s *Service) Send(ctx context.Context, message Message) error {
	var reader io.Reader

	reader, err := s.prepareMessage(message)
	if err != nil {
		return errors.Wrap(err, "unisender send")
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, s.config.APIURI, reader)
	if err != nil {
		return errors.Wrap(err, "unisender send")
	}

	r.Header.Set("X-API-KEY", s.config.APIKey)

	resp, err := s.client.Do(r)
	if err != nil {
		return errors.Wrap(err, "unisender send")
	}

	if resp != nil && resp.Body != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}

	if s.log.GetLevel() == zerolog.DebugLevel {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "unisender send")
		}

		s.log.Debug().Str("unisender answer", string(body)).Msg("message was sent")
	}

	return nil
}

func (s *Service) prepareMessage(message Message,
) (io.Reader, error) {
	s.log.Debug().Str("subject", message.Subject).Msg("prepared")
	s.log.Debug().Interface("addresses", message.Recipients).Msg("prepared")
	s.log.Debug().Str("senderEmail", message.SenderEmail).Msg("prepared")
	s.log.Debug().Str("html", message.Body.HTML).Msg("prepared")

	md := MessageDecorator{
		Message: message,
	}

	bt, err := json.Marshal(md)
	if err != nil {
		return nil, errors.Wrap(err, "prepare message")
	}

	return bytes.NewReader(bt), nil
}
