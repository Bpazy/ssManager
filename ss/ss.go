package ss

type Client interface {
	QueryPorts() []string
	AddPortPassword(port, password string)
	DeletePort(port string) error
	Restart() error
}
