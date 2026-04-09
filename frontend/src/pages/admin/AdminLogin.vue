<template>
  <div class="admin-auth">
    <section class="admin-auth__hero">
      <div class="admin-auth__brand">
        <span class="admin-auth__mark">AD</span>
        <div>
          <h1>Bill Software</h1>
          <p>管理员后台入口</p>
        </div>
      </div>

      <div class="admin-auth__copy">
        <h2>更克制，但更清楚。</h2>
        <p>
          管理端只服务唯一管理员账号，用来查看用户基础信息、确认启用状态，
          并在后续批次继续承接用户详情与账单数据查看。
        </p>
      </div>

      <div class="admin-auth__highlights">
        <article>
          <strong>用户管理</strong>
          <span>查看账号、昵称、手机号、邮箱与启用状态。</span>
        </article>
        <article>
          <strong>认证独立</strong>
          <span>管理员登录只读取管理端认证逻辑，不回落到用户体系。</span>
        </article>
        <article>
          <strong>入口分离</strong>
          <span>用户端与管理端保持不同路由前缀和不同存储 key。</span>
        </article>
      </div>
    </section>

    <section class="admin-auth__panel">
      <div class="admin-auth__card page-card">
        <span class="admin-auth__eyebrow">管理端登录</span>
        <h3>管理员登录</h3>
        <p class="admin-auth__description">请使用固定管理员账号进入后台。</p>

        <form class="admin-auth__form" @submit.prevent="submit">
          <label class="admin-auth__field">
            <span>管理员账号</span>
            <input
              v-model.trim="form.username"
              type="text"
              autocomplete="username"
              placeholder="请输入管理员账号"
            />
          </label>

          <label class="admin-auth__field">
            <span>密码</span>
            <div class="admin-auth__password">
              <input
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                placeholder="请输入密码"
              />
              <button type="button" class="admin-auth__toggle" @click="showPassword = !showPassword">
                {{ showPassword ? "隐藏" : "显示" }}
              </button>
            </div>
          </label>

          <button type="submit" class="admin-auth__submit">登录后台</button>
        </form>

        <div class="admin-auth__footer">
          <span>需要返回普通用户入口？</span>
          <button type="button" class="admin-auth__link" @click="$router.push('/user/login')">
            回到用户登录
          </button>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import { adminLogin } from "@/api/userAuth";

export default {
  name: "AdminLogin",
  data() {
    return {
      form: {
        username: "",
        password: ""
      },
      showPassword: false
    };
  },
  methods: {
    extractToken(result) {
      if (result && result.data && result.data.token) {
        return result.data.token;
      }

      if (result && result.token) {
        return result.token;
      }

      return "mock-admin-token";
    },
    buildErrorMessage(error) {
      if (error && error.response && error.response.data && error.response.data.message) {
        return error.response.data.message;
      }

      return "管理员登录接口当前为脚手架版本，请联通后端后使用";
    },
    submit() {
      if (!this.form.username || !this.form.password) {
        this.$message.warning("请完整填写管理员登录信息");
        return;
      }

      adminLogin(this.form)
        .then(
          function (result) {
            const token = this.extractToken(result);

            localStorage.setItem("bill_admin_token", token);
            this.$message.success("管理员登录成功");
            this.$router.push("/admin/dashboard");
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
.admin-auth {
  min-height: 100vh;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 500px;
  background: linear-gradient(135deg, rgba(246, 246, 246, 0.98) 0%, rgba(255, 255, 255, 0.92) 100%);
}

.admin-auth__hero,
.admin-auth__panel {
  padding: 48px;
}

.admin-auth__hero {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 24px;
}

.admin-auth__brand {
  display: flex;
  align-items: center;
  gap: 16px;
}

.admin-auth__mark {
  width: 68px;
  height: 68px;
  border-radius: 22px;
  background: #171717;
  color: #ffffff;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 800;
  letter-spacing: 0.08em;
  box-shadow: var(--shadow-soft);
}

.admin-auth__brand h1 {
  margin: 0;
  font-size: 42px;
}

.admin-auth__brand p {
  margin: 8px 0 0;
  color: var(--text-subtle);
  font-size: 18px;
}

.admin-auth__copy h2 {
  margin: 0 0 14px;
  font-size: 38px;
  line-height: 1.15;
}

.admin-auth__copy p {
  margin: 0;
  max-width: 540px;
  color: var(--text-subtle);
  line-height: 1.8;
}

.admin-auth__highlights {
  display: grid;
  gap: 14px;
  max-width: 560px;
}

.admin-auth__highlights article {
  padding: 20px 22px;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
}

.admin-auth__highlights strong {
  display: block;
  margin-bottom: 6px;
  font-size: 17px;
}

.admin-auth__highlights span {
  color: var(--text-subtle);
  line-height: 1.7;
  font-size: 14px;
}

.admin-auth__panel {
  display: flex;
  align-items: center;
  justify-content: center;
  border-left: 1px solid rgba(229, 231, 235, 0.86);
  background: rgba(255, 255, 255, 0.88);
}

.admin-auth__card {
  width: 100%;
  max-width: 390px;
  padding: 34px;
}

.admin-auth__eyebrow {
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

.admin-auth__card h3 {
  margin: 18px 0 8px;
  font-size: 34px;
}

.admin-auth__description {
  margin: 0;
  color: var(--text-subtle);
}

.admin-auth__form {
  display: flex;
  flex-direction: column;
  gap: 18px;
  margin-top: 28px;
}

.admin-auth__field {
  display: flex;
  flex-direction: column;
  gap: 8px;
  color: var(--text-main);
  font-weight: 600;
  font-size: 14px;
}

.admin-auth__field input {
  width: 100%;
  height: 46px;
  padding: 0 16px;
  border-radius: 16px;
  border: 1px solid var(--border-color);
  background: var(--card-muted);
  color: var(--text-main);
}

.admin-auth__field input:focus {
  outline: none;
  border-color: rgba(23, 23, 23, 0.56);
  box-shadow: 0 0 0 4px rgba(23, 23, 23, 0.08);
}

.admin-auth__password {
  position: relative;
}

.admin-auth__password input {
  padding-right: 72px;
}

.admin-auth__toggle {
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

.admin-auth__submit {
  height: 48px;
  border: none;
  border-radius: 16px;
  background: #171717;
  color: #ffffff;
  font-weight: 700;
  box-shadow: var(--shadow-sm);
}

.admin-auth__footer {
  margin-top: 22px;
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-subtle);
  font-size: 14px;
}

.admin-auth__link {
  border: none;
  background: transparent;
  padding: 0;
  color: var(--text-main);
  font-weight: 600;
}

@media (max-width: 1080px) {
  .admin-auth {
    grid-template-columns: 1fr;
  }

  .admin-auth__panel {
    border-left: none;
    border-top: 1px solid rgba(229, 231, 235, 0.86);
  }
}

@media (max-width: 640px) {
  .admin-auth__hero,
  .admin-auth__panel {
    padding: 24px;
  }

  .admin-auth__brand h1 {
    font-size: 32px;
  }

  .admin-auth__copy h2 {
    font-size: 28px;
  }
}
</style>
