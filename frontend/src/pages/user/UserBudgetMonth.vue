<template>
  <div class="finance-page budget-month-page">
    <div v-if="errorMessage" class="finance-inline-notice">
      <span>{{ errorMessage }}</span>
      <button type="button" @click="loadBudget">重新加载</button>
    </div>

    <section class="finance-toolbar">
      <div class="finance-toolbar__meta">
        <h1 class="page-title">{{ budgetData.label }}</h1>
        <p class="page-description">{{ budgetData.notice }}</p>
      </div>

      <div class="finance-toolbar__actions">
        <button class="finance-button finance-button--ghost" type="button" @click="openYearBudget">查看年预算</button>
        <button
          class="finance-button finance-button--primary"
          type="button"
          :disabled="!canEditBudget || isLoading"
          @click="openBudgetDialog"
        >
          {{ budgetActionLabel }}
        </button>
      </div>
    </section>

    <section class="finance-hero finance-hero--soft">
      <div class="finance-hero__eyebrow">
        <span class="finance-hero__eyebrow-mark">预</span>
        <span>本月预算总览</span>
      </div>

      <div class="finance-hero__headline">
        <h2>预算页不是另一张账单页，它应该直接告诉你哪里还能花、哪里该收手</h2>
        <p>{{ isLoading ? "正在同步当前月份预算数据…" : heroDescription }}</p>
      </div>

      <div class="finance-stat-grid finance-stat-grid--triple">
        <article class="finance-stat-card">
          <span>本月预算</span>
          <strong>{{ formatCurrency(budgetData.overview.totalBudget) }}</strong>
          <small>围绕支出类型分配额度</small>
        </article>
        <article class="finance-stat-card">
          <span>本月支出</span>
          <strong class="finance-tone-expense">{{ formatCurrency(budgetData.overview.totalExpense) }}</strong>
          <small>当前月度已发生支出</small>
        </article>
        <article class="finance-stat-card">
          <span>剩余预算</span>
          <strong class="finance-tone-income">{{ formatCurrency(budgetData.overview.remaining) }}</strong>
          <small>仍可主动安排的空间</small>
        </article>
      </div>

      <div class="budget-month-page__hero-progress">
        <div class="budget-month-page__hero-progress-meta">
          <span>预算使用进度</span>
          <strong>{{ formatPercent(budgetData.overview.percentage) }}</strong>
        </div>
        <div class="finance-progress">
          <span class="budget-month-page__hero-bar" :style="{ width: progressWidth(budgetData.overview.percentage) }"></span>
        </div>
      </div>
    </section>

    <section class="finance-grid-3">
      <article v-for="item in budgetData.highlights" :key="item.label" class="page-card budget-month-page__highlight">
        <span>{{ item.label }}</span>
        <strong>{{ item.value }}</strong>
        <p>{{ item.hint }}</p>
      </article>
    </section>

    <section class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>分类预算</h3>
          <p>{{ categoryPanelDescription }}</p>
        </div>
      </header>

      <div v-if="isLoading && !budgetData.categories.length" class="budget-month-page__empty">
        <strong>正在加载本月预算分类…</strong>
        <p>预算与支出明细同步后会显示在这里。</p>
      </div>

      <div v-else-if="!budgetData.categories.length" class="budget-month-page__empty">
        <strong>当前还没有可展示的预算分类</strong>
        <p>设置本月预算后，这里会按支出分类显示预算、已用和剩余额度。</p>
      </div>

      <div v-else class="budget-month-page__categories">
        <article v-for="item in budgetData.categories" :key="item.name" class="budget-card">
          <div class="budget-card__header">
            <div class="budget-card__title">
              <span class="finance-badge">{{ item.badge }}</span>
              <div>
                <strong>{{ item.name }}</strong>
                <p>{{ item.note }}</p>
              </div>
            </div>

            <button
              class="finance-button finance-button--ghost budget-card__action"
              type="button"
              :disabled="!canEditBudget || isLoading"
              @click="openBudgetDialog"
            >
              调整
            </button>
          </div>

          <div class="budget-card__numbers">
            <div>
              <span>预算</span>
              <strong>{{ formatCurrency(item.budget) }}</strong>
            </div>
            <div>
              <span>已用</span>
              <strong class="finance-tone-expense">{{ formatCurrency(item.expense) }}</strong>
            </div>
            <div>
              <span>{{ item.remaining >= 0 ? "剩余" : "超支" }}</span>
              <strong :class="statusTone(item.status)">
                {{ item.remaining >= 0 ? formatCurrency(item.remaining) : formatCurrency(Math.abs(item.remaining)) }}
              </strong>
            </div>
          </div>

          <div class="budget-card__progress">
            <div class="budget-card__progress-meta">
              <span>{{ statusLabel(item.status) }}</span>
              <strong :class="statusTone(item.status)">{{ formatPercent(item.percentage) }}</strong>
            </div>
            <div class="finance-progress">
              <span :class="statusBar(item.status)" :style="{ width: progressWidth(item.percentage) }"></span>
            </div>
          </div>
        </article>
      </div>
    </section>
  </div>

  <el-dialog
    v-model="budgetDialogVisible"
    :title="budgetActionLabel"
    width="680px"
    destroy-on-close
  >
    <form class="budget-form" @submit.prevent="submitBudget">
      <label class="budget-form__field">
        <span>本月总预算</span>
        <input
          v-model.number="budgetForm.totalBudget"
          class="budget-form__input"
          type="number"
          min="0"
          step="0.01"
          placeholder="请输入本月总预算"
        />
      </label>

      <div class="budget-form__section">
        <div class="budget-form__section-head">
          <div>
            <strong>分类预算</strong>
            <p>按支出分类分别设置额度，保存后会即时回刷当前页面。</p>
          </div>
        </div>

        <div v-if="!budgetForm.categories.length" class="budget-form__empty">
          当前没有可设置的分类预算，后端返回分类后即可在这里编辑。
        </div>

        <div v-else class="budget-form__category-list">
          <label v-for="item in budgetForm.categories" :key="item.key" class="budget-form__category-item">
            <div class="budget-form__category-meta">
              <div class="budget-form__category-title">
                <span class="finance-badge">{{ item.badge }}</span>
                <div>
                  <strong>{{ item.name }}</strong>
                  <p>已用 {{ formatCurrency(item.expense) }}</p>
                </div>
              </div>
              <span class="budget-form__category-remaining">
                {{ item.remaining >= 0 ? "剩余" : "超支" }}
                {{ formatCurrency(Math.abs(item.remaining)) }}
              </span>
            </div>
            <input
              v-model.number="item.budget"
              class="budget-form__input"
              type="number"
              min="0"
              step="0.01"
              :placeholder="'设置 ' + item.name + ' 预算'"
            />
          </label>
        </div>
      </div>

      <footer class="budget-form__footer">
        <button class="budget-form__button budget-form__button--ghost" type="button" @click="budgetDialogVisible = false">
          取消
        </button>
        <button class="budget-form__button budget-form__button--primary" type="submit" :disabled="isSaving">
          {{ isSaving ? "保存中..." : "保存预算" }}
        </button>
      </footer>
    </form>
  </el-dialog>
