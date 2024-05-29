package main

import (
	"github.com/gin-gonic/gin"
	"github.com/eldhopaulose/ReviewRoster/src/config/db"
	"github.com/eldhopaulose/ReviewRoster/src/initialize"
	"github.com/eldhopaulose/ReviewRoster/src/controllers"
	"github.com/eldhopaulose/ReviewRoster/src/middleware"
)

func main() {
	// Load environment variables
	initialize.LoadEnv()

	// Connect to the database
	db.ConnectDB()
	db.Migrate()

	r := gin.Default()

	//Auth
	r.POST("/login", controllers.LoginUserController)
	r.POST("/signup", controllers.CreateUserController)


	r.GET("/", controllers.GetAllBooksController)
	r.GET("/:id", controllers.GetBookController)

	// Middleware
	r.Use(middleware.AuthMiddleware())
	r.POST("/", controllers.CreateBookController)
	r.PUT("/:id", controllers.UpdateBookController)
	r.DELETE("/:id", controllers.DeleteBookController)
	r.GET("/users", controllers.GetAllUsersController)

	r.Run(":8080")                                // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
