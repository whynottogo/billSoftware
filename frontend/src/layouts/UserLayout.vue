<template>
  <div class="page-shell user-shell">
    <aside class="user-shell__sidebar">
      <div class="user-shell__brand">
        <span class="user-shell__mark">B</span>
        <div>
          <strong>Bill Software</strong>
          <p>轻松记账</p>
        </div>
      </div>

      <nav class="user-shell__nav">
        <button
          v-for="item in navItems"
          :key="item.path"
          :class="['user-shell__nav-item', activePath === item.path ? 'is-active' : '']"
          @click="navigate(item.path)"
        >
          <span class="user-shell__nav-badge">{{ item.badge }}</span>
          <span>{{ item.label }}</span>
        </button>
      </nav>

      <div class="user-shell__footer">
        <div class="user-shell__session">
          <p>当前会话</p>
          <strong>{{ sessionLabel }}</strong>
          <span>用户端与管理端登录态相互隔离</span>
        </div>
        <button class="user-shell__logout" @click="logout">退出登录</button>
      </div>
    </aside>

    <main class="user-shell__main">
      <router-view />
    </main>
  </div>
</template>

<script>
export default {
  name: "UserLayout",
  data() {
    return {
      navItems: [
        { path: "/user/ledger", label: "当月收支", badge: "收支" },
        { path: "/user/bills/month", label: "账单", badge: "账单" },
        { path: "/user/budget/month", label: "预算", badge: "预算" },
        { path: "/user/assets", label: "资产管家", badge: "资产" },
        { path: "/user/charts/expense", label: "图表", badge: "图表" },
        { path: "/user/profile", label: "个人信息", badge: "资料" },
        { path: "/user/families", label: "家庭功能", badge: "家庭" }
      ]
    };
  },
  computed: {
    activePath() {
      const path = this.$route.path;

      if (path.indexOf("/user/budget") === 0) return "/user/budget/month";
      if (path.indexOf("/user/bills") === 0) return "/user/bills/month";
      if (path.indexOf("/user/charts") === 0) return "/user/charts/expense";
      if (path.indexOf("/user/assets") === 0) return "/user/assets";
      if (path.indexOf("/user/families") === 0) return "/user/families";
      if (path.indexOf("/user/profile") === 0) return "/user/profile";

      return "/user/ledger";
    },
    sessionLabel() {
      const rawProfile = localStorage.getItem("bill_user_profile");

      if (!rawProfile) {
        return "用户端工作区";
      }

      try {
        const profile = JSON.parse(rawProfile);
        return profile.account || profile.username || "用户端工作区";
      } catch (error) {
        return "用户端工作区";
      }
    }
  },
  methods: {
    navigate(path) {
      if (this.$route.path !== path) {
        this.$router.push(path);
      }
    },
    logout() {
      localStorage.removeItem("bill_user_token");
      localStorage.removeItem("bill_user_profile");
      this.$router.push("/user/login");
    }
  }
};
</script>

<style scoped>
.user-shell {
  display: flex;
  min-height: 100vh;
}

.user-shell__sidebar {
  width: 248px;
  padding: 24px 18px 20px;
  border-right: 1px solid var(--border-color);
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(18px);
  display: flex;
  flex-direction: column;
}

.user-shell__brand {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 4px 6px 24px;
}

.user-shell__mark {
  width: 48px;
  height: 48px;
  border-radius: 16px;
  background: var(--brand-color);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  font-weight: 800;
  color: var(--text-main);
  box-shadow: var(--shadow-sm);
}

.user-shell__brand strong {
  display: block;
  font-size: 17px;
}

.user-shell__brand p {
  margin: 4px 0 0;
  color: var(--text-muted);
  font-size: 12px;
}

.user-shell__nav {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.user-shell__nav-item {
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
  transition: all 0.2s ease;
  text-align: left;
}

.user-shell__nav-item:hover {
  background: rgba(246, 211, 74, 0.18);
  color: var(--text-main);
}

.user-shell__nav-item.is-active {
  background: var(--brand-color);
  color: var(--text-main);
  box-shadow: var(--shadow-sm);
}

.user-shell__nav-badge {
  min-width: 38px;
  height: 34px;
  padding: 0 10px;
  border-radius: 12px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: rgba(23, 23, 23, 0.06);
  font-size: 12px;
  font-weight: 700;
}

.user-shell__nav-item.is-active .user-shell__nav-badge {
  background: rgba(255, 255, 255, 0.55);
}

.user-shell__footer {
  margin-top: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding-top: 22px;
}

.user-shell__session {
  padding: 16px;
  border-radius: 20px;
  background: linear-gradient(180deg, #fff8db 0%, #fffef7 100%);
  border: 1px solid rgba(246, 211, 74, 0.55);
}

.user-shell__session p,
.user-shell__session span {
  margin: 0;
  color: var(--text-subtle);
  font-size: 12px;
  line-height: 1.6;
}

.user-shell__session strong {
  display: block;
  margin: 4px 0;
  font-size: 15px;
}

.user-shell__logout {
  min-height: 46px;
  border-radius: 16px;
  border: 1px solid var(--border-color);
  background: #ffffff;
  color: var(--text-main);
  font-weight: 600;
}

.user-shell__main {
  flex: 1;
  padding: 28px 30px;
  overflow-y: auto;
}

@media (max-width: 960px) {
  .user-shell {
    flex-direction: column;
  }

  .user-shell__sidebar {
    width: auto;
    border-right: none;
    border-bottom: 1px solid var(--border-color);
  }

  .user-shell__nav {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .user-shell__main {
    padding: 20px;
  }
}
</style>
