<template>
  <div class="finance-page bill-detail-page">
    <section class="finance-toolbar">
      <div class="finance-toolbar__actions">
        <button class="finance-button finance-button--ghost" @click="goBack">返回月账单列表</button>
        <button class="finance-button finance-button--ghost" :disabled="!monthData.previousKey" @click="goToMonth(monthData.previousKey)">
          上一月
        </button>
        <button class="finance-button finance-button--ghost" :disabled="!monthData.nextKey" @click="goToMonth(monthData.nextKey)">
          下一月
        </button>
      </div>

      <div class="finance-toolbar__actions">
        <span class="finance-pill">{{ monthData.days }} 天 · {{ monthData.records }} 笔记录</span>
        <button class="finance-button finance-button--primary" @click="goYearlyBill">查看年账单</button>
      </div>
    </section>

    <section class="finance-hero">
      <div class="finance-hero__eyebrow">
        <span class="finance-hero__eyebrow-mark">明</span>
        <span>{{ monthData.label }}</span>
      </div>

      <div class="finance-hero__headline">
        <h1>{{ monthData.label }}账单详情</h1>
        <p>{{ monthData.highlight }}</p>
      </div>

      <div class="finance-stat-grid">
        <article class="finance-stat-card">
          <span>本月结余</span>
          <strong>{{ formatCurrency(monthData.balance) }}</strong>
          <small>{{ monthData.note }}</small>
        </article>
        <article class="finance-stat-card">
          <span>上月结余</span>
          <strong>{{ formatCurrency(monthData.previousBalance) }}</strong>
          <small>{{ monthData.previousLabel || "已到最早月份" }}</small>
        </article>
        <article class="finance-stat-card">
          <span>本月收入</span>
          <strong class="finance-tone-income">{{ formatCurrency(monthData.income) }}</strong>
          <small>收入维持在稳定区间</small>
        </article>
        <article class="finance-stat-card">
          <span>本月支出</span>
          <strong class="finance-tone-expense">{{ formatCurrency(monthData.expense) }}</strong>
          <small>可继续按分类收口</small>
        </article>
      </div>
    </section>

    <section class="finance-grid-3">
      <article v-for="item in monthData.highlightCards" :key="item.label" class="page-card bill-detail-page__highlight">
        <span>{{ item.label }}</span>
        <strong>{{ item.value }}</strong>
        <p>{{ item.hint }}</p>
      </article>
    </section>

    <section class="finance-grid-2">
      <article class="page-card finance-panel">
        <header class="finance-panel__header">
          <div>
            <h3>支出分类分布</h3>
            <p>先看结构，再决定接下来优先收口哪一类支出。</p>
          </div>
        </header>

        <SimpleDonutChart
          :segments="monthData.categorySplit"
          center-title="总支出"
          :center-value="formatCurrency(monthData.expense)"
        />

        <div class="bill-detail-page__legend">
          <div v-for="item in monthData.categorySplit" :key="item.name" class="bill-detail-page__legend-item">
            <span class="bill-detail-page__legend-dot" :style="{ backgroundColor: item.color }"></span>
            <span>{{ item.name }}</span>
            <strong>{{ item.percent }}%</strong>
          </div>
        </div>
      </article>

      <article class="page-card finance-panel">
        <header class="finance-panel__header">
          <div>
            <h3>支出排行</h3>
            <p>从金额和笔数双维度看最值得优先关注的分类。</p>
          </div>
        </header>

        <div class="bill-detail-page__ranking">
          <article v-for="(item, index) in monthData.ranking" :key="item.name" class="ranking-row">
            <div class="ranking-row__index">{{ index + 1 }}</div>
            <div class="ranking-row__body">
              <div class="ranking-row__topline">
                <div class="ranking-row__title">
                  <span class="finance-badge">{{ item.badge }}</span>
                  <div>
                    <strong>{{ item.name }}</strong>
                    <small>{{ item.count }} 笔</small>
                  </div>
                </div>
                <div class="ranking-row__amount">
                  <strong>{{ formatCurrency(item.value) }}</strong>
                  <span :class="item.trend === 'up' ? 'finance-tone-expense' : 'finance-tone-income'">
                    {{ item.trend === "up" ? "较上月上升" : "较上月回落" }}
                  </span>
                </div>
              </div>

              <div class="finance-progress">
                <span class="ranking-row__progress" :style="{ width: item.percent + '%' }"></span>
              </div>
            </div>
          </article>
        </div>
      </article>
    </section>

    <article class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>本月支出趋势</h3>
          <p>用一条线看出本月什么时候开始偏离正常节奏。</p>
        </div>
      </header>
      <SimpleLineChart
        :labels="monthData.dailyTrend.map(function(item) { return item.label; })"
        :series="trendSeries"
        :height="280"
      />
    </article>

    <article class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>近 6 个月收支对比</h3>
          <p>当前月份与前 5 个月放在一起，更容易看清收入和支出的变化惯性。</p>
        </div>
      </header>
      <SimpleGroupedBarChart
        :labels="monthData.comparison.map(function(item) { return item.label; })"
        :series="comparisonSeries"
        :height="300"
      />
    </article>

    <section class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>记账成就</h3>
          <p>这不是装饰信息，它帮助你看见自己已经建立起来的习惯。</p>
        </div>
      </header>

      <div class="finance-grid-4">
        <article v-for="item in monthData.achievements" :key="item.title" class="bill-detail-page__achievement">
          <span class="bill-detail-page__achievement-mark" :class="toneClass(item.tone)">{{ achievementMark(item.tone) }}</span>
          <strong>{{ item.value }}</strong>
          <p>{{ item.title }}</p>
          <small>{{ item.hint }}</small>
        </article>
      </div>
    </section>
  </div>
