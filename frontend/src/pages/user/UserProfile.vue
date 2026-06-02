<template>
  <div class="profile-page">
    <section class="profile-hero page-card">
      <div class="profile-hero__badge">
        <span>资料中心</span>
        <strong>USER PROFILE</strong>
      </div>
      <h1>把你的账号信息维护清楚</h1>
      <p>
        账户与手机号用于身份识别，保持只读。昵称、邮箱和头像可随时更新，当前版本已接入真实接口并
        与当前登录账号同步。
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
          <span>支持头像上传、压缩预览和资料编辑</span>
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
              <button class="finance-button finance-button--ghost" :disabled="savingProfile" @click="triggerAvatarUpload">上传头像</button>
              <button class="finance-button finance-button--ghost" :disabled="savingProfile" @click="resetAvatar">恢复默认</button>
            </div>
            <small>{{ avatarHint }}</small>
          </div>
        </div>

        <input
          ref="avatarInput"
          class="profile-avatar__input"
          type="file"
          accept="image/jpeg,image/png,image/webp"
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
          <button class="finance-button finance-button--primary" :disabled="savingProfile" @click="saveProfile">
            {{ savingProfile ? "保存中..." : "保存资料" }}
          </button>
          <button class="finance-button" :disabled="savingProfile" @click="restoreProfile">还原本次修改</button>
        </div>
      </article>
    </section>
  </div>
</template>

<script>
import {
  buildUserProfileError,
  getUserProfile,
  normalizeUserProfilePayload,
  persistUserSessionProfile,
  updateUserProfile
} from "@/api/userProfile";

const MAX_COMPRESSED_AVATAR_BYTES = 2 * 1024 * 1024;
const ALLOWED_AVATAR_TYPES = ["image/jpeg", "image/png", "image/webp"];

function cloneProfile(profile) {
  return Object.assign({}, profile || {});
}

function createAvatarPlaceholder(fileName) {
  if (!fileName) {
    return "未命名头像";
  }

  return fileName.replace(/\.[^.]+$/, "");
}

function readFileAsDataUrl(file) {
  return new Promise(function(resolve, reject) {
    const reader = new FileReader();

    reader.onload = function(loadEvent) {
      const result = loadEvent && loadEvent.target ? loadEvent.target.result : "";

      if (!result) {
        reject(new Error("头像读取失败"));
        return;
      }

      resolve(String(result));
    };
    reader.onerror = function() {
      reject(new Error("头像读取失败"));
    };
    reader.readAsDataURL(file);
  });
}

function loadImage(source) {
  return new Promise(function(resolve, reject) {
    const image = new Image();

    image.onload = function() {
      resolve(image);
    };
    image.onerror = function() {
      reject(new Error("头像解析失败"));
    };
    image.src = source;
  });
}

function estimateDataUrlBytes(dataUrl) {
  const base64Body = String(dataUrl || "").split(",")[1] || "";
  return Math.ceil((base64Body.length * 3) / 4);
}

function canvasToDataUrl(canvas, quality) {
  return canvas.toDataURL("image/jpeg", quality);
}

function isAllowedAvatar(file) {
  return ALLOWED_AVATAR_TYPES.indexOf(file.type) !== -1;
}

