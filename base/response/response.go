package response

import "github.com/gin-gonic/gin"

type BaseResponse struct {
	Status  string      `json:"status"`
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func Success(c *gin.Context, message string, result interface{}) {
	c.JSON(200, BaseResponse{
		Status:  "success",
		Message: message,
		Result:  result,
	})
}

func Created(c *gin.Context, message string, result interface{}) {
	c.JSON(201, BaseResponse{
		Status:  "success",
		Message: message,
		Result:  result,
	})
}

func Error(c *gin.Context, httpCode int, code string, message string) {
	c.JSON(httpCode, BaseResponse{
		Status:  "fail",
		Code:    code,
		Message: message,
		Result:  nil,
	})
}

func BadRequest(c *gin.Context, code string, message string) {
	Error(c, 400, code, message)
}

func NotFound(c *gin.Context, code string, message string) {
	Error(c, 404, code, message)
}

func Unauthorized(c *gin.Context, code string, message string) {
	Error(c, 401, code, message)
}

func InternalServerError(c *gin.Context, code string, message string) {
	Error(c, 500, code, message)
}
