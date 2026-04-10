<template>
  <div class="finance-page bills-year-page">
    <section class="finance-toolbar">
      <div class="finance-switcher">
        <button :disabled="!hasPreviousYear" @click="prevYear">‹</button>
        <div>
          <strong>{{ selectedYear }} 年度账单</strong>
          <p class="bills-year-page__switcher-note">可直接切换年份查看全年节奏。</p>
        </div>
        <button :disabled="!hasNextYear" @click="nextYear">›</button>
      </div>

      <div class="finance-toolbar__actions">
        <span class="finance-pill">{{ yearData.summary.months }} 个月 · {{ yearData.summary.days }} 天</span>
        <button class="finance-button finance-button--ghost" @click="goMonthList">查看月账单</button>
      </div>
    </section>

    <div v-if="errorMessage" class="finance-inline-notice">
      <span>{{ errorMessage }}</span>
      <button type="button" @click="retryCurrentYear">重试</button>
    </div>

    <section class="finance-hero finance-hero--soft">
      <div class="finance-hero__eyebrow">
        <span class="finance-hero__eyebrow-mark">年</span>
        <span>{{ selectedYear }} 年账单总览</span>
      </div>

      <div class="finance-hero__headline">
        <h1>年度视角不是汇总页，而是你消费节奏的总控制台</h1>
        <p>{{ yearData.summary.insight }}</p>
      </div>

      <div class="finance-stat-grid finance-stat-grid--triple">
        <article class="finance-stat-card">
          <span>总结余</span>
          <strong>{{ formatCurrency(yearData.summary.balance) }}</strong>
          <small>共 {{ yearData.summary.records }} 笔账单</small>
        </article>
        <article class="finance-stat-card">
          <span>总收入</span>
          <strong class="finance-tone-income">{{ formatCurrency(yearData.summary.income) }}</strong>
          <small>收入曲线整体平稳</small>
        </article>
        <article class="finance-stat-card">
          <span>总支出</span>
          <strong class="finance-tone-expense">{{ formatCurrency(yearData.summary.expense) }}</strong>
          <small>年度节奏一眼可见</small>
        </article>
      </div>
    </section>

    <article class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>{{ selectedYear }} 年月度趋势</h3>
          <p>收入、支出和结余放到一张图里，最适合判断全年高低点。</p>
        </div>
      </header>
      <div v-if="isLoading && !yearData.months.length" class="bills-year-page__placeholder">
        正在加载年度趋势...
      </div>
      <div v-else-if="!yearData.months.length" class="bills-year-page__placeholder">
        当前年份暂无月度数据，记账后会自动生成趋势。
      </div>
      <SimpleLineChart
        :labels="lineLabels"
        :series="lineSeries"
        :height="320"
      />
    </article>

    <section class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>历年账单</h3>
          <p>年度之间的差别不是数字多少，而是哪一年开始形成更好的消费结构。</p>
        </div>
      </header>

      <div class="bills-year-page__history">
        <article v-if="isLoading && !history.length" class="year-card year-card--placeholder">
          <div class="year-card__year">
            <strong>--</strong>
            <span>正在加载历年账单</span>
          </div>
        </article>

        <article v-else-if="!history.length" class="year-card year-card--placeholder">
          <div class="year-card__year">
            <strong>--</strong>
            <span>暂无历年账单</span>
          </div>
        </article>

        <article
          v-for="item in history"
          :key="item.year"
          :class="['year-card', item.year === selectedYear ? 'is-active' : '']"
          @click="selectYear(item.year)"
        >
          <div class="year-card__year">
            <strong>{{ item.year }}</strong>
            <span>{{ item.year === selectedYear ? "当前查看" : "点击切换" }}</span>
          </div>

          <div class="year-card__stats">
            <div>
              <span>总结余</span>
              <strong>{{ formatCurrency(item.summary.balance) }}</strong>
            </div>
            <div>
              <span>总收入</span>
              <strong class="finance-tone-income">{{ formatCurrency(item.summary.income) }}</strong>
            </div>
            <div>
              <span>总支出</span>
              <strong class="finance-tone-expense">{{ formatCurrency(item.summary.expense) }}</strong>
            </div>
            <div>
              <span>记账情况</span>
              <strong>{{ item.summary.months }}个月 · {{ item.summary.days }}天</strong>
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
  buildBillSummaryError,
  createEmptyBillYearData,
  formatBillCurrency,
  getUserBillYear,
  getUserBillYears,
  normalizeBillYearPayload,
  normalizeBillYearsPayload
} from "@/api/userBillSummary";

