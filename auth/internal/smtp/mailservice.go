package smtp

import (
	"crypto/tls"
	"fmt"
	"log"
	"strconv"

	gomail "gopkg.in/mail.v2"
)

type MailService struct {
	from               string
	password           string
	smtpHost           string
	smtpPort           int
	insecureSkipVerify bool
}

func NewMailService(from, password, smtpHost, smtpPort, insecureSkipVerify string) *MailService {
	if from == "" {
		log.Panic(fmt.Errorf("error creating mail service: value of from is empty"))
	}
	if password == "" {
		log.Panic(fmt.Errorf("error creating mail service: value of password is empty"))
	}
	if smtpHost == "" {
		log.Panic(fmt.Errorf("error creating mail service: value of smtpHost is empty"))
	}

	if smtpPort == "" {
		log.Panic(fmt.Errorf("error creating mail service: value of smtpPort is empty"))
	}
	smtpPortInt, err := strconv.Atoi(smtpPort)
	if err != nil {
		log.Panic(fmt.Errorf("error creating mail service: value of smtpPort is not an integer"))
	}

	if insecureSkipVerify == "" {
		log.Panic(fmt.Errorf("error creating mail service: value of insecureSkipVerify is empty"))
	}
	insecureSkipVerifyBool, err := strconv.ParseBool(insecureSkipVerify)
	if err != nil {
		log.Panic(fmt.Errorf("error creating mail service: value of insecureSkipVerify is not a boolean"))
	}
	return &MailService{
		from:               from,
		password:           password,
		smtpHost:           smtpHost,
		smtpPort:           smtpPortInt,
		insecureSkipVerify: insecureSkipVerifyBool,
	}
}

// ref:
// - https://www.loginradius.com/blog/engineering/sending-emails-with-golang/
// - https://ethereal.email/create
func (m *MailService) Send(to, subject, message string) error {
	msg := gomail.NewMessage()

	// Set E-Mail sender
	msg.SetHeader("From", m.from)

	// Set E-Mail receivers
	msg.SetHeader("To", to)

	// Set E-Mail subject
	msg.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	msg.SetBody("text/plain", message)

	// Settings for SMTP server
	d := gomail.NewDialer(m.smtpHost, m.smtpPort, m.from, m.password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: m.insecureSkipVerify}

	// Now send E-Mail
	if err := d.DialAndSend(msg); err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("\nSent mail:")
	fmt.Println("\tFrom:", m.from)
	fmt.Println("\tTo:", to)
	fmt.Println("\tSubject:", subject)
	fmt.Printf("\n\t%v", message)
	fmt.Println("\n ")

	return nil

}
