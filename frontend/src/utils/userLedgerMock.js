var DEFAULT_EXPENSE_CATEGORIES = [
  { name: "餐饮", badge: "餐" },
  { name: "交通", badge: "交" },
  { name: "购物", badge: "购" },
  { name: "居家", badge: "居" },
  { name: "娱乐", badge: "娱" }
];

var DEFAULT_INCOME_CATEGORIES = [
  { name: "工资", badge: "工" },
  { name: "兼职", badge: "兼" },
  { name: "奖金", badge: "奖" },
  { name: "投资", badge: "投" }
];

var INITIAL_ENTRIES = [
  {
    id: "entry-1",
    type: "expense",
    categoryName: "餐饮",
    badge: "餐",
    amount: 258,
    note: "午餐聚会",
    date: "2026-04-09",
    time: "12:30",
    imageName: ""
  },
  {
    id: "entry-2",
    type: "expense",
    categoryName: "交通",
    badge: "交",
    amount: 128,
    note: "打车",
    date: "2026-04-09",
    time: "09:15",
    imageName: ""
  },
  {
    id: "entry-3",
    type: "expense",
    categoryName: "奶茶",
    badge: "茶",
    amount: 72,
    note: "下午茶",
    date: "2026-04-09",
    time: "15:20",
    imageName: ""
  },
  {
    id: "entry-4",
    type: "expense",
    categoryName: "购物",
    badge: "购",
    amount: 980,
    note: "日用品采购",
    date: "2026-04-08",
    time: "19:30",
    imageName: ""
  },
  {
    id: "entry-5",
    type: "expense",
    categoryName: "娱乐",
    badge: "娱",
    amount: 300,
    note: "看电影",
    date: "2026-04-08",
    time: "20:00",
    imageName: ""
  },
  {
    id: "entry-6",
    type: "income",
    categoryName: "工资",
    badge: "工",
    amount: 4200,
    note: "月度工资",
    date: "2026-03-28",
    time: "10:30",
    imageName: ""
  },
  {
    id: "entry-7",
    type: "income",
    categoryName: "兼职",
    badge: "兼",
    amount: 300,
    note: "周末项目收入",
    date: "2026-03-28",
    time: "21:00",
    imageName: ""
  },
  {
    id: "entry-8",
    type: "expense",
    categoryName: "居家",
    badge: "居",
    amount: 260,
    note: "家居用品补货",
    date: "2026-03-26",
    time: "11:20",
    imageName: ""
  },
  {
    id: "entry-9",
    type: "expense",
    categoryName: "蔬菜",
    badge: "蔬",
    amount: 120,
    note: "买菜",
    date: "2026-03-26",
    time: "18:10",
    imageName: ""
  },
  {
    id: "entry-10",
    type: "expense",
    categoryName: "社交",
    badge: "社",
    amount: 180,
    note: "朋友聚餐",
    date: "2026-03-26",
    time: "20:30",
    imageName: ""
  }
];

var storeByUser = {};

function clone(value) {
  return JSON.parse(JSON.stringify(value));
}

function pad(value) {
  return String(value).padStart(2, "0");
}

function getTodayDate() {
  var now = new Date();
  return now.getFullYear() + "-" + pad(now.getMonth() + 1) + "-" + pad(now.getDate());
}

function getCurrentTime() {
  var now = new Date();
  return pad(now.getHours()) + ":" + pad(now.getMinutes());
}

function monthLabel(monthKey) {
  var parts = monthKey.split("-");
  return parts[0] + "年" + Number(parts[1]) + "月";
}

function parseDateTime(entry) {
  var withTime = entry.time || "00:00";
  return new Date(entry.date + "T" + withTime + ":00");
}

function weekdayLabel(dateText) {
  var weekMap = ["星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"];
  return weekMap[new Date(dateText + "T00:00:00").getDay()];
}

function formatMoney(value) {
  return "¥" + Number(value || 0).toLocaleString("zh-CN");
}

function resolveUserKey() {
  if (typeof window === "undefined") {
    return "local-user";
  }

  try {
    var profileRaw = window.localStorage.getItem("bill_user_profile");

    if (!profileRaw) {
      return "local-user";
    }

    var profile = JSON.parse(profileRaw);
    return profile.account || profile.username || profile.phone || "local-user";
  } catch (error) {
    return "local-user";
  }
}

