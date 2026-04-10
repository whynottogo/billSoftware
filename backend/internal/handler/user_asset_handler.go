package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"

	"billsoftware/backend/internal/model"
	"billsoftware/backend/internal/response"
)

type UserAssetHandler struct {
	engine *xorm.Engine
}

type assetMonthlyChangeRow struct {
	AccountID      uint64  `xorm:"account_id"`
	MonthlyChanged float64 `xorm:"monthly_change"`
}

type assetCategoryPayload struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Tone     string  `json:"tone"`
	Total    float64 `json:"total"`
	Accounts []gin.H `json:"accounts"`
}

func NewUserAssetHandler(engine *xorm.Engine) *UserAssetHandler {
	return &UserAssetHandler{engine: engine}
}

func (h *UserAssetHandler) List(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	accounts := make([]model.AssetAccount, 0)
	if err := h.engine.Where("user_id = ?", userID).Asc("id").Find(&accounts); err != nil {
		response.Fail(c, http.StatusInternalServerError, "query asset accounts failed")
		return
	}

	monthlyChanges, err := h.queryMonthlyChanges(userID, time.Now())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query asset monthly changes failed")
		return
	}

	categories, summary := buildAssetOverviewPayload(accounts, monthlyChanges)
	response.Success(c, gin.H{
		"summary":    summary,
		"categories": categories,
	})
}

func (h *UserAssetHandler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	payload, err := parseJSONMap(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid asset payload")
		return
	}

	accountType := normalizeAssetAccountType(readString(payload, "account_type", "type"))
	name := strings.TrimSpace(readString(payload, "account_name", "name"))
	remark := strings.TrimSpace(readString(payload, "remark", "note"))
	cardNo := strings.TrimSpace(readString(payload, "card_no", "cardNo"))
	provider := normalizeAssetProvider(readString(payload, "provider", "sub_type"))
	balance, hasBalance := readFloat(payload, "balance")

	if err := validateAssetPayload(accountType, name, remark, cardNo, provider, balance, hasBalance); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	account := &model.AssetAccount{
		UserID:        userID,
		AccountType:   accountType,
		SubType:       provider,
		Name:          name,
		Remark:        remark,
		CardNo:        cardNo,
		BalanceNature: inferBalanceNature(accountType),
		Balance:       round2(balance),
	}

	session := h.engine.NewSession()
	defer func() {
		_ = session.Close()
	}()

	if err := session.Begin(); err != nil {
		response.Fail(c, http.StatusInternalServerError, "begin transaction failed")
		return
	}

	if _, err := session.Insert(account); err != nil {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "create asset account failed")
		return
	}

	if err := insertAssetLog(session, account.ID, userID, "adjust", 0, account.Balance, account.Balance, "创建账户", now); err != nil {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "create asset account log failed")
		return
	}

	if err := session.Commit(); err != nil {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "commit transaction failed")
		return
	}

	response.Success(c, gin.H{
		"account": buildAssetAccountPayload(*account, accountCategoryMeta(*account), account.Balance),
	})
}

