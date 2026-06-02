package email

type Sender interface {
	SendVerificationCode(to, code string) error
}
