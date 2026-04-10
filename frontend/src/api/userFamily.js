import request from "@/utils/request";

const FAMILY_COLORS = ["#f6d34a", "#6bcf7c", "#4d96ff", "#ff8b8b", "#9b8cff", "#34d399", "#f97316"];

function extractPayload(result) {
  if (result && Object.prototype.hasOwnProperty.call(result, "data")) {
    return result.data;
  }

  return result || {};
}

function firstDefined(values, fallback) {
  var list = Array.isArray(values) ? values : [];

  for (var index = 0; index < list.length; index += 1) {
    var current = list[index];

    if (current !== undefined && current !== null && current !== "") {
      return current;
    }
  }

  return fallback;
}

function toNumber(value) {
  var normalized = Number(value || 0);

  if (!Number.isFinite(normalized)) {
    return 0;
  }

  return normalized;
}

function toArray(value) {
  if (Array.isArray(value)) {
    return value;
  }

  return [];
}

function pad(value) {
  return String(value).padStart(2, "0");
}

function hashText(value) {
  var hash = 0;
  var text = String(value || "");

  for (var index = 0; index < text.length; index += 1) {
    hash = (hash << 5) - hash + text.charCodeAt(index);
    hash |= 0;
  }

  return Math.abs(hash);
}

function pickFamilyColor(seedText, index) {
  var seed = hashText(seedText || index);
  return FAMILY_COLORS[(seed + index) % FAMILY_COLORS.length];
}

function formatDateLabel(value) {
  var raw = String(value || "").trim();

  if (!raw) {
    return "--";
  }

  if (/^\d{4}-\d{2}-\d{2}/.test(raw)) {
    return raw.slice(0, 10);
  }

  var normalized = new Date(raw);

  if (Number.isNaN(normalized.getTime())) {
    return raw;
  }

  return normalized.getFullYear() + "-" + pad(normalized.getMonth() + 1) + "-" + pad(normalized.getDate());
}

function formatMonthLabel(monthKey) {
  var key = String(monthKey || "");
  var parts = key.split("-");

  if (parts.length !== 2) {
    return key || "暂无月份";
  }

  return parts[0] + "年" + Number(parts[1]) + "月";
}

function formatYearLabel(yearKey) {
  var key = String(yearKey || "").trim();

  if (!key) {
    return "暂无年份";
  }

  return key + "年";
}

function normalizeFamilyId(value) {
  return String(value || "").trim().toUpperCase();
}

function formatFamilyCurrency(value) {
  return "¥" + Number(value || 0).toLocaleString("zh-CN");
}

function formatFamilyPercent(value) {
  return Number(value || 0).toFixed(1) + "%";
}

function buildInviteLink(inviteLink, familyId, inviteCode) {
  var explicit = String(inviteLink || "").trim();

  if (explicit) {
    return explicit;
  }

  var code = String(inviteCode || familyId || "").trim();

  if (!code) {
    return "";
  }

  if (typeof window !== "undefined" && window.location && window.location.origin) {
    return window.location.origin + "/invite/" + code;
  }

  return "/invite/" + code;
}

function buildCreatorName(source, members) {
  var owner = firstDefined(
    [
      source.creator,
      source.creator_name,
      source.creatorName,
      source.owner_name,
      source.ownerName,
      source.owner && source.owner.name,
      source.owner && source.owner.nickname,
      source.created_by_name,
      source.createdByName
    ],
    ""
  );

  if (owner) {
    return owner;
  }

  var memberList = Array.isArray(members) ? members : [];
  var creatorMember = memberList.find(function(member) {
    return member.role === "创建人";
  });

  return creatorMember ? creatorMember.name : "未知创建人";
}

function normalizePeriodKey(value, fallbackType) {
  var raw = String(value || "").trim();

  if (!raw) {
    return "";
  }

  if (fallbackType === "year") {
    return raw.slice(0, 4);
  }

  return raw.slice(0, 7);
}

function buildPeriodLabel(periodType, key) {
  return periodType === "year" ? formatYearLabel(key) : formatMonthLabel(key);
}