func (h *UserAssetHandler) Update(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	accountID, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil || accountID == 0 {
		response.Fail(c, http.StatusBadRequest, "invalid asset account id")
		return
	}

	payload, err := parseJSONMap(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid asset payload")
		return
	}

	account := &model.AssetAccount{}
	has, err := h.engine.Where("id = ? AND user_id = ?", accountID, userID).Get(account)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query asset account failed")
		return
	}
	if !has {
		response.Fail(c, http.StatusNotFound, "asset account not found")
		return
	}

	name := strings.TrimSpace(readString(payload, "account_name", "name"))
	if name == "" {
		name = account.Name
	}
	remark := readOptionalString(payload, account.Remark, "remark", "note")
	cardNo := readOptionalString(payload, account.CardNo, "card_no", "cardNo")
	provider := normalizeAssetProvider(readOptionalString(payload, account.SubType, "provider", "sub_type"))
	balance, hasBalance := readFloat(payload, "balance")
	if !hasBalance {
		balance = account.Balance
	}

	if err := validateAssetPayload(account.AccountType, name, strings.TrimSpace(remark), strings.TrimSpace(cardNo), provider, balance, true); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	updateMap := map[string]any{
		"name":           name,
		"remark":         strings.TrimSpace(remark),
		"card_no":        strings.TrimSpace(cardNo),
		"sub_type":       provider,
		"balance_nature": inferBalanceNature(account.AccountType),
		"balance":        round2(balance),
		"updated_at":     time.Now(),
	}

	session := h.engine.NewSession()
	defer func() {
		_ = session.Close()
	}()

	if err := session.Begin(); err != nil {
		response.Fail(c, http.StatusInternalServerError, "begin transaction failed")
		return
	}

	if _, err := session.Table(new(model.AssetAccount)).Where("id = ? AND user_id = ?", accountID, userID).Update(updateMap); err != nil {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "update asset account failed")
		return
	}

	if round2(balance) != round2(account.Balance) {
		change := round2(balance - account.Balance)
		if err := insertAssetLog(session, accountID, userID, "settings_adjust", account.Balance, change, balance, "更新账户设置", time.Now()); err != nil {
			_ = session.Rollback()
			response.Fail(c, http.StatusInternalServerError, "create asset account log failed")
			return
		}
	}

	if err := session.Commit(); err != nil {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "commit transaction failed")
		return
	}

	account.Name = name
	account.Remark = strings.TrimSpace(remark)
	account.CardNo = strings.TrimSpace(cardNo)
	account.SubType = provider
	account.BalanceNature = inferBalanceNature(account.AccountType)
	account.Balance = round2(balance)
	account.UpdatedAt = time.Now()

	response.Success(c, gin.H{
		"account": buildAssetAccountPayload(*account, accountCategoryMeta(*account), 0),
	})
}

func (h *UserAssetHandler) Detail(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	accountID, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil || accountID == 0 {
		response.Fail(c, http.StatusBadRequest, "invalid asset account id")
		return
	}

	account := &model.AssetAccount{}
	has, err := h.engine.Where("id = ? AND user_id = ?", accountID, userID).Get(account)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query asset account failed")
		return
	}
	if !has {
		response.Fail(c, http.StatusNotFound, "asset account not found")
		return
	}

	monthlyChanges, err := h.queryMonthlyChanges(userID, time.Now())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query asset monthly changes failed")
		return
	}

	logs := make([]model.AssetAccountLog, 0)
	if err := h.engine.Where("account_id = ? AND user_id = ?", accountID, userID).Desc("log_date", "id").Find(&logs); err != nil {
		response.Fail(c, http.StatusInternalServerError, "query asset account logs failed")
		return
	}

	category := accountCategoryMeta(*account)
	records := make([]gin.H, 0, len(logs))
	for _, item := range logs {
		records = append(records, buildAssetLogPayload(item))
	}

	response.Success(c, gin.H{
		"account":  buildAssetAccountPayload(*account, category, monthlyChanges[account.ID]),
		"category": gin.H{"id": category.ID, "category_id": category.ID, "name": category.Name, "category_name": category.Name},
		"records":  records,
	})
}

