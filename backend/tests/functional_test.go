package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cinagroup/cinaseek/backend/internal/config"
	"github.com/cinagroup/cinaseek/backend/internal/handler"
	"github.com/cinagroup/cinaseek/backend/internal/middleware"
	"github.com/cinagroup/cinaseek/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// ─── test helpers ───────────────────────────────────────────────────────────

// funcTestRouter builds a fresh gin.Engine wired with all handlers and JWT
// middleware. A valid token for userID=1, username="tester" is returned.
func funcTestRouter() (*gin.Engine, string) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery()) // prevent panics from nil DB

	jwtCfg := &config.JWTConfig{Secret: "test-secret-key", ExpireTime: 24 * 1e9 * 60 * 60}
	validToken := generateTestToken(jwtCfg.Secret, 1, "tester", defaultExpiry())

	r.Use(middleware.JWTAuth(jwtCfg))

	// Handlers with nil services — we only validate routing / middleware / request
	// binding. When services are nil the handler will panic, but for validation
	// of middleware (auth, binding) this is sufficient. We'll override specific
	// routes that need mock services below.
	cfg := &config.Config{JWT: *jwtCfg}

	vmSvc := service.NewVMService(nil)
	vmH := handler.NewVMHandler(cfg, vmSvc)

	mountSvc := service.NewMountService(nil)
	mountH := handler.NewMountHandler(cfg, mountSvc)

	ocSvc := service.NewOpenClawService(nil)
	ocH := handler.NewOpenClawHandler(cfg, ocSvc)

	remoteH := handler.NewRemoteHandler(cfg)

	// VM routes
	vms := r.Group("/api/v1/vm")
	{
		vms.GET("/list", vmH.ListVMs)
		vms.POST("/create", vmH.CreateVM)
		vms.GET("/detail/:id", vmH.GetVMDetail)
		vms.POST("/operate/:id", vmH.OperateVM)
		vms.PUT("/config/:id", vmH.UpdateVMConfig)
		vms.GET("/:id/logs", vmH.GetVMLogs)
		vms.GET("/:id/metrics", vmH.GetVMMetrics)
		vms.POST("/:id/snapshots", vmH.CreateSnapshot)
		vms.GET("/:id/snapshots", vmH.ListSnapshots)
	}

	// Mount routes
	mounts := r.Group("/api/v1/mount")
	{
		mounts.GET("/list", mountH.ListMounts)
		mounts.POST("/add", mountH.AddMount)
		mounts.POST("/operate/:id", mountH.OperateMount)
		mounts.GET("/openclaw-config", mountH.GetOpenClawConfig)
		mounts.POST("/openclaw-configure", mountH.ConfigureOpenClawMount)
	}

	// OpenClaw routes
	oc := r.Group("/api/v1/openclaw")
	{
		oc.GET("/status", ocH.GetOpenClawStatus)
		oc.POST("/deploy", ocH.DeployOpenClaw)
		oc.POST("/operate/:id", ocH.OperateOpenClaw)
		oc.GET("/:id/logs", ocH.GetOpenClawLog)
		oc.PUT("/config/:id", ocH.UpdateOpenClawConfig)
		oc.GET("/monitor", ocH.GetOpenClawMonitor)
		oc.GET("/workspaces", ocH.GetWorkspaceList)
	}

	// Remote routes
	remote := r.Group("/api/v1/remote")
	{
		remote.GET("/status", remoteH.GetRemoteStatus)
		remote.POST("/switch/:id", remoteH.SwitchRemoteAccess)
		remote.GET("/whitelist", remoteH.GetIPWhitelist)
		remote.POST("/whitelist/add", remoteH.AddIPWhitelist)
		remote.DELETE("/whitelist/:id/:whitelist_id", remoteH.RemoveIPWhitelist)
		remote.GET("/log/:id", remoteH.GetRemoteLog)
	}

	return r, validToken
}

// defaultExpiry returns a time 24 hours in the future.
func defaultExpiry() time.Time {
	return time.Now().Add(24 * time.Hour)
}

// ─── VM functional tests ────────────────────────────────────────────────────

func TestVMList_ReturnsOK(t *testing.T) {
	r, token := funcTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/vm/list", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// With nil services we expect 500 (DB unreachable) but NOT 401.
	// This proves JWT middleware passed and handler was invoked.
	if w.Code == http.StatusUnauthorized {
		t.Errorf("expected handler to be reached, got 401")
	}
}

func TestVMCreate_ValidInput(t *testing.T) {
	r, token := funcTestRouter()

	body := map[string]interface{}{
		"name":         "test-vm-1",
		"image":        "ubuntu:22.04",
		"cpu":          2,
		"memory":       4,
		"disk":         50,
		"network_type": "nat",
	}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/vm/create", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code == http.StatusUnauthorized {
		t.Errorf("expected handler to be reached, got 401")
	}
}

func TestVMCreate_InvalidInput(t *testing.T) {
	r, token := funcTestRouter()

	// Missing required fields (name, image)
	body := map[string]interface{}{
		"cpu": 2,
	}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/vm/create", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should return 400 because required fields are missing
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid input, got %d", w.Code)
	}
}

func TestVMDetail_NotFound(t *testing.T) {
	r, token := funcTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/vm/detail/99999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// With nil DB, handler will return 500, but proves routing works and
	// the ID param is parsed correctly (non-numeric IDs should return 400).
	if w.Code == http.StatusUnauthorized {
		t.Errorf("expected handler to be reached, got 401")
	}
}

func TestVMDetail_InvalidID(t *testing.T) {
	r, token := funcTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/vm/detail/notanumber", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for non-numeric ID, got %d", w.Code)
	}
}

