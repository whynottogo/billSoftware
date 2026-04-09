<template>
  <div class="page-shell admin-shell">
    <aside class="admin-shell__sidebar">
      <div class="admin-shell__brand">
        <span class="admin-shell__mark">AD</span>
        <div>
          <strong>Bill Software</strong>
          <p>管理后台</p>
        </div>
      </div>

      <nav class="admin-shell__nav">
        <button
          v-for="item in navItems"
          :key="item.label"
          :class="['admin-shell__nav-item', activeLabel === item.label ? 'is-active' : '']"
          @click="navigate(item)"
        >
          <span class="admin-shell__nav-badge">{{ item.badge }}</span>
          <span>{{ item.label }}</span>
        </button>
      </nav>

      <div class="admin-shell__footer">
        <div class="admin-shell__session">
          <p>当前角色</p>
          <strong>唯一管理员</strong>
          <span>管理员与普通用户会话严格隔离</span>
        </div>
        <button class="admin-shell__logout" @click="logout">退出后台</button>
      </div>
    </aside>

    <main class="admin-shell__main">
      <router-view />
    </main>
  </div>
</template>

<script>
export default {
  name: "AdminLayout",
  data() {
    return {
      navItems: [
        { label: "数据概览", badge: "概览", path: "/admin/dashboard" },
        { label: "用户管理", badge: "用户", path: "/admin/users" },
        { label: "待审批用户", badge: "审批", path: "/admin/approvals" },
        { label: "家庭管理", badge: "家庭", path: "/admin/families" }
      ]
    };
  },
  computed: {
    activeLabel() {
      var path = this.$route.path;

      if (path.indexOf("/admin/dashboard") === 0) {
        return "数据概览";
      }

      if (path.indexOf("/admin/users") === 0) {
        return "用户管理";
      }

      if (path.indexOf("/admin/approvals") === 0) {
        return "待审批用户";
      }

      if (path.indexOf("/admin/families") === 0) {
        return "家庭管理";
      }

      return "数据概览";
    }
  },
  methods: {
    navigate(item) {
      if (this.$route.path !== item.path) {
        this.$router.push(item.path);
      }
    },
    logout() {
      localStorage.removeItem("bill_admin_token");
      this.$router.push("/admin/login");
    }
  }
};
</script>

<style scoped>
.admin-shell {
  display: flex;
  min-height: 100vh;
}

.admin-shell__sidebar {
  width: 248px;
  padding: 24px 18px 20px;
  border-right: 1px solid var(--border-color);
  background: rgba(255, 255, 255, 0.96);
  display: flex;
  flex-direction: column;
}

.admin-shell__brand {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 4px 6px 24px;
}

.admin-shell__mark {
  width: 48px;
  height: 48px;
  border-radius: 16px;
  background: #171717;
  color: #ffffff;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 800;
  letter-spacing: 0.06em;
  box-shadow: var(--shadow-sm);
}

.admin-shell__brand strong {
  display: block;
  font-size: 17px;
}

.admin-shell__brand p {
  margin: 4px 0 0;
  color: var(--text-muted);
  font-size: 12px;
}

.admin-shell__nav {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.admin-shell__nav-item {
  border: none;
  background: transparent;
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  min-height: 52px;
  padding: 10px 12px;
  border-radius: 18px;
  color: var(--text-subtle);
  font-weight: 600;
  text-align: left;
  transition: all 0.2s ease;
}

.admin-shell__nav-item:hover {
  background: rgba(23, 23, 23, 0.06);
  color: var(--text-main);
}

.admin-shell__nav-item.is-active {
  background: #171717;
  color: #ffffff;
  box-shadow: var(--shadow-sm);
}

.admin-shell__nav-badge {
  min-width: 38px;
  height: 34px;
  padding: 0 10px;
  border-radius: 12px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: rgba(23, 23, 23, 0.07);
  font-size: 12px;
  font-weight: 700;
}

.admin-shell__nav-item.is-active .admin-shell__nav-badge {
  background: rgba(255, 255, 255, 0.12);
}

.admin-shell__footer {
  margin-top: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding-top: 22px;
}

.admin-shell__session {
  padding: 16px;
  border-radius: 20px;
  background: linear-gradient(180deg, #f7f7f8 0%, #ffffff 100%);
  border: 1px solid var(--border-color);
}

.admin-shell__session p,
.admin-shell__session span {
  margin: 0;
  color: var(--text-subtle);
  font-size: 12px;
  line-height: 1.6;
}

.admin-shell__session strong {
  display: block;
  margin: 4px 0;
  font-size: 15px;
}

.admin-shell__logout {
  min-height: 46px;
  border-radius: 16px;
  border: 1px solid var(--border-color);
  background: #ffffff;
  color: var(--text-main);
  font-weight: 600;
}

.admin-shell__main {
  flex: 1;
  padding: 28px 30px;
  overflow-y: auto;
}

@media (max-width: 960px) {
  .admin-shell {
    flex-direction: column;
  }

  .admin-shell__sidebar {
    width: auto;
    border-right: none;
    border-bottom: 1px solid var(--border-color);
  }

  .admin-shell__nav {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .admin-shell__main {
    padding: 20px;
  }
}
</style>
