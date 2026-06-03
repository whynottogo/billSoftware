<template>
  <div class="admin-family-detail-page">
    <section class="admin-family-detail-page__header page-card">
      <div>
        <span class="admin-family-detail-page__eyebrow">管理端 / 家庭管理 / 详情</span>
        <h1>{{ detail.family.name }}</h1>
        <p>{{ detail.family.id }} · 创建者 {{ detail.family.creator }} · 创建于 {{ detail.family.createdAt }}</p>
      </div>
      <div class="admin-family-detail-page__actions">
        <button class="finance-button finance-button--ghost" @click="goBack">返回家庭列表</button>
      </div>
    </section>

    <section v-if="loading" class="page-card admin-family-detail-page__state">
      正在加载家庭详情...
    </section>

    <section v-else-if="notFound" class="page-card admin-family-detail-page__state">
      <strong>未找到该家庭</strong>
      <p>家庭 ID 不存在或已被删除。</p>
      <button class="finance-button finance-button--ghost" @click="goBack">返回家庭列表</button>
    </section>

    <section v-else-if="loadError" class="page-card admin-family-detail-page__state is-error">
      <p>{{ loadError }}</p>
      <button class="finance-button finance-button--ghost" @click="fetchDetail">重新加载</button>
    </section>

    <template v-else>
      <section class="admin-family-detail-page__summary">
        <article class="admin-family-stat">
          <span>成员数量</span>
          <strong>{{ detail.family.memberCount }} 位</strong>
        </article>
        <article class="admin-family-stat">
          <span>账单记录</span>
          <strong>{{ detail.incomeExpense.records }} 笔</strong>
        </article>
        <article class="admin-family-stat">
          <span>家庭状态</span>
          <strong>{{ detail.family.status }}</strong>
        </article>
      </section>

      <section class="admin-family-detail-page__metrics">
        <article class="page-card admin-family-panel">
          <header>
            <span>收支汇总</span>
            <strong>{{ formatCurrency(detail.incomeExpense.balance) }}</strong>
          </header>
          <div class="admin-family-detail-page__metric-grid">
            <p><span>收入</span><strong class="is-income">{{ formatCurrency(detail.incomeExpense.income) }}</strong></p>
            <p><span>支出</span><strong class="is-expense">{{ formatCurrency(detail.incomeExpense.expense) }}</strong></p>
          </div>
        </article>

        <article class="page-card admin-family-panel">
          <header>
            <span>资产汇总</span>
            <strong>{{ formatCurrency(detail.assetSummary.netAssets) }}</strong>
          </header>
          <div class="admin-family-detail-page__metric-grid">
            <p><span>资产</span><strong>{{ formatCurrency(detail.assetSummary.totalAssets) }}</strong></p>
            <p><span>负债</span><strong>{{ formatCurrency(detail.assetSummary.totalLiabilities) }}</strong></p>
            <p><span>账户数</span><strong>{{ detail.assetSummary.accountCount }} 个</strong></p>
          </div>
        </article>
      </section>

      <section class="admin-family-detail-page__content">
        <article class="page-card admin-family-panel">
          <header>
            <span>成员</span>
            <strong>{{ detail.members.length }} 位</strong>
          </header>
          <div v-if="detail.members.length" class="admin-family-detail-page__members">
            <div v-for="member in detail.members" :key="member.userId" class="admin-family-member">
              <span class="admin-family-member__dot" :style="{ backgroundColor: member.color }"></span>
              <div>
                <strong>{{ member.name }}</strong>
                <p>{{ member.role }} · 加入于 {{ member.joinedAt }}</p>
              </div>
            </div>
          </div>
          <p v-else class="admin-family-detail-page__empty">暂无成员。</p>
        </article>

        <article class="page-card admin-family-panel">
          <header>
            <span>最近账单</span>
            <strong>{{ detail.recentBills.length }} 笔</strong>
          </header>
          <div v-if="detail.recentBills.length" class="admin-family-detail-page__bills">
            <div v-for="bill in detail.recentBills" :key="bill.id" class="admin-family-bill">
              <div>
                <strong>{{ bill.memberName }}</strong>
                <p>{{ bill.recordDate }} · {{ bill.remark || "无备注" }}</p>
              </div>
              <span :class="bill.recordType === 'income' ? 'is-income' : 'is-expense'">
                {{ bill.recordType === "income" ? "+" : "-" }}{{ formatCurrency(bill.amount) }}
              </span>
            </div>
          </div>
          <p v-else class="admin-family-detail-page__empty">暂无最近账单。</p>
        </article>
      </section>
    </template>
  </div>
</template>

<script>
import {
  buildAdminFamilyError,
  formatAdminFamilyCurrency,
  getAdminFamilyDetail,
  normalizeDetail
} from "@/api/adminFamilies";

