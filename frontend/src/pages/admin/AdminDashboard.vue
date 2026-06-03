<template>
  <div class="finance-page admin-dashboard-page">
    <section class="finance-hero finance-hero--soft admin-dashboard-page__hero">
      <div class="finance-hero__eyebrow">
        <span class="finance-hero__eyebrow-mark">概</span>
        <span>管理后台 / 数据概览</span>
      </div>

      <div class="finance-hero__headline">
        <h1>一屏查看运营走势和待处理事项</h1>
        <p>核心指标来自管理端实时统计，便于快速判断用户、家庭和账单变化。</p>
      </div>

      <div v-if="isLoading" class="finance-stat-grid">
        <article v-for="index in 4" :key="'stat-loading-' + index" class="finance-stat-card admin-dashboard-page__stat-loading">
          <span></span>
          <strong></strong>
          <small></small>
        </article>
      </div>

      <div v-else class="finance-stat-grid">
        <article v-for="item in dashboard.stats" :key="item.key" class="finance-stat-card">
          <span>{{ item.label }}</span>
          <strong>{{ item.value }}</strong>
          <small :class="trendClass(item.tone)">{{ item.trend }}</small>
        </article>
      </div>
    </section>

    <section v-if="isLoading" class="finance-grid-2">
      <article v-for="index in 2" :key="'panel-loading-' + index" class="page-card finance-panel admin-dashboard-page__panel-loading">
        <span></span>
        <strong></strong>
        <i></i>
      </article>
    </section>

    <section v-else-if="errorMessage" class="page-card finance-panel admin-dashboard-page__error">
      <h3>数据加载失败</h3>
      <p>{{ errorMessage }}</p>
      <button class="finance-button finance-button--primary" @click="fetchDashboard">重试</button>
    </section>

    <template v-else>
      <section class="finance-grid-2">
        <article class="page-card finance-panel">
          <header class="finance-panel__header">
            <div>
              <h3>用户增长趋势</h3>
              <p>最近 6 个月新增用户曲线。</p>
            </div>
          </header>
          <SimpleLineChart
            v-if="hasUserGrowthData"
            :labels="dashboard.userGrowth.labels"
            :series="dashboard.userGrowth.series"
            :height="300"
            value-unit="count"
          />
          <div v-else class="admin-dashboard-page__empty">暂无用户增长数据</div>
        </article>

        <article class="page-card finance-panel">
          <header class="finance-panel__header">
            <div>
              <h3>账单数量统计</h3>
              <p>最近 7 天账单提交情况。</p>
            </div>
          </header>
          <SimpleGroupedBarChart
            v-if="hasBillTrendData"
            :labels="dashboard.billTrend.labels"
            :series="dashboard.billTrend.series"
            :height="300"
            value-unit="count"
          />
          <div v-else class="admin-dashboard-page__empty">暂无账单数量数据</div>
        </article>
      </section>

      <section class="page-card finance-panel">
        <header class="finance-panel__header">
          <div>
            <h3>最近注册用户</h3>
            <p>最新 5 位注册用户，包含审批状态参考。</p>
          </div>
          <button class="finance-button finance-button--ghost" @click="openUsers">进入用户管理</button>
        </header>

        <div v-if="dashboard.recentUsers.length" class="admin-dashboard-page__user-list">
          <article v-for="user in dashboard.recentUsers" :key="user.id" class="admin-dashboard-page__user-row">
            <div class="admin-dashboard-page__user-main">
              <span class="admin-dashboard-page__avatar">{{ user.name.slice(0, 1) }}</span>
              <div>
                <strong>{{ user.name }}</strong>
                <p>{{ user.email || user.username }}</p>
              </div>
            </div>
            <div class="admin-dashboard-page__user-meta">
              <span>{{ user.registerDate }}</span>
              <em :class="statusClass(user.status)">{{ statusLabel(user.status) }}</em>
            </div>
          </article>
        </div>
        <div v-else class="admin-dashboard-page__empty">暂无最近注册用户</div>
      </section>
    </template>
  </div>
</template>

<script>
import SimpleGroupedBarChart from "@/components/SimpleGroupedBarChart.vue";
import SimpleLineChart from "@/components/SimpleLineChart.vue";
import { getAdminDashboardOverview } from "@/api/adminDashboard";

function emptyDashboard() {
  return {
    stats: [
      { key: "totalUsers", label: "总用户数", value: "0", trend: "实时统计", tone: "normal" },
      { key: "pendingUsers", label: "待审批用户", value: "0", trend: "已清零", tone: "normal" },
      { key: "todayBills", label: "今日账单数", value: "0", trend: "今日", tone: "normal" },
      { key: "families", label: "家庭数量", value: "0", trend: "实时统计", tone: "normal" }
    ],
    userGrowth: {
      labels: [],
      series: [{ name: "新增用户", color: "#171717", values: [] }]
    },
    billTrend: {
      labels: [],
      series: [{ name: "账单数", color: "#f6d34a", values: [] }]
    },
    recentUsers: []
  };
}