function buildDefaultCategories() {
  var nextId = 1;
  var categories = {
    expense: [],
    income: []
  };

  DEFAULT_EXPENSE_CATEGORIES.forEach(function(item) {
    categories.expense.push({
      id: "expense-default-" + nextId,
      type: "expense",
      name: item.name,
      badge: item.badge,
      isDefault: true
    });
    nextId += 1;
  });

  DEFAULT_INCOME_CATEGORIES.forEach(function(item) {
    categories.income.push({
      id: "income-default-" + nextId,
      type: "income",
      name: item.name,
      badge: item.badge,
      isDefault: true
    });
    nextId += 1;
  });

  return {
    categories: categories,
    nextCategoryIndex: nextId
  };
}

function createInitialState() {
  var categoryBundle = buildDefaultCategories();
  var entries = clone(INITIAL_ENTRIES);
  var maxEntryNumber = entries.reduce(function(maxValue, item) {
    var numberValue = Number(String(item.id).replace(/[^\d]/g, ""));
    return numberValue > maxValue ? numberValue : maxValue;
  }, 10);

  return {
    categories: categoryBundle.categories,
    nextCategoryIndex: categoryBundle.nextCategoryIndex,
    nextEntryIndex: maxEntryNumber + 1,
    entries: entries
  };
}

function ensureState() {
  var userKey = resolveUserKey();

  if (!storeByUser[userKey]) {
    storeByUser[userKey] = createInitialState();
  }

  return storeByUser[userKey];
}

function getMonthKey(dateText) {
  return dateText.slice(0, 7);
}

function buildMonthGroups(entries) {
  var grouped = {};

  entries.forEach(function(item) {
    if (!grouped[item.date]) {
      grouped[item.date] = {
        date: item.date,
        weekday: weekdayLabel(item.date),
        totalIncome: 0,
        totalExpense: 0,
        items: []
      };
    }

    if (item.type === "income") {
      grouped[item.date].totalIncome += Number(item.amount || 0);
    } else {
      grouped[item.date].totalExpense += Number(item.amount || 0);
    }

    grouped[item.date].items.push({
      id: item.id,
      badge: item.badge || item.categoryName.slice(0, 1),
      category: item.categoryName,
      time: item.time,
      note: item.note,
      amount: item.amount,
      type: item.type,
      imageName: item.imageName || ""
    });
  });

  return Object.keys(grouped)
    .sort(function(a, b) {
      return new Date(b + "T00:00:00") - new Date(a + "T00:00:00");
    })
    .map(function(dateKey) {
      var group = grouped[dateKey];
      group.items.sort(function(a, b) {
        return b.time.localeCompare(a.time);
      });
      return group;
    });
}

function buildMonthCategories(monthEntries, categories) {
  var statsMap = {};

  monthEntries.forEach(function(item) {
    var name = item.categoryName;

    if (!statsMap[name]) {
      statsMap[name] = {
        name: name,
        badge: item.badge || name.slice(0, 1),
        count: 0
      };
    }

    statsMap[name].count += 1;
  });

  ["expense", "income"].forEach(function(type) {
    categories[type].forEach(function(item) {
      if (!statsMap[item.name]) {
        statsMap[item.name] = {
          name: item.name,
          badge: item.badge,
          count: 0
        };
      }
    });
  });

  return Object.keys(statsMap)
    .map(function(key) {
      return statsMap[key];
    })
    .sort(function(a, b) {
      if (b.count !== a.count) {
        return b.count - a.count;
      }
      return a.name.localeCompare(b.name, "zh-CN");
    })
    .slice(0, 6);
}

function buildMonthOverview(monthEntries, summary, categories) {
  var days = {};
  var maxCategory = categories[0] || { name: "暂无", count: 0 };
  var expenseDays = 0;

  monthEntries.forEach(function(item) {
    days[item.date] = true;
    if (item.type === "expense") {
      expenseDays += 1;
    }
  });

  var recordedDays = Object.keys(days).length || 1;
  var averageExpense = summary.expense / recordedDays;

  return [
    {
      label: "记账天数",
      value: recordedDays + "天",
      progress: Math.min(100, recordedDays * 12)
    },
    {
      label: "日均支出",
      value: formatMoney(Math.round(averageExpense)),
      progress: Math.min(100, expenseDays * 8)
    },
    {
      label: "高频分类",
      value: maxCategory.name + (maxCategory.count > 0 ? " · " + maxCategory.count + "笔" : ""),
      progress: Math.min(100, maxCategory.count * 14 || 20)
    }
  ];
}

