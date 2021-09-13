package controllers

import (
	"net/http"
	"strconv"

	"github.com/aysf/gojwt/lib/database"
	"github.com/aysf/gojwt/middlewares"
	"github.com/aysf/gojwt/models"
	"github.com/labstack/echo/v4"
)

// public or not authorized handler

func LoginUsersController(c echo.Context) error {
	user := models.User{}
	c.Bind(&user)

	users, err := database.LoginUsers(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "failed to login",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success login",
		"data":    users,
	})
}

func GetUsersController(c echo.Context) error {
	users, err := database.GetUsers()
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "failed to show users",
			"data":    "",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get all users",
		"data":    users,
	})
}

func GetUserByIdController(c echo.Context) error {
	loggedUserId := middlewares.ExtractTokenUserId(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "user id invalid",
			"data":    "",
		})
	}

	if loggedUserId != id {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "user id invalid",
			"data":    "",
		})
	}

	users, err := database.GetUserById((id))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "unauthorized access !",
			"data":    "",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "successfully loaded all users data",
		"data":    users,
	})
}

// 	CREATE FUNCTION

func CreateUserController(c echo.Context) error {
	inputUser := new(models.User)
	if err := c.Bind(inputUser); err != nil {
		return err
	}
	user, err := database.CreateUser(inputUser)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "a user has been created !",
		"data":    user,
	})
}

// DELETE FUNCTION

func DeleteUserController(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := database.DeleteUser(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "failed to delete ",
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "delete user succeed",
	})
}

func UpdateUserController(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	userInput := new(models.User)
	if err := c.Bind(userInput); err != nil {
		return err
	}
	user, err := database.UpdateUser(id, userInput)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "failed to update user",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "successfully updated",
		"data":    user,
	})
}

func GetUserDetailControllersTesting() echo.HandlerFunc {
	return GetUserByIdController
}
