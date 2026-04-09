<template>
  <div class="finance-page budget-year-page">
    <section class="finance-toolbar">
      <div class="finance-switcher">
        <button :disabled="!hasPreviousYear" @click="prevYear">‹</button>
        <div>
          <strong>{{ selectedYear }} 年预算</strong>
          <p class="budget-year-page__switcher-note">预算是独立视图，不复制月预算的信息结构。</p>
        </div>
        <button :disabled="!hasNextYear" @click="nextYear">›</button>
      </div>

      <div class="finance-toolbar__actions">
        <button class="finance-button finance-button--ghost" @click="openMonthBudget">查看月预算</button>
        <button class="finance-button finance-button--primary" @click="editBudget('年度预算')">调整年度预算</button>
      </div>
    </section>

    <section class="finance-hero finance-hero--soft">
      <div class="finance-hero__eyebrow">
        <span class="finance-hero__eyebrow-mark">年</span>
        <span>{{ selectedYear }} 年预算总览</span>
      </div>

      <div class="finance-hero__headline">
        <h2>年度预算最重要的是看节奏，而不是盯着某一个分类的瞬时高点</h2>
        <p>{{ yearData.overview.note }}</p>
      </div>

      <div class="budget-year-page__hero-progress">
        <div class="budget-year-page__hero-progress-meta">
          <span>预算执行进度</span>
          <strong :class="yearData.overview.remaining < 0 ? 'finance-tone-expense' : 'finance-tone-income'">
            {{ formatPercent(yearData.overview.percentage) }}
          </strong>
        </div>
        <div class="finance-progress">
          <span :class="yearData.overview.remaining < 0 ? 'budget-year-page__hero-bar--danger' : 'budget-year-page__hero-bar--safe'" :style="{ width: progressWidth(yearData.overview.percentage) }"></span>
        </div>
      </div>

      <div class="finance-stat-grid finance-stat-grid--triple">
        <article class="finance-stat-card">
          <span>本年预算</span>
          <strong>{{ formatCurrency(yearData.overview.totalBudget) }}</strong>
          <small>已执行 {{ yearData.overview.months }} 个月</small>
        </article>
        <article class="finance-stat-card">
          <span>本年支出</span>
          <strong class="finance-tone-expense">{{ formatCurrency(yearData.overview.totalExpense) }}</strong>
          <small>预算节奏需要持续盯盘</small>
        </article>
        <article class="finance-stat-card">
          <span>剩余预算</span>
          <strong :class="yearData.overview.remaining < 0 ? 'finance-tone-expense' : 'finance-tone-income'">
            {{ formatCurrency(Math.abs(yearData.overview.remaining)) }}
          </strong>
          <small>{{ yearData.overview.remaining >= 0 ? "仍有可调空间" : "当前已经透支" }}</small>
        </article>
      </div>
    </section>

    <article class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>月度预算执行情况</h3>
          <p>每个月的预算与实际支出并列展示，更容易发现全年失控点。</p>
        </div>
      </header>
      <SimpleGroupedBarChart
        :labels="executionLabels"
        :series="executionSeries"
        :height="320"
      />
    </article>

    <section class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>分类年度预算汇总</h3>
          <p>年预算看的是长期分配结果，所以每一类都应该有明确状态。</p>
        </div>
      </header>

      <div class="budget-year-page__categories">
        <article v-for="item in yearData.categories" :key="item.name" class="year-budget-card">
          <div class="year-budget-card__header">
            <div class="year-budget-card__title">
              <span class="finance-badge">{{ item.badge }}</span>
              <div>
                <strong>{{ item.name }}</strong>
                <p>{{ formatCurrency(item.expense) }} / {{ formatCurrency(item.budget) }}</p>
              </div>
            </div>

            <button class="finance-button finance-button--ghost year-budget-card__action" @click="editBudget(item.name)">
              调整
            </button>
          </div>

          <div class="year-budget-card__meta">
            <div>
              <span>执行比例</span>
              <strong :class="statusTone(item.status)">{{ formatPercent(item.percentage) }}</strong>
            </div>
            <div>
              <span>{{ item.remaining >= 0 ? "剩余" : "超支" }}</span>
              <strong :class="statusTone(item.status)">
                {{ formatCurrency(Math.abs(item.remaining)) }}
              </strong>
            </div>
          </div>

          <div class="finance-progress">
            <span :class="statusBar(item.status)" :style="{ width: progressWidth(item.percentage) }"></span>
          </div>

          <p class="year-budget-card__status" :class="statusTone(item.status)">
            {{ statusLabel(item.status) }}
          </p>
        </article>
      </div>
    </section>
  </div>
