<template>
  <div class="kickout-view">
    <section class="kickout-card page-card">
      <span class="kickout-card__tag">SESSION NOTICE</span>
      <h1>会话已失效</h1>
      <p>{{ description }}</p>

      <div class="kickout-highlight">
        <strong>统一提示文案</strong>
        <span>账号已在其他地方登录，请重新登录。</span>
      </div>

      <div class="kickout-actions">
        <button class="finance-button finance-button--primary" @click="goLogin">立即登录</button>
        <button class="finance-button" @click="goUserLoginWithCountdown">
          {{ countdownLabel }}
        </button>
      </div>
    </section>
  </div>
</template>

<script>
import { clearUserSession } from "@/utils/userProfileMock";

export default {
  name: "UserSessionKickout",
  data() {
    return {
      timer: null,
      countdown: 5
    };
  },
  computed: {
    description() {
      const reason = this.$route.query.reason || "kicked";

      if (reason === "kicked") {
        return "当前账号在其他设备完成了新的登录，本地会话已自动失效。";
      }

      return "当前登录会话已过期，请重新登录后继续使用。";
    },
    countdownLabel() {
      return `${this.countdown}s 后自动前往登录`;
    }
  },
  mounted() {
    clearUserSession("kicked");
    this.timer = window.setInterval(
      function () {
        this.countdown -= 1;
        if (this.countdown <= 0) {
          this.goLogin();
        }
      }.bind(this),
      1000
    );
  },
  beforeUnmount() {
    if (this.timer) {
      window.clearInterval(this.timer);
      this.timer = null;
    }
  },
  methods: {
    goLogin() {
      this.goUserLoginWithCountdown();
    },
    goUserLoginWithCountdown() {
      if (this.timer) {
        window.clearInterval(this.timer);
        this.timer = null;
      }

      this.$router.push({
        path: "/user/login",
        query: {
          reason: "kicked"
        }
      });
    }
  }
};
</script>

<style scoped>
.kickout-view {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background:
    radial-gradient(circle at 15% 18%, rgba(246, 211, 74, 0.22) 0%, rgba(255, 255, 255, 0) 35%),
    linear-gradient(180deg, rgba(255, 248, 224, 0.95) 0%, rgba(245, 246, 248, 0.96) 100%);
}

.kickout-card {
  width: min(680px, 100%);
  padding: 34px;
}

.kickout-card__tag {
  display: inline-flex;
  align-items: center;
  min-height: 32px;
  border-radius: 999px;
  padding: 0 12px;
  background: rgba(255, 244, 214, 0.8);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.06em;
}

.kickout-card h1 {
  margin: 14px 0 0;
  font-size: 34px;
}

.kickout-card p {
  margin: 10px 0 0;
  color: var(--text-subtle);
  line-height: 1.7;
}

.kickout-highlight {
  margin-top: 20px;
  border-radius: 18px;
  border: 1px solid rgba(246, 211, 74, 0.45);
  background: rgba(255, 250, 230, 0.84);
  padding: 15px 16px;
}

.kickout-highlight strong,
.kickout-highlight span {
  display: block;
}

.kickout-highlight span {
  margin-top: 8px;
  color: var(--text-subtle);
}

.kickout-actions {
  margin-top: 20px;
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}
</style>
