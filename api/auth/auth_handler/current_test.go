package auth_handler_test

import (
	"encoding/json"
	"fmt"
	"goauth/auth/auth_handler"
	"goauth/user/user_domain"
	"goauth/user/user_mocks"
	"goauth/utils/apperrors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCurrent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const USERID uint = 101

	t.Run("Not Found", func(t *testing.T) {
		// 1* Mock the response
		e := apperrors.NewNotFound("user", fmt.Sprintf("ID(%v)", USERID))

		// 2* mock user auth_service
		mockUserService := new(user_mocks.MockUserService)
		mockUserService.On("Get", USERID).Return(nil, e)
		// 3* Create the router and the server
		rr := httptest.NewRecorder()
		router := gin.Default()
		group := router.Group("/api/v1")

		// mock authentication middleware
		group.Use(func(ctx *gin.Context) {
			ctx.Set("user", &user_domain.User{
				Model: gorm.Model{
					ID: USERID,
				},
			})
		})

		auth_handler.NewHandler(&auth_handler.Config{
			RouterGroup: group,
			UserService: mockUserService,
		})
		//  crete http request
		req, err := http.NewRequest(http.MethodGet, "/api/v1/auth/current", nil)
		// 4* Test starts here
		assert.NoError(t, err)

		// expected response
		resBody, err := json.Marshal(gin.H{
			"error": e,
		})
		assert.NoError(t, err)
		// start the server
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, resBody, rr.Body.Bytes())
	})

	t.Run("No User in context", func(t *testing.T) {

		mockRsp := apperrors.NewInternal()

		router := gin.Default()
		groupV1 := router.Group("/api/v1")

		mockUserService := new(user_mocks.MockUserService)
		mockUserService.On("Get", USERID).Return(nil, mockRsp)

		auth_handler.NewHandler(&auth_handler.Config{
			RouterGroup: groupV1,
			UserService: mockUserService,
		})

		req, err := http.NewRequest(http.MethodGet, "/api/v1/auth/current", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		resBody, err := json.Marshal(gin.H{
			"error": mockRsp,
		})

		assert.NoError(t, err)
		assert.Equal(t, resBody, rr.Body.Bytes())

	})

	t.Run("Success", func(t *testing.T) {

		mockUserRsp := &user_domain.User{
			Email:    "kfekairi@digitrans.link",
			Password: "super-secret",
			Model: gorm.Model{
				ID: USERID,
			},
		}

		mockUserService := new(user_mocks.MockUserService)
		mockUserService.On("Get", USERID).Return(mockUserRsp, nil)

		// Response recorder for getting http response
		rr := httptest.NewRecorder()
		router := gin.Default()
		groupV1 := router.Group("/api/v1")

		// auth_mocks the auth middleware
		groupV1.Use(func(ctx *gin.Context) {
			ctx.Set("user", &user_domain.User{
				Model: gorm.Model{ID: USERID},
			})
		})

		auth_handler.NewHandler(&auth_handler.Config{
			RouterGroup: groupV1,
			UserService: mockUserService,
		})

		req, err := http.NewRequest(http.MethodGet, "/api/v1/auth/current", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, req)

		resBody, err := json.Marshal(gin.H{
			"user": mockUserRsp,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, resBody, rr.Body.Bytes())
		mockUserService.AssertExpectations(t)
	})

}
