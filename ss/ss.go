package ss

type Client interface {
	QueryPortPasswords() map[string]string
	AddPortPassword(port, password string)
	DeletePort(port string) error
	Restart() error
}
