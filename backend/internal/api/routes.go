package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"learn-hub/config"
	"learn-hub/internal/api/handler"
	"learn-hub/internal/middleware"
	"learn-hub/pkg/oss"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// 公开路由
	public := router.Group("/api/v1")
	{
		// 认证路由
		auth := public.Group("/auth")
		{
			authHandler := handler.NewAuthHandler(db, cfg)
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/refresh", authHandler.RefreshToken)
		}
	}

	// 受保护的路由
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		// 初始化 OSS 客户端
		ossClient, err := oss.NewOSSClient(cfg.OSS)
		if err != nil {
			// 如果 OSS 初始化失败，使用 nil（可以在 handler 中处理）
			ossClient = nil
		}

		// 文件路由
		files := protected.Group("/files")
		{
			fileHandler := handler.NewFileHandler(db, cfg, ossClient)
			files.POST("/upload", fileHandler.Upload)
			files.POST("/delete", fileHandler.Delete)
			files.POST("/presigned-url", fileHandler.GetPresignedURL)
		}

		// 用户路由
		users := protected.Group("/users")
		{
			userHandler := handler.NewUserHandler(db)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.GET("/profile/me", userHandler.GetProfile)
		}

		// 资料路由
		materials := protected.Group("/materials")
		{
			materialHandler := handler.NewMaterialHandler(db)
			materials.GET("", materialHandler.ListMaterials)
			materials.GET("/:id", materialHandler.GetMaterial)
			materials.POST("", middleware.RequirePermission("materials:create"), materialHandler.CreateMaterial)
			materials.PUT("/:id", middleware.RequirePermission("materials:update"), materialHandler.UpdateMaterial)
			materials.DELETE("/:id", middleware.RequirePermission("materials:delete"), materialHandler.DeleteMaterial)
		}

		// 考试路由
		exams := protected.Group("/exams")
		{
			examHandler := handler.NewExamHandler(db)
			exams.GET("", examHandler.ListExams)
			exams.GET("/:id", examHandler.GetExam)
			exams.POST("/:id/start", examHandler.StartExam)
			exams.POST("/:id/submit", examHandler.SubmitExam)
			exams.GET("/:id/records", examHandler.GetExamRecords)
			exams.POST("", middleware.RequirePermission("exams:manage"), examHandler.CreateExam)
			exams.PUT("/:id", middleware.RequirePermission("exams:manage"), examHandler.UpdateExam)
			exams.DELETE("/:id", middleware.RequirePermission("exams:manage"), examHandler.DeleteExam)
		}

		// 题库路由
		questions := protected.Group("/questions")
		{
			questionHandler := handler.NewQuestionHandlerImpl(db)
			questions.GET("", questionHandler.ListQuestions)
			questions.GET("/:id", questionHandler.GetQuestion)
			questions.POST("", middleware.RequirePermission("questions:manage"), questionHandler.CreateQuestion)
			questions.PUT("/:id", middleware.RequirePermission("questions:manage"), questionHandler.UpdateQuestion)
			questions.DELETE("/:id", middleware.RequirePermission("questions:manage"), questionHandler.DeleteQuestion)
		}

		// 学习记录路由
		courseRecords := protected.Group("/course-records")
		{
			recordHandler := handler.NewCourseRecordHandlerImpl(db)
			courseRecords.GET("", recordHandler.ListRecords)
			courseRecords.GET("/:id", recordHandler.GetRecord)
			courseRecords.PUT("/:id", recordHandler.UpdateRecord)
		}

		// 菜单路由
		menus := protected.Group("/menus")
		{
			menuHandler := handler.NewMenuHandlerImpl(db)
			menus.GET("", menuHandler.GetMenus)
		}

		// 导入导出路由
		importExport := protected.Group("/import-export")
		{
			importExportHandler := handler.NewImportExportHandler(db)
			importExport.POST("/questions", middleware.RequirePermission("questions:manage"), importExportHandler.ImportQuestions)
			importExport.POST("/users", middleware.RequirePermission("users:manage"), importExportHandler.ImportUsers)
			importExport.GET("/exam-scores", middleware.RequirePermission("users:view"), importExportHandler.ExportExamScores)
		}

		// 管理员路由
		admin := protected.Group("/admin")
		admin.Use(middleware.RequirePermission("users:manage"))
		{
			adminHandler := handler.NewAdminHandlerImpl(db)
			admin.GET("/users", adminHandler.ListUsers)
			admin.GET("/users/:id", adminHandler.GetUser)
			admin.POST("/users", adminHandler.CreateUser)
			admin.PUT("/users/:id", adminHandler.UpdateUser)
			admin.DELETE("/users/:id", adminHandler.DeleteUser)

			admin.GET("/roles", adminHandler.ListRoles)
			admin.POST("/roles", adminHandler.CreateRole)
			admin.PUT("/roles/:id", adminHandler.UpdateRole)
			admin.DELETE("/roles/:id", adminHandler.DeleteRole)

			admin.GET("/permissions", adminHandler.ListPermissions)
			admin.POST("/permissions", adminHandler.CreatePermission)
		}
	}
}
