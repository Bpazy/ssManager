package ws

type Type string

const (
	Speed Type = "Speed"
	Usage Type = "Usage"
)

type H struct {
	Type Type        `json:"type"`
	Data interface{} `json:"data"`
}
