package entity

import "time"

type UserCourseAccess struct {
	UserID    string `gorm:"type:uuid"` // ← UUID как string
	CourseID  uint
	Unlocked  bool
	CreatedAt time.Time
}

func (UserCourseAccess) TableName() string {
	return "user_access_course"
}
