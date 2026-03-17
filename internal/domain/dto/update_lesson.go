package dto

type UpdateLessonRequest struct {
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description"`
	Content       string `json:"content"`
	OrderPosition int    `json:"order_position"`
}
