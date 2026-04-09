<template>
  <div class="profile-page">
    <section class="profile-hero page-card">
      <div class="profile-hero__badge">
        <span>资料中心</span>
        <strong>USER PROFILE</strong>
      </div>
      <h1>把你的账号信息维护清楚</h1>
      <p>
        账户与手机号用于身份识别，保持只读。昵称、邮箱和头像可随时更新，当前版本使用本地 mock
        状态保存，后续可直接替换为真实接口。
      </p>
      <div class="profile-hero__actions">
        <button class="finance-button finance-button--primary" @click="goPasswordPage">修改密码</button>
        <button class="finance-button" @click="triggerAvatarUpload">更新头像</button>
      </div>
    </section>

    <section class="profile-content">
      <article class="profile-panel page-card">
        <div class="profile-panel__header">
          <h2>头像与基础信息</h2>
          <span>支持头像占位上传和资料编辑</span>
        </div>

        <div class="profile-avatar">
          <div class="profile-avatar__preview">
            <img v-if="profile.avatar" :src="profile.avatar" alt="用户头像" />
            <span v-else>{{ profileInitial }}</span>
          </div>
          <div class="profile-avatar__meta">
            <strong>{{ profile.nickname }}</strong>
            <p>账号：{{ profile.account }}</p>
            <div class="profile-avatar__buttons">
              <button class="finance-button finance-button--ghost" @click="triggerAvatarUpload">上传头像</button>
              <button class="finance-button finance-button--ghost" @click="resetAvatar">恢复默认</button>
            </div>
            <small>{{ avatarHint }}</small>
          </div>
        </div>

        <input
          ref="avatarInput"
          class="profile-avatar__input"
          type="file"
          accept="image/*"
          @change="onAvatarChange"
        />

        <div class="profile-form-grid">
          <label class="profile-form__field">
            <span>账号</span>
            <input :value="profile.account" type="text" readonly />
            <small>账号仅展示，不可修改</small>
          </label>

          <label class="profile-form__field">
            <span>手机号</span>
            <input :value="profile.phone" type="text" readonly />
            <small>手机号用于登录识别，不支持修改</small>
          </label>

          <label class="profile-form__field">
            <span>昵称</span>
            <input v-model.trim="profile.nickname" type="text" placeholder="请输入昵称" />
            <small>昵称会同步展示在页面欢迎区</small>
          </label>

          <label class="profile-form__field">
            <span>邮箱</span>
            <input v-model.trim="profile.email" type="email" placeholder="请输入邮箱" />
            <small>用于消息触达和账号找回</small>
          </label>
        </div>

        <div class="profile-panel__footer">
          <button class="finance-button finance-button--primary" @click="saveProfile">保存资料</button>
          <button class="finance-button" @click="restoreProfile">还原本次修改</button>
        </div>
      </article>

      <article class="session-panel page-card">
        <div class="session-panel__header">
          <h3>会话安全</h3>
          <p>统一管理主动退出与会话失效提示文案。</p>
        </div>
        <div class="session-panel__tip">
          <strong>统一提示文案</strong>
          <span>账号已在其他地方登录，请重新登录。</span>
        </div>
        <div class="session-panel__actions">
          <button class="finance-button finance-button--primary" @click="logoutNow">主动退出登录</button>
          <button class="finance-button" @click="simulateKickout">模拟被挤下线提示</button>
        </div>
      </article>
    </section>
  </div>
</template>

<script>
import {
  getUserProfileMock,
  saveUserProfileMock,
  clearUserSession,
  createAvatarPlaceholder
} from "@/utils/userProfileMock";

const FALLBACK_AVATAR =
  "https://images.unsplash.com/photo-1527980965255-d3b416303d12?auto=format&fit=crop&w=240&q=80";

export default {
  name: "UserProfile",
  data() {
    return {
      profile: {
        account: "",
        username: "",
        nickname: "",
        phone: "",
        email: "",
        avatar: "",
        avatarCompressed: "",
        avatarUpdatedAt: ""
      },
      initialSnapshot: null,
      avatarHint: "头像大于 2M 时会显示“已模拟压缩处理”的占位提示。"
    };
  },
  computed: {
    profileInitial() {
      const nickname = (this.profile.nickname || "").trim();
      const account = (this.profile.account || "").trim();
      const seed = nickname || account || "U";
      return seed.slice(0, 1).toUpperCase();
    }
  },
  created() {
    this.restoreFromStorage();
  },
  methods: {
    restoreFromStorage() {
      const saved = getUserProfileMock();
      this.profile = Object.assign({}, saved);
      this.initialSnapshot = Object.assign({}, saved);
    },
    triggerAvatarUpload() {
      this.$refs.avatarInput.click();
    },
    onAvatarChange(event) {
      const file = event.target.files && event.target.files[0];

      if (!file) {
        return;
      }

      const reader = new FileReader();
      reader.onload = (loadEvent) => {
        const previewUrl = loadEvent.target && loadEvent.target.result ? String(loadEvent.target.result) : "";

        if (!previewUrl) {
          this.$message.error("头像读取失败，请重试");
          return;
        }

        const isLargeFile = file.size > 2 * 1024 * 1024;
        this.profile.avatar = previewUrl;
        this.profile.avatarCompressed = isLargeFile ? previewUrl : "";
        this.profile.avatarUpdatedAt = new Date().toISOString();
        this.avatarHint = isLargeFile
          ? `已为 ${createAvatarPlaceholder(file.name)} 执行模拟压缩处理（前端占位）`
          : `已更新头像：${createAvatarPlaceholder(file.name)}`;
        this.$message.success("头像占位已更新");
      };
      reader.readAsDataURL(file);
    },
    resetAvatar() {
      this.profile.avatar = FALLBACK_AVATAR;
      this.profile.avatarCompressed = "";
      this.profile.avatarUpdatedAt = "";
      this.avatarHint = "已恢复默认头像，可重新上传。";
      this.$message.success("头像已恢复默认");
    },
    saveProfile() {
      if (!this.profile.nickname || !this.profile.email) {
        this.$message.warning("昵称和邮箱不能为空");
        return;
      }

      const saved = saveUserProfileMock(this.profile);
      this.profile = Object.assign({}, saved);
      this.initialSnapshot = Object.assign({}, saved);
      this.$message.success("资料已保存（mock）");
    },
    restoreProfile() {
      if (!this.initialSnapshot) {
        this.restoreFromStorage();
        return;
      }

      this.profile = Object.assign({}, this.initialSnapshot);
      this.avatarHint = "已还原为最近一次保存的资料。";
      this.$message.success("已还原本次修改");
    },
    goPasswordPage() {
      this.$router.push("/user/profile/password");
    },
    logoutNow() {
      clearUserSession("logout");
      this.$router.push({
        path: "/user/login",
        query: {
          reason: "logout"
        }
      });
    },
    simulateKickout() {
      clearUserSession("kicked");
      this.$router.push({
        path: "/user/session-kickout",
        query: {
          reason: "kicked",
          from: "profile"
        }
      });
    }
  }
};
</script>

