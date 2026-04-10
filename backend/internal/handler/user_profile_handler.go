package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"

	"billsoftware/backend/internal/model"
	"billsoftware/backend/internal/response"
	"billsoftware/backend/internal/service"
)

type UserProfileHandler struct {
	engine *xorm.Engine
}

func NewUserProfileHandler(engine *xorm.Engine) *UserProfileHandler {
	return &UserProfileHandler{engine: engine}
}

func (h *UserProfileHandler) GetProfile(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	user := &model.User{}
	has, err := h.engine.ID(userID).Get(user)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query user profile failed")
		return
	}
	if !has {
		response.Fail(c, http.StatusNotFound, "user profile not found")
		return
	}

	response.Success(c, buildUserProfilePayload(*user))
}

func (h *UserProfileHandler) UpdateProfile(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	payload, err := parseJSONMap(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid profile payload")
		return
	}

	nickname := strings.TrimSpace(readString(payload, "nickname"))
	email := strings.TrimSpace(readString(payload, "email"))
	avatar := readString(payload, "avatar")
	avatarOriginal := readString(payload, "avatar_original", "avatarOriginal")
	avatarCompressed := readString(payload, "avatar_compressed", "avatarCompressed")

	if nickname == "" {
		response.Fail(c, http.StatusBadRequest, "nickname is required")
		return
	}
	if avatarOriginal == "" && avatar != "" {
		avatarOriginal = avatar
	}
	if avatarCompressed == "" && avatar != "" && avatarOriginal == "" {
		avatarCompressed = avatar
	}

	updateMap := map[string]any{
		"nickname":          nickname,
		"email":             email,
		"avatar_original":   avatarOriginal,
		"avatar_compressed": avatarCompressed,
		"updated_at":        time.Now(),
	}

	if _, err := h.engine.Table(new(model.User)).ID(userID).Update(updateMap); err != nil {
		response.Fail(c, http.StatusInternalServerError, "update user profile failed")
		return
	}

	user := &model.User{}
	if _, err := h.engine.ID(userID).Get(user); err != nil {
		response.Fail(c, http.StatusInternalServerError, "query user profile failed")
		return
	}

	response.Success(c, buildUserProfilePayload(*user))
}

func (h *UserProfileHandler) UpdatePassword(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	payload, err := parseJSONMap(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid password payload")
		return
	}

	currentPassword := strings.TrimSpace(readString(payload, "current_password", "currentPassword"))
	newPassword := strings.TrimSpace(readString(payload, "new_password", "newPassword"))
	confirmPassword := strings.TrimSpace(readString(payload, "confirm_password", "confirmPassword"))

	if currentPassword == "" || newPassword == "" || confirmPassword == "" {
		response.Fail(c, http.StatusBadRequest, "current_password, new_password and confirm_password are required")
		return
	}
	if len(newPassword) < 6 {
		response.Fail(c, http.StatusBadRequest, "new password length must be at least 6")
		return
	}
	if newPassword != confirmPassword {
		response.Fail(c, http.StatusBadRequest, "new password and confirm password do not match")
		return
	}

	user := &model.User{}
	has, err := h.engine.ID(userID).Get(user)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query user failed")
		return
	}
	if !has {
		response.Fail(c, http.StatusNotFound, "user not found")
		return
	}
	if user.PasswordHash != service.HashPassword(currentPassword) {
		response.Fail(c, http.StatusBadRequest, "current password is incorrect")
		return
	}

	session := h.engine.NewSession()
	defer func() {
		_ = session.Close()
	}()

	if err := session.Begin(); err != nil {
		response.Fail(c, http.StatusInternalServerError, "begin transaction failed")
		return
	}

	if _, err := session.Table(new(model.User)).ID(userID).Update(map[string]any{
		"password_hash": service.HashPassword(newPassword),
		"updated_at":    time.Now(),
	}); err != nil {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "update user password failed")
		return
	}

	if _, err := session.Where("user_id = ? AND is_active = ?", userID, 1).Cols("is_active").Update(&model.UserSession{
		IsActive: 0,
	}); err != nil {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "invalidate user sessions failed")
		return
	}

	if err := session.Commit(); err != nil {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "commit transaction failed")
		return
	}

	response.Success(c, gin.H{
		"user_id": userID,
		"message": "password updated, please login again",
	})
}

func buildUserProfilePayload(user model.User) gin.H {
	avatar := user.AvatarCompressed
	if strings.TrimSpace(avatar) == "" {
		avatar = user.AvatarOriginal
	}

	updatedAt := ""
	if !user.UpdatedAt.IsZero() {
		updatedAt = user.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	return gin.H{
		"id":                user.ID,
		"account":           user.Username,
		"username":          user.Username,
		"nickname":          user.Nickname,
		"phone":             user.Phone,
		"email":             user.Email,
		"avatar":            avatar,
		"avatar_original":   user.AvatarOriginal,
		"avatar_compressed": user.AvatarCompressed,
		"updated_at":        updatedAt,
	}
}
