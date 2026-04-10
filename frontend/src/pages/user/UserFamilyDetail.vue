<template>
  <div class="finance-page user-family-detail-page">
    <div v-if="errorMessage" class="finance-inline-notice">
      <span>{{ errorMessage }}</span>
      <button type="button" @click="retryDetail">重新加载</button>
    </div>

    <section v-if="family" class="finance-toolbar">
      <div class="finance-toolbar__meta">
        <h1 class="page-title">{{ family.name }}</h1>
        <p class="page-description">{{ family.slogan }}</p>
      </div>
      <div class="finance-toolbar__actions">
        <button class="finance-button finance-button--ghost" @click="goBack">返回家庭列表</button>
        <button class="finance-button finance-button--primary" :disabled="!family.inviteLink" @click="copyInviteLink">
          复制邀请链接
        </button>
      </div>
    </section>

    <section v-if="detailLoading && !family" class="page-card user-family-detail-page__state">
      <h2>正在加载家庭详情</h2>
      <p>月度、年度汇总和成员占比会在真实接口返回后展示。</p>
    </section>

    <template v-else-if="family">
      <section class="finance-hero">
        <div class="finance-hero__eyebrow">
          <span class="finance-hero__eyebrow-mark">统</span>
          <span>家庭统计总览</span>
        </div>
        <div class="finance-hero__headline">
          <h2>同一个家庭里谁贡献了收入、谁拉高了支出，一眼就能看出来。</h2>
          <p>
            创建人：{{ family.creator }} · 成员数：{{ family.members.length }} · 家庭 ID：{{ family.id }}
          </p>
        </div>
        <div class="user-family-detail-page__filters">
          <label>
            月份
            <el-select v-model="selectedMonth" size="large" :disabled="detailLoading || !family.monthOptions.length" @change="handlePeriodChange('month')">
              <el-option v-for="item in family.monthOptions" :key="item.key" :label="item.label" :value="item.key" />
            </el-select>
          </label>
          <label>
            年份
            <el-select v-model="selectedYear" size="large" :disabled="detailLoading || !family.yearOptions.length" @change="handlePeriodChange('year')">
              <el-option v-for="item in family.yearOptions" :key="item.key" :label="item.label" :value="item.key" />
            </el-select>
          </label>
        </div>
      </section>

      <section class="user-family-detail-page__summary-grid">
        <article class="page-card summary-card">
          <header>
            <h3>{{ monthSummary.label }} 汇总</h3>
            <p>{{ monthSummary.note }}</p>
          </header>
          <div class="summary-card__metrics">
            <button
              :class="['summary-card__metric', activeShareKey === 'month-income' ? 'is-active' : '']"
              :disabled="shareLoading"
              @click="setShareMode('month-income')"
            >
              <span>月收入</span>
              <strong class="finance-tone-income">{{ formatCurrency(monthSummary.income) }}</strong>
            </button>
            <button
              :class="['summary-card__metric', activeShareKey === 'month-expense' ? 'is-active' : '']"
              :disabled="shareLoading"
              @click="setShareMode('month-expense')"
            >
              <span>月支出</span>
              <strong class="finance-tone-expense">{{ formatCurrency(monthSummary.expense) }}</strong>
            </button>
            <div class="summary-card__metric summary-card__metric--static">
              <span>月结余</span>
              <strong>{{ formatCurrency(monthSummary.balance) }}</strong>
            </div>
          </div>
        </article>

        <article class="page-card summary-card">
          <header>
            <h3>{{ yearSummary.label }} 汇总</h3>
            <p>{{ yearSummary.note }}</p>
          </header>
          <div class="summary-card__metrics">
            <button
              :class="['summary-card__metric', activeShareKey === 'year-income' ? 'is-active' : '']"
              :disabled="shareLoading"
              @click="setShareMode('year-income')"
            >
              <span>年收入</span>
              <strong class="finance-tone-income">{{ formatCurrency(yearSummary.income) }}</strong>
            </button>
            <button
              :class="['summary-card__metric', activeShareKey === 'year-expense' ? 'is-active' : '']"
              :disabled="shareLoading"
              @click="setShareMode('year-expense')"
            >
              <span>年支出</span>
              <strong class="finance-tone-expense">{{ formatCurrency(yearSummary.expense) }}</strong>
            </button>
            <div class="summary-card__metric summary-card__metric--static">
              <span>年结余</span>
              <strong>{{ formatCurrency(yearSummary.balance) }}</strong>
            </div>
          </div>
        </article>
      </section>

      <div v-if="shareErrorMessage" class="finance-inline-notice">
        <span>{{ shareErrorMessage }}</span>
        <button type="button" @click="loadShareData">重新加载占比</button>
      </div>

      <div v-if="shareLoading" class="user-family-detail-page__share-tip">正在同步成员占比数据…</div>
      <FamilyMemberSharePanel :title="shareData.title" :total="shareData.total" :rows="shareData.rows" />

      <section class="page-card user-family-detail-page__member-panel">
        <header>
          <h3>家庭成员</h3>
          <p>当前展示成员角色与颜色标识，便于对照占比图。</p>
        </header>
        <div class="user-family-detail-page__member-list">
          <article v-for="member in family.members" :key="member.userId || member.name" class="member-row">
            <div class="member-row__identity">
              <span class="member-row__dot" :style="{ backgroundColor: member.color }"></span>
              <div>
                <strong>{{ member.name }}</strong>
                <p>{{ member.role }}</p>
              </div>
            </div>
            <span class="member-row__tag">{{ member.role === "创建人" ? "Owner" : "Member" }}</span>
          </article>
        </div>
      </section>
    </template>

    <section v-else-if="isNotFound" class="page-card user-family-detail-page__missing">
      <h2>未找到该家庭</h2>
      <p>当前家庭 ID 可能不存在或已失效，请返回列表重新选择。</p>
      <button class="finance-button finance-button--primary" @click="goBack">返回家庭列表</button>
    </section>

    <section v-else class="page-card user-family-detail-page__state">
      <h2>家庭详情暂时不可用</h2>
      <p>接口暂未返回可展示的数据，你可以稍后重试或返回列表重新进入。</p>
      <button class="finance-button finance-button--primary" @click="retryDetail">重新加载</button>
    </section>
  </div>
