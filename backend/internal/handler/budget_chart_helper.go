package handler

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"xorm.io/xorm"

	"billsoftware/backend/internal/model"
)

const (
	budgetTypeMonth = "month"
	budgetTypeYear  = "year"
)

var (
	errBudgetInvalidCategory = errors.New("category_id must belong to current user expense categories")
	errBudgetInvalidAmount   = errors.New("budget amount must be greater than or equal to 0")
)

type budgetChartCore struct {
	engine *xorm.Engine
}

type budgetOverview struct {
	TotalBudget  float64 `json:"totalBudget"`
	TotalExpense float64 `json:"totalExpense"`
	Remaining    float64 `json:"remaining"`
	Percentage   float64 `json:"percentage"`
}

type budgetHighlight struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Hint  string `json:"hint"`
}

type budgetCategoryCard struct {
	CategoryID uint64  `json:"categoryId"`
	Name       string  `json:"name"`
	Badge      string  `json:"badge"`
	Budget     float64 `json:"budget"`
	Expense    float64 `json:"expense"`
	Remaining  float64 `json:"remaining"`
	Percentage float64 `json:"percentage"`
	Status     string  `json:"status"`
	Note       string  `json:"note"`
	IsDefault  bool    `json:"isDefault"`
}

type monthBudgetResponse struct {
	Month      string               `json:"month"`
	Label      string               `json:"label"`
	Notice     string               `json:"notice"`
	Overview   budgetOverview       `json:"overview"`
	Highlights []budgetHighlight    `json:"highlights"`
	Categories []budgetCategoryCard `json:"categories"`
}

type yearBudgetOverview struct {
	TotalBudget  float64 `json:"totalBudget"`
	TotalExpense float64 `json:"totalExpense"`
	Remaining    float64 `json:"remaining"`
	Percentage   float64 `json:"percentage"`
	Months       int     `json:"months"`
	Note         string  `json:"note"`
}

type budgetExecutionItem struct {
	Label   string  `json:"label"`
	Budget  float64 `json:"budget"`
	Expense float64 `json:"expense"`
}

type yearBudgetResponse struct {
	Year             int                   `json:"year"`
	Overview         yearBudgetOverview    `json:"overview"`
	MonthlyExecution []budgetExecutionItem `json:"monthlyExecution"`
	Categories       []budgetCategoryCard  `json:"categories"`
}

type chartTrendPoint struct {
	Label  string  `json:"label"`
	Amount float64 `json:"amount"`
}

type expenseChartSummary struct {
	YearlyExpense  float64 `json:"yearlyExpense"`
	MonthlyAverage float64 `json:"monthlyAverage"`
	Records        int     `json:"records"`
	MonthKey       string  `json:"monthKey"`
	MonthExpense   float64 `json:"monthExpense"`
}

type incomeChartSummary struct {
	YearlyIncome   float64 `json:"yearlyIncome"`
	MonthlyAverage float64 `json:"monthlyAverage"`
	Records        int     `json:"records"`
}

type expenseChartResponse struct {
	Year       int                 `json:"year"`
	Summary    expenseChartSummary `json:"summary"`
	MonthTrend []chartTrendPoint   `json:"monthTrend"`
	YearTrend  []chartTrendPoint   `json:"yearTrend"`
	Ranking    []billRankingItem   `json:"ranking"`
}

type incomeChartResponse struct {
	Year      int                `json:"year"`
	Summary   incomeChartSummary `json:"summary"`
	YearTrend []chartTrendPoint  `json:"yearTrend"`
	Ranking   []billRankingItem  `json:"ranking"`
}

type budgetAmountRow struct {
	CategoryID uint64  `xorm:"category_id"`
	Amount     float64 `xorm:"amount"`
}

type chartMonthAggregateRow struct {
	MonthNo int     `xorm:"month_no"`
	Amount  float64 `xorm:"amount"`
}

type chartDailyAggregateRow struct {
	Label  string  `xorm:"label"`
	Amount float64 `xorm:"amount"`
}

