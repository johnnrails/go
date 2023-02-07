package common

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestConnectionToDatabase(t *testing.T) {
	asserts := assert.New(t)
	db := Init()

	_, err := os.Stat("./../gorm.db")
	asserts.NoError(err, "DB Should Exist")
	asserts.NoError(db.DB().Ping(), "Should be able to connect to DB with Ping")

	connection := GetDB()
	asserts.NoError(connection.DB().Ping(), "DB Should be able to ping")
	db.Close()

	// Test DB Exceptions
	os.Chmod("./../gorm.db", 0000)
	db = Init()
	asserts.Error(db.DB().Ping(), "Should not be able to ping with no permissions")
	db.Close()

	os.Chmod("./../gorm.db", 0644)
}

func TestConnectingTestDatabase(t *testing.T) {
	asserts := assert.New(t)
	// Test create & close DB
	db := TestDBInit()
	_, err := os.Stat("./../gorm_test.db")
	asserts.NoError(err, "Db should exist")
	asserts.NoError(db.DB().Ping(), "Db should be able to ping")
	db.Close()

	// Test testDB exceptions
	os.Chmod("./../gorm_test.db", 0000)
	db = TestDBInit()
	_, err = os.Stat("./../gorm_test.db")
	asserts.NoError(err, "Db should exist")
	asserts.Error(db.DB().Ping(), "Db should not be able to ping")
	os.Chmod("./../gorm_test.db", 0644)

	// Test close delete DB
	TestDBFree(db)
	_, err = os.Stat("./../gorm_test.db")

	asserts.Error(err, "Db should not exist")
}

func TestGenToken(t *testing.T) {
	asserts := assert.New(t)
	token := GenToken(2)
	asserts.Len(token, 115, "JWT's length should be 115")
}

func TestNewValidatorError(t *testing.T) {
	asserts := assert.New(t)

	type Login struct {
		Username string `form:"username" json:"username" binding:"exists,alphanum,min=8,max=255"`
		Password string `form:"password" json:"password" binding:"exists,min=8,max=255"`
	}

	var requestTests = []struct {
		bodyData     string
		expectedCode int
		response     string
		msg          string
	}{
		{
			`{"username": "wangzitian0","password": "0123456789"}`,
			http.StatusOK,
			`{"status":"you are logged in"}`,
			"valid data and should return StatusCreated",
		},
		{
			`{"username": "wangzitian0","password": "01234567866"}`,
			http.StatusUnauthorized,
			`{"errors":{"user":"wrong username or password"}}`,
			"wrong login status should return StatusUnauthorized",
		},
		{
			`{"username": "wangzitian0","password": "0122"}`,
			http.StatusUnprocessableEntity,
			`{"errors":{"Password":"{min: 8}"}}`,
			"invalid password of too short and should return StatusUnprocessableEntity",
		},
		{
			`{"username": "_wangzitian0","password": "0123456789"}`,
			http.StatusUnprocessableEntity,
			`{"errors":{"Username":"{key: alphanum}"}}`,
			"invalid username of non alphanum and should return StatusUnprocessableEntity",
		},
	}

	r := gin.Default()

	r.POST("/login", func(ctx *gin.Context) {
		var json Login
		if err := Bind(ctx, &json); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, NewValidatorError(err))
		} else {
			if json.Username == "wangzitian0" && json.Password == "0123456789" {
				ctx.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			} else {
				ctx.JSON(http.StatusUnauthorized, NewError("user", errors.New("wrong username or password")))
			}
		}
	})

	for _, testData := range requestTests {
		bodyData := testData.bodyData
		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(bodyData))
		req.Header.Set("Content-Type", "application/json")
		asserts.NoError(err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		asserts.Equal(testData.expectedCode, w.Code, "Response Status - "+testData.msg)
		asserts.Regexp(testData.response, w.Body.String(), "Response Content - "+testData.msg)
	}
}
