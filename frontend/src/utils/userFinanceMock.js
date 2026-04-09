function createAchievement(title, value, tone, hint) {
  return {
    title: title,
    value: value,
    tone: tone,
    hint: hint
  };
}

const billYears = [
  {
    year: 2026,
    summary: {
      balance: 96832,
      income: 215600,
      expense: 118768,
      months: 4,
      days: 99,
      records: 607,
      insight: "收入保持稳定，4 月开始主动压低高波动消费。"
    },
    months: [
      {
        key: "2026-04",
        label: "2026年4月",
        year: 2026,
        month: 4,
        income: 18500,
        expense: 8932,
        balance: 9568,
        days: 9,
        records: 68,
        status: "本月",
        note: "连续 9 天保持记账，支出峰值集中在餐饮与购物。",
        highlight: "本月已经明显收缩了娱乐型支出。",
        categorySplit: [
          { name: "餐饮", value: 3200, percent: 35.8, color: "#f6d34a" },
          { name: "购物", value: 2800, percent: 31.4, color: "#6bcf7c" },
          { name: "交通", value: 1280, percent: 14.3, color: "#4d96ff" },
          { name: "娱乐", value: 980, percent: 11.0, color: "#ff8b8b" },
          { name: "其他", value: 672, percent: 7.5, color: "#9b8cff" }
        ],
        ranking: [
          { name: "餐饮", value: 3200, percent: 35.8, count: 23, trend: "up", badge: "餐" },
          { name: "购物", value: 2800, percent: 31.4, count: 15, trend: "down", badge: "购" },
          { name: "交通", value: 1280, percent: 14.3, count: 18, trend: "up", badge: "交" },
          { name: "娱乐", value: 980, percent: 11.0, count: 8, trend: "up", badge: "娱" },
          { name: "其他", value: 672, percent: 7.5, count: 4, trend: "down", badge: "杂" }
        ],
        dailyTrend: [
          { label: "04-01", amount: 856 },
          { label: "04-02", amount: 1280 },
          { label: "04-03", amount: 680 },
          { label: "04-04", amount: 1420 },
          { label: "04-05", amount: 2856 },
          { label: "04-06", amount: 320 },
          { label: "04-07", amount: 428 },
          { label: "04-08", amount: 1280 },
          { label: "04-09", amount: 458 }
        ],
        highlightCards: [
          { label: "支出峰值", value: "¥2,856", hint: "04-05 周末聚会" },
          { label: "最高频分类", value: "餐饮", hint: "23 笔" },
          { label: "预算占用", value: "59.5%", hint: "剩余空间充足" }
        ],
        achievements: [
          createAchievement("连续记账", "9天", "danger", "保持连续性"),
          createAchievement("本月记账", "68笔", "brand", "比上月快 12 笔"),
          createAchievement("记账天数", "9/30天", "success", "节奏稳定"),
          createAchievement("预算达成", "85%", "info", "已收口超支分类")
        ]
      },
      {
        key: "2026-03",
        label: "2026年3月",
        year: 2026,
        month: 3,
        income: 24800,
        expense: 15280,
        balance: 9520,
        days: 31,
        records: 186,
        status: "已结账",
        note: "3 月有一次集中采买，购物和居家支出较高。",
        highlight: "收入保持稳定，但购物支出冲高导致结余收窄。",
        categorySplit: [
          { name: "购物", value: 4820, percent: 31.5, color: "#f6d34a" },
          { name: "餐饮", value: 3580, percent: 23.4, color: "#6bcf7c" },
          { name: "居家", value: 2460, percent: 16.1, color: "#4d96ff" },
          { name: "交通", value: 2210, percent: 14.5, color: "#ff8b8b" },
          { name: "其他", value: 2210, percent: 14.5, color: "#9b8cff" }
        ],
        ranking: [
          { name: "购物", value: 4820, percent: 31.5, count: 24, trend: "up", badge: "购" },
          { name: "餐饮", value: 3580, percent: 23.4, count: 26, trend: "up", badge: "餐" },
          { name: "居家", value: 2460, percent: 16.1, count: 12, trend: "up", badge: "居" },
          { name: "交通", value: 2210, percent: 14.5, count: 18, trend: "down", badge: "交" },
          { name: "其他", value: 2210, percent: 14.5, count: 9, trend: "down", badge: "杂" }
        ],
        dailyTrend: [
          { label: "03-05", amount: 480 },
          { label: "03-09", amount: 920 },
          { label: "03-12", amount: 1380 },
          { label: "03-16", amount: 860 },
          { label: "03-20", amount: 2220 },
          { label: "03-23", amount: 1180 },
          { label: "03-27", amount: 840 },
          { label: "03-31", amount: 760 }
        ],
        highlightCards: [
          { label: "最大单笔", value: "¥1,980", hint: "生活用品补货" },
          { label: "高频分类", value: "餐饮", hint: "26 笔" },
          { label: "预算占用", value: "101.9%", hint: "购物轻微超支" }
        ],
        achievements: [
          createAchievement("连续记账", "31天", "success", "整月无中断"),
          createAchievement("本月记账", "186笔", "brand", "账单密度高"),
          createAchievement("月结余", "¥9,520", "info", "保持正向"),
          createAchievement("预算预警", "2类", "warning", "购物与餐饮靠前")
        ]
      },
      {
        key: "2026-02",
        label: "2026年2月",
        year: 2026,
        month: 2,
        income: 24800,
        expense: 12680,
        balance: 12120,
        days: 28,
        records: 158,
        status: "已结账",
        note: "2 月节奏平稳，交通和社交活动有所下降。",
        highlight: "结余保持全年高位，支出控制表现最好。",
        categorySplit: [
          { name: "餐饮", value: 2980, percent: 23.5, color: "#f6d34a" },
          { name: "购物", value: 2760, percent: 21.8, color: "#6bcf7c" },
          { name: "交通", value: 1840, percent: 14.5, color: "#4d96ff" },
          { name: "社交", value: 1620, percent: 12.8, color: "#ff8b8b" },
          { name: "其他", value: 3480, percent: 27.4, color: "#9b8cff" }
        ],
        ranking: [
          { name: "餐饮", value: 2980, percent: 23.5, count: 20, trend: "down", badge: "餐" },
          { name: "购物", value: 2760, percent: 21.8, count: 14, trend: "down", badge: "购" },
          { name: "交通", value: 1840, percent: 14.5, count: 16, trend: "down", badge: "交" },
          { name: "社交", value: 1620, percent: 12.8, count: 9, trend: "up", badge: "社" },
          { name: "其他", value: 3480, percent: 27.4, count: 12, trend: "up", badge: "杂" }
        ],
        dailyTrend: [
          { label: "02-03", amount: 680 },
          { label: "02-06", amount: 820 },
          { label: "02-10", amount: 1080 },
          { label: "02-14", amount: 1320 },
          { label: "02-18", amount: 980 },
          { label: "02-22", amount: 1650 },
          { label: "02-26", amount: 760 },
          { label: "02-28", amount: 540 }
        ],
        highlightCards: [
          { label: "支出低谷", value: "¥540", hint: "02-28" },
          { label: "高频分类", value: "餐饮", hint: "20 笔" },
          { label: "预算占用", value: "84.5%", hint: "控制优秀" }
        ],
        achievements: [
          createAchievement("连续记账", "28天", "success", "整月无缺口"),
          createAchievement("月结余", "¥12,120", "brand", "全年最佳"),
          createAchievement("日均支出", "¥453", "info", "显著下降"),
          createAchievement("预算达成", "96%", "success", "留出充足余量")
        ]
      },
      {
        key: "2026-01",
        label: "2026年1月",
        year: 2026,
        month: 1,
        income: 32500,
        expense: 18560,
        balance: 13940,
        days: 31,
        records: 195,
        status: "已结账",
        note: "年初收入高，但节假日消费明显增加。",
        highlight: "收入拉高了结余，但购物、出行和聚会支出集中释放。",
        categorySplit: [
          { name: "购物", value: 4860, percent: 26.2, color: "#f6d34a" },
          { name: "餐饮", value: 4520, percent: 24.4, color: "#6bcf7c" },
          { name: "交通", value: 3240, percent: 17.5, color: "#4d96ff" },
          { name: "旅行", value: 2880, percent: 15.5, color: "#ff8b8b" },
          { name: "其他", value: 3060, percent: 16.4, color: "#9b8cff" }
        ],
        ranking: [
          { name: "购物", value: 4860, percent: 26.2, count: 28, trend: "up", badge: "购" },
          { name: "餐饮", value: 4520, percent: 24.4, count: 31, trend: "up", badge: "餐" },
          { name: "交通", value: 3240, percent: 17.5, count: 22, trend: "up", badge: "交" },
          { name: "旅行", value: 2880, percent: 15.5, count: 7, trend: "up", badge: "旅" },
          { name: "其他", value: 3060, percent: 16.4, count: 11, trend: "down", badge: "杂" }
        ],
        dailyTrend: [
          { label: "01-03", amount: 1280 },
          { label: "01-07", amount: 760 },
          { label: "01-11", amount: 1680 },
          { label: "01-15", amount: 2260 },
          { label: "01-19", amount: 1380 },
          { label: "01-23", amount: 2980 },
          { label: "01-27", amount: 2080 },
          { label: "01-31", amount: 1160 }
        ],
        highlightCards: [
          { label: "最高单日", value: "¥2,980", hint: "01-23" },
          { label: "高频分类", value: "餐饮", hint: "31 笔" },
          { label: "预算占用", value: "123.7%", hint: "假期支出抬升" }
        ],
        achievements: [
          createAchievement("连续记账", "31天", "success", "开年稳定"),
          createAchievement("月收入", "¥32,500", "brand", "全年峰值"),
          createAchievement("记账总笔", "195笔", "info", "记录完整"),
          createAchievement("预算提醒", "3类", "warning", "假期消费偏高")
        ]
      }
    ]
  },
  {
    year: 2025,
    summary: {
      balance: 185640,
      income: 297600,
      expense: 186320,
      months: 12,
      days: 342,
      records: 2156,
      insight: "年度支出维持在可控区间，四季度结余改善明显。"
    },
    months: [
      {
        key: "2025-12",
        label: "2025年12月",
        year: 2025,
        month: 12,
        income: 28600,
        expense: 16840,
        balance: 11760,
        days: 31,
        records: 203,
        status: "已结账",
        note: "年末聚会带来餐饮上扬，但整体仍留出健康结余。",
        highlight: "节庆消费增加，但购物控制比去年更稳。",
        categorySplit: [
          { name: "餐饮", value: 4660, percent: 27.7, color: "#f6d34a" },
          { name: "购物", value: 3820, percent: 22.7, color: "#6bcf7c" },
          { name: "交通", value: 2720, percent: 16.2, color: "#4d96ff" },
          { name: "社交", value: 2480, percent: 14.7, color: "#ff8b8b" },
          { name: "其他", value: 3160, percent: 18.7, color: "#9b8cff" }
        ],
        ranking: [
          { name: "餐饮", value: 4660, percent: 27.7, count: 33, trend: "up", badge: "餐" },
          { name: "购物", value: 3820, percent: 22.7, count: 21, trend: "down", badge: "购" },
          { name: "交通", value: 2720, percent: 16.2, count: 19, trend: "up", badge: "交" },
          { name: "社交", value: 2480, percent: 14.7, count: 13, trend: "up", badge: "社" },
          { name: "其他", value: 3160, percent: 18.7, count: 10, trend: "down", badge: "杂" }
        ],
        dailyTrend: [
          { label: "12-04", amount: 880 },
          { label: "12-09", amount: 1320 },
          { label: "12-13", amount: 1580 },
          { label: "12-18", amount: 980 },
          { label: "12-22", amount: 1820 },
          { label: "12-26", amount: 2280 },
          { label: "12-29", amount: 1240 },
          { label: "12-31", amount: 980 }
        ],
        highlightCards: [
          { label: "最高单笔", value: "¥1,880", hint: "年末聚餐" },
          { label: "高频分类", value: "餐饮", hint: "33 笔" },
          { label: "预算占用", value: "93.4%", hint: "控制尚可" }
        ],
        achievements: [
          createAchievement("连续记账", "31天", "success", "全年收尾"),
          createAchievement("月结余", "¥11,760", "brand", "保持健康"),
          createAchievement("记账总笔", "203笔", "info", "信息完整"),
          createAchievement("预算达成", "94%", "success", "未超支")
        ]
      },
      {
        key: "2025-11",
        label: "2025年11月",
        year: 2025,
        month: 11,
        income: 24800,
        expense: 14280,
        balance: 10520,
        days: 30,
        records: 178,
        status: "已结账",
        note: "11 月购物回落，社交和居家支出更稳定。",
        highlight: "购物热度开始下降，结余继续改善。",
        categorySplit: [
          { name: "餐饮", value: 3260, percent: 22.8, color: "#f6d34a" },
          { name: "购物", value: 3040, percent: 21.3, color: "#6bcf7c" },
          { name: "交通", value: 2180, percent: 15.3, color: "#4d96ff" },
          { name: "居家", value: 2520, percent: 17.6, color: "#ff8b8b" },
          { name: "其他", value: 3280, percent: 23.0, color: "#9b8cff" }
        ],
        ranking: [
          { name: "餐饮", value: 3260, percent: 22.8, count: 28, trend: "up", badge: "餐" },
          { name: "购物", value: 3040, percent: 21.3, count: 16, trend: "down", badge: "购" },
          { name: "交通", value: 2180, percent: 15.3, count: 18, trend: "down", badge: "交" },
          { name: "居家", value: 2520, percent: 17.6, count: 10, trend: "up", badge: "居" },
          { name: "其他", value: 3280, percent: 23.0, count: 9, trend: "up", badge: "杂" }
        ],
        dailyTrend: [
          { label: "11-03", amount: 620 },
          { label: "11-08", amount: 1020 },
          { label: "11-12", amount: 860 },
          { label: "11-17", amount: 1480 },
          { label: "11-21", amount: 980 },
          { label: "11-24", amount: 1640 },
          { label: "11-28", amount: 1320 },
          { label: "11-30", amount: 760 }
        ],
        highlightCards: [
          { label: "最大单日", value: "¥1,640", hint: "11-24" },
          { label: "高频分类", value: "餐饮", hint: "28 笔" },
          { label: "预算占用", value: "88.2%", hint: "继续改善" }
        ],
        achievements: [
          createAchievement("连续记账", "30天", "success", "持续稳定"),
          createAchievement("月结余", "¥10,520", "brand", "回升明显"),
          createAchievement("日均支出", "¥476", "info", "维持中位"),
          createAchievement("预算达成", "95%", "success", "控制良好")
        ]
      },
      {
        key: "2025-10",
        label: "2025年10月",
        year: 2025,
        month: 10,
        income: 24800,
        expense: 13896,
        balance: 10904,
        days: 31,
        records: 192,
        status: "已结账",
        note: "10 月假期出行带来短期抬升，但整体节奏稳。",
        highlight: "旅行支出短期放大，但被餐饮回落抵消。",
        categorySplit: [
          { name: "旅行", value: 3180, percent: 22.9, color: "#f6d34a" },
          { name: "餐饮", value: 2960, percent: 21.3, color: "#6bcf7c" },
          { name: "交通", value: 2420, percent: 17.4, color: "#4d96ff" },
          { name: "购物", value: 2360, percent: 17.0, color: "#ff8b8b" },
          { name: "其他", value: 2976, percent: 21.4, color: "#9b8cff" }
        ],
        ranking: [
          { name: "旅行", value: 3180, percent: 22.9, count: 8, trend: "up", badge: "旅" },
          { name: "餐饮", value: 2960, percent: 21.3, count: 27, trend: "down", badge: "餐" },
          { name: "交通", value: 2420, percent: 17.4, count: 18, trend: "up", badge: "交" },
          { name: "购物", value: 2360, percent: 17.0, count: 14, trend: "down", badge: "购" },
          { name: "其他", value: 2976, percent: 21.4, count: 10, trend: "up", badge: "杂" }
        ],
        dailyTrend: [
          { label: "10-02", amount: 1680 },
          { label: "10-05", amount: 2260 },
          { label: "10-09", amount: 980 },
          { label: "10-13", amount: 860 },
          { label: "10-18", amount: 1360 },
          { label: "10-22", amount: 980 },
          { label: "10-27", amount: 1520 },
          { label: "10-31", amount: 640 }
        ],
        highlightCards: [
          { label: "旅行花费", value: "¥3,180", hint: "短期拉高" },
          { label: "高频分类", value: "餐饮", hint: "27 笔" },
          { label: "预算占用", value: "90.6%", hint: "出行可控" }
        ],
        achievements: [
          createAchievement("连续记账", "31天", "success", "覆盖假期"),
          createAchievement("结余", "¥10,904", "brand", "健康收口"),
          createAchievement("旅行预算", "91%", "info", "控制住了"),
          createAchievement("消费节奏", "平稳", "success", "波动不大")
        ]
      },
      {
        key: "2025-09",
        label: "2025年9月",
        year: 2025,
        month: 9,
        income: 24800,
        expense: 14300,
        balance: 10500,
        days: 30,
        records: 181,
        status: "已结账",
        note: "9 月是四季度的起点，节奏偏稳。",
        highlight: "总体维持在预算内，支出结构均衡。",
        categorySplit: [
          { name: "餐饮", value: 3180, percent: 22.2, color: "#f6d34a" },
          { name: "购物", value: 2800, percent: 19.6, color: "#6bcf7c" },
          { name: "交通", value: 2240, percent: 15.7, color: "#4d96ff" },
          { name: "社交", value: 1980, percent: 13.8, color: "#ff8b8b" },
          { name: "其他", value: 4100, percent: 28.7, color: "#9b8cff" }
        ],
        ranking: [
          { name: "餐饮", value: 3180, percent: 22.2, count: 24, trend: "up", badge: "餐" },
          { name: "购物", value: 2800, percent: 19.6, count: 15, trend: "down", badge: "购" },
          { name: "交通", value: 2240, percent: 15.7, count: 16, trend: "up", badge: "交" },
          { name: "社交", value: 1980, percent: 13.8, count: 10, trend: "down", badge: "社" },
          { name: "其他", value: 4100, percent: 28.7, count: 11, trend: "up", badge: "杂" }
        ],
        dailyTrend: [
          { label: "09-02", amount: 720 },
          { label: "09-07", amount: 980 },
          { label: "09-11", amount: 1460 },
          { label: "09-16", amount: 860 },
          { label: "09-20", amount: 1620 },
          { label: "09-24", amount: 1120 },
          { label: "09-27", amount: 1280 },
          { label: "09-30", amount: 760 }
        ],
        highlightCards: [
          { label: "日均支出", value: "¥477", hint: "节奏均衡" },
          { label: "高频分类", value: "餐饮", hint: "24 笔" },
          { label: "预算占用", value: "89.8%", hint: "留有余量" }
        ],
        achievements: [
          createAchievement("连续记账", "30天", "success", "稳定起步"),
          createAchievement("月结余", "¥10,500", "brand", "表现平稳"),
          createAchievement("日均支出", "¥477", "info", "控制正常"),
          createAchievement("预算达成", "95%", "success", "状态良好")
        ]
      }
    ]
  },
  {
    year: 2024,
    summary: {
      balance: 174360,
      income: 285000,
      expense: 178920,
      months: 12,
      days: 336,
      records: 2089,
      insight: "2024 年消费结构更分散，年底开始明显形成收口意识。"
    },
    months: [
      {
        key: "2024-12",
        label: "2024年12月",
        year: 2024,
        month: 12,
        income: 27200,
        expense: 16540,
        balance: 10660,
        days: 31,
        records: 190,
        status: "已结账",
        note: "2024 年末开始收口高波动支出。",
        highlight: "年底已出现预算意识，结余趋稳。",
        categorySplit: [
          { name: "餐饮", value: 4100, percent: 24.8, color: "#f6d34a" },
          { name: "购物", value: 3520, percent: 21.3, color: "#6bcf7c" },
          { name: "交通", value: 2560, percent: 15.5, color: "#4d96ff" },
          { name: "社交", value: 1980, percent: 12.0, color: "#ff8b8b" },
          { name: "其他", value: 4380, percent: 26.4, color: "#9b8cff" }
        ],
        ranking: [
          { name: "餐饮", value: 4100, percent: 24.8, count: 28, trend: "up", badge: "餐" },
          { name: "购物", value: 3520, percent: 21.3, count: 18, trend: "down", badge: "购" },
          { name: "交通", value: 2560, percent: 15.5, count: 16, trend: "up", badge: "交" },
          { name: "社交", value: 1980, percent: 12.0, count: 9, trend: "down", badge: "社" },
          { name: "其他", value: 4380, percent: 26.4, count: 12, trend: "up", badge: "杂" }
        ],
        dailyTrend: [
          { label: "12-03", amount: 860 },
          { label: "12-08", amount: 1160 },
          { label: "12-12", amount: 1540 },
          { label: "12-17", amount: 980 },
          { label: "12-21", amount: 1760 },
          { label: "12-25", amount: 1480 },
          { label: "12-28", amount: 1240 },
          { label: "12-31", amount: 880 }
        ],
        highlightCards: [
          { label: "最高单日", value: "¥1,760", hint: "12-21" },
          { label: "高频分类", value: "餐饮", hint: "28 笔" },
          { label: "预算占用", value: "92.1%", hint: "已开始收口" }
        ],
        achievements: [
          createAchievement("连续记账", "31天", "success", "稳定收官"),
          createAchievement("月结余", "¥10,660", "brand", "回升明显"),
          createAchievement("支出结构", "均衡", "info", "分类分散"),
          createAchievement("预算意识", "增强", "success", "收口清晰")
        ]
      },
      {
        key: "2024-11",
        label: "2024年11月",
        year: 2024,
        month: 11,
        income: 23500,
        expense: 14980,
        balance: 8520,
        days: 30,
        records: 174,
        status: "已结账",
        note: "11 月维持中位水平，支出结构逐步稳定。",
        highlight: "餐饮和购物仍靠前，但超支幅度收窄。",
        categorySplit: [
          { name: "餐饮", value: 3620, percent: 24.2, color: "#f6d34a" },
          { name: "购物", value: 3180, percent: 21.2, color: "#6bcf7c" },
          { name: "交通", value: 2480, percent: 16.6, color: "#4d96ff" },
          { name: "娱乐", value: 2020, percent: 13.5, color: "#ff8b8b" },
          { name: "其他", value: 3680, percent: 24.5, color: "#9b8cff" }
        ],
        ranking: [
          { name: "餐饮", value: 3620, percent: 24.2, count: 25, trend: "up", badge: "餐" },
          { name: "购物", value: 3180, percent: 21.2, count: 17, trend: "up", badge: "购" },
          { name: "交通", value: 2480, percent: 16.6, count: 16, trend: "down", badge: "交" },
          { name: "娱乐", value: 2020, percent: 13.5, count: 8, trend: "down", badge: "娱" },
          { name: "其他", value: 3680, percent: 24.5, count: 11, trend: "up", badge: "杂" }
        ],
        dailyTrend: [
          { label: "11-02", amount: 720 },
          { label: "11-07", amount: 980 },
          { label: "11-11", amount: 1880 },
          { label: "11-15", amount: 920 },
          { label: "11-19", amount: 1460 },
          { label: "11-23", amount: 1180 },
          { label: "11-27", amount: 1380 },
          { label: "11-30", amount: 760 }
        ],
        highlightCards: [
          { label: "购物峰值", value: "¥1,880", hint: "双 11 集中" },
          { label: "高频分类", value: "餐饮", hint: "25 笔" },
          { label: "预算占用", value: "97.5%", hint: "接近上限" }
        ],
        achievements: [
          createAchievement("连续记账", "30天", "success", "节奏稳定"),
          createAchievement("双11控制", "可接受", "info", "波动有限"),
          createAchievement("月结余", "¥8,520", "brand", "仍为正向"),
          createAchievement("预算提醒", "2类", "warning", "购物偏高")
        ]
      },
      {
        key: "2024-10",
        label: "2024年10月",
        year: 2024,
        month: 10,
        income: 23600,
        expense: 15120,
        balance: 8480,
        days: 31,
        records: 176,
        status: "已结账",
        note: "10 月出行增加，但未出现大幅失控。",
        highlight: "节假日预算仍偏紧，开始形成预算提醒意识。",
        categorySplit: [
          { name: "旅行", value: 3320, percent: 22.0, color: "#f6d34a" },
          { name: "餐饮", value: 3140, percent: 20.8, color: "#6bcf7c" },
          { name: "购物", value: 2820, percent: 18.7, color: "#4d96ff" },
          { name: "交通", value: 2480, percent: 16.4, color: "#ff8b8b" },
          { name: "其他", value: 3360, percent: 22.2, color: "#9b8cff" }
        ],
        ranking: [
          { name: "旅行", value: 3320, percent: 22.0, count: 7, trend: "up", badge: "旅" },
          { name: "餐饮", value: 3140, percent: 20.8, count: 26, trend: "down", badge: "餐" },
          { name: "购物", value: 2820, percent: 18.7, count: 15, trend: "up", badge: "购" },
          { name: "交通", value: 2480, percent: 16.4, count: 17, trend: "up", badge: "交" },
          { name: "其他", value: 3360, percent: 22.2, count: 11, trend: "down", badge: "杂" }
        ],
        dailyTrend: [
          { label: "10-01", amount: 1820 },
          { label: "10-05", amount: 2260 },
          { label: "10-09", amount: 980 },
          { label: "10-14", amount: 860 },
          { label: "10-18", amount: 1280 },
          { label: "10-22", amount: 980 },
          { label: "10-27", amount: 1520 },
          { label: "10-30", amount: 820 }
        ],
        highlightCards: [
          { label: "旅行预算", value: "¥3,320", hint: "国庆假期" },
          { label: "高频分类", value: "餐饮", hint: "26 笔" },
          { label: "预算占用", value: "100.3%", hint: "轻微超支" }
        ],
        achievements: [
          createAchievement("连续记账", "31天", "success", "覆盖节假日"),
          createAchievement("旅行支出", "可控", "info", "波动有限"),
          createAchievement("月结余", "¥8,480", "brand", "稳定输出"),
          createAchievement("预算提醒", "1类", "warning", "轻微超支")
        ]
      },
      {
        key: "2024-09",
        label: "2024年9月",
        year: 2024,
        month: 9,
        income: 22800,
        expense: 14260,
        balance: 8540,
        days: 30,
        records: 167,
        status: "已结账",
        note: "9 月支出结构分散，是 2024 年较有代表性的稳态月份。",
        highlight: "作为基线月份，消费结构最均衡。",
        categorySplit: [
          { name: "餐饮", value: 3180, percent: 22.3, color: "#f6d34a" },
          { name: "购物", value: 2840, percent: 19.9, color: "#6bcf7c" },
          { name: "交通", value: 2380, percent: 16.7, color: "#4d96ff" },
          { name: "娱乐", value: 1840, percent: 12.9, color: "#ff8b8b" },
          { name: "其他", value: 4020, percent: 28.2, color: "#9b8cff" }
        ],
        ranking: [
          { name: "餐饮", value: 3180, percent: 22.3, count: 22, trend: "up", badge: "餐" },
          { name: "购物", value: 2840, percent: 19.9, count: 14, trend: "up", badge: "购" },
          { name: "交通", value: 2380, percent: 16.7, count: 15, trend: "down", badge: "交" },
          { name: "娱乐", value: 1840, percent: 12.9, count: 7, trend: "down", badge: "娱" },
          { name: "其他", value: 4020, percent: 28.2, count: 10, trend: "up", badge: "杂" }
        ],
        dailyTrend: [
          { label: "09-03", amount: 680 },
          { label: "09-08", amount: 980 },
          { label: "09-12", amount: 1360 },
          { label: "09-16", amount: 820 },
          { label: "09-20", amount: 1440 },
          { label: "09-24", amount: 1120 },
          { label: "09-27", amount: 1280 },
          { label: "09-30", amount: 760 }
        ],
        highlightCards: [
          { label: "基线月份", value: "稳态", hint: "适合作为对比" },
          { label: "高频分类", value: "餐饮", hint: "22 笔" },
          { label: "预算占用", value: "95.1%", hint: "略高但可控" }
        ],
        achievements: [
          createAchievement("连续记账", "30天", "success", "稳态持续"),
          createAchievement("月结余", "¥8,540", "brand", "保持正向"),
          createAchievement("结构均衡", "优秀", "info", "分布均匀"),
          createAchievement("预算达成", "95%", "success", "可控范围")
        ]
      }
    ]
  }
];

