<template>
  <div class="finance-page budget-year-page">
    <div v-if="errorMessage" class="finance-inline-notice">
      <span>{{ errorMessage }}</span>
      <button type="button" @click="retryCurrentYear">重新加载</button>
    </div>

    <section class="finance-toolbar">
      <div class="finance-switcher">
        <button :disabled="!hasPreviousYear" @click="prevYear">‹</button>
        <div>
          <strong>{{ selectedYear }} 年预算</strong>
          <p class="budget-year-page__switcher-note">预算是独立视图，不复制月预算的信息结构。</p>
        </div>
        <button :disabled="!hasNextYear" @click="nextYear">›</button>
      </div>

      <div class="finance-toolbar__actions">
        <button class="finance-button finance-button--ghost" type="button" @click="openMonthBudget">查看月预算</button>
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
        <span class="finance-hero__eyebrow-mark">年</span>
        <span>{{ selectedYear }} 年预算总览</span>
      </div>

      <div class="finance-hero__headline">
        <h2>年度预算最重要的是看节奏，而不是盯着某一个分类的瞬时高点</h2>
        <p>{{ yearData.overview.note }}</p>
      </div>

      <div class="budget-year-page__hero-progress">
        <div class="budget-year-page__hero-progress-meta">
          <span>预算执行进度</span>
          <strong :class="yearData.overview.remaining < 0 ? 'finance-tone-expense' : 'finance-tone-income'">
            {{ formatPercent(yearData.overview.percentage) }}
          </strong>
        </div>
        <div class="finance-progress">
          <span
            :class="yearData.overview.remaining < 0 ? 'budget-year-page__hero-bar--danger' : 'budget-year-page__hero-bar--safe'"
            :style="{ width: progressWidth(yearData.overview.percentage) }"
          ></span>
        </div>
      </div>

      <div class="finance-stat-grid finance-stat-grid--triple">
        <article class="finance-stat-card">
          <span>本年预算</span>
          <strong>{{ formatCurrency(yearData.overview.totalBudget) }}</strong>
          <small>已执行 {{ yearData.overview.months }} 个月</small>
        </article>
        <article class="finance-stat-card">
          <span>本年支出</span>
          <strong class="finance-tone-expense">{{ formatCurrency(yearData.overview.totalExpense) }}</strong>
          <small>预算节奏需要持续盯盘</small>
        </article>
        <article class="finance-stat-card">
          <span>剩余预算</span>
          <strong :class="yearData.overview.remaining < 0 ? 'finance-tone-expense' : 'finance-tone-income'">
            {{ formatCurrency(Math.abs(yearData.overview.remaining)) }}
          </strong>
          <small>{{ yearData.overview.remaining >= 0 ? "仍有可调空间" : "当前已经透支" }}</small>
        </article>
      </div>
    </section>

    <article class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>月度预算执行情况</h3>
          <p>每个月的预算与实际支出并列展示，更容易发现全年失控点。</p>
        </div>
      </header>
      <div v-if="isLoading && !executionLabels.length" class="budget-year-page__empty">
        <strong>正在加载年度预算执行图…</strong>
        <p>月度预算与支出同步完成后会显示在这里。</p>
      </div>
      <div v-else-if="!executionLabels.length" class="budget-year-page__empty">
        <strong>当前年份暂无月度预算执行数据</strong>
        <p>设置当前年度预算后，这里会显示每个月的预算与实际支出。</p>
      </div>
      <SimpleGroupedBarChart
        v-else
        :labels="executionLabels"
        :series="executionSeries"
        :height="320"
      />
    </article>

    <section class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>分类年度预算汇总</h3>
          <p>{{ categoryPanelDescription }}</p>
        </div>
      </header>

      <div v-if="isLoading && !yearData.categories.length" class="budget-year-page__empty">
        <strong>正在加载年度分类预算…</strong>
        <p>预算与分类支出同步后会显示在这里。</p>
      </div>
      <div v-else-if="!yearData.categories.length" class="budget-year-page__empty">
        <strong>当前还没有可展示的年度分类预算</strong>
        <p>设置当年预算后，这里会按支出分类显示预算执行状态。</p>
      </div>

      <div v-else class="budget-year-page__categories">
        <article v-for="item in yearData.categories" :key="item.key" class="year-budget-card">
          <div class="year-budget-card__header">
            <div class="year-budget-card__title">
              <span class="finance-badge">{{ item.badge }}</span>
              <div>
                <strong>{{ item.name }}</strong>
                <p>{{ formatCurrency(item.expense) }} / {{ formatCurrency(item.budget) }}</p>
              </div>
            </div>

            <button
              class="finance-button finance-button--ghost year-budget-card__action"
              type="button"
              :disabled="!canEditBudget || isLoading"
              @click="openBudgetDialog"
            >
              调整
            </button>
          </div>

          <div class="year-budget-card__meta">
            <div>
              <span>执行比例</span>
              <strong :class="statusTone(item.status)">{{ formatPercent(item.percentage) }}</strong>
            </div>
            <div>
              <span>{{ item.remaining >= 0 ? "剩余" : "超支" }}</span>
              <strong :class="statusTone(item.status)">
                {{ formatCurrency(Math.abs(item.remaining)) }}
              </strong>
            </div>
          </div>

          <div class="finance-progress">
            <span :class="statusBar(item.status)" :style="{ width: progressWidth(item.percentage) }"></span>
          </div>

          <p class="year-budget-card__status" :class="statusTone(item.status)">
            {{ statusLabel(item.status) }}
          </p>
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
        <span>本年总预算</span>
        <input
          v-model.number="budgetForm.totalBudget"
          class="budget-form__input"
          type="number"
          min="0"
          step="0.01"
          placeholder="请输入本年总预算"
        />
      </label>

      <div class="budget-form__section">
        <div class="budget-form__section-head">
          <div>
            <strong>分类预算</strong>
            <p>按支出分类分别设置年度额度，保存后会即时回刷当前页面。</p>
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
              :placeholder="'设置 ' + item.name + ' 年度预算'"
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
import SimpleGroupedBarChart from "@/components/SimpleGroupedBarChart.vue";
import {
  getYearBudget,
  getYearBudgetOptions,
  updateCurrentYearBudget
} from "@/api/userBudget";

