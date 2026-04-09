const FAMILY_COLORS = ["#f6d34a", "#6bcf7c", "#4d96ff", "#ff8b8b", "#9b8cff", "#34d399", "#f97316"];

const familyStore = {
  nextId: 4900,
  families: [
    {
      id: "FAM-4821",
      name: "周末小家",
      slogan: "把每周的生活支出一起看清楚。",
      creator: "林芷",
      createdAt: "2025-08-12",
      inviteCode: "FAM-4821",
      inviteLink: "https://bill.local/invite/FAM-4821",
      monthOptions: [
        { key: "2026-04", label: "2026年4月", note: "清明假期后恢复常规节奏。" },
        { key: "2026-03", label: "2026年3月", note: "聚餐和出行都偏高。" },
        { key: "2026-02", label: "2026年2月", note: "春节回家，交通占比提升。" }
      ],
      yearOptions: [
        { key: "2026", label: "2026年", note: "当前年度仍在持续记账中。" },
        { key: "2025", label: "2025年", note: "上一年度已完成结账。" }
      ],
      members: [
        {
          name: "林芷",
          role: "创建人",
          color: "#f6d34a",
          monthStats: {
            "2026-04": { income: 18200, expense: 7260 },
            "2026-03": { income: 17600, expense: 7380 },
            "2026-02": { income: 16900, expense: 6920 }
          },
          yearStats: {
            "2026": { income: 121000, expense: 54600 },
            "2025": { income: 242500, expense: 132800 }
          }
        },
        {
          name: "许川",
          role: "成员",
          color: "#6bcf7c",
          monthStats: {
            "2026-04": { income: 13600, expense: 9120 },
            "2026-03": { income: 13400, expense: 9580 },
            "2026-02": { income: 12800, expense: 8440 }
          },
          yearStats: {
            "2026": { income: 90300, expense: 61500 },
            "2025": { income: 186200, expense: 117400 }
          }
        },
        {
          name: "周澈",
          role: "成员",
          color: "#4d96ff",
          monthStats: {
            "2026-04": { income: 9200, expense: 4860 },
            "2026-03": { income: 8600, expense: 4380 },
            "2026-02": { income: 9300, expense: 3980 }
          },
          yearStats: {
            "2026": { income: 58200, expense: 28600 },
            "2025": { income: 120500, expense: 65400 }
          }
        }
      ]
    },
    {
      id: "FAM-4768",
      name: "三餐计划组",
      slogan: "专门追踪餐饮预算和买菜节奏。",
      creator: "陈愿",
      createdAt: "2025-11-03",
      inviteCode: "FAM-4768",
      inviteLink: "https://bill.local/invite/FAM-4768",
      monthOptions: [
        { key: "2026-04", label: "2026年4月", note: "外卖次数在下降。" },
        { key: "2026-03", label: "2026年3月", note: "春季聚餐频率提升。" },
        { key: "2026-02", label: "2026年2月", note: "春节采购开销明显。" }
      ],
      yearOptions: [
        { key: "2026", label: "2026年", note: "以餐饮预算收口为主。" },
        { key: "2025", label: "2025年", note: "全年消费波动较小。" }
      ],
      members: [
        {
          name: "陈愿",
          role: "创建人",
          color: "#ff8b8b",
          monthStats: {
            "2026-04": { income: 7800, expense: 4280 },
            "2026-03": { income: 7600, expense: 4960 },
            "2026-02": { income: 7300, expense: 5520 }
          },
          yearStats: {
            "2026": { income: 48800, expense: 30900 },
            "2025": { income: 96100, expense: 61100 }
          }
        },
        {
          name: "苏芽",
          role: "成员",
          color: "#9b8cff",
          monthStats: {
            "2026-04": { income: 11200, expense: 5180 },
            "2026-03": { income: 10600, expense: 5860 },
            "2026-02": { income: 9800, expense: 6120 }
          },
          yearStats: {
            "2026": { income: 71100, expense: 34700 },
            "2025": { income: 142200, expense: 68900 }
          }
        }
      ]
    }
  ]
};

function deepClone(value) {
  return JSON.parse(JSON.stringify(value));
}

function formatCurrency(value) {
  return "¥" + Number(value || 0).toLocaleString("zh-CN");
}