const monthlyBudget = {
  key: "2026-04",
  label: "2026年4月预算",
  notice: "仅支持设置当前月份预算，年预算请进入独立页面查看。",
  overview: {
    totalBudget: 15000,
    totalExpense: 8932,
    remaining: 6068,
    percentage: 59.5
  },
  highlights: [
    { label: "已超支分类", value: "2个", hint: "餐饮、购物" },
    { label: "可控分类", value: "4个", hint: "仍有余量" },
    { label: "预算节奏", value: "稳中可控", hint: "本月可继续收口" }
  ],
  categories: [
    { name: "餐饮", badge: "餐", budget: 3000, expense: 3200, remaining: -200, percentage: 106.7, status: "over", note: "午餐外食频率高于预期。" },
    { name: "购物", badge: "购", budget: 2500, expense: 2800, remaining: -300, percentage: 112.0, status: "over", note: "集中补货拉高了本月消费。" },
    { name: "交通", badge: "交", budget: 1500, expense: 1280, remaining: 220, percentage: 85.3, status: "warning", note: "本月已经接近上限。" },
    { name: "娱乐", badge: "娱", budget: 1200, expense: 980, remaining: 220, percentage: 81.7, status: "warning", note: "再新增娱乐支出需谨慎。" },
    { name: "运动", badge: "运", budget: 800, expense: 0, remaining: 800, percentage: 0, status: "safe", note: "本月暂无该分类支出。" },
    { name: "通讯", badge: "讯", budget: 300, expense: 168, remaining: 132, percentage: 56.0, status: "safe", note: "维持在正常波动范围。" }
  ]
};

