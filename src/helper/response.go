package helper

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PaginateMeta struct {
	PerPage     int         `json:"per_page"`
	CurrentPage int         `json:"current_page"`
	TotalRow    int64       `json:"total_row"`
	TotalPage   int         `json:"total_page"`
	Info        interface{} `json:"info,omitempty"`
}

type ApiResponse struct {
	Success  bool        `json:"success"`
	Message  string      `json:"message"`
	Errors   interface{} `json:"errors"`
	Metadata interface{} `json:"metadata"`
	Data     interface{} `json:"data"`
}

func Result(c *fiber.Ctx, code int, res interface{}, err error) error {
	var errorObj *ErrorResponse
	var metadata interface{}

	if err != nil && code != fiber.StatusOK {
		message := parseError(err)

		errorObj = &ErrorResponse{
			Code:    code,
			Message: message,
		}
	}

	if data, ok := res.(map[string]interface{}); ok {
		if rows, ok1 := data["rows"]; ok1 {
			if count, ok2 := data["count"]; ok2 {

				page, _ := strconv.Atoi(c.Query("page", "1"))
				perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

				totalCount := count.(int64)

				metadata = PaginateMeta{
					PerPage:     perPage,
					CurrentPage: page,
					TotalRow:    totalCount,
					TotalPage:   int((totalCount + int64(perPage) - 1) / int64(perPage)),
				}

				res = rows
			}
		}
	}

	return c.Status(code).JSON(ApiResponse{
		Success:  code == fiber.StatusOK,
		Message:  ifElse(code == fiber.StatusOK, "Success", "Failed"),
		Errors:   errorObj,
		Metadata: metadata,
		Data:     res,
	})
}

func parseError(err error) string {
	goEnv := os.Getenv("GO_ENV")
	debugMode := goEnv == "local" || goEnv == "development" || goEnv == "staging"

	if err != nil {
		if debugMode {
			return err.Error()
		}
		return "INTERNAL ERROR"
	}

	return "Unknown error"
}

func ifElse[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

func PaginateArray[T any](arr []T, page, perPage int) []T {
	start := (page - 1) * perPage
	end := start + perPage

	if start > len(arr) {
		return []T{}
	}
	if end > len(arr) {
		end = len(arr)
	}

	return arr[start:end]
}

func ValidationMessage(err error) string {
	var msg string
	for _, e := range err.(validator.ValidationErrors) {
		msg += fmt.Sprintf("Field %s is %s; ", e.Field(), e.Tag())
	}
	return msg
}
