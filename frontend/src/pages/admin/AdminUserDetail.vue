<template>
  <div class="admin-user-detail-page">
    <section class="admin-user-detail-page__header page-card">
      <div>
        <span class="admin-user-detail-page__eyebrow">管理端 / 用户管理 / 详情</span>
        <h1>{{ detail.profile.nickname }}（{{ detail.profile.username }}）</h1>
        <p>当前为只读视图，按月/按年查看单用户收入、支出与结余汇总。</p>
      </div>
      <div class="admin-user-detail-page__actions">
        <button class="finance-button finance-button--ghost" @click="goBack">返回用户列表</button>
      </div>
    </section>

    <section class="admin-user-detail-page__summary">
      <article class="admin-summary-card">
        <span>状态</span>
        <strong :class="detail.profile.status === '启用' ? 'is-enabled' : 'is-disabled'">{{ detail.profile.status }}</strong>
      </article>
      <article class="admin-summary-card">
        <span>注册时间</span>
        <strong>{{ detail.profile.registerDate }}</strong>
      </article>
      <article class="admin-summary-card">
        <span>账单记录</span>
        <strong>{{ detail.profile.billCount }} 笔</strong>
      </article>
    </section>

    <section class="admin-user-detail-page__panel page-card">
      <div class="admin-user-detail-page__tabs">
        <button
          :class="['tab-btn', activeTab === 'month' ? 'is-active' : '']"
          @click="activeTab = 'month'"
        >
          月账单
        </button>
        <button
          :class="['tab-btn', activeTab === 'year' ? 'is-active' : '']"
          @click="activeTab = 'year'"
        >
          年账单
        </button>
      </div>

      <div v-if="activeTab === 'month'" class="admin-user-detail-page__tab-content">
        <div class="admin-user-detail-page__filter-row">
          <label>
            <span>月份</span>
            <select v-model="selectedMonth">
              <option v-for="item in detail.monthOptions" :key="item" :value="item">
                {{ item }}
              </option>
            </select>
          </label>
        </div>

        <div class="admin-user-detail-page__stats-grid">
          <article class="admin-detail-stat">
            <span>收入</span>
            <strong class="is-income">{{ formatCurrency(monthlyData.income) }}</strong>
          </article>
          <article class="admin-detail-stat">
            <span>支出</span>
            <strong class="is-expense">{{ formatCurrency(monthlyData.expense) }}</strong>
          </article>
          <article class="admin-detail-stat">
            <span>结余</span>
            <strong>{{ formatCurrency(monthlyData.balance) }}</strong>
          </article>
          <article class="admin-detail-stat">
            <span>记录数</span>
            <strong>{{ monthlyData.records }} 笔</strong>
          </article>
        </div>

        <p class="admin-user-detail-page__insight">{{ monthlyData.insight }}</p>
      </div>

      <div v-else class="admin-user-detail-page__tab-content">
        <div class="admin-user-detail-page__filter-row">
          <label>
            <span>年份</span>
            <select v-model="selectedYear">
              <option v-for="item in detail.yearOptions" :key="item" :value="item">
                {{ item }}
              </option>
            </select>
          </label>
        </div>

        <div class="admin-user-detail-page__stats-grid">
          <article class="admin-detail-stat">
            <span>收入</span>
            <strong class="is-income">{{ formatCurrency(yearlyData.income) }}</strong>
          </article>
          <article class="admin-detail-stat">
            <span>支出</span>
            <strong class="is-expense">{{ formatCurrency(yearlyData.expense) }}</strong>
          </article>
          <article class="admin-detail-stat">
            <span>结余</span>
            <strong>{{ formatCurrency(yearlyData.balance) }}</strong>
          </article>
          <article class="admin-detail-stat">
            <span>覆盖月份</span>
            <strong>{{ yearlyData.months }} 个月</strong>
          </article>
        </div>

        <p class="admin-user-detail-page__insight">{{ yearlyData.insight }}</p>
      </div>
    </section>
  </div>