type chartYearOptionRow struct {
	YearNo int `xorm:"year_no"`
}

func newBudgetChartCore(engine *xorm.Engine) *budgetChartCore {
	return &budgetChartCore{engine: engine}
}

func (c *budgetChartCore) GetMonthBudget(userID uint64, now time.Time) (monthBudgetResponse, error) {
	month := beginningOfMonth(now)
	monthKey := month.Format("2006-01")
	response := monthBudgetResponse{
		Month:      monthKey,
		Label:      fmt.Sprintf("%s预算", monthLabel(monthKey)),
		Notice:     "仅支持设置当前月份预算，历史月份只读。",
		Highlights: make([]budgetHighlight, 0),
		Categories: make([]budgetCategoryCard, 0),
	}

	budget, budgetItems, err := c.loadBudgetWithItems(userID, budgetTypeMonth, monthKey)
	if err != nil {
		return response, err
	}

	expenseTotal, expenseByCategory, err := c.loadExpenseTotalsByCategory(userID, month, month.AddDate(0, 1, 0))
	if err != nil {
		return response, err
	}

	categories, err := c.loadExpenseCategories(userID)
	if err != nil {
		return response, err
	}

	response.Categories = buildBudgetCategoryCards(categories, budgetItems, expenseByCategory, budget)
	response.Overview = buildMonthBudgetOverview(budget.TotalAmount, expenseTotal)
	response.Highlights = buildMonthBudgetHighlights(response.Categories, response.Overview.Percentage)

	return response, nil
}

func (c *budgetChartCore) UpsertCurrentMonthBudget(userID uint64, totalAmount *float64, inputs *[]budgetCategoryInput) error {
	now := time.Now()
	return c.upsertBudget(userID, budgetTypeMonth, beginningOfMonth(now).Format("2006-01"), beginningOfMonth(now), beginningOfMonth(now).AddDate(0, 1, 0), totalAmount, inputs)
}

func (c *budgetChartCore) ListBudgetYears(userID uint64) ([]int, error) {
	rows := make([]chartYearOptionRow, 0)
	if err := c.engine.SQL(`
		SELECT year_no
		FROM (
			SELECT CAST(period_key AS UNSIGNED) AS year_no
			FROM budgets
			WHERE user_id = ? AND budget_type = 'year'
			UNION
			SELECT CAST(LEFT(period_key, 4) AS UNSIGNED) AS year_no
			FROM budgets
			WHERE user_id = ? AND budget_type = 'month'
			UNION
			SELECT CAST(DATE_FORMAT(record_date, '%Y') AS UNSIGNED) AS year_no
			FROM ledger_records
			WHERE user_id = ?
		) years
		WHERE year_no IS NOT NULL AND year_no > 0
		ORDER BY year_no DESC
	`, userID, userID, userID).Find(&rows); err != nil {
		return nil, err
	}

	years := make([]int, 0, len(rows)+1)
	seen := map[int]bool{}
	currentYear := time.Now().Year()
	years = append(years, currentYear)
	seen[currentYear] = true
	for _, row := range rows {
		if row.YearNo <= 0 || seen[row.YearNo] {
			continue
		}
		years = append(years, row.YearNo)
		seen[row.YearNo] = true
	}
	sort.Slice(years, func(i, j int) bool { return years[i] > years[j] })
	return years, nil
}

