<template>
  <div class="finance-page user-chart-page">
    <section class="finance-toolbar">
      <div class="finance-switcher">
        <button :disabled="!hasPreviousYear" @click="prevYear">‹</button>
        <div>
          <strong>{{ selectedYear }} {{ isExpenseMode ? "支出图表" : "收入图表" }}</strong>
          <p class="user-chart-page__switcher-note">
            {{ isExpenseMode ? "支持月/年双趋势与分类排行" : "聚焦年收入趋势与收入结构" }}
          </p>
        </div>
        <button :disabled="!hasNextYear" @click="nextYear">›</button>
      </div>

      <div class="finance-toolbar__actions">
        <button
          :class="['finance-button', isExpenseMode ? 'finance-button--primary' : 'finance-button--ghost']"
          @click="switchMode('expense')"
        >
          支出图表
        </button>
        <button
          :class="['finance-button', !isExpenseMode ? 'finance-button--primary' : 'finance-button--ghost']"
          @click="switchMode('income')"
        >
          收入图表
        </button>
      </div>
    </section>

    <section class="finance-hero finance-hero--soft">
      <div class="finance-hero__eyebrow">
        <span class="finance-hero__eyebrow-mark">{{ isExpenseMode ? "支" : "收" }}</span>
        <span>{{ selectedYear }} 年{{ isExpenseMode ? "支出" : "收入" }}总览</span>
      </div>

      <div class="finance-hero__headline">
        <h1>{{ heroTitle }}</h1>
        <p>{{ heroDescription }}</p>
      </div>

      <div class="finance-stat-grid finance-stat-grid--triple">
        <article v-for="item in heroCards" :key="item.label" class="finance-stat-card">
          <span>{{ item.label }}</span>
          <strong :class="item.tone">{{ item.value }}</strong>
          <small>{{ item.note }}</small>
        </article>
      </div>
    </section>

    <section v-if="isExpenseMode" class="finance-grid-2">
      <article class="page-card finance-panel">
        <header class="finance-panel__header">
          <div>
            <h3>月统计（日趋势）</h3>
            <p>按天查看支出波动，快速定位本月异常高点。</p>
          </div>
        </header>
        <SimpleLineChart :labels="expenseMonthLabels" :series="expenseMonthSeries" :height="290" />
      </article>

      <article class="page-card finance-panel">
        <header class="finance-panel__header">
          <div>
            <h3>年统计（月趋势）</h3>
            <p>把全年每个月支出拉成一条曲线，观察年度消费节奏。</p>
          </div>
        </header>
        <SimpleLineChart :labels="expenseYearLabels" :series="expenseYearSeries" :height="290" />
      </article>
    </section>

    <article v-if="!isExpenseMode" class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>{{ selectedYear }} 年收入趋势</h3>
          <p>V1 聚焦年视角，按月查看收入变化。</p>
        </div>
      </header>
      <SimpleLineChart :labels="incomeYearLabels" :series="incomeYearSeries" :height="320" />
    </article>

    <section class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>{{ isExpenseMode ? "支出前 10 类型排行榜" : "收入前 10 类型排行榜" }}</h3>
          <p>{{ isExpenseMode ? "按支出类型聚合并展示占比" : "按收入类型聚合并展示占比" }}</p>
        </div>
      </header>

      <div class="user-chart-page__ranking">
        <article v-for="(item, index) in currentRanking" :key="item.name" class="chart-rank-row">
          <div class="chart-rank-row__index">{{ index + 1 }}</div>
          <div class="chart-rank-row__body">
            <div class="chart-rank-row__topline">
              <div class="chart-rank-row__title">
                <span class="finance-badge">{{ item.badge }}</span>
                <strong>{{ item.name }}</strong>
              </div>
              <div class="chart-rank-row__value">
                <strong :class="isExpenseMode ? 'finance-tone-expense' : 'finance-tone-income'">
                  {{ formatCurrency(item.value) }}
                </strong>
                <span>{{ item.percent }}%</span>
              </div>
            </div>
            <div class="finance-progress">
              <span class="chart-rank-row__progress" :style="{ width: item.percent + '%' }"></span>
            </div>
          </div>
        </article>
      </div>
    </section>
  </div>
</template>

<script>
import SimpleLineChart from "@/components/SimpleLineChart.vue";
import {
  formatChartCurrency,
  getChartYears,
  getExpenseChartYear,
  getIncomeChartYear
} from "@/utils/userChartMock";

