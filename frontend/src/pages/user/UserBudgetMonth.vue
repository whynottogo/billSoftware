<template>
  <div class="finance-page budget-month-page">
    <section class="finance-toolbar">
      <div class="finance-toolbar__meta">
        <h1 class="page-title">{{ budgetData.label }}</h1>
        <p class="page-description">{{ budgetData.notice }}</p>
      </div>

      <div class="finance-toolbar__actions">
        <button class="finance-button finance-button--ghost" @click="openYearBudget">查看年预算</button>
        <button class="finance-button finance-button--primary" @click="editBudget('本月总预算')">设置预算</button>
      </div>
    </section>

    <section class="finance-hero finance-hero--soft">
      <div class="finance-hero__eyebrow">
        <span class="finance-hero__eyebrow-mark">预</span>
        <span>本月预算总览</span>
      </div>

      <div class="finance-hero__headline">
        <h2>预算页不是另一张账单页，它应该直接告诉你哪里还能花、哪里该收手</h2>
        <p>本月支出已进入后半程，但还有足够余量可以精细调整。</p>
      </div>

      <div class="finance-stat-grid finance-stat-grid--triple">
        <article class="finance-stat-card">
          <span>本月预算</span>
          <strong>{{ formatCurrency(budgetData.overview.totalBudget) }}</strong>
          <small>围绕支出类型分配额度</small>
        </article>
        <article class="finance-stat-card">
          <span>本月支出</span>
          <strong class="finance-tone-expense">{{ formatCurrency(budgetData.overview.totalExpense) }}</strong>
          <small>当前月度已发生支出</small>
        </article>
        <article class="finance-stat-card">
          <span>剩余预算</span>
          <strong class="finance-tone-income">{{ formatCurrency(budgetData.overview.remaining) }}</strong>
          <small>仍可主动安排的空间</small>
        </article>
      </div>

      <div class="budget-month-page__hero-progress">
        <div class="budget-month-page__hero-progress-meta">
          <span>预算使用进度</span>
          <strong>{{ formatPercent(budgetData.overview.percentage) }}</strong>
        </div>
        <div class="finance-progress">
          <span class="budget-month-page__hero-bar" :style="{ width: budgetData.overview.percentage + '%' }"></span>
        </div>
      </div>
    </section>

    <section class="finance-grid-3">
      <article v-for="item in budgetData.highlights" :key="item.label" class="page-card budget-month-page__highlight">
        <span>{{ item.label }}</span>
        <strong>{{ item.value }}</strong>
        <p>{{ item.hint }}</p>
      </article>
    </section>

    <section class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>分类预算</h3>
          <p>每一类预算都单独看状态，这样才能知道要不要继续压缩。</p>
        </div>
      </header>

      <div class="budget-month-page__categories">
        <article v-for="item in budgetData.categories" :key="item.name" class="budget-card">
          <div class="budget-card__header">
            <div class="budget-card__title">
              <span class="finance-badge">{{ item.badge }}</span>
              <div>
                <strong>{{ item.name }}</strong>
                <p>{{ item.note }}</p>
              </div>
            </div>

            <button class="finance-button finance-button--ghost budget-card__action" @click="editBudget(item.name)">
              调整
            </button>
          </div>

          <div class="budget-card__numbers">
            <div>
              <span>预算</span>
              <strong>{{ formatCurrency(item.budget) }}</strong>
            </div>
            <div>
              <span>已用</span>
              <strong class="finance-tone-expense">{{ formatCurrency(item.expense) }}</strong>
            </div>
            <div>
              <span>{{ item.remaining >= 0 ? "剩余" : "超支" }}</span>
              <strong :class="statusTone(item.status)">
                {{ item.remaining >= 0 ? formatCurrency(item.remaining) : formatCurrency(Math.abs(item.remaining)) }}
              </strong>
            </div>
          </div>

          <div class="budget-card__progress">
            <div class="budget-card__progress-meta">
              <span>{{ statusLabel(item.status) }}</span>
              <strong :class="statusTone(item.status)">{{ formatPercent(item.percentage) }}</strong>
            </div>
            <div class="finance-progress">
              <span :class="statusBar(item.status)" :style="{ width: progressWidth(item.percentage) }"></span>
            </div>
          </div>
        </article>
      </div>
    </section>
  </div>
</template>

<script>
import { ElMessage } from "element-plus";

import { formatCurrency, formatPercent, getMonthlyBudget } from "@/utils/userFinanceMock";

export default {
  name: "UserBudgetMonth",
  computed: {
    budgetData() {
      return getMonthlyBudget();
    }
  },
  methods: {
    formatCurrency: formatCurrency,
    formatPercent: formatPercent,
    openYearBudget() {
      this.$router.push("/user/budget/year");
    },
    editBudget(name) {
      ElMessage.info(name + " 的预算设置将在后续联调批次接入。");
    },
    statusLabel(status) {
      var labels = {
        over: "已超支",
        warning: "接近上限",
        safe: "仍有余量"
      };

      return labels[status] || "正常";
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
        over: "budget-card__bar--danger",
        warning: "budget-card__bar--warning",
        safe: "budget-card__bar--safe"
      };

      return bars[status] || "budget-card__bar--safe";
    },
    progressWidth(value) {
      return Math.min(value, 100) + "%";
    }
  }
};
</script>

<style scoped>
.budget-month-page__hero-progress {
  margin-top: 24px;
  padding: 20px 22px;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.88);
}

.budget-month-page__hero-progress-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.budget-month-page__hero-progress-meta span {
  color: var(--text-subtle);
}

.budget-month-page__hero-progress-meta strong {
  font-size: 22px;
}

.budget-month-page__hero-bar {
  background: linear-gradient(90deg, #22c55e 0%, #f6d34a 100%);
}

.budget-month-page__highlight {
  padding: 22px;
}

.budget-month-page__highlight span {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
}

.budget-month-page__highlight strong {
  display: block;
  margin-top: 12px;
  font-size: 26px;
}

.budget-month-page__highlight p {
  margin: 10px 0 0;
  color: var(--text-subtle);
  line-height: 1.6;
}

.budget-month-page__categories {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.budget-card {
  padding: 22px;
  border-radius: 24px;
  border: 1px solid var(--border-color);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98) 0%, rgba(251, 251, 252, 0.96) 100%);
}

.budget-card__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.budget-card__title {
  display: flex;
  gap: 14px;
}

.budget-card__title strong {
  display: block;
  font-size: 18px;
}

.budget-card__title p {
  margin: 8px 0 0;
  color: var(--text-subtle);
  line-height: 1.6;
}

.budget-card__action {
  min-height: 38px;
  padding: 0 14px;
}

.budget-card__numbers {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
  margin-top: 20px;
}

.budget-card__numbers span {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
}

.budget-card__numbers strong {
  display: block;
  margin-top: 8px;
  font-size: 18px;
}

.budget-card__progress {
  margin-top: 20px;
}

.budget-card__progress-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.budget-card__bar--danger {
  background: var(--danger-color);
}

.budget-card__bar--warning {
  background: var(--warning-color);
}

.budget-card__bar--safe {
  background: var(--success-color);
}

@media (max-width: 960px) {
  .budget-month-page__categories,
  .budget-card__numbers {
    grid-template-columns: minmax(0, 1fr);
  }

  .budget-card__header {
    flex-direction: column;
  }
}
</style>