</template>

<script>
import { getCurrentMonthBudget, updateCurrentMonthBudget } from "@/api/userBudget";

function pad(value) {
  return String(value).padStart(2, "0");
}

function getCurrentMonthKey() {
  var now = new Date();
  return now.getFullYear() + "-" + pad(now.getMonth() + 1);
}

function getCurrentMonthLabel(monthKey) {
  var parts = String(monthKey || "").split("-");

  if (parts.length !== 2) {
    return "本月预算";
  }

  return parts[0] + "年" + Number(parts[1]) + "月预算";
}

function toNumber(value) {
  var parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : 0;
}

function toBadge(name) {
  return String(name || "预").slice(0, 1);
}

function formatBudgetCurrency(value) {
  return "¥" + Number(value || 0).toLocaleString("zh-CN");
}

function formatBudgetPercent(value) {
  var parsed = Number(value || 0).toFixed(1);
  return parsed.replace(".0", "") + "%";
}

function buildBudgetError(error, fallback) {
  if (error && error.response && error.response.data && error.response.data.message) {
    return error.response.data.message;
  }

  return fallback;
}

function extractPayload(result) {
  if (result && typeof result === "object" && result.data && typeof result.data === "object") {
    return result.data;
  }

  return result || {};
}

function resolveStatus(percentage) {
  if (percentage >= 100) {
    return "over";
  }

  if (percentage >= 80) {
    return "warning";
  }

  return "safe";
}

