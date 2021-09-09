package database

import (
	"github.com/aysf/gojwt/config"
	"github.com/aysf/gojwt/models"
)

func CreateBook(book *models.Book) (interface{}, error) {
	if err := config.DB.Create(&book).Error; err != nil {
		return nil, err
	}
	return book, nil
}

func GetBooks() (interface{}, error) {
	books := new([]models.Book)
	if err := config.DB.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func GetBookByID(id int) (interface{}, error) {
	book := new(models.Book)
	if err := config.DB.Find(book, id).Error; err != nil {
		return nil, err
	}
	return book, nil
}

func UpdateBook(bookUpdate *models.Book, id int) (interface{}, error) {
	book := new(models.Book)
	tx := config.DB.First(&book)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected > 0 {
		book.Author = bookUpdate.Author
		book.Title = bookUpdate.Title
		book.Publisher = bookUpdate.Publisher
	}
	config.DB.Save(&book)
	return book, nil

}

func DeleteBook(id int) (interface{}, error) {
	book := new(models.Book)
	if err := config.DB.Delete(&book, id).Error; err != nil {
		return nil, err
	}
	return book, nil
}
