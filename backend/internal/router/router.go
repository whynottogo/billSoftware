package router

import (
	"github.com/gin-gonic/gin"
	"xorm.io/xorm"

	"billsoftware/backend/internal/config"
	"billsoftware/backend/internal/handler"
	"billsoftware/backend/internal/middleware"
	"billsoftware/backend/internal/storage"
)

func NewHTTPRouter(cfg *config.AppConfig, engine *xorm.Engine, objectStorage *storage.ObjectStorage) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	httpRouter := gin.New()
	httpRouter.Use(gin.Logger())
	httpRouter.Use(gin.Recovery())

	healthHandler := handler.NewHealthHandler(engine)
	userAuthHandler := handler.NewUserAuthHandler(engine)
	adminAuthHandler := handler.NewAdminAuthHandler()
	userHandler := handler.NewUserHandler(engine)
	userFileHandler := handler.NewUserFileHandler(engine, objectStorage, cfg.MinIO)
	userLedgerHandler := handler.NewUserLedgerHandler(engine, objectStorage)
	userBillHandler := handler.NewUserBillHandler(engine)
	userBudgetHandler := handler.NewUserBudgetHandler(engine)
	userChartHandler := handler.NewUserChartHandler(engine)
	userAssetHandler := handler.NewUserAssetHandler(engine)
	userProfileHandler := handler.NewUserProfileHandler(engine)
	userFamilyHandler := handler.NewUserFamilyHandler(engine)
	adminUserBillHandler := handler.NewAdminUserBillHandler(engine)
	adminDashboardHandler := handler.NewAdminDashboardHandler(engine)
	adminApprovalHandler := handler.NewAdminApprovalHandler(engine)
	adminFamilyHandler := handler.NewAdminFamilyHandler(engine)

	httpRouter.GET("/api/health", healthHandler.Ping)

	userGroup := httpRouter.Group("/api/user")
	{
		authGroup := userGroup.Group("/auth")
		authGroup.POST("/register", userAuthHandler.Register)
		authGroup.POST("/login", userAuthHandler.Login)
	}

	adminGroup := httpRouter.Group("/api/admin")
	{
		authGroup := adminGroup.Group("/auth")
		authGroup.POST("/login", adminAuthHandler.Login)

		protectedGroup := adminGroup.Group("")
		protectedGroup.Use(middleware.AdminAuthRequired())
		protectedGroup.GET("/users", userHandler.List)
		protectedGroup.PUT("/users/:id/status", userHandler.ChangeStatus)
		protectedGroup.GET("/users/:id/summary", userHandler.Summary)
		protectedGroup.GET("/users/:id/bills/overview", adminUserBillHandler.GetOverview)
		protectedGroup.GET("/dashboard/overview", adminDashboardHandler.Overview)
		protectedGroup.GET("/approvals/users", adminApprovalHandler.ListPendingUsers)
		protectedGroup.POST("/approvals/users/batch-approve", adminApprovalHandler.BatchApproveUsers)
		protectedGroup.POST("/approvals/users/batch-reject", adminApprovalHandler.BatchRejectUsers)
		protectedGroup.POST("/approvals/users/:id/approve", adminApprovalHandler.ApproveUser)
		protectedGroup.POST("/approvals/users/:id/reject", adminApprovalHandler.RejectUser)
		protectedGroup.GET("/families/summary", adminFamilyHandler.Summary)
		protectedGroup.GET("/families", adminFamilyHandler.List)
		protectedGroup.GET("/families/:familyId", adminFamilyHandler.Detail)
	}

	protectedUserGroup := userGroup.Group("")
	protectedUserGroup.Use(middleware.UserAuthRequired(engine))
	protectedUserGroup.GET("/ledger", userLedgerHandler.GetLedger)
	protectedUserGroup.POST("/ledger", userLedgerHandler.CreateLedger)
	protectedUserGroup.DELETE("/ledger/:id", userLedgerHandler.DeleteLedger)
	protectedUserGroup.POST("/files", userFileHandler.Upload)
	protectedUserGroup.GET("/categories", userLedgerHandler.ListCategories)
	protectedUserGroup.POST("/categories", userLedgerHandler.CreateCategory)
	protectedUserGroup.DELETE("/categories/:id", userLedgerHandler.DeleteCategory)
	protectedUserGroup.GET("/bills/years", userBillHandler.ListYears)
	protectedUserGroup.GET("/bills/year/:year", userBillHandler.GetYearDetail)
	protectedUserGroup.GET("/bills/month/:month", userBillHandler.GetMonthDetail)
	protectedUserGroup.GET("/budgets/month/current", userBudgetHandler.GetCurrentMonth)
	protectedUserGroup.PUT("/budgets/month/current", userBudgetHandler.UpdateCurrentMonth)
	protectedUserGroup.GET("/budgets/year/options", userBudgetHandler.ListYearOptions)
	protectedUserGroup.GET("/budgets/year/:year", userBudgetHandler.GetYear)
	protectedUserGroup.PUT("/budgets/year/current", userBudgetHandler.UpdateCurrentYear)
	protectedUserGroup.GET("/charts/years", userChartHandler.ListYears)
	protectedUserGroup.GET("/charts/expense/:year", userChartHandler.GetExpenseYear)
	protectedUserGroup.GET("/charts/income/:year", userChartHandler.GetIncomeYear)
	protectedUserGroup.GET("/assets", userAssetHandler.List)
	protectedUserGroup.POST("/assets", userAssetHandler.Create)
	protectedUserGroup.PUT("/assets/:id", userAssetHandler.Update)
	protectedUserGroup.GET("/assets/:id", userAssetHandler.Detail)
	protectedUserGroup.POST("/assets/:id/operations", userAssetHandler.CreateOperation)
	protectedUserGroup.GET("/profile", userProfileHandler.GetProfile)
	protectedUserGroup.PUT("/profile", userProfileHandler.UpdateProfile)
	protectedUserGroup.PUT("/profile/password", userProfileHandler.UpdatePassword)
	protectedUserGroup.GET("/families", userFamilyHandler.List)
	protectedUserGroup.POST("/families", userFamilyHandler.Create)
	protectedUserGroup.POST("/families/join", userFamilyHandler.Join)
	protectedUserGroup.POST("/families/join-by-link", userFamilyHandler.JoinByLink)
	protectedUserGroup.POST("/families/:familyId/leave", userFamilyHandler.Leave)
	protectedUserGroup.GET("/families/:familyId", userFamilyHandler.GetDetail)
	protectedUserGroup.GET("/families/:familyId/member-share", userFamilyHandler.GetMemberShare)

	return httpRouter
}
