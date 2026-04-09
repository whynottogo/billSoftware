function clone(value) {
  return JSON.parse(JSON.stringify(value));
}

function toCurrency(number) {
  return Number(number || 0);
}

function formatCurrency(amount) {
  var number = Number(amount || 0);
  var absNumber = Math.abs(number);
  var formatted = absNumber.toLocaleString("zh-CN", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  });

  return (number < 0 ? "-¥" : "¥") + formatted;
}

function createRecord(id, monthKey, dateLabel, action, change, balanceAfter, note, source) {
  return {
    id: id,
    monthKey: monthKey,
    dateLabel: dateLabel,
    action: action,
    change: change,
    balanceAfter: balanceAfter,
    note: note,
    source: source
  };
}

const accountTypes = [
  { label: "现金", value: "cash" },
  { label: "银行卡", value: "bankCard" },
  { label: "信用卡", value: "creditCard" },
  { label: "虚拟账户", value: "virtual" },
  { label: "投资账户", value: "investment" },
  { label: "负债账户", value: "liability" }
];

const virtualProviders = [
  { label: "微信", value: "wechat" },
  { label: "支付宝", value: "alipay" }
];

const assetSeed = [
  {
    id: "category-liquid",
    name: "流动资产",
    tone: "info",
    accounts: [
      {
        id: "acc-cmb-6225",
        name: "招商银行卡",
        type: "bankCard",
        typeLabel: "银行卡",
        remark: "日常收支主账户",
        balance: 36850.3,
        cardNo: "6225 **** **** 1902",
        provider: "",
        direction: "asset",
        monthlyChange: 1820.4
      },
      {
        id: "acc-cash-wallet",
        name: "现金钱包",
        type: "cash",
        typeLabel: "现金",
        remark: "零散备用金",
        balance: 1260,
        cardNo: "",
        provider: "",
        direction: "asset",
        monthlyChange: -120
      },
      {
        id: "acc-wechat-main",
        name: "微信零钱",
        type: "virtual",
        typeLabel: "虚拟账户",
        remark: "日常支付",
        balance: 2580.9,
        cardNo: "",
        provider: "wechat",
        direction: "asset",
        monthlyChange: 318.4
      },
      {
        id: "acc-alipay-main",
        name: "支付宝余额",
        type: "virtual",
        typeLabel: "虚拟账户",
        remark: "线上消费",
        balance: 1940.2,
        cardNo: "",
        provider: "alipay",
        direction: "asset",
        monthlyChange: -65
      }
    ]
  },
  {
    id: "category-investment",
    name: "投资资产",
    tone: "success",
    accounts: [
      {
        id: "acc-fund-growth",
        name: "基金定投账户",
        type: "investment",
        typeLabel: "投资账户",
        remark: "长期定投",
        balance: 84200,
        cardNo: "",
        provider: "",
        direction: "asset",
        monthlyChange: 2780
      },
      {
        id: "acc-bond-safe",
        name: "稳健债券账户",
        type: "investment",
        typeLabel: "投资账户",
        remark: "低波动仓位",
        balance: 36500,
        cardNo: "",
        provider: "",
        direction: "asset",
        monthlyChange: 420
      }
    ]
  },
  {
    id: "category-liability",
    name: "负债账户",
    tone: "danger",
    accounts: [
      {
        id: "acc-credit-9812",
        name: "招商信用卡",
        type: "creditCard",
        typeLabel: "信用卡",
        remark: "每月 8 日还款",
        balance: 6800,
        cardNo: "9812 **** **** 8846",
        provider: "",
        direction: "liability",
        monthlyChange: -560
      },
      {
        id: "acc-huabei-virtual",
        name: "花呗",
        type: "liability",
        typeLabel: "负债账户",
        remark: "次月 10 日还款",
        balance: 2380,
        cardNo: "",
        provider: "",
        direction: "liability",
        monthlyChange: -210
      }
    ]
  }
];