const yearlyBudgets = [
  {
    year: 2026,
    overview: {
      totalBudget: 180000,
      totalExpense: 118768,
      remaining: 61232,
      percentage: 66.0,
      months: 4,
      note: "4 月开始出现更强的预算意识，超支分类正在减少。"
    },
    monthlyExecution: [
      { label: "1月", budget: 15000, expense: 18560 },
      { label: "2月", budget: 15000, expense: 12680 },
      { label: "3月", budget: 15000, expense: 15280 },
      { label: "4月", budget: 15000, expense: 8932 },
      { label: "5月", budget: 15000, expense: 0 },
      { label: "6月", budget: 15000, expense: 0 },
      { label: "7月", budget: 15000, expense: 0 },
      { label: "8月", budget: 15000, expense: 0 },
      { label: "9月", budget: 15000, expense: 0 },
      { label: "10月", budget: 15000, expense: 0 },
      { label: "11月", budget: 15000, expense: 0 },
      { label: "12月", budget: 15000, expense: 0 }
    ],
    categories: [
      { name: "餐饮", badge: "餐", budget: 36000, expense: 38400, remaining: -2400, percentage: 106.7, status: "over" },
      { name: "购物", badge: "购", budget: 24000, expense: 18320, remaining: 5680, percentage: 76.3, status: "safe" },
      { name: "交通", badge: "交", budget: 18000, expense: 16240, remaining: 1760, percentage: 90.2, status: "warning" },
      { name: "娱乐", badge: "娱", budget: 12000, expense: 8960, remaining: 3040, percentage: 74.7, status: "safe" },
      { name: "运动", badge: "运", budget: 12000, expense: 5680, remaining: 6320, percentage: 47.3, status: "safe" },
      { name: "通讯", badge: "讯", budget: 6000, expense: 3920, remaining: 2080, percentage: 65.3, status: "safe" }
    ]
  },
  {
    year: 2025,
    overview: {
      totalBudget: 168000,
      totalExpense: 186320,
      remaining: -18320,
      percentage: 110.9,
      months: 12,
      note: "2025 年四季度控制已经开始见效，但全年仍被上半年拉高。"
    },
    monthlyExecution: [
      { label: "1月", budget: 14000, expense: 15620 },
      { label: "2月", budget: 14000, expense: 14860 },
      { label: "3月", budget: 14000, expense: 15240 },
      { label: "4月", budget: 14000, expense: 16420 },
      { label: "5月", budget: 14000, expense: 15980 },
      { label: "6月", budget: 14000, expense: 15160 },
      { label: "7月", budget: 14000, expense: 16280 },
      { label: "8月", budget: 14000, expense: 15420 },
      { label: "9月", budget: 14000, expense: 14300 },
      { label: "10月", budget: 14000, expense: 13896 },
      { label: "11月", budget: 14000, expense: 14280 },
      { label: "12月", budget: 14000, expense: 16840 }
    ],
    categories: [
      { name: "餐饮", badge: "餐", budget: 32000, expense: 35400, remaining: -3400, percentage: 110.6, status: "over" },
      { name: "购物", badge: "购", budget: 28000, expense: 29780, remaining: -1780, percentage: 106.4, status: "over" },
      { name: "交通", badge: "交", budget: 22000, expense: 21240, remaining: 760, percentage: 96.5, status: "warning" },
      { name: "娱乐", badge: "娱", budget: 14000, expense: 12620, remaining: 1380, percentage: 90.1, status: "warning" },
      { name: "社交", badge: "社", budget: 12000, expense: 10980, remaining: 1020, percentage: 91.5, status: "warning" },
      { name: "其他", badge: "杂", budget: 60000, expense: 56300, remaining: 3700, percentage: 93.8, status: "warning" }
    ]
  },
  {
    year: 2024,
    overview: {
      totalBudget: 156000,
      totalExpense: 178920,
      remaining: -22920,
      percentage: 114.7,
      months: 12,
      note: "2024 年预算使用偏激进，是后续开始修正预算策略的起点。"
    },
    monthlyExecution: [
      { label: "1月", budget: 13000, expense: 14680 },
      { label: "2月", budget: 13000, expense: 13860 },
      { label: "3月", budget: 13000, expense: 15420 },
      { label: "4月", budget: 13000, expense: 14980 },
      { label: "5月", budget: 13000, expense: 15140 },
      { label: "6月", budget: 13000, expense: 14780 },
      { label: "7月", budget: 13000, expense: 15360 },
      { label: "8月", budget: 13000, expense: 14620 },
      { label: "9月", budget: 13000, expense: 14260 },
      { label: "10月", budget: 13000, expense: 15120 },
      { label: "11月", budget: 13000, expense: 14980 },
      { label: "12月", budget: 13000, expense: 16540 }
    ],
    categories: [
      { name: "餐饮", badge: "餐", budget: 30000, expense: 33240, remaining: -3240, percentage: 110.8, status: "over" },
      { name: "购物", badge: "购", budget: 24000, expense: 27220, remaining: -3220, percentage: 113.4, status: "over" },
      { name: "交通", badge: "交", budget: 20000, expense: 21460, remaining: -1460, percentage: 107.3, status: "over" },
      { name: "娱乐", badge: "娱", budget: 12000, expense: 12840, remaining: -840, percentage: 107.0, status: "over" },
      { name: "旅行", badge: "旅", budget: 10000, expense: 10980, remaining: -980, percentage: 109.8, status: "over" },
      { name: "其他", badge: "杂", budget: 60000, expense: 73180, remaining: -13180, percentage: 122.0, status: "over" }
    ]
  }
];

