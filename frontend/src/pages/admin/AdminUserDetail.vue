<template>
  <div class="admin-user-detail-page">
    <section class="admin-user-detail-page__header page-card">
      <div>
        <span class="admin-user-detail-page__eyebrow">管理端 / 用户管理 / 详情</span>
        <h1>{{ profileTitle }}</h1>
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
      <div v-if="loading" class="admin-user-detail-page__loading">正在加载用户账单汇总...</div>

      <div v-else-if="loadError" class="admin-user-detail-page__error">
        <p>{{ loadError }}</p>
        <button class="finance-button finance-button--ghost" @click="fetchDetail">重新加载</button>
      </div>

      <div class="admin-user-detail-page__tabs">
        <button
          :class="['tab-btn', activeTab === 'month' ? 'is-active' : '']"
          :disabled="!detail.monthOptions.length"
          @click="activeTab = 'month'"
        >
          月账单
        </button>
        <button
          :class="['tab-btn', activeTab === 'year' ? 'is-active' : '']"
          :disabled="!detail.yearOptions.length"
          @click="activeTab = 'year'"
        >
          年账单
        </button>
      </div>

      <div v-if="activeTab === 'month'" class="admin-user-detail-page__tab-content">
        <div class="admin-user-detail-page__filter-row">
          <label>
            <span>月份</span>
            <select v-model="selectedMonth" :disabled="!detail.monthOptions.length">
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
            <select v-model="selectedYear" :disabled="!detail.yearOptions.length">
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
import { getAdminUserBillsOverview } from "@/api/adminUserBills";

function createEmptyDetail(userId) {
  return {
    profile: {
      userId: String(userId || ""),
      username: "-",
      nickname: "用户账单详情",
      status: "-",
      registerDate: "-",
      billCount: 0
    },
    monthOptions: [],
    yearOptions: [],
    monthly: {},
    yearly: {}
  };
}

function createEmptyMonthData() {
  return {
    income: 0,
    expense: 0,
    balance: 0,
    records: 0,
    insight: "当前月份暂无账单汇总。"
  };
}

function createEmptyYearData() {
  return {
    income: 0,
    expense: 0,
    balance: 0,
    months: 0,
    insight: "当前年份暂无账单汇总。"
  };
}

function formatCurrency(value) {
  return `¥${Number(value || 0).toLocaleString("zh-CN")}`;
}

function asText(value, fallback) {
  if (value === undefined || value === null || value === "") {
    return fallback;
  }

  return String(value);
}

function parseDate(input) {
  if (!input) {
    return null;
  }

  const normalizedInput =
    typeof input === "string" && /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}(:\d{2})?$/.test(input)
      ? input.replace(" ", "T")
      : input;
  const date = new Date(normalizedInput);
  if (Number.isNaN(date.getTime())) {
    return null;
  }

  return date;
}

function pad(value) {
  return value < 10 ? `0${value}` : String(value);
}

function formatDate(input) {
  const date = parseDate(input);
  if (!date) {
    return "-";
  }

  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`;
}

function buildErrorMessage(error, fallback) {
  if (error && error.response && error.response.data && error.response.data.message) {
    return error.response.data.message;
  }

  return fallback;
}

function normalizeMetric(source, fallback) {
  const target = Object.assign({}, fallback);
  const payload = source || {};

  Object.keys(target).forEach(function(key) {
    if (payload[key] !== undefined && payload[key] !== null && payload[key] !== "") {
      target[key] = payload[key];
    }
  });

  return target;
}

function normalizeDetail(payload, userId) {
  const safePayload = payload || {};
  const safeProfile = safePayload.profile || {};
  const detail = createEmptyDetail(userId);

  detail.profile = {
    userId: asText(safeProfile.userId, detail.profile.userId),
    username: asText(safeProfile.username, "-"),
    nickname: asText(safeProfile.nickname, "用户账单详情"),
    status: Number(safeProfile.status) === 1 ? "启用" : safeProfile.status === 0 ? "禁用" : asText(safeProfile.status, "-"),
    registerDate: formatDate(safeProfile.registerDate),
    billCount: Number(safeProfile.billCount || 0)
  };

  detail.monthOptions = Array.isArray(safePayload.monthOptions) ? safePayload.monthOptions.slice() : [];
  detail.yearOptions = Array.isArray(safePayload.yearOptions) ? safePayload.yearOptions.slice() : [];

  detail.monthly = {};
  detail.monthOptions.forEach(function(month) {
    detail.monthly[month] = normalizeMetric(
      safePayload.monthly && safePayload.monthly[month],
      createEmptyMonthData()
    );
  });

  detail.yearly = {};
  detail.yearOptions.forEach(function(year) {
    detail.yearly[year] = normalizeMetric(
      safePayload.yearly && safePayload.yearly[year],
      createEmptyYearData()
    );
  });

  return detail;
}

export default {
  name: "AdminUserDetail",
  data() {
    return {
      detail: createEmptyDetail(this.$route.params.userId),
      activeTab: "month",
      selectedMonth: "",
      selectedYear: "",
      loading: false,
      loadError: ""
    };
  },
  computed: {
    profileTitle() {
      const nickname = this.detail.profile.nickname || "用户账单详情";
      const username = this.detail.profile.username || "-";
      return `${nickname}（${username}）`;
    },
    monthlyData() {
      return this.detail.monthly[this.selectedMonth] || createEmptyMonthData();
    },
    yearlyData() {
      return this.detail.yearly[this.selectedYear] || createEmptyYearData();
    }
  },
  created() {
    this.fetchDetail();
  },
  watch: {
    "$route.params.userId": function() {
      this.fetchDetail();
    },
    "detail.monthOptions": function(options) {
      if (!options.length) {
        this.selectedMonth = "";
        return;
      }

      if (options.indexOf(this.selectedMonth) === -1) {
        this.selectedMonth = options[0];
      }
    },
    "detail.yearOptions": function(options) {
      if (!options.length) {
        this.selectedYear = "";
        return;
      }

      if (options.indexOf(this.selectedYear) === -1) {
        this.selectedYear = options[0];
      }
    }
  },
  methods: {
    formatCurrency: formatCurrency,
    fetchDetail() {
      const userId = this.$route.params.userId;

      this.loading = true;
      this.loadError = "";
      this.detail = createEmptyDetail(userId);
      this.selectedMonth = "";
      this.selectedYear = "";

      return getAdminUserBillsOverview(userId)
        .then(
          function(result) {
            this.detail = normalizeDetail(result && result.data ? result.data : result, userId);
          }.bind(this)
        )
        .catch(
          function(error) {
            const message = buildErrorMessage(error, "用户账单详情加载失败，请稍后重试");

            this.loadError = message;
            this.$message.error(message);
          }.bind(this)
        )
        .finally(
          function() {
            this.loading = false;
          }.bind(this)
        );
    },
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

.admin-user-detail-page__loading,
.admin-user-detail-page__error {
  margin-bottom: 20px;
  padding: 16px 18px;
  border-radius: 16px;
  border: 1px solid var(--border-color);
  background: rgba(245, 246, 248, 0.9);
}

.admin-user-detail-page__loading {
  color: var(--text-subtle);
}

.admin-user-detail-page__error {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  color: var(--danger-color);
}

.admin-user-detail-page__error p {
  margin: 0;
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

.tab-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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