function normalizePeriodOption(item, index, periodType) {
  var source = item || {};
  var key = normalizePeriodKey(
    firstDefined([source.key, source.value, source.period_key, source.periodKey, source.period, source.month, source.year], ""),
    periodType
  );
  var income = toNumber(
    firstDefined(
      [
        source.income,
        source.totalIncome,
        source.total_income,
        source.monthIncome,
        source.month_income,
        source.yearIncome,
        source.year_income
      ],
      0
    )
  );
  var expense = toNumber(
    firstDefined(
      [
        source.expense,
        source.totalExpense,
        source.total_expense,
        source.monthExpense,
        source.month_expense,
        source.yearExpense,
        source.year_expense
      ],
      0
    )
  );
  var balance = toNumber(firstDefined([source.balance, source.totalBalance, source.total_balance], income - expense));

  return {
    key: key || String(index + 1),
    label: firstDefined([source.label, source.name, source.title], buildPeriodLabel(periodType, key)),
    note: firstDefined([source.note, source.description, source.remark, source.summary], "暂无说明"),
    income: income,
    expense: expense,
    balance: balance
  };
}

function normalizePeriodOptions(source, periodType) {
  var normalized = [];

  if (Array.isArray(source)) {
    normalized = source.map(function(item, index) {
      return normalizePeriodOption(item, index, periodType);
    });
  } else if (source && typeof source === "object") {
    normalized = Object.keys(source).map(function(key, index) {
      var item = source[key];
      var normalizedItem = normalizePeriodOption(
        typeof item === "object" && item !== null ? Object.assign({ key: key }, item) : { key: key, value: item },
        index,
        periodType
      );

      if (!normalizedItem.key) {
        normalizedItem.key = key;
      }

      return normalizedItem;
    });
  }

  return normalized.filter(function(item) {
    return !!item.key;
  }).sort(function(left, right) {
    return String(right.key).localeCompare(String(left.key));
  });
}

function normalizeMember(member, index, familyId) {
  var source = member || {};
  var name = firstDefined(
    [
      source.name,
      source.nickname,
      source.member_name,
      source.memberName,
      source.user_name,
      source.userName,
      source.username,
      source.display_name,
      source.displayName
    ],
    "成员" + (index + 1)
  );
  var role = firstDefined(
    [
      source.role,
      source.member_role,
      source.memberRole,
      source.is_creator ? "创建人" : "",
      source.is_owner ? "创建人" : "",
      source.isOwner ? "创建人" : ""
    ],
    "成员"
  );

  return {
    userId: firstDefined([source.user_id, source.userId, source.id], ""),
    name: name,
    role: role || "成员",
    color: firstDefined([source.color, source.avatar_color, source.avatarColor], pickFamilyColor(name + familyId, index))
  };
}

function normalizeMembers(source, familyId) {
  return toArray(source).map(function(member, index) {
    return normalizeMember(member, index, familyId);
  });
}

function normalizeFamilySummary(item) {
  var source = item || {};
  var familyId = normalizeFamilyId(
    firstDefined([source.id, source.family_id, source.familyId, source.family_code, source.familyCode, source.code], "")
  );
  var monthOptions = normalizePeriodOptions(
    firstDefined([source.monthOptions, source.month_options, source.months, source.month_summary, source.monthSummary], []),
    "month"
  );
  var yearOptions = normalizePeriodOptions(
    firstDefined([source.yearOptions, source.year_options, source.years, source.year_summary, source.yearSummary], []),
    "year"
  );
  var members = normalizeMembers(
    firstDefined([source.members, source.member_list, source.memberList, source.users, source.user_list, source.userList], []),
    familyId
  );
  var latestMonth = monthOptions[0] || {};
  var latestYear = yearOptions[0] || {};
  var creator = buildCreatorName(source, members);
  var inviteCode = normalizeFamilyId(firstDefined([source.inviteCode, source.invite_code, source.code], familyId));
  var memberCount = toNumber(
    firstDefined([source.memberCount, source.member_count, source.members_count, source.membersCount], members.length)
  );
  var monthIncome = toNumber(
    firstDefined(
      [
        source.monthIncome,
        source.month_income,
        source.currentMonthIncome,
        source.current_month_income,
        source.totalMonthIncome,
        source.total_month_income,
        latestMonth.income
      ],
      0
    )
  );
  var monthExpense = toNumber(
    firstDefined(
      [
        source.monthExpense,
        source.month_expense,
        source.currentMonthExpense,
        source.current_month_expense,
        source.totalMonthExpense,
        source.total_month_expense,
        latestMonth.expense
      ],
      0
    )
  );
  var monthBalance = toNumber(
    firstDefined([source.monthBalance, source.month_balance, source.currentMonthBalance, latestMonth.balance], monthIncome - monthExpense)
  );
  var yearBalance = toNumber(
    firstDefined([source.yearBalance, source.year_balance, source.currentYearBalance, latestYear.balance], 0)
  );

  return {
    id: familyId || inviteCode,
    name: firstDefined([source.name, source.family_name, source.familyName], "未命名家庭"),
    slogan: firstDefined([source.slogan, source.description, source.intro, source.remark], "暂无家庭介绍"),
    creator: creator,
    createdAt: formatDateLabel(
      firstDefined([source.createdAt, source.created_at, source.create_time, source.createTime, source.created_date], "")
    ),
    inviteCode: inviteCode || familyId,
    inviteLink: buildInviteLink(
      firstDefined([source.inviteLink, source.invite_link, source.share_link, source.shareLink, source.invitation_link], ""),
      familyId,
      inviteCode
    ),
    monthIncome: monthIncome,
    monthExpense: monthExpense,
    monthBalance: monthBalance,
    yearBalance: yearBalance,
    memberCount: memberCount || members.length,
    members: members,
    monthOptions: monthOptions,
    yearOptions: yearOptions
  };
}

