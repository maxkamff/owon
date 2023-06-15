package email

import "net/smtp"



func SendEmail(to []string, message []byte) error{
	from := "maxkamff77@gmail.com"
	password := "qljvcfvmozsniiwx"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	return err

}