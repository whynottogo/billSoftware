<template>
  <section class="page-card family-share-panel">
    <header class="family-share-panel__header">
      <div>
        <h3>{{ title }}</h3>
        <p>点击月收入、月支出、年收入、年支出后，这里会同步展示成员占比。</p>
      </div>
      <strong>{{ totalLabel }}</strong>
    </header>

    <div v-if="rows.length" class="family-share-panel__content">
      <SimpleDonutChart :segments="segments" :center-title="'总计'" :center-value="totalLabel" :size="220" />

      <div class="family-share-panel__list">
        <article v-for="item in rows" :key="item.name" class="family-share-panel__row">
          <div class="family-share-panel__member">
            <span class="family-share-panel__dot" :style="{ backgroundColor: item.color }"></span>
            <div>
              <strong>{{ item.name }}</strong>
              <p>{{ item.role }}</p>
            </div>
          </div>
          <div class="family-share-panel__value">
            <strong>{{ formatCurrency(item.value) }}</strong>
            <span>{{ formatPercent(item.percent) }}</span>
          </div>
        </article>
      </div>
    </div>

    <el-empty v-else description="当前口径暂无可展示的成员占比数据" />
  </section>
</template>

<script>
import SimpleDonutChart from "@/components/SimpleDonutChart.vue";
import { formatFamilyCurrency, formatFamilyPercent } from "@/api/userFamily";

export default {
  name: "FamilyMemberSharePanel",
  components: {
    SimpleDonutChart
  },
  props: {
    title: {
      type: String,
      default: "成员占比"
    },
    total: {
      type: Number,
      default: 0
    },
    rows: {
      type: Array,
      default: function() {
        return [];
      }
    }
  },
  computed: {
    totalLabel() {
      return formatFamilyCurrency(this.total);
    },
    segments() {
      return this.rows.map(function(item) {
        return {
          name: item.name,
          value: item.value,
          color: item.color
        };
      });
    }
  },
  methods: {
    formatCurrency: formatFamilyCurrency,
    formatPercent: formatFamilyPercent
  }
};
</script>

<style scoped>
.family-share-panel {
  padding: 24px;
}

.family-share-panel__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.family-share-panel__header h3 {
  margin: 0;
  font-size: 22px;
}

.family-share-panel__header p {
  margin: 10px 0 0;
  color: var(--text-subtle);
  line-height: 1.6;
}

.family-share-panel__header strong {
  font-size: 26px;
  line-height: 1.2;
}

.family-share-panel__content {
  margin-top: 16px;
  display: grid;
  grid-template-columns: 250px minmax(0, 1fr);
  gap: 18px;
}

.family-share-panel__list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.family-share-panel__row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  padding: 16px;
  border-radius: 16px;
  border: 1px solid var(--border-color);
  background: rgba(255, 255, 255, 0.9);
}

.family-share-panel__member {
  display: flex;
  align-items: center;
  gap: 12px;
}

.family-share-panel__dot {
  width: 12px;
  height: 12px;
  border-radius: 999px;
}

.family-share-panel__member strong {
  display: block;
  font-size: 15px;
}

.family-share-panel__member p {
  margin: 4px 0 0;
  color: var(--text-muted);
  font-size: 12px;
}

.family-share-panel__value {
  text-align: right;
}

.family-share-panel__value strong {
  display: block;
  font-size: 18px;
}

.family-share-panel__value span {
  display: block;
  margin-top: 4px;
  color: var(--text-subtle);
  font-size: 12px;
}

@media (max-width: 960px) {
  .family-share-panel__content {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