function normalizeOverview(source, families) {
  var familyList = Array.isArray(families) ? families : [];
  var overviewSource = source || {};

  return {
    familyCount: toNumber(firstDefined([overviewSource.familyCount, overviewSource.family_count], familyList.length)),
    totalMembers: toNumber(
      firstDefined(
        [overviewSource.totalMembers, overviewSource.total_members],
        familyList.reduce(function(sum, family) {
          return sum + toNumber(family.memberCount);
        }, 0)
      )
    ),
    joinedCount: toNumber(firstDefined([overviewSource.joinedCount, overviewSource.joined_count], familyList.length))
  };
}

function resolveFamilyPayload(payload) {
  if (Array.isArray(payload)) {
    return payload;
  }

  return firstDefined([payload.list, payload.families, payload.rows, payload.items], []);
}

function deriveFallbackPeriodOptions(source, periodType) {
  var fallback = normalizePeriodOptions(
    firstDefined(
      [
        source[periodType === "year" ? "periods" : "periods"],
        source[periodType === "year" ? "year_stats" : "month_stats"],
        source[periodType === "year" ? "yearStats" : "monthStats"]
      ],
      []
    ),
    periodType
  );

  return fallback;
}

function buildSinglePeriodOption(source, periodType) {
  if (!source || typeof source !== "object") {
    return [];
  }

  var key = normalizePeriodKey(
    firstDefined([source.key, source.value, source.periodKey, source.period_key, source.month, source.year], ""),
    periodType
  );

  if (!key) {
    return [];
  }

  return [normalizePeriodOption(source, 0, periodType)];
}

export function getUserFamilies() {
  return request.get("/user/families");
}

export function createUserFamily(payload) {
  return request.post("/user/families", payload);
}

export function joinUserFamilyById(payload) {
  return request.post("/user/families/join", payload);
}

export function joinUserFamilyByInviteLink(payload) {
  return request.post("/user/families/join-by-link", payload);
}

export function leaveUserFamily(familyId, payload) {
  return request.post("/user/families/" + encodeURIComponent(familyId) + "/leave", payload);
}

export function getUserFamilyDetail(familyId) {
  return request.get("/user/families/" + encodeURIComponent(familyId));
}

export function getUserFamilyMemberShare(familyId, params) {
  return request.get("/user/families/" + encodeURIComponent(familyId) + "/member-share", {
    params: params
  });
}

export function buildUserFamilyError(error, fallback) {
  if (error && error.response && error.response.data) {
    return error.response.data.message || error.response.data.msg || fallback;
  }

  return fallback;
}

export function normalizeFamilyListPayload(result) {
  var payload = extractPayload(result);
  var families = resolveFamilyPayload(payload).map(function(item) {
    return normalizeFamilySummary(item);
  }).filter(function(item) {
    return !!item.id;
  });

  return {
    families: families,
    overview: normalizeOverview(payload.overview || payload.summary, families)
  };
}