</template>

<script>
import { formatAdminCurrency, getAdminUserDetailMock } from "@/utils/adminUserMock";

export default {
  name: "AdminUserDetail",
  data() {
    const detail = getAdminUserDetailMock(this.$route.params.userId);

    return {
      detail: detail,
      activeTab: "month",
      selectedMonth: detail.monthOptions[0],
      selectedYear: detail.yearOptions[0]
    };
  },
  computed: {
    monthlyData() {
      return this.detail.monthly[this.selectedMonth];
    },
    yearlyData() {
      return this.detail.yearly[this.selectedYear];
    }
  },
  methods: {
    formatCurrency: formatAdminCurrency,
    goBack() {
      this.$router.push("/admin/users");
    }
  }
};
</script>

<style scoped>
.admin-user-detail-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.admin-user-detail-page__header {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 24px;
}

.admin-user-detail-page__eyebrow {
  display: inline-flex;
  align-items: center;
  min-height: 34px;
  padding: 0 14px;
  border-radius: 999px;
  background: rgba(23, 23, 23, 0.08);
  color: var(--text-main);
  font-size: 13px;
  font-weight: 700;
}

.admin-user-detail-page__header h1 {
  margin: 16px 0 10px;
  font-size: 30px;
}

.admin-user-detail-page__header p {
  margin: 0;
  color: var(--text-subtle);
  line-height: 1.7;
}

.admin-user-detail-page__summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
}

.admin-summary-card {
  padding: 20px;
  border-radius: 20px;
  border: 1px solid var(--border-color);
  background: rgba(255, 255, 255, 0.95);
}

.admin-summary-card span {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
}

.admin-summary-card strong {
  display: block;
  margin-top: 8px;
  font-size: 24px;
}

.admin-summary-card .is-enabled {
  color: var(--success-color);
}

.admin-summary-card .is-disabled {
  color: var(--danger-color);
}

.admin-user-detail-page__panel {
  padding: 24px;
}

.admin-user-detail-page__tabs {
  display: inline-flex;
  padding: 4px;
  border-radius: 14px;
  background: rgba(23, 23, 23, 0.06);
}

.tab-btn {
  min-height: 38px;
  border: none;
  background: transparent;
  border-radius: 10px;
  padding: 0 14px;
  color: var(--text-subtle);
  font-weight: 600;
}

.tab-btn.is-active {
  background: #171717;
  color: #ffffff;
}

.admin-user-detail-page__tab-content {
  margin-top: 20px;
}

.admin-user-detail-page__filter-row label {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 180px;
  font-size: 13px;
  color: var(--text-subtle);
  font-weight: 600;
}

.admin-user-detail-page__filter-row select {
  height: 40px;
  border-radius: 12px;
  border: 1px solid var(--border-color);
  padding: 0 12px;
}

.admin-user-detail-page__stats-grid {
  margin-top: 16px;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.admin-detail-stat {
  padding: 18px;
  border-radius: 16px;
  border: 1px solid var(--border-color);
  background: rgba(245, 246, 248, 0.9);
}

.admin-detail-stat span {
  display: block;
  font-size: 12px;
  color: var(--text-muted);
}

.admin-detail-stat strong {
  display: block;
  margin-top: 8px;
  font-size: 22px;
}

.admin-detail-stat .is-income {
  color: var(--success-color);
}

.admin-detail-stat .is-expense {
  color: var(--danger-color);
}

.admin-user-detail-page__insight {
  margin: 18px 0 0;
  padding: 14px 16px;
  border-radius: 14px;
  background: rgba(246, 211, 74, 0.14);
  color: var(--text-main);
}

@media (max-width: 980px) {
  .admin-user-detail-page__header {
    flex-direction: column;
  }

  .admin-user-detail-page__summary,
  .admin-user-detail-page__stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