function pad(value) {
  return String(value).padStart(2, "0");
}

function getCurrentYear() {
  return new Date().getFullYear();
}

function toNumber(value) {
  var parsed = Number(value);

  if (!Number.isFinite(parsed)) {
    return 0;
  }

  return parsed;
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

function extractPayload(result) {
  if (result && typeof result === "object" && result.data && typeof result.data === "object") {
    return result.data;
  }

  return result || {};
}

function buildBudgetError(error, fallback) {
  if (error && error.response && error.response.data && error.response.data.message) {
    return error.response.data.message;
  }

  return fallback;
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
    return "当前分类年度预算已经超支，建议全年计划及时收口。";
  }

  if (item.percentage >= 80) {
    return "当前分类年度预算接近上限，需要持续关注。";
  }

  if (item.expense <= 0) {
    return "当前年度暂未发生该分类支出。";
  }

  return "当前分类年度预算执行正常，可继续按年度计划安排。";
}

function mapYearCategory(item, index) {
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

function mapExecutionItem(item, index) {
  var source = item || {};
  var month = toNumber(source.month || source.order || index + 1);
  var label = source.label || month + "月";

  return {
    label: label,
    budget: toNumber(source.budget || source.totalBudget),
    expense: toNumber(source.expense || source.totalExpense || source.usedAmount)
  };
}

function buildYearNote(year, overview) {
  if (overview.note) {
    return overview.note;
  }

  if (overview.totalBudget <= 0) {
    return year + " 年预算尚未设置，可以先录入年度预算和分类额度。";
  }

  if (overview.remaining < 0) {
    return year + " 年预算已经透支，需要重新分配剩余月份的预算节奏。";
  }

  return year + " 年预算执行正常，可结合月度预算执行情况继续调整。";
}

function createEmptyYearBudget(year) {
  var resolvedYear = Number(year || getCurrentYear());
  var overview = {
    totalBudget: 0,
    totalExpense: 0,
    remaining: 0,
    percentage: 0,
    months: 0,
    note: buildYearNote(resolvedYear, {})
  };

  return {
    year: resolvedYear,
    overview: overview,
    monthlyExecution: [],
    categories: [],
    canEdit: resolvedYear === getCurrentYear()
  };
}

function normalizeYearOptionsPayload(result) {
  var payload = extractPayload(result);
  var years = Array.isArray(payload) ? payload : payload.years;
  var normalizedYears = Array.from(new Set((Array.isArray(years) ? years : []).map(function(item) {
    return Number(item);
  }).filter(function(item) {
    return Number.isFinite(item) && item > 0;
  }))).sort(function(left, right) {
    return right - left;
  });

  if (!normalizedYears.length) {
    normalizedYears = [getCurrentYear()];
  }

  return normalizedYears;
}

function normalizeYearBudgetPayload(result, fallbackYear) {
  var payload = extractPayload(result);
  var year = Number(payload.year || fallbackYear || getCurrentYear());
  var overviewSource = payload.overview || {};
  var overview = {
    totalBudget: toNumber(overviewSource.totalBudget || overviewSource.budget || payload.totalBudget),
    totalExpense: toNumber(overviewSource.totalExpense || overviewSource.expense || payload.totalExpense),
    remaining: 0,
    percentage: 0,
    months: toNumber(overviewSource.months || overviewSource.executedMonths),
    note: ""
  };

  overview.remaining = Object.prototype.hasOwnProperty.call(overviewSource, "remaining")
    ? toNumber(overviewSource.remaining)
    : overview.totalBudget - overview.totalExpense;
  overview.percentage = Object.prototype.hasOwnProperty.call(overviewSource, "percentage")
    ? toNumber(overviewSource.percentage)
    : overview.totalBudget > 0
      ? Number(((overview.totalExpense / overview.totalBudget) * 100).toFixed(1))
      : 0;
  overview.note = buildYearNote(year, overviewSource);

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
    year: year,
    overview: overview,
    monthlyExecution: Array.isArray(payload.monthlyExecution) ? payload.monthlyExecution.map(mapExecutionItem) : [],
    categories: Array.isArray(payload.categories) ? payload.categories.map(mapYearCategory) : [],
    canEdit: typeof rawEditable === "boolean" ? rawEditable : year === getCurrentYear()
  };
}