</template>

<script>
import { ElMessage } from "element-plus";

import {
  buildUserFamilyError,
  formatFamilyCurrency,
  getUserFamilyDetail,
  getUserFamilyMemberShare,
  normalizeFamilyDetailPayload,
  normalizeFamilySharePayload
} from "@/api/userFamily";
import FamilyMemberSharePanel from "@/components/FamilyMemberSharePanel.vue";

function createEmptyShareData() {
  return {
    title: "暂无成员占比",
    total: 0,
    rows: []
  };
}

function findOption(options, target) {
  if (!options || !options.length) {
    return {
      key: "",
      label: "暂无数据",
      note: "暂无说明",
      income: 0,
      expense: 0,
      balance: 0
    };
  }

  var matched = options.find(function(item) {
    return item.key === target;
  });

  return matched || options[0];
}

export default {
  name: "UserFamilyDetail",
  components: {
    FamilyMemberSharePanel
  },
  data() {
    return {
      family: null,
      selectedMonth: "",
      selectedYear: "",
      activeShareKey: "month-income",
      detailLoading: false,
      shareLoading: false,
      isNotFound: false,
      errorMessage: "",
      shareErrorMessage: "",
      shareData: createEmptyShareData(),
      detailRequestKey: 0,
      shareRequestKey: 0
    };
  },
  computed: {
    familyId() {
      return this.$route.params.familyId;
    },
    monthSummary() {
      return findOption(this.family ? this.family.monthOptions : [], this.selectedMonth);
    },
    yearSummary() {
      return findOption(this.family ? this.family.yearOptions : [], this.selectedYear);
    },
    shareParams() {
      var mapping = {
        "month-income": { periodType: "month", metricType: "income", periodKey: this.selectedMonth },
        "month-expense": { periodType: "month", metricType: "expense", periodKey: this.selectedMonth },
        "year-income": { periodType: "year", metricType: "income", periodKey: this.selectedYear },
        "year-expense": { periodType: "year", metricType: "expense", periodKey: this.selectedYear }
      };

      return mapping[this.activeShareKey] || mapping["month-income"];
    }
  },
  watch: {
    familyId: {
      immediate: true,
      handler() {
        this.loadFamily();
      }
    }
  },
  methods: {
    formatCurrency: formatFamilyCurrency,
    loadFamily() {
      var vm = this;
      var requestKey = this.detailRequestKey + 1;

      this.detailRequestKey = requestKey;
      this.detailLoading = true;
      this.family = null;
      this.isNotFound = false;
      this.errorMessage = "";
      this.shareErrorMessage = "";
      this.shareData = createEmptyShareData();

      if (!this.familyId) {
        this.family = null;
        this.selectedMonth = "";
        this.selectedYear = "";
        this.activeShareKey = "month-income";
        this.detailLoading = false;
        return;
      }

      return getUserFamilyDetail(this.familyId)
        .then(function(result) {
          var detail;

          if (requestKey !== vm.detailRequestKey) {
            return;
          }

          detail = normalizeFamilyDetailPayload(result, vm.familyId);

          if (!detail || !detail.id) {
            vm.family = null;
            vm.isNotFound = true;
            vm.selectedMonth = "";
            vm.selectedYear = "";
            vm.activeShareKey = "month-income";
            return;
          }

          vm.family = detail;
          vm.selectedMonth = detail.monthOptions.length ? detail.monthOptions[0].key : "";
          vm.selectedYear = detail.yearOptions.length ? detail.yearOptions[0].key : "";
          vm.activeShareKey = detail.monthOptions.length ? "month-income" : "year-income";
          return vm.loadShareData();
        })
        .catch(function(error) {
          if (requestKey !== vm.detailRequestKey) {
            return;
          }

          vm.family = null;
          vm.selectedMonth = "";
          vm.selectedYear = "";
          vm.activeShareKey = "month-income";

          if (error && error.response && (error.response.status === 403 || error.response.status === 404)) {
            vm.isNotFound = true;
          } else {
            vm.errorMessage = buildUserFamilyError(error, "家庭详情加载失败，请稍后重试。");
          }
        })
        .then(function() {
          if (requestKey === vm.detailRequestKey) {
            vm.detailLoading = false;
          }
        });
    },
    loadShareData() {
      var vm = this;

      if (!this.family) {
        this.shareData = createEmptyShareData();
        return;
      }

      if (!this.shareParams.periodKey) {
        this.shareData = createEmptyShareData();
        return;
      }

      var requestKey = this.shareRequestKey + 1;
      var params = {
        periodType: this.shareParams.periodType,
        period_type: this.shareParams.periodType,
        metricType: this.shareParams.metricType,
        metric_type: this.shareParams.metricType,
        periodKey: this.shareParams.periodKey,
        period_key: this.shareParams.periodKey,
        month: this.shareParams.periodType === "month" ? this.shareParams.periodKey : "",
        year: this.shareParams.periodType === "year" ? this.shareParams.periodKey : ""
      };

      this.shareRequestKey = requestKey;
      this.shareLoading = true;
      this.shareErrorMessage = "";

      return getUserFamilyMemberShare(this.family.id, params)
        .then(function(result) {
          if (requestKey !== vm.shareRequestKey) {
            return;
          }

          vm.shareData = normalizeFamilySharePayload(result, vm.shareParams, vm.family);
        })
        .catch(function(error) {
          if (requestKey !== vm.shareRequestKey) {
            return;
          }

          vm.shareData = createEmptyShareData();
          vm.shareErrorMessage = buildUserFamilyError(error, "成员占比加载失败，请稍后重试。");
        })
        .then(function() {
          if (requestKey === vm.shareRequestKey) {
            vm.shareLoading = false;
          }
        });
    },
    setShareMode(mode) {
      this.activeShareKey = mode;
      this.loadShareData();
    },
    handlePeriodChange(periodType) {
      if (periodType === "month" && this.activeShareKey.indexOf("month-") === 0) {
        this.loadShareData();
      }

      if (periodType === "year" && this.activeShareKey.indexOf("year-") === 0) {
        this.loadShareData();
      }
    },
    goBack() {
      this.$router.push("/user/families");
    },
    retryDetail() {
      this.loadFamily();
    },
    copyInviteLink() {
      if (!this.family) {
        return;
      }

      var text = this.family.inviteLink || "";

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
.user-family-detail-page__state {
  padding: 40px 30px;
  text-align: center;
}

.user-family-detail-page__state h2 {
  margin: 0;
  font-size: 30px;
}

.user-family-detail-page__state p {
  margin: 12px auto 0;
  max-width: 560px;
  color: var(--text-subtle);
  line-height: 1.7;
}

.user-family-detail-page__state button {
  margin-top: 18px;
}

.user-family-detail-page__filters {
  margin-top: 18px;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.user-family-detail-page__filters label {
  display: flex;
  flex-direction: column;
  gap: 8px;
  color: var(--text-subtle);
  font-size: 13px;
  min-width: 220px;
}

.user-family-detail-page__summary-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.user-family-detail-page__share-tip {
  margin: 0 0 12px;
  color: var(--text-subtle);
  font-size: 13px;
}

.summary-card {
  padding: 22px;
}

.summary-card h3 {
  margin: 0;
  font-size: 20px;
}

.summary-card p {
  margin: 8px 0 0;
  color: var(--text-subtle);
  line-height: 1.6;
}

.summary-card__metrics {
  margin-top: 14px;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.summary-card__metric {
  border: 1px solid var(--border-color);
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.92);
  padding: 12px;
  text-align: left;
  transition: 0.2s ease;
}

.summary-card__metric:hover {
  box-shadow: var(--shadow-sm);
  transform: translateY(-1px);
}

.summary-card__metric span {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
}

.summary-card__metric strong {
  display: block;
  margin-top: 7px;
  font-size: 18px;
}

.summary-card__metric.is-active {
  border-color: rgba(246, 211, 74, 0.85);
  background: rgba(246, 211, 74, 0.2);
}

.summary-card__metric--static {
  cursor: default;
}

.summary-card__metric--static:hover {
  transform: none;
  box-shadow: none;
}

.user-family-detail-page__member-panel {
  padding: 22px;
}

.user-family-detail-page__member-panel h3 {
  margin: 0;
  font-size: 20px;
}

.user-family-detail-page__member-panel p {
  margin: 8px 0 0;
  color: var(--text-subtle);
}

.user-family-detail-page__member-list {
  margin-top: 14px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.member-row {
  border: 1px solid var(--border-color);
  border-radius: 14px;
  padding: 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.member-row__identity {
  display: flex;
  align-items: center;
  gap: 10px;
}

.member-row__dot {
  width: 12px;
  height: 12px;
  border-radius: 999px;
}

.member-row__identity strong {
  display: block;
  font-size: 15px;
}

.member-row__identity p {
  margin: 4px 0 0;
  color: var(--text-muted);
  font-size: 12px;
}

.member-row__tag {
  min-height: 28px;
  padding: 0 10px;
  border-radius: 999px;
  background: rgba(246, 211, 74, 0.24);
  display: inline-flex;
  align-items: center;
  font-size: 12px;
  font-weight: 700;
}

.user-family-detail-page__missing {
  padding: 40px 30px;
  text-align: center;
}

.user-family-detail-page__missing h2 {
  margin: 0;
  font-size: 30px;
}

.user-family-detail-page__missing p {
  margin: 12px auto 0;
  max-width: 560px;
  color: var(--text-subtle);
  line-height: 1.7;
}

.user-family-detail-page__missing button {
  margin-top: 18px;
}

@media (max-width: 1200px) {
  .summary-card__metrics {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 960px) {
  .user-family-detail-page__summary-grid {
    grid-template-columns: minmax(0, 1fr);
  }

  .user-family-detail-page__member-list {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
