package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
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
	Message string        `json:"message"`
	Data    []models.User `json:"data"`
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
			expectMessage: "success get all users",
		},
	}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

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
			// need additional function
			name:          "get unauthorized id",
			path:          "/users",
			id:            2,
			expectCode:    http.StatusBadRequest, // temporary
			expectMessage: "user id invalid",     // temporary
		},
	}

	user_test := models.User{
		Password: "123",
		Email:    "alta@gmail.com",
	}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()

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
		fmt.Println("RESPONSE resBODY -->>", res)

		assert.Equal(t, tc.expectCode, res.Code)
		resBody := res.Body.String()

		var response UserResponse
		json.Unmarshal([]byte(resBody), &response)

		assert.Equal(t, tc.expectMessage, response.Message)

	}
}

func TestCreateUserController(t *testing.T) {
	var testCases = []struct {
		name          string
		reqBody       string
		expectCode    int
		expectMessage string
	}{
		{
			name:          "success created user",
			reqBody:       `{"firstName": "Ananto", "lastName":"Wicaksono", "email": "aw@test.com", "password": "ay123"}`,
			expectCode:    http.StatusOK,
			expectMessage: "created user success !",
		},
		// add validator !
		// {
		// 	name:          "missing information",
		// 	reqBody:       `{"firstName": "Budi", "lastName":"Anduk", "email": "ba@test.com", "password": "ba123"}`,
		// 	expectCode:    http.StatusBadRequest,
		// 	expectMessage: "created user success !",
		// },
	}

	e := InitEchoTestAPI()

	type userForm struct {
		FirstName string
		LastName  string
		Email     string
		Password  string
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte(tc.reqBody)))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)

		var request userForm
		json.Unmarshal([]byte(tc.reqBody), &request)

		if assert.NoError(t, CreateUserController(c)) {

			var response UserResponse
			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, tc.expectCode, res.Code)
			assert.Equal(t, request.Email, response.Data[0].Email)

		}
	}

}

func TestDeleteUserController(t *testing.T) {
	// define test cases
	var testCases = []struct {
		name          string
		path          string
		id            int
		expectCode    int
		expectMessage string
	}{
		{
			name:          "delete succeeded",
			path:          "/users",
			id:            2,
			expectCode:    http.StatusOK,
			expectMessage: "delete user succeed",
		},
		{
			name:          "user not found",
			path:          "/users",
			id:            2,
			expectCode:    http.StatusOK,
			expectMessage: "delete user succeed",
		},
	}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()
	for _, tc := range testCases {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		res := httptest.NewRecorder()

		c := e.NewContext(req, res)
		c.SetPath(tc.path)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(tc.id))
		if assert.NoError(t, DeleteUserController(c)) {
			assert.Equal(t, res.Code, tc.expectCode)
			body := res.Body.String()

			var response UserResponse
			json.Unmarshal([]byte(body), &response)
			assert.Equal(t, tc.expectMessage, response.Message)
		}

	}
}

func TestUpdateUserController(t *testing.T) {
	var testCases = []struct {
		name          string
		reqBody       string
		path          string
		id            int
		expectCode    int
		expectMessage string
	}{
		{
			name:          "success update user",
			reqBody:       `{"firstName": "Ananto", "lastName":"Wicaksono", "email": "aw@test.com", "password": "ay123"}`,
			path:          "/users",
			id:            1,
			expectCode:    http.StatusOK,
			expectMessage: "updated user success !",
		},
	}

	user_test := models.User{
		Password: "123",
		Email:    "alta@gmail.com",
	}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()

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
		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(tc.reqBody))
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		var request models.User
		json.Unmarshal([]byte(tc.reqBody), &request)

		res := httptest.NewRecorder()

		c := e.NewContext(req, res)
		c.SetPath(tc.path)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(tc.id))

		if assert.NoError(t, CreateUserController(c)) {
			assert.Equal(t, tc.expectCode, res.Code)

			body := res.Body.String()
			var response UserResponse
			json.Unmarshal([]byte(body), &response)

			assert.Equal(t, request.Email, response.Data[0].Email)
		}
	}

}

func TestLoginUsersController(t *testing.T) {

}

/*
go test ./controllers -coverprofile=cover.out && go tool cover -html=cover.out
go test -v -coverpkg=./controllers -coverprofile=profile.cov ./controllers && go tool cover -func profile.cov
*/