function createEmptyDetail(familyId) {
  return {
    family: {
      id: String(familyId || ""),
      name: "家庭详情",
      creator: "-",
      memberCount: 0,
      billCount: 0,
      totalAssets: 0,
      createdAt: "-",
      status: "-"
    },
    members: [],
    incomeExpense: {
      income: 0,
      expense: 0,
      balance: 0,
      records: 0
    },
    assetSummary: {
      totalAssets: 0,
      totalLiabilities: 0,
      netAssets: 0,
      accountCount: 0
    },
    recentBills: []
  };
}

export default {
  name: "AdminFamilyDetail",
  data() {
    return {
      detail: createEmptyDetail(this.$route.params.familyId),
      loading: false,
      loadError: "",
      notFound: false
    };
  },
  created() {
    this.fetchDetail();
  },
  watch: {
    "$route.params.familyId": function() {
      this.fetchDetail();
    }
  },
  methods: {
    formatCurrency: formatAdminFamilyCurrency,
    fetchDetail() {
      var familyId = this.$route.params.familyId;
      this.loading = true;
      this.loadError = "";
      this.notFound = false;
      this.detail = createEmptyDetail(familyId);

      return getAdminFamilyDetail(familyId)
        .then(
          function(result) {
            this.detail = normalizeDetail(result, familyId);
          }.bind(this)
        )
        .catch(
          function(error) {
            if (error && error.response && error.response.status === 404) {
              this.notFound = true;
              return;
            }
            this.loadError = buildAdminFamilyError(error, "家庭详情加载失败，请稍后重试");
          }.bind(this)
        )
        .finally(
          function() {
            this.loading = false;
          }.bind(this)
        );
    },
    goBack() {
      this.$router.push("/admin/families");
    }
  }
};
</script>

<style scoped>
.admin-family-detail-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.admin-family-detail-page__header {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 24px;
}

.admin-family-detail-page__eyebrow {
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

.admin-family-detail-page__header h1 {
  margin: 16px 0 10px;
  font-size: 30px;
}

.admin-family-detail-page__header p {
  margin: 0;
  color: var(--text-subtle);
  line-height: 1.7;
}

.admin-family-detail-page__state {
  min-height: 180px;
  padding: 28px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: var(--text-subtle);
  text-align: center;
}

.admin-family-detail-page__state p {
  margin: 0;
}

.admin-family-detail-page__state.is-error {
  color: var(--danger-color);
}

.admin-family-detail-page__summary,
.admin-family-detail-page__metrics,
.admin-family-detail-page__content {
  display: grid;
  gap: 14px;
}

.admin-family-detail-page__summary {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.admin-family-detail-page__metrics,
.admin-family-detail-page__content {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.admin-family-stat,
.admin-family-panel {
  padding: 20px;
}

.admin-family-stat {
  border-radius: 18px;
  border: 1px solid var(--border-color);
  background: rgba(255, 255, 255, 0.95);
}

.admin-family-stat span,
.admin-family-panel header span,
.admin-family-detail-page__metric-grid span {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
}

.admin-family-stat strong {
  display: block;
  margin-top: 8px;
  font-size: 24px;
}

.admin-family-panel header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 14px;
}

.admin-family-panel header strong {
  font-size: 24px;
}

.admin-family-detail-page__metric-grid {
  margin-top: 18px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.admin-family-detail-page__metric-grid p {
  margin: 0;
  padding: 14px;
  border-radius: 14px;
  border: 1px solid var(--border-color);
  background: rgba(245, 246, 248, 0.9);
}

.admin-family-detail-page__metric-grid strong {
  display: block;
  margin-top: 8px;
  font-size: 20px;
}

.admin-family-detail-page__members,
.admin-family-detail-page__bills {
  margin-top: 18px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.admin-family-member,
.admin-family-bill {
  min-height: 58px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  padding: 12px 14px;
  border-radius: 14px;
  border: 1px solid var(--border-color);
  background: rgba(251, 251, 252, 0.94);
}

.admin-family-member {
  justify-content: flex-start;
}

.admin-family-member__dot {
  width: 12px;
  height: 12px;
  border-radius: 999px;
  flex: 0 0 auto;
}

.admin-family-member strong,
.admin-family-bill strong {
  display: block;
}

.admin-family-member p,
.admin-family-bill p,
.admin-family-detail-page__empty {
  margin: 4px 0 0;
  color: var(--text-subtle);
  font-size: 13px;
}

.is-income {
  color: var(--success-color);
}

.is-expense {
  color: var(--danger-color);
}

@media (max-width: 980px) {
  .admin-family-detail-page__header {
    flex-direction: column;
  }

  .admin-family-detail-page__summary,
  .admin-family-detail-page__metrics,
  .admin-family-detail-page__content {
    grid-template-columns: 1fr;
  }
}
</style>