function buildMonthData(monthKey, state) {
  var monthEntries = state.entries
    .filter(function(item) {
      return getMonthKey(item.date) === monthKey;
    })
    .slice()
    .sort(function(a, b) {
      return parseDateTime(b) - parseDateTime(a);
    });

  var summary = monthEntries.reduce(function(result, item) {
    if (item.type === "income") {
      result.income += Number(item.amount || 0);
    } else {
      result.expense += Number(item.amount || 0);
    }

    return result;
  }, { income: 0, expense: 0 });

  summary.balance = summary.income - summary.expense;

  var categories = buildMonthCategories(monthEntries, state.categories);

  return {
    key: monthKey,
    label: monthLabel(monthKey),
    summary: summary,
    groups: buildMonthGroups(monthEntries),
    categories: categories,
    overview: buildMonthOverview(monthEntries, summary, categories)
  };
}

function getLedgerMonths() {
  var state = ensureState();
  var keyMap = {};

  state.entries.forEach(function(item) {
    keyMap[getMonthKey(item.date)] = true;
  });

  if (Object.keys(keyMap).length === 0) {
    keyMap[getMonthKey(getTodayDate())] = true;
  }

  return Object.keys(keyMap)
    .sort(function(a, b) {
      return b.localeCompare(a);
    })
    .map(function(monthKey) {
      return buildMonthData(monthKey, state);
    });
}

function getLedgerCategoryGroups() {
  var state = ensureState();
  return clone(state.categories);
}

function getLedgerCategoryOptions(type) {
  var state = ensureState();
  return clone(state.categories[type] || []);
}

function createLedgerCategory(payload) {
  var state = ensureState();
  var type = payload.type === "income" ? "income" : "expense";
  var name = String(payload.name || "").trim();
  var badge = String(payload.badge || "").trim();

  if (!name) {
    throw new Error("分类名称不能为空");
  }

  var exists = state.categories[type].some(function(item) {
    return item.name === name;
  });

  if (exists) {
    throw new Error("分类名称已存在");
  }

  var category = {
    id: type + "-custom-" + state.nextCategoryIndex,
    type: type,
    name: name,
    badge: (badge || name.slice(0, 1)).slice(0, 2),
    isDefault: false
  };

  state.nextCategoryIndex += 1;
  state.categories[type].push(category);

  return clone(category);
}

function removeLedgerCategory(payload) {
  var state = ensureState();
  var type = payload.type === "income" ? "income" : "expense";
  var id = payload.id;
  var index = state.categories[type].findIndex(function(item) {
    return item.id === id;
  });

  if (index < 0) {
    return false;
  }

  if (state.categories[type][index].isDefault) {
    throw new Error("默认分类不允许删除");
  }

  state.categories[type].splice(index, 1);
  return true;
}

function addLedgerEntry(payload) {
  var state = ensureState();
  var type = payload.type === "income" ? "income" : "expense";
  var amount = Number(payload.amount || 0);
  var note = String(payload.note || "").trim();
  var date = payload.date || getTodayDate();
  var imageName = String(payload.imageName || "").trim();
  var category = state.categories[type].find(function(item) {
    return item.id === payload.categoryId;
  });

  if (!category) {
    throw new Error("请选择分类");
  }

  if (!(amount > 0)) {
    throw new Error("金额必须大于 0");
  }

  var entry = {
    id: "entry-" + state.nextEntryIndex,
    type: type,
    categoryName: category.name,
    badge: category.badge,
    amount: Math.round(amount * 100) / 100,
    note: note || (type === "income" ? "新增收入" : "新增支出"),
    date: date,
    time: getCurrentTime(),
    imageName: imageName
  };

  state.nextEntryIndex += 1;
  state.entries.push(entry);

  return clone(entry);
}

function formatLedgerCurrency(value) {
  return formatMoney(value);
}

export {
  addLedgerEntry,
  createLedgerCategory,
  formatLedgerCurrency,
  getLedgerCategoryGroups,
  getLedgerCategoryOptions,
  getLedgerMonths,
  getTodayDate,
  removeLedgerCategory
};