export default {
  name: "UserCharts",
  components: {
    SimpleLineChart
  },
  data() {
    return {
      selectedYear: getChartYears()[0]
    };
  },
  computed: {
    years() {
      return getChartYears();
    },
    isExpenseMode() {
      return this.$route.path.indexOf("/user/charts/income") !== 0;
    },
    expenseData() {
      return getExpenseChartYear(this.selectedYear);
    },
    incomeData() {
      return getIncomeChartYear(this.selectedYear);
    },
    hasPreviousYear() {
      return this.years.indexOf(this.selectedYear) > 0;
    },
    hasNextYear() {
      return this.years.indexOf(this.selectedYear) < this.years.length - 1;
    },
    heroTitle() {
      if (this.isExpenseMode) {
        return "把花钱节奏画成两条曲线，一眼看懂月内波动和全年趋势";
      }

      return "收入不是一条直线，先看全年趋势再做目标拆分";
    },
    heroDescription() {
      if (this.isExpenseMode) {
        return "月趋势用按天折线，年趋势用按月折线，配合前十排行能快速定位重点支出类型。";
      }

      return "收入页按需求保留年视角，配合前十类型排行，先建立稳定增长的观察框架。";
    },
    heroCards() {
      if (this.isExpenseMode) {
        return [
          {
            label: "年总支出",
            value: this.formatCurrency(this.expenseData.summary.yearlyExpense),
            note: "总记录 " + this.expenseData.summary.records + " 笔",
            tone: "finance-tone-expense"
          },
          {
            label: "月均支出",
            value: this.formatCurrency(this.expenseData.summary.monthlyAverage),
            note: "可用于预算校准",
            tone: ""
          },
          {
            label: "高频分类",
            value: this.expenseData.ranking[0].name,
            note: "占比 " + this.expenseData.ranking[0].percent + "%",
            tone: ""
          }
        ];
      }

      return [
        {
          label: "年总收入",
          value: this.formatCurrency(this.incomeData.summary.yearlyIncome),
          note: "总记录 " + this.incomeData.summary.records + " 笔",
          tone: "finance-tone-income"
        },
        {
          label: "月均收入",
          value: this.formatCurrency(this.incomeData.summary.monthlyAverage),
          note: "可作为年度目标基线",
          tone: ""
        },
        {
          label: "主要收入来源",
          value: this.incomeData.ranking[0].name,
          note: "占比 " + this.incomeData.ranking[0].percent + "%",
          tone: ""
        }
      ];
    },
    expenseMonthLabels() {
      return this.expenseData.monthTrend.map(function(_, index) {
        return index + 1 + "日";
      });
    },
    expenseMonthSeries() {
      return [
        {
          name: "当月支出",
          color: "#ef4444",
          values: this.expenseData.monthTrend
        }
      ];
    },
    expenseYearLabels() {
      return this.expenseData.yearTrend.map(function(_, index) {
        return index + 1 + "月";
      });
    },
    expenseYearSeries() {
      return [
        {
          name: "年度支出",
          color: "#f97316",
          values: this.expenseData.yearTrend
        }
      ];
    },
    incomeYearLabels() {
      return this.incomeData.yearTrend.map(function(_, index) {
        return index + 1 + "月";
      });
    },
    incomeYearSeries() {
      return [
        {
          name: "年度收入",
          color: "#22c55e",
          values: this.incomeData.yearTrend
        }
      ];
    },
    currentRanking() {
      return this.isExpenseMode ? this.expenseData.ranking : this.incomeData.ranking;
    }
  },
  methods: {
    formatCurrency: formatChartCurrency,
    prevYear() {
      var index = this.years.indexOf(this.selectedYear);

      if (index > 0) {
        this.selectedYear = this.years[index - 1];
      }
    },
    nextYear() {
      var index = this.years.indexOf(this.selectedYear);

      if (index < this.years.length - 1) {
        this.selectedYear = this.years[index + 1];
      }
    },
    switchMode(mode) {
      var targetPath = mode === "income" ? "/user/charts/income" : "/user/charts/expense";

      if (this.$route.path !== targetPath) {
        this.$router.push(targetPath);
      }
    }
  }
};
</script>

<style scoped>
.user-chart-page__switcher-note {
  margin: 8px 0 0;
  color: var(--text-subtle);
  font-size: 13px;
}

.user-chart-page__ranking {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.chart-rank-row {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 18px;
  border-radius: 18px;
  border: 1px solid var(--border-color);
  background: rgba(255, 255, 255, 0.95);
}

.chart-rank-row__index {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  background: rgba(246, 211, 74, 0.35);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
}

.chart-rank-row__body {
  flex: 1;
}

.chart-rank-row__topline {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  margin-bottom: 8px;
}

.chart-rank-row__title {
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.chart-rank-row__value {
  display: inline-flex;
  align-items: baseline;
  gap: 10px;
}

.chart-rank-row__value span {
  color: var(--text-muted);
  font-size: 13px;
}

.chart-rank-row__progress {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(135deg, var(--brand-color) 0%, var(--brand-hover) 100%);
}

@media (max-width: 920px) {
  .chart-rank-row {
    flex-direction: column;
    align-items: flex-start;
  }

  .chart-rank-row__topline {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
