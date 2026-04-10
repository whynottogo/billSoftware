import request from "@/utils/request";

const assetTypeOptions = [
  { label: "现金", value: "cash" },
  { label: "银行卡", value: "bankCard" },
  { label: "信用卡", value: "creditCard" },
  { label: "虚拟账户", value: "virtual" },
  { label: "投资账户", value: "investment" },
  { label: "负债账户", value: "liability" }
];

const virtualProviderOptions = [
  { label: "微信", value: "wechat" },
  { label: "支付宝", value: "alipay" }
];

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

function toStringId(value, fallback) {
  var target = value;

  if (target === undefined || target === null || target === "") {
    target = fallback;
  }

  if (target === undefined || target === null || target === "") {
    return "";
  }

  return String(target);
}

function toPayloadId(value) {
  var normalized = toStringId(value, "");

  if (!normalized) {
    return normalized;
  }

  if (/^\d+$/.test(normalized)) {
    return Number(normalized);
  }

  return normalized;
}

function cloneOptions(list) {
  return list.map(function(item) {
    return Object.assign({}, item);
  });
}

function formatAssetCurrency(value) {
  var number = toNumber(value);
  var formatted = Math.abs(number).toLocaleString("zh-CN", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  });

  return (number < 0 ? "-¥" : "¥") + formatted;
}

function normalizeDateLabel(value) {
  if (!value) {
    return "";
  }

  var text = String(value).trim();
  var match = text.match(/(\d{4})[-/](\d{2})[-/](\d{2})[ T](\d{2}):(\d{2})/);

  if (match) {
    return match[2] + "-" + match[3] + " " + match[4] + ":" + match[5];
  }

  var shortMatch = text.match(/(\d{2})-(\d{2}) (\d{2}):(\d{2})/);
  if (shortMatch) {
    return shortMatch[1] + "-" + shortMatch[2] + " " + shortMatch[3] + ":" + shortMatch[4];
  }

  return text.slice(0, 16);
}

function getMonthKeyFromDate(value) {
  if (!value) {
    return "";
  }

  var text = String(value).trim();
  var match = text.match(/(\d{4})[-/](\d{2})/);

  if (!match) {
    return "";
  }

  return match[1] + "-" + match[2];
}

function getMonthLabel(monthKey) {
  if (!monthKey) {
    return "全部月份";
  }

  var parts = String(monthKey).split("-");
  if (parts.length !== 2) {
    return "全部月份";
  }

  return parts[0] + " 年 " + parts[1] + " 月";
}

function monthKeysFromRecords(records) {
  var exists = {};
  var keys = [];

  (Array.isArray(records) ? records : []).forEach(function(record) {
    if (record.monthKey && !exists[record.monthKey]) {
      exists[record.monthKey] = true;
      keys.push(record.monthKey);
    }
  });

  return keys.sort(function(left, right) {
    return left > right ? -1 : 1;
  });
}

function normalizeType(value) {
  var target = String(value || "").trim();
  var aliases = {
    cash: "cash",
    bank: "bankCard",
    bankcard: "bankCard",
    bank_card: "bankCard",
    credit: "creditCard",
    creditcard: "creditCard",
    credit_card: "creditCard",
    virtual: "virtual",
    ewallet: "virtual",
    wallet: "virtual",
    investment: "investment",
    invest: "investment",
    liability: "liability",
    debt: "liability"
  };

  return aliases[target.toLowerCase()] || target || "cash";
}

function inferDirection(type, source) {
  var explicitDirection = String(
    source.direction || source.account_direction || source.balance_direction || ""
  ).toLowerCase();

  if (explicitDirection === "liability" || explicitDirection === "debt") {
    return "liability";
  }

  if (explicitDirection === "asset") {
    return "asset";
  }

  var normalizedType = normalizeType(type);
  if (normalizedType === "creditCard" || normalizedType === "liability") {
    return "liability";
  }

  return "asset";
}

function getTypeLabel(type) {
  var match = assetTypeOptions.find(function(item) {
    return item.value === normalizeType(type);
  });

  return match ? match.label : "账户";
}

function normalizeCategoryTone(category, accounts) {
  var source = category || {};
  var tone = source.tone || source.colorTone || "";

  if (tone) {
    return tone;
  }

  var name = String(source.name || source.category_name || "").toLowerCase();
  var hasLiability = (Array.isArray(accounts) ? accounts : []).some(function(item) {
    return item.direction === "liability";
  });

  if (name.indexOf("负债") >= 0 || hasLiability) {
    return "danger";
  }

  if (name.indexOf("投资") >= 0) {
    return "success";
  }

  return "info";
}