function formatPercent(value) {
  return Number(value || 0).toFixed(1) + "%";
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

function getMemberValueByPeriod(member, periodType, periodKey, metricType) {
  var metric = metricType === "income" ? "income" : "expense";
  var periodBucket = periodType === "year" ? member.yearStats : member.monthStats;

  if (!periodBucket || !periodBucket[periodKey]) {
    return 0;
  }

  return Number(periodBucket[periodKey][metric] || 0);
}

function refreshFamilyTotals(family) {
  family.monthOptions.forEach(function(option) {
    var monthIncome = 0;
    var monthExpense = 0;

    family.members.forEach(function(member) {
      monthIncome += getMemberValueByPeriod(member, "month", option.key, "income");
      monthExpense += getMemberValueByPeriod(member, "month", option.key, "expense");
    });

    option.income = monthIncome;
    option.expense = monthExpense;
    option.balance = monthIncome - monthExpense;
  });

  family.yearOptions.forEach(function(option) {
    var yearIncome = 0;
    var yearExpense = 0;

    family.members.forEach(function(member) {
      yearIncome += getMemberValueByPeriod(member, "year", option.key, "income");
      yearExpense += getMemberValueByPeriod(member, "year", option.key, "expense");
    });

    option.income = yearIncome;
    option.expense = yearExpense;
    option.balance = yearIncome - yearExpense;
  });

  family.memberCount = family.members.length;
}

function getFamilyById(familyId) {
  return familyStore.families.find(function(item) {
    return item.id === familyId;
  });
}

function normalizeFamilyId(value) {
  return String(value || "").trim().toUpperCase();
}

function hydrateFamilies() {
  familyStore.families.forEach(function(family) {
    refreshFamilyTotals(family);
  });
}

function getSessionMemberName() {
  var rawProfile = localStorage.getItem("bill_user_profile");

  if (!rawProfile) {
    return "playwright-user";
  }

  try {
    var profile = JSON.parse(rawProfile);
    return profile.nickname || profile.account || profile.username || "playwright-user";
  } catch (error) {
    return "playwright-user";
  }
}

function createGeneratedMember(name, family) {
  var seed = hashText(name + family.id);
  var color = FAMILY_COLORS[seed % FAMILY_COLORS.length];
  var monthStats = {};
  var yearStats = {};

  family.monthOptions.forEach(function(option, index) {
    var incomeBase = 2800 + (seed % 1200) + index * 120;
    var expenseBase = 1400 + (seed % 800) + index * 90;

    monthStats[option.key] = {
      income: incomeBase,
      expense: expenseBase
    };
  });

  family.yearOptions.forEach(function(option, index) {
    var incomeBase = 28000 + (seed % 12000) + index * 1300;
    var expenseBase = 13600 + (seed % 7000) + index * 900;

    yearStats[option.key] = {
      income: incomeBase,
      expense: expenseBase
    };
  });

  return {
    name: name,
    role: "成员",
    color: color,
    monthStats: monthStats,
    yearStats: yearStats
  };
}

function buildFamilySummary(family) {
  var latestMonth = family.monthOptions[0];
  var latestYear = family.yearOptions[0];

  return {
    id: family.id,
    name: family.name,
    slogan: family.slogan,
    creator: family.creator,
    createdAt: family.createdAt,
    inviteCode: family.inviteCode,
    inviteLink: family.inviteLink,
    memberCount: family.memberCount,
    monthIncome: latestMonth.income,
    monthExpense: latestMonth.expense,
    monthBalance: latestMonth.balance,
    yearIncome: latestYear.income,
    yearExpense: latestYear.expense,
    yearBalance: latestYear.balance,
    members: family.members.map(function(member) {
      return {
        name: member.name,
        role: member.role,
        color: member.color
      };
    })
  };
}

function listFamilies() {
  return deepClone(
    familyStore.families.map(function(family) {
      return buildFamilySummary(family);
    })
  );
}

function getFamilyOverview() {
  var families = familyStore.families;
  var totalMembers = families.reduce(function(sum, family) {
    return sum + family.members.length;
  }, 0);
  var currentUser = getSessionMemberName();
  var joinedCount = families.filter(function(family) {
    return family.members.some(function(member) {
      return member.name === currentUser;
    });
  }).length;

  return {
    familyCount: families.length,
    totalMembers: totalMembers,
    joinedCount: joinedCount
  };
}

function createFamily(payload) {
  var familyName = String(payload.name || "").trim();
  var creator = String(payload.creator || getSessionMemberName()).trim();
  var slogan = String(payload.slogan || "").trim();

  if (!familyName) {
    return {
      ok: false,
      message: "家庭名称不能为空。"
    };
  }

  familyStore.nextId += 1;

  var familyId = "FAM-" + familyStore.nextId;
  var today = new Date();
  var month = String(today.getMonth() + 1).padStart(2, "0");
  var day = String(today.getDate()).padStart(2, "0");
  var year = String(today.getFullYear());
  var previousYear = String(today.getFullYear() - 1);
  var monthKey = year + "-" + month;
  var previousMonthKey = year + "-" + String(Math.max(1, Number(month) - 1)).padStart(2, "0");

  var family = {
    id: familyId,
    name: familyName,
    slogan: slogan || "新家庭刚刚建立，欢迎补齐成员并开始记账。",
    creator: creator,
    createdAt: year + "-" + month + "-" + day,
    inviteCode: familyId,
    inviteLink: "https://bill.local/invite/" + familyId,
    monthOptions: [
      { key: monthKey, label: year + "年" + Number(month) + "月", note: "新建家庭首月记录。" },
      { key: previousMonthKey, label: year + "年" + Number(previousMonthKey.split("-")[1]) + "月", note: "补录上月数据。" }
    ],
    yearOptions: [
      { key: year, label: year + "年", note: "当前年度统计。" },
      { key: previousYear, label: previousYear + "年", note: "上一年度统计。" }
    ],
    members: []
  };

  family.members.push(createGeneratedMember(creator, family));
  refreshFamilyTotals(family);
  familyStore.families.unshift(family);

  return {
    ok: true,
    family: deepClone(buildFamilySummary(family))
  };
}

function joinFamilyById(payload) {
  var familyId = normalizeFamilyId(payload.familyId);
  var memberName = String(payload.memberName || getSessionMemberName()).trim();

  if (!familyId) {
    return {
      ok: false,
      message: "请输入家庭 ID。"
    };
  }

  var family = getFamilyById(familyId);

  if (!family) {
    return {
      ok: false,
      message: "未找到该家庭，请检查家庭 ID。"
    };
  }

  if (!memberName) {
    return {
      ok: false,
      message: "请填写成员名称。"
    };
  }

  var exists = family.members.some(function(member) {
    return member.name === memberName;
  });

  if (exists) {
    return {
      ok: false,
      message: memberName + " 已经在家庭中。"
    };
  }

  family.members.push(createGeneratedMember(memberName, family));
  refreshFamilyTotals(family);

  return {
    ok: true,
    family: deepClone(buildFamilySummary(family))
  };
}

function parseFamilyIdFromInvite(link) {
  var text = String(link || "");
  var matched = text.match(/FAM-\d{4,}/i);

  if (!matched) {
    return "";
  }

  return normalizeFamilyId(matched[0]);
}

function joinFamilyByInviteLink(payload) {
  var familyId = parseFamilyIdFromInvite(payload.inviteLink);

  if (!familyId) {
    return {
      ok: false,
      message: "邀请链接格式不正确，示例：https://bill.local/invite/FAM-4821"
    };
  }

  return joinFamilyById({
    familyId: familyId,
    memberName: payload.memberName
  });
}

function leaveFamily(payload) {
  var familyId = normalizeFamilyId(payload.familyId);
  var memberName = String(payload.memberName || getSessionMemberName()).trim();
  var family = getFamilyById(familyId);

  if (!family) {
    return {
      ok: false,
      message: "家庭不存在，无法退出。"
    };
  }

  if (!memberName) {
    return {
      ok: false,
      message: "请先确认当前成员名称。"
    };
  }

  if (family.creator === memberName) {
    return {
      ok: false,
      message: "当前成员是创建人，原型阶段不支持创建人直接退出。"
    };
  }

  var targetIndex = family.members.findIndex(function(member) {
    return member.name === memberName;
  });

  if (targetIndex < 0) {
    return {
      ok: false,
      message: memberName + " 不在该家庭中。"
    };
  }

  family.members.splice(targetIndex, 1);
  refreshFamilyTotals(family);

  return {
    ok: true,
    family: deepClone(buildFamilySummary(family))
  };
}

function getFamilyDetail(familyId) {
  var family = getFamilyById(normalizeFamilyId(familyId));

  if (!family) {
    return null;
  }

  return deepClone(family);
}

function getFamilyMemberShare(familyId, periodType, periodKey, metricType) {
  var family = getFamilyById(normalizeFamilyId(familyId));

  if (!family) {
    return {
      title: "暂无成员占比",
      total: 0,
      rows: []
    };
  }

  var rows = family.members.map(function(member) {
    return {
      name: member.name,
      role: member.role,
      color: member.color,
      value: getMemberValueByPeriod(member, periodType, periodKey, metricType)
    };
  });

  var total = rows.reduce(function(sum, item) {
    return sum + item.value;
  }, 0);

  rows.forEach(function(item) {
    item.percent = total ? (item.value / total) * 100 : 0;
  });

  rows.sort(function(left, right) {
    return right.value - left.value;
  });

  var metricText = metricType === "income" ? "收入" : "支出";
  var periodText = periodType === "year" ? "年度" : "月度";

  return deepClone({
    title: periodKey + " " + periodText + metricText + "成员占比",
    total: total,
    rows: rows
  });
}

function getDefaultFamilyId() {
  if (!familyStore.families.length) {
    return "";
  }

  return familyStore.families[0].id;
}

hydrateFamilies();

export {
  formatCurrency,
  formatPercent,
  listFamilies,
  getFamilyOverview,
  createFamily,
  joinFamilyById,
  joinFamilyByInviteLink,
  leaveFamily,
  getFamilyDetail,
  getFamilyMemberShare,
  getDefaultFamilyId
};