func (h *UserAssetHandler) CreateOperation(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	accountID, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil || accountID == 0 {
		response.Fail(c, http.StatusBadRequest, "invalid asset account id")
		return
	}

	payload, err := parseJSONMap(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid asset operation payload")
		return
	}

	actionType := normalizeAssetActionType(readString(payload, "action_type", "actionType", "type"))
	note := strings.TrimSpace(readString(payload, "note", "remark"))
	amount, hasAmount := readFloat(payload, "amount")
	targetBalance, hasTargetBalance := readFloat(payload, "target_balance", "targetBalance")

	if note == "" {
		response.Fail(c, http.StatusBadRequest, "note is required")
		return
	}
	if !hasAmount || amount <= 0 {
		response.Fail(c, http.StatusBadRequest, "amount must be greater than 0")
		return
	}
	if actionType == "" {
		response.Fail(c, http.StatusBadRequest, "action_type is required")
		return
	}

	account := &model.AssetAccount{}
	has, err := h.engine.Where("id = ? AND user_id = ?", accountID, userID).Get(account)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query asset account failed")
		return
	}
	if !has {
		response.Fail(c, http.StatusNotFound, "asset account not found")
		return
	}

	before := round2(account.Balance)
	change := 0.0
	after := before

	switch actionType {
	case "adjust":
		if hasTargetBalance {
			after = round2(targetBalance)
		} else {
			after = round2(amount)
		}
		change = round2(after - before)
	case "increase":
		change = round2(amount)
		after = round2(before + change)
	case "decrease":
		change = round2(-amount)
		after = round2(before + change)
	default:
		response.Fail(c, http.StatusBadRequest, "action_type must be adjust, increase or decrease")
		return
	}

	if after < 0 {
		response.Fail(c, http.StatusBadRequest, "result balance must be greater than or equal to 0")
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

	if _, err := session.Table(new(model.AssetAccount)).Where("id = ? AND user_id = ?", accountID, userID).Update(map[string]any{
		"balance":    after,
		"updated_at": time.Now(),
	}); err != nil {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "update asset balance failed")
		return
	}

	if err := insertAssetLog(session, accountID, userID, actionType, before, change, after, note, time.Now()); err != nil {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "create asset account log failed")
		return
	}

	if err := session.Commit(); err != nil {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "commit transaction failed")
		return
	}

	response.Success(c, gin.H{
		"account_id": accountID,
		"balance":    after,
	})
}

func (h *UserAssetHandler) queryMonthlyChanges(userID uint64, baseTime time.Time) (map[uint64]float64, error) {
	start := time.Date(baseTime.Year(), baseTime.Month(), 1, 0, 0, 0, 0, baseTime.Location())
	end := start.AddDate(0, 1, 0)

	rows := make([]assetMonthlyChangeRow, 0)
	err := h.engine.SQL(`
		SELECT
			account_id,
			COALESCE(SUM(amount_change), 0) AS monthly_change
		FROM asset_account_logs
		WHERE user_id = ? AND log_date >= ? AND log_date < ?
		GROUP BY account_id
	`, userID, start, end).Find(&rows)
	if err != nil {
		return nil, err
	}

	result := make(map[uint64]float64, len(rows))
	for _, item := range rows {
		result[item.AccountID] = round2(item.MonthlyChanged)
	}
	return result, nil
}

func buildAssetOverviewPayload(accounts []model.AssetAccount, monthlyChanges map[uint64]float64) ([]assetCategoryPayload, gin.H) {
	categories := []assetCategoryPayload{
		{ID: "category-liquid", Name: "流动资产", Tone: "info", Accounts: make([]gin.H, 0)},
		{ID: "category-investment", Name: "投资资产", Tone: "success", Accounts: make([]gin.H, 0)},
		{ID: "category-liability", Name: "负债账户", Tone: "danger", Accounts: make([]gin.H, 0)},
	}
	categoryIndex := map[string]int{
		"category-liquid":     0,
		"category-investment": 1,
		"category-liability":  2,
	}

	totalAsset := 0.0
	totalLiability := 0.0
	totalMonthlyChange := 0.0
	latestUpdatedAt := ""

	for _, account := range accounts {
		category := accountCategoryMeta(account)
		index := categoryIndex[category.ID]
		monthlyChange := monthlyChanges[account.ID]
		categories[index].Accounts = append(categories[index].Accounts, buildAssetAccountPayload(account, category, monthlyChange))
		categories[index].Total = round2(categories[index].Total + account.Balance)
		totalMonthlyChange = round2(totalMonthlyChange + monthlyChange)
		if category.ID == "category-liability" {
			totalLiability = round2(totalLiability + account.Balance)
		} else {
			totalAsset = round2(totalAsset + account.Balance)
		}
		updatedAt := account.UpdatedAt.Format("2006-01-02 15:04")
		if updatedAt > latestUpdatedAt {
			latestUpdatedAt = updatedAt
		}
	}

	summary := gin.H{
		"netAsset":        round2(totalAsset - totalLiability),
		"net_asset":       round2(totalAsset - totalLiability),
		"totalAsset":      totalAsset,
		"total_asset":     totalAsset,
		"totalLiability":  totalLiability,
		"total_liability": totalLiability,
		"monthlyChange":   totalMonthlyChange,
		"monthly_change":  totalMonthlyChange,
		"accountCount":    len(accounts),
		"account_count":   len(accounts),
		"updatedAt":       latestUpdatedAt,
		"updated_at":      latestUpdatedAt,
	}

	return categories, summary
}

