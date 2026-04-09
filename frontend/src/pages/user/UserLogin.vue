<template>
  <div class="auth-view auth-view--user">
    <section class="auth-view__hero">
      <div class="auth-view__glow auth-view__glow--primary"></div>
      <div class="auth-view__glow auth-view__glow--soft"></div>
      <div class="auth-view__hero-content">
        <div class="auth-brand">
          <span class="auth-brand__mark">B</span>
          <div>
            <h1>Bill Software</h1>
            <p>轻松记账，明智理财</p>
          </div>
        </div>

        <div class="auth-copy">
          <h2>把记账这件事，做得清楚一点。</h2>
          <p>
            面向个人与家庭的桌面端记账工作台，帮助你快速录入每一笔收支，
            用清晰的统计卡和分类视角看懂自己的财务节奏。
          </p>
        </div>

        <div class="feature-list">
          <article class="feature-card">
            <strong>趋势看板</strong>
            <p>月度收入、支出与结余统一汇总，默认首页就是你的账本概览。</p>
          </article>
          <article class="feature-card">
            <strong>分类统计</strong>
            <p>账单、预算和图表采用一致的分类体系，方便快速回看消费结构。</p>
          </article>
          <article class="feature-card">
            <strong>家庭协作</strong>
            <p>支持家庭维度的汇总查看，但不干扰个人账本的独立管理。</p>
          </article>
        </div>
      </div>
    </section>

    <section class="auth-view__panel">
      <div class="auth-panel page-card">
        <span class="auth-panel__eyebrow">用户端登录</span>
        <h3>欢迎回来</h3>
        <p class="auth-panel__description">支持“用户名或手机号 + 密码”登录。</p>

        <div v-if="notice" :class="['auth-notice', notice.type === 'success' ? 'is-success' : 'is-warning']">
          <strong>{{ notice.title }}</strong>
          <span>{{ notice.text }}</span>
        </div>

        <form class="auth-form" @submit.prevent="submit">
          <label class="auth-form__field">
            <span>用户名或手机号</span>
            <input
              v-model.trim="form.account"
              type="text"
              autocomplete="username"
              placeholder="请输入用户名或手机号"
            />
          </label>

          <label class="auth-form__field">
            <span>密码</span>
            <div class="auth-form__password">
              <input
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                placeholder="请输入密码"
              />
              <button type="button" class="auth-form__toggle" @click="showPassword = !showPassword">
                {{ showPassword ? "隐藏" : "显示" }}
              </button>
            </div>
          </label>

          <div class="auth-form__meta">
            <label class="auth-checkbox">
              <input v-model="rememberMe" type="checkbox" />
              <span>记住我</span>
            </label>
            <button type="button" class="auth-link" @click="$router.push('/admin/login')">
              管理员登录入口
            </button>
          </div>

          <button type="submit" class="auth-form__submit">登录</button>
        </form>

        <div class="auth-panel__footer">
          <span>还没有账户？</span>
          <button type="button" class="auth-link" @click="$router.push('/user/register')">立即注册</button>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import { userLogin } from "@/api/userAuth";
import { consumeSessionNotice } from "@/utils/userProfileMock";

export default {
  name: "UserLogin",
  data() {
    return {
      form: {
        account: "",
        password: ""
      },
      rememberMe: true,
      showPassword: false,
      cachedSessionNotice: null
    };
  },
  created() {
    this.cachedSessionNotice = consumeSessionNotice();
  },
  computed: {
    notice() {
      const query = this.$route.query || {};

      if (query.registered === "1") {
        return {
          type: "success",
          title: "注册已提交",
          text: "账号创建成功，需等待管理员启用后才能登录。"
        };
      }

      if (query.reason === "kicked") {
        return {
          type: "warning",
          title: "会话已失效",
          text: "账号已在其他地方登录，请重新登录。"
        };
      }

      if (query.reason === "logout") {
        return {
          type: "success",
          title: "已退出登录",
          text: "你的账号已安全退出，如需继续使用请重新登录。"
        };
      }

      if (query.reason === "password_reset") {
        return {
          type: "warning",
          title: "请重新登录",
          text: "密码已修改，为保护账号安全请重新登录。"
        };
      }

      if (this.cachedSessionNotice && this.cachedSessionNotice.text) {
        return this.cachedSessionNotice;
      }

      return null;
    }
  },
  methods: {
    extractToken(result) {
      if (result && result.data && result.data.token) {
        return result.data.token;
      }

      if (result && result.token) {
        return result.token;
      }

      return "mock-user-token";
    },
    buildErrorMessage(error) {
      if (error && error.response && error.response.data && error.response.data.message) {
        return error.response.data.message;
      }

      return "用户登录接口当前为脚手架版本，请联通后端后使用";
    },
    submit() {
      if (!this.form.account || !this.form.password) {
        this.$message.warning("请完整填写登录信息");
        return;
      }

      userLogin(this.form)
        .then(
          function (result) {
            const token = this.extractToken(result);

            localStorage.setItem("bill_user_token", token);
            localStorage.setItem(
              "bill_user_profile",
              JSON.stringify({
                account: this.form.account
              })
            );
            this.$message.success("登录成功");
            this.$router.push("/user/ledger");
          }.bind(this)
        )
        .catch(
          function (error) {
            this.$message.error(this.buildErrorMessage(error));
          }.bind(this)
        );
    }
  }
};
</script>

<style scoped>
.auth-view {
  min-height: 100vh;
  display: grid;
  grid-template-columns: minmax(0, 1.1fr) 520px;
  background: linear-gradient(135deg, rgba(255, 250, 232, 0.96) 0%, rgba(245, 246, 248, 0.96) 100%);
}

