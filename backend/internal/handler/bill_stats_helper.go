package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"xorm.io/xorm"
)

var billColorPalette = []string{
	"#f6d34a",
	"#6bcf7c",
	"#4d96ff",
	"#ff8b8b",
	"#9b8cff",
	"#34d399",
	"#f97316",
	"#0ea5e9",
}

type billStatsCore struct {
	engine *xorm.Engine
}

type billPeriodSummary struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

type billHighlightCard struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Hint  string `json:"hint"`
}

type billAchievement struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Tone  string `json:"tone"`
	Hint  string `json:"hint"`
}

type billCategorySplitItem struct {
	Name    string  `json:"name"`
	Value   float64 `json:"value"`
	Percent float64 `json:"percent"`
	Color   string  `json:"color"`
}

type billRankingItem struct {
	Name    string  `json:"name"`
	Value   float64 `json:"value"`
	Percent float64 `json:"percent"`
	Count   int     `json:"count"`
	Trend   string  `json:"trend"`
	Badge   string  `json:"badge"`
}

type billTrendPoint struct {
	Label  string  `json:"label"`
	Amount float64 `json:"amount"`
}

type billComparisonItem struct {
	Label   string  `json:"label"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

type billMonthListItem struct {
	Key            string              `json:"key"`
	Label          string              `json:"label"`
	Year           int                 `json:"year"`
	Month          int                 `json:"month"`
	Income         float64             `json:"income"`
	Expense        float64             `json:"expense"`
	Balance        float64             `json:"balance"`
	Days           int                 `json:"days"`
	Records        int                 `json:"records"`
	Status         string              `json:"status"`
	Note           string              `json:"note"`
	Highlight      string              `json:"highlight"`
	HighlightCards []billHighlightCard `json:"highlightCards"`
}

type billYearHistoryItem struct {
	Year    int `json:"year"`
	Summary struct {
		Balance float64 `json:"balance"`
		Income  float64 `json:"income"`
		Expense float64 `json:"expense"`
		Months  int     `json:"months"`
		Days    int     `json:"days"`
		Records int     `json:"records"`
		Insight string  `json:"insight"`
	} `json:"summary"`
}

type billYearDetailResponse struct {
	Year    int `json:"year"`
	Summary struct {
		Balance float64 `json:"balance"`
		Income  float64 `json:"income"`
		Expense float64 `json:"expense"`
		Months  int     `json:"months"`
		Days    int     `json:"days"`
		Records int     `json:"records"`
		Insight string  `json:"insight"`
	} `json:"summary"`
	Months []billMonthListItem `json:"months"`
}

type billMonthDetailResponse struct {
	Key             string                  `json:"key"`
	Label           string                  `json:"label"`
	PreviousKey     string                  `json:"previousKey"`
	NextKey         string                  `json:"nextKey"`
	PreviousLabel   string                  `json:"previousLabel"`
	PreviousBalance float64                 `json:"previousBalance"`
	Income          float64                 `json:"income"`
	Expense         float64                 `json:"expense"`
	Balance         float64                 `json:"balance"`
	Days            int                     `json:"days"`
	Records         int                     `json:"records"`
	Note            string                  `json:"note"`
	Highlight       string                  `json:"highlight"`
	CategorySplit   []billCategorySplitItem `json:"categorySplit"`
	Ranking         []billRankingItem       `json:"ranking"`
	DailyTrend      []billTrendPoint        `json:"dailyTrend"`
	Comparison      []billComparisonItem    `json:"comparison"`
	HighlightCards  []billHighlightCard     `json:"highlightCards"`
	Achievements    []billAchievement       `json:"achievements"`
}

type adminUserBillOverviewResponse struct {
	Profile struct {
		UserID       uint64 `json:"userId"`
		Username     string `json:"username"`
		Nickname     string `json:"nickname"`
		Status       string `json:"status"`
		RegisterDate string `json:"registerDate"`
		BillCount    int    `json:"billCount"`
	} `json:"profile"`
	MonthOptions []string                     `json:"monthOptions"`
	YearOptions  []string                     `json:"yearOptions"`
	Monthly      map[string]adminMonthSummary `json:"monthly"`
	Yearly       map[string]adminYearSummary  `json:"yearly"`
}

type adminMonthSummary struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
	Records int     `json:"records"`
	Insight string  `json:"insight"`
}

type adminYearSummary struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
	Months  int     `json:"months"`
	Insight string  `json:"insight"`
}

type periodAggregateRow struct {
	PeriodKey string  `xorm:"period_key"`
	Income    float64 `xorm:"income"`
	Expense   float64 `xorm:"expense"`
	Records   int     `xorm:"records"`
	Days      int     `xorm:"days"`
	Months    int     `xorm:"months"`
}

type categoryAggregateRow struct {
	CategoryName string  `xorm:"category_name"`
	Amount       float64 `xorm:"amount"`
	RecordCount  int     `xorm:"record_count"`
}

type dailyExpenseRow struct {
	DateKey string  `xorm:"date_key"`
	Amount  float64 `xorm:"amount"`
}

type userProfileRow struct {
	ID        uint64    `xorm:"id"`
	Username  string    `xorm:"username"`
	Nickname  string    `xorm:"nickname"`
	Status    int       `xorm:"status"`
	CreatedAt time.Time `xorm:"created_at"`
	BillCount int       `xorm:"bill_count"`
}

func newBillStatsCore(engine *xorm.Engine) *billStatsCore {
	return &billStatsCore{engine: engine}
}

func (c *billStatsCore) ListBillYears(userID uint64) ([]int, []billYearHistoryItem, error) {
	rows := make([]periodAggregateRow, 0)
	if err := c.engine.SQL(`
		SELECT
			DATE_FORMAT(record_date, '%Y') AS period_key,
			COALESCE(SUM(CASE WHEN record_type = 'income' THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN record_type = 'expense' THEN amount ELSE 0 END), 0) AS expense,
			COUNT(*) AS records,
			COUNT(DISTINCT record_date) AS days,
			COUNT(DISTINCT DATE_FORMAT(record_date, '%Y-%m')) AS months
		FROM ledger_records
		WHERE user_id = ?
		GROUP BY DATE_FORMAT(record_date, '%Y')
		ORDER BY period_key DESC
	`, userID).Find(&rows); err != nil {
		return nil, nil, err
	}

	years := make([]int, 0, len(rows))
	history := make([]billYearHistoryItem, 0, len(rows))
	for _, row := range rows {
		year, err := strconv.Atoi(row.PeriodKey)
		if err != nil {
			continue
		}
		years = append(years, year)
		item := billYearHistoryItem{Year: year}
		item.Summary.Income = round2(row.Income)
		item.Summary.Expense = round2(row.Expense)
		item.Summary.Balance = round2(row.Income - row.Expense)
		item.Summary.Months = row.Months
		item.Summary.Days = row.Days
		item.Summary.Records = row.Records
		item.Summary.Insight = buildYearInsight(item.Summary.Income, item.Summary.Expense, item.Summary.Balance, item.Summary.Months)
		history = append(history, item)
	}

	if len(years) == 0 {
		currentYear := time.Now().Year()
		years = append(years, currentYear)
		item := billYearHistoryItem{Year: currentYear}
		item.Summary.Insight = "当前年份暂无账单记录。"
		history = append(history, item)
	}

	return years, history, nil
}

func (c *billStatsCore) GetBillYearDetail(userID uint64, year int) (billYearDetailResponse, error) {
	response := billYearDetailResponse{Year: year, Months: make([]billMonthListItem, 0)}

	rows := make([]periodAggregateRow, 0)
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(1, 0, 0)
	if err := c.engine.SQL(`
		SELECT
			DATE_FORMAT(record_date, '%Y-%m') AS period_key,
			COALESCE(SUM(CASE WHEN record_type = 'income' THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN record_type = 'expense' THEN amount ELSE 0 END), 0) AS expense,
			COUNT(*) AS records,
			COUNT(DISTINCT record_date) AS days
		FROM ledger_records
		WHERE user_id = ? AND record_date >= ? AND record_date < ?
		GROUP BY DATE_FORMAT(record_date, '%Y-%m')
		ORDER BY period_key DESC
	`, userID, start.Format("2006-01-02"), end.Format("2006-01-02")).Find(&rows); err != nil {
		return response, err
	}

	totalMonths := 0
	totalDays := 0
	totalRecords := 0
	totalIncome := 0.0
	totalExpense := 0.0
	currentMonthKey := time.Now().Format("2006-01")
	for _, row := range rows {
		monthItem := buildBillMonthListItem(row, currentMonthKey)
		response.Months = append(response.Months, monthItem)
		totalMonths += 1
		totalDays += row.Days
		totalRecords += row.Records
		totalIncome += row.Income
		totalExpense += row.Expense
	}

	response.Summary.Income = round2(totalIncome)
	response.Summary.Expense = round2(totalExpense)
	response.Summary.Balance = round2(totalIncome - totalExpense)
	response.Summary.Months = totalMonths
	response.Summary.Days = totalDays
	response.Summary.Records = totalRecords
	response.Summary.Insight = buildYearInsight(response.Summary.Income, response.Summary.Expense, response.Summary.Balance, totalMonths)

	return response, nil
}

func (c *billStatsCore) GetBillMonthDetail(userID uint64, monthKey string) (billMonthDetailResponse, error) {
	response := buildEmptyBillMonthDetail(monthKey)

	monthStart, monthEnd, normalizedMonthKey, err := resolveMonthRange(monthKey)
	if err != nil {
		return response, err
	}
	response.Key = normalizedMonthKey
	response.Label = monthLabel(normalizedMonthKey)

	monthSummary, err := c.loadMonthAggregate(userID, normalizedMonthKey)
	if err != nil {
		return response, err
	}
	response.Income = round2(monthSummary.Income)
	response.Expense = round2(monthSummary.Expense)
	response.Balance = round2(monthSummary.Income - monthSummary.Expense)
	response.Days = monthSummary.Days
	response.Records = monthSummary.Records
	response.Note = buildMonthNote(response.Income, response.Expense, response.Days, response.Records)
	response.Highlight = buildMonthHighlight(response.Income, response.Expense, response.Balance)

	monthKeys, err := c.listMonthKeys(userID)
	if err != nil {
		return response, err
	}
	previousKey, nextKey := adjacentMonthKeys(monthKeys, normalizedMonthKey)
	response.PreviousKey = previousKey
	response.NextKey = nextKey
	response.PreviousLabel = monthLabel(previousKey)
	if previousKey != "" {
		previousSummary, err := c.loadMonthAggregate(userID, previousKey)
		if err != nil {
			return response, err
		}
		response.PreviousBalance = round2(previousSummary.Income - previousSummary.Expense)
	}

	categoryRows := make([]categoryAggregateRow, 0)
	if err := c.engine.SQL(`
		SELECT
			COALESCE(uc.name, '未分类') AS category_name,
			COALESCE(SUM(lr.amount), 0) AS amount,
			COUNT(*) AS record_count
		FROM ledger_records lr
		LEFT JOIN user_categories uc ON uc.id = lr.category_id AND uc.user_id = lr.user_id
		WHERE lr.user_id = ? AND lr.record_type = 'expense' AND lr.record_date >= ? AND lr.record_date < ?
		GROUP BY COALESCE(uc.name, '未分类')
		ORDER BY amount DESC, record_count DESC, category_name ASC
	`, userID, monthStart.Format("2006-01-02"), monthEnd.Format("2006-01-02")).Find(&categoryRows); err != nil {
		return response, err
	}

	previousCategoryTotals, err := c.loadPreviousMonthCategoryTotals(userID, normalizedMonthKey)
	if err != nil {
		return response, err
	}

	response.CategorySplit = buildCategorySplit(categoryRows, response.Expense)
	response.Ranking = buildRanking(categoryRows, response.Expense, previousCategoryTotals)

	dailyRows := make([]dailyExpenseRow, 0)
	if err := c.engine.SQL(`
		SELECT
			DATE_FORMAT(record_date, '%m-%d') AS date_key,
			COALESCE(SUM(amount), 0) AS amount
		FROM ledger_records
		WHERE user_id = ? AND record_type = 'expense' AND record_date >= ? AND record_date < ?
		GROUP BY DATE_FORMAT(record_date, '%m-%d')
		ORDER BY date_key ASC
	`, userID, monthStart.Format("2006-01-02"), monthEnd.Format("2006-01-02")).Find(&dailyRows); err != nil {
		return response, err
	}
	response.DailyTrend = buildDailyTrend(dailyRows)

	response.Comparison, err = c.buildSixMonthComparison(userID, normalizedMonthKey)
	if err != nil {
		return response, err
	}

	response.HighlightCards = buildMonthHighlightCards(response.Expense, response.DailyTrend, response.Ranking, response.Income)
	response.Achievements = buildMonthAchievements(monthStart, response.Days, response.Records, response.Ranking)

	return response, nil
}

func (c *billStatsCore) GetAdminUserBillOverview(userID uint64) (adminUserBillOverviewResponse, error) {
	response := adminUserBillOverviewResponse{
		MonthOptions: make([]string, 0),
		YearOptions:  make([]string, 0),
		Monthly:      make(map[string]adminMonthSummary),
		Yearly:       make(map[string]adminYearSummary),
	}

	profile := userProfileRow{}
	has, err := c.engine.SQL(`
		SELECT
			u.id,
			u.username,
			u.nickname,
			u.status,
			u.created_at,
			COUNT(lr.id) AS bill_count
		FROM users u
		LEFT JOIN ledger_records lr ON lr.user_id = u.id
		WHERE u.id = ?
		GROUP BY u.id, u.username, u.nickname, u.status, u.created_at
	`, userID).Get(&profile)
	if err != nil {
		return response, err
	}
	if !has {
		return response, errBillUserNotFound
	}

	response.Profile.UserID = profile.ID
	response.Profile.Username = profile.Username
	response.Profile.Nickname = profile.Nickname
	if profile.Status == 1 {
		response.Profile.Status = "启用"
	} else {
		response.Profile.Status = "禁用"
	}
	response.Profile.RegisterDate = profile.CreatedAt.Format("2006-01-02 15:04")
	response.Profile.BillCount = profile.BillCount

	monthRows := make([]periodAggregateRow, 0)
	if err := c.engine.SQL(`
		SELECT
			DATE_FORMAT(record_date, '%Y-%m') AS period_key,
			COALESCE(SUM(CASE WHEN record_type = 'income' THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN record_type = 'expense' THEN amount ELSE 0 END), 0) AS expense,
			COUNT(*) AS records,
			COUNT(DISTINCT record_date) AS days
		FROM ledger_records
		WHERE user_id = ?
		GROUP BY DATE_FORMAT(record_date, '%Y-%m')
		ORDER BY period_key DESC
	`, userID).Find(&monthRows); err != nil {
		return response, err
	}

	if len(monthRows) == 0 {
		currentMonth := time.Now().Format("2006-01")
		response.MonthOptions = append(response.MonthOptions, currentMonth)
		response.Monthly[currentMonth] = adminMonthSummary{
			Insight: "当前月份暂无账单记录。",
		}
	} else {
		for _, row := range monthRows {
			response.MonthOptions = append(response.MonthOptions, row.PeriodKey)
			response.Monthly[row.PeriodKey] = adminMonthSummary{
				Income:  round2(row.Income),
				Expense: round2(row.Expense),
				Balance: round2(row.Income - row.Expense),
				Records: row.Records,
				Insight: buildAdminMonthInsight(row.Income, row.Expense, row.Records),
			}
		}
	}

	years, history, err := c.ListBillYears(userID)
	if err != nil {
		return response, err
	}
	for _, year := range years {
		response.YearOptions = append(response.YearOptions, strconv.Itoa(year))
	}
	for _, item := range history {
		key := strconv.Itoa(item.Year)
		response.Yearly[key] = adminYearSummary{
			Income:  item.Summary.Income,
			Expense: item.Summary.Expense,
			Balance: item.Summary.Balance,
			Months:  item.Summary.Months,
			Insight: item.Summary.Insight,
		}
	}

	return response, nil
}

func (c *billStatsCore) loadMonthAggregate(userID uint64, monthKey string) (periodAggregateRow, error) {
	row := periodAggregateRow{PeriodKey: monthKey}
	start, end, _, err := resolveMonthRange(monthKey)
	if err != nil {
		return row, err
	}
	_, err = c.engine.SQL(`
		SELECT
			? AS period_key,
			COALESCE(SUM(CASE WHEN record_type = 'income' THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN record_type = 'expense' THEN amount ELSE 0 END), 0) AS expense,
			COUNT(*) AS records,
			COUNT(DISTINCT record_date) AS days
		FROM ledger_records
		WHERE user_id = ? AND record_date >= ? AND record_date < ?
	`, monthKey, userID, start.Format("2006-01-02"), end.Format("2006-01-02")).Get(&row)
	return row, err
}

func (c *billStatsCore) listMonthKeys(userID uint64) ([]string, error) {
	rows := make([]periodAggregateRow, 0)
	if err := c.engine.SQL(`
		SELECT DATE_FORMAT(record_date, '%Y-%m') AS period_key
		FROM ledger_records
		WHERE user_id = ?
		GROUP BY DATE_FORMAT(record_date, '%Y-%m')
		ORDER BY period_key DESC
	`, userID).Find(&rows); err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(rows))
	for _, row := range rows {
		keys = append(keys, row.PeriodKey)
	}
	return keys, nil
}

func (c *billStatsCore) loadPreviousMonthCategoryTotals(userID uint64, monthKey string) (map[string]float64, error) {
	totals := make(map[string]float64)
	monthTime, err := time.ParseInLocation("2006-01", monthKey, time.Local)
	if err != nil {
		return totals, nil
	}
	previousKey := time.Date(monthTime.Year(), monthTime.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, -1, 0).Format("2006-01")
	start, end, _, err := resolveMonthRange(previousKey)
	if err != nil {
		return totals, err
	}

	rows := make([]categoryAggregateRow, 0)
	if err := c.engine.SQL(`
		SELECT
			COALESCE(uc.name, '未分类') AS category_name,
			COALESCE(SUM(lr.amount), 0) AS amount,
			COUNT(*) AS record_count
		FROM ledger_records lr
		LEFT JOIN user_categories uc ON uc.id = lr.category_id AND uc.user_id = lr.user_id
		WHERE lr.user_id = ? AND lr.record_type = 'expense' AND lr.record_date >= ? AND lr.record_date < ?
		GROUP BY COALESCE(uc.name, '未分类')
	`, userID, start.Format("2006-01-02"), end.Format("2006-01-02")).Find(&rows); err != nil {
		return totals, err
	}

	for _, row := range rows {
		totals[row.CategoryName] = round2(row.Amount)
	}
	return totals, nil
}

func (c *billStatsCore) buildSixMonthComparison(userID uint64, monthKey string) ([]billComparisonItem, error) {
	monthTime, err := time.ParseInLocation("2006-01", monthKey, time.Local)
	if err != nil {
		return nil, err
	}

	items := make([]billComparisonItem, 0, 6)
	for offset := 5; offset >= 0; offset -= 1 {
		target := time.Date(monthTime.Year(), monthTime.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, -offset, 0)
		targetKey := target.Format("2006-01")
		row, err := c.loadMonthAggregate(userID, targetKey)
		if err != nil {
			return nil, err
		}
		items = append(items, billComparisonItem{
			Label:   fmt.Sprintf("%d月", int(target.Month())),
			Income:  round2(row.Income),
			Expense: round2(row.Expense),
			Balance: round2(row.Income - row.Expense),
		})
	}
	return items, nil
}

func buildBillMonthListItem(row periodAggregateRow, currentMonthKey string) billMonthListItem {
	monthNumber := monthNumberFromKey(row.PeriodKey)
	income := round2(row.Income)
	expense := round2(row.Expense)
	balance := round2(row.Income - row.Expense)
	item := billMonthListItem{
		Key:            row.PeriodKey,
		Label:          monthLabel(row.PeriodKey),
		Year:           yearNumberFromKey(row.PeriodKey),
		Month:          monthNumber,
		Income:         income,
		Expense:        expense,
		Balance:        balance,
		Days:           row.Days,
		Records:        row.Records,
		Status:         "已结账",
		Note:           buildMonthNote(income, expense, row.Days, row.Records),
		Highlight:      buildMonthHighlight(income, expense, balance),
		HighlightCards: buildMonthListHighlightCards(income, expense, balance, row.Days, row.Records),
	}
	if row.PeriodKey == currentMonthKey {
		item.Status = "本月"
	}
	return item
}

func buildEmptyBillMonthDetail(monthKey string) billMonthDetailResponse {
	key := strings.TrimSpace(monthKey)
	if key == "" {
		key = time.Now().Format("2006-01")
	}
	return billMonthDetailResponse{
		Key:            key,
		Label:          monthLabel(key),
		CategorySplit:  make([]billCategorySplitItem, 0),
		Ranking:        make([]billRankingItem, 0),
		DailyTrend:     make([]billTrendPoint, 0),
		Comparison:     make([]billComparisonItem, 0),
		HighlightCards: make([]billHighlightCard, 0),
		Achievements:   make([]billAchievement, 0),
		Note:           "当前月份暂无账单记录。",
		Highlight:      "暂无可展示的账单统计。",
	}
}

func buildCategorySplit(rows []categoryAggregateRow, totalExpense float64) []billCategorySplitItem {
	items := make([]billCategorySplitItem, 0, len(rows))
	for index, row := range rows {
		percent := 0.0
		if totalExpense > 0 {
			percent = round2(row.Amount / totalExpense * 100)
		}
		items = append(items, billCategorySplitItem{
			Name:    row.CategoryName,
			Value:   round2(row.Amount),
			Percent: percent,
			Color:   billColorPalette[index%len(billColorPalette)],
		})
	}
	return items
}

func buildRanking(rows []categoryAggregateRow, totalExpense float64, previousTotals map[string]float64) []billRankingItem {
	items := make([]billRankingItem, 0, len(rows))
	for index, row := range rows {
		percent := 0.0
		if totalExpense > 0 {
			percent = round2(row.Amount / totalExpense * 100)
		}
		trend := "down"
		if row.Amount > previousTotals[row.CategoryName] {
			trend = "up"
		}
		items = append(items, billRankingItem{
			Name:    row.CategoryName,
			Value:   round2(row.Amount),
			Percent: percent,
			Count:   row.RecordCount,
			Trend:   trend,
			Badge:   badgeFromName(row.CategoryName),
		})
		if index >= 9 {
			break
		}
	}
	return items
}

func buildDailyTrend(rows []dailyExpenseRow) []billTrendPoint {
	items := make([]billTrendPoint, 0, len(rows))
	for _, row := range rows {
		items = append(items, billTrendPoint{
			Label:  row.DateKey,
			Amount: round2(row.Amount),
		})
	}
	return items
}

func buildMonthListHighlightCards(income float64, expense float64, balance float64, days int, records int) []billHighlightCard {
	coverage := 0.0
	if expense > 0 {
		coverage = round2(income / expense * 100)
	}
	return []billHighlightCard{
		{
			Label: "月结余",
			Value: formatMoney(balance),
			Hint:  "收入减去支出后的结果",
		},
		{
			Label: "记账密度",
			Value: fmt.Sprintf("%d天 · %d笔", days, records),
			Hint:  "当月实际发生的记录情况",
		},
		{
			Label: "收入覆盖",
			Value: fmt.Sprintf("%.0f%%", coverage),
			Hint:  "收入对支出的覆盖比例",
		},
	}
}

func buildMonthHighlightCards(totalExpense float64, trend []billTrendPoint, ranking []billRankingItem, income float64) []billHighlightCard {
	peakValue := "¥0"
	peakHint := "本月暂无支出记录"
	if len(trend) > 0 {
		maxItem := trend[0]
		for _, item := range trend[1:] {
			if item.Amount > maxItem.Amount {
				maxItem = item
			}
		}
		peakValue = formatMoney(maxItem.Amount)
		peakHint = maxItem.Label
	}

	topCategoryValue := "暂无"
	topCategoryHint := "本月暂无支出分类"
	if len(ranking) > 0 {
		topCategoryValue = ranking[0].Name
		topCategoryHint = fmt.Sprintf("%d 笔 · %s", ranking[0].Count, formatMoney(ranking[0].Value))
	}

	balanceRate := 0.0
	if income > 0 {
		balanceRate = round2((income - totalExpense) / income * 100)
	}

	return []billHighlightCard{
		{
			Label: "支出峰值",
			Value: peakValue,
			Hint:  peakHint,
		},
		{
			Label: "高频分类",
			Value: topCategoryValue,
			Hint:  topCategoryHint,
		},
		{
			Label: "结余率",
			Value: fmt.Sprintf("%.0f%%", balanceRate),
			Hint:  "以收入为基准计算当前结余占比",
		},
	}
}

func buildMonthAchievements(monthStart time.Time, recordedDays int, records int, ranking []billRankingItem) []billAchievement {
	consecutiveDays := recordedDays
	if consecutiveDays > daysInMonth(monthStart) {
		consecutiveDays = daysInMonth(monthStart)
	}

	topCategory := "暂无"
	topHint := "本月暂无支出分类"
	if len(ranking) > 0 {
		topCategory = ranking[0].Name
		topHint = fmt.Sprintf("占比 %.1f%%", ranking[0].Percent)
	}

	return []billAchievement{
		{
			Title: "连续记账",
			Value: fmt.Sprintf("%d天", consecutiveDays),
			Tone:  "success",
			Hint:  "按本月连续记账天数估算",
		},
		{
			Title: "本月记账",
			Value: fmt.Sprintf("%d笔", records),
			Tone:  "brand",
			Hint:  "本月全部收支记录总数",
		},
		{
			Title: "记账天数",
			Value: fmt.Sprintf("%d/%d天", recordedDays, daysInMonth(monthStart)),
			Tone:  "info",
			Hint:  "本月有记录的实际天数",
		},
		{
			Title: "高频分类",
			Value: topCategory,
			Tone:  "warning",
			Hint:  topHint,
		},
	}
}

func buildYearInsight(income float64, expense float64, balance float64, months int) string {
	if months == 0 {
		return "当前年份暂无账单记录。"
	}
	if balance >= 0 {
		return "当前年度收入覆盖支出，整体保持结余。"
	}
	return "当前年度支出高于收入，需要重点关注高波动月份。"
}

func buildMonthNote(income float64, expense float64, days int, records int) string {
	if records == 0 {
		return "当前月份暂无账单记录。"
	}
	if income >= expense {
		return fmt.Sprintf("本月记账 %d 天，共 %d 笔，收入覆盖支出。", days, records)
	}
	return fmt.Sprintf("本月记账 %d 天，共 %d 笔，支出高于收入。", days, records)
}

func buildMonthHighlight(income float64, expense float64, balance float64) string {
	if income == 0 && expense == 0 {
		return "当前月份暂无可展示的收支趋势。"
	}
	if balance >= 0 {
		return "本月整体保持结余，可继续观察高频支出分类。"
	}
	return "本月结余为负，建议优先收口高频支出分类。"
}

func buildAdminMonthInsight(income float64, expense float64, records int) string {
	if records == 0 {
		return "该月份暂无账单记录。"
	}
	if income >= expense {
		return "该月份收入覆盖支出。"
	}
	return "该月份支出高于收入。"
}

func adjacentMonthKeys(monthKeys []string, current string) (string, string) {
	for index, item := range monthKeys {
		if item != current {
			continue
		}
		previousKey := ""
		nextKey := ""
		if index < len(monthKeys)-1 {
			previousKey = monthKeys[index+1]
		}
		if index > 0 {
			nextKey = monthKeys[index-1]
		}
		return previousKey, nextKey
	}
	return "", ""
}

func monthNumberFromKey(key string) int {
	parts := strings.Split(key, "-")
	if len(parts) != 2 {
		return 0
	}
	value, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0
	}
	return value
}

func yearNumberFromKey(key string) int {
	parts := strings.Split(key, "-")
	if len(parts) == 0 {
		return 0
	}
	value, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}
	return value
}

func daysInMonth(monthStart time.Time) int {
	return monthStart.AddDate(0, 1, -1).Day()
}
