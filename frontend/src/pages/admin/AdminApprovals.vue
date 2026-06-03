<template>
  <div class="finance-page admin-approvals-page">
    <section class="finance-hero finance-hero--soft admin-approvals-page__hero">
      <div class="finance-hero__eyebrow">
        <span class="finance-hero__eyebrow-mark">审</span>
        <span>管理后台 / 待审批用户</span>
      </div>

      <div class="finance-hero__headline">
        <h1>集中处理新注册用户的审批动作</h1>
        <p>待审批清单已接入真实接口，批准后用户可登录，拒绝后用户继续保持不可登录。</p>
      </div>

      <div class="finance-stat-grid finance-stat-grid--triple">
        <article class="finance-stat-card">
          <span>待审批用户</span>
          <strong>{{ summary.pendingCount }}</strong>
          <small>{{ loading ? "正在同步" : "来自真实用户表" }}</small>
        </article>
        <article class="finance-stat-card">
          <span>当前已选择</span>
          <strong>{{ selectedCount }}</strong>
          <small>支持批量批准或拒绝</small>
        </article>
        <article class="finance-stat-card">
          <span>最近操作</span>
          <strong>{{ recentActionLabel }}</strong>
          <small>{{ recentActionHint }}</small>
        </article>
      </div>
    </section>

    <section class="page-card finance-panel">
      <header class="finance-panel__header admin-approvals-page__panel-header">
        <div>
          <h3>待审批列表</h3>
          <p>字段包含昵称、邮箱、手机号、注册时间和申请用途。</p>
        </div>

        <button class="finance-button finance-button--ghost" :disabled="loading" @click="loadUsers">
          {{ loading ? "加载中..." : "重试" }}
        </button>
      </header>

      <div v-if="errorMessage" class="admin-approvals-page__state is-error">
        <strong>审批列表加载失败</strong>
        <span>{{ errorMessage }}</span>
        <button class="finance-button finance-button--primary" :disabled="loading" @click="loadUsers">重试</button>
      </div>

      <div v-else-if="loading && users.length === 0" class="admin-approvals-page__state">
        <strong>正在加载待审批用户</strong>
        <span>请稍候，系统正在读取管理端审批接口。</span>
      </div>

      <div v-else-if="users.length === 0" class="admin-approvals-page__state">
        <strong>暂无待审批用户</strong>
        <span>新注册且尚未处理的用户会显示在这里。</span>
      </div>

      <div v-else class="admin-approvals-page__list">
        <article v-for="user in users" :key="user.id" class="admin-approvals-page__row">
          <label class="admin-approvals-page__check">
            <input
              :checked="isSelected(user.id)"
              :disabled="isRowBusy(user.id)"
              type="checkbox"
              @change="toggleSelect(user.id)"
            />
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
            <button
              class="finance-button finance-button--primary"
              :disabled="isRowBusy(user.id)"
              @click="approve(user)"
            >
              {{ rowActionLabel(user.id, "批准") }}
            </button>
            <button
              class="finance-button finance-button--ghost"
              :disabled="isRowBusy(user.id)"
              @click="reject(user)"
            >
              {{ rowActionLabel(user.id, "拒绝") }}
            </button>
          </div>
        </article>
      </div>
    </section>

    <section class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>批量操作</h3>
          <p>批量动作只处理当前勾选的待审批用户。</p>
        </div>
      </header>

      <div class="admin-approvals-page__bulk">
        <button class="finance-button finance-button--primary" :disabled="!canRunBatch" @click="approveBatch">
          {{ batchBusy ? "处理中..." : "批量批准" }}
        </button>
        <button class="finance-button finance-button--ghost" :disabled="!canRunBatch" @click="rejectBatch">
          {{ batchBusy ? "处理中..." : "批量拒绝" }}
        </button>
        <span>{{ selectedCount }} 位用户已勾选</span>
      </div>
    </section>
  </div>
</template>

<script>
import { ElMessageBox } from "element-plus";
import {
  approveApprovalUser,
  batchApproveApprovalUsers,
  batchRejectApprovalUsers,
  listPendingApprovalUsers,
  rejectApprovalUser
} from "@/api/adminApprovals";

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

  var date = new Date(input);
  if (Number.isNaN(date.getTime())) {
    return null;
  }

  return date;
}

