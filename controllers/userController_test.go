package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/aysf/gojwt/config"
	"github.com/aysf/gojwt/constants"
	"github.com/aysf/gojwt/middlewares"
	"github.com/aysf/gojwt/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func InitEchoTestAPI() *echo.Echo {

	config.InitDBTest()
	e := echo.New()
	return e
}

type UserResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func InsertDataUserForGetUsers() error {
	mock_user1 := models.User{
		FirstName: "Alta",
		LastName:  "Accedemy",
		Password:  "123",
		Email:     "alta@gmail.com",
	}
	mock_user2 := models.User{
		FirstName: "Roberto",
		LastName:  "Carlos",
		Password:  "123",
		Email:     "rc@gmail.com",
	}

	if err := config.DB.Save(&mock_user1).Error; err != nil {
		return err
	}
	if err := config.DB.Save(&mock_user2).Error; err != nil {
		return err
	}

	return nil
}

func TestGetUsersController(t *testing.T) {
	var testCases = []struct {
		name          string
		path          string
		expectCode    int
		expectMessage string
	}{
		{
			name:          "get user normal",
			path:          "/users",
			expectCode:    http.StatusOK,
			expectMessage: "success get all users-",
		},
	}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	type UserResponse struct {
		Message string
		Data    []models.User
	}

	for _, tc := range testCases {
		c.SetPath(tc.path)

		if assert.NoError(t, GetUsersController(c)) {
			assert.Equal(t, tc.expectCode, rec.Code)
			body := rec.Body.String()

			var response UserResponse
			err := json.Unmarshal([]byte(body), &response)

			if err != nil {
				assert.Error(t, err, "error")
			}
			assert.Equal(t, tc.expectMessage, response.Message)
		}
	}
}

func TestGetUserByIdController(t *testing.T) {
	var testCases = []struct {
		name          string
		path          string
		id            int
		expectCode    int
		expectMessage string
	}{
		{
			name:          "get correct id",
			path:          "/users",
			id:            1,
			expectCode:    http.StatusOK,
			expectMessage: "successfully loaded all users data",
		},
		{
			name:          "get invalid id",
			path:          "/users",
			id:            23,
			expectCode:    http.StatusBadRequest,
			expectMessage: "user id invalid",
		},
		{
			name:          "get unauthorized id",
			path:          "/users",
			id:            2,
			expectCode:    http.StatusUnauthorized,
			expectMessage: "unauthorized access !",
		},
	}

	user_test := models.User{
		Password: "123",
		Email:    "alta@gmail.com",
	}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()

	// create token
	var user models.User
	tx := config.DB.Where("email = ? AND password = ?", user_test.Email, user_test.Password).First(&user)
	if tx.Error != nil {
		t.Error(tx.Error)
	}
	token, err := middlewares.CreateToken(int(user.ID))
	if err != nil {
		panic(err)
	}

	for _, tc := range testCases {

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		res := httptest.NewRecorder()
		c := e.NewContext(req, res)
		c.SetPath(tc.path)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(tc.id))

		middleware.JWT([]byte(constants.SECRET_JWT))(GetUserDetailControllersTesting())(c)

		assert.Equal(t, tc.expectCode, res.Code)
		resBody := res.Body.String()

		var response UserResponse
		json.Unmarshal([]byte(resBody), &response)
		assert.Equal(t, tc.expectMessage, response.Status)

	}
}

// func TestCreateUserController(t *testing.T) {
// 	var testCases = []struct {
// 		name          string
// 		path          string
// 		expectCode    int
// 		expectMessage string
// 	}{
// 		{
// 			name:          "create user normal",
// 			path:          "/users",
// 			expectCode:    http.StatusOK,
// 			expectMessage: "user found",
// 		},
// 		{
// 			name:          "invalid user id parameter",
// 			path:          "/users",
// 			expectCode:    http.StatusBadRequest,
// 			expectMessage: "user id invalid",
// 		},
// 	}
// }

/*
go test ./controllers -coverprofile=cover.out && go tool cover -html=cover.out
go test -v -coverpkg=./controllers -coverprofile=profile.cov ./controllers && go tool cover -func profile.cov
*/

func GetUserDetailControllersTesting() echo.HandlerFunc {
	return GetUserByIdController
}