export function normalizeFamilyMutationPayload(result) {
  var payload = extractPayload(result);
  var familySource = firstDefined([payload.family, payload.item, payload.detail], payload);
  var family = null;

  if (familySource && typeof familySource === "object" && !Array.isArray(familySource)) {
    var normalizedFamily = normalizeFamilySummary(familySource);
    family = normalizedFamily.id ? normalizedFamily : null;
  }

  return {
    family: family,
    message: firstDefined([payload.message, payload.msg], "")
  };
}

export function normalizeFamilyDetailPayload(result, fallbackFamilyId) {
  var payload = extractPayload(result);
  var raw = firstDefined([payload.family, payload.detail, payload.item], payload);

  if (!raw || typeof raw !== "object" || Array.isArray(raw)) {
    return null;
  }

  var family = normalizeFamilySummary(
    Object.assign({}, raw, {
      id: firstDefined([raw.id, raw.family_id, raw.familyId], fallbackFamilyId)
    })
  );
  var monthOptions = family.monthOptions.length
    ? family.monthOptions
    : deriveFallbackPeriodOptions(raw, "month").length
      ? deriveFallbackPeriodOptions(raw, "month")
      : buildSinglePeriodOption(raw.currentMonth || raw.current_month || raw.month, "month");
  var yearOptions = family.yearOptions.length
    ? family.yearOptions
    : deriveFallbackPeriodOptions(raw, "year").length
      ? deriveFallbackPeriodOptions(raw, "year")
      : buildSinglePeriodOption(raw.currentYear || raw.current_year || raw.year, "year");

  return Object.assign({}, family, {
    monthOptions: monthOptions,
    yearOptions: yearOptions,
    memberCount: family.members.length || family.memberCount
  });
}

function normalizeShareRows(source, family) {
  var familyMembers = family && Array.isArray(family.members) ? family.members : [];
  var rows = [];

  if (Array.isArray(source)) {
    rows = source;
  } else if (source && typeof source === "object") {
    rows = Object.keys(source).map(function(key) {
      var item = source[key];
      return typeof item === "object" && item !== null
        ? Object.assign({ name: key }, item)
        : { name: key, value: item };
    });
  }

  return rows.map(function(item, index) {
    var sourceItem = item || {};
    var name = firstDefined(
      [
        sourceItem.name,
        sourceItem.member_name,
        sourceItem.memberName,
        sourceItem.user_name,
        sourceItem.userName,
        sourceItem.nickname
      ],
      "成员" + (index + 1)
    );
    var matchedMember = familyMembers.find(function(member) {
      return member.name === name;
    });

    return {
      name: name,
      role: firstDefined([sourceItem.role, sourceItem.member_role, sourceItem.memberRole], matchedMember ? matchedMember.role : "成员"),
      color: firstDefined([sourceItem.color, sourceItem.avatar_color, sourceItem.avatarColor], matchedMember ? matchedMember.color : pickFamilyColor(name, index)),
      value: toNumber(
        firstDefined(
          [sourceItem.value, sourceItem.amount, sourceItem.total, sourceItem.expense, sourceItem.income, sourceItem.share],
          0
        )
      ),
      percent: toNumber(firstDefined([sourceItem.percent, sourceItem.ratio], 0))
    };
  });
}

function buildShareTitle(periodType, periodKey, metricType) {
  var periodLabel = buildPeriodLabel(periodType, periodKey);
  return periodLabel + (metricType === "income" ? "成员收入占比" : "成员支出占比");
}

export function normalizeFamilySharePayload(result, params, family) {
  var payload = extractPayload(result);
  var rows = normalizeShareRows(
    firstDefined([payload.rows, payload.list, payload.items, payload.members, payload.member_share, payload.memberShare], []),
    family
  );
  var total = toNumber(firstDefined([payload.total, payload.amount, payload.value], 0));

  if (!total) {
    total = rows.reduce(function(sum, item) {
      return sum + item.value;
    }, 0);
  }

  return {
    title: firstDefined([payload.title, payload.label], buildShareTitle(params.periodType, params.periodKey, params.metricType)),
    total: total,
    rows: rows.map(function(item) {
      var percent = item.percent;

      if (!percent && total > 0) {
        percent = Math.round((item.value / total) * 1000) / 10;
      }

      return Object.assign({}, item, {
        percent: percent
      });
    })
  };
}

export { formatFamilyCurrency, formatFamilyPercent };
