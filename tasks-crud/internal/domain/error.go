package domain

// ErrorResponse структура для ошибок API
// @Description Стандартная структура ответа при ошибках
type ErrorResponse struct {
    Error   string    `json:"error" example:"Not found"`
    Details string    `json:"details,omitempty" example:"Task with id 5 not found"`
    Status  int       `json:"status" example:"404"`
    Time    string    `json:"timestamp" example:"2023-12-08T13:17:27Z"`
}