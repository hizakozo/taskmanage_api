package mail

import (
	"fmt"
	"net/smtp"
)

func SendMail(mailAddress string, message string) error {
	auth := smtp.PlainAuth(
		"",
		"sample@gmail.com", // foo@gmail.com
		"sample",
		"smtp.gmail.com",
	)

	fmt.Print(message)

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"sample@gmail.com", // 送信元
		[]string{mailAddress}, // 送信先
		[]byte(message),
	)
	return err
}