</template>

<script>
import { getUserBillMonthDetail } from "@/api/userBillDetail";
import SimpleDonutChart from "@/components/SimpleDonutChart.vue";
import SimpleGroupedBarChart from "@/components/SimpleGroupedBarChart.vue";
import SimpleLineChart from "@/components/SimpleLineChart.vue";

const CATEGORY_COLORS = ["#f6d34a", "#6bcf7c", "#4d96ff", "#ff8b8b", "#9b8cff", "#fb7185"];

function pad(value) {
  return String(value).padStart(2, "0");
}

function getCurrentMonthKey() {
  var now = new Date();
  return now.getFullYear() + "-" + pad(now.getMonth() + 1);
}

function shiftMonthKey(monthKey, offset) {
  var parts = String(monthKey || "").split("-");
  var year = Number(parts[0]);
  var month = Number(parts[1]);
  var date = new Date(year, month - 1 + offset, 1);

  return date.getFullYear() + "-" + pad(date.getMonth() + 1);
}

function monthLabel(monthKey) {
  var parts = String(monthKey || "").split("-");
  if (parts.length !== 2) {
    return monthKey || "暂无数据";
  }

  return parts[0] + "年" + Number(parts[1]) + "月";
}

function formatCurrency(value) {
  return "¥" + Number(value || 0).toLocaleString("zh-CN");
}

function normalizeAmount(value) {
  var amount = Number(value || 0);
  return Number.isFinite(amount) ? amount : 0;
}

function normalizePercent(value) {
  var percent = Number(value || 0);
  return Number.isFinite(percent) ? percent : 0;
}

function resolveMonthKey(monthKey) {
  return /^\d{4}-\d{2}$/.test(String(monthKey || "")) ? String(monthKey) : getCurrentMonthKey();
}

function extractPayload(result) {
  if (result && result.data) {
    return result.data;
  }

  return result || {};
}

function buildDailyTrendPlaceholder(monthKey) {
  return [
    {
      label: monthKey.slice(5) || "-",
      amount: 0
    }
  ];
}

function buildComparisonPlaceholder(monthKey) {
  var list = [];
  var index;

  for (index = 5; index >= 0; index -= 1) {
    list.push({
      label: shiftMonthKey(monthKey, -index).slice(5),
      income: 0,
      expense: 0,
      balance: 0
    });
  }

  return list;
}

function buildMonthDetailPlaceholder(monthKey) {
  return {
    key: monthKey,
    label: monthLabel(monthKey),
    previousKey: null,
    nextKey: null,
    previousLabel: "",
    previousBalance: 0,
    income: 0,
    expense: 0,
    balance: 0,
    days: 0,
    records: 0,
    note: "本月暂无详细账单数据",
    highlight: "当前月份还没有形成可分析的账单特征。",
    categorySplit: [],
    ranking: [],
    dailyTrend: buildDailyTrendPlaceholder(monthKey),
    comparison: buildComparisonPlaceholder(monthKey),
    highlightCards: [],
    achievements: []
  };
}

