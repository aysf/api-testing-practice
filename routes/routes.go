package routes

import (
	"github.com/aysf/gojwt/constants"
	"github.com/aysf/gojwt/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()

	// not auth
	/*
		[POST] create user!, [POST] login user!
		[GET] book, [GET] book by id
	*/
	e.POST("/users", controllers.CreateUserController)
	e.POST("/login", controllers.LoginUsersController)
	e.GET("/books", controllers.GetBooksController)
	e.GET("/books/:id", controllers.GetBookByIdController)

	// auth
	/*
		get users!, get user by id!, delete user by id !, update user by id !
		post, delete, update book
	*/
	eJwt := e.Group("/success")
	eJwt.Use(middleware.JWT([]byte(constants.SECRET_JWT)))
	// user
	eJwt.GET("/users", controllers.GetUsersController)
	eJwt.GET("/users/:id", controllers.GetUserByIdController)
	eJwt.DELETE("/users/:id", controllers.DeleteUserController)
	eJwt.PUT("/users/:id", controllers.UpdateUserController)
	// book
	eJwt.POST("/books", controllers.CreateBookController)
	eJwt.DELETE("/books/:id", controllers.DeleteBookController)
	eJwt.PUT("/books/:id", controllers.UpdateBookController)
	return e
}
