<template>
  <div class="finance-page user-families-page">
    <div v-if="errorMessage" class="finance-inline-notice">
      <span>{{ errorMessage }}</span>
      <button type="button" @click="refreshFamilies">重新加载</button>
    </div>

    <section class="finance-toolbar">
      <div class="finance-toolbar__meta">
        <h1 class="page-title">家庭列表</h1>
        <p class="page-description">你可以创建家庭、通过家庭 ID 或邀请链接加入，并随时进入家庭详情查看汇总。</p>
      </div>
      <div class="finance-toolbar__actions">
        <button class="finance-button finance-button--ghost" :disabled="actionLocked" @click="openJoinDialog">加入家庭</button>
        <button class="finance-button finance-button--primary" :disabled="actionLocked" @click="openCreateDialog">创建家庭</button>
      </div>
    </section>

    <section class="finance-hero finance-hero--soft">
      <div class="finance-hero__eyebrow">
        <span class="finance-hero__eyebrow-mark">家</span>
        <span>家庭协同记账</span>
      </div>
      <div class="finance-hero__headline">
        <h2>把同一家庭的收支放在一个上下文里看，月度与年度节奏会更清楚。</h2>
        <p>{{ heroDescription }}</p>
      </div>
      <div class="finance-stat-grid finance-stat-grid--triple">
        <article class="finance-stat-card">
          <span>家庭数量</span>
          <strong>{{ overview.familyCount }}</strong>
          <small>可进入详情继续查看汇总</small>
        </article>
        <article class="finance-stat-card">
          <span>总成员数</span>
          <strong>{{ overview.totalMembers }}</strong>
          <small>跨家庭累计成员规模</small>
        </article>
        <article class="finance-stat-card">
          <span>我已加入</span>
          <strong>{{ overview.joinedCount }}</strong>
          <small>按当前登录资料自动识别</small>
        </article>
      </div>
    </section>

    <section v-if="isLoading && !families.length" class="page-card user-families-page__placeholder">
      <strong>正在同步家庭列表…</strong>
      <p>创建、加入和退出后的最新数据会在这里刷新。</p>
    </section>

    <section v-else-if="families.length" class="user-families-page__list">
      <article v-for="family in families" :key="family.id" class="page-card family-card">
        <header class="family-card__header">
          <div>
            <div class="family-card__title">
              <span class="finance-badge">家庭</span>
              <h3>{{ family.name }}</h3>
            </div>
            <p>{{ family.slogan }}</p>
          </div>
          <div class="family-card__meta">
            <span>ID {{ family.id }}</span>
            <button
              class="finance-button finance-button--ghost family-card__copy"
              :disabled="!family.inviteLink"
              @click="copyInviteLink(family)"
            >
              复制邀请链接
            </button>
          </div>
        </header>

        <div class="family-card__stats">
          <div>
            <span>本月收入</span>
            <strong class="finance-tone-income">{{ formatCurrency(family.monthIncome) }}</strong>
          </div>
          <div>
            <span>本月支出</span>
            <strong class="finance-tone-expense">{{ formatCurrency(family.monthExpense) }}</strong>
          </div>
          <div>
            <span>本月结余</span>
            <strong>{{ formatCurrency(family.monthBalance) }}</strong>
          </div>
          <div>
            <span>年度结余</span>
            <strong>{{ formatCurrency(family.yearBalance) }}</strong>
          </div>
        </div>

        <div class="family-card__members">
          <span>创建人：{{ family.creator }}</span>
          <span>成员数：{{ family.memberCount }}</span>
          <span>创建时间：{{ family.createdAt }}</span>
        </div>

        <div class="family-card__avatars">
          <span
            v-for="member in family.members"
            :key="family.id + '-' + member.name"
            class="family-card__avatar"
            :style="{ backgroundColor: member.color }"
          >
            {{ member.name.slice(0, 1) }}
          </span>
        </div>

        <footer class="family-card__actions">
          <button class="finance-button finance-button--primary" :disabled="actionLocked" @click="openFamilyDetail(family.id)">
            查看详情
          </button>
          <button class="finance-button finance-button--ghost" :disabled="actionLocked" @click="openJoinDialog(family.id)">
            按 ID 加入
          </button>
          <button
            class="finance-button finance-button--ghost"
            :disabled="actionLocked || leavingFamilyId === family.id"
            @click="leaveFamilyAsCurrentUser(family.id)"
          >
            {{ leavingFamilyId === family.id ? "退出中..." : "退出家庭" }}
          </button>
        </footer>
      </article>
    </section>

    <UserFamilyEmptyState v-else @create="openCreateDialog" @join="openJoinDialog" />

    <el-dialog v-model="createDialogVisible" title="创建家庭" width="520px">
      <el-form label-position="top">
        <el-form-item label="家庭名称">
          <el-input v-model="createForm.name" maxlength="24" show-word-limit placeholder="例如：周末小家" />
        </el-form-item>
        <el-form-item label="家庭介绍">
          <el-input
            v-model="createForm.slogan"
            type="textarea"
            :rows="3"
            maxlength="80"
            show-word-limit
            placeholder="补充一句家庭记账目标，例如：一起压低餐饮支出波动。"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <button class="finance-button finance-button--ghost" :disabled="isCreating" @click="createDialogVisible = false">取消</button>
        <button class="finance-button finance-button--primary" :disabled="isCreating" @click="submitCreateFamily">
          {{ isCreating ? "创建中..." : "创建" }}
        </button>
      </template>
    </el-dialog>

    <el-dialog v-model="joinDialogVisible" title="加入家庭" width="540px">
      <el-form label-position="top">
        <el-form-item label="加入方式">
          <el-radio-group v-model="joinForm.mode">
            <el-radio label="id">家庭 ID</el-radio>
            <el-radio label="link">邀请链接</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="joinForm.mode === 'id'" label="家庭 ID">
          <el-input v-model="joinForm.familyId" placeholder="例如：FAM-4821" />
        </el-form-item>
        <el-form-item v-else label="邀请链接">
          <el-input v-model="joinForm.inviteLink" placeholder="例如：https://bill.local/invite/FAM-4821" />
        </el-form-item>
        <el-form-item label="成员名称">
          <el-input v-model="joinForm.memberName" placeholder="默认使用当前登录资料昵称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <button class="finance-button finance-button--ghost" :disabled="isJoining" @click="joinDialogVisible = false">取消</button>
        <button class="finance-button finance-button--primary" :disabled="isJoining" @click="submitJoinFamily">
          {{ isJoining ? "加入中..." : "加入" }}
        </button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ElMessage, ElMessageBox } from "element-plus";