function normalizeMonthDetail(result, fallbackMonthKey) {
  var payload = extractPayload(result);
  var monthKey = resolveMonthKey(payload.key || fallbackMonthKey);
  var categorySplit = Array.isArray(payload.categorySplit) ? payload.categorySplit : [];
  var ranking = Array.isArray(payload.ranking) ? payload.ranking : [];
  var dailyTrend = Array.isArray(payload.dailyTrend) ? payload.dailyTrend : [];
  var comparison = Array.isArray(payload.comparison) ? payload.comparison : [];
  var highlightCards = Array.isArray(payload.highlightCards) ? payload.highlightCards : [];
  var achievements = Array.isArray(payload.achievements) ? payload.achievements : [];

  return {
    key: monthKey,
    label: payload.label || monthLabel(monthKey),
    previousKey: payload.previousKey || null,
    nextKey: payload.nextKey || null,
    previousLabel: payload.previousLabel || (payload.previousKey ? monthLabel(payload.previousKey) : ""),
    previousBalance: normalizeAmount(payload.previousBalance),
    income: normalizeAmount(payload.income),
    expense: normalizeAmount(payload.expense),
    balance: normalizeAmount(payload.balance),
    days: Number(payload.days || 0),
    records: Number(payload.records || 0),
    note: payload.note || "本月暂无额外说明",
    highlight: payload.highlight || "当前月份账单详情已同步到真实数据。",
    categorySplit: categorySplit.map(function(item, index) {
      return {
        name: item.name || "未分类",
        value: normalizeAmount(item.value),
        percent: normalizePercent(item.percent),
        color: item.color || CATEGORY_COLORS[index % CATEGORY_COLORS.length]
      };
    }),
    ranking: ranking.map(function(item) {
      return {
        name: item.name || "未分类",
        value: normalizeAmount(item.value),
        percent: normalizePercent(item.percent),
        count: Number(item.count || 0),
        trend: item.trend === "up" ? "up" : "down",
        badge: item.badge || String(item.name || "未").slice(0, 1)
      };
    }),
    dailyTrend: (dailyTrend.length ? dailyTrend : buildDailyTrendPlaceholder(monthKey)).map(function(item) {
      return {
        label: item.label || "-",
        amount: normalizeAmount(item.amount)
      };
    }),
    comparison: (comparison.length ? comparison : buildComparisonPlaceholder(monthKey)).map(function(item) {
      return {
        label: item.label || "-",
        income: normalizeAmount(item.income),
        expense: normalizeAmount(item.expense),
        balance: normalizeAmount(item.balance)
      };
    }),
    highlightCards: highlightCards.map(function(item) {
      return {
        label: item.label || "-",
        value: item.value || "-",
        hint: item.hint || ""
      };
    }),
    achievements: achievements.map(function(item) {
      return {
        title: item.title || "-",
        value: item.value || "-",
        tone: item.tone || "info",
        hint: item.hint || ""
      };
    })
  };
}

function buildErrorMessage(error, fallback) {
  if (error && error.response && error.response.data && error.response.data.message) {
    return error.response.data.message;
  }

  return fallback;
}

