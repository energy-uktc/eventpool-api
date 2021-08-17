package mail_service

import (
	"log"
	"net/smtp"

	"github.com/energy-uktc/eventpool-api/config"
)

var (
	from     string
	password string
	smtpHost string
	smtpPort string
)

func init() {
	from = readSmtpValueFromConfig("from").(string)
	password = readSmtpValueFromConfig("password").(string)
	smtpHost = readSmtpValueFromConfig("host").(string)
	smtpPort = readSmtpValueFromConfig("port").(string)
}

func readSmtpValueFromConfig(key string) interface{} {
	value := config.GetPath("smtp." + key)
	if value == nil {
		log.Fatalf("SMTP %s is not configured", key)
	}
	return value
}

func send(receivers []string, message string) {
	messageBytes := []byte(message)
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, receivers, messageBytes)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Email Sent Successfully!")

}