function buildAvatarPayload(file) {
  if (!isAllowedAvatar(file)) {
    return Promise.reject(new Error("仅支持 JPG、PNG、WebP 格式头像"));
  }

  return readFileAsDataUrl(file).then(function(originalDataUrl) {
    if (file.size <= 2 * 1024 * 1024) {
      return {
        previewUrl: originalDataUrl,
        originalDataUrl: originalDataUrl,
        compressedDataUrl: originalDataUrl,
        compressed: false
      };
    }

    return loadImage(originalDataUrl).then(function(image) {
      const canvas = document.createElement("canvas");
      const context = canvas.getContext("2d");

      if (!context) {
        throw new Error("当前浏览器不支持头像压缩，请更换 2M 以内图片");
      }

      let scale = 1;
      let quality = 0.88;
      let compressedDataUrl = originalDataUrl;

      for (let index = 0; index < 12; index += 1) {
        canvas.width = Math.max(1, Math.round(image.naturalWidth * scale));
        canvas.height = Math.max(1, Math.round(image.naturalHeight * scale));
        context.clearRect(0, 0, canvas.width, canvas.height);
        context.drawImage(image, 0, 0, canvas.width, canvas.height);
        compressedDataUrl = canvasToDataUrl(canvas, quality);

        if (estimateDataUrlBytes(compressedDataUrl) <= MAX_COMPRESSED_AVATAR_BYTES) {
          break;
        }

        scale *= 0.78;
        quality = Math.max(0.38, quality - 0.08);
      }

      if (estimateDataUrlBytes(compressedDataUrl) > MAX_COMPRESSED_AVATAR_BYTES) {
        throw new Error("头像压缩后仍超过 2M，请更换图片");
      }

      return {
        previewUrl: compressedDataUrl,
        originalDataUrl: originalDataUrl,
        compressedDataUrl: compressedDataUrl,
        compressed: true
      };
    });
  });
}

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
        avatarOriginal: "",
        avatarCompressed: "",
        avatarUpdatedAt: ""
      },
      initialSnapshot: null,
      avatarHint: "头像大于 2M 时会自动压缩后再保存。",
      loadingProfile: false,
      savingProfile: false
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
    this.fetchProfile();
  },
  methods: {
    syncSessionProfile(profile) {
      persistUserSessionProfile(profile);
    },
    applyProfile(profile) {
      const normalizedProfile = cloneProfile(profile);

      this.profile = normalizedProfile;
      this.initialSnapshot = cloneProfile(normalizedProfile);
      this.syncSessionProfile(normalizedProfile);
      this.avatarHint = normalizedProfile.avatarCompressed
        ? "当前头像已包含压缩版本，保存后会同步更新。"
        : "资料已从真实接口读取，可继续编辑昵称、邮箱和头像。";
    },
    fetchProfile() {
      this.loadingProfile = true;

      getUserProfile()
        .then(
          function(result) {
            const normalizedProfile = normalizeUserProfilePayload(result, this.profile);
            this.applyProfile(normalizedProfile);
          }.bind(this)
        )
        .catch(
          function(error) {
            this.$message.error(buildUserProfileError(error, "个人信息加载失败，请稍后重试"));
          }.bind(this)
        )
        .finally(
          function() {
            this.loadingProfile = false;
          }.bind(this)
        );
    },
    triggerAvatarUpload() {
      if (this.$refs.avatarInput) {
        this.$refs.avatarInput.click();
      }
    },
    onAvatarChange(event) {
      const file = event.target.files && event.target.files[0];

      if (!file) {
        return;
      }

      this.avatarHint = "正在处理头像，请稍候...";

      buildAvatarPayload(file)
        .then(
          function(avatarPayload) {
            this.profile.avatar = avatarPayload.previewUrl;
            this.profile.avatarOriginal = avatarPayload.originalDataUrl;
            this.profile.avatarCompressed = avatarPayload.compressedDataUrl;
            this.profile.avatarUpdatedAt = new Date().toISOString();
            this.avatarHint = avatarPayload.compressed
              ? `已为 ${createAvatarPlaceholder(file.name)} 生成压缩头像，保存后将同步到真实接口`
              : `已更新头像：${createAvatarPlaceholder(file.name)}，保存后生效`;
            this.$message.success("头像已处理，请点击保存资料生效");
          }.bind(this)
        )
        .catch(
          function(error) {
            this.avatarHint = "头像处理失败，请重新选择 JPG、PNG 或 WebP 图片。";
            this.$message.error(error && error.message ? error.message : "头像读取失败，请重试");
          }.bind(this)
        )
        .finally(function() {
          if (event.target) {
            event.target.value = "";
          }
        });
    },
    resetAvatar() {
      this.profile.avatar = "";
      this.profile.avatarOriginal = "";
      this.profile.avatarCompressed = "";
      this.profile.avatarUpdatedAt = new Date().toISOString();
      this.avatarHint = "已恢复为默认头像状态，保存后生效。";
      this.$message.success("头像已恢复默认状态");
    },
    saveProfile() {
      if (!this.profile.nickname || !this.profile.email) {
        this.$message.warning("昵称和邮箱不能为空");
        return;
      }

      if (this.savingProfile) {
        return;
      }

      this.savingProfile = true;

      updateUserProfile({
        nickname: this.profile.nickname,
        email: this.profile.email,
        avatar: this.profile.avatar || "",
        avatar_original: this.profile.avatarOriginal || this.profile.avatar || "",
        avatar_compressed: this.profile.avatarCompressed || "",
        avatarOriginal: this.profile.avatarOriginal || this.profile.avatar || "",
        avatarCompressed: this.profile.avatarCompressed || ""
      })
        .then(
          function(result) {
            const normalizedProfile = normalizeUserProfilePayload(result, this.profile);
            this.applyProfile(normalizedProfile);
            this.$message.success("资料已保存");
          }.bind(this)
        )
        .catch(
          function(error) {
            this.$message.error(buildUserProfileError(error, "资料保存失败，请稍后重试"));
          }.bind(this)
        )
        .finally(
          function() {
            this.savingProfile = false;
          }.bind(this)
        );
    },
    restoreProfile() {
      if (!this.initialSnapshot) {
        this.fetchProfile();
        return;
      }

      this.profile = cloneProfile(this.initialSnapshot);
      this.avatarHint = "已还原为最近一次保存的资料。";
      this.$message.success("已还原本次修改");
    },
    goPasswordPage() {
      this.$router.push("/user/profile/password");
    }
  }
};
</script>

