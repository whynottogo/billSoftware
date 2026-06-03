<template>
  <div class="finance-page admin-families-page">
    <section class="finance-hero finance-hero--soft admin-families-page__hero">
      <div class="finance-hero__eyebrow">
        <span class="finance-hero__eyebrow-mark">家</span>
        <span>管理后台 / 家庭管理</span>
      </div>

      <div class="finance-hero__headline">
        <h1>统一查看家庭规模与资产分布</h1>
        <p>家庭统计、成员规模、账单数量和资产汇总均来自真实后台接口。</p>
      </div>

      <div class="finance-stat-grid">
        <article v-for="item in statCards" :key="item.label" class="finance-stat-card">
          <span>{{ item.label }}</span>
          <strong>{{ item.value }}</strong>
          <small>{{ item.note }}</small>
        </article>
      </div>
    </section>

    <section class="page-card finance-panel admin-families-page__filters">
      <label>
        <span>搜索家庭</span>
        <input
          v-model.trim="keywordInput"
          type="text"
          placeholder="输入家庭名称、家庭 ID 或创建人"
          @keyup.enter="applySearch"
        />
      </label>
      <div class="admin-families-page__filter-actions">
        <button class="finance-button finance-button--primary" :disabled="loadingList" @click="applySearch">搜索</button>
        <button class="finance-button finance-button--ghost" :disabled="loadingList && !keyword" @click="resetSearch">重置</button>
      </div>
    </section>

    <section v-if="loadError" class="page-card finance-panel admin-families-page__state is-error">
      <p>{{ loadError }}</p>
      <button class="finance-button finance-button--ghost" @click="loadAll">重新加载</button>
    </section>

    <section v-else-if="loadingList && !families.length" class="page-card finance-panel admin-families-page__state">
      正在加载家庭列表...
    </section>

    <section v-else-if="!families.length" class="page-card finance-panel admin-families-page__state">
      <strong>当前没有匹配的家庭</strong>
      <p>可以调整搜索条件后重新查询。</p>
    </section>

    <section v-else class="finance-grid-2">
      <article v-for="family in families" :key="family.id" class="page-card finance-panel admin-families-page__family-card">
        <header class="admin-families-page__family-header">
          <div>
            <h3>{{ family.name }}</h3>
            <p>创建者：{{ family.creator }} · {{ family.id }}</p>
          </div>
          <span :class="['admin-families-page__status', family.status === '活跃' ? 'is-active' : '']">
            {{ family.status }}
          </span>
        </header>

        <div class="admin-families-page__metrics">
          <article>
            <span>成员数量</span>
            <strong>{{ family.memberCount }} 位</strong>
          </article>
          <article>
            <span>账单数量</span>
            <strong>{{ family.billCount }} 笔</strong>
          </article>
        </div>

        <footer class="admin-families-page__footer">
          <div>
            <span>家庭净资产</span>
            <strong>{{ formatCurrency(family.totalAssets) }}</strong>
            <small>创建于 {{ family.createdAt }}</small>
          </div>
          <div class="admin-families-page__actions">
            <button class="finance-button finance-button--ghost" @click="openDetail(family)">查看</button>
            <button class="finance-button finance-button--ghost" @click="showMorePlaceholder(family)">更多</button>
          </div>
        </footer>
      </article>
    </section>

    <section class="page-card finance-panel admin-families-page__pagination">
      <div>
        <strong>共 {{ pagination.total }} 个家庭</strong>
        <p>当前第 {{ pagination.page }} / {{ totalPages }} 页</p>
      </div>
      <div class="admin-families-page__page-buttons">
        <button class="finance-button finance-button--ghost" :disabled="!canPrev || loadingList" @click="changePage(pagination.page - 1)">上一页</button>
        <button
          v-for="page in pageButtons"
          :key="page"
          :class="['finance-button', page === pagination.page ? 'finance-button--primary' : 'finance-button--ghost']"
          :disabled="loadingList"
          @click="changePage(page)"
        >
          {{ page }}
        </button>
        <button class="finance-button finance-button--ghost" :disabled="!canNext || loadingList" @click="changePage(pagination.page + 1)">下一页</button>
      </div>
    </section>
  </div>
</template>

<script>
import {
  buildAdminFamilyError,
  formatAdminFamilyCurrency,
  getAdminFamilySummary,
  listAdminFamilies,
  normalizeList,
  normalizeSummary
} from "@/api/adminFamilies";

function defaultSummary() {
  return {
    totalFamilies: 0,
    totalMembers: 0,
    averageMembers: 0,
    activeFamilies: 0,
    activeDays: 30
  };
}