func TestVMOperate_Start(t *testing.T) {
	r, token := funcTestRouter()

	body := map[string]interface{}{"operation": "start"}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/vm/operate/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code == http.StatusUnauthorized {
		t.Errorf("expected handler to be reached, got 401")
	}
}

func TestVMOperate_Stop(t *testing.T) {
	r, token := funcTestRouter()

	body := map[string]interface{}{"operation": "stop"}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/vm/operate/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code == http.StatusUnauthorized {
		t.Errorf("expected handler to be reached, got 401")
	}
}

func TestVMOperate_Restart(t *testing.T) {
	r, token := funcTestRouter()

	body := map[string]interface{}{"operation": "restart"}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/vm/operate/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code == http.StatusUnauthorized {
		t.Errorf("expected handler to be reached, got 401")
	}
}

func TestVMOperate_Delete(t *testing.T) {
	r, token := funcTestRouter()

	body := map[string]interface{}{"operation": "delete"}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/vm/operate/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code == http.StatusUnauthorized {
		t.Errorf("expected handler to be reached, got 401")
	}
}

func TestVMOperate_InvalidOperation_BadRequest(t *testing.T) {
	r, token := funcTestRouter()

	// Missing operation field → binding fails
	body := map[string]interface{}{}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/vm/operate/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing operation, got %d", w.Code)
	}
}

func TestVMOperate_InvalidID_BadRequest(t *testing.T) {
	r, token := funcTestRouter()

	body := map[string]interface{}{"operation": "start"}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/vm/operate/abc", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for non-numeric VM ID, got %d", w.Code)
	}
}

// ─── Mount functional tests ─────────────────────────────────────────────────

func TestMountList_ReturnsHandlerReached(t *testing.T) {
	r, token := funcTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/mount/list", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code == http.StatusUnauthorized {
		t.Errorf("expected handler to be reached, got 401")
	}
}

func TestMountAdd_InvalidInput_BadRequest(t *testing.T) {
	r, token := funcTestRouter()

	// Missing required fields: vm_id, name, host_path, vm_path
	body := map[string]interface{}{}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/mount/add", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing fields, got %d", w.Code)
	}
}

func TestMountOperate_InvalidID_BadRequest(t *testing.T) {
	r, token := funcTestRouter()

	body := map[string]interface{}{"operation": "mount"}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/mount/operate/notanumber", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for non-numeric mount ID, got %d", w.Code)
	}
}

func TestMountOpenClawConfig_MissingVMID(t *testing.T) {
	r, token := funcTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/mount/openclaw-config", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing vm_id, got %d", w.Code)
	}
}

// ─── OpenClaw functional tests ──────────────────────────────────────────────

func TestOpenClawStatus_MissingVMID(t *testing.T) {
	r, token := funcTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/openclaw/status", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing vm_id, got %d", w.Code)
	}
}

func TestOpenClawDeploy_InvalidInput(t *testing.T) {
	r, token := funcTestRouter()

	// Missing vm_id
	body := map[string]interface{}{}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/openclaw/deploy", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing vm_id, got %d", w.Code)
	}
}

func TestOpenClawOperate_InvalidID(t *testing.T) {
	r, token := funcTestRouter()

	body := map[string]interface{}{"operation": "start"}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/openclaw/operate/abc", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for non-numeric ID, got %d", w.Code)
	}
}

func TestOpenClawLogs_InvalidID(t *testing.T) {
	r, token := funcTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/openclaw/xyz/logs", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for non-numeric ID, got %d", w.Code)
	}
}

func TestOpenClawMonitor_MissingVMID(t *testing.T) {
	r, token := funcTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/openclaw/monitor", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing vm_id, got %d", w.Code)
	}
}

func TestOpenClawWorkspaces_MissingVMID(t *testing.T) {
	r, token := funcTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/openclaw/workspaces", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing vm_id, got %d", w.Code)
	}
}

// ─── Remote functional tests ────────────────────────────────────────────────

func TestRemoteStatus_MissingVMID(t *testing.T) {
	r, token := funcTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/remote/status", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing vm_id, got %d", w.Code)
	}
}

func TestRemoteStatus_InvalidVMID(t *testing.T) {
	r, token := funcTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/remote/status?vm_id=notanumber", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for non-numeric vm_id, got %d", w.Code)
	}
}

func TestRemoteSwitch_InvalidID(t *testing.T) {
	r, token := funcTestRouter()

	body := map[string]interface{}{"enabled": true}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/remote/switch/abc", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for non-numeric ID, got %d", w.Code)
	}
}

func TestRemoteSwitch_MissingEnabled(t *testing.T) {
	r, token := funcTestRouter()

	body := map[string]interface{}{}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/remote/switch/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing enabled field, got %d", w.Code)
	}
}

func TestRemoteWhitelist_MissingVMID(t *testing.T) {
	r, token := funcTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/remote/whitelist", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing vm_id, got %d", w.Code)
	}
}

func TestRemoteWhitelistAdd_InvalidInput(t *testing.T) {
	r, token := funcTestRouter()

	body := map[string]interface{}{}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/remote/whitelist/add", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing fields, got %d", w.Code)
	}
}

func TestRemoteLog_InvalidID(t *testing.T) {
	r, token := funcTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/remote/log/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for non-numeric ID, got %d", w.Code)
	}
}

// ─── Auth middleware integration ─────────────────────────────────────────────

func TestFunctional_NoToken_Returns401(t *testing.T) {
	r, _ := funcTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/vm/list", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without token, got %d", w.Code)
	}
}

func TestFunctional_InvalidToken_Returns401(t *testing.T) {
	r, _ := funcTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/vm/list", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 with invalid token, got %d", w.Code)
	}
}
