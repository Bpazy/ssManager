package result

const (
	statusOk   = "Ok"
	statusFail = "Failed"
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

func Fail(msg string, data interface{}) *Result {
	r := Result{}
	r.Code = statusFail
	r.Msg = msg
	r.Data = data
	return &r
}
