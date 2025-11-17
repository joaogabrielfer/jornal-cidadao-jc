package main

import(
	"github/jornal-cidadao-jc/internal/handlers"
	"github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.Engine, handler *handlers.Handler){
	api := router.Group("/api")
	{
		api.POST("/users", handler.CreateUser)

		api.GET("/charges", handler.GetChargesList)
		api.GET("/charges/random", handler.GetRandomCharge)

		api.GET("/materias", handler.GetArticles)
		api.GET("/materia/:id", handler.GetArticleByID)

		api.POST("/enquete/votar/:id", handler.UpdateVoteCount)
	}

	router.GET("/", handler.GetIndexPage)
	router.GET("/charge", handler.GetNoIdChargePage)
	router.GET("/charge/:id", handler.GetChargePage)
	router.GET("/cadastro", handler.GetSignupPage)
	router.GET("/login", handler.GetLoginPage)
	router.GET("/ultimas", handler.ShowJornalCidadaoDashboard)

	admin := router.Group("/admin")
	{
		admin.GET("/", handler.GetAdminPage)
		admin.GET("/adicionar-charge", handler.GetUploadChargePage)
		admin.GET("/charges", handler.GetDeleteChargePage)
		admin.GET("/users", handler.GetUsersAdminPage)
		admin.GET("/materia", handler.GetUploadArticlePage)
		admin.GET("/materias", handler.GetArticlesPage)
		admin.GET("/materias/:id/edit", handler.GetUpdateArticlePage)

		adminApi := admin.Group("/api")
		{
			adminApi.POST("/charge", handler.UploadCharge)
			adminApi.DELETE("/charge/:id", handler.DeleteCharge)

			adminApi.GET("/users", handler.GetUsers)
			adminApi.DELETE("/user/:id", handler.DeleteUser)

			adminApi.POST("/materia", handler.UploadArticle)
			adminApi.PUT("/materia/:id", handler.UpdateArticle)
			adminApi.DELETE("/materia/:id", handler.DeleteArticle)
		}
	}

}
