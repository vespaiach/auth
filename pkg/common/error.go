package common

type ErrCode int

const (
	ErrGetData ErrCode = iota
	ErrExecData
	ErrDataNotFound
	ErrDataFailValidation
)

type AppErr struct {
	err     error
	Code    ErrCode
	Payload map[string]string
}

func NewAppErr(err error, code ErrCode) *AppErr {
	return &AppErr{err: err, Code: code}
}

func (ae *AppErr) Error() string {
	return ae.err.Error()
}

func (ae *AppErr) AddMsg(key string, msg string) {
	ae.Payload[key] = msg
}
