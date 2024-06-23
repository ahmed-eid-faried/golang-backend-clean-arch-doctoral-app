package http

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"testing"
//
//

// 	"github.com/stretchr/testify/assert"

// 	"main/internal/user/dto"
// 	"main/internal/user/model"
// 	"main/pkg/jtoken"
// )

// // Login
// // =================================================================================================

// func TestUserAPI_LoginSuccess(t *testing.T) {
// 	dbTest.Create(context.Background(), &model.User{
// 		Email:    "login@test.com",
// 		Password: "test123456",
// 	})

// 	user := &dto.LoginReq{
// 		Email:    "login@test.com",
// 		Password: "test123456",
// 	}
// 	writer := makeRequest("POST", "/auth/login", user, "")
// 	assert.Equal(t, http.StatusOK, writer.Code)
// }

// func TestUserAPI_LoginInvalidFieldType(t *testing.T) {
// 	user := map[string]interface{}{
// 		"email":    1,
// 		"password": "test123456",
// 	}
// 	writer := makeRequest("POST", "/auth/login", user, "")
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusBadRequest, writer.Code)
// 	assert.Equal(t, "Invalid parameters", response["error"]["message"])
// }

// func TestUserAPI_LoginInvalidEmailFormat(t *testing.T) {
// 	user := &dto.LoginReq{
// 		Email:    "invalid",
// 		Password: "test123456",
// 	}
// 	writer := makeRequest("POST", "/auth/login", user, "")
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusInternalServerError, writer.Code)
// 	assert.Equal(t, "Something went wrong", response["error"]["message"])
// }

// func TestUserAPI_LoginInvalidPassword(t *testing.T) {
// 	user := &dto.LoginReq{
// 		Email:    "test@test.com",
// 		Password: "test",
// 	}
// 	writer := makeRequest("POST", "/auth/login", user, "")
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusInternalServerError, writer.Code)
// 	assert.Equal(t, "Something went wrong", response["error"]["message"])
// }

// func TestUserAPI_LoginUserNotFound(t *testing.T) {
// 	user := &dto.LoginReq{
// 		Email:    "notfound@test.com",
// 		Password: "test123456",
// 	}
// 	writer := makeRequest("POST", "/auth/login", user, "")
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusInternalServerError, writer.Code)
// 	assert.Equal(t, "Something went wrong", response["error"]["message"])
// }

// func TestUserAPI_LoginUserWrongPassword(t *testing.T) {
// 	user := &dto.LoginReq{
// 		Email:    "test@test.com",
// 		Password: "test1234567",
// 	}
// 	writer := makeRequest("POST", "/auth/login", user, "")
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusInternalServerError, writer.Code)
// 	assert.Equal(t, "Something went wrong", response["error"]["message"])
// }

// // Register
// // =================================================================================================

// func TestUserAPI_RegisterSuccess(t *testing.T) {
// 	defer cleanData()

// 	user := &dto.RegisterReq{
// 		Email:    "register@test.com",
// 		Password: "test123456",
// 	}
// 	writer := makeRequest("POST", "/auth/register", user, "")
// 	assert.Equal(t, http.StatusOK, writer.Code)
// }

// func TestUserAPI_RegisterInvalidFieldType(t *testing.T) {
// 	user := map[string]interface{}{
// 		"email":    1,
// 		"password": "test123456",
// 	}
// 	writer := makeRequest("POST", "/auth/register", user, "")
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusBadRequest, writer.Code)
// 	assert.Equal(t, "Invalid parameters", response["error"]["message"])
// }

// func TestUserAPI_RegisterInvalidEmail(t *testing.T) {
// 	user := map[string]interface{}{
// 		"email":    "invalid",
// 		"password": "test123456",
// 	}
// 	writer := makeRequest("POST", "/auth/register", user, "")
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusInternalServerError, writer.Code)
// 	assert.Equal(t, "Something went wrong", response["error"]["message"])
// }

// func TestUserAPI_RegisterInvalidPassword(t *testing.T) {
// 	user := map[string]interface{}{
// 		"email":    "register@test.com",
// 		"password": "test",
// 	}
// 	writer := makeRequest("POST", "/auth/register", user, "")
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusInternalServerError, writer.Code)
// 	assert.Equal(t, "Something went wrong", response["error"]["message"])
// }

// func TestUserAPI_RegisterEmailExist(t *testing.T) {
// 	defer cleanData()

// 	dbTest.Create(context.Background(), &model.User{
// 		Email:    "emailexist@test.com",
// 		Password: "password",
// 	})

// 	user := map[string]interface{}{
// 		"email":    "emailexist@test.com",
// 		"password": "test123456",
// 	}
// 	writer := makeRequest("POST", "/auth/register", user, "")
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusInternalServerError, writer.Code)
// 	assert.Equal(t, "Something went wrong", response["error"]["message"])
// }

