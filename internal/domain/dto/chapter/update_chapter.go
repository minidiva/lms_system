package chapter

type UpdateChapterRequest struct {
	CourseID      string `json:"course_id" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description"`
	OrderPosition int    `json:"order_position"`
}