export default {
  name: "UserBillsYear",
  components: {
    SimpleLineChart
  },
  data() {
    return {
      selectedYear: new Date().getFullYear(),
      years: [],
      history: [],
      yearData: createEmptyBillYearData(),
      isLoading: false,
      errorMessage: "",
      activeYearRequestKey: 0
    };
  },
  created() {
    this.bootstrap();
  },
  computed: {
    hasPreviousYear() {
      return this.years.indexOf(this.selectedYear) < this.years.length - 1;
    },
    hasNextYear() {
      return this.years.indexOf(this.selectedYear) > 0;
    },
    lineLabels() {
      return this.yearData.months.slice().sort(function(a, b) {
        return a.month - b.month;
      }).map(function(item) {
        return item.month + "月";
      });
    },
    lineSeries() {
      var months = this.yearData.months.slice().sort(function(a, b) {
        return a.month - b.month;
      });

      return [
        {
          name: "收入",
          color: "#22c55e",
          values: months.map(function(item) {
            return item.income;
          })
        },
        {
          name: "支出",
          color: "#ef4444",
          values: months.map(function(item) {
            return item.expense;
          })
        },
        {
          name: "结余",
          color: "#f6d34a",
          values: months.map(function(item) {
            return item.balance;
          })
        }
      ];
    }
  },
  methods: {
    formatCurrency: formatBillCurrency,
    bootstrap() {
      this.isLoading = true;
      this.errorMessage = "";

      return getUserBillYears()
        .then(
          function(result) {
            var normalized = normalizeBillYearsPayload(result);
            var fallbackYear = normalized.years[0] || new Date().getFullYear();

            this.years = normalized.years;
            this.history = normalized.history;
            this.selectedYear = fallbackYear;

            if (!normalized.years.length) {
              this.yearData = createEmptyBillYearData(fallbackYear);
              return null;
            }

            return this.loadYearData(fallbackYear, true);
          }.bind(this)
        )
        .catch(
          function(error) {
            this.years = [];
            this.history = [];
            this.yearData = createEmptyBillYearData(this.selectedYear);
            this.errorMessage = buildBillSummaryError(error, "年度账单列表加载失败，请稍后重试。");
            this.$message.error(this.errorMessage);
          }.bind(this)
        )
        .finally(
          function() {
            this.isLoading = false;
          }.bind(this)
        );
    },
    loadYearData(year, silent) {
      var requestKey = Date.now();

      this.activeYearRequestKey = requestKey;
      this.isLoading = true;
      this.errorMessage = "";

      return getUserBillYear(year)
        .then(
          function(result) {
            if (this.activeYearRequestKey !== requestKey) {
              return;
            }

            this.yearData = normalizeBillYearPayload(result, year);
          }.bind(this)
        )
        .catch(
          function(error) {
            if (this.activeYearRequestKey !== requestKey) {
              return;
            }

            this.yearData = createEmptyBillYearData(year);
            this.errorMessage = buildBillSummaryError(error, year + " 年年度账单加载失败，请稍后重试。");

            if (!silent) {
              this.$message.error(this.errorMessage);
            }
          }.bind(this)
        )
        .finally(
          function() {
            if (this.activeYearRequestKey === requestKey) {
              this.isLoading = false;
            }
          }.bind(this)
        );
    },
    prevYear() {
      var index = this.years.indexOf(this.selectedYear);

      if (index < this.years.length - 1) {
        this.selectedYear = this.years[index + 1];
        this.loadYearData(this.selectedYear);
      }
    },
    nextYear() {
      var index = this.years.indexOf(this.selectedYear);

      if (index > 0) {
        this.selectedYear = this.years[index - 1];
        this.loadYearData(this.selectedYear);
      }
    },
    selectYear(year) {
      this.selectedYear = year;
      this.loadYearData(year);
    },
    goMonthList() {
      this.$router.push("/user/bills/month");
    },
    retryCurrentYear() {
      if (!this.years.length) {
        this.bootstrap();
        return;
      }

      this.loadYearData(this.selectedYear);
    }
  }
};
</script>

<style scoped>
.finance-inline-notice {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 18px;
  border: 1px solid rgba(239, 68, 68, 0.18);
  background: rgba(255, 244, 244, 0.92);
  color: #9f2f2f;
}

.finance-inline-notice button {
  border: none;
  background: transparent;
  color: inherit;
  font-weight: 700;
  cursor: pointer;
}

.bills-year-page__switcher-note {
  margin: 8px 0 0;
  color: var(--text-subtle);
  font-size: 13px;
}

.bills-year-page__placeholder {
  margin: 0 0 18px;
  padding: 18px 20px;
  border-radius: 20px;
  background: rgba(255, 248, 228, 0.68);
  color: var(--text-subtle);
}

.bills-year-page__history {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.year-card--placeholder {
  cursor: default;
}

.year-card {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 24px;
  border-radius: 24px;
  border: 1px solid var(--border-color);
  background: rgba(255, 255, 255, 0.96);
  transition: 0.2s ease;
  cursor: pointer;
}

.year-card:hover,
.year-card.is-active {
  border-color: rgba(246, 211, 74, 0.72);
  box-shadow: var(--shadow-sm);
}

.year-card__year {
  min-width: 140px;
}

.year-card__year strong,
.year-card__year span {
  display: block;
}

.year-card__year strong {
  font-size: 32px;
}

.year-card__year span {
  margin-top: 8px;
  color: var(--text-muted);
  font-size: 12px;
}

.year-card__stats {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

.year-card__stats span {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
}

.year-card__stats strong {
  display: block;
  margin-top: 8px;
  font-size: 18px;
  line-height: 1.5;
}

@media (max-width: 960px) {
  .year-card {
    flex-direction: column;
    align-items: flex-start;
  }

  .year-card__stats {
    width: 100%;
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

</style>