.auth-view__hero,
.auth-view__panel {
  position: relative;
  padding: 48px;
}

.auth-view__hero {
  overflow: hidden;
  display: flex;
  align-items: center;
}

.auth-view__glow {
  position: absolute;
  border-radius: 999px;
  filter: blur(22px);
}

.auth-view__glow--primary {
  top: -120px;
  left: -80px;
  width: 320px;
  height: 320px;
  background: rgba(246, 211, 74, 0.42);
}

.auth-view__glow--soft {
  right: 40px;
  bottom: -80px;
  width: 260px;
  height: 260px;
  background: rgba(255, 255, 255, 0.88);
}

.auth-view__hero-content {
  position: relative;
  z-index: 1;
  max-width: 720px;
}

.auth-brand {
  display: flex;
  align-items: center;
  gap: 18px;
}

.auth-brand__mark {
  width: 72px;
  height: 72px;
  border-radius: 22px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: var(--brand-color);
  color: var(--text-main);
  font-size: 32px;
  font-weight: 800;
  box-shadow: var(--shadow-soft);
}

.auth-brand h1 {
  margin: 0;
  font-size: 48px;
  line-height: 1.1;
}

.auth-brand p {
  margin: 8px 0 0;
  color: var(--text-subtle);
  font-size: 18px;
}

.auth-copy {
  margin: 48px 0 28px;
}

.auth-copy h2 {
  margin: 0 0 16px;
  font-size: 40px;
  line-height: 1.15;
}

.auth-copy p {
  margin: 0;
  max-width: 560px;
  color: var(--text-subtle);
  line-height: 1.8;
  font-size: 16px;
}

.feature-list {
  display: grid;
  gap: 16px;
  max-width: 580px;
}

.feature-card {
  padding: 22px 24px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.72);
  border: 1px solid rgba(255, 255, 255, 0.92);
  box-shadow: var(--shadow-sm);
}

.feature-card strong {
  display: block;
  margin-bottom: 8px;
  font-size: 18px;
}

.feature-card p {
  margin: 0;
  color: var(--text-subtle);
  line-height: 1.7;
}

.auth-view__panel {
  display: flex;
  align-items: center;
  justify-content: center;
  border-left: 1px solid rgba(229, 231, 235, 0.86);
  background: rgba(255, 255, 255, 0.72);
}

.auth-panel {
  width: 100%;
  max-width: 408px;
  padding: 36px;
}

.auth-panel__eyebrow {
  display: inline-flex;
  align-items: center;
  min-height: 34px;
  padding: 0 14px;
  border-radius: 999px;
  background: var(--brand-soft);
  color: var(--text-main);
  font-size: 13px;
  font-weight: 700;
}

.auth-panel h3 {
  margin: 18px 0 8px;
  font-size: 34px;
}

.auth-panel__description {
  margin: 0;
  color: var(--text-subtle);
  line-height: 1.7;
}

.auth-notice {
  margin-top: 22px;
  padding: 14px 16px;
  border-radius: 18px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  border: 1px solid transparent;
}

.auth-notice strong {
  font-size: 14px;
}

.auth-notice span {
  color: var(--text-subtle);
  font-size: 13px;
  line-height: 1.6;
}

.auth-notice.is-success {
  background: rgba(34, 197, 94, 0.08);
  border-color: rgba(34, 197, 94, 0.18);
}

.auth-notice.is-warning {
  background: rgba(245, 158, 11, 0.08);
  border-color: rgba(245, 158, 11, 0.18);
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
  margin-top: 28px;
}

.auth-form__field {
  display: flex;
  flex-direction: column;
  gap: 8px;
  color: var(--text-main);
  font-weight: 600;
  font-size: 14px;
}

.auth-form__field input {
  width: 100%;
  height: 46px;
  padding: 0 16px;
  border-radius: 16px;
  border: 1px solid var(--border-color);
  background: var(--card-muted);
  color: var(--text-main);
}

.auth-form__field input:focus {
  outline: none;
  border-color: rgba(232, 191, 37, 0.9);
  box-shadow: 0 0 0 4px rgba(246, 211, 74, 0.18);
}

.auth-form__password {
  position: relative;
}

.auth-form__password input {
  padding-right: 72px;
}

.auth-form__toggle {
  position: absolute;
  top: 50%;
  right: 10px;
  transform: translateY(-50%);
  border: none;
  background: transparent;
  color: var(--text-subtle);
  font-size: 13px;
  font-weight: 600;
}

.auth-form__meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  color: var(--text-subtle);
  font-size: 13px;
}

.auth-checkbox {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.auth-checkbox input {
  width: 16px;
  height: 16px;
  accent-color: var(--brand-hover);
}

.auth-link {
  border: none;
  background: transparent;
  padding: 0;
  color: var(--text-main);
  font-weight: 600;
}

.auth-form__submit {
  height: 48px;
  border: none;
  border-radius: 16px;
  background: linear-gradient(135deg, var(--brand-color) 0%, var(--brand-hover) 100%);
  color: var(--text-main);
  font-weight: 700;
  box-shadow: var(--shadow-sm);
}

.auth-panel__footer {
  margin-top: 22px;
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-subtle);
  font-size: 14px;
}

@media (max-width: 1100px) {
  .auth-view {
    grid-template-columns: 1fr;
  }

  .auth-view__panel {
    border-left: none;
    border-top: 1px solid rgba(229, 231, 235, 0.86);
  }
}

@media (max-width: 640px) {
  .auth-view__hero,
  .auth-view__panel {
    padding: 24px;
  }

  .auth-brand h1 {
    font-size: 34px;
  }

  .auth-copy h2 {
    font-size: 28px;
  }
}
</style>
