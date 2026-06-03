package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"

	"billsoftware/backend/internal/model"
	"billsoftware/backend/internal/response"
)

type AdminApprovalHandler struct {
	engine *xorm.Engine
}

type approvalRequest struct {
	Remark string `json:"remark"`
}

type batchApprovalRequest struct {
	UserIDs []uint64 `json:"user_ids"`
	Remark  string   `json:"remark"`
}

func NewAdminApprovalHandler(engine *xorm.Engine) *AdminApprovalHandler {
	return &AdminApprovalHandler{engine: engine}
}

func (h *AdminApprovalHandler) ListPendingUsers(c *gin.Context) {
	users := make([]model.User, 0)
	if err := h.engine.Where("approval_status = ?", model.UserApprovalPending).Desc("id").Find(&users); err != nil {
		response.Fail(c, http.StatusInternalServerError, "query pending approval users failed")
		return
	}

	list := make([]gin.H, 0, len(users))
	for _, user := range users {
		list = append(list, buildApprovalUser(user))
	}

	response.Success(c, gin.H{
		"list": list,
		"summary": gin.H{
			"pending_count": len(list),
		},
	})
}

func (h *AdminApprovalHandler) ApproveUser(c *gin.Context) {
	userID, err := parseApprovalUserID(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.updateApproval(userID, model.UserApprovalApproved, ""); err != nil {
		writeApprovalError(c, err, "approve user failed")
		return
	}

	response.Success(c, gin.H{
		"user_id":         userID,
		"status":          1,
		"approval_status": model.UserApprovalApproved,
	})
}

func (h *AdminApprovalHandler) RejectUser(c *gin.Context) {
	userID, err := parseApprovalUserID(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req approvalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid reject payload")
		return
	}

	remark := strings.TrimSpace(req.Remark)
	if err := h.updateApproval(userID, model.UserApprovalRejected, remark); err != nil {
		writeApprovalError(c, err, "reject user failed")
		return
	}

	response.Success(c, gin.H{
		"user_id":         userID,
		"status":          0,
		"approval_status": model.UserApprovalRejected,
		"approval_remark": remark,
	})
}

func (h *AdminApprovalHandler) BatchApproveUsers(c *gin.Context) {
	var req batchApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid batch approve payload")
		return
	}

	userIDs := uniqueApprovalUserIDs(req.UserIDs)
	if len(userIDs) == 0 {
		response.Fail(c, http.StatusBadRequest, "user_ids is required")
		return
	}

	affected, err := h.batchUpdateApproval(userIDs, model.UserApprovalApproved, "")
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "batch approve users failed")
		return
	}

	response.Success(c, gin.H{
		"user_ids":        userIDs,
		"affected":        affected,
		"status":          1,
		"approval_status": model.UserApprovalApproved,
	})
}

func (h *AdminApprovalHandler) BatchRejectUsers(c *gin.Context) {
	var req batchApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid batch reject payload")
		return
	}

	userIDs := uniqueApprovalUserIDs(req.UserIDs)
	if len(userIDs) == 0 {
		response.Fail(c, http.StatusBadRequest, "user_ids is required")
		return
	}

	remark := strings.TrimSpace(req.Remark)
	affected, err := h.batchUpdateApproval(userIDs, model.UserApprovalRejected, remark)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "batch reject users failed")
		return
	}

	response.Success(c, gin.H{
		"user_ids":        userIDs,
		"affected":        affected,
		"status":          0,
		"approval_status": model.UserApprovalRejected,
		"approval_remark": remark,
	})
}

func (h *AdminApprovalHandler) updateApproval(userID uint64, approvalStatus string, remark string) error {
	session := h.engine.NewSession()
	defer func() {
		_ = session.Close()
	}()

	if err := session.Begin(); err != nil {
		return err
	}

	user := &model.User{}
	has, err := session.ID(userID).Get(user)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	if !has {
		_ = session.Rollback()
		return errApprovalUserNotFound
	}

	now := time.Now()
	update := &model.User{
		Status:            approvalStatusToUserStatus(approvalStatus),
		ApprovalStatus:    approvalStatus,
		ApprovalUpdatedAt: &now,
		ApprovalRemark:    strings.TrimSpace(remark),
	}

	if _, err := session.ID(userID).Cols("status", "approval_status", "approval_updated_at", "approval_remark").Update(update); err != nil {
		_ = session.Rollback()
		return err
	}

	if approvalStatus != model.UserApprovalApproved {
		if _, err := session.Where("user_id = ? AND is_active = ?", userID, 1).Cols("is_active").Update(&model.UserSession{
			IsActive: 0,
		}); err != nil {
			_ = session.Rollback()
			return err
		}
	}

	if err := session.Commit(); err != nil {
		_ = session.Rollback()
		return err
	}

	return nil
}

func (h *AdminApprovalHandler) batchUpdateApproval(userIDs []uint64, approvalStatus string, remark string) (int64, error) {
	session := h.engine.NewSession()
	defer func() {
		_ = session.Close()
	}()

	if err := session.Begin(); err != nil {
		return 0, err
	}

	now := time.Now()
	update := &model.User{
		Status:            approvalStatusToUserStatus(approvalStatus),
		ApprovalStatus:    approvalStatus,
		ApprovalUpdatedAt: &now,
		ApprovalRemark:    strings.TrimSpace(remark),
	}

	affected, err := session.In("id", userIDs).Where("approval_status = ?", model.UserApprovalPending).Cols(
		"status",
		"approval_status",
		"approval_updated_at",
		"approval_remark",
	).Update(update)
	if err != nil {
		_ = session.Rollback()
		return 0, err
	}

	if approvalStatus != model.UserApprovalApproved {
		if _, err := session.In("user_id", userIDs).Where("is_active = ?", 1).Cols("is_active").Update(&model.UserSession{
			IsActive: 0,
		}); err != nil {
			_ = session.Rollback()
			return 0, err
		}
	}

	if err := session.Commit(); err != nil {
		_ = session.Rollback()
		return 0, err
	}

	return affected, nil
}

func buildApprovalUser(user model.User) gin.H {
	return gin.H{
		"id":                  user.ID,
		"username":            user.Username,
		"nickname":            user.Nickname,
		"phone":               user.Phone,
		"email":               user.Email,
		"status":              user.Status,
		"approval_status":     user.ApprovalStatus,
		"approval_updated_at": user.ApprovalUpdatedAt,
		"approval_remark":     user.ApprovalRemark,
		"application_purpose": "注册后等待管理员审批",
		"created_at":          user.CreatedAt,
		"updated_at":          user.UpdatedAt,
	}
}

func approvalStatusToUserStatus(approvalStatus string) int {
	if approvalStatus == model.UserApprovalApproved {
		return 1
	}

	return 0
}

func parseApprovalUserID(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}

func uniqueApprovalUserIDs(userIDs []uint64) []uint64 {
	seen := make(map[uint64]bool)
	result := make([]uint64, 0, len(userIDs))

	for _, userID := range userIDs {
		if userID == 0 || seen[userID] {
			continue
		}

		seen[userID] = true
		result = append(result, userID)
	}

	return result
}

func writeApprovalError(c *gin.Context, err error, fallback string) {
	if errors.Is(err, errApprovalUserNotFound) {
		response.Fail(c, http.StatusNotFound, "user not found")
		return
	}

	response.Fail(c, http.StatusInternalServerError, fallback)
}

var errApprovalUserNotFound = errors.New("approval user not found")
