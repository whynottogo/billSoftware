<template>
  <div class="admin-users-page">
    <section class="admin-users-page__header page-card">
      <div>
        <span class="admin-users-page__eyebrow">管理端 / 用户管理</span>
        <h1>用户列表</h1>
        <p>支持启用/禁用状态管理，并可进入只读账单详情查看单用户月度与年度汇总。</p>
      </div>
    </section>

    <section class="admin-users-page__summary">
      <article class="admin-stat">
        <span>启用用户</span>
        <strong>{{ enabledCount }}</strong>
      </article>
      <article class="admin-stat">
        <span>禁用用户</span>
        <strong>{{ disabledCount }}</strong>
      </article>
      <article class="admin-stat">
        <span>今日新增</span>
        <strong>{{ todayNewCount }}</strong>
      </article>
    </section>

    <section class="admin-users-page__filters page-card">
      <label class="admin-filter">
        <span>搜索用户</span>
        <input v-model.trim="keyword" type="text" placeholder="搜索账号、昵称、手机号或邮箱" />
      </label>

      <label class="admin-filter admin-filter--compact">
        <span>状态筛选</span>
        <select v-model="status">
          <option value="all">全部状态</option>
          <option value="启用">启用</option>
          <option value="禁用">禁用</option>
        </select>
      </label>
    </section>

    <section class="admin-users-table page-card">
      <table>
        <thead>
          <tr>
            <th>账号</th>
            <th>昵称</th>
            <th>手机号</th>
            <th>邮箱</th>
            <th>状态</th>
            <th>注册时间</th>
            <th>最近登录</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in filteredRows" :key="row.userId">
            <td>
              <div class="user-cell">
                <span class="user-cell__avatar">{{ row.avatarLabel }}</span>
                <div>
                  <strong>{{ row.username }}</strong>
                  <small>{{ row.billCount }} 笔账单</small>
                </div>
              </div>
            </td>
            <td>{{ row.nickname }}</td>
            <td>{{ row.phone }}</td>
            <td>{{ row.email }}</td>
            <td>
              <span :class="['status-pill', row.status === '启用' ? 'is-enabled' : 'is-disabled']">
                {{ row.status }}
              </span>
            </td>
            <td>{{ row.registerDate }}</td>
            <td>{{ row.lastLogin }}</td>
            <td>
              <div class="table-actions">
                <button class="table-actions__detail" @click="openDetail(row.userId)">查看详情</button>
                <button
                  class="table-actions__toggle"
                  :class="[row.status === '启用' ? 'is-danger' : 'is-success', isStatusUpdating(row.userId) ? 'is-pending' : '']"
                  :disabled="isStatusUpdating(row.userId)"
                  @click="toggleStatus(row.userId)"
                >
                  {{ isStatusUpdating(row.userId) ? "处理中..." : row.status === "启用" ? "禁用" : "启用" }}
                </button>
              </div>
            </td>
          </tr>
          <tr v-if="filteredRows.length === 0">
            <td colspan="8" class="admin-users-table__empty">当前筛选条件下没有用户数据。</td>
          </tr>
        </tbody>
      </table>
    </section>
  </div>
</template>

<script>
import { changeAdminUserStatus, listAdminUsers } from "@/api/adminUsers";

function asText(value, fallback) {
  if (value === undefined || value === null || value === "") {
    return fallback;
  }

  return String(value);
}

function toAvatarLabel(username) {
  const text = asText(username, "U");
  return text.slice(0, 1).toUpperCase();
}