const recordSeed = {
  "acc-cmb-6225": [
    createRecord("cmb-1", "2026-04", "04-09 20:35", "增加", 500, 36850.3, "工资补发", "手动操作"),
    createRecord("cmb-2", "2026-04", "04-07 09:10", "减少", -420, 36350.3, "房租转出", "手动操作"),
    createRecord("cmb-3", "2026-03", "03-25 18:44", "调整", 180, 36770.3, "对账后补记", "账户设置"),
    createRecord("cmb-4", "2026-03", "03-12 11:18", "增加", 2000, 36590.3, "项目回款", "手动操作"),
    createRecord("cmb-5", "2026-02", "02-08 14:30", "减少", -860, 34590.3, "信用卡还款", "手动操作")
  ],
  "acc-cash-wallet": [
    createRecord("cash-1", "2026-04", "04-08 16:12", "减少", -60, 1260, "午餐现金支出", "手动操作"),
    createRecord("cash-2", "2026-03", "03-23 08:18", "增加", 200, 1320, "零钱补充", "手动操作")
  ],
  "acc-wechat-main": [
    createRecord("wechat-1", "2026-04", "04-09 12:11", "减少", -88, 2580.9, "咖啡订阅", "手动操作"),
    createRecord("wechat-2", "2026-03", "03-31 21:04", "增加", 300, 2668.9, "朋友转账", "手动操作")
  ],
  "acc-alipay-main": [
    createRecord("alipay-1", "2026-04", "04-06 09:25", "减少", -68, 1940.2, "地铁月卡", "手动操作"),
    createRecord("alipay-2", "2026-03", "03-14 14:20", "增加", 120, 2008.2, "红包入账", "手动操作")
  ],
  "acc-fund-growth": [
    createRecord("fund-1", "2026-04", "04-05 10:00", "增加", 1200, 84200, "月度定投", "手动操作"),
    createRecord("fund-2", "2026-03", "03-05 10:00", "增加", 1200, 83000, "月度定投", "手动操作")
  ],
  "acc-bond-safe": [
    createRecord("bond-1", "2026-04", "04-01 09:10", "增加", 500, 36500, "债券收益结转", "手动操作"),
    createRecord("bond-2", "2026-03", "03-01 09:10", "增加", 480, 36000, "债券收益结转", "手动操作")
  ],
  "acc-credit-9812": [
    createRecord("credit-1", "2026-04", "04-03 20:42", "减少", -580, 6800, "商超消费", "手动操作"),
    createRecord("credit-2", "2026-03", "03-09 09:12", "调整", 120, 7380, "账单校正", "账户设置")
  ],
  "acc-huabei-virtual": [
    createRecord("huabei-1", "2026-04", "04-04 18:20", "减少", -220, 2380, "外卖消费", "手动操作"),
    createRecord("huabei-2", "2026-03", "03-08 10:40", "减少", -180, 2600, "交通充值", "手动操作")
  ]
};

function computeSummary(categories) {
  var totalAsset = 0;
  var totalLiability = 0;
  var monthlyChange = 0;
  var accountCount = 0;

  categories.forEach(function(category) {
    category.accounts.forEach(function(account) {
      accountCount += 1;
      monthlyChange += toCurrency(account.monthlyChange);

      if (account.direction === "liability") {
        totalLiability += toCurrency(account.balance);
      } else {
        totalAsset += toCurrency(account.balance);
      }
    });
  });

  return {
    totalAsset: totalAsset,
    totalLiability: totalLiability,
    netAsset: totalAsset - totalLiability,
    monthlyChange: monthlyChange,
    accountCount: accountCount,
    updatedAt: "2026-04-09 21:00"
  };
}

function enrichCategories(categories) {
  return categories.map(function(category) {
    var total = category.accounts.reduce(function(sum, account) {
      return sum + toCurrency(account.balance);
    }, 0);

    return Object.assign({}, category, {
      total: total
    });
  });
}

function buildOverview() {
  var categories = enrichCategories(clone(assetSeed));

  return {
    summary: computeSummary(categories),
    categories: categories
  };
}

function getAssetOverview() {
  return buildOverview();
}

function findAccountById(accountId) {
  var overview = buildOverview();
  var targetAccount = null;
  var targetCategory = null;

  overview.categories.forEach(function(category) {
    category.accounts.forEach(function(account) {
      if (account.id === accountId) {
        targetAccount = account;
        targetCategory = {
          id: category.id,
          name: category.name,
          tone: category.tone
        };
      }
    });
  });

  if (!targetAccount && overview.categories.length > 0 && overview.categories[0].accounts.length > 0) {
    targetAccount = overview.categories[0].accounts[0];
    targetCategory = {
      id: overview.categories[0].id,
      name: overview.categories[0].name,
      tone: overview.categories[0].tone
    };
  }

  return {
    account: targetAccount,
    category: targetCategory
  };
}

function getAssetAccountDetail(accountId) {
  var accountMatch = findAccountById(accountId);
  var account = accountMatch.account;
  var category = accountMatch.category;
  var records = clone(recordSeed[account.id] || []);

  return {
    account: account,
    category: category,
    records: records
  };
}

function getAssetTypeOptions() {
  return clone(accountTypes);
}

function getVirtualProviderOptions() {
  return clone(virtualProviders);
}

function monthKeysFromRecords(records) {
  var exists = {};
  var keys = [];

  records.forEach(function(record) {
    if (!exists[record.monthKey]) {
      exists[record.monthKey] = true;
      keys.push(record.monthKey);
    }
  });

  return keys.sort(function(a, b) {
    return a > b ? -1 : 1;
  });
}

function getMonthLabel(monthKey) {
  if (!monthKey) {
    return "全部月份";
  }

  var parts = String(monthKey).split("-");
  return parts[0] + " 年 " + parts[1] + " 月";
}

export {
  formatCurrency,
  getAssetOverview,
  getAssetAccountDetail,
  getAssetTypeOptions,
  getVirtualProviderOptions,
  monthKeysFromRecords,
  getMonthLabel
};
