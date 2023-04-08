package web

type SuccessfulResponse struct {
	Data any `json:"data"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

var (
	InvalidTokenResponse = ErrorResponse{
		Status:  401,
		Code:    "InvalidTokenError",
		Message: "The user token is not valid",
	}

	InvalidRequestResponse = ErrorResponse{
		Status:  400,
		Code:    "RequestError",
		Message: "The request is not valid",
	}

	NotFoundResponse = ErrorResponse{
		Status:  404,
		Code:    "NotFoundError",
		Message: "The requested resource was not found",
	}
)
