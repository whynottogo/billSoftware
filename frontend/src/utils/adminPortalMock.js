function formatAdminCurrency(value) {
  return "¥" + Number(value || 0).toLocaleString("zh-CN");
}

const dashboardData = {
  stats: [
    {
      key: "totalUsers",
      label: "总用户数",
      value: 520,
      trend: "+12.5%",
      tone: "up"
    },
    {
      key: "pendingUsers",
      label: "待审批用户",
      value: 12,
      trend: "需处理",
      tone: "alert"
    },
    {
      key: "todayBills",
      label: "今日账单数",
      value: 367,
      trend: "+8.3%",
      tone: "up"
    },
    {
      key: "families",
      label: "家庭数量",
      value: 45,
      trend: "+5.2%",
      tone: "up"
    }
  ],
  userGrowth: {
    labels: ["1月", "2月", "3月", "4月", "5月", "6月"],
    series: [
      {
        name: "新增用户",
        color: "#171717",
        values: [120, 185, 245, 320, 410, 520]
      }
    ]
  },
  billTrend: {
    labels: ["04-04", "04-05", "04-06", "04-07", "04-08", "04-09"],
    series: [
      {
        name: "账单数",
        color: "#f6d34a",
        values: [256, 312, 289, 345, 298, 367]
      }
    ]
  },
  recentUsers: [
    {
      id: 1,
      name: "李明",
      email: "liming@email.com",
      status: "pending",
      registerDate: "2026-04-09 10:30"
    },
    {
      id: 2,
      name: "王芳",
      email: "wangfang@email.com",
      status: "pending",
      registerDate: "2026-04-09 09:15"
    },
    {
      id: 3,
      name: "陈浩",
      email: "chenhao@email.com",
      status: "active",
      registerDate: "2026-04-08 16:20"
    },
    {
      id: 4,
      name: "张伟",
      email: "zhangwei@email.com",
      status: "active",
      registerDate: "2026-04-08 14:45"
    },
    {
      id: 5,
      name: "刘洋",
      email: "liuyang@email.com",
      status: "inactive",
      registerDate: "2026-04-07 11:30"
    }
  ]
};

const pendingApprovalsData = {
  summary: {
    pendingCount: 12,
    selectedCount: 0
  },
  users: [
    {
      id: 1,
      name: "李明",
      email: "liming@email.com",
      phone: "138****1234",
      registerDate: "2026-04-09 10:30",
      reason: "个人记账"
    },
    {
      id: 2,
      name: "王芳",
      email: "wangfang@email.com",
      phone: "139****5678",
      registerDate: "2026-04-09 09:15",
      reason: "家庭财务管理"
    },
    {
      id: 3,
      name: "赵强",
      email: "zhaoqiang@email.com",
      phone: "136****9012",
      registerDate: "2026-04-09 08:45",
      reason: "个人记账"
    },
    {
      id: 4,
      name: "孙丽",
      email: "sunli@email.com",
      phone: "137****3456",
      registerDate: "2026-04-08 20:30",
      reason: "家庭财务管理"
    },
    {
      id: 5,
      name: "周敏",
      email: "zhoumin@email.com",
      phone: "135****7890",
      registerDate: "2026-04-08 18:15",
      reason: "个人记账"
    },
    {
      id: 6,
      name: "吴浩",
      email: "wuhao@email.com",
      phone: "133****2345",
      registerDate: "2026-04-08 16:00",
      reason: "个人记账"
    },
    {
      id: 7,
      name: "郑芳",
      email: "zhengfang@email.com",
      phone: "134****6789",
      registerDate: "2026-04-08 14:30",
      reason: "家庭财务管理"
    },
    {
      id: 8,
      name: "冯伟",
      email: "fengwei@email.com",
      phone: "188****0123",
      registerDate: "2026-04-08 11:20",
      reason: "个人记账"
    }
  ]
};

const familyManagementData = {
  stats: [
    {
      label: "家庭总数",
      value: "45"
    },
    {
      label: "总成员数",
      value: "156"
    },
    {
      label: "平均成员数",
      value: "3.5"
    },
    {
      label: "活跃家庭",
      value: "42"
    }
  ],
  families: [
    {
      id: "FAM-1001",
      name: "张家",
      creator: "张小明",
      members: 3,
      billCount: 456,
      totalAssets: 125680,
      createDate: "2026-01-15",
      status: "活跃"
    },
    {
      id: "FAM-1002",
      name: "李家",
      creator: "李明华",
      members: 4,
      billCount: 328,
      totalAssets: 98500,
      createDate: "2026-02-08",
      status: "活跃"
    },
    {
      id: "FAM-1003",
      name: "王家",
      creator: "王建国",
      members: 2,
      billCount: 189,
      totalAssets: 67800,
      createDate: "2026-02-20",
      status: "活跃"
    },
    {
      id: "FAM-1004",
      name: "陈家",
      creator: "陈浩",
      members: 5,
      billCount: 612,
      totalAssets: 156300,
      createDate: "2026-01-22",
      status: "活跃"
    },
    {
      id: "FAM-1005",
      name: "刘家",
      creator: "刘洋",
      members: 2,
      billCount: 145,
      totalAssets: 45200,
      createDate: "2026-03-10",
      status: "活跃"
    },
    {
      id: "FAM-1006",
      name: "赵家",
      creator: "赵丽",
      members: 3,
      billCount: 278,
      totalAssets: 89600,
      createDate: "2026-02-28",
      status: "活跃"
    }
  ],
  pagination: {
    total: 45,
    currentPage: 1,
    pages: [1, 2, 3]
  }
};

function getAdminDashboardData() {
  return dashboardData;
}

function getPendingApprovalsData() {
  return pendingApprovalsData;
}

function getFamilyManagementData() {
  return familyManagementData;
}

export {
  formatAdminCurrency,
  getAdminDashboardData,
  getPendingApprovalsData,
  getFamilyManagementData
};
