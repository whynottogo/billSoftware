<template>
  <div class="finance-page bills-list-page">
    <section class="finance-toolbar">
      <div class="finance-switcher">
        <button :disabled="!hasPreviousYear" @click="prevYear">‹</button>
        <div>
          <strong>{{ selectedYear }} 年账单</strong>
          <p class="bills-list-page__switcher-note">{{ yearData.summary.insight }}</p>
        </div>
        <button :disabled="!hasNextYear" @click="nextYear">›</button>
      </div>

      <div class="finance-toolbar__actions">
        <span class="finance-pill">{{ yearData.summary.months }} 个月 · {{ yearData.summary.records }} 笔记录</span>
        <button class="finance-button finance-button--ghost" @click="openLatestMonth">查看最近账单</button>
        <button class="finance-button finance-button--primary" @click="openYearlyBill">查看年账单</button>
      </div>
    </section>

    <section class="finance-hero">
      <div class="finance-hero__eyebrow">
        <span class="finance-hero__eyebrow-mark">账</span>
        <span>{{ selectedYear }} 年度账单总览</span>
      </div>

      <div class="finance-hero__headline">
        <h1>把每个月的收支起伏收进同一张年度时间线里</h1>
        <p>{{ yearData.summary.insight }}</p>
      </div>

      <div class="finance-stat-grid finance-stat-grid--triple">
        <article class="finance-stat-card">
          <span>本年结余</span>
          <strong>{{ formatCurrency(yearData.summary.balance) }}</strong>
          <small>累计记账 {{ yearData.summary.days }} 天</small>
        </article>
        <article class="finance-stat-card">
          <span>年收入</span>
          <strong class="finance-tone-income">{{ formatCurrency(yearData.summary.income) }}</strong>
          <small>收入节奏稳定</small>
        </article>
        <article class="finance-stat-card">
          <span>年支出</span>
          <strong class="finance-tone-expense">{{ formatCurrency(yearData.summary.expense) }}</strong>
          <small>可直接进入月份明细查看结构</small>
        </article>
      </div>
    </section>

    <section class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>月度账单</h3>
          <p>按月份查看收入、支出、结余和记账密度。</p>
        </div>
      </header>

      <div class="bills-list-page__list">
        <article
          v-for="month in yearData.months"
          :key="month.key"
          class="bill-row"
          @click="openMonth(month.key)"
        >
          <div class="bill-row__month">
            <div class="bill-row__month-number">{{ month.month }}</div>
            <div class="bill-row__month-meta">
              <strong>{{ month.label }}</strong>
              <p>{{ month.note }}</p>
            </div>
            <span v-if="month.status === '本月'" class="bill-row__tag">本月</span>
          </div>

          <div class="bill-row__stats">
            <div>
              <span>收入</span>
              <strong class="finance-tone-income">{{ formatCurrency(month.income) }}</strong>
            </div>
            <div>
              <span>支出</span>
              <strong class="finance-tone-expense">{{ formatCurrency(month.expense) }}</strong>
            </div>
            <div>
              <span>结余</span>
              <strong>{{ formatCurrency(month.balance) }}</strong>
            </div>
            <div>
              <span>记账情况</span>
              <strong>{{ month.days }}天 · {{ month.records }}笔</strong>
            </div>
          </div>

          <div class="bill-row__arrow">›</div>
        </article>
      </div>
    </section>

    <section class="finance-grid-3">
      <article
        v-for="item in latestMonth.highlightCards"
        :key="item.label"
        class="page-card bills-list-page__insight"
      >
        <span>{{ item.label }}</span>
        <strong>{{ item.value }}</strong>
        <p>{{ item.hint }}</p>
      </article>
    </section>
  </div>
</template>

<script>
import { formatCurrency, getBillYear, getBillYears } from "@/utils/userFinanceMock";

export default {
  name: "UserBillsMonth",
  data() {
    return {
      selectedYear: getBillYears()[0]
    };
  },
  computed: {
    years() {
      return getBillYears();
    },
    yearData() {
      return getBillYear(this.selectedYear);
    },
    latestMonth() {
      return this.yearData.months[0];
    },
    hasPreviousYear() {
      return this.years.indexOf(this.selectedYear) > 0;
    },
    hasNextYear() {
      return this.years.indexOf(this.selectedYear) < this.years.length - 1;
    }
  },
  methods: {
    formatCurrency: formatCurrency,
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
    openMonth(monthKey) {
      this.$router.push("/user/bills/month/" + monthKey);
    },
    openLatestMonth() {
      this.openMonth(this.latestMonth.key);
    },
    openYearlyBill() {
      this.$router.push("/user/bills/year");
    }
  }
};
</script>

<style scoped>
.bills-list-page__switcher-note {
  margin: 8px 0 0;
  color: var(--text-subtle);
  font-size: 13px;
}

.bills-list-page__list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.bill-row {
  display: flex;
  align-items: center;
  gap: 22px;
  padding: 24px;
  border-radius: 24px;
  border: 1px solid var(--border-color);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.96) 0%, rgba(251, 251, 252, 0.96) 100%);
  transition: 0.2s ease;
  cursor: pointer;
}

.bill-row:hover {
  transform: translateY(-2px);
  border-color: rgba(246, 211, 74, 0.7);
  box-shadow: var(--shadow-sm);
}

.bill-row__month {
  display: flex;
  align-items: center;
  gap: 16px;
  min-width: 260px;
}

.bill-row__month-number {
  width: 68px;
  height: 68px;
  border-radius: 24px;
  background: rgba(246, 211, 74, 0.2);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: 800;
}

.bill-row__month-meta strong {
  display: block;
  font-size: 18px;
}

.bill-row__month-meta p {
  margin: 8px 0 0;
  color: var(--text-subtle);
  line-height: 1.6;
}

.bill-row__tag {
  min-height: 28px;
  padding: 0 12px;
  border-radius: 999px;
  background: rgba(246, 211, 74, 0.3);
  display: inline-flex;
  align-items: center;
  font-size: 12px;
  font-weight: 700;
}

.bill-row__stats {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

.bill-row__stats span {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
}

.bill-row__stats strong {
  display: block;
  margin-top: 8px;
  font-size: 18px;
  line-height: 1.4;
}

.bill-row__arrow {
  font-size: 28px;
  color: var(--text-muted);
}

.bills-list-page__insight {
  padding: 22px;
}

.bills-list-page__insight span {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
}

.bills-list-page__insight strong {
  display: block;
  margin-top: 12px;
  font-size: 26px;
}

.bills-list-page__insight p {
  margin: 10px 0 0;
  color: var(--text-subtle);
  line-height: 1.6;
}

@media (max-width: 960px) {
  .bill-row {
    flex-direction: column;
    align-items: flex-start;
  }

  .bill-row__month,
  .bill-row__stats {
    width: 100%;
  }

  .bill-row__stats {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .bill-row__arrow {
    display: none;
  }
}
</style>
