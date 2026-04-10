import request from "@/utils/request";

function extractPayload(result) {
  if (result && result.data) {
    return result.data;
  }

  return result || {};
}

function toNumber(value) {
  var normalized = Number(value || 0);

  if (!Number.isFinite(normalized)) {
    return 0;
  }

  return normalized;
}

function toYear(value, fallbackYear) {
  var year = Number(value || fallbackYear || new Date().getFullYear());

  if (!Number.isFinite(year) || year <= 0) {
    return Number(fallbackYear || new Date().getFullYear());
  }

  return year;
}

function buildBadge(name) {
  return String(name || "未").slice(0, 1);
}

function normalizeTrendItem(item, index, suffix) {
  if (typeof item === "number") {
    return {
      label: index + 1 + suffix,
      value: toNumber(item)
    };
  }

  var source = item || {};
  var day = toNumber(source.day || source.date || source.index);
  var month = toNumber(source.month || source.period || source.order);
  var label = source.label || source.name || "";

  if (!label) {
    if (suffix === "日" && day) {
      label = day + "日";
    } else if (suffix === "月" && month) {
      label = month + "月";
    } else {
      label = index + 1 + suffix;
    }
  }

  return {
    label: label,
    value: toNumber(source.value || source.amount || source.total || source.expense || source.income)
  };
}

function normalizeTrend(list, suffix) {
  var source = Array.isArray(list) ? list : [];

  return source.map(function(item, index) {
    return normalizeTrendItem(item, index, suffix);
  });
}

function normalizeRanking(list) {
  var source = Array.isArray(list) ? list : [];
  var normalized = source.map(function(item) {
    var rankingItem = item || {};
    var name = rankingItem.name || rankingItem.category_name || rankingItem.categoryName || rankingItem.label || "未命名";

    return {
      name: name,
      badge: rankingItem.badge || rankingItem.category_badge || rankingItem.categoryBadge || buildBadge(name),
      value: toNumber(rankingItem.value || rankingItem.amount || rankingItem.total)
    };
  });
  var total = normalized.reduce(function(sum, item) {
    return sum + item.value;
  }, 0);

  return normalized.slice(0, 10).map(function(item, index) {
    var sourceItem = source[index] || {};
    var percent = sourceItem.percent;

    if (percent === undefined || percent === null || percent === "") {
      percent = total ? Math.round((item.value / total) * 1000) / 10 : 0;
    }

    return Object.assign({}, item, {
      percent: toNumber(percent)
    });
  });
}

function normalizeExpenseSummary(summary, monthTrend, yearTrend, ranking) {
  var source = summary || {};
  var yearlyExpense = toNumber(
    source.yearlyExpense || source.totalExpense || source.expenseTotal || source.yearTotal || source.total
  );
  var records = toNumber(source.records || source.recordCount || source.count);
  var monthlyAverage = toNumber(
    source.monthlyAverage || source.avgMonthlyExpense || source.average || source.monthAverage
  );

  if (!monthlyAverage) {
    var months = yearTrend.length || 12;
    monthlyAverage = months ? Math.round((yearlyExpense / months) * 100) / 100 : 0;
  }

  return {
    yearlyExpense: yearlyExpense,
    monthlyAverage: monthlyAverage,
    records: records,
    peakMonthExpense: toNumber(source.peakMonthExpense),
    peakDayExpense: toNumber(source.peakDayExpense),
    topCategory: source.topCategory || (ranking[0] && ranking[0].name) || "暂无数据",
    activeDays: toNumber(source.activeDays || monthTrend.filter(function(item) {
      return item.value > 0;
    }).length)
  };
}

function normalizeIncomeSummary(summary, yearTrend, ranking) {
  var source = summary || {};
  var yearlyIncome = toNumber(
    source.yearlyIncome || source.totalIncome || source.incomeTotal || source.yearTotal || source.total
  );
  var records = toNumber(source.records || source.recordCount || source.count);
  var monthlyAverage = toNumber(
    source.monthlyAverage || source.avgMonthlyIncome || source.average || source.monthAverage
  );

  if (!monthlyAverage) {
    var months = yearTrend.length || 12;
    monthlyAverage = months ? Math.round((yearlyIncome / months) * 100) / 100 : 0;
  }

  return {
    yearlyIncome: yearlyIncome,
    monthlyAverage: monthlyAverage,
    records: records,
    bestMonthIncome: toNumber(source.bestMonthIncome),
    topCategory: source.topCategory || (ranking[0] && ranking[0].name) || "暂无数据"
  };
}

export function getUserChartYears() {
  return request.get("/user/charts/years");
}

export function getUserExpenseChartYear(year) {
  return request.get("/user/charts/expense/" + year);
}

export function getUserIncomeChartYear(year) {
  return request.get("/user/charts/income/" + year);
}

export function formatChartCurrency(value) {
  return "¥" + Number(value || 0).toLocaleString("zh-CN");
}

export function buildUserChartsError(error, fallback) {
  if (error && error.response && error.response.data && error.response.data.message) {
    return error.response.data.message;
  }

  return fallback;
}

export function normalizeChartYearsPayload(result) {
  var payload = extractPayload(result);
  var years = Array.isArray(payload) ? payload : payload.years;

  return Array.from(new Set((Array.isArray(years) ? years : []).map(function(item) {
    return Number(item);
  }).filter(function(item) {
    return Number.isFinite(item);
  }))).sort(function(left, right) {
    return right - left;
  });
}

export function createEmptyExpenseChartData(year) {
  var normalizedYear = toYear(year);

  return {
    year: normalizedYear,
    summary: normalizeExpenseSummary({}, [], [], []),
    monthTrend: [],
    yearTrend: [],
    ranking: []
  };
}

export function createEmptyIncomeChartData(year) {
  var normalizedYear = toYear(year);

  return {
    year: normalizedYear,
    summary: normalizeIncomeSummary({}, [], []),
    yearTrend: [],
    ranking: []
  };
}

export function normalizeExpenseChartPayload(result, fallbackYear) {
  var payload = extractPayload(result);
  var year = toYear(payload.year, fallbackYear);
  var monthTrend = normalizeTrend(payload.monthTrend, "日");
  var yearTrend = normalizeTrend(payload.yearTrend, "月");
  var ranking = normalizeRanking(payload.ranking);

  return {
    year: year,
    summary: normalizeExpenseSummary(payload.summary, monthTrend, yearTrend, ranking),
    monthTrend: monthTrend,
    yearTrend: yearTrend,
    ranking: ranking
  };
}

export function normalizeIncomeChartPayload(result, fallbackYear) {
  var payload = extractPayload(result);
  var year = toYear(payload.year, fallbackYear);
  var yearTrend = normalizeTrend(payload.yearTrend, "月");
  var ranking = normalizeRanking(payload.ranking);

  return {
    year: year,
    summary: normalizeIncomeSummary(payload.summary, yearTrend, ranking),
    yearTrend: yearTrend,
    ranking: ranking
  };
}
