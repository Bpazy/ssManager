package base64

import "encoding/base64"

const (
	base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
)

var coder = base64.NewEncoding(base64Table)

func Encode(src []byte) string {
	return coder.EncodeToString(src)
}

func Decode(s string) ([]byte, error) {
	src := []byte(s)
	dst := make([]byte, len(src))
	_, err := coder.Decode(dst, src)
	if err != nil {
		return nil, err
	}
	return dst, nil
}