function asNumber(value) {
  var number = Number(value || 0);
  return Number.isFinite(number) ? number : 0;
}

function formatCount(value) {
  return asNumber(value).toLocaleString("zh-CN");
}

function buildTrend(items, name, color) {
  var rows = Array.isArray(items) ? items : [];

  return {
    labels: rows.map(function(item) {
      return item && item.label ? item.label : "-";
    }),
    series: [
      {
        name: name,
        color: color,
        values: rows.map(function(item) {
          return asNumber(item && item.value);
        })
      }
    ]
  };
}

function formatDateTime(value) {
  if (!value) {
    return "-";
  }

  var date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return String(value);
  }

  var year = date.getFullYear();
  var month = String(date.getMonth() + 1).padStart(2, "0");
  var day = String(date.getDate()).padStart(2, "0");
  var hour = String(date.getHours()).padStart(2, "0");
  var minute = String(date.getMinutes()).padStart(2, "0");

  return year + "-" + month + "-" + day + " " + hour + ":" + minute;
}

function normalizeUserStatus(status, approvalStatus) {
  if (approvalStatus === "pending") {
    return "pending";
  }

  if (approvalStatus === "rejected") {
    return "inactive";
  }

  var numericStatus = Number(status);

  if (numericStatus === 0) {
    return "pending";
  }

  if (numericStatus === 1) {
    return "active";
  }

  return "inactive";
}

export default {
  name: "AdminDashboard",
  components: {
    SimpleLineChart: SimpleLineChart,
    SimpleGroupedBarChart: SimpleGroupedBarChart
  },
  data() {
    return {
      dashboard: emptyDashboard(),
      isLoading: false,
      errorMessage: ""
    };
  },
  computed: {
    hasUserGrowthData() {
      return this.hasTrendData(this.dashboard.userGrowth);
    },
    hasBillTrendData() {
      return this.hasTrendData(this.dashboard.billTrend);
    }
  },
  mounted() {
    this.fetchDashboard();
  },
  methods: {
    fetchDashboard() {
      this.isLoading = true;
      this.errorMessage = "";

      getAdminDashboardOverview()
        .then(
          function(result) {
            var payload = result && result.data ? result.data : result;
            this.dashboard = this.normalizeDashboard(payload || {});
          }.bind(this)
        )
        .catch(
          function(error) {
            this.dashboard = emptyDashboard();
            this.errorMessage = this.buildErrorMessage(error);
          }.bind(this)
        )
        .finally(
          function() {
            this.isLoading = false;
          }.bind(this)
        );
    },
    normalizeDashboard(data) {
      var pendingUsers = asNumber(data.pending_users);

      return {
        stats: [
          {
            key: "totalUsers",
            label: "总用户数",
            value: formatCount(data.total_users),
            trend: "实时统计",
            tone: "normal"
          },
          {
            key: "pendingUsers",
            label: "待审批用户",
            value: formatCount(pendingUsers),
            trend: pendingUsers > 0 ? "需处理" : "已清零",
            tone: pendingUsers > 0 ? "alert" : "normal"
          },
          {
            key: "todayBills",
            label: "今日账单数",
            value: formatCount(data.today_bills),
            trend: "今日",
            tone: "up"
          },
          {
            key: "families",
            label: "家庭数量",
            value: formatCount(data.families),
            trend: "实时统计",
            tone: "normal"
          }
        ],
        userGrowth: buildTrend(data.user_growth, "新增用户", "#171717"),
        billTrend: buildTrend(data.bill_trend, "账单数", "#f6d34a"),
        recentUsers: this.normalizeRecentUsers(data.recent_users)
      };
    },
    normalizeRecentUsers(users) {
      if (!Array.isArray(users)) {
        return [];
      }

      return users.map(function(user) {
        var name = user.nickname || user.username || "未命名用户";

        return {
          id: user.id,
          username: user.username || "",
          name: name,
          email: user.email || "",
          status: normalizeUserStatus(user.status, user.approval_status),
          registerDate: formatDateTime(user.created_at)
        };
      });
    },
    hasTrendData(trend) {
      if (!trend || !trend.series || !trend.series.length) {
        return false;
      }

      return trend.series.some(function(item) {
        return Array.isArray(item.values) && item.values.some(function(value) {
          return asNumber(value) > 0;
        });
      });
    },
    buildErrorMessage(error) {
      if (error && error.response && error.response.data && error.response.data.message) {
        return error.response.data.message;
      }

      return "暂时无法获取后台概览数据";
    },
    trendClass(tone) {
      if (tone === "up") {
        return "admin-dashboard-page__trend admin-dashboard-page__trend--up";
      }

      if (tone === "alert") {
        return "admin-dashboard-page__trend admin-dashboard-page__trend--alert";
      }

      return "admin-dashboard-page__trend";
    },
    statusLabel(status) {
      if (status === "pending") {
        return "待审批";
      }

      if (status === "active") {
        return "已启用";
      }

      return "已禁用";
    },
    statusClass(status) {
      return "admin-dashboard-page__status admin-dashboard-page__status--" + status;
    },
    openUsers() {
      this.$router.push("/admin/users");
    }
  }
};
</script>