function normalizeAssetAccount(account, category) {
  var source = account || {};
  var type = normalizeType(source.type || source.account_type || source.accountType);
  var direction = inferDirection(type, source);
  var categoryId = toStringId(
    source.categoryId || source.category_id,
    category && (category.id || category.category_id)
  );
  var categoryName = source.categoryName || source.category_name || (category && category.name) || "未分类";

  return {
    id: toStringId(source.id || source.account_id),
    name: source.name || source.account_name || source.title || "未命名账户",
    type: type,
    typeLabel: source.typeLabel || source.type_label || getTypeLabel(type),
    remark: source.remark || source.note || source.description || "",
    balance: toNumber(
      source.balance || source.current_balance || source.amount || source.total_balance
    ),
    cardNo: source.cardNo || source.card_no || source.card || "",
    provider: source.provider || source.provider_code || source.channel || "",
    direction: direction,
    monthlyChange: toNumber(
      source.monthlyChange || source.monthly_change || source.monthChange || source.delta
    ),
    categoryId: categoryId,
    categoryName: categoryName
  };
}

function normalizeAssetCategory(category) {
  var source = category || {};
  var accountsSource = source.accounts || source.items || source.list || [];
  var accounts = (Array.isArray(accountsSource) ? accountsSource : []).map(function(item) {
    return normalizeAssetAccount(item, source);
  });
  var total = toNumber(source.total || source.total_amount || source.balance_total);

  if (!total) {
    total = accounts.reduce(function(sum, item) {
      return sum + toNumber(item.balance);
    }, 0);
  }

  return {
    id: toStringId(source.id || source.category_id || source.categoryId || source.code || source.name),
    name: source.name || source.category_name || "未分类",
    tone: normalizeCategoryTone(source, accounts),
    total: total,
    accounts: accounts
  };
}

function normalizeSummary(summary, categories) {
  var source = summary || {};
  var categoryList = Array.isArray(categories) ? categories : [];
  var assetTotal = categoryList.reduce(function(sum, category) {
    return sum + category.accounts.reduce(function(accountSum, account) {
      if (account.direction === "liability") {
        return accountSum;
      }

      return accountSum + toNumber(account.balance);
    }, 0);
  }, 0);
  var liabilityTotal = categoryList.reduce(function(sum, category) {
    return sum + category.accounts.reduce(function(accountSum, account) {
      if (account.direction !== "liability") {
        return accountSum;
      }

      return accountSum + toNumber(account.balance);
    }, 0);
  }, 0);
  var accountCount = categoryList.reduce(function(sum, category) {
    return sum + category.accounts.length;
  }, 0);
  var monthlyChange = categoryList.reduce(function(sum, category) {
    return sum + category.accounts.reduce(function(accountSum, account) {
      return accountSum + toNumber(account.monthlyChange);
    }, 0);
  }, 0);

  return {
    netAsset: toNumber(source.netAsset || source.net_asset || source.netWorth || assetTotal - liabilityTotal),
    totalAsset: toNumber(source.totalAsset || source.total_asset || source.assetTotal || assetTotal),
    totalLiability: toNumber(
      source.totalLiability || source.total_liability || source.liabilityTotal || liabilityTotal
    ),
    monthlyChange: toNumber(
      source.monthlyChange || source.monthly_change || source.monthChange || monthlyChange
    ),
    accountCount: toNumber(source.accountCount || source.account_count || accountCount),
    updatedAt: source.updatedAt || source.updated_at || source.last_updated_at || source.sync_at || ""
  };
}

function normalizeRecordAction(source) {
  var raw = String(source.action || source.action_label || source.action_type || source.type || "").toLowerCase();

  if (raw === "increase" || raw === "credit" || raw === "deposit" || raw === "增加") {
    return "增加";
  }

  if (raw === "decrease" || raw === "debit" || raw === "withdraw" || raw === "减少") {
    return "减少";
  }

  return "调整";
}

function normalizeRecordChange(source, account) {
  if (source.change !== undefined || source.delta !== undefined || source.change_amount !== undefined) {
    return toNumber(source.change || source.delta || source.change_amount);
  }

  var amount = toNumber(source.amount || source.operation_amount || source.value);
  var rawType = String(source.action_type || source.type || source.action || "").toLowerCase();

  if (rawType === "decrease" || rawType === "debit" || rawType === "withdraw" || rawType === "减少") {
    return -Math.abs(amount);
  }

  if (rawType === "increase" || rawType === "credit" || rawType === "deposit" || rawType === "增加") {
    return Math.abs(amount);
  }

  if (rawType === "adjust" || rawType === "调整") {
    var balanceBefore = toNumber(source.balance_before || source.before_balance);
    var balanceAfter = toNumber(
      source.balanceAfter || source.balance_after || source.after_balance || account.balance
    );

    if (balanceBefore || balanceAfter) {
      return balanceAfter - balanceBefore;
    }
  }

  return amount;
}