func buildAssetAccountPayload(account model.AssetAccount, category assetCategoryMetaPayload, monthlyChange float64) gin.H {
	frontendType := assetFrontendType(account.AccountType)
	typeLabel := assetTypeLabel(account.AccountType)
	provider := normalizeAssetProvider(account.SubType)
	cardNo := strings.TrimSpace(account.CardNo)

	return gin.H{
		"id":              account.ID,
		"account_id":      account.ID,
		"name":            account.Name,
		"account_name":    account.Name,
		"type":            frontendType,
		"account_type":    account.AccountType,
		"typeLabel":       typeLabel,
		"type_label":      typeLabel,
		"remark":          account.Remark,
		"note":            account.Remark,
		"balance":         round2(account.Balance),
		"current_balance": round2(account.Balance),
		"cardNo":          cardNo,
		"card_no":         cardNo,
		"provider":        provider,
		"direction":       directionFromBalanceNature(account.BalanceNature),
		"monthlyChange":   round2(monthlyChange),
		"monthly_change":  round2(monthlyChange),
		"categoryId":      category.ID,
		"category_id":     category.ID,
		"categoryName":    category.Name,
		"category_name":   category.Name,
	}
}

func buildAssetLogPayload(log model.AssetAccountLog) gin.H {
	actionType, actionLabel, source := assetLogMeta(log.ChangeType)
	dateLabel := log.LogDate.Format("01-02 15:04")
	monthKey := log.LogDate.Format("2006-01")

	return gin.H{
		"id":              log.ID,
		"record_id":       log.ID,
		"operation_id":    log.ID,
		"monthKey":        monthKey,
		"month_key":       monthKey,
		"dateLabel":       dateLabel,
		"date_label":      dateLabel,
		"created_at":      log.LogDate.Format("2006-01-02 15:04:05"),
		"occurred_at":     log.LogDate.Format("2006-01-02 15:04:05"),
		"action":          actionLabel,
		"action_label":    actionLabel,
		"action_type":     actionType,
		"change":          round2(log.AmountChange),
		"delta":           round2(log.AmountChange),
		"change_amount":   round2(log.AmountChange),
		"balanceAfter":    round2(log.AmountAfter),
		"balance_after":   round2(log.AmountAfter),
		"after_balance":   round2(log.AmountAfter),
		"current_balance": round2(log.AmountAfter),
		"note":            log.Remark,
		"remark":          log.Remark,
		"source":          source,
		"source_label":    source,
		"operator":        source,
	}
}

func insertAssetLog(session *xorm.Session, accountID uint64, userID uint64, changeType string, before float64, change float64, after float64, remark string, logTime time.Time) error {
	record := &model.AssetAccountLog{
		AccountID:    accountID,
		UserID:       userID,
		ChangeType:   changeType,
		AmountBefore: round2(before),
		AmountChange: round2(change),
		AmountAfter:  round2(after),
		Remark:       strings.TrimSpace(remark),
		LogDate:      logTime,
	}
	_, err := session.Insert(record)
	return err
}

func validateAssetPayload(accountType string, name string, remark string, cardNo string, provider string, balance float64, hasBalance bool) error {
	if accountType == "" {
		return httpError("account_type is required")
	}
	if name == "" {
		return httpError("name is required")
	}
	if remark == "" {
		return httpError("remark is required")
	}
	if !hasBalance {
		return httpError("balance is required")
	}
	if round2(balance) < 0 {
		return httpError("balance must be greater than or equal to 0")
	}
	if accountType == "bank_card" || accountType == "credit_card" {
		if strings.TrimSpace(cardNo) == "" {
			return httpError("card_no is required for bank card or credit card")
		}
	}
	if accountType == "virtual" && provider == "" {
		return httpError("provider is required for virtual account")
	}
	return nil
}

type assetCategoryMetaPayload struct {
	ID   string
	Name string
	Tone string
}

