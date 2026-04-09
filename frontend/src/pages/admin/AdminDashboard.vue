<template>
  <div class="finance-page admin-dashboard-page">
    <section class="finance-hero finance-hero--soft admin-dashboard-page__hero">
      <div class="finance-hero__eyebrow">
        <span class="finance-hero__eyebrow-mark">概</span>
        <span>管理后台 / 数据概览</span>
      </div>

      <div class="finance-hero__headline">
        <h1>一屏查看运营走势和待处理事项</h1>
        <p>对齐 Figma Make 的后台门户节奏，先完成原型级数据可视化承载。</p>
      </div>

      <div class="finance-stat-grid">
        <article v-for="item in dashboard.stats" :key="item.key" class="finance-stat-card">
          <span>{{ item.label }}</span>
          <strong>{{ item.value }}</strong>
          <small :class="trendClass(item.tone)">{{ item.trend }}</small>
        </article>
      </div>
    </section>

    <section class="finance-grid-2">
      <article class="page-card finance-panel">
        <header class="finance-panel__header">
          <div>
            <h3>用户增长趋势</h3>
            <p>最近 6 个月新增用户曲线。</p>
          </div>
        </header>
        <SimpleLineChart :labels="dashboard.userGrowth.labels" :series="dashboard.userGrowth.series" :height="300" />
      </article>

      <article class="page-card finance-panel">
        <header class="finance-panel__header">
          <div>
            <h3>账单数量统计</h3>
            <p>最近 6 天账单提交情况。</p>
          </div>
        </header>
        <SimpleGroupedBarChart :labels="dashboard.billTrend.labels" :series="dashboard.billTrend.series" :height="300" />
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

      <div class="admin-dashboard-page__user-list">
        <article v-for="user in dashboard.recentUsers" :key="user.id" class="admin-dashboard-page__user-row">
          <div class="admin-dashboard-page__user-main">
            <span class="admin-dashboard-page__avatar">{{ user.name.slice(0, 1) }}</span>
            <div>
              <strong>{{ user.name }}</strong>
              <p>{{ user.email }}</p>
            </div>
          </div>
          <div class="admin-dashboard-page__user-meta">
            <span>{{ user.registerDate }}</span>
            <em :class="statusClass(user.status)">{{ statusLabel(user.status) }}</em>
          </div>
        </article>
      </div>
    </section>
  </div>
</template>

<script>
import SimpleGroupedBarChart from "@/components/SimpleGroupedBarChart.vue";
import SimpleLineChart from "@/components/SimpleLineChart.vue";
import { getAdminDashboardData } from "@/utils/adminPortalMock";

export default {
  name: "AdminDashboard",
  components: {
    SimpleLineChart: SimpleLineChart,
    SimpleGroupedBarChart: SimpleGroupedBarChart
  },
  data() {
    return {
      dashboard: getAdminDashboardData()
    };
  },
  methods: {
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
}

.admin-dashboard-page__avatar {
  width: 44px;
  height: 44px;
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
}

.admin-dashboard-page__user-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--text-subtle);
  font-size: 13px;
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
}
</style>
