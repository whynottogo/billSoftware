var CHART_YEARS = [2026, 2025, 2024];

var EXPENSE_DATA_BY_YEAR = {
  2026: {
    summary: {
      yearlyExpense: 96840,
      monthlyAverage: 8070,
      records: 412
    },
    monthTrend: [620, 780, 540, 830, 910, 640, 700, 840, 760, 690, 880, 720],
    yearTrend: [7320, 6880, 7410, 8932, 8120, 7660, 7850, 8410, 8230, 7920, 8040, 8068],
    ranking: [
      { name: "餐饮", badge: "餐", value: 23240 },
      { name: "购物", badge: "购", value: 18450 },
      { name: "交通", badge: "交", value: 12180 },
      { name: "娱乐", badge: "娱", value: 10380 },
      { name: "居家", badge: "居", value: 7980 },
      { name: "医疗", badge: "医", value: 6760 },
      { name: "社交", badge: "社", value: 5920 },
      { name: "宠物", badge: "宠", value: 4470 },
      { name: "学习", badge: "学", value: 3830 },
      { name: "旅行", badge: "旅", value: 3630 }
    ]
  },
  2025: {
    summary: {
      yearlyExpense: 88520,
      monthlyAverage: 7377,
      records: 366
    },
    monthTrend: [520, 690, 560, 710, 800, 620, 680, 760, 690, 640, 770, 660],
    yearTrend: [6610, 6360, 6920, 7180, 7320, 7450, 7130, 7610, 7580, 7340, 7460, 7560],
    ranking: [
      { name: "餐饮", badge: "餐", value: 21460 },
      { name: "购物", badge: "购", value: 16980 },
      { name: "交通", badge: "交", value: 10930 },
      { name: "娱乐", badge: "娱", value: 9420 },
      { name: "居家", badge: "居", value: 7050 },
      { name: "医疗", badge: "医", value: 6110 },
      { name: "社交", badge: "社", value: 5200 },
      { name: "宠物", badge: "宠", value: 3840 },
      { name: "学习", badge: "学", value: 3450 },
      { name: "旅行", badge: "旅", value: 3080 }
    ]
  },
  2024: {
    summary: {
      yearlyExpense: 80130,
      monthlyAverage: 6678,
      records: 338
    },
    monthTrend: [480, 610, 520, 660, 740, 560, 620, 690, 630, 590, 700, 610],
    yearTrend: [6020, 5800, 6280, 6540, 6710, 6950, 6630, 7080, 6920, 6740, 6860, 7600],
    ranking: [
      { name: "餐饮", badge: "餐", value: 19340 },
      { name: "购物", badge: "购", value: 15360 },
      { name: "交通", badge: "交", value: 10120 },
      { name: "娱乐", badge: "娱", value: 8500 },
      { name: "居家", badge: "居", value: 6580 },
      { name: "医疗", badge: "医", value: 5440 },
      { name: "社交", badge: "社", value: 4920 },
      { name: "宠物", badge: "宠", value: 3560 },
      { name: "学习", badge: "学", value: 3090 },
      { name: "旅行", badge: "旅", value: 3220 }
    ]
  }
};

var INCOME_DATA_BY_YEAR = {
  2026: {
    summary: {
      yearlyIncome: 196300,
      monthlyAverage: 16358,
      records: 128
    },
    yearTrend: [14800, 15600, 16100, 18500, 16200, 15800, 16600, 17100, 16800, 16000, 17300, 17500],
    ranking: [
      { name: "工资", badge: "工", value: 145000 },
      { name: "兼职", badge: "兼", value: 18200 },
      { name: "奖金", badge: "奖", value: 14100 },
      { name: "投资", badge: "投", value: 10800 },
      { name: "稿费", badge: "稿", value: 3200 },
      { name: "返现", badge: "返", value: 2100 },
      { name: "租金", badge: "租", value: 1200 },
      { name: "礼金", badge: "礼", value: 840 },
      { name: "退税", badge: "税", value: 620 },
      { name: "其他", badge: "其", value: 240 }
    ]
  },
  2025: {
    summary: {
      yearlyIncome: 183400,
      monthlyAverage: 15283,
      records: 116
    },
    yearTrend: [13900, 14500, 15200, 16800, 15100, 14800, 15400, 15900, 15700, 15100, 16200, 16800],
    ranking: [
      { name: "工资", badge: "工", value: 136800 },
      { name: "兼职", badge: "兼", value: 16200 },
      { name: "奖金", badge: "奖", value: 12300 },
      { name: "投资", badge: "投", value: 9800 },
      { name: "稿费", badge: "稿", value: 2800 },
      { name: "返现", badge: "返", value: 1800 },
      { name: "租金", badge: "租", value: 900 },
      { name: "礼金", badge: "礼", value: 640 },
      { name: "退税", badge: "税", value: 490 },
      { name: "其他", badge: "其", value: 470 }
    ]
  },
  2024: {
    summary: {
      yearlyIncome: 169800,
      monthlyAverage: 14150,
      records: 104
    },
    yearTrend: [12800, 13200, 13800, 15200, 14100, 13600, 14500, 14900, 14600, 14200, 15100, 15800],
    ranking: [
      { name: "工资", badge: "工", value: 126700 },
      { name: "兼职", badge: "兼", value: 14300 },
      { name: "奖金", badge: "奖", value: 10900 },
      { name: "投资", badge: "投", value: 8600 },
      { name: "稿费", badge: "稿", value: 2400 },
      { name: "返现", badge: "返", value: 1600 },
      { name: "租金", badge: "租", value: 820 },
      { name: "礼金", badge: "礼", value: 560 },
      { name: "退税", badge: "税", value: 420 },
      { name: "其他", badge: "其", value: 500 }
    ]
  }
};

function clone(value) {
  return JSON.parse(JSON.stringify(value));
}

function formatChartCurrency(value) {
  return "¥" + Number(value || 0).toLocaleString("zh-CN");
}

function withPercent(list) {
  var total = list.reduce(function(sum, item) {
    return sum + Number(item.value || 0);
  }, 0);

  if (!total) {
    return list.map(function(item) {
      return Object.assign({}, item, { percent: 0 });
    });
  }

  return list.map(function(item) {
    return Object.assign({}, item, {
      percent: Math.round((Number(item.value || 0) / total) * 1000) / 10
    });
  });
}

function getChartYears() {
  return CHART_YEARS.slice();
}

function getExpenseChartYear(year) {
  var targetYear = Number(year || CHART_YEARS[0]);
  var data = EXPENSE_DATA_BY_YEAR[targetYear] || EXPENSE_DATA_BY_YEAR[CHART_YEARS[0]];
  var cloned = clone(data);
  cloned.year = targetYear;
  cloned.ranking = withPercent(cloned.ranking);
  return cloned;
}

function getIncomeChartYear(year) {
  var targetYear = Number(year || CHART_YEARS[0]);
  var data = INCOME_DATA_BY_YEAR[targetYear] || INCOME_DATA_BY_YEAR[CHART_YEARS[0]];
  var cloned = clone(data);
  cloned.year = targetYear;
  cloned.ranking = withPercent(cloned.ranking);
  return cloned;
}

export {
  formatChartCurrency,
  getChartYears,
  getExpenseChartYear,
  getIncomeChartYear
};
