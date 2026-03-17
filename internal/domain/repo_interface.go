package domain

import (
	"context"
	"lms_system/internal/domain/entity"
)

type MainRepositoryInterface interface {
	Course() CourseRepositoryInterface
	Chapter() ChapterRepositoryInterface
	Lesson() LessonRepositoryInterface
	UserCourseAccess() UserCourseAccessInterface
	Attachment() AttachmentRepositoryInterface
}

type CourseRepositoryInterface interface {
	CreateCourse(ctx context.Context, entity entity.Course) (id uint, err error)
	UpdateCourseById(ctx context.Context, course entity.Course) error
	DeleteCourseById(ctx context.Context, courseId uint) error
	GetCourseById(ctx context.Context, id uint) (*entity.CourseAggregate, error)
	GetAllCourses(ctx context.Context) ([]entity.Course, error)
}

type ChapterRepositoryInterface interface {
	CreateChapter(ctx context.Context, courseId uint, entity *entity.Chapter) (id uint, err error)
	UpdateChapterById(ctx context.Context, chapter *entity.Chapter) error
	DeleteChapterById(ctx context.Context, chapterId uint) error
	GetChapterById(ctx context.Context, id uint) (*entity.Chapter, error)
	GetChaptersByCourseId(ctx context.Context, id uint) ([]entity.Chapter, error)
}

type LessonRepositoryInterface interface {
	CreateLesson(ctx context.Context, chapterId uint, entity entity.Lesson) (id uint, err error)
	UpdateLessonById(ctx context.Context, lesson entity.Lesson) error
	DeleteLessonById(ctx context.Context, lessonId uint) error
	GetLessonById(ctx context.Context, id uint) (*entity.Lesson, error)
	GetAllLessonsByChapterId(ctx context.Context, chapterId uint) ([]entity.Lesson, error)
}

type UserCourseAccessInterface interface {
	CreateUserCourseAccess(ctx context.Context, userAccess entity.UserCourseAccess) error
	GetByUserIdAndCourseId(ctx context.Context, userId string, courseId uint) (*entity.UserCourseAccess, error)
	GetAllByUserId(ctx context.Context, userId string) ([]entity.UserCourseAccess, error)
	UpdateAccess(ctx context.Context, access *entity.UserCourseAccess) error
}

type AttachmentRepositoryInterface interface {
	CreateAttachment(ctx context.Context, attachment *entity.Attachment) error
	GetAttachmentById(ctx context.Context, id uint) (*entity.Attachment, error)
	GetAttachmentsByLessonId(ctx context.Context, lessonId uint) ([]entity.Attachment, error)
	DeleteAttachment(ctx context.Context, id uint) error
}
