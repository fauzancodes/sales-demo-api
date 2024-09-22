package dto

type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}

type FindParameter struct {
	BaseFilter string
	Filter     string
	Order      string
	Limit      int
	Offset     int
}
