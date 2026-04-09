const STORAGE_KEY = "bill_admin_user_mock_rows";

const DEFAULT_ROWS = [
  {
    userId: "u-zhangxiaoming",
    username: "zhangxiaoming",
    nickname: "张小明",
    phone: "13800000001",
    email: "zhangxm@email.com",
    status: "启用",
    registerDate: "2026-01-15",
    lastLogin: "2026-04-09 10:30",
    billCount: 156,
    todayNew: false
  },
  {
    userId: "u-liming",
    username: "liming",
    nickname: "李明",
    phone: "13800000002",
    email: "liming@email.com",
    status: "禁用",
    registerDate: "2026-04-09",
    lastLogin: "-",
    billCount: 0,
    todayNew: true
  },
  {
    userId: "u-wangfang",
    username: "wangfang",
    nickname: "王芳",
    phone: "13800000003",
    email: "wangfang@email.com",
    status: "禁用",
    registerDate: "2026-04-09",
    lastLogin: "-",
    billCount: 0,
    todayNew: true
  },
  {
    userId: "u-chenhao",
    username: "chenhao",
    nickname: "陈浩",
    phone: "13800000004",
    email: "chenhao@email.com",
    status: "启用",
    registerDate: "2026-03-20",
    lastLogin: "2026-04-08 18:20",
    billCount: 89,
    todayNew: false
  },
  {
    userId: "u-zhoumin",
    username: "zhoumin",
    nickname: "周敏",
    phone: "13800000005",
    email: "zhoumin@email.com",
    status: "启用",
    registerDate: "2026-03-05",
    lastLogin: "2026-04-07 16:30",
    billCount: 67,
    todayNew: false
  }
];

const USER_BILL_DETAIL = {
  "u-zhangxiaoming": {
    profile: {
      userId: "u-zhangxiaoming",
      username: "zhangxiaoming",
      nickname: "张小明",
      status: "启用",
      registerDate: "2026-01-15",
      billCount: 156
    },
    monthOptions: ["2026-04", "2026-03", "2026-02"],
    yearOptions: ["2026", "2025"],
    monthly: {
      "2026-04": {
        income: 18500,
        expense: 8932,
        balance: 9568,
        records: 68,
        insight: "4 月支出结构已明显收口。"
      },
      "2026-03": {
        income: 24800,
        expense: 15280,
        balance: 9520,
        records: 186,
        insight: "3 月购物集中释放，已在 4 月回落。"
      },
      "2026-02": {
        income: 24800,
        expense: 12680,
        balance: 12120,
        records: 158,
        insight: "2 月整体节奏平稳，结余表现较优。"
      }
    },
    yearly: {
      "2026": {
        income: 215600,
        expense: 118768,
        balance: 96832,
        months: 4,
        insight: "本年仍保持净结余。"
      },
      "2025": {
        income: 297600,
        expense: 186320,
        balance: 111280,
        months: 12,
        insight: "去年四季度结余改善明显。"
      }
    }
  }
};

function safeParse(raw) {
  if (!raw) {
    return null;
  }

  try {
    return JSON.parse(raw);
  } catch (error) {
    return null;
  }
}

function cloneRows(rows) {
  return rows.map(function(item) {
    return Object.assign({}, item);
  });
}

function writeRows(rows) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(rows));
}

function readRows() {
  const storedRows = safeParse(localStorage.getItem(STORAGE_KEY));

  if (!storedRows || !Array.isArray(storedRows)) {
    writeRows(DEFAULT_ROWS);
    return cloneRows(DEFAULT_ROWS);
  }

  return cloneRows(storedRows);
}

export function getAdminUsersMock() {
  return readRows();
}

export function toggleAdminUserStatus(userId) {
  const rows = readRows();
  const nextRows = rows.map(function(row) {
    if (row.userId !== userId) {
      return row;
    }

    const nextStatus = row.status === "启用" ? "禁用" : "启用";
    const nextLastLogin = nextStatus === "禁用" ? "-" : row.lastLogin === "-" ? "2026-04-09 21:49" : row.lastLogin;

    return Object.assign({}, row, {
      status: nextStatus,
      lastLogin: nextLastLogin
    });
  });

  writeRows(nextRows);
  return cloneRows(nextRows);
}

function getFallbackUser(userId) {
  const rows = readRows();
  const match = rows.find(function(item) {
    return item.userId === userId;
  }) || rows[0];

  return {
    profile: {
      userId: match.userId,
      username: match.username,
      nickname: match.nickname,
      status: match.status,
      registerDate: match.registerDate,
      billCount: match.billCount
    },
    monthOptions: ["2026-04", "2026-03"],
    yearOptions: ["2026", "2025"],
    monthly: {
      "2026-04": { income: 14800, expense: 9200, balance: 5600, records: 52, insight: "账单分布稳定。" },
      "2026-03": { income: 15200, expense: 9800, balance: 5400, records: 49, insight: "支出较上月略高。" }
    },
    yearly: {
      "2026": { income: 168000, expense: 109000, balance: 59000, months: 4, insight: "年度仍保持结余。" },
      "2025": { income: 240000, expense: 176000, balance: 64000, months: 12, insight: "去年全年账单稳定。" }
    }
  };
}

export function getAdminUserDetailMock(userId) {
  const detail = USER_BILL_DETAIL[userId] || getFallbackUser(userId);
  const rows = readRows();
  const row = rows.find(function(item) {
    return item.userId === detail.profile.userId;
  });

  const mergedProfile = row
    ? Object.assign({}, detail.profile, {
        status: row.status,
        nickname: row.nickname,
        username: row.username,
        billCount: row.billCount
      })
    : Object.assign({}, detail.profile);

  return Object.assign({}, detail, {
    profile: mergedProfile
  });
}

export function formatAdminCurrency(value) {
  return "¥" + Number(value || 0).toLocaleString("zh-CN");
}