import {
  buildUserFamilyError,
  createUserFamily,
  formatFamilyCurrency,
  getUserFamilies,
  joinUserFamilyById,
  joinUserFamilyByInviteLink,
  leaveUserFamily,
  normalizeFamilyListPayload,
  normalizeFamilyMutationPayload
} from "@/api/userFamily";
import UserFamilyEmptyState from "@/components/UserFamilyEmptyState.vue";

function createEmptyOverview() {
  return {
    familyCount: 0,
    totalMembers: 0,
    joinedCount: 0
  };
}

function buildCreatePayload(form) {
  var name = String(form.name || "").trim();
  var slogan = String(form.slogan || "").trim();

  return {
    name: name,
    family_name: name,
    slogan: slogan,
    description: slogan,
    intro: slogan
  };
}

function buildJoinPayload(form, memberName) {
  var normalizedMemberName = String(memberName || "").trim();

  return {
    familyId: String(form.familyId || "").trim(),
    family_id: String(form.familyId || "").trim(),
    inviteLink: String(form.inviteLink || "").trim(),
    invite_link: String(form.inviteLink || "").trim(),
    memberName: normalizedMemberName,
    member_name: normalizedMemberName
  };
}

export default {
  name: "UserFamilies",
  components: {
    UserFamilyEmptyState
  },
  data() {
    return {
      families: [],
      overview: createEmptyOverview(),
      isLoading: false,
      isCreating: false,
      isJoining: false,
      leavingFamilyId: "",
      errorMessage: "",
      createDialogVisible: false,
      joinDialogVisible: false,
      createForm: {
        name: "",
        slogan: ""
      },
      joinForm: {
        mode: "id",
        familyId: "",
        inviteLink: "",
        memberName: ""
      }
    };
  },
  created() {
    this.refreshFamilies();
  },
  computed: {
    actionLocked() {
      return this.isLoading || this.isCreating || this.isJoining || !!this.leavingFamilyId;
    },
    heroDescription() {
      if (this.isLoading) {
        return "正在从真实家庭服务同步你当前账号可访问的家庭列表。";
      }

      if (!this.families.length) {
        return "当前还没有已加入的家庭，创建一个新家庭或通过家庭 ID、邀请链接加入都可以。";
      }

      return "当前展示的是基于真实接口返回的家庭与汇总数据，创建、加入和退出都会立即回刷列表。";
    }
  },
  methods: {
    formatCurrency: formatFamilyCurrency,
    refreshFamilies() {
      var vm = this;

      this.isLoading = true;
      this.errorMessage = "";

      return getUserFamilies()
        .then(function(result) {
        var normalized = normalizeFamilyListPayload(result);

          vm.families = normalized.families;
          vm.overview = normalized.overview;
        })
        .catch(function(error) {
          vm.families = [];
          vm.overview = createEmptyOverview();
          vm.errorMessage = buildUserFamilyError(error, "家庭列表加载失败，请稍后重试。");
        })
        .then(function() {
          vm.isLoading = false;
        });
    },
    getCurrentMemberName() {
      var rawProfile = localStorage.getItem("bill_user_profile");

      if (!rawProfile) {
        return "playwright-user";
      }

      try {
        var profile = JSON.parse(rawProfile);
        return profile.nickname || profile.account || profile.username || "playwright-user";
      } catch (error) {
        return "playwright-user";
      }
    },
    openCreateDialog() {
      this.createForm = {
        name: "",
        slogan: ""
      };
      this.createDialogVisible = true;
    },
    openJoinDialog(familyId) {
      this.joinForm = {
        mode: "id",
        familyId: familyId || "",
        inviteLink: "",
        memberName: this.getCurrentMemberName()
      };
      this.joinDialogVisible = true;
    },
    submitCreateFamily() {
      var vm = this;

      if (!String(this.createForm.name || "").trim()) {
        ElMessage.warning("请先填写家庭名称");
        return;
      }

      this.isCreating = true;

      createUserFamily(buildCreatePayload(this.createForm))
        .then(function(result) {
          var normalized = normalizeFamilyMutationPayload(result);

          vm.createDialogVisible = false;
          return vm.refreshFamilies().then(function() {
            ElMessage.success("家庭创建成功" + (normalized.family && normalized.family.id ? "，已生成家庭 ID：" + normalized.family.id : ""));
          });
        })
        .catch(function(error) {
          ElMessage.warning(buildUserFamilyError(error, "家庭创建失败，请稍后重试。"));
        })
        .then(function() {
          vm.isCreating = false;
        });
    },
    submitJoinFamily() {
      var vm = this;

      if (this.joinForm.mode === "id" && !String(this.joinForm.familyId || "").trim()) {
        ElMessage.warning("请先填写家庭 ID");
        return;
      }

      if (this.joinForm.mode === "link" && !String(this.joinForm.inviteLink || "").trim()) {
        ElMessage.warning("请先填写邀请链接");
        return;
      }

      this.isJoining = true;

      var payload = buildJoinPayload(this.joinForm, this.joinForm.memberName || this.getCurrentMemberName());
      var joinRequest =
        this.joinForm.mode === "id" ? joinUserFamilyById(payload) : joinUserFamilyByInviteLink(payload);

      joinRequest
        .then(function(result) {
          var normalized = normalizeFamilyMutationPayload(result);

          vm.joinDialogVisible = false;
          return vm.refreshFamilies().then(function() {
            ElMessage.success("加入家庭成功" + (normalized.family && normalized.family.name ? "：" + normalized.family.name : ""));
          });
        })
        .catch(function(error) {
          ElMessage.warning(buildUserFamilyError(error, "加入家庭失败，请稍后重试。"));
        })
        .then(function() {
          vm.isJoining = false;
        });
    },
    leaveFamilyAsCurrentUser(familyId) {
      var vm = this;
      var memberName = this.getCurrentMemberName();

      ElMessageBox.confirm("确认以 " + memberName + " 身份退出该家庭吗？", "退出家庭", {
        confirmButtonText: "确认退出",
        cancelButtonText: "取消",
        type: "warning"
      })
        .then(function() {
          vm.leavingFamilyId = familyId;

          return leaveUserFamily(familyId, {
            memberName: memberName,
            member_name: memberName
          })
            .then(function() {
              return vm.refreshFamilies();
            })
            .then(function() {
              ElMessage.success("已退出家庭");
            })
            .catch(function(error) {
              ElMessage.warning(buildUserFamilyError(error, "退出家庭失败，请稍后重试。"));
            })
            .then(function() {
              vm.leavingFamilyId = "";
            });
        })
        .catch(function() {});
    },
    openFamilyDetail(familyId) {
      this.$router.push("/user/families/" + familyId);
    },
    copyInviteLink(family) {
      var text = family && family.inviteLink ? family.inviteLink : "";

      if (!text) {
        ElMessage.info("当前家庭暂无可复制的邀请链接");
        return;
      }

      if (navigator && navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard
          .writeText(text)
          .then(function() {
            ElMessage.success("邀请链接已复制");
          })
          .catch(function() {
            ElMessage.info("复制失败，请手动复制： " + text);
          });
        return;
      }

      ElMessage.info("当前浏览器不支持自动复制，请手动复制： " + text);
    }
  }
};
</script>