function buildCategoryNote(item) {
  if (item.note) {
    return item.note;
  }

  if (item.remaining < 0) {
    return "当前分类预算已经超支，建议本月优先收口。";
  }

  if (item.percentage >= 80) {
    return "当前分类预算接近上限，后续支出需更谨慎。";
  }

  if (item.expense <= 0) {
    return "本月暂未发生该分类支出。";
  }

  return "当前分类预算仍有余量，可继续按计划安排。";
}

function buildMonthHighlights(categories, overview) {
  var overCount = 0;
  var safeCount = 0;
  var warningCount = 0;

  categories.forEach(function(item) {
    if (item.status === "over") {
      overCount += 1;
      return;
    }

    if (item.status === "warning") {
      warningCount += 1;
      return;
    }

    safeCount += 1;
  });

  return [
    {
      label: "已超支分类",
      value: overCount + "个",
      hint: overCount ? "建议优先回收超支分类。" : "当前没有分类超支。"
    },
    {
      label: "可控分类",
      value: safeCount + "个",
      hint: safeCount ? "仍有余量的分类可以继续跟踪。" : "当前没有明显宽松的分类。"
    },
    {
      label: "预算节奏",
      value: formatBudgetPercent(overview.percentage),
      hint: warningCount ? "有 " + warningCount + " 个分类接近上限。" : "当前预算节奏整体平稳。"
    }
  ];
}

function mapMonthCategory(item, index) {
  var budget = toNumber(item.budget || item.totalBudget);
  var expense = toNumber(item.expense || item.totalExpense || item.usedAmount);
  var remaining = Object.prototype.hasOwnProperty.call(item, "remaining") ? toNumber(item.remaining) : budget - expense;
  var percentage = Object.prototype.hasOwnProperty.call(item, "percentage")
    ? toNumber(item.percentage)
    : budget > 0
      ? Number(((expense / budget) * 100).toFixed(1))
      : expense > 0
        ? 100
        : 0;
  var status = item.status || resolveStatus(percentage);

  return {
    key: item.id || item.categoryId || item.category_id || "category-" + index,
    id: item.id || item.categoryId || item.category_id || "",
    name: item.name || "未命名分类",
    badge: item.badge || toBadge(item.name),
    note: buildCategoryNote({
      note: item.note,
      remaining: remaining,
      percentage: percentage,
      expense: expense
    }),
    budget: budget,
    expense: expense,
    remaining: remaining,
    percentage: percentage,
    status: status
  };
}

function createEmptyBudgetData(monthKey) {
  var resolvedMonth = monthKey || getCurrentMonthKey();
  var overview = {
    totalBudget: 0,
    totalExpense: 0,
    remaining: 0,
    percentage: 0
  };

  return {
    month: resolvedMonth,
    label: getCurrentMonthLabel(resolvedMonth),
    notice: "仅支持设置当前月份预算，保存后会即时刷新当前页。",
    overview: overview,
    highlights: buildMonthHighlights([], overview),
    categories: [],
    canEdit: true
  };
}

