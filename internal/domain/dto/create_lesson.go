package dto

type CreateLessonRequest struct {
	ChapterID     uint   `json:"chapter_id" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description"`
	Content       string `json:"content"`
	OrderPosition int    `json:"order_position"`
}