<style scoped>
.admin-dashboard-page {
  gap: 20px;
}

.admin-dashboard-page__hero {
  padding: 28px;
}

.admin-dashboard-page__trend {
  color: var(--text-subtle);
}

.admin-dashboard-page__trend--up {
  color: #15803d;
}

.admin-dashboard-page__trend--alert {
  color: #b45309;
}

.admin-dashboard-page__stat-loading span,
.admin-dashboard-page__stat-loading strong,
.admin-dashboard-page__stat-loading small,
.admin-dashboard-page__panel-loading span,
.admin-dashboard-page__panel-loading strong,
.admin-dashboard-page__panel-loading i {
  display: block;
  border-radius: 999px;
  background: linear-gradient(90deg, rgba(229, 231, 235, 0.65), rgba(255, 255, 255, 0.92), rgba(229, 231, 235, 0.65));
  background-size: 220% 100%;
  animation: admin-dashboard-loading 1.35s ease-in-out infinite;
}

.admin-dashboard-page__stat-loading span {
  width: 72px;
  height: 14px;
}

.admin-dashboard-page__stat-loading strong {
  width: 94px;
  height: 34px;
  margin: 18px 0 10px;
}

.admin-dashboard-page__stat-loading small {
  width: 64px;
  height: 12px;
}

.admin-dashboard-page__panel-loading {
  min-height: 390px;
}

.admin-dashboard-page__panel-loading span {
  width: 116px;
  height: 18px;
}

.admin-dashboard-page__panel-loading strong {
  width: 220px;
  height: 13px;
  margin: 12px 0 34px;
}

.admin-dashboard-page__panel-loading i {
  width: 100%;
  height: 260px;
  border-radius: 18px;
}

.admin-dashboard-page__error {
  align-items: flex-start;
  gap: 12px;
}

.admin-dashboard-page__error h3 {
  margin: 0;
}

.admin-dashboard-page__error p {
  margin: 0;
  color: var(--text-subtle);
}

.admin-dashboard-page__empty {
  min-height: 260px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px dashed var(--border-color);
  border-radius: 18px;
  color: var(--text-subtle);
  background: rgba(255, 255, 255, 0.62);
  font-weight: 600;
}

.admin-dashboard-page__user-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.admin-dashboard-page__user-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 16px 18px;
  border-radius: 18px;
  border: 1px solid var(--border-color);
  background: #ffffff;
}

.admin-dashboard-page__user-main {
  display: flex;
  align-items: center;
  gap: 14px;
  min-width: 0;
}

.admin-dashboard-page__avatar {
  width: 44px;
  height: 44px;
  flex: 0 0 44px;
  border-radius: 14px;
  background: rgba(23, 23, 23, 0.1);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
}

.admin-dashboard-page__user-main strong {
  display: block;
}

.admin-dashboard-page__user-main p {
  margin: 6px 0 0;
  color: var(--text-muted);
  font-size: 13px;
  word-break: break-all;
}

.admin-dashboard-page__user-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--text-subtle);
  font-size: 13px;
  white-space: nowrap;
}

.admin-dashboard-page__status {
  min-height: 30px;
  padding: 0 12px;
  border-radius: 999px;
  display: inline-flex;
  align-items: center;
  font-size: 12px;
  font-style: normal;
  font-weight: 700;
}

.admin-dashboard-page__status--pending {
  background: rgba(245, 158, 11, 0.15);
  color: #b45309;
}

.admin-dashboard-page__status--active {
  background: rgba(34, 197, 94, 0.15);
  color: #15803d;
}

.admin-dashboard-page__status--inactive {
  background: rgba(156, 163, 175, 0.2);
  color: #4b5563;
}

@keyframes admin-dashboard-loading {
  0% {
    background-position: 180% 0;
  }

  100% {
    background-position: -40% 0;
  }
}

@media (max-width: 1080px) {
  .admin-dashboard-page__hero .finance-stat-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 860px) {
  .admin-dashboard-page__user-row {
    flex-direction: column;
    align-items: flex-start;
  }

  .admin-dashboard-page__user-meta {
    flex-wrap: wrap;
    white-space: normal;
  }
}
</style>
