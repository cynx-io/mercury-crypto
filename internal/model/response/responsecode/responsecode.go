package responsecode

type ResponseCode string

func (r ResponseCode) String() string {
	return string(r)
}

const (

	// Expected Error
	CodeSuccess             ResponseCode = "00"
	CodeValidationError     ResponseCode = "VE"
	CodeAuthenticationError ResponseCode = "AU"
	CodeNotAllowed          ResponseCode = "NA"
	CodeNotFound            ResponseCode = "NF"
	CodeInvalidCredentials  ResponseCode = "IC"

	// Internal
	CodeInternalError ResponseCode = "I-IE"

	// External Errors
	CodeCoinGeckoError ResponseCode = "E-CG"

	// DB Error
	CodeTblUserError ResponseCode = "TBLUSR"
)

var ResponseCodeNames = map[ResponseCode]string{
	CodeSuccess:             "Success",
	CodeValidationError:     "Validation Error",
	CodeAuthenticationError: "Authentication Error",
	CodeInternalError:       "Internal Error",
	CodeNotAllowed:          "Not Allowed",
	CodeNotFound:            "Not Found",
	CodeInvalidCredentials:  "Invalid Credentials",

	CodeTblUserError: "TblUser Error",
}
