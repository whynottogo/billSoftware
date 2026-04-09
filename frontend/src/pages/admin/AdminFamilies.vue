<template>
  <div class="finance-page admin-families-page">
    <section class="finance-hero finance-hero--soft admin-families-page__hero">
      <div class="finance-hero__eyebrow">
        <span class="finance-hero__eyebrow-mark">家</span>
        <span>管理后台 / 家庭管理</span>
      </div>

      <div class="finance-hero__headline">
        <h1>统一查看家庭规模与资产分布</h1>
        <p>基于 Make 的家庭管理原型重建后台页面，先以静态数据完成结构承载。</p>
      </div>

      <div class="finance-stat-grid">
        <article v-for="item in dataSource.stats" :key="item.label" class="finance-stat-card">
          <span>{{ item.label }}</span>
          <strong>{{ item.value }}</strong>
          <small>门户阶段 mock 统计</small>
        </article>
      </div>
    </section>

    <section class="finance-grid-2">
      <article
        v-for="family in dataSource.families"
        :key="family.id"
        class="page-card finance-panel admin-families-page__family-card"
      >
        <header class="admin-families-page__family-header">
          <div>
            <h3>{{ family.name }}</h3>
            <p>创建者：{{ family.creator }}</p>
          </div>
          <span class="admin-families-page__status">{{ family.status }}</span>
        </header>

        <div class="admin-families-page__metrics">
          <article>
            <span>成员数量</span>
            <strong>{{ family.members }} 位</strong>
          </article>
          <article>
            <span>账单数量</span>
            <strong>{{ family.billCount }} 笔</strong>
          </article>
        </div>

        <footer class="admin-families-page__footer">
          <div>
            <span>家庭总资产</span>
            <strong>{{ formatAdminCurrency(family.totalAssets) }}</strong>
            <small>创建于 {{ family.createDate }}</small>
          </div>
          <div class="admin-families-page__actions">
            <button class="finance-button finance-button--ghost" @click="preview('查看详情', family)">查看</button>
            <button class="finance-button finance-button--ghost" @click="preview('更多操作', family)">更多</button>
          </div>
        </footer>
      </article>
    </section>

    <section class="page-card finance-panel admin-families-page__pagination">
      <div>
        <strong>共 {{ dataSource.pagination.total }} 个家庭</strong>
        <p>当前第 {{ dataSource.pagination.currentPage }} 页</p>
      </div>
      <div class="admin-families-page__page-buttons">
        <button class="finance-button finance-button--ghost" @click="preview('上一页', null)">上一页</button>
        <button
          v-for="page in dataSource.pagination.pages"
          :key="page"
          :class="['finance-button', page === dataSource.pagination.currentPage ? 'finance-button--primary' : 'finance-button--ghost']"
          @click="preview('切换页码 ' + page, null)"
        >
          {{ page }}
        </button>
        <button class="finance-button finance-button--ghost" @click="preview('下一页', null)">下一页</button>
      </div>
    </section>
  </div>
</template>

<script>
import { formatAdminCurrency, getFamilyManagementData } from "@/utils/adminPortalMock";

export default {
  name: "AdminFamilies",
  data() {
    return {
      dataSource: getFamilyManagementData()
    };
  },
  methods: {
    formatAdminCurrency: formatAdminCurrency,
    preview(action, family) {
      if (family) {
        this.$message.info(action + "：" + family.name + "（原型操作）");
        return;
      }

      this.$message.info(action + "（原型操作）");
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
  background: rgba(34, 197, 94, 0.14);
  color: #15803d;
  font-size: 12px;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
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

.admin-families-page__actions .finance-button {
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

.admin-families-page__page-buttons .finance-button {
  min-height: 38px;
  padding: 0 14px;
  border-radius: 12px;
}

@media (max-width: 1080px) {
  .admin-families-page__hero .finance-stat-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 920px) {
  .admin-families-page__footer,
  .admin-families-page__pagination {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