function normalizeMonthBudgetPayload(result) {
  var payload = extractPayload(result);
  var month = payload.month || getCurrentMonthKey();
  var overviewSource = payload.overview || {};
  var overview = {
    totalBudget: toNumber(overviewSource.totalBudget || overviewSource.budget || payload.totalBudget),
    totalExpense: toNumber(overviewSource.totalExpense || overviewSource.expense || payload.totalExpense),
    remaining: 0,
    percentage: 0
  };

  overview.remaining = Object.prototype.hasOwnProperty.call(overviewSource, "remaining")
    ? toNumber(overviewSource.remaining)
    : overview.totalBudget - overview.totalExpense;
  overview.percentage = Object.prototype.hasOwnProperty.call(overviewSource, "percentage")
    ? toNumber(overviewSource.percentage)
    : overview.totalBudget > 0
      ? Number(((overview.totalExpense / overview.totalBudget) * 100).toFixed(1))
      : 0;

  var categories = Array.isArray(payload.categories) ? payload.categories.map(mapMonthCategory) : [];
  var highlights = Array.isArray(payload.highlights) && payload.highlights.length
    ? payload.highlights.map(function(item) {
        return {
          label: item.label || "",
          value: item.value || "",
          hint: item.hint || ""
        };
      })
    : buildMonthHighlights(categories, overview);
  var rawEditable = payload.canEdit;

  if (typeof rawEditable !== "boolean") {
    rawEditable = payload.can_edit;
  }

  if (typeof rawEditable !== "boolean") {
    rawEditable = payload.editable;
  }

  if (typeof rawEditable !== "boolean") {
    rawEditable = payload.isEditable;
  }

  return {
    month: month,
    label: payload.label || getCurrentMonthLabel(month),
    notice: payload.notice || "仅支持设置当前月份预算，保存后会即时刷新当前页。",
    overview: overview,
    highlights: highlights,
    categories: categories,
    canEdit: typeof rawEditable === "boolean" ? rawEditable : month === getCurrentMonthKey()
  };
}