export default {
  name: "UserBudgetYear",
  components: {
    SimpleGroupedBarChart
  },
  data() {
    return {
      years: [],
      selectedYear: getCurrentYear(),
      yearData: createEmptyYearBudget(),
      isLoading: false,
      isSaving: false,
      errorMessage: "",
      budgetDialogVisible: false,
      budgetForm: {
        totalBudget: 0,
        categories: []
      },
      activeYearRequestKey: 0
    };
  },
  computed: {
    canEditBudget() {
      return Boolean(this.yearData.canEdit);
    },
    budgetActionLabel() {
      return this.yearData.overview.totalBudget > 0 ? "调整年度预算" : "设置年度预算";
    },
    hasPreviousYear() {
      return this.years.length > 1 && this.years.indexOf(this.selectedYear) < this.years.length - 1;
    },
    hasNextYear() {
      return this.years.length > 1 && this.years.indexOf(this.selectedYear) > 0;
    },
    executionLabels() {
      return this.yearData.monthlyExecution.map(function(item) {
        return item.label;
      });
    },
    executionSeries() {
      return [
        {
          name: "预算",
          color: "#f6d34a",
          values: this.yearData.monthlyExecution.map(function(item) {
            return item.budget;
          })
        },
        {
          name: "实际支出",
          color: "#ef4444",
          values: this.yearData.monthlyExecution.map(function(item) {
            return item.expense;
          })
        }
      ];
    },
    categoryPanelDescription() {
      if (!this.canEditBudget) {
        return "当前年份预算处于只读状态，仅支持查看预算执行结果。";
      }

      return "年预算看的是长期分配结果，所以每一类都应该有明确状态。";
    }
  },
  created() {
    this.bootstrap();
  },
  methods: {
    formatCurrency: formatBudgetCurrency,
    formatPercent: formatBudgetPercent,
    bootstrap() {
      this.isLoading = true;
      this.errorMessage = "";

      return getYearBudgetOptions()
        .then(
          function(result) {
            var normalizedYears = normalizeYearOptionsPayload(result);
            var currentYear = getCurrentYear();
            var fallbackYear = normalizedYears.indexOf(currentYear) >= 0 ? currentYear : normalizedYears[0];

            this.years = normalizedYears;
            this.selectedYear = fallbackYear;

            return this.loadYearBudget(fallbackYear, true);
          }.bind(this)
        )
        .catch(
          function(error) {
            this.years = [getCurrentYear()];
            this.selectedYear = getCurrentYear();
            this.yearData = createEmptyYearBudget(this.selectedYear);
            this.errorMessage = buildBudgetError(error, "年度预算年份加载失败，请稍后重试。");
            this.$message.error(this.errorMessage);
          }.bind(this)
        )
        .finally(
          function() {
            this.isLoading = false;
          }.bind(this)
        );
    },
    loadYearBudget(year, silent) {
      var targetYear = Number(year || this.selectedYear || getCurrentYear());
      var requestKey = Date.now();

      this.activeYearRequestKey = requestKey;
      this.isLoading = true;
      this.errorMessage = "";

      return getYearBudget(targetYear)
        .then(
          function(result) {
            if (this.activeYearRequestKey !== requestKey) {
              return;
            }

            this.yearData = normalizeYearBudgetPayload(result, targetYear);
          }.bind(this)
        )
        .catch(
          function(error) {
            if (this.activeYearRequestKey !== requestKey) {
              return;
            }

            this.yearData = createEmptyYearBudget(targetYear);
            this.errorMessage = buildBudgetError(error, targetYear + " 年预算加载失败，请稍后重试。");

            if (!silent) {
              this.$message.error(this.errorMessage);
            }
          }.bind(this)
        )
        .finally(
          function() {
            if (this.activeYearRequestKey === requestKey) {
              this.isLoading = false;
            }
          }.bind(this)
        );
    },
    prevYear() {
      var index = this.years.indexOf(this.selectedYear);

      if (index > -1 && index < this.years.length - 1) {
        this.selectedYear = this.years[index + 1];
        this.loadYearBudget(this.selectedYear);
      }
    },
    nextYear() {
      var index = this.years.indexOf(this.selectedYear);

      if (index > 0) {
        this.selectedYear = this.years[index - 1];
        this.loadYearBudget(this.selectedYear);
      }
    },
    retryCurrentYear() {
      if (!this.years.length) {
        this.bootstrap();
        return;
      }

      this.loadYearBudget(this.selectedYear);
    },
    openMonthBudget() {
      this.$router.push("/user/budget/month");
    },
    openBudgetDialog() {
      if (!this.canEditBudget) {
        this.$message.info("当前年份预算为只读状态，暂不支持编辑。");
        return;
      }

      this.budgetForm = {
        totalBudget: this.yearData.overview.totalBudget,
        categories: this.yearData.categories.map(function(item) {
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
        year: this.selectedYear,
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
        this.$message.warning("本年总预算不能小于 0");
        return;
      }

      if (invalidCategory) {
        this.$message.warning("分类预算不能小于 0");
        return;
      }

      this.isSaving = true;

      updateCurrentYearBudget(this.buildUpdatePayload())
        .then(
          function() {
            this.budgetDialogVisible = false;
            return this.loadYearBudget(this.selectedYear);
          }.bind(this)
        )
        .then(
          function() {
            this.$message.success("本年预算已更新。");
          }.bind(this)
        )
        .catch(
          function(error) {
            this.$message.error(buildBudgetError(error, "本年预算保存失败，请稍后重试。"));
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
        over: "年度预算已经超支",
        warning: "预算接近上限，需要持续关注",
        safe: "预算执行正常"
      };

      return labels[status] || "预算执行正常";
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
        over: "budget-year-page__bar--danger",
        warning: "budget-year-page__bar--warning",
        safe: "budget-year-page__bar--safe"
      };

      return bars[status] || "budget-year-page__bar--safe";
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

.budget-year-page__switcher-note {
  margin: 8px 0 0;
  color: var(--text-subtle);
  font-size: 13px;
}

.budget-year-page__hero-progress {
  margin-top: 24px;
  padding: 20px 22px;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.88);
}

.budget-year-page__hero-progress-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.budget-year-page__hero-bar--danger {
  background: var(--danger-color);
}

.budget-year-page__hero-bar--safe {
  background: linear-gradient(90deg, #22c55e 0%, #f6d34a 100%);
}

.budget-year-page__empty {
  padding: 24px;
  border-radius: 22px;
  background: rgba(247, 247, 248, 0.92);
}

.budget-year-page__empty strong {
  display: block;
  font-size: 17px;
}

.budget-year-page__empty p {
  margin: 10px 0 0;
  color: var(--text-subtle);
  line-height: 1.6;
}

.budget-year-page__categories {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.year-budget-card {
  padding: 22px;
  border-radius: 24px;
  border: 1px solid var(--border-color);
  background: rgba(255, 255, 255, 0.96);
}

.year-budget-card__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.year-budget-card__title {
  display: flex;
  gap: 14px;
}

.year-budget-card__title strong {
  display: block;
  font-size: 18px;
}

.year-budget-card__title p {
  margin: 8px 0 0;
  color: var(--text-subtle);
  line-height: 1.6;
}

.year-budget-card__action {
  min-height: 38px;
  padding: 0 14px;
}

.year-budget-card__meta {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
  margin: 20px 0 14px;
}

.year-budget-card__meta span {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
}

.year-budget-card__meta strong {
  display: block;
  margin-top: 8px;
  font-size: 18px;
}

.year-budget-card__status {
  margin: 12px 0 0;
  line-height: 1.6;
  font-size: 13px;
  font-weight: 700;
}

.budget-year-page__bar--danger {
  background: var(--danger-color);
}

.budget-year-page__bar--warning {
  background: var(--warning-color);
}

.budget-year-page__bar--safe {
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
  .budget-year-page__categories,
  .year-budget-card__meta {
    grid-template-columns: minmax(0, 1fr);
  }

  .year-budget-card__header,
  .budget-form__category-meta {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