func accountCategoryMeta(account model.AssetAccount) assetCategoryMetaPayload {
	if directionFromBalanceNature(account.BalanceNature) == "liability" {
		return assetCategoryMetaPayload{ID: "category-liability", Name: "负债账户", Tone: "danger"}
	}
	if account.AccountType == "investment" {
		return assetCategoryMetaPayload{ID: "category-investment", Name: "投资资产", Tone: "success"}
	}
	return assetCategoryMetaPayload{ID: "category-liquid", Name: "流动资产", Tone: "info"}
}

func directionFromBalanceNature(value string) string {
	target := strings.ToLower(strings.TrimSpace(value))
	if target == "liability" || target == "debt" {
		return "liability"
	}
	return "asset"
}

func inferBalanceNature(accountType string) string {
	switch accountType {
	case "credit_card", "liability":
		return "liability"
	default:
		return "asset"
	}
}

func assetFrontendType(accountType string) string {
	switch accountType {
	case "bank_card":
		return "bankCard"
	case "credit_card":
		return "creditCard"
	default:
		return accountType
	}
}

func assetTypeLabel(accountType string) string {
	switch accountType {
	case "cash":
		return "现金"
	case "bank_card":
		return "银行卡"
	case "credit_card":
		return "信用卡"
	case "virtual":
		return "虚拟账户"
	case "investment":
		return "投资账户"
	case "liability":
		return "负债账户"
	default:
		return "账户"
	}
}

func normalizeAssetAccountType(value string) string {
	target := strings.ToLower(strings.TrimSpace(value))
	switch target {
	case "cash":
		return "cash"
	case "bankcard", "bank_card", "bank":
		return "bank_card"
	case "creditcard", "credit_card", "credit":
		return "credit_card"
	case "virtual", "wallet", "ewallet":
		return "virtual"
	case "investment", "invest":
		return "investment"
	case "liability", "debt":
		return "liability"
	default:
		return target
	}
}

func normalizeAssetProvider(value string) string {
	target := strings.ToLower(strings.TrimSpace(value))
	switch target {
	case "wechat":
		return "wechat"
	case "alipay":
		return "alipay"
	default:
		return ""
	}
}

func normalizeAssetActionType(value string) string {
	target := strings.ToLower(strings.TrimSpace(value))
	switch target {
	case "adjust":
		return "adjust"
	case "increase", "credit", "deposit":
		return "increase"
	case "decrease", "debit", "withdraw":
		return "decrease"
	default:
		return ""
	}
}

func assetLogMeta(changeType string) (string, string, string) {
	target := strings.ToLower(strings.TrimSpace(changeType))
	source := "手动操作"
	switch target {
	case "settings_adjust":
		return "adjust", "调整", "账户设置"
	case "increase":
		return "increase", "增加", source
	case "decrease":
		return "decrease", "减少", source
	default:
		return "adjust", "调整", source
	}
}

func parseJSONMap(c *gin.Context) (map[string]any, error) {
	payload := make(map[string]any)
	if err := c.ShouldBindJSON(&payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func readString(payload map[string]any, keys ...string) string {
	for _, key := range keys {
		raw, exists := payload[key]
		if !exists || raw == nil {
			continue
		}
		switch value := raw.(type) {
		case string:
			if strings.TrimSpace(value) != "" {
				return value
			}
		case float64:
			return strconv.FormatFloat(value, 'f', -1, 64)
		case int:
			return strconv.Itoa(value)
		case int64:
			return strconv.FormatInt(value, 10)
		}
	}
	return ""
}

func readOptionalString(payload map[string]any, fallback string, keys ...string) string {
	value := readString(payload, keys...)
	if value == "" {
		return fallback
	}
	return value
}

func readFloat(payload map[string]any, keys ...string) (float64, bool) {
	for _, key := range keys {
		raw, exists := payload[key]
		if !exists || raw == nil {
			continue
		}
		switch value := raw.(type) {
		case float64:
			return value, true
		case float32:
			return float64(value), true
		case int:
			return float64(value), true
		case int64:
			return float64(value), true
		case string:
			parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
			if err == nil {
				return parsed, true
			}
		}
	}
	return 0, false
}

type simpleHTTPError struct {
	message string
}

func (e simpleHTTPError) Error() string {
	return e.message
}

func httpError(message string) error {
	return simpleHTTPError{message: message}
}
