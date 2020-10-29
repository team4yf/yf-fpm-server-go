//Package errno manager the errors
package errno

const (
	defaultBizCode = -9999
)

var (
	//UnDefinedErr undefined error
	UnDefinedErr = &BizError{Code: -1000, Message: ""}
	//OAuthOnlySupportClientErr about
	OAuthOnlySupportClientErr = &BizError{Code: -901, Message: "only support grant_type=client_credentials"}
	//OAuthClientAuthErr err
	OAuthClientAuthErr = &BizError{Code: -902, Message: "id or secret error!"}
)

// BizError error struct
type BizError struct {
	Code    int    `json:"errno"`
	Message string `json:"message,omitempty"`
}

func (err BizError) Error() string {
	return err.Message
}

//New create a new bizError
func New(code int, message string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
	}
}

//Wrap wrap a error
func Wrap(err error) *BizError {
	switch err.(type) {
	case *BizError:
		return err.(*BizError)
	}
	return &BizError{
		Code:    defaultBizCode,
		Message: err.Error(),
	}
}

//Wraps wrap a error string
func Wraps(err string) *BizError {
	return &BizError{
		Code:    defaultBizCode,
		Message: err,
	}
}
