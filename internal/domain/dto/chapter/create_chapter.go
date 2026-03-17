package chapter

type CreateChapterRequest struct {
	CourseID      uint   `json:"course_id" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description"`
	OrderPosition int    `json:"order_position"`
}
