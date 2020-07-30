package mail

import (
	"fmt"
	"net/smtp"
	"taskmanage_api/src/constants"
)

func SendMail(mailAddress string, message string) error {
	auth := smtp.PlainAuth(
		"",
		constants.Params.MailUserName,
		constants.Params.MailPassword,
		"smtp.gmail.com",
	)

	fmt.Print(message)

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"development@gmail.com", // 送信元
		[]string{mailAddress},   // 送信先
		[]byte(message),
	)
	return err
}
