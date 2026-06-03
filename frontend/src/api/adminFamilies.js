import request from "@/utils/request";

function firstDefined(values, fallback) {
  for (var index = 0; index < values.length; index += 1) {
    if (values[index] !== undefined && values[index] !== null && values[index] !== "") {
      return values[index];
    }
  }
  return fallback;
}

function toNumber(value, fallback) {
  var parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : fallback;
}

function unwrap(result) {
  return result && result.data ? result.data : result || {};
}

function formatAdminFamilyCurrency(value) {
  return "¥" + Number(value || 0).toLocaleString("zh-CN", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  });
}

function buildAdminFamilyError(error, fallback) {
  if (error && error.response && error.response.data && error.response.data.message) {
    return error.response.data.message;
  }
  return fallback;
}

function normalizeSummary(result) {
  var payload = unwrap(result);
  return {
    totalFamilies: toNumber(firstDefined([payload.totalFamilies, payload.total_families], 0), 0),
    totalMembers: toNumber(firstDefined([payload.totalMembers, payload.total_members], 0), 0),
    averageMembers: toNumber(firstDefined([payload.averageMembers, payload.average_members], 0), 0),
    activeFamilies: toNumber(firstDefined([payload.activeFamilies, payload.active_families], 0), 0),
    activeDays: toNumber(firstDefined([payload.activeDays, payload.active_days], 30), 30)
  };
}

function normalizeFamilyRow(row) {
  var source = row || {};
  var id = String(firstDefined([source.id, source.familyId, source.family_id, source.family_uid], ""));
  return {
    id: id,
    name: String(firstDefined([source.name, source.familyName, source.family_name], "未命名家庭")),
    creator: String(firstDefined([source.creator, source.creatorName, source.creator_name], "-")),
    memberCount: toNumber(firstDefined([source.memberCount, source.members, source.member_count], 0), 0),
    billCount: toNumber(firstDefined([source.billCount, source.bill_count], 0), 0),
    totalAssets: toNumber(firstDefined([source.totalAssets, source.total_assets], 0), 0),
    createdAt: String(firstDefined([source.createdAt, source.created_at], "-")),
    status: String(firstDefined([source.status], "待记账"))
  };
}

function normalizePagination(source) {
  var payload = source || {};
  return {
    page: toNumber(firstDefined([payload.page, payload.currentPage, payload.current_page], 1), 1),
    pageSize: toNumber(firstDefined([payload.pageSize, payload.page_size], 10), 10),
    total: toNumber(firstDefined([payload.total], 0), 0),
    totalPages: toNumber(firstDefined([payload.totalPages, payload.total_pages], 0), 0)
  };
}

function normalizeList(result) {
  var payload = unwrap(result);
  var rows = firstDefined([payload.list, payload.families, payload.rows, payload.items], []);
  return {
    families: Array.isArray(rows) ? rows.map(normalizeFamilyRow) : [],
    pagination: normalizePagination(payload.pagination)
  };
}

function normalizeMember(row) {
  var source = row || {};
  return {
    userId: String(firstDefined([source.userId, source.user_id], "")),
    name: String(firstDefined([source.name, source.nickname, source.username], "-")),
    role: String(firstDefined([source.role], "成员")),
    joinedAt: String(firstDefined([source.joinedAt, source.joined_at], "-")),
    color: String(firstDefined([source.color], "#f6d34a"))
  };
}

function normalizeRecentBill(row) {
  var source = row || {};
  return {
    id: String(firstDefined([source.id], "")),
    memberName: String(firstDefined([source.memberName, source.member_name], "-")),
    recordType: String(firstDefined([source.recordType, source.record_type], "")),
    amount: toNumber(firstDefined([source.amount], 0), 0),
    remark: String(firstDefined([source.remark], "")),
    recordDate: String(firstDefined([source.recordDate, source.record_date], "-"))
  };
}

function normalizeDetail(result, fallbackFamilyId) {
  var payload = unwrap(result);
  var family = normalizeFamilyRow(
    Object.assign({}, payload.family || {}, {
      id: firstDefined([payload.family && payload.family.id, payload.family && payload.family.family_uid], fallbackFamilyId)
    })
  );
  var incomeExpense = payload.income_expense || payload.incomeExpense || {};
  var assetSummary = payload.asset_summary || payload.assetSummary || {};
  var members = Array.isArray(payload.members) ? payload.members.map(normalizeMember) : [];
  var recentBills = Array.isArray(payload.recent_bills)
    ? payload.recent_bills.map(normalizeRecentBill)
    : Array.isArray(payload.recentBills)
      ? payload.recentBills.map(normalizeRecentBill)
      : [];

  return {
    family: Object.assign({}, family, {
      memberCount: members.length || family.memberCount
    }),
    members: members,
    incomeExpense: {
      income: toNumber(firstDefined([incomeExpense.income], 0), 0),
      expense: toNumber(firstDefined([incomeExpense.expense], 0), 0),
      balance: toNumber(firstDefined([incomeExpense.balance], 0), 0),
      records: toNumber(firstDefined([incomeExpense.records], 0), 0)
    },
    assetSummary: {
      totalAssets: toNumber(firstDefined([assetSummary.totalAssets, assetSummary.total_assets], 0), 0),
      totalLiabilities: toNumber(firstDefined([assetSummary.totalLiabilities, assetSummary.total_liabilities], 0), 0),
      netAssets: toNumber(firstDefined([assetSummary.netAssets, assetSummary.net_assets], 0), 0),
      accountCount: toNumber(firstDefined([assetSummary.accountCount, assetSummary.account_count], 0), 0)
    },
    recentBills: recentBills
  };
}

export function getAdminFamilySummary() {
  return request.get("/admin/families/summary");
}

export function listAdminFamilies(params) {
  return request.get("/admin/families", {
    params: params || {}
  });
}

export function getAdminFamilyDetail(familyId) {
  return request.get("/admin/families/" + encodeURIComponent(familyId));
}

export {
  buildAdminFamilyError,
  formatAdminFamilyCurrency,
  normalizeDetail,
  normalizeList,
  normalizeSummary
};