<style scoped>
.profile-page {
  display: flex;
  flex-direction: column;
  gap: 14px;
  min-height: calc(100vh - 56px);
}

.profile-hero {
  padding: 20px 24px;
  background:
    radial-gradient(circle at right top, rgba(255, 255, 255, 0.64) 0%, rgba(255, 255, 255, 0) 46%),
    linear-gradient(140deg, rgba(246, 211, 74, 0.94) 0%, rgba(255, 248, 220, 0.92) 54%, rgba(255, 255, 255, 0.95) 100%);
  border: 1px solid rgba(246, 211, 74, 0.4);
}

.profile-content {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  gap: 14px;
  flex: 1;
}

.profile-panel {
  padding: 18px;
  height: 100%;
}

.profile-hero {
  padding: 20px 24px;
  background:
    radial-gradient(circle at right top, rgba(255, 255, 255, 0.64) 0%, rgba(255, 255, 255, 0) 46%),
    linear-gradient(140deg, rgba(246, 211, 74, 0.94) 0%, rgba(255, 248, 220, 0.92) 54%, rgba(255, 255, 255, 0.95) 100%);
  border: 1px solid rgba(246, 211, 74, 0.4);
}

.profile-hero__badge {
  width: fit-content;
  display: inline-flex;
  align-items: center;
  gap: 10px;
  min-height: 28px;
  padding: 0 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.76);
  font-size: 11px;
  font-weight: 700;
}

.profile-hero__badge strong {
  font-size: 10px;
  letter-spacing: 0.08em;
}

.profile-hero h1 {
  margin: 10px 0 0;
  font-size: 24px;
  line-height: 1.2;
}

.profile-hero p {
  margin: 8px 0 0;
  max-width: 760px;
  line-height: 1.5;
  font-size: 13px;
  color: rgba(23, 23, 23, 0.74);
}

.profile-hero__actions {
  margin-top: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.profile-panel {
  padding: 18px;
}

.profile-panel__header h2 {
  margin: 0;
  font-size: 18px;
}

.profile-panel__header span {
  display: block;
  margin-top: 4px;
  font-size: 13px;
  color: var(--text-muted);
  line-height: 1.4;
}

.profile-avatar {
  margin-top: 14px;
  display: flex;
  gap: 14px;
  align-items: center;
  padding: 12px;
  border-radius: 14px;
  background: rgba(255, 250, 233, 0.68);
  border: 1px solid rgba(246, 211, 74, 0.32);
}

.profile-avatar__preview {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  overflow: hidden;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: rgba(246, 211, 74, 0.44);
  color: var(--text-main);
  font-size: 22px;
  font-weight: 800;
  flex-shrink: 0;
}

.profile-avatar__preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.profile-avatar__meta strong {
  display: block;
  font-size: 15px;
}

.profile-avatar__meta p {
  margin: 4px 0 0;
  font-size: 13px;
  color: var(--text-subtle);
}

.profile-avatar__buttons {
  margin-top: 8px;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.profile-avatar__meta small {
  display: block;
  margin-top: 6px;
  font-size: 12px;
  color: var(--text-muted);
  line-height: 1.4;
}

.profile-avatar__input {
  display: none;
}

.profile-form-grid {
  margin-top: 14px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.profile-form__field {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.profile-form__field span {
  font-size: 12px;
  font-weight: 700;
  color: var(--text-subtle);
}

.profile-form__field input {
  min-height: 38px;
  border-radius: 10px;
  border: 1px solid var(--border-color);
  padding: 0 12px;
  background: #fff;
  font-size: 14px;
}

.profile-form__field input[readonly] {
  background: rgba(249, 250, 251, 0.9);
  color: var(--text-muted);
}

.profile-form__field small {
  font-size: 11px;
  color: var(--text-muted);
}

.profile-panel__footer {
  margin-top: 14px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

@media (max-width: 1080px) {
  .profile-content {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 860px) {
  .profile-hero {
    padding: 16px 18px;
  }

  .profile-hero h1 {
    font-size: 20px;
  }

  .profile-form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
