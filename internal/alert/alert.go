package alert

type Alert interface {
	Send(header, body string) error
}
