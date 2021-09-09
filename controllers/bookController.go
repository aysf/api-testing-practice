package controllers

import (
	"net/http"
	"strconv"

	"github.com/aysf/gojwt/lib/database"
	"github.com/aysf/gojwt/models"
	"github.com/labstack/echo/v4"
)

func GetBooksController(c echo.Context) error {
	books, err := database.GetBooks()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "failed to show books",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get all books",
		"data":    books,
	})
}

func GetBookByIdController(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	book, err := database.GetBookByID(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "query success",
		"book":    book,
	})
}

func CreateBookController(c echo.Context) error {
	bookInput := new(models.Book)
	if err := c.Bind(bookInput); err != nil {
		return err
	}
	book, err := database.CreateBook(bookInput)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success create book",
		"book":    book,
	})
}

func DeleteBookController(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := database.DeleteBook(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "failed to delete ",
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "successfully delete book id " + c.Param("id"),
	})
}

func UpdateBookController(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	bookInput := new(models.Book)
	if err := c.Bind(bookInput); err != nil {
		return err
	}
	book, err := database.UpdateBook(bookInput, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "failed to update data",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "successfully updated",
		"book":    book,
	})
}
