package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cinagroup/cinaseek/backend/internal/cinaclaw"
	"github.com/cinagroup/cinaseek/backend/internal/config"
	"github.com/cinagroup/cinaseek/backend/internal/handler"
	"github.com/cinagroup/cinaseek/backend/internal/middleware"
	"github.com/cinagroup/cinaseek/backend/internal/service"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
	token  string
	cfg    *config.Config
)

// TestMain 测试初始化
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	// 加载配置
	cfg = config.Load()

	// 创建路由器
	router = gin.New()
	router.Use(middleware.Logger())
	router.Use(middleware.CORS(&cfg.CORS))

	// 创建处理器
	authHandler := handler.NewAuthHandler(cfg)
	clientMgr := cinaclaw.NewClientManager("/var/run/cinaclaw.sock")
	vmService := service.NewVMService(clientMgr)
	vmHandler := handler.NewVMHandler(cfg, vmService)

	// 注册测试路由
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	vms := router.Group("/api/v1/vm")
	vms.Use(middleware.JWTAuth(&cfg.JWT))
	{
		vms.GET("/list", vmHandler.ListVMs)
		vms.POST("/create", vmHandler.CreateVM)
	}

	os.Exit(m.Run())
}

// skipIfNoDB skips tests that require a database connection.
func skipIfNoDB(t *testing.T) {
	t.Helper()
	// Skip DB-dependent tests if no database is available
	if os.Getenv("DB_HOST") == "" && os.Getenv("TEST_WITH_DB") == "" {
		t.Skip("Skipping: database not available (set TEST_WITH_DB=1 to run)")
	}
}

// TestRegister 测试用户注册
func TestRegister(t *testing.T) {
	skipIfNoDB(t)

	data := map[string]interface{}{
		"username":         "testuser",
		"email":            "test@example.com",
		"password":         "Test1234",
		"confirm_password": "Test1234",
	}

	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"].(float64) != 0 {
		t.Errorf("Expected code 0, got %v", response["code"])
	}
}

// TestLogin 测试用户登录
func TestLogin(t *testing.T) {
	skipIfNoDB(t)

	data := map[string]interface{}{
		"username": "testuser",
		"password": "Test1234",
	}

	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"].(float64) != 0 {
		t.Errorf("Expected code 0, got %v", response["code"])
	}

	// 保存 token 供其他测试使用
	if data, ok := response["data"].(map[string]interface{}); ok {
		token = data["token"].(string)
	}
}

// TestListVMs 测试获取虚拟机列表
func TestListVMs(t *testing.T) {
	skipIfNoDB(t)

	req, _ := http.NewRequest("GET", "/api/v1/vm/list", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"].(float64) != 0 {
		t.Errorf("Expected code 0, got %v", response["code"])
	}
}

// TestCreateVM 测试创建虚拟机
func TestCreateVM(t *testing.T) {
	skipIfNoDB(t)

	data := map[string]interface{}{
		"name":        "test-vm",
		"image":       "ubuntu:22.04",
		"cpu":         2,
		"memory":      4,
		"disk":        50,
		"network_type": "nat",
	}

	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/api/v1/vm/create", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"].(float64) != 0 {
		t.Errorf("Expected code 0, got %v", response["code"])
	}
}

// TestUnauthorized 测试未授权访问
func TestUnauthorized(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/vm/list", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}