function flattenBillMonths() {
  var months = [];

  billYears.forEach(function(yearItem) {
    yearItem.months.forEach(function(monthItem) {
      months.push(monthItem);
    });
  });

  months.sort(function(a, b) {
    return a.key.localeCompare(b.key);
  });

  return months;
}

export function formatCurrency(value) {
  return "¥" + Number(value || 0).toLocaleString("zh-CN");
}

export function formatSignedCurrency(value) {
  var amount = Number(value || 0);
  var prefix = amount >= 0 ? "+" : "-";
  return prefix + formatCurrency(Math.abs(amount));
}

export function formatPercent(value) {
  var parsed = Number(value || 0).toFixed(1);
  return parsed.replace(".0", "") + "%";
}

export function getBillYears() {
  return billYears.map(function(item) {
    return item.year;
  });
}

export function getBillYear(year) {
  var matched = billYears.find(function(item) {
    return item.year === year;
  });

  return matched || billYears[0];
}

export function getBillYearHistory() {
  return billYears;
}

export function getBillMonth(monthKey) {
  var months = flattenBillMonths();
  var matched = months.find(function(item) {
    return item.key === monthKey;
  });
  var current = matched || months[months.length - 1];
  var index = months.findIndex(function(item) {
    return item.key === current.key;
  });
  var previous = index > 0 ? months[index - 1] : null;
  var next = index < months.length - 1 ? months[index + 1] : null;
  var comparison = months.slice(Math.max(0, index - 5), index + 1).map(function(item) {
    return {
      label: item.month + "月",
      income: item.income,
      expense: item.expense,
      balance: item.balance
    };
  });

  return Object.assign({}, current, {
    previousKey: previous ? previous.key : "",
    previousLabel: previous ? previous.label : "",
    previousBalance: previous ? previous.balance : 0,
    nextKey: next ? next.key : "",
    nextLabel: next ? next.label : "",
    comparison: comparison
  });
}

export function getMonthlyBudget() {
  return monthlyBudget;
}

export function getBudgetYears() {
  return yearlyBudgets.map(function(item) {
    return item.year;
  });
}

export function getYearlyBudget(year) {
  var matched = yearlyBudgets.find(function(item) {
    return item.year === year;
  });

  return matched || yearlyBudgets[0];
}
