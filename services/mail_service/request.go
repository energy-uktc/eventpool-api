package mail_service

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/energy-uktc/eventpool-api/config"
)

type Request struct {
	to      []string
	subject string
	body    string
}

func VerificationCodeRequest(name string, email string, verificationCode string, mobileAppUrl string) {
	request := &Request{
		to:      []string{email},
		subject: "Eventpool: Verify your email",
	}
	valuesMap := make(map[string]string)
	valuesMap["name"] = name
	valuesMap["url"] = fmt.Sprintf("%s/web/auth/verify?code=%s&mobileLink=%s", config.Properties.Hostname, verificationCode, mobileAppUrl)

	if err := request.parseTemplate("./templates/mail/verify_code.html", valuesMap); err == nil {
		mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		subject := "Subject: " + request.subject + "\n"
		message := subject + mime + "\n" + request.body
		send(request.to, message)
	} else {
		log.Println(err)
	}
}

func ResetPasswordRequest(name string, email string, verificationCode string, mobileAppUrl string) {
	request := &Request{
		to:      []string{email},
		subject: "Eventpool: Password reset",
	}

	valuesMap := make(map[string]string)
	valuesMap["name"] = name
	valuesMap["url"] = fmt.Sprintf("%s/web/auth/resetPassword?code=%s&mobileLink=%s", config.Properties.Hostname, verificationCode, mobileAppUrl)

	if err := request.parseTemplate("./templates/mail/password_reset.html", valuesMap); err == nil {
		mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		subject := "Subject: " + request.subject + "\n"
		message := subject + mime + "\n" + request.body
		send(request.to, message)
	} else {
		log.Println(err)
	}

}

func (r *Request) parseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
