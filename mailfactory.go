package kiam

type MailFactory interface {
	MakeEmail(key string, param map[string]interface{}) MailMessage
}

type MailMessage struct {
	Subject string
	Body    string
}