// // GetMe
// // =================================================================================================

// func TestUserAPI_GetMeSuccess(t *testing.T) {
// 	writer := makeRequest("GET", "/auth/me", nil, accessToken())
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusOK, writer.Code)
// 	assert.Equal(t, "test@test.com", response["result"]["email"])
// 	assert.Equal(t, "", response["result"]["password"])
// }

// func TestUserAPI_GetMeUnauthorized(t *testing.T) {
// 	writer := makeRequest("GET", "/auth/me", nil, "")
// 	assert.Equal(t, http.StatusUnauthorized, writer.Code)
// }

// func TestUserAPI_GetMeUserNotFound(t *testing.T) {
// 	token := jtoken.GenerateAccessToken(map[string]interface{}{
// 		"id": "user-not-found",
// 	})

// 	writer := makeRequest("GET", "/auth/me", nil, token)
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusInternalServerError, writer.Code)
// 	assert.Equal(t, "Something went wrong", response["error"]["message"])
// }

// func TestUserAPI_GetMeInvalidTokenType(t *testing.T) {
// 	writer := makeRequest("GET", "/auth/me", nil, refreshToken())
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusUnauthorized, writer.Code)
// }

// // Refresh Token
// // =================================================================================================

// func TestUserAPI_RefreshTokenSuccess(t *testing.T) {
// 	writer := makeRequest("POST", "/auth/refresh", nil, refreshToken())
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusOK, writer.Code)
// 	assert.NotNil(t, response["result"]["access_token"])
// }

// func TestUserAPI_RefreshTokenUnauthorized(t *testing.T) {
// 	writer := makeRequest("POST", "/auth/refresh", nil, "")
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusUnauthorized, writer.Code)
// }

// func TestUserAPI_RefreshTokenInvalidTokenType(t *testing.T) {
// 	writer := makeRequest("POST", "/auth/refresh", nil, accessToken())
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusUnauthorized, writer.Code)
// }

// func TestUserAPI_RefreshTokenUserNotFound(t *testing.T) {
// 	token := jtoken.GenerateRefreshToken(map[string]interface{}{
// 		"id": "user-not-found",
// 	})

// 	writer := makeRequest("POST", "/auth/refresh", nil, token)
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusInternalServerError, writer.Code)
// 	assert.Equal(t, "Something went wrong", response["error"]["message"])
// }

// // Change Password
// // =================================================================================================

// func TestUserAPI_UpdateUserSuccess(t *testing.T) {
// 	defer cleanData()

// 	user := model.User{Email: "UpdateUser1@gmail.com", Password: "123456"}
// 	dbTest.Create(context.Background(), &user)

// 	token := jtoken.GenerateAccessToken(map[string]interface{}{
// 		"id": user.ID,
// 	})

// 	req := &dto.UpdateUserReq{
// 		Password:    "123456",
// 		NewPassword: "new123456",
// 	}

// 	writer := makeRequest("PUT", "/auth/update-user", req, token)
// 	assert.Equal(t, http.StatusOK, writer.Code)
// }

// func TestUserAPI_UpdateUserUnauthorized(t *testing.T) {
// 	req := &dto.UpdateUserReq{
// 		Password:    "123456",
// 		NewPassword: "new123456",
// 	}

// 	writer := makeRequest("PUT", "/auth/update-user", req, "")
// 	assert.Equal(t, http.StatusUnauthorized, writer.Code)
// }

// func TestUserAPI_UpdateUserIsWrong(t *testing.T) {
// 	req := &dto.UpdateUserReq{
// 		Password:    "wrong123456",
// 		NewPassword: "new123456",
// 	}

// 	writer := makeRequest("PUT", "/auth/update-user", req, accessToken())
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusInternalServerError, writer.Code)
// 	assert.Equal(t, "Something went wrong", response["error"]["message"])
// }

// func TestUserAPI_UpdateUserInvalidNewPassword(t *testing.T) {
// 	req := &dto.UpdateUserReq{
// 		Password:    "test123456",
// 		NewPassword: "new",
// 	}

// 	writer := makeRequest("PUT", "/auth/update-user", req, accessToken())
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusInternalServerError, writer.Code)
// 	assert.Equal(t, "Something went wrong", response["error"]["message"])
// }

// func TestUserAPI_UpdateUserInvalidFieldType(t *testing.T) {
// 	req := map[string]interface{}{
// 		"password":     1,
// 		"new_password": "new",
// 	}

// 	writer := makeRequest("PUT", "/auth/update-user", req, accessToken())
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	assert.Equal(t, http.StatusBadRequest, writer.Code)
// 	assert.Equal(t, "Invalid parameters", response["error"]["message"])
// }
