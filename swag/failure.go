package swag

type Validation struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Validation Failed"`
}

type NotFound struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" example:"Resource Not Found"`
}

type BadRequest struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad Request"`
}

type SystemError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Internal Server Error"`
}
