package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cinagroup/cinaseek/backend/internal/cinaclaw"
	"github.com/cinagroup/cinaseek/backend/internal/config"
	"github.com/cinagroup/cinaseek/backend/internal/handler"
	"github.com/cinagroup/cinaseek/backend/internal/middleware"
	"github.com/cinagroup/cinaseek/backend/internal/repository"
	"github.com/cinagroup/cinaseek/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 解析命令行参数
	var configPath string
	flag.StringVar(&configPath, "config", "", "配置文件路径")
	flag.Parse()

	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	if err := repository.InitDB(&cfg.Database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 创建 Gin 引擎
	r := gin.New()

	// 使用中间件
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(&cfg.CORS))
	r.Use(gin.Recovery())

	// 创建 gRPC 客户端管理器
	clientMgr := cinaclaw.NewClientManager("/var/run/cinaclaw.sock")

	// 创建服务层
	vmService := service.NewVMService(clientMgr)
	mountService := service.NewMountService(clientMgr)
	openclawService := service.NewOpenClawService(clientMgr)

	// 创建处理器
	authHandler := handler.NewAuthHandler(cfg)
	oauthHandler := handler.NewOAuthHandler(cfg)
	vmHandler := handler.NewVMHandler(cfg, vmService)
	mountHandler := handler.NewMountHandler(cfg, mountService)
	openclawHandler := handler.NewOpenClawHandler(cfg, openclawService)
	remoteHandler := handler.NewRemoteHandler(cfg)
	systemHandler := handler.NewSystemHandler(cfg)
	adminHandler := handler.NewAdminHandler(cfg)

	// 注册路由
	registerRoutes(r, authHandler, oauthHandler, vmHandler, mountHandler, openclawHandler, remoteHandler, systemHandler, adminHandler, cfg)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 优雅关闭
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}

		log.Println("Server exited")
	}()

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func registerRoutes(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	oauthHandler *handler.OAuthHandler,
	vmHandler *handler.VMHandler,
	mountHandler *handler.MountHandler,
	openclawHandler *handler.OpenClawHandler,
	remoteHandler *handler.RemoteHandler,
	systemHandler *handler.SystemHandler,
	adminHandler *handler.AdminHandler,
	cfg *config.Config,
) {
	// 健康检查
	r.GET("/health", systemHandler.HealthCheck)
	r.GET("/api/version", systemHandler.GetSystemVersion)

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// OAuth 认证模块（无需登录）
		oauth := v1.Group("/oauth")
		{
			oauth.GET("/authorize", oauthHandler.OAuthRedirect)
			oauth.GET("/callback", oauthHandler.OAuthCallback)
			oauth.GET("/providers", oauthHandler.GetOAuthProviders)
			oauth.POST("/logout", oauthHandler.Logout)
		}

		// 传统认证模块（兼容模式，可选）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/reset-pwd", authHandler.ResetPassword)
		}

		// 需要认证的路由
		authorized := v1.Group("")
		authorized.Use(middleware.JWTAuth(&cfg.JWT))
		{
			// 认证模块
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/user-info", authHandler.GetUserInfo)
			auth.PUT("/user-info", authHandler.UpdateUserInfo)
			auth.PUT("/user-pwd", authHandler.UpdatePassword)
			auth.GET("/login-logs", authHandler.GetLoginLogs)
			auth.GET("/sessions", authHandler.GetSessions)
			auth.POST("/sessions/revoke", authHandler.RevokeSession)

			// 虚拟机管理模块
			vms := authorized.Group("/vm")
			{
				vms.GET("/list", vmHandler.ListVMs)
				vms.GET("/detail/:id", vmHandler.GetVMDetail)
				vms.POST("/create", vmHandler.CreateVM)
				vms.POST("/operate/:id", vmHandler.OperateVM)
				vms.PUT("/update-config/:id", vmHandler.UpdateVMConfig)
				vms.GET("/snapshots/:id", vmHandler.ListSnapshots)
				vms.POST("/snapshot/:id", vmHandler.CreateSnapshot)
				vms.POST("/snapshot/:id/restore", vmHandler.RestoreSnapshot)
				vms.DELETE("/snapshot/:id/:snapshot_id", vmHandler.DeleteSnapshot)
				vms.GET("/logs/:id", vmHandler.GetVMLogs)
				vms.GET("/metrics/:id", vmHandler.GetVMMetrics)
			}

			// 目录挂载模块
			mounts := authorized.Group("/mount")
			{
				mounts.GET("/list", mountHandler.ListMounts)
				mounts.POST("/add", mountHandler.AddMount)
				mounts.POST("/operate/:id", mountHandler.OperateMount)
				mounts.GET("/openclaw-config", mountHandler.GetOpenClawConfig)
				mounts.POST("/openclaw-config", mountHandler.ConfigureOpenClawMount)
			}

			// OpenClaw 管理模块
			openclaw := authorized.Group("/openclaw")
			{
				openclaw.GET("/status", openclawHandler.GetOpenClawStatus)
				openclaw.POST("/deploy", openclawHandler.DeployOpenClaw)
				openclaw.POST("/operate/:id", openclawHandler.OperateOpenClaw)
				openclaw.GET("/log/:id", openclawHandler.GetOpenClawLog)
				openclaw.PUT("/config/:id", openclawHandler.UpdateOpenClawConfig)
				openclaw.GET("/monitor", openclawHandler.GetOpenClawMonitor)
				openclaw.GET("/workspace", openclawHandler.GetWorkspaceList)
			}

			// 远程访问模块
			remote := authorized.Group("/remote")
			{
				remote.GET("/status", remoteHandler.GetRemoteStatus)
				remote.PUT("/switch/:id", remoteHandler.SwitchRemoteAccess)
				remote.GET("/ip-whitelist", remoteHandler.GetIPWhitelist)
				remote.POST("/ip-whitelist", remoteHandler.AddIPWhitelist)
				remote.DELETE("/ip-whitelist/:id/:whitelist_id", remoteHandler.RemoveIPWhitelist)
				remote.GET("/log/:id", remoteHandler.GetRemoteLog)
			}

			// 用户中心模块
			user := authorized.Group("/user")
			{
				user.PUT("/update-info", authHandler.UpdateUserInfo)
				user.PUT("/update-pwd", authHandler.UpdatePassword)
			}

			// 系统模块
			system := authorized.Group("/system")
			{
				system.GET("/setting", systemHandler.GetSystemSetting)
				system.PUT("/setting", middleware.AdminRequired(), systemHandler.UpdateSystemSetting)
				system.GET("/version", systemHandler.GetSystemVersion)
				system.GET("/dashboard", systemHandler.GetDashboard)
				system.GET("/statistics", systemHandler.GetStatistics)
				system.GET("/search", systemHandler.SearchVMs)
				system.POST("/batch-vm", systemHandler.BatchOperateVMs)
			}

			// 管理员模块
			admin := authorized.Group("/admin")
			admin.Use(middleware.AdminRequired())
			{
				admin.GET("/users", adminHandler.ListUsers)
				admin.GET("/users/:id", adminHandler.GetUser)
				admin.GET("/stats", adminHandler.GetSystemStats)

				// Root 专属路由
				root := admin.Group("")
				root.Use(middleware.RootRequired())
				{
					root.PUT("/users/:id/role", adminHandler.UpdateUserRole)
					root.DELETE("/users/:id", adminHandler.DeleteUser)
				}
			}
		}
	}
}