export default {
  name: "UserBillDetail",
  components: {
    SimpleDonutChart,
    SimpleGroupedBarChart,
    SimpleLineChart
  },
  data() {
    return {
      monthData: buildMonthDetailPlaceholder(resolveMonthKey(this.$route.params.month)),
      detailLoading: false,
      lastRequestId: 0
    };
  },
  computed: {
    trendSeries() {
      return [
        {
          name: "支出",
          color: "#ef4444",
          values: this.monthData.dailyTrend.map(function(item) {
            return item.amount;
          })
        }
      ];
    },
    comparisonSeries() {
      return [
        {
          name: "收入",
          color: "#22c55e",
          values: this.monthData.comparison.map(function(item) {
            return item.income;
          })
        },
        {
          name: "支出",
          color: "#ef4444",
          values: this.monthData.comparison.map(function(item) {
            return item.expense;
          })
        }
      ];
    }
  },
  watch: {
    "$route.params.month": {
      immediate: true,
      handler(nextMonth) {
        this.loadMonthDetail(nextMonth);
      }
    }
  },
  methods: {
    formatCurrency: formatCurrency,
    loadMonthDetail(routeMonth) {
      var monthKey = resolveMonthKey(routeMonth);
      var requestId = this.lastRequestId + 1;

      this.lastRequestId = requestId;
      this.detailLoading = true;
      this.monthData = buildMonthDetailPlaceholder(monthKey);

      return getUserBillMonthDetail(monthKey)
        .then(
          function(result) {
            if (requestId !== this.lastRequestId) {
              return;
            }

            this.monthData = normalizeMonthDetail(result, monthKey);
          }.bind(this)
        )
        .catch(
          function(error) {
            if (requestId !== this.lastRequestId) {
              return;
            }

            this.monthData = buildMonthDetailPlaceholder(monthKey);
            this.$message.error(buildErrorMessage(error, "月账单详情加载失败，请稍后重试。"));
          }.bind(this)
        )
        .finally(
          function() {
            if (requestId === this.lastRequestId) {
              this.detailLoading = false;
            }
          }.bind(this)
        );
    },
    goBack() {
      this.$router.push("/user/bills/month");
    },
    goToMonth(monthKey) {
      if (monthKey) {
        this.$router.push("/user/bills/month/" + monthKey);
      }
    },
    goYearlyBill() {
      this.$router.push("/user/bills/year");
    },
    toneClass(tone) {
      return "tone-" + tone;
    },
    achievementMark(tone) {
      var marks = {
        danger: "燃",
        brand: "账",
        success: "稳",
        info: "进",
        warning: "警"
      };

      return marks[tone] || "记";
    }
  }
};
</script>

<style scoped>
.bill-detail-page__highlight {
  padding: 22px;
}

.bill-detail-page__highlight span {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
}

.bill-detail-page__highlight strong {
  display: block;
  margin-top: 12px;
  font-size: 28px;
}

.bill-detail-page__highlight p {
  margin: 10px 0 0;
  color: var(--text-subtle);
  line-height: 1.6;
}

.bill-detail-page__legend {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin-top: 18px;
}

.bill-detail-page__legend-item {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 42px;
  padding: 0 12px;
  border-radius: 14px;
  background: rgba(245, 246, 248, 0.92);
}

.bill-detail-page__legend-item strong {
  margin-left: auto;
  font-size: 13px;
}

.bill-detail-page__legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 999px;
}

.bill-detail-page__ranking {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.ranking-row {
  display: flex;
  gap: 14px;
  padding: 18px;
  border-radius: 20px;
  background: rgba(245, 246, 248, 0.9);
}

.ranking-row__index {
  width: 34px;
  height: 34px;
  border-radius: 12px;
  background: rgba(246, 211, 74, 0.24);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
}

.ranking-row__body {
  flex: 1;
}

.ranking-row__topline {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.ranking-row__title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.ranking-row__title strong,
.ranking-row__amount strong {
  display: block;
}

.ranking-row__title small,
.ranking-row__amount span {
  display: block;
  margin-top: 6px;
  color: var(--text-muted);
  font-size: 12px;
}

.ranking-row__amount {
  text-align: right;
}

.ranking-row__progress {
  background: linear-gradient(90deg, #ef4444 0%, #f6d34a 100%);
}

.bill-detail-page__achievement {
  padding: 22px;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.96);
  border: 1px solid var(--border-color);
  text-align: center;
}

.bill-detail-page__achievement strong {
  display: block;
  margin-top: 14px;
  font-size: 26px;
}

.bill-detail-page__achievement p,
.bill-detail-page__achievement small {
  display: block;
}

.bill-detail-page__achievement p {
  margin: 8px 0 0;
  font-weight: 700;
}

.bill-detail-page__achievement small {
  margin-top: 8px;
  color: var(--text-muted);
  line-height: 1.6;
}

.bill-detail-page__achievement-mark {
  width: 52px;
  height: 52px;
  border-radius: 18px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 800;
}

.tone-danger {
  background: rgba(239, 68, 68, 0.12);
  color: var(--danger-color);
}

.tone-brand {
  background: rgba(246, 211, 74, 0.2);
  color: #c79400;
}

.tone-success {
  background: rgba(34, 197, 94, 0.12);
  color: var(--success-color);
}

.tone-info {
  background: rgba(59, 130, 246, 0.12);
  color: var(--info-color);
}

.tone-warning {
  background: rgba(245, 158, 11, 0.12);
  color: var(--warning-color);
}

@media (max-width: 960px) {
  .ranking-row__topline {
    flex-direction: column;
    align-items: flex-start;
  }

  .ranking-row__amount {
    text-align: left;
  }

  .bill-detail-page__legend {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
