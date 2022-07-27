package main

import (
	"context"
	"net/mail"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/cadyrov/unisenderclient"
)

func main() {
	cnf := unisenderclient.Config{
		APIKey:      "65z84unrezrkyc1xocjnvlfivunf",                                  // Fill api key
		APIURI:      "https://eu1.unione.io/ru/transactional/api/v1/email/send.json", // api uri
		SenderEmail: "info@test.com",                                                 // email from account of unisender
		Timeout:     time.Minute,                                                     // base client http timeout
	}

	log := zerolog.New(os.Stdout)
	log.WithLevel(zerolog.DebugLevel)

	serv, err := unisenderclient.New(cnf, &log)
	if err != nil {
		log.Err(err).Msg("create service")
	}

	a, err := mail.ParseAddress("good@recipient.com")
	if err != nil {
		log.Err(err).Msg("parse address")

		return
	}

	msg := serv.NewMessage("info", "<h4>hi</h4>", "goodSender", "en", *a)

	if err = serv.Send(context.Background(), msg); err != nil {
		log.Err(err).Msg("parse address")
	}
}