</template>

<script>
import { ElMessage } from "element-plus";

import SimpleGroupedBarChart from "@/components/SimpleGroupedBarChart.vue";
import { formatCurrency, formatPercent, getBudgetYears, getYearlyBudget } from "@/utils/userFinanceMock";

export default {
  name: "UserBudgetYear",
  components: {
    SimpleGroupedBarChart
  },
  data() {
    return {
      selectedYear: getBudgetYears()[0]
    };
  },
  computed: {
    years() {
      return getBudgetYears();
    },
    yearData() {
      return getYearlyBudget(this.selectedYear);
    },
    hasPreviousYear() {
      return this.years.indexOf(this.selectedYear) > 0;
    },
    hasNextYear() {
      return this.years.indexOf(this.selectedYear) < this.years.length - 1;
    },
    executionLabels() {
      return this.yearData.monthlyExecution.map(function(item) {
        return item.label;
      });
    },
    executionSeries() {
      return [
        {
          name: "预算",
          color: "#f6d34a",
          values: this.yearData.monthlyExecution.map(function(item) {
            return item.budget;
          })
        },
        {
          name: "实际支出",
          color: "#ef4444",
          values: this.yearData.monthlyExecution.map(function(item) {
            return item.expense;
          })
        }
      ];
    }
  },
  methods: {
    formatCurrency: formatCurrency,
    formatPercent: formatPercent,
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
    openMonthBudget() {
      this.$router.push("/user/budget/month");
    },
    editBudget(name) {
      ElMessage.info(name + " 的年度预算调整将在后续联调批次接入。");
    },
    statusLabel(status) {
      var labels = {
        over: "年度预算已经超支",
        warning: "预算接近上限，需要持续关注",
        safe: "预算执行正常"
      };

      return labels[status] || "预算执行正常";
    },
    statusTone(status) {
      var tones = {
        over: "finance-tone-expense",
        warning: "finance-tone-warning",
        safe: "finance-tone-income"
      };

      return tones[status] || "";
    },
    statusBar(status) {
      var bars = {
        over: "budget-year-page__bar--danger",
        warning: "budget-year-page__bar--warning",
        safe: "budget-year-page__bar--safe"
      };

      return bars[status] || "budget-year-page__bar--safe";
    },
    progressWidth(value) {
      return Math.min(value, 100) + "%";
    }
  }
};
</script>

<style scoped>
.budget-year-page__switcher-note {
  margin: 8px 0 0;
  color: var(--text-subtle);
  font-size: 13px;
}

.budget-year-page__hero-progress {
  margin-top: 24px;
  padding: 20px 22px;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.88);
}

.budget-year-page__hero-progress-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.budget-year-page__hero-bar--danger {
  background: var(--danger-color);
}

.budget-year-page__hero-bar--safe {
  background: linear-gradient(90deg, #22c55e 0%, #f6d34a 100%);
}

.budget-year-page__categories {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.year-budget-card {
  padding: 22px;
  border-radius: 24px;
  border: 1px solid var(--border-color);
  background: rgba(255, 255, 255, 0.96);
}

.year-budget-card__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.year-budget-card__title {
  display: flex;
  gap: 14px;
}

.year-budget-card__title strong {
  display: block;
  font-size: 18px;
}

.year-budget-card__title p {
  margin: 8px 0 0;
  color: var(--text-subtle);
  line-height: 1.6;
}

.year-budget-card__action {
  min-height: 38px;
  padding: 0 14px;
}

.year-budget-card__meta {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
  margin: 20px 0 14px;
}

.year-budget-card__meta span {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
}

.year-budget-card__meta strong {
  display: block;
  margin-top: 8px;
  font-size: 18px;
}

.year-budget-card__status {
  margin: 12px 0 0;
  line-height: 1.6;
  font-size: 13px;
  font-weight: 700;
}

.budget-year-page__bar--danger {
  background: var(--danger-color);
}

.budget-year-page__bar--warning {
  background: var(--warning-color);
}

.budget-year-page__bar--safe {
  background: var(--success-color);
}

@media (max-width: 960px) {
  .budget-year-page__categories,
  .year-budget-card__meta {
    grid-template-columns: minmax(0, 1fr);
  }

  .year-budget-card__header {
    flex-direction: column;
  }
}

</style>