function normalizeAssetRecord(record, account) {
  var source = record || {};
  var monthKey = source.monthKey || source.month_key || getMonthKeyFromDate(source.created_at || source.createdAt || source.occurred_at || source.occurredAt || source.operation_time);

  return {
    id: toStringId(source.id || source.record_id || source.operation_id),
    monthKey: monthKey,
    dateLabel: source.dateLabel || source.date_label || normalizeDateLabel(
      source.created_at || source.createdAt || source.occurred_at || source.occurredAt || source.operation_time
    ),
    action: source.action || source.action_label || normalizeRecordAction(source),
    change: normalizeRecordChange(source, account),
    balanceAfter: toNumber(
      source.balanceAfter || source.balance_after || source.after_balance || source.current_balance || account.balance
    ),
    note: source.note || source.remark || source.description || "",
    source: source.source || source.source_label || source.operator || "手动操作"
  };
}

function sortRecords(records) {
  return records.slice().sort(function(left, right) {
    if (left.monthKey === right.monthKey) {
      return left.dateLabel > right.dateLabel ? -1 : 1;
    }

    return left.monthKey > right.monthKey ? -1 : 1;
  });
}

function normalizeAssetOverviewPayload(result) {
  var payload = extractPayload(result);
  var categoriesSource = payload.categories || payload.list || [];
  var categories = (Array.isArray(categoriesSource) ? categoriesSource : []).map(normalizeAssetCategory);

  return {
    summary: normalizeSummary(payload.summary || payload.overview, categories),
    categories: categories
  };
}

function normalizeAssetDetailPayload(result) {
  var payload = extractPayload(result);
  var rawCategory = payload.category || payload.categoryInfo || {};
  var account = normalizeAssetAccount(payload.account || payload.item || payload.detail, rawCategory);
  var category = {
    id: account.categoryId || toStringId(rawCategory.id || rawCategory.category_id),
    name: account.categoryName || rawCategory.name || rawCategory.category_name || "未分类",
    tone: normalizeCategoryTone(rawCategory, [account])
  };
  var recordsSource = payload.records || payload.operations || payload.logs || payload.history || [];
  var records = sortRecords((Array.isArray(recordsSource) ? recordsSource : []).map(function(item) {
    return normalizeAssetRecord(item, account);
  }));

  return {
    account: account,
    category: category,
    records: records
  };
}

function buildAssetPayload(model) {
  var source = model || {};
  var payload = {
    type: normalizeType(source.type),
    account_type: normalizeType(source.type),
    name: String(source.name || "").trim(),
    account_name: String(source.name || "").trim(),
    remark: String(source.remark || "").trim(),
    note: String(source.remark || "").trim(),
    balance: toNumber(source.balance),
    category_id: toPayloadId(source.categoryId),
    categoryId: toPayloadId(source.categoryId),
    card_no: String(source.cardNo || "").trim(),
    cardNo: String(source.cardNo || "").trim(),
    provider: source.provider || ""
  };

  return payload;
}

function buildAssetUpdatePayload(account, overrides) {
  var merged = Object.assign({}, account || {}, overrides || {});

  return buildAssetPayload({
    type: merged.type,
    name: merged.name,
    remark: merged.remark,
    balance: merged.balance,
    categoryId: merged.categoryId,
    cardNo: merged.cardNo,
    provider: merged.provider
  });
}

function buildAssetOperationPayload(model) {
  var source = model || {};
  var actionType = String(source.actionType || source.action_type || source.type || "adjust");
  var amount = toNumber(source.amount);
  var payload = {
    action_type: actionType,
    actionType: actionType,
    type: actionType,
    amount: amount,
    note: String(source.note || "").trim(),
    remark: String(source.note || "").trim()
  };

  if (actionType === "adjust") {
    payload.target_balance = amount;
    payload.targetBalance = amount;
  }

  return payload;
}

function buildUserAssetsError(error, fallback) {
  if (error && error.response && error.response.data && error.response.data.message) {
    return error.response.data.message;
  }

  return fallback;
}

function getAssetTypeOptions() {
  return cloneOptions(assetTypeOptions);
}

function getVirtualProviderOptions() {
  return cloneOptions(virtualProviderOptions);
}

function getUserAssets() {
  return request.get("/user/assets");
}

function createUserAsset(payload) {
  return request.post("/user/assets", payload);
}

function updateUserAsset(id, payload) {
  return request.put("/user/assets/" + id, payload);
}

function getUserAssetDetail(id) {
  return request.get("/user/assets/" + id);
}

function createUserAssetOperation(id, payload) {
  return request.post("/user/assets/" + id + "/operations", payload);
}

export {
  buildAssetOperationPayload,
  buildAssetPayload,
  buildAssetUpdatePayload,
  buildUserAssetsError,
  createUserAsset,
  createUserAssetOperation,
  formatAssetCurrency,
  getAssetTypeOptions,
  getMonthLabel,
  getUserAssetDetail,
  getUserAssets,
  getVirtualProviderOptions,
  monthKeysFromRecords,
  normalizeAssetDetailPayload,
  normalizeAssetOverviewPayload,
  updateUserAsset
};
