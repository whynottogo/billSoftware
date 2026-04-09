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
import SimpleDonutChart from "@/components/SimpleDonutChart.vue";
import SimpleGroupedBarChart from "@/components/SimpleGroupedBarChart.vue";
import SimpleLineChart from "@/components/SimpleLineChart.vue";
import { formatCurrency, getBillMonth } from "@/utils/userFinanceMock";

export default {
  name: "UserBillDetail",
  components: {
    SimpleDonutChart,
    SimpleGroupedBarChart,
    SimpleLineChart
  },
  computed: {
    monthData() {
      return getBillMonth(this.$route.params.month);
    },
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
  methods: {
    formatCurrency: formatCurrency,
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
