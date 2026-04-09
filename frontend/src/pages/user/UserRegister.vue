<template>
  <div class="register-view">
    <section class="register-view__hero">
      <div class="register-view__badge">用户注册</div>
      <h1>创建你的记账账号</h1>
      <p>
        完成基础信息后，账号会先进入“等待管理员启用”的状态。
        首批版本不自动登录，启用后请回到用户端登录页进入系统。
      </p>

      <div class="register-view__tips">
        <article>
          <strong>账号体系独立</strong>
          <span>普通用户账号与管理员账号完全分离，避免串端登录。</span>
        </article>
        <article>
          <strong>字段一次填全</strong>
          <span>用户名、昵称、手机号、邮箱、密码五项信息均为首批必填。</span>
        </article>
        <article>
          <strong>启用后再登录</strong>
          <span>注册成功后需管理员启用，启用前无法进入用户端首页。</span>
        </article>
      </div>
    </section>

    <section class="register-view__panel">
      <div class="register-card page-card">
        <h2>填写注册信息</h2>
        <p class="register-card__description">所有字段均与需求文档保持一致。</p>

        <form class="register-form" @submit.prevent="submit">
          <label class="register-form__field">
            <span>用户名</span>
            <input v-model.trim="form.username" type="text" autocomplete="username" placeholder="请输入用户名" />
          </label>
          <label class="register-form__field">
            <span>昵称</span>
            <input v-model.trim="form.nickname" type="text" autocomplete="name" placeholder="请输入昵称" />
          </label>
          <label class="register-form__field">
            <span>手机号</span>
            <input v-model.trim="form.phone" type="text" autocomplete="tel" placeholder="请输入手机号" />
          </label>
          <label class="register-form__field">
            <span>邮箱</span>
            <input v-model.trim="form.email" type="email" autocomplete="email" placeholder="请输入邮箱" />
          </label>
          <label class="register-form__field register-form__field--full">
            <span>密码</span>
            <div class="register-form__password">
              <input
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="new-password"
                placeholder="请输入密码"
              />
              <button type="button" class="register-form__toggle" @click="showPassword = !showPassword">
                {{ showPassword ? "隐藏" : "显示" }}
              </button>
            </div>
          </label>

          <div class="register-card__notice">
            <strong>注册规则</strong>
            <span>注册成功后不自动登录，管理员启用前无法进入用户端。</span>
          </div>

          <button type="submit" class="register-form__submit">提交注册</button>
        </form>

        <div class="register-card__footer">
          <span>已经有账户？</span>
          <button type="button" class="register-link" @click="$router.push('/user/login')">返回登录</button>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import { userRegister } from "@/api/userAuth";

export default {
  name: "UserRegister",
  data() {
    return {
      form: {
        username: "",
        nickname: "",
        phone: "",
        email: "",
        password: ""
      },
      showPassword: false
    };
  },
  methods: {
    buildErrorMessage(error) {
      if (error && error.response && error.response.data && error.response.data.message) {
        return error.response.data.message;
      }

      return "用户注册接口当前为脚手架版本，请联通后端后使用";
    },
    submit() {
      if (!this.form.username || !this.form.nickname || !this.form.phone || !this.form.email || !this.form.password) {
        this.$message.warning("请完整填写注册信息");
        return;
      }

      userRegister(this.form)
        .then(
          function () {
            this.$message.success("注册请求已提交，等待管理员启用");
            this.$router.push({
              path: "/user/login",
              query: {
                registered: "1"
              }
            });
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
.register-view {
  min-height: 100vh;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 620px;
  background: linear-gradient(180deg, rgba(255, 248, 220, 0.92) 0%, rgba(245, 246, 248, 0.96) 100%);
}

.register-view__hero,
.register-view__panel {
  padding: 48px;
}

.register-view__hero {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 18px;
}

.register-view__badge {
  display: inline-flex;
  align-items: center;
  min-height: 34px;
  padding: 0 14px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.72);
  color: var(--text-main);
  font-size: 13px;
  font-weight: 700;
  width: fit-content;
}

.register-view__hero h1 {
  margin: 0;
  font-size: 44px;
  line-height: 1.12;
}

.register-view__hero p {
  margin: 0;
  max-width: 520px;
  color: var(--text-subtle);
  line-height: 1.8;
}

.register-view__tips {
  display: grid;
  gap: 14px;
  margin-top: 18px;
  max-width: 540px;
}

.register-view__tips article {
  padding: 20px 22px;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.76);
  border: 1px solid rgba(255, 255, 255, 0.9);
  box-shadow: var(--shadow-sm);
}

.register-view__tips strong {
  display: block;
  margin-bottom: 6px;
  font-size: 17px;
}

.register-view__tips span {
  color: var(--text-subtle);
  line-height: 1.7;
  font-size: 14px;
}

.register-view__panel {
  display: flex;
  align-items: center;
  justify-content: center;
  border-left: 1px solid rgba(229, 231, 235, 0.9);
  background: rgba(255, 255, 255, 0.72);
}

.register-card {
  width: 100%;
  padding: 34px;
}

.register-card h2 {
  margin: 0;
  font-size: 34px;
}

.register-card__description {
  margin: 10px 0 0;
  color: var(--text-subtle);
}

.register-form {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px 16px;
  margin-top: 26px;
}

.register-form__field {
  display: flex;
  flex-direction: column;
  gap: 8px;
  color: var(--text-main);
  font-weight: 600;
  font-size: 14px;
}

.register-form__field input {
  width: 100%;
  height: 46px;
  padding: 0 16px;
  border-radius: 16px;
  border: 1px solid var(--border-color);
  background: var(--card-muted);
  color: var(--text-main);
}

.register-form__field input:focus {
  outline: none;
  border-color: rgba(232, 191, 37, 0.9);
  box-shadow: 0 0 0 4px rgba(246, 211, 74, 0.18);
}

.register-form__field--full {
  grid-column: 1 / -1;
}

.register-form__password {
  position: relative;
}

.register-form__password input {
  padding-right: 72px;
}

.register-form__toggle {
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

.register-card__notice {
  grid-column: 1 / -1;
  padding: 16px 18px;
  border-radius: 18px;
  background: rgba(246, 211, 74, 0.12);
  border: 1px solid rgba(246, 211, 74, 0.3);
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.register-card__notice strong {
  font-size: 14px;
}

.register-card__notice span {
  color: var(--text-subtle);
  font-size: 13px;
  line-height: 1.6;
}

.register-form__submit {
  grid-column: 1 / -1;
  height: 48px;
  border: none;
  border-radius: 16px;
  background: linear-gradient(135deg, var(--brand-color) 0%, var(--brand-hover) 100%);
  color: var(--text-main);
  font-weight: 700;
  box-shadow: var(--shadow-sm);
}

.register-card__footer {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 22px;
  color: var(--text-subtle);
}

.register-link {
  border: none;
  background: transparent;
  padding: 0;
  color: var(--text-main);
  font-weight: 600;
}

@media (max-width: 1100px) {
  .register-view {
    grid-template-columns: 1fr;
  }

  .register-view__panel {
    border-left: none;
    border-top: 1px solid rgba(229, 231, 235, 0.9);
  }
}

@media (max-width: 720px) {
  .register-view__hero,
  .register-view__panel {
    padding: 24px;
  }

  .register-form {
    grid-template-columns: 1fr;
  }
}
</style>
