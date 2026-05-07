package domain

import (
	"context"
	"lms_system/internal/domain/common"
	"lms_system/internal/domain/dto"
	"lms_system/internal/domain/entity"
	"mime/multipart"
)

type ServiceInterface interface {
	UserServiceInterface
	AdminServiceInterface
	AttachmentServiceInterface
}

type UserServiceInterface interface {
	BuyCourse(ctx context.Context, request dto.BuyCourseRequest) error
	GetAllCourses(ctx context.Context) ([]entity.Course, error)
	GetCourseById(ctx context.Context, id uint) (entity.CourseAggregate, error)
	GetChaptersInfoByCourseId(ctx context.Context, id uint) ([]entity.ChapterInfoAggregate, error)
	GetLessonById(ctx context.Context, id uint, userUUID string, role common.UserRole) (entity.Lesson, error)
}

type AdminServiceInterface interface {
	CreateCourse(ctx context.Context, entity entity.Course) (id uint, err error)
	CreateChapter(ctx context.Context, courseId uint, entity entity.Chapter) (id uint, err error)
	CreateLesson(ctx context.Context, chapterId uint, entity entity.Lesson) (id uint, err error)
	UpdateCourseById(ctx context.Context, course entity.Course) error
	UpdateChapterById(ctx context.Context, chapter entity.Chapter) error
	UpdateLessonById(ctx context.Context, lesson entity.Lesson) error
	DeleteCourseById(ctx context.Context, courseId uint) error
	DeleteChapterById(ctx context.Context, chapterId uint) error
	DeleteLessonById(ctx context.Context, lessonId uint) error
}

type AttachmentServiceInterface interface {
	UploadAttachment(ctx context.Context, lessonId uint, file multipart.File, fileHeader *multipart.FileHeader) (*entity.Attachment, error)
	GetAttachment(ctx context.Context, attachmentId uint, userId string, role common.UserRole) (*entity.Attachment, error)
	DownloadAttachment(ctx context.Context, attachmentId uint, userId string, role common.UserRole) (string, error)
	GetLessonAttachments(ctx context.Context, lessonId uint) ([]entity.Attachment, error)
	DeleteAttachment(ctx context.Context, attachmentId uint) error
	CheckUserAccessToLesson(ctx context.Context, userId string, role common.UserRole, lessonId uint) (bool, error)
}
