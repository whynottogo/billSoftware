<template>
  <div class="finance-page user-families-page">
    <section class="finance-toolbar">
      <div class="finance-toolbar__meta">
        <h1 class="page-title">家庭列表</h1>
        <p class="page-description">你可以创建家庭、通过家庭 ID 或邀请链接加入，并随时进入家庭详情查看汇总。</p>
      </div>
      <div class="finance-toolbar__actions">
        <button class="finance-button finance-button--ghost" @click="openJoinDialog">加入家庭</button>
        <button class="finance-button finance-button--primary" @click="openCreateDialog">创建家庭</button>
      </div>
    </section>

    <section class="finance-hero finance-hero--soft">
      <div class="finance-hero__eyebrow">
        <span class="finance-hero__eyebrow-mark">家</span>
        <span>家庭协同记账</span>
      </div>
      <div class="finance-hero__headline">
        <h2>把同一家庭的收支放在一个上下文里看，月度与年度节奏会更清楚。</h2>
        <p>当前展示的是原型级交互，后续可以把创建、加入和退出动作接入真实家庭服务。</p>
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

    <section v-if="families.length" class="user-families-page__list">
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
            <button class="finance-button finance-button--ghost family-card__copy" @click="copyInviteLink(family)">
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
          <button class="finance-button finance-button--primary" @click="openFamilyDetail(family.id)">查看详情</button>
          <button class="finance-button finance-button--ghost" @click="openJoinDialog(family.id)">按 ID 加入</button>
          <button class="finance-button finance-button--ghost" @click="leaveFamilyAsCurrentUser(family.id)">退出家庭</button>
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
        <button class="finance-button finance-button--ghost" @click="createDialogVisible = false">取消</button>
        <button class="finance-button finance-button--primary" @click="submitCreateFamily">创建</button>
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
        <button class="finance-button finance-button--ghost" @click="joinDialogVisible = false">取消</button>
        <button class="finance-button finance-button--primary" @click="submitJoinFamily">加入</button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ElMessage, ElMessageBox } from "element-plus";

import UserFamilyEmptyState from "@/components/UserFamilyEmptyState.vue";
import {
  formatCurrency,
  listFamilies,
  getFamilyOverview,
  createFamily,
  joinFamilyById,
  joinFamilyByInviteLink,
  leaveFamily
} from "@/utils/userFamilyMock";

export default {
  name: "UserFamilies",
  components: {
    UserFamilyEmptyState
  },
  data() {
    return {
      families: [],
      overview: {
        familyCount: 0,
        totalMembers: 0,
        joinedCount: 0
      },
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
  methods: {
    formatCurrency: formatCurrency,
    refreshFamilies() {
      this.families = listFamilies();
      this.overview = getFamilyOverview();
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
      var result = createFamily({
        name: this.createForm.name,
        slogan: this.createForm.slogan,
        creator: this.getCurrentMemberName()
      });

      if (!result.ok) {
        ElMessage.warning(result.message);
        return;
      }

      this.createDialogVisible = false;
      this.refreshFamilies();
      ElMessage.success("家庭创建成功，已生成家庭 ID：" + result.family.id);
    },
    submitJoinFamily() {
      var payload = {
        familyId: this.joinForm.familyId,
        inviteLink: this.joinForm.inviteLink,
        memberName: this.joinForm.memberName || this.getCurrentMemberName()
      };
      var result = this.joinForm.mode === "id" ? joinFamilyById(payload) : joinFamilyByInviteLink(payload);

      if (!result.ok) {
        ElMessage.warning(result.message);
        return;
      }

      this.joinDialogVisible = false;
      this.refreshFamilies();
      ElMessage.success("加入家庭成功：" + result.family.name);
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
          var result = leaveFamily({
            familyId: familyId,
            memberName: memberName
          });

          if (!result.ok) {
            ElMessage.warning(result.message);
            return;
          }

          vm.refreshFamilies();
          ElMessage.success("已退出家庭");
        })
        .catch(function() {});
    },
    openFamilyDetail(familyId) {
      this.$router.push("/user/families/" + familyId);
    },
    copyInviteLink(family) {
      var text = family.inviteLink;

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
