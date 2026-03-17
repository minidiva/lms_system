package dto

type BuyCourseRequest struct {
	CourseId uint   `json:"course_id"`
	UserUUID string // ← UUID напрямую из Keycloak, не приходит из JSON
}