export default {
  name: "UserBudgetMonth",
  data() {
    return {
      isLoading: false,
      isSaving: false,
      errorMessage: "",
      budgetDialogVisible: false,
      budgetData: createEmptyBudgetData(),
      budgetForm: {
        totalBudget: 0,
        categories: []
      }
    };
  },
  computed: {
    canEditBudget() {
      return Boolean(this.budgetData.canEdit);
    },
    budgetActionLabel() {
      return this.budgetData.overview.totalBudget > 0 ? "调整预算" : "设置预算";
    },
    heroDescription() {
      if (!this.budgetData.categories.length) {
        return "当前还没有设置细分预算，可以先设置总预算，再按支出分类补充分配。";
      }

      return "本月支出已进入后半程，但还有足够余量可以精细调整。";
    },
    categoryPanelDescription() {
      if (!this.canEditBudget) {
        return "当前月份预算处于只读状态，仅支持查看预算与实际支出。";
      }

      return "每一类预算都单独看状态，这样才能知道要不要继续压缩。";
    }
  },
  created() {
    this.loadBudget();
  },
  methods: {
    formatCurrency: formatBudgetCurrency,
    formatPercent: formatBudgetPercent,
    loadBudget() {
      this.isLoading = true;
      this.errorMessage = "";

      return getCurrentMonthBudget()
        .then(
          function(result) {
            this.budgetData = normalizeMonthBudgetPayload(result);
          }.bind(this)
        )
        .catch(
          function(error) {
            this.budgetData = createEmptyBudgetData();
            this.errorMessage = buildBudgetError(error, "月预算加载失败，请稍后重试。");
            this.$message.error(this.errorMessage);
          }.bind(this)
        )
        .finally(
          function() {
            this.isLoading = false;
          }.bind(this)
        );
    },
    openYearBudget() {
      this.$router.push("/user/budget/year");
    },
    openBudgetDialog() {
      if (!this.canEditBudget) {
        this.$message.info("当前月份预算为只读状态，暂不支持编辑。");
        return;
      }

      this.budgetForm = {
        totalBudget: this.budgetData.overview.totalBudget,
        categories: this.budgetData.categories.map(function(item) {
          return {
            key: item.key,
            id: item.id,
            name: item.name,
            badge: item.badge,
            expense: item.expense,
            remaining: item.remaining,
            budget: item.budget
          };
        })
      };
      this.budgetDialogVisible = true;
    },
    buildUpdatePayload() {
      var categoryBudgets = this.budgetForm.categories.filter(function(item) {
        return item.id !== "" && item.id !== null && item.id !== undefined;
      }).map(function(item) {
        return {
          category_id: Number(item.id),
          categoryId: Number(item.id),
          amount: toNumber(item.budget),
          budget: toNumber(item.budget),
          name: item.name
        };
      });
      var totalBudget = toNumber(this.budgetForm.totalBudget);

      return {
        month: this.budgetData.month,
        total_amount: totalBudget,
        totalAmount: totalBudget,
        total_budget: totalBudget,
        totalBudget: totalBudget,
        categories: categoryBudgets,
        items: categoryBudgets,
        category_budgets: categoryBudgets,
        categoryBudgets: categoryBudgets
      };
    },
    submitBudget() {
      var invalidCategory = this.budgetForm.categories.some(function(item) {
        return Number(item.budget) < 0;
      });

      if (Number(this.budgetForm.totalBudget) < 0) {
        this.$message.warning("本月总预算不能小于 0");
        return;
      }

      if (invalidCategory) {
        this.$message.warning("分类预算不能小于 0");
        return;
      }

      this.isSaving = true;

      updateCurrentMonthBudget(this.buildUpdatePayload())
        .then(
          function() {
            this.budgetDialogVisible = false;
            return this.loadBudget();
          }.bind(this)
        )
        .then(
          function() {
            this.$message.success("本月预算已更新。");
          }.bind(this)
        )
        .catch(
          function(error) {
            this.$message.error(buildBudgetError(error, "本月预算保存失败，请稍后重试。"));
          }.bind(this)
        )
        .finally(
          function() {
            this.isSaving = false;
          }.bind(this)
        );
    },
    statusLabel(status) {
      var labels = {
        over: "已超支",
        warning: "接近上限",
        safe: "仍有余量"
      };

      return labels[status] || "正常";
    },
    statusTone(status) {
      var tones = {
        over: "finance-tone-expense",
        warning: "finance-tone-warning",
        safe: "finance-tone-income"
      };

      return tones[status] || "";
    },
    statusBar(status) {
      var bars = {
        over: "budget-card__bar--danger",
        warning: "budget-card__bar--warning",
        safe: "budget-card__bar--safe"
      };

      return bars[status] || "budget-card__bar--safe";
    },
    progressWidth(value) {
      return Math.min(Math.max(Number(value || 0), 0), 100) + "%";
    }
  }
};
</script>

<style scoped>
.finance-inline-notice {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 18px;
  border: 1px solid rgba(239, 68, 68, 0.18);
  background: rgba(255, 244, 244, 0.92);
  color: #9f2f2f;
}

.finance-inline-notice button {
  border: none;
  background: transparent;
  color: inherit;
  font-weight: 700;
  cursor: pointer;
}

.budget-month-page__hero-progress {
  margin-top: 24px;
  padding: 20px 22px;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.88);
}

.budget-month-page__hero-progress-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.budget-month-page__hero-progress-meta span {
  color: var(--text-subtle);
}

.budget-month-page__hero-progress-meta strong {
  font-size: 22px;
}