func (c *budgetChartCore) GetYearBudget(userID uint64, year int) (yearBudgetResponse, error) {
	response := yearBudgetResponse{
		Year:             year,
		MonthlyExecution: make([]budgetExecutionItem, 0, 12),
		Categories:       make([]budgetCategoryCard, 0),
	}

	yearKey := strconv.Itoa(year)
	budget, budgetItems, err := c.loadBudgetWithItems(userID, budgetTypeYear, yearKey)
	if err != nil {
		return response, err
	}

	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(1, 0, 0)
	expenseTotal, expenseByCategory, err := c.loadExpenseTotalsByCategory(userID, start, end)
	if err != nil {
		return response, err
	}

	categories, err := c.loadExpenseCategories(userID)
	if err != nil {
		return response, err
	}

	monthExpenseMap, monthsWithExpense, err := c.loadMonthlyExpenseMap(userID, year)
	if err != nil {
		return response, err
	}

	monthBudgetMap, err := c.loadMonthBudgetMap(userID, year)
	if err != nil {
		return response, err
	}

	defaultMonthBudget := 0.0
	if budget.TotalAmount > 0 {
		defaultMonthBudget = round2(budget.TotalAmount / 12)
	}

	for month := 1; month <= 12; month += 1 {
		monthKey := fmt.Sprintf("%04d-%02d", year, month)
		budgetValue := defaultMonthBudget
		if override, ok := monthBudgetMap[monthKey]; ok {
			budgetValue = override
		}
		response.MonthlyExecution = append(response.MonthlyExecution, budgetExecutionItem{
			Label:   fmt.Sprintf("%d月", month),
			Budget:  round2(budgetValue),
			Expense: round2(monthExpenseMap[month]),
		})
	}

	response.Categories = buildBudgetCategoryCards(categories, budgetItems, expenseByCategory, budget)
	response.Overview = yearBudgetOverview{
		TotalBudget:  round2(budget.TotalAmount),
		TotalExpense: round2(expenseTotal),
		Remaining:    round2(budget.TotalAmount - expenseTotal),
		Percentage:   budgetUsagePercent(budget.TotalAmount, expenseTotal),
		Months:       monthsWithExpense,
		Note:         buildYearBudgetNote(budget.TotalAmount, expenseTotal, monthsWithExpense),
	}

	return response, nil
}

func (c *budgetChartCore) UpsertCurrentYearBudget(userID uint64, totalAmount *float64, inputs *[]budgetCategoryInput) error {
	now := time.Now()
	start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	return c.upsertBudget(userID, budgetTypeYear, strconv.Itoa(now.Year()), start, start.AddDate(1, 0, 0), totalAmount, inputs)
}

func (c *budgetChartCore) ListChartYears(userID uint64) ([]int, error) {
	rows := make([]chartYearOptionRow, 0)
	if err := c.engine.SQL(`
		SELECT CAST(DATE_FORMAT(record_date, '%Y') AS UNSIGNED) AS year_no
		FROM ledger_records
		WHERE user_id = ?
		GROUP BY DATE_FORMAT(record_date, '%Y')
		ORDER BY year_no DESC
	`, userID).Find(&rows); err != nil {
		return nil, err
	}

	years := make([]int, 0, len(rows)+1)
	currentYear := time.Now().Year()
	years = append(years, currentYear)
	seen := map[int]bool{currentYear: true}
	for _, row := range rows {
		if row.YearNo <= 0 || seen[row.YearNo] {
			continue
		}
		years = append(years, row.YearNo)
		seen[row.YearNo] = true
	}
	sort.Slice(years, func(i, j int) bool { return years[i] > years[j] })
	return years, nil
}

func (c *budgetChartCore) GetExpenseChart(userID uint64, year int) (expenseChartResponse, error) {
	response := expenseChartResponse{
		Year:       year,
		MonthTrend: make([]chartTrendPoint, 0),
		YearTrend:  make([]chartTrendPoint, 0, 12),
		Ranking:    make([]billRankingItem, 0),
	}

	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(1, 0, 0)
	yearlyExpense, records, err := c.loadTypePeriodSummary(userID, recordTypeExpense, start, end)
	if err != nil {
		return response, err
	}

	monthRows, err := c.loadYearTrendRows(userID, recordTypeExpense, year)
	if err != nil {
		return response, err
	}
	response.YearTrend = buildYearTrend(monthRows)

	categoryRows, err := c.loadCategoryRankingRows(userID, recordTypeExpense, start, end)
	if err != nil {
		return response, err
	}
	response.Ranking = buildBudgetChartRanking(categoryRows, yearlyExpense)

	monthKey, monthExpense, err := c.findExpenseTrendMonth(userID, year)
	if err != nil {
		return response, err
	}
	if monthKey != "" {
		response.MonthTrend, err = c.loadExpenseMonthTrend(userID, monthKey)
		if err != nil {
			return response, err
		}
	}

	response.Summary = expenseChartSummary{
		YearlyExpense:  round2(yearlyExpense),
		MonthlyAverage: round2(yearlyExpense / 12),
		Records:        records,
		MonthKey:       monthKey,
		MonthExpense:   round2(monthExpense),
	}

	return response, nil
}