<style scoped>
.profile-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.profile-hero {
  padding: 30px;
  background:
    radial-gradient(circle at right top, rgba(255, 255, 255, 0.64) 0%, rgba(255, 255, 255, 0) 46%),
    linear-gradient(140deg, rgba(246, 211, 74, 0.94) 0%, rgba(255, 248, 220, 0.92) 54%, rgba(255, 255, 255, 0.95) 100%);
  border: 1px solid rgba(246, 211, 74, 0.4);
}

.profile-hero__badge {
  width: fit-content;
  display: inline-flex;
  align-items: center;
  gap: 12px;
  min-height: 34px;
  padding: 0 14px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.76);
  font-size: 12px;
  font-weight: 700;
}

.profile-hero__badge strong {
  font-size: 11px;
  letter-spacing: 0.08em;
}

.profile-hero h1 {
  margin: 16px 0 0;
  font-size: 34px;
  line-height: 1.2;
}

.profile-hero p {
  margin: 12px 0 0;
  max-width: 760px;
  line-height: 1.7;
  color: rgba(23, 23, 23, 0.74);
}

.profile-hero__actions {
  margin-top: 20px;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.profile-content {
  display: grid;
  grid-template-columns: minmax(0, 1.4fr) minmax(0, 1fr);
  gap: 18px;
}

.profile-panel,
.session-panel {
  padding: 24px;
}

.profile-panel__header h2,
.profile-panel__header span,
.session-panel__header h3,
.session-panel__header p {
  margin: 0;
}

.profile-panel__header span,
.session-panel__header p {
  display: block;
  margin-top: 8px;
  color: var(--text-muted);
  line-height: 1.6;
}

.profile-avatar {
  margin-top: 18px;
  display: flex;
  gap: 18px;
  align-items: center;
  padding: 16px;
  border-radius: 20px;
  background: rgba(255, 250, 233, 0.68);
  border: 1px solid rgba(246, 211, 74, 0.32);
}

.profile-avatar__preview {
  width: 78px;
  height: 78px;
  border-radius: 24px;
  overflow: hidden;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: rgba(246, 211, 74, 0.44);
  color: var(--text-main);
  font-size: 28px;
  font-weight: 800;
}

.profile-avatar__preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.profile-avatar__meta strong {
  display: block;
  font-size: 18px;
}

.profile-avatar__meta p {
  margin: 6px 0 0;
  color: var(--text-subtle);
}

.profile-avatar__buttons {
  margin-top: 12px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.profile-avatar__meta small {
  display: block;
  margin-top: 10px;
  color: var(--text-muted);
  line-height: 1.5;
}

.profile-avatar__input {
  display: none;
}

.profile-form-grid {
  margin-top: 18px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.profile-form__field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.profile-form__field span {
  font-size: 13px;
  font-weight: 700;
  color: var(--text-subtle);
}

.profile-form__field input {
  min-height: 46px;
  border-radius: 14px;
  border: 1px solid var(--border-color);
  padding: 0 14px;
  background: #fff;
}

.profile-form__field input[readonly] {
  background: rgba(249, 250, 251, 0.9);
  color: var(--text-muted);
}

.profile-form__field small {
  font-size: 12px;
  color: var(--text-muted);
}

.profile-panel__footer {
  margin-top: 20px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.session-panel__tip {
  margin-top: 16px;
  padding: 14px 16px;
  border-radius: 16px;
  background: rgba(255, 244, 214, 0.72);
  border: 1px solid rgba(246, 211, 74, 0.42);
}

.session-panel__tip strong,
.session-panel__tip span {
  display: block;
}

.session-panel__tip span {
  margin-top: 8px;
  color: var(--text-subtle);
  line-height: 1.6;
}

.session-panel__actions {
  margin-top: 18px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

@media (max-width: 1080px) {
  .profile-content {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 860px) {
  .profile-hero {
    padding: 24px;
  }

  .profile-hero h1 {
    font-size: 28px;
  }

  .profile-form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
