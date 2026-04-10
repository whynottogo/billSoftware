<template>
  <div class="password-page">
    <section class="password-hero page-card">
      <span class="password-hero__tag">PASSWORD SECURITY</span>
      <h1>修改登录密码</h1>
      <p>本页已接入真实接口，提交后会执行原密码校验并在成功后清理登录态，重新回到登录页承接。</p>
    </section>

    <section class="password-panel page-card">
      <header class="password-panel__header">
        <h2>密码表单</h2>
        <p>请输入当前密码和新密码，修改成功后需要重新登录。</p>
      </header>

      <form class="password-form" @submit.prevent="submit">
        <label class="password-form__field">
          <span>原密码</span>
          <input
            v-model.trim="form.currentPassword"
            type="password"
            autocomplete="current-password"
            placeholder="请输入原密码"
          />
        </label>

        <label class="password-form__field">
          <span>新密码</span>
          <input
            v-model.trim="form.nextPassword"
            type="password"
            autocomplete="new-password"
            placeholder="请输入新密码（至少 6 位）"
          />
        </label>

        <label class="password-form__field">
          <span>确认新密码</span>
          <input
            v-model.trim="form.confirmPassword"
            type="password"
            autocomplete="new-password"
            placeholder="请再次输入新密码"
          />
        </label>

        <div class="password-form__tips">
          <article>
            <strong>规则 1</strong>
            <span>原密码、新密码、确认新密码均为必填。</span>
          </article>
          <article>
            <strong>规则 2</strong>
            <span>新密码和确认密码必须一致。</span>
          </article>
          <article>
            <strong>规则 3</strong>
            <span>修改成功后会清空登录态并跳回登录页。</span>
          </article>
        </div>

        <div class="password-form__actions">
          <button type="submit" class="finance-button finance-button--primary">提交修改</button>
          <button type="button" class="finance-button" @click="backProfile">返回个人信息</button>
        </div>
      </form>
    </section>
  </div>
</template>

<script>
import { buildUserProfileError, clearUserSession, updateUserPassword } from "@/api/userProfile";

export default {
  name: "UserProfilePassword",
  data() {
    return {
      form: {
        currentPassword: "",
        nextPassword: "",
        confirmPassword: ""
      },
      submitting: false
    };
  },
  methods: {
    backProfile() {
      this.$router.push("/user/profile");
    },
    submit() {
      if (!this.form.currentPassword || !this.form.nextPassword || !this.form.confirmPassword) {
        this.$message.warning("请完整填写原密码、新密码和确认密码");
        return;
      }

      if (this.form.nextPassword.length < 6) {
        this.$message.warning("新密码长度至少为 6 位");
        return;
      }

      if (this.form.nextPassword !== this.form.confirmPassword) {
        this.$message.warning("两次输入的新密码不一致");
        return;
      }

      if (this.submitting) {
        return;
      }

      this.submitting = true;

      updateUserPassword({
        current_password: this.form.currentPassword,
        new_password: this.form.nextPassword,
        confirm_password: this.form.confirmPassword,
        currentPassword: this.form.currentPassword,
        newPassword: this.form.nextPassword,
        confirmPassword: this.form.confirmPassword
      })
        .then(
          function() {
            clearUserSession("password_reset");
            this.$message.success("密码修改成功，请重新登录");
            this.$router.push({
              path: "/user/login",
              query: {
                reason: "password_reset"
              }
            });
          }.bind(this)
        )
        .catch(
          function(error) {
            this.$message.warning(buildUserProfileError(error, "密码修改失败，请稍后重试"));
          }.bind(this)
        )
        .finally(
          function() {
            this.submitting = false;
          }.bind(this)
        );
    }
  }
};
</script>

<style scoped>
.password-page {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.password-hero {
  padding: 26px 28px;
  background:
    radial-gradient(circle at 85% 0%, rgba(255, 255, 255, 0.48) 0%, rgba(255, 255, 255, 0) 38%),
    linear-gradient(135deg, rgba(252, 242, 207, 0.96) 0%, rgba(255, 255, 255, 0.96) 70%);
}

.password-hero__tag {
  display: inline-flex;
  min-height: 32px;
  align-items: center;
  padding: 0 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.78);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.06em;
}

.password-hero h1 {
  margin: 12px 0 0;
  font-size: 32px;
}

.password-hero p {
  margin: 10px 0 0;
  color: var(--text-subtle);
  line-height: 1.7;
}

.password-panel {
  padding: 24px;
}

.password-panel__header h2,
.password-panel__header p {
  margin: 0;
}

.password-panel__header p {
  margin-top: 8px;
  color: var(--text-muted);
}

.password-form {
  margin-top: 18px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.password-form__field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.password-form__field span {
  font-size: 13px;
  font-weight: 700;
  color: var(--text-subtle);
}

.password-form__field input {
  min-height: 48px;
  border-radius: 14px;
  border: 1px solid var(--border-color);
  background: #fff;
  padding: 0 14px;
}

.password-form__tips {
  margin-top: 4px;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.password-form__tips article {
  border-radius: 14px;
  border: 1px solid rgba(246, 211, 74, 0.45);
  background: rgba(255, 249, 229, 0.76);
  padding: 12px;
}

.password-form__tips strong,
.password-form__tips span {
  display: block;
}

.password-form__tips span {
  margin-top: 7px;
  color: var(--text-subtle);
  font-size: 13px;
  line-height: 1.55;
}

.password-form__actions {
  margin-top: 6px;
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

@media (max-width: 980px) {
  .password-form__tips {
    grid-template-columns: 1fr;
  }
}
</style>
