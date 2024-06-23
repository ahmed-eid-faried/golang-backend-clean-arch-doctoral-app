package http

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"
//
//

// 	"github.com/gin-gonic/gin"
// 	"github.com/quangdangfit/gocommon/logger"
// 	"github.com/quangdangfit/gocommon/validation"

// 	// orderModel "main/internal/order/model"
// 	// productModel "main/internal/product/model"
// 	httpServer "main/internal/server/http"
// 	"main/internal/user/dto"
// 	userModel "main/internal/user/model"
// 	"main/pkg/config"
// 	"main/pkg/dbs"
// 	"main/pkg/redis"
// 	"main/pkg/utils"
// )

// // Global variables for the test router, test database, and test cache
// var (
// 	testRouter *gin.Engine
// 	dbTest     dbs.IDatabase
// 	testCache  redis.IRedis
// )

// // TestMain is the entry point for the test suite. It sets the Gin mode to test,
// // sets up the test environment, runs the tests, and then tears down the test
// // environment.
// func TestMain(m *testing.M) {
// 	gin.SetMode(gin.TestMode)
// 	setup()
// 	exitCode := m.Run()
// 	teardown()

// 	os.Exit(exitCode)
// }

// // setup initializes the test environment, including the database, cache, and
// // HTTP server.
// func setup() {
// 	// Load the configuration from the environment
// 	cfg := config.LoadConfig()

// 	// Initialize the logger
// 	logger.Initialize(config.ProductionEnv)

// 	// Connect to the database
// 	var err error
// 	dbTest, err = dbs.NewDatabase(cfg.DatabaseURI)
// 	if err != nil {
// 		logger.Fatal("Cannot connect to database", err)
// 	}

// 	// Perform database migration
// 	err = dbTest.AutoMigrate(&userModel.User{})
// 	if err != nil {
// 		logger.Fatal("Database migration fail", err)
// 	}

// 	// Initialize the validator
// 	validator := validation.New()

// 	// Initialize the cache
// 	testCache = redis.New(redis.Config{
// 		Address:  cfg.RedisURI,
// 		Password: cfg.RedisPassword,
// 		Database: cfg.RedisDB,
// 	})

// 	// Initialize the HTTP server
// 	server := httpServer.NewServer(validator, dbTest, testCache)
// 	_ = server.MapRoutes()
// 	testRouter = server.GetEngine()

// 	// Create a test user
// 	dbTest.Create(context.Background(), &userModel.User{
// 		Email:    "test@test.com",
// 		Password: "test123456",
// 	})
// }

// // teardown cleans up the test environment, including dropping the test database
// // table.
// func teardown() {
// 	migrator := dbTest.GetDB().Migrator()
// 	migrator.DropTable(&userModel.User{})
// }

// // makeRequest creates and sends an HTTP request to the test router, and returns
// // the response.
// func makeRequest(method, url string, body interface{}, token string) *httptest.ResponseRecorder {
// 	requestBody, _ := json.Marshal(body)
// 	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
// 	if token != "" {
// 		request.Header.Add("Authorization", "Bearer "+token)
// 	}
// 	writer := httptest.NewRecorder()
// 	testRouter.ServeHTTP(writer, request)
// 	return writer
// }

// // accessToken returns the access token for the test user.
// func accessToken() string {
// 	user := dto.LoginReq{
// 		Email:    "test@test.com",
// 		Password: "test123456",
// 	}

// 	writer := makeRequest("POST", "/auth/login", user, "")
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	return response["result"]["access_token"]
// }

// // refreshToken returns the refresh token for the test user.
// func refreshToken() string {
// 	user := dto.LoginReq{
// 		Email:    "test@test.com",
// 		Password: "test123456",
// 	}

// 	writer := makeRequest("POST", "/auth/login", user, "")
// 	var response map[string]map[string]string
// 	_ = json.Unmarshal(writer.Body.Bytes(), &response)
// 	return response["result"]["refresh_token"]
// }

// // parseResponseResult parses the "result" field from the given response data
// // and copies it to the given result object.
// func parseResponseResult(resData []byte, result interface{}) {
// 	var response map[string]interface{}
// 	_ = json.Unmarshal(resData, &response)
// 	utils.Copy(result, response["result"])
// }

// // cleanData deletes the given records from the database and removes any
// // matching keys from the cache.
// func cleanData(records ...interface{}) {
// 	// 	// dbTest.GetDB().Where("1 = 1").Delete(&orderModel.OrderLine{})
// 	// 	// dbTest.GetDB().Where("1 = 1").Delete(&productModel.Product{})
// 	// 	// dbTest.GetDB().Where("1 = 1").Delete(&orderModel.Order{})
// 	for _, record := range records {
// 		dbTest.Delete(context.Background(), record)
// 	}

// 	testCache.RemovePattern("*")
// }