function pad(value) {
  return value < 10 ? "0" + value : String(value);
}

function formatDate(input) {
  var date = parseDate(input);
  if (!date) {
    return "-";
  }

  return (
    date.getFullYear() +
    "-" +
    pad(date.getMonth() + 1) +
    "-" +
    pad(date.getDate()) +
    " " +
    pad(date.getHours()) +
    ":" +
    pad(date.getMinutes())
  );
}

function extractList(result) {
  if (result && result.data && Array.isArray(result.data.list)) {
    return result.data.list;
  }

  if (result && Array.isArray(result.list)) {
    return result.list;
  }

  return [];
}

function extractPendingCount(result, fallback) {
  if (result && result.data && result.data.summary && typeof result.data.summary.pending_count === "number") {
    return result.data.summary.pending_count;
  }

  if (result && result.summary && typeof result.summary.pending_count === "number") {
    return result.summary.pending_count;
  }

  return fallback;
}

function buildErrorMessage(error, fallback) {
  if (error && error.response && error.response.data && error.response.data.message) {
    return error.response.data.message;
  }

  return fallback;
}

function mapUser(user) {
  var username = asText(user && user.username, "-");
  var nickname = asText(user && user.nickname, username);

  return {
    id: String(user && user.id ? user.id : ""),
    name: nickname,
    username: username,
    email: asText(user && user.email, "-"),
    phone: asText(user && user.phone, "-"),
    registerDate: formatDate(user && user.created_at),
    reason: asText(user && user.application_purpose, "注册后等待管理员审批")
  };
}