func (c *budgetChartCore) GetIncomeChart(userID uint64, year int) (incomeChartResponse, error) {
	response := incomeChartResponse{
		Year:      year,
		YearTrend: make([]chartTrendPoint, 0, 12),
		Ranking:   make([]billRankingItem, 0),
	}

	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(1, 0, 0)
	yearlyIncome, records, err := c.loadTypePeriodSummary(userID, recordTypeIncome, start, end)
	if err != nil {
		return response, err
	}

	monthRows, err := c.loadYearTrendRows(userID, recordTypeIncome, year)
	if err != nil {
		return response, err
	}
	response.YearTrend = buildYearTrend(monthRows)

	categoryRows, err := c.loadCategoryRankingRows(userID, recordTypeIncome, start, end)
	if err != nil {
		return response, err
	}
	response.Ranking = buildBudgetChartRanking(categoryRows, yearlyIncome)

	response.Summary = incomeChartSummary{
		YearlyIncome:   round2(yearlyIncome),
		MonthlyAverage: round2(yearlyIncome / 12),
		Records:        records,
	}

	return response, nil
}

func (c *budgetChartCore) upsertBudget(userID uint64, budgetType string, periodKey string, start time.Time, end time.Time, totalAmount *float64, inputs *[]budgetCategoryInput) error {
	expenseTotal, expenseByCategory, err := c.loadExpenseTotalsByCategory(userID, start, end)
	if err != nil {
		return err
	}

	validCategories, err := c.loadExpenseCategoryMap(userID)
	if err != nil {
		return err
	}

	normalizedItems := make([]budgetCategoryInput, 0)
	if inputs != nil {
		seen := make(map[uint64]bool)
		for _, item := range *inputs {
			if item.Amount < 0 {
				return errBudgetInvalidAmount
			}
			if item.CategoryID == 0 || validCategories[item.CategoryID] == nil {
				return errBudgetInvalidCategory
			}
			if seen[item.CategoryID] {
				continue
			}
			seen[item.CategoryID] = true
			normalizedItems = append(normalizedItems, budgetCategoryInput{
				CategoryID: item.CategoryID,
				Amount:     round2(item.Amount),
			})
		}
	}

	session := c.engine.NewSession()
	defer func() {
		_ = session.Close()
	}()

	if err := session.Begin(); err != nil {
		return err
	}

	budget := &model.Budget{}
	has, err := session.Where("user_id = ? AND budget_type = ? AND period_key = ?", userID, budgetType, periodKey).Get(budget)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	if !has {
		budget = &model.Budget{
			UserID:      userID,
			BudgetType:  budgetType,
			PeriodKey:   periodKey,
			TotalAmount: 0,
			UsedAmount:  round2(expenseTotal),
		}
		if totalAmount != nil {
			budget.TotalAmount = round2(*totalAmount)
		}
		if _, err := session.Insert(budget); err != nil {
			_ = session.Rollback()
			return err
		}
	} else {
		if totalAmount != nil {
			budget.TotalAmount = round2(*totalAmount)
		}
		budget.UsedAmount = round2(expenseTotal)
		if _, err := session.ID(budget.ID).Cols("total_amount", "used_amount").Update(budget); err != nil {
			_ = session.Rollback()
			return err
		}
	}

	if inputs != nil {
		if _, err := session.Where("budget_id = ?", budget.ID).Delete(&model.BudgetItem{}); err != nil {
			_ = session.Rollback()
			return err
		}

		if len(normalizedItems) > 0 {
			items := make([]model.BudgetItem, 0, len(normalizedItems))
			for _, item := range normalizedItems {
				items = append(items, model.BudgetItem{
					BudgetID:   budget.ID,
					CategoryID: item.CategoryID,
					Amount:     round2(item.Amount),
					UsedAmount: round2(expenseByCategory[item.CategoryID]),
				})
			}
			if _, err := session.Insert(&items); err != nil {
				_ = session.Rollback()
				return err
			}
		}
	} else {
		existingItems := make([]model.BudgetItem, 0)
		if err := session.Where("budget_id = ?", budget.ID).Find(&existingItems); err != nil {
			_ = session.Rollback()
			return err
		}
		for _, item := range existingItems {
			item.UsedAmount = round2(expenseByCategory[item.CategoryID])
			if _, err := session.ID(item.ID).Cols("used_amount").Update(&item); err != nil {
				_ = session.Rollback()
				return err
			}
		}
	}

	if err := session.Commit(); err != nil {
		_ = session.Rollback()
		return err
	}
	return nil
}

