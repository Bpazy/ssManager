package result

const (
	statusOk   = "Ok"
	statusFail = "Fail"
)

type Result struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Ok(msg string, data interface{}) *Result {
	r := Result{}
	r.Code = statusOk
	r.Msg = msg
	r.Data = data
	return &r
}
