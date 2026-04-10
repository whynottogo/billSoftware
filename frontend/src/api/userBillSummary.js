import request from "@/utils/request";

function padMonth(value) {
  return String(value).padStart(2, "0");
}

function extractPayload(result) {
  if (result && result.data) {
    return result.data;
  }

  return result || {};
}

function normalizeSummary(summary) {
  var source = summary || {};

  return {
    balance: Number(source.balance || 0),
    income: Number(source.income || 0),
    expense: Number(source.expense || 0),
    months: Number(source.months || 0),
    days: Number(source.days || 0),
    records: Number(source.records || 0),
    insight: source.insight || "当前年份的账单汇总正在同步。"
  };
}

function isCurrentMonthKey(monthKey) {
  var now = new Date();
  var currentKey = now.getFullYear() + "-" + padMonth(now.getMonth() + 1);

  return monthKey === currentKey;
}

function buildDefaultHighlightCards(month) {
  return [
    {
      label: "本月结余",
      value: formatBillCurrency(month.balance),
      hint: month.highlight || "可继续进入账单详情查看分类结构"
    },
    {
      label: "记账密度",
      value: month.days + "天",
      hint: month.records + " 笔记录"
    },
    {
      label: "月度状态",
      value: month.status,
      hint: month.note
    }
  ];
}

function normalizeHighlightCards(cards, month) {
  var source = Array.isArray(cards) ? cards : [];

  if (!source.length) {
    return buildDefaultHighlightCards(month);
  }

  return source.map(function(item, index) {
    var fallback = buildDefaultHighlightCards(month)[index] || buildDefaultHighlightCards(month)[0];

    return {
      label: item && item.label ? item.label : fallback.label,
      value: item && item.value ? item.value : fallback.value,
      hint: item && item.hint ? item.hint : fallback.hint
    };
  });
}

function normalizeMonth(item, fallbackYear) {
  var source = item || {};
  var year = Number(source.year || fallbackYear || new Date().getFullYear());
  var month = Number(source.month || String(source.key || "").slice(5, 7) || 1);
  var key = source.key || (year + "-" + padMonth(month));
  var status = source.status || (isCurrentMonthKey(key) ? "本月" : "已结账");
  var normalized = {
    key: key,
    label: source.label || (year + "年" + month + "月"),
    year: year,
    month: month,
    income: Number(source.income || 0),
    expense: Number(source.expense || 0),
    balance: Number(source.balance || 0),
    days: Number(source.days || 0),
    records: Number(source.records || 0),
    status: status,
    note: source.note || "当前月份暂无补充说明。",
    highlight: source.highlight || "可以进入月账单详情查看更细的结构。"
  };

  normalized.highlightCards = normalizeHighlightCards(source.highlightCards, normalized);

  return normalized;
}

function sortYearsDesc(left, right) {
  return Number(right) - Number(left);
}

function sortMonthsDesc(left, right) {
  return String(right.key || "").localeCompare(String(left.key || ""));
}

export function getUserBillYears() {
  return request.get("/user/bills/years");
}

export function getUserBillYear(year) {
  return request.get("/user/bills/year/" + year);
}

export function formatBillCurrency(value) {
  return "¥" + Number(value || 0).toLocaleString("zh-CN");
}

export function buildBillSummaryError(error, fallback) {
  if (error && error.response && error.response.data && error.response.data.message) {
    return error.response.data.message;
  }

  return fallback;
}

export function createEmptyBillYearData(year) {
  return {
    year: Number(year || new Date().getFullYear()),
    summary: normalizeSummary(),
    months: []
  };
}

export function normalizeBillYearsPayload(result) {
  var payload = extractPayload(result);
  var history = Array.isArray(payload.history) ? payload.history : [];
  var years = Array.isArray(payload.years) ? payload.years : [];
  var normalizedHistory = history
    .map(function(item) {
      return {
        year: Number(item && item.year),
        summary: normalizeSummary(item && item.summary)
      };
    })
    .filter(function(item) {
      return Number.isFinite(item.year);
    })
    .sort(function(left, right) {
      return right.year - left.year;
    });
  var normalizedYears = years
    .map(function(item) {
      return Number(item);
    })
    .filter(function(item) {
      return Number.isFinite(item);
    });

  if (!normalizedYears.length) {
    normalizedYears = normalizedHistory.map(function(item) {
      return item.year;
    });
  }

  normalizedYears = Array.from(new Set(normalizedYears)).sort(sortYearsDesc);

  return {
    years: normalizedYears,
    history: normalizedHistory
  };
}

export function normalizeBillYearPayload(result, fallbackYear) {
  var payload = extractPayload(result);
  var year = Number(payload.year || fallbackYear || new Date().getFullYear());
  var months = Array.isArray(payload.months) ? payload.months : [];

  return {
    year: year,
    summary: normalizeSummary(payload.summary),
    months: months.map(function(item) {
      return normalizeMonth(item, year);
    }).sort(sortMonthsDesc)
  };
}