<style scoped>
.user-families-page__placeholder {
  padding: 34px 28px;
  text-align: center;
}

.user-families-page__placeholder strong {
  display: block;
  font-size: 20px;
}

.user-families-page__placeholder p {
  margin: 10px 0 0;
  color: var(--text-subtle);
}

.user-families-page__list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.family-card {
  padding: 24px;
}

.family-card__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 18px;
}

.family-card__title {
  display: flex;
  align-items: center;
  gap: 10px;
}

.family-card__title h3 {
  margin: 0;
  font-size: 22px;
}

.family-card__header p {
  margin: 10px 0 0;
  color: var(--text-subtle);
  line-height: 1.6;
}

.family-card__meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8px;
}

.family-card__meta span {
  color: var(--text-subtle);
  font-size: 13px;
}

.family-card__copy {
  min-height: 36px;
}

.family-card__stats {
  margin-top: 16px;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.family-card__stats div {
  padding: 14px;
  border-radius: 16px;
  border: 1px solid var(--border-color);
  background: rgba(255, 255, 255, 0.86);
}

.family-card__stats span {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
}

.family-card__stats strong {
  display: block;
  margin-top: 8px;
  font-size: 20px;
}

.family-card__members {
  margin-top: 14px;
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  color: var(--text-subtle);
  font-size: 13px;
}

.family-card__avatars {
  margin-top: 14px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.family-card__avatar {
  width: 34px;
  height: 34px;
  border-radius: 999px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #111827;
  font-weight: 700;
  font-size: 14px;
}

.family-card__actions {
  margin-top: 16px;
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

@media (max-width: 1200px) {
  .family-card__stats {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 780px) {
  .family-card__header {
    flex-direction: column;
  }

  .family-card__meta {
    align-items: flex-start;
  }
}
</style>