.budget-month-page__hero-bar {
  background: linear-gradient(90deg, #22c55e 0%, #f6d34a 100%);
}

.budget-month-page__highlight {
  padding: 22px;
}

.budget-month-page__highlight span {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
}

.budget-month-page__highlight strong {
  display: block;
  margin-top: 12px;
  font-size: 26px;
}

.budget-month-page__highlight p {
  margin: 10px 0 0;
  color: var(--text-subtle);
  line-height: 1.6;
}

.budget-month-page__categories {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.budget-month-page__empty {
  padding: 24px;
  border-radius: 22px;
  background: rgba(247, 247, 248, 0.92);
}

.budget-month-page__empty strong {
  display: block;
  font-size: 17px;
}

.budget-month-page__empty p {
  margin: 10px 0 0;
  color: var(--text-subtle);
  line-height: 1.6;
}

.budget-card {
  padding: 22px;
  border-radius: 24px;
  border: 1px solid var(--border-color);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98) 0%, rgba(251, 251, 252, 0.96) 100%);
}

.budget-card__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.budget-card__title {
  display: flex;
  gap: 14px;
}

.budget-card__title strong {
  display: block;
  font-size: 18px;
}

.budget-card__title p {
  margin: 8px 0 0;
  color: var(--text-subtle);
  line-height: 1.6;
}

.budget-card__action {
  min-height: 38px;
  padding: 0 14px;
}

.budget-card__numbers {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
  margin-top: 20px;
}

.budget-card__numbers span {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
}

.budget-card__numbers strong {
  display: block;
  margin-top: 8px;
  font-size: 18px;
}

.budget-card__progress {
  margin-top: 20px;
}

.budget-card__progress-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.budget-card__bar--danger {
  background: var(--danger-color);
}

.budget-card__bar--warning {
  background: var(--warning-color);
}

.budget-card__bar--safe {
  background: var(--success-color);
}

.budget-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.budget-form__field,
.budget-form__category-item {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.budget-form__field span,
.budget-form__category-item span {
  color: var(--text-subtle);
}

.budget-form__input {
  width: 100%;
  min-height: 44px;
  padding: 0 14px;
  border-radius: 14px;
  border: 1px solid var(--border-color);
  background: rgba(255, 255, 255, 0.96);
  color: var(--text-main);
  font-size: 15px;
}

.budget-form__section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.budget-form__section-head strong {
  display: block;
  font-size: 18px;
}

.budget-form__section-head p {
  margin: 8px 0 0;
  color: var(--text-subtle);
  line-height: 1.6;
}

.budget-form__category-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
  max-height: 360px;
  overflow: auto;
  padding-right: 4px;
}

.budget-form__category-item {
  padding: 16px;
  border-radius: 18px;
  background: rgba(247, 247, 248, 0.92);
}

.budget-form__category-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.budget-form__category-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.budget-form__category-title strong {
  display: block;
}

.budget-form__category-title p {
  margin: 6px 0 0;
  color: var(--text-subtle);
  line-height: 1.5;
}

.budget-form__category-remaining {
  font-size: 13px;
  color: var(--text-subtle);
}

.budget-form__empty {
  padding: 18px;
  border-radius: 18px;
  background: rgba(247, 247, 248, 0.92);
  color: var(--text-subtle);
  line-height: 1.6;
}

.budget-form__footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.budget-form__button {
  min-width: 112px;
  min-height: 42px;
  padding: 0 18px;
  border-radius: 14px;
  border: 1px solid transparent;
  font-weight: 700;
  cursor: pointer;
}

.budget-form__button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.budget-form__button--ghost {
  border-color: var(--border-color);
  background: #fff;
  color: var(--text-main);
}

.budget-form__button--primary {
  background: var(--brand-color);
  color: #fff;
}

@media (max-width: 960px) {
  .budget-month-page__categories,
  .budget-card__numbers {
    grid-template-columns: minmax(0, 1fr);
  }

  .budget-card__header,
  .budget-form__category-meta {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
