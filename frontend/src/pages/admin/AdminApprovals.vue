<template>
  <div class="finance-page admin-approvals-page">
    <section class="finance-hero finance-hero--soft admin-approvals-page__hero">
      <div class="finance-hero__eyebrow">
        <span class="finance-hero__eyebrow-mark">审</span>
        <span>管理后台 / 待审批用户</span>
      </div>

      <div class="finance-hero__headline">
        <h1>集中处理新注册用户的审批动作</h1>
        <p>页面优先落地审批工作流信息结构，审批接口后续联调接入。</p>
      </div>

      <div class="finance-stat-grid finance-stat-grid--triple">
        <article class="finance-stat-card">
          <span>待审批用户</span>
          <strong>{{ summary.pendingCount }}</strong>
          <small>可逐条批准或拒绝</small>
        </article>
        <article class="finance-stat-card">
          <span>当前已选择</span>
          <strong>{{ selectedCount }}</strong>
          <small>支持批量操作</small>
        </article>
        <article class="finance-stat-card">
          <span>最近操作</span>
          <strong>{{ recentActionLabel }}</strong>
          <small>仅记录原型操作轨迹</small>
        </article>
      </div>
    </section>

    <section class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>待审批列表</h3>
          <p>字段包含昵称、联系方式、注册时间和申请用途。</p>
        </div>
      </header>

      <div class="admin-approvals-page__list">
        <article v-for="user in users" :key="user.id" class="admin-approvals-page__row">
          <label class="admin-approvals-page__check">
            <input :checked="isSelected(user.id)" type="checkbox" @change="toggleSelect(user.id)" />
          </label>

          <div class="admin-approvals-page__main">
            <div class="admin-approvals-page__title">
              <strong>{{ user.name }}</strong>
              <em>待审批</em>
            </div>

            <div class="admin-approvals-page__meta-grid">
              <p>邮箱：{{ user.email }}</p>
              <p>手机：{{ user.phone }}</p>
              <p>注册：{{ user.registerDate }}</p>
              <p>用途：{{ user.reason }}</p>
            </div>
          </div>

          <div class="admin-approvals-page__actions">
            <button class="finance-button finance-button--primary" @click="approve(user)">批准</button>
            <button class="finance-button finance-button--ghost" @click="reject(user)">拒绝</button>
          </div>
        </article>
      </div>
    </section>

    <section class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>批量操作</h3>
          <p>当前阶段仅提供原型级反馈，不触发真实后端审批。</p>
        </div>
      </header>

      <div class="admin-approvals-page__bulk">
        <button class="finance-button finance-button--primary" :disabled="selectedCount === 0" @click="approveBatch">
          批量批准
        </button>
        <button class="finance-button finance-button--ghost" :disabled="selectedCount === 0" @click="rejectBatch">
          批量拒绝
        </button>
        <span>{{ selectedCount }} 位用户已勾选</span>
      </div>
    </section>
  </div>
</template>

<script>
import { getPendingApprovalsData } from "@/utils/adminPortalMock";

export default {
  name: "AdminApprovals",
  data() {
    var source = getPendingApprovalsData();

    return {
      users: source.users,
      summary: source.summary,
      selectedIds: [],
      recentActionLabel: "无"
    };
  },
  computed: {
    selectedCount() {
      return this.selectedIds.length;
    }
  },
  methods: {
    isSelected(userId) {
      return this.selectedIds.indexOf(userId) !== -1;
    },
    toggleSelect(userId) {
      var index = this.selectedIds.indexOf(userId);

      if (index === -1) {
        this.selectedIds.push(userId);
      } else {
        this.selectedIds.splice(index, 1);
      }
    },
    approve(user) {
      this.recentActionLabel = "批准 " + user.name;
      this.$message.success("已批准 " + user.name + "（原型操作）");
    },
    reject(user) {
      this.recentActionLabel = "拒绝 " + user.name;
      this.$message.warning("已拒绝 " + user.name + "（原型操作）");
    },
    approveBatch() {
      this.recentActionLabel = "批量批准";
      this.$message.success("已批量批准 " + this.selectedCount + " 位用户（原型操作）");
      this.selectedIds = [];
    },
    rejectBatch() {
      this.recentActionLabel = "批量拒绝";
      this.$message.warning("已批量拒绝 " + this.selectedCount + " 位用户（原型操作）");
      this.selectedIds = [];
    }
  }
};
</script>

<style scoped>
.admin-approvals-page {
  gap: 20px;
}

.admin-approvals-page__hero {
  padding: 28px;
}

.admin-approvals-page__list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.admin-approvals-page__row {
  display: grid;
  grid-template-columns: 36px minmax(0, 1fr) auto;
  gap: 14px;
  padding: 16px;
  border-radius: 18px;
  border: 1px solid var(--border-color);
  background: #ffffff;
}

.admin-approvals-page__check {
  display: inline-flex;
  justify-content: center;
  padding-top: 4px;
}

.admin-approvals-page__check input {
  width: 18px;
  height: 18px;
}

.admin-approvals-page__title {
  display: flex;
  align-items: center;
  gap: 10px;
}

.admin-approvals-page__title strong {
  font-size: 17px;
}

.admin-approvals-page__title em {
  min-height: 26px;
  padding: 0 10px;
  border-radius: 999px;
  background: rgba(245, 158, 11, 0.16);
  color: #b45309;
  display: inline-flex;
  align-items: center;
  font-size: 12px;
  font-style: normal;
  font-weight: 700;
}

.admin-approvals-page__meta-grid {
  margin-top: 10px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px 16px;
}

.admin-approvals-page__meta-grid p {
  margin: 0;
  color: var(--text-subtle);
  font-size: 13px;
}

.admin-approvals-page__actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.admin-approvals-page__actions .finance-button {
  min-height: 38px;
  padding: 0 14px;
  border-radius: 12px;
}

.admin-approvals-page__bulk {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
}

.admin-approvals-page__bulk span {
  color: var(--text-subtle);
  font-size: 13px;
}

@media (max-width: 980px) {
  .admin-approvals-page__row {
    grid-template-columns: 1fr;
  }

  .admin-approvals-page__meta-grid {
    grid-template-columns: 1fr;
  }

  .admin-approvals-page__actions {
    justify-content: flex-start;
  }
}
</style>