func (c *budgetChartCore) loadExpenseCategories(userID uint64) ([]model.UserCategory, error) {
	categories := make([]model.UserCategory, 0)
	if err := c.engine.Where("user_id = ? AND category_type = ?", userID, recordTypeExpense).Asc("sort_order", "id").Find(&categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *budgetChartCore) loadExpenseCategoryMap(userID uint64) (map[uint64]*model.UserCategory, error) {
	categories, err := c.loadExpenseCategories(userID)
	if err != nil {
		return nil, err
	}
	result := make(map[uint64]*model.UserCategory, len(categories))
	for index := range categories {
		result[categories[index].ID] = &categories[index]
	}
	return result, nil
}

func (c *budgetChartCore) loadBudgetWithItems(userID uint64, budgetType string, periodKey string) (model.Budget, map[uint64]float64, error) {
	budget := model.Budget{
		UserID:      userID,
		BudgetType:  budgetType,
		PeriodKey:   periodKey,
		TotalAmount: 0,
		UsedAmount:  0,
	}

	has, err := c.engine.Where("user_id = ? AND budget_type = ? AND period_key = ?", userID, budgetType, periodKey).Get(&budget)
	if err != nil {
		return budget, nil, err
	}

	items := make(map[uint64]float64)
	if !has {
		return budget, items, nil
	}

	rows := make([]model.BudgetItem, 0)
	if err := c.engine.Where("budget_id = ?", budget.ID).Find(&rows); err != nil {
		return budget, nil, err
	}
	for _, row := range rows {
		items[row.CategoryID] = round2(row.Amount)
	}
	return budget, items, nil
}

func (c *budgetChartCore) loadExpenseTotalsByCategory(userID uint64, start time.Time, end time.Time) (float64, map[uint64]float64, error) {
	rows := make([]budgetAmountRow, 0)
	if err := c.engine.SQL(`
		SELECT category_id, COALESCE(SUM(amount), 0) AS amount
		FROM ledger_records
		WHERE user_id = ? AND record_type = 'expense' AND record_date >= ? AND record_date < ?
		GROUP BY category_id
	`, userID, start.Format("2006-01-02"), end.Format("2006-01-02")).Find(&rows); err != nil {
		return 0, nil, err
	}

	total := 0.0
	result := make(map[uint64]float64, len(rows))
	for _, row := range rows {
		value := round2(row.Amount)
		total += value
		result[row.CategoryID] = value
	}

	return round2(total), result, nil
}

func (c *budgetChartCore) loadMonthlyExpenseMap(userID uint64, year int) (map[int]float64, int, error) {
	rows := make([]chartMonthAggregateRow, 0)
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(1, 0, 0)
	if err := c.engine.SQL(`
		SELECT CAST(DATE_FORMAT(record_date, '%m') AS UNSIGNED) AS month_no, COALESCE(SUM(amount), 0) AS amount
		FROM ledger_records
		WHERE user_id = ? AND record_type = 'expense' AND record_date >= ? AND record_date < ?
		GROUP BY DATE_FORMAT(record_date, '%m')
	`, userID, start.Format("2006-01-02"), end.Format("2006-01-02")).Find(&rows); err != nil {
		return nil, 0, err
	}

	result := make(map[int]float64, len(rows))
	for _, row := range rows {
		result[row.MonthNo] = round2(row.Amount)
	}
	return result, len(rows), nil
}

func (c *budgetChartCore) loadMonthBudgetMap(userID uint64, year int) (map[string]float64, error) {
	rows := make([]model.Budget, 0)
	if err := c.engine.Where("user_id = ? AND budget_type = ? AND period_key >= ? AND period_key <= ?", userID, budgetTypeMonth, fmt.Sprintf("%04d-01", year), fmt.Sprintf("%04d-12", year)).Find(&rows); err != nil {
		return nil, err
	}

	result := make(map[string]float64, len(rows))
	for _, row := range rows {
		result[row.PeriodKey] = round2(row.TotalAmount)
	}
	return result, nil
}

func (c *budgetChartCore) loadTypePeriodSummary(userID uint64, recordType string, start time.Time, end time.Time) (float64, int, error) {
	var row struct {
		Amount  float64 `xorm:"amount"`
		Records int     `xorm:"records"`
	}
	_, err := c.engine.SQL(`
		SELECT COALESCE(SUM(amount), 0) AS amount, COUNT(*) AS records
		FROM ledger_records
		WHERE user_id = ? AND record_type = ? AND record_date >= ? AND record_date < ?
	`, userID, recordType, start.Format("2006-01-02"), end.Format("2006-01-02")).Get(&row)
	if err != nil {
		return 0, 0, err
	}
	return round2(row.Amount), row.Records, nil
}

func (c *budgetChartCore) loadYearTrendRows(userID uint64, recordType string, year int) ([]chartMonthAggregateRow, error) {
	rows := make([]chartMonthAggregateRow, 0)
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(1, 0, 0)
	if err := c.engine.SQL(`
		SELECT CAST(DATE_FORMAT(record_date, '%m') AS UNSIGNED) AS month_no, COALESCE(SUM(amount), 0) AS amount
		FROM ledger_records
		WHERE user_id = ? AND record_type = ? AND record_date >= ? AND record_date < ?
		GROUP BY DATE_FORMAT(record_date, '%m')
		ORDER BY month_no ASC
	`, userID, recordType, start.Format("2006-01-02"), end.Format("2006-01-02")).Find(&rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func (c *budgetChartCore) loadCategoryRankingRows(userID uint64, recordType string, start time.Time, end time.Time) ([]categoryAggregateRow, error) {
	rows := make([]categoryAggregateRow, 0)
	if err := c.engine.SQL(`
		SELECT
			COALESCE(uc.name, '未分类') AS category_name,
			COALESCE(SUM(lr.amount), 0) AS amount,
			COUNT(*) AS record_count
		FROM ledger_records lr
		LEFT JOIN user_categories uc ON uc.id = lr.category_id AND uc.user_id = lr.user_id
		WHERE lr.user_id = ? AND lr.record_type = ? AND lr.record_date >= ? AND lr.record_date < ?
		GROUP BY COALESCE(uc.name, '未分类')
		ORDER BY amount DESC, record_count DESC, category_name ASC
	`, userID, recordType, start.Format("2006-01-02"), end.Format("2006-01-02")).Find(&rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func (c *budgetChartCore) findExpenseTrendMonth(userID uint64, year int) (string, float64, error) {
	rows := make([]periodAggregateRow, 0)
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(1, 0, 0)
	if err := c.engine.SQL(`
		SELECT DATE_FORMAT(record_date, '%Y-%m') AS period_key, COALESCE(SUM(amount), 0) AS expense
		FROM ledger_records
		WHERE user_id = ? AND record_type = 'expense' AND record_date >= ? AND record_date < ?
		GROUP BY DATE_FORMAT(record_date, '%Y-%m')
		ORDER BY period_key DESC
		LIMIT 1
	`, userID, start.Format("2006-01-02"), end.Format("2006-01-02")).Find(&rows); err != nil {
		return "", 0, err
	}

	if len(rows) == 0 {
		return "", 0, nil
	}
	return rows[0].PeriodKey, round2(rows[0].Expense), nil
}

func (c *budgetChartCore) loadExpenseMonthTrend(userID uint64, monthKey string) ([]chartTrendPoint, error) {
	start, end, _, err := resolveMonthRange(monthKey)
	if err != nil {
		return nil, err
	}

	rows := make([]chartDailyAggregateRow, 0)
	if err := c.engine.SQL(`
		SELECT DATE_FORMAT(record_date, '%m-%d') AS label, COALESCE(SUM(amount), 0) AS amount
		FROM ledger_records
		WHERE user_id = ? AND record_type = 'expense' AND record_date >= ? AND record_date < ?
		GROUP BY DATE_FORMAT(record_date, '%m-%d')
		ORDER BY label ASC
	`, userID, start.Format("2006-01-02"), end.Format("2006-01-02")).Find(&rows); err != nil {
		return nil, err
	}

	points := make([]chartTrendPoint, 0, len(rows))
	for _, row := range rows {
		points = append(points, chartTrendPoint{
			Label:  row.Label,
			Amount: round2(row.Amount),
		})
	}
	return points, nil
}

func buildBudgetCategoryCards(categories []model.UserCategory, budgetByCategory map[uint64]float64, expenseByCategory map[uint64]float64, budget model.Budget) []budgetCategoryCard {
	if budgetByCategory == nil {
		budgetByCategory = make(map[uint64]float64)
	}
	if expenseByCategory == nil {
		expenseByCategory = make(map[uint64]float64)
	}

	items := make([]budgetCategoryCard, 0, len(categories))
	for _, category := range categories {
		budgetValue := round2(budgetByCategory[category.ID])
		expenseValue := round2(expenseByCategory[category.ID])
		remaining := round2(budgetValue - expenseValue)
		percentage := budgetUsagePercent(budgetValue, expenseValue)
		status := budgetStatus(budgetValue, expenseValue)
		items = append(items, budgetCategoryCard{
			CategoryID: category.ID,
			Name:       category.Name,
			Badge:      badgeFromName(category.Name),
			Budget:     budgetValue,
			Expense:    expenseValue,
			Remaining:  remaining,
			Percentage: percentage,
			Status:     status,
			Note:       buildBudgetCategoryNote(status, budgetValue, expenseValue),
			IsDefault:  category.IsSystem == 1,
		})
	}
	return items
}

func buildMonthBudgetOverview(totalBudget float64, totalExpense float64) budgetOverview {
	totalBudget = round2(totalBudget)
	totalExpense = round2(totalExpense)
	return budgetOverview{
		TotalBudget:  totalBudget,
		TotalExpense: totalExpense,
		Remaining:    round2(totalBudget - totalExpense),
		Percentage:   budgetUsagePercent(totalBudget, totalExpense),
	}
}

func buildMonthBudgetHighlights(categories []budgetCategoryCard, totalPercentage float64) []budgetHighlight {
	overCount := 0
	controllableCount := 0
	for _, item := range categories {
		if item.Status == "over" {
			overCount += 1
			continue
		}
		if item.Status == "warning" || item.Status == "safe" {
			controllableCount += 1
		}
	}

	paceValue := "预算未开始"
	paceHint := "当前月份还没有预算或支出数据。"
	switch {
	case totalPercentage >= 100:
		paceValue = "支出偏快"
		paceHint = "整体预算已经超出，需要先收口高频分类。"
	case totalPercentage >= 80:
		paceValue = "临近上限"
		paceHint = "本月预算进入预警区，新增支出要更谨慎。"
	case totalPercentage > 0:
		paceValue = "稳中可控"
		paceHint = "当前预算仍有余量，可以继续观察高波动分类。"
	}

	return []budgetHighlight{
		{
			Label: "已超支分类",
			Value: fmt.Sprintf("%d个", overCount),
			Hint:  "优先处理已经超出预算的支出分类。",
		},
		{
			Label: "可控分类",
			Value: fmt.Sprintf("%d个", controllableCount),
			Hint:  "仍处在安全或预警区间，可继续追踪。",
		},
		{
			Label: "预算节奏",
			Value: paceValue,
			Hint:  paceHint,
		},
	}
}

func buildYearBudgetNote(totalBudget float64, totalExpense float64, months int) string {
	switch {
	case totalBudget == 0 && totalExpense == 0:
		return "当前年度暂无预算设置，也暂无支出记录。"
	case totalBudget == 0 && totalExpense > 0:
		return "当前年度已有真实支出，但还没有设置年度预算。"
	case totalExpense > totalBudget:
		return "当前年度预算已被突破，建议优先检查高波动分类。"
	case months == 0:
		return "当前年度预算已设置，等待实际支出逐步进入。"
	default:
		return "年度预算与实际支出已形成闭环，可继续按月份跟踪节奏。"
	}
}

func buildBudgetChartRanking(rows []categoryAggregateRow, total float64) []billRankingItem {
	items := make([]billRankingItem, 0, len(rows))
	for index, row := range rows {
		percent := 0.0
		if total > 0 {
			percent = round2(row.Amount / total * 100)
		}
		items = append(items, billRankingItem{
			Name:    row.CategoryName,
			Value:   round2(row.Amount),
			Percent: percent,
			Count:   row.RecordCount,
			Trend:   "flat",
			Badge:   badgeFromName(row.CategoryName),
		})
		if index >= 9 {
			break
		}
	}
	return items
}

func buildYearTrend(rows []chartMonthAggregateRow) []chartTrendPoint {
	amounts := make(map[int]float64, len(rows))
	for _, row := range rows {
		amounts[row.MonthNo] = round2(row.Amount)
	}

	points := make([]chartTrendPoint, 0, 12)
	for month := 1; month <= 12; month += 1 {
		points = append(points, chartTrendPoint{
			Label:  fmt.Sprintf("%d月", month),
			Amount: round2(amounts[month]),
		})
	}
	return points
}

func budgetUsagePercent(budget float64, expense float64) float64 {
	if budget <= 0 {
		if expense <= 0 {
			return 0
		}
		return 100
	}
	return round2(expense / budget * 100)
}

func budgetStatus(budget float64, expense float64) string {
	switch {
	case expense > budget && expense > 0:
		return "over"
	case budget > 0 && expense >= budget*0.8:
		return "warning"
	default:
		return "safe"
	}
}

func buildBudgetCategoryNote(status string, budget float64, expense float64) string {
	switch {
	case budget <= 0 && expense > 0:
		return "当前分类尚未设置预算，但已经产生真实支出。"
	case budget <= 0 && expense <= 0:
		return "当前分类还没有预算设置，也没有产生支出。"
	case status == "over":
		return "当前分类已经超出预算，建议优先收口。"
	case status == "warning":
		return "当前分类接近预算上限，需要继续盯住。"
	case expense <= 0:
		return "当前分类预算已预留，本期暂无支出。"
	default:
		return "当前分类仍有预算余量，整体处于可控区间。"
	}
}

func beginningOfMonth(target time.Time) time.Time {
	return time.Date(target.Year(), target.Month(), 1, 0, 0, 0, 0, time.Local)
}
