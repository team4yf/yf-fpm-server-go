//Package errno manager the errors
package errno

const (
	defaultBizCode = -9999
)

var (
	UnDefinedErr = &BizError{Code: -1000, Message: ""}
	//OAuth about
	OAuthOnlySupportClientErr = &BizError{Code: -901, Message: "only support grant_type=client_credentials"}
	OAuthClientAuthErr        = &BizError{Code: -902, Message: "id or secret error!"}
)

type BizError struct {
	Code    int    `json:"errno"`
	Message string `json:"message,omitempty"`
}

func New(code int, message string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
	}
}

func Wrap(err error) *BizError {
	return &BizError{
		Code:    defaultBizCode,
		Message: err.Error(),
	}
}

func Wraps(err string) *BizError {
	return &BizError{
		Code:    defaultBizCode,
		Message: err,
	}
}