export default {
  name: "AdminApprovals",
  data() {
    return {
      users: [],
      summary: {
        pendingCount: 0
      },
      selectedIds: [],
      busyMap: {},
      loading: false,
      batchBusy: false,
      errorMessage: "",
      recentActionLabel: "无",
      recentActionHint: "等待审批动作"
    };
  },
  computed: {
    selectedCount() {
      return this.selectedIds.length;
    },
    canRunBatch() {
      return this.selectedCount > 0 && !this.batchBusy && !this.loading;
    }
  },
  created() {
    this.loadUsers();
  },
  methods: {
    loadUsers() {
      this.loading = true;
      this.errorMessage = "";

      return listPendingApprovalUsers()
        .then(
          function(result) {
            var list = extractList(result).map(mapUser);
            this.users = list;
            this.summary.pendingCount = extractPendingCount(result, list.length);
            this.selectedIds = this.selectedIds.filter(
              function(userId) {
                return list.some(function(user) {
                  return user.id === userId;
                });
              }.bind(this)
            );
          }.bind(this)
        )
        .catch(
          function(error) {
            this.errorMessage = buildErrorMessage(error, "待审批用户加载失败，请检查后端服务");
            this.users = [];
            this.summary.pendingCount = 0;
          }.bind(this)
        )
        .finally(
          function() {
            this.loading = false;
          }.bind(this)
        );
    },
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
    isRowBusy(userId) {
      return Boolean(this.busyMap[userId]) || this.batchBusy;
    },
    setRowBusy(userId, value) {
      var nextMap = Object.assign({}, this.busyMap);
      if (value) {
        nextMap[userId] = true;
      } else {
        delete nextMap[userId];
      }

      this.busyMap = nextMap;
    },
    rowActionLabel(userId, label) {
      return this.isRowBusy(userId) ? "处理中..." : label;
    },
    removeUsers(userIds) {
      this.users = this.users.filter(function(user) {
        return userIds.indexOf(user.id) === -1;
      });
      this.selectedIds = this.selectedIds.filter(function(userId) {
        return userIds.indexOf(userId) === -1;
      });
      this.summary.pendingCount = this.users.length;
    },
    approve(user) {
      if (this.isRowBusy(user.id)) {
        return;
      }

      this.setRowBusy(user.id, true);
      approveApprovalUser(user.id)
        .then(
          function() {
            this.removeUsers([user.id]);
            this.recentActionLabel = "批准 " + user.name;
            this.recentActionHint = "已允许用户端登录";
            this.$message.success("已批准 " + user.name);
          }.bind(this)
        )
        .catch(
          function(error) {
            this.$message.error(buildErrorMessage(error, "批准失败，请稍后重试"));
          }.bind(this)
        )
        .finally(
          function() {
            this.setRowBusy(user.id, false);
          }.bind(this)
        );
    },
    reject(user) {
      if (this.isRowBusy(user.id)) {
        return;
      }

      ElMessageBox.prompt("请输入拒绝原因", "拒绝用户", {
        confirmButtonText: "确认拒绝",
        cancelButtonText: "取消",
        inputPlaceholder: "例如：资料不完整",
        inputType: "textarea"
      })
        .then(
          function(result) {
            var remark = result && result.value ? result.value : "";
            this.setRowBusy(user.id, true);

            return rejectApprovalUser(user.id, remark)
              .then(
                function() {
                  this.removeUsers([user.id]);
                  this.recentActionLabel = "拒绝 " + user.name;
                  this.recentActionHint = "用户端登录保持禁用";
                  this.$message.warning("已拒绝 " + user.name);
                }.bind(this)
              )
              .catch(
                function(error) {
                  this.$message.error(buildErrorMessage(error, "拒绝失败，请稍后重试"));
                }.bind(this)
              )
              .finally(
                function() {
                  this.setRowBusy(user.id, false);
                }.bind(this)
              );
          }.bind(this)
        )
        .catch(function() {});
    },
    approveBatch() {
      if (!this.canRunBatch) {
        return;
      }

      var userIds = this.selectedIds.slice();
      this.batchBusy = true;
      batchApproveApprovalUsers(userIds)
        .then(
          function() {
            this.removeUsers(userIds);
            this.recentActionLabel = "批量批准";
            this.recentActionHint = "已处理 " + userIds.length + " 位用户";
            this.$message.success("已批量批准 " + userIds.length + " 位用户");
          }.bind(this)
        )
        .catch(
          function(error) {
            this.$message.error(buildErrorMessage(error, "批量批准失败，请稍后重试"));
          }.bind(this)
        )
        .finally(
          function() {
            this.batchBusy = false;
          }.bind(this)
        );
    },
    rejectBatch() {
      if (!this.canRunBatch) {
        return;
      }

      var userIds = this.selectedIds.slice();
      ElMessageBox.prompt("请输入统一拒绝原因", "批量拒绝用户", {
        confirmButtonText: "确认拒绝",
        cancelButtonText: "取消",
        inputPlaceholder: "例如：批量资料待补充",
        inputType: "textarea"
      })
        .then(
          function(result) {
            var remark = result && result.value ? result.value : "";
            this.batchBusy = true;

            return batchRejectApprovalUsers(userIds, remark)
              .then(
                function() {
                  this.removeUsers(userIds);
                  this.recentActionLabel = "批量拒绝";
                  this.recentActionHint = "已处理 " + userIds.length + " 位用户";
                  this.$message.warning("已批量拒绝 " + userIds.length + " 位用户");
                }.bind(this)
              )
              .catch(
                function(error) {
                  this.$message.error(buildErrorMessage(error, "批量拒绝失败，请稍后重试"));
                }.bind(this)
              )
              .finally(
                function() {
                  this.batchBusy = false;
                }.bind(this)
              );
          }.bind(this)
        )
        .catch(function() {});
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

.admin-approvals-page__panel-header {
  gap: 16px;
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
  border-radius: 8px;
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
  accent-color: #2563eb;
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

.admin-approvals-page__actions .finance-button,
.admin-approvals-page__panel-header .finance-button {
  min-height: 38px;
  padding: 0 14px;
  border-radius: 8px;
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

.admin-approvals-page__state {
  min-height: 180px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border: 1px dashed var(--border-color);
  border-radius: 8px;
  background: #fafafa;
  text-align: center;
}

.admin-approvals-page__state strong {
  color: var(--text-main);
  font-size: 17px;
}

.admin-approvals-page__state span {
  color: var(--text-subtle);
  font-size: 13px;
}

.admin-approvals-page__state.is-error {
  border-color: rgba(220, 38, 38, 0.28);
  background: rgba(254, 242, 242, 0.82);
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
