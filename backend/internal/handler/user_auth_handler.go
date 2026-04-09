package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"

	"billsoftware/backend/internal/model"
	"billsoftware/backend/internal/response"
	"billsoftware/backend/internal/service"
)

type UserAuthHandler struct {
	engine *xorm.Engine
}

type RegisterRequest struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func NewUserAuthHandler(engine *xorm.Engine) *UserAuthHandler {
	return &UserAuthHandler{engine: engine}
}

func (h *UserAuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid register payload")
		return
	}

	if req.Username == "" || req.Nickname == "" || req.Phone == "" || req.Password == "" {
		response.Fail(c, http.StatusBadRequest, "username, nickname, phone and password are required")
		return
	}

	usernameExists, err := h.engine.Where("username = ?", req.Username).Exist(&model.User{})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "check username failed")
		return
	}
	if usernameExists {
		response.Fail(c, http.StatusConflict, "username already exists")
		return
	}

	phoneExists, err := h.engine.Where("phone = ?", req.Phone).Exist(&model.User{})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "check phone failed")
		return
	}
	if phoneExists {
		response.Fail(c, http.StatusConflict, "phone already exists")
		return
	}

	user := &model.User{
		Username:     req.Username,
		Nickname:     req.Nickname,
		Phone:        req.Phone,
		Email:        req.Email,
		PasswordHash: service.HashPassword(req.Password),
		Status:       0,
	}

	if _, err := h.engine.Insert(user); err != nil {
		response.Fail(c, http.StatusInternalServerError, "create user failed")
		return
	}

	if err := h.cloneDefaultCategories(user.ID); err != nil {
		response.Fail(c, http.StatusInternalServerError, "init user categories failed")
		return
	}

	response.Success(c, gin.H{
		"user_id": user.ID,
		"status":  "disabled",
		"message": "registered successfully, waiting for admin enable",
	})
}

func (h *UserAuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid login payload")
		return
	}

	user := &model.User{}
	has, err := h.engine.Where("username = ? OR phone = ?", req.Account, req.Account).Get(user)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query user failed")
		return
	}
	if !has {
		response.Fail(c, http.StatusUnauthorized, "username or phone is incorrect")
		return
	}

	if user.Status != 1 {
		response.Fail(c, http.StatusForbidden, "user is disabled")
		return
	}

	if user.PasswordHash != service.HashPassword(req.Password) {
		response.Fail(c, http.StatusUnauthorized, "password is incorrect")
		return
	}

	token, err := service.GenerateToken()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "generate token failed")
		return
	}

	_, _ = h.engine.Where("user_id = ? AND is_active = ?", user.ID, 1).Cols("is_active").Update(&model.UserSession{
		IsActive: 0,
	})

	session := &model.UserSession{
		UserID:       user.ID,
		SessionToken: token,
		ClientType:   "pc",
		IsActive:     1,
		LoginAt:      time.Now(),
	}

	if _, err := h.engine.Insert(session); err != nil {
		response.Fail(c, http.StatusInternalServerError, "create user session failed")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"profile": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"phone":    user.Phone,
			"email":    user.Email,
			"status":   user.Status,
		},
	})
}

func (h *UserAuthHandler) cloneDefaultCategories(userID uint64) error {
	templates := make([]model.CategoryTemplate, 0)
	if err := h.engine.Asc("category_type", "sort_order", "id").Find(&templates); err != nil {
		return err
	}

	records := make([]model.UserCategory, 0, len(templates))
	if len(templates) > 0 {
		for _, item := range templates {
			records = append(records, model.UserCategory{
				UserID:       userID,
				CategoryType: item.CategoryType,
				Name:         item.Name,
				SortOrder:    item.SortOrder,
				IsSystem:     1,
			})
		}
	} else {
		expenseDefaults := []string{
			"早午晚餐", "餐饮", "购物", "日用", "奶茶", "交通", "蔬菜", "水果", "零食", "运动", "娱乐",
			"通讯", "服饰", "美容", "住房", "居家", "孩子", "长辈", "社交", "旅行", "烟酒", "数码",
			"汽车", "医疗", "学习", "宠物", "礼物", "办公", "维修", "捐赠", "彩票", "亲友", "快递",
		}
		incomeDefaults := []string{"工资", "兼职", "理财", "礼金", "其他"}

		sortOrder := 1
		for _, name := range expenseDefaults {
			records = append(records, model.UserCategory{
				UserID:       userID,
				CategoryType: "expense",
				Name:         name,
				SortOrder:    sortOrder,
				IsSystem:     1,
			})
			sortOrder += 1
		}

		sortOrder = 1
		for _, name := range incomeDefaults {
			records = append(records, model.UserCategory{
				UserID:       userID,
				CategoryType: "income",
				Name:         name,
				SortOrder:    sortOrder,
				IsSystem:     1,
			})
			sortOrder += 1
		}
	}

	if len(records) == 0 {
		return nil
	}

	_, err := h.engine.Insert(&records)
	return err
}