function parseDate(input) {
  if (!input) {
    return null;
  }

  const date = new Date(input);
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

function isToday(input) {
  const date = parseDate(input);
  if (!date) {
    return false;
  }

  const now = new Date();
  return (
    date.getFullYear() === now.getFullYear() &&
    date.getMonth() === now.getMonth() &&
    date.getDate() === now.getDate()
  );
}

function mapUser(user) {
  const username = asText(user && user.username, "-");
  const createdAt = user ? user.created_at : null;

  return {
    userId: String(user && user.id ? user.id : ""),
    username: username,
    nickname: asText(user && user.nickname, "-"),
    phone: asText(user && user.phone, "-"),
    email: asText(user && user.email, "-"),
    status: user && Number(user.status) === 1 ? "启用" : "禁用",
    registerDate: formatDate(createdAt),
    lastLogin: "-",
    billCount: 0,
    todayNew: isToday(createdAt),
    avatarLabel: toAvatarLabel(username),
    avatarCompressed: asText(user && user.avatar_compressed, "")
  };
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

function buildErrorMessage(error, fallback) {
  if (error && error.response && error.response.data && error.response.data.message) {
    return error.response.data.message;
  }

  return fallback;
}

export default {
  name: "AdminUsers",
  data() {
    return {
      keyword: "",
      status: "all",
      rows: [],
      statusUpdatingMap: {}
    };
  },
  computed: {
    enabledCount() {
      return this.rows.filter(function(row) {
        return row.status === "启用";
      }).length;
    },
    disabledCount() {
      return this.rows.filter(function(row) {
        return row.status === "禁用";
      }).length;
    },
    todayNewCount() {
      return this.rows.filter(function(row) {
        return row.todayNew;
      }).length;
    },
    filteredRows() {
      const keyword = this.keyword.toLowerCase();

      return this.rows.filter(
        function(row) {
          const matchesStatus = this.status === "all" || row.status === this.status;
          const source = [row.username, row.nickname, row.phone, row.email].join(" ").toLowerCase();
          const matchesKeyword = !keyword || source.indexOf(keyword) !== -1;

          return matchesStatus && matchesKeyword;
        }.bind(this)
      );
    }
  },
  created() {
    this.refreshRows();
  },
  methods: {
    refreshRows() {
      return listAdminUsers()
        .then(
          function(result) {
            this.rows = extractList(result).map(mapUser);
          }.bind(this)
        )
        .catch(
          function(error) {
            this.rows = [];
            this.$message.error(buildErrorMessage(error, "用户列表加载失败，请检查后端服务"));
          }.bind(this)
        );
    },
    openDetail(userId) {
      this.$router.push(`/admin/users/${userId}`);
    },
    isStatusUpdating(userId) {
      return Boolean(this.statusUpdatingMap[userId]);
    },
    toggleStatus(userId) {
      const row = this.rows.find(function(item) {
        return item.userId === userId;
      });

      if (!row || this.isStatusUpdating(userId)) {
        return;
      }

      const nextStatus = row.status === "启用" ? 0 : 1;

      this.statusUpdatingMap = Object.assign({}, this.statusUpdatingMap, {
        [userId]: true
      });

      changeAdminUserStatus(userId, nextStatus)
        .then(
          function() {
            this.rows = this.rows.map(function(item) {
              if (item.userId !== userId) {
                return item;
              }

              return Object.assign({}, item, {
                status: nextStatus === 1 ? "启用" : "禁用"
              });
            });

            if (nextStatus === 0) {
              this.$message.warning("用户已禁用，用户端将不可登录。");
            } else {
              this.$message.success("用户已启用，可恢复用户端登录。");
            }
          }.bind(this)
        )
        .catch(
          function(error) {
            this.$message.error(buildErrorMessage(error, "状态更新失败，请稍后重试"));
          }.bind(this)
        )
        .finally(
          function() {
            const nextMap = Object.assign({}, this.statusUpdatingMap);
            delete nextMap[userId];
            this.statusUpdatingMap = nextMap;
          }.bind(this)
        );
    }
  }
};
</script>

<style scoped>
.admin-users-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.admin-users-page__header,
.admin-users-page__filters,
.admin-users-table {
  padding: 24px;
}

.admin-users-page__eyebrow {
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

.admin-users-page__header h1 {
  margin: 16px 0 10px;
  font-size: 34px;
}

.admin-users-page__header p {
  margin: 0;
  color: var(--text-subtle);
  line-height: 1.7;
}

.admin-users-page__summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.admin-stat {
  padding: 22px 24px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
}

.admin-stat span {
  display: block;
  margin-bottom: 12px;
  color: var(--text-muted);
  font-size: 14px;
}

.admin-stat strong {
  font-size: 26px;
}

.admin-users-page__filters {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 220px;
  gap: 16px;
}

.admin-filter {
  display: flex;
  flex-direction: column;
  gap: 8px;
  color: var(--text-main);
  font-size: 14px;
  font-weight: 600;
}

.admin-filter input,
.admin-filter select {
  width: 100%;
  height: 44px;
  padding: 0 16px;
  border-radius: 16px;
  border: 1px solid var(--border-color);
  background: var(--card-muted);
  color: var(--text-main);
}

.admin-filter input:focus,
.admin-filter select:focus {
  outline: none;
  border-color: rgba(23, 23, 23, 0.52);
  box-shadow: 0 0 0 4px rgba(23, 23, 23, 0.08);
}

.admin-users-table {
  overflow: hidden;
}

.admin-users-table table {
  width: 100%;
  border-collapse: collapse;
}

.admin-users-table thead {
  background: rgba(23, 23, 23, 0.04);
}

.admin-users-table th,
.admin-users-table td {
  padding: 16px 14px;
  border-bottom: 1px solid rgba(229, 231, 235, 0.95);
  text-align: left;
  vertical-align: middle;
}

.admin-users-table th {
  color: var(--text-muted);
  font-size: 13px;
  font-weight: 700;
}

.admin-users-table td {
  font-size: 14px;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-cell__avatar {
  width: 42px;
  height: 42px;
  border-radius: 16px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: rgba(23, 23, 23, 0.08);
  color: var(--text-main);
  font-weight: 700;
}

.user-cell strong {
  display: block;
}

.user-cell small {
  display: block;
  margin-top: 4px;
  color: var(--text-muted);
}

.status-pill {
  display: inline-flex;
  align-items: center;
  min-height: 32px;
  padding: 0 12px;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 700;
}

.status-pill.is-enabled {
  background: rgba(34, 197, 94, 0.1);
  color: var(--success-color);
}

.status-pill.is-disabled {
  background: rgba(239, 68, 68, 0.1);
  color: var(--danger-color);
}

.table-actions {
  display: flex;
  gap: 8px;
}

.table-actions button {
  min-height: 34px;
  padding: 0 12px;
  border-radius: 12px;
  border: 1px solid var(--border-color);
  background: #ffffff;
  color: var(--text-main);
  font-size: 13px;
  font-weight: 600;
}

.table-actions__toggle.is-danger {
  border-color: rgba(239, 68, 68, 0.45);
  color: var(--danger-color);
}

.table-actions__toggle.is-success {
  border-color: rgba(34, 197, 94, 0.45);
  color: var(--success-color);
}

.table-actions__toggle.is-pending,
.table-actions__toggle:disabled {
  opacity: 0.65;
  cursor: not-allowed;
}

.admin-users-table__empty {
  text-align: center;
  color: var(--text-muted);
  padding: 28px 0;
}

@media (max-width: 980px) {
  .admin-users-page__summary,
  .admin-users-page__filters {
    grid-template-columns: 1fr;
  }

  .admin-users-table {
    overflow-x: auto;
  }

  .admin-users-table table {
    min-width: 860px;
  }
}
</style>
