package routes

import (
	"learn-hub-backend/controllers"
	"learn-hub-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// CORS中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 初始化控制器
	authCtrl := controllers.NewAuthController()
	menuCtrl := controllers.NewMenuController()
	ruleCtrl := controllers.NewRuleController()
	userCtrl := controllers.NewUserController()
	roleCtrl := controllers.NewRoleController()
	courseCtrl := controllers.NewCourseController()
	learnCtrl := controllers.NewLearnController()
	progressCtrl := controllers.NewProgressController()

	// API路由组
	api := r.Group("/api")
	{
		// 认证相关（不需要token）
		auth := api.Group("/login")
		{
			auth.POST("/account", authCtrl.Login)
			auth.GET("/captcha", authCtrl.GetCaptcha)
		}

		// 需要认证的路由
		api.Use(middleware.AuthMiddleware())
		{
			// 用户相关
			api.GET("/currentUser", authCtrl.GetCurrentUser)
			api.POST("/login/outLogin", authCtrl.Logout)

			// 菜单权限相关
			api.GET("/menu/list", menuCtrl.GetMenuList)
			api.GET("/user/permissions", menuCtrl.GetUserPermissions)

			// 规则管理相关
			api.GET("/rule", ruleCtrl.GetRuleList)
			api.POST("/rule", ruleCtrl.CreateOrUpdateRule)

			// 用户管理相关
			api.GET("/user/list", userCtrl.GetUserList)
			api.POST("/user", userCtrl.CreateUser)
			api.PUT("/user/:id", userCtrl.UpdateUser)
			api.DELETE("/user/:id", userCtrl.DeleteUser)
			api.GET("/user/roles", userCtrl.GetAllRoles)

			// 角色管理相关
			api.GET("/role/list", roleCtrl.GetRoleList)
			api.POST("/role", roleCtrl.CreateRole)
			api.PUT("/role/:id", roleCtrl.UpdateRole)
			api.DELETE("/role/:id", roleCtrl.DeleteRole)
			api.GET("/role/menus", roleCtrl.GetAllMenus)

			// 课程管理相关（管理员）
			api.GET("/course/list", courseCtrl.GetCourseList)
			api.GET("/course/:id", courseCtrl.GetCourseDetail)
			api.POST("/course", courseCtrl.CreateCourse)
			api.PUT("/course/:id", courseCtrl.UpdateCourse)
			api.DELETE("/course/:id", courseCtrl.DeleteCourse)
			api.POST("/course/:id/publish", courseCtrl.PublishCourse)

			// 学习相关（员工端）
			api.GET("/learn/courses", learnCtrl.GetCourses)
			api.GET("/learn/course/:id", learnCtrl.GetCourseDetail)
			api.POST("/learn/course/:id/progress", learnCtrl.UpdateProgress)
			api.POST("/learn/course/:id/complete", learnCtrl.CompleteCourse)

			// 学习进度查看（管理员）
			api.GET("/admin/course/:id/progress", progressCtrl.GetCourseProgress)
			api.GET("/admin/user/:id/progress", progressCtrl.GetUserProgress)
		}
	}

	return r
}
