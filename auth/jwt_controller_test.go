package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/System-Glitch/goyave/v2"
	"github.com/System-Glitch/goyave/v2/config"
	"github.com/System-Glitch/goyave/v2/database"
	"golang.org/x/crypto/bcrypt"
)

const testUserPassword = "secret"

type JWTControllerTestSuite struct {
	goyave.TestSuite
}

func (suite *JWTControllerTestSuite) SetupSuite() {
	config.Set("dbConnection", "mysql")
	database.ClearRegisteredModels()
	database.RegisterModel(&TestUser{})

	database.Migrate()
}

func (suite *JWTControllerTestSuite) SetupTest() {
	password, err := bcrypt.GenerateFromPassword([]byte(testUserPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user := &TestUser{
		Name:     "Admin",
		Email:    "johndoe@example.org",
		Password: string(password),
	}
	database.GetConnection().Create(user)
}

func (suite *JWTControllerTestSuite) TestLogin() {
	controller := NewJWTController(&TestUser{})
	suite.NotNil(controller)

	request := suite.CreateTestRequest(nil)
	request.Data = map[string]interface{}{
		"username": "johndoe@example.org",
		"password": testUserPassword,
	}
	writer := httptest.NewRecorder()
	response := suite.CreateTestResponse(writer)

	controller.Login(response, request)
	result := writer.Result()
	suite.Equal(http.StatusOK, result.StatusCode)

	json := map[string]string{}
	err := suite.GetJSONBody(result, &json)
	suite.Nil(err)

	if err == nil {
		token, ok := json["token"]
		suite.True(ok)
		suite.NotEmpty(token)
	}

	request.Data = map[string]interface{}{
		"username": "johndoe@example.org",
		"password": "wrongpassword",
	}
	writer = httptest.NewRecorder()
	response = suite.CreateTestResponse(writer)

	controller.Login(response, request)
	result = writer.Result()
	suite.Equal(http.StatusUnprocessableEntity, result.StatusCode)

	json = map[string]string{}
	err = suite.GetJSONBody(result, &json)
	suite.Nil(err)

	if err == nil {
		errMessage, ok := json["validationError"]
		suite.True(ok)
		suite.Equal("These credentials don't match our records.", errMessage)
	}
}

func (suite *JWTControllerTestSuite) TestLoginPanic() {
	suite.Panics(func() {
		request := suite.CreateTestRequest(nil)
		request.Data = map[string]interface{}{
			"username": "johndoe@example.org",
			"password": testUserPassword,
		}
		controller := NewJWTController(&TestUserOverride{})
		controller.Login(suite.CreateTestResponse(httptest.NewRecorder()), request)
	})
}

func (suite *JWTControllerTestSuite) TestRoutes() {
	suite.RunServer(func(router *goyave.Router) {
		suite.NotNil(JWTRoutes(router, &TestUser{}))
	}, func() {
		json, err := json.Marshal(map[string]string{
			"username": "johndoe@example.org",
			"password": testUserPassword,
		})
		if err != nil {
			panic(err)
		}

		headers := map[string]string{"Content-Type": "application/json"}
		resp, err := suite.Post("/auth/login", headers, strings.NewReader(string(json)))
		suite.Nil(err)
		suite.NotNil(resp)
		if resp != nil {
			suite.Equal(200, resp.StatusCode)
		}
	})
}

func (suite *JWTControllerTestSuite) TearDownTest() {
	suite.ClearDatabase()
}

func (suite *JWTControllerTestSuite) TearDownSuite() {
	database.GetConnection().DropTable(&TestUser{})
	database.ClearRegisteredModels()
}

func TestJWTControllerSuite(t *testing.T) {
	goyave.RunTest(t, new(JWTControllerTestSuite))
}