export default {
  name: "AdminFamilies",
  data() {
    return {
      summary: defaultSummary(),
      families: [],
      pagination: {
        page: 1,
        pageSize: 10,
        total: 0,
        totalPages: 0
      },
      keyword: "",
      keywordInput: "",
      loadingSummary: false,
      loadingList: false,
      loadError: ""
    };
  },
  computed: {
    totalPages() {
      return this.pagination.totalPages || 1;
    },
    canPrev() {
      return this.pagination.page > 1;
    },
    canNext() {
      return this.pagination.page < this.totalPages;
    },
    statCards() {
      return [
        { label: "家庭总数", value: this.summary.totalFamilies, note: "真实家庭数量" },
        { label: "总成员数", value: this.summary.totalMembers, note: "已加入家庭成员" },
        { label: "平均成员数", value: this.summary.averageMembers, note: "每个家庭平均规模" },
        { label: "活跃家庭数", value: this.summary.activeFamilies, note: "近 " + this.summary.activeDays + " 天有账单" }
      ];
    },
    pageButtons() {
      var total = this.totalPages;
      var current = this.pagination.page;
      var start = Math.max(1, current - 2);
      var end = Math.min(total, start + 4);
      start = Math.max(1, end - 4);
      var pages = [];
      for (var page = start; page <= end; page += 1) {
        pages.push(page);
      }
      return pages;
    }
  },
  created() {
    this.loadAll();
  },
  methods: {
    formatCurrency: formatAdminFamilyCurrency,
    loadAll() {
      return Promise.all([this.fetchSummary(), this.fetchFamilies()]);
    },
    fetchSummary() {
      this.loadingSummary = true;
      return getAdminFamilySummary()
        .then(
          function(result) {
            this.summary = normalizeSummary(result);
          }.bind(this)
        )
        .catch(
          function(error) {
            this.loadError = buildAdminFamilyError(error, "家庭统计加载失败，请检查后端服务");
          }.bind(this)
        )
        .finally(
          function() {
            this.loadingSummary = false;
          }.bind(this)
        );
    },
    fetchFamilies() {
      this.loadingList = true;
      this.loadError = "";
      return listAdminFamilies({
        page: this.pagination.page,
        page_size: this.pagination.pageSize,
        keyword: this.keyword
      })
        .then(
          function(result) {
            var payload = normalizeList(result);
            this.families = payload.families;
            this.pagination = Object.assign({}, this.pagination, payload.pagination);
          }.bind(this)
        )
        .catch(
          function(error) {
            this.families = [];
            this.loadError = buildAdminFamilyError(error, "家庭列表加载失败，请稍后重试");
          }.bind(this)
        )
        .finally(
          function() {
            this.loadingList = false;
          }.bind(this)
        );
    },
    applySearch() {
      this.keyword = this.keywordInput;
      this.pagination.page = 1;
      this.fetchFamilies();
    },
    resetSearch() {
      this.keyword = "";
      this.keywordInput = "";
      this.pagination.page = 1;
      this.fetchFamilies();
    },
    changePage(page) {
      if (page < 1 || page > this.totalPages || page === this.pagination.page) {
        return;
      }
      this.pagination.page = page;
      this.fetchFamilies();
    },
    openDetail(family) {
      if (!family || !family.id) {
        return;
      }
      this.$router.push("/admin/families/" + encodeURIComponent(family.id));
    },
    showMorePlaceholder(family) {
      this.$message.info((family ? family.name : "家庭") + " 的更多操作将在后续批次开放");
    }
  }
};
</script>

<style scoped>
.admin-families-page {
  gap: 20px;
}

.admin-families-page__hero {
  padding: 28px;
}

.admin-families-page__filters {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
  padding: 18px;
}

.admin-families-page__filters label {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
  color: var(--text-subtle);
  font-size: 13px;
  font-weight: 700;
}

.admin-families-page__filters input {
  height: 44px;
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 0 14px;
  background: rgba(251, 251, 252, 0.92);
  color: var(--text-main);
}

.admin-families-page__filter-actions {
  display: flex;
  gap: 10px;
}

.admin-families-page__state {
  min-height: 150px;
  padding: 28px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: var(--text-subtle);
  text-align: center;
}

.admin-families-page__state p {
  margin: 0;
}

.admin-families-page__state.is-error {
  color: var(--danger-color);
}

.admin-families-page__family-card {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.admin-families-page__family-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.admin-families-page__family-header h3 {
  margin: 0;
  font-size: 22px;
}

.admin-families-page__family-header p {
  margin: 8px 0 0;
  color: var(--text-subtle);
}

.admin-families-page__status {
  min-height: 30px;
  padding: 0 12px;
  border-radius: 999px;
  background: rgba(245, 158, 11, 0.14);
  color: #b45309;
  font-size: 12px;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  white-space: nowrap;
}

.admin-families-page__status.is-active {
  background: rgba(34, 197, 94, 0.14);
  color: #15803d;
}

.admin-families-page__metrics {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.admin-families-page__metrics article {
  padding: 14px 16px;
  border-radius: 16px;
  border: 1px solid var(--border-color);
  background: rgba(251, 251, 252, 0.92);
}

.admin-families-page__metrics span {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
}

.admin-families-page__metrics strong {
  display: block;
  margin-top: 8px;
  font-size: 22px;
}

.admin-families-page__footer {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 12px;
}

.admin-families-page__footer span,
.admin-families-page__footer small {
  display: block;
}

.admin-families-page__footer span {
  color: var(--text-muted);
  font-size: 12px;
}

.admin-families-page__footer strong {
  display: block;
  margin-top: 8px;
  font-size: 28px;
}

.admin-families-page__footer small {
  margin-top: 8px;
  color: var(--text-subtle);
  font-size: 12px;
}

.admin-families-page__actions {
  display: flex;
  gap: 10px;
}

.admin-families-page__actions .finance-button,
.admin-families-page__page-buttons .finance-button {
  min-height: 38px;
  padding: 0 14px;
  border-radius: 12px;
}

.admin-families-page__pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.admin-families-page__pagination p {
  margin: 8px 0 0;
  color: var(--text-subtle);
}

.admin-families-page__page-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

@media (max-width: 1080px) {
  .admin-families-page__hero .finance-stat-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 920px) {
  .admin-families-page__filters,
  .admin-families-page__footer,
  .admin-families-page__pagination {
    flex-direction: column;
    align-items: flex-start;
  }

  .admin-families-page__filter-actions {
    width: 100%;
    flex-wrap: wrap;
  }
}
</style>
