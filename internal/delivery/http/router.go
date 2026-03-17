package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "lms_system/docs"
	"lms_system/internal/delivery/http/handlers/attachment"
	"lms_system/internal/delivery/http/handlers/auth"
	"lms_system/internal/delivery/http/handlers/chapter"
	"lms_system/internal/delivery/http/handlers/course"
	"lms_system/internal/delivery/http/handlers/file"
	"lms_system/internal/delivery/http/handlers/lesson"
	"lms_system/internal/delivery/http/middleware"
	"lms_system/internal/domain"
	"lms_system/internal/domain/common"
)

func NewRouter(service domain.ServiceInterface, authService domain.AuthServiceInterface, fileService domain.FileServiceInterface) *chi.Mux {
	router := chi.NewRouter()

	// Global middleware
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)

	// Create handlers
	courseHandler := course.NewHandler(service)
	chapterHandler := chapter.NewHandler(service)
	lessonHandler := lesson.NewHandler(service)
	authHandler := auth.NewHandler(authService)
	fileHandler := file.NewHandler(fileService)
	attachmentHandler := attachment.NewHandler(service)

	router.Route("/api/v1", func(r chi.Router) {

		// Auth routes
		// Tested
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.Refresh)
		})

		// Public routes (no auth)
		// Not tested
		r.Route("/public", func(r chi.Router) {
			r.Get("/courses", courseHandler.GetAllCourses)
			r.Get("/courses/{id}", courseHandler.GetCourseById)
			r.Get("/courses/{id}/chapters", courseHandler.GetChaptersInfoByCourseId)
		})

		// User routes (user auth required)
		r.Route("/user", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)

			r.Put("/profile", authHandler.UpdateProfile)
			r.Put("/change-password", authHandler.ChangePassword)

			r.Route("/courses", func(r chi.Router) {
				r.Get("/", courseHandler.GetAllCourses)
				r.Get("/{id}", courseHandler.GetCourseById)
				r.Post("/buy", courseHandler.BuyCourse) // +
				r.Get("/{id}/chapters", courseHandler.GetChaptersInfoByCourseId)
			})

			r.Route("/lessons", func(r chi.Router) {
				r.Get("/{id}", lessonHandler.GetLessonById)
				r.Get("/{lessonId}/attachments", attachmentHandler.GetLessonAttachments)
			})

			r.Route("/files", func(r chi.Router) {
				r.Get("/download", fileHandler.DownloadFile)
				r.Get("/url", fileHandler.GetFileURL)
			})

			r.Get("/attachments/{attachmentId}/download", attachmentHandler.DownloadAttachment)
		})

		// Attachment management (admin and teacher auth)
		r.Route("/attachments", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)

			r.Group(func(r chi.Router) {
				r.Use(middleware.OnlyRoles(common.RoleAdmin, common.RoleTeacher))
				r.Post("/lessons/{lessonId}/upload", attachmentHandler.UploadAttachment)
			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.OnlyRoles(common.RoleAdmin))
				r.Delete("/{attachmentId}", attachmentHandler.DeleteAttachment)
			})
		})

		// User management
		r.Route("/admin", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)
			r.Use(middleware.OnlyRoles(common.RoleAdmin))
			r.Post("/register", authHandler.RegisterUser)
		})

		// Course management
		r.Route("/courses", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)
			r.Group(func(r chi.Router) {
				r.Use(middleware.OnlyRoles(common.RoleAdmin))
				r.Post("/create", courseHandler.CreateCourse)
				r.Delete("/delete/{id}", courseHandler.DeleteCourseById)
			})
			r.Group(func(r chi.Router) {
				r.Use(middleware.OnlyRoles(common.RoleTeacher, common.RoleAdmin))
				r.Put("/update/{id}", courseHandler.UpdateCourseById)
			})
		})
		// Chapter management
		r.Route("/chapters", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)

			r.Group(func(r chi.Router) {
				r.Use(middleware.OnlyRoles(common.RoleAdmin))
				r.Delete("/delete/{id}", chapterHandler.DeleteChapterById)
			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.OnlyRoles(common.RoleAdmin, common.RoleTeacher))
				r.Post("/create", chapterHandler.CreateChapterStandalone)
				r.Put("/update/{id}", chapterHandler.UpdateChapterById)
			})
		})

		// Lesson management
		r.Route("/lessons", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)

			r.Group(func(r chi.Router) {
				r.Use(middleware.OnlyRoles(common.RoleAdmin))
				r.Delete("/{id}", lessonHandler.DeleteLessonById)
			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.OnlyRoles(common.RoleAdmin, common.RoleTeacher))
				r.Post("/", lessonHandler.CreateLessonStandalone)
				r.Put("/{id}", lessonHandler.UpdateLessonById)
			})
		})

		// File management
		r.Route("/files", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)
			r.Use(middleware.OnlyRoles(common.RoleAdmin, common.RoleTeacher))

			r.Post("/upload", fileHandler.UploadFile)
			r.Delete("/", fileHandler.DeleteFile)
			r.Post("/courses/{courseId}/upload", fileHandler.UploadCourseFile)
			r.Post("/lessons/{lessonId}/upload", fileHandler.UploadLessonFile)
		})
	})
	// Swagger documentation
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8082/swagger/doc.json"),
	))

	// Serve static swagger files
	router.Handle("/swagger/swagger.yaml", http.FileServer(http.Dir("./docs/")))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("LMS system is running"))
		if err != nil {
			return
		}
	})

	return router
}
