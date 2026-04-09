<template>
  <div class="ledger-page">
    <section class="ledger-page__toolbar">
      <div class="month-switcher">
        <button :disabled="monthIndex === 0" @click="prevMonth">‹</button>
        <strong>{{ currentMonth.label }}</strong>
        <button :disabled="monthIndex === months.length - 1" @click="nextMonth">›</button>
      </div>

      <div class="ledger-page__actions">
        <button class="ledger-page__action ledger-page__action--expense" @click="openAction('支出')">
          记一笔支出
        </button>
        <button class="ledger-page__action ledger-page__action--income" @click="openAction('收入')">
          记一笔收入
        </button>
        <button class="ledger-page__action ledger-page__action--category" @click="openAction('分类管理')">
          分类管理
        </button>
      </div>
    </section>

    <section class="ledger-page__summary">
      <article class="summary-card">
        <span>本月收入</span>
        <strong class="is-income">{{ formatMoney(currentMonth.summary.income) }}</strong>
      </article>
      <article class="summary-card">
        <span>本月支出</span>
        <strong class="is-expense">{{ formatMoney(currentMonth.summary.expense) }}</strong>
      </article>
      <article class="summary-card">
        <span>结余</span>
        <strong>{{ formatMoney(currentMonth.summary.balance) }}</strong>
      </article>
    </section>

    <section class="ledger-page__content">
      <div class="ledger-feed page-card">
        <article v-if="!currentMonth.groups.length" class="ledger-empty">
          <strong>当前月份还没有收支记录</strong>
          <p>点击右上角“记一笔支出”或“记一笔收入”开始记录。</p>
        </article>

        <article v-for="group in currentMonth.groups" :key="group.date" class="ledger-group">
          <header class="ledger-group__header">
            <div>
              <h2>{{ group.date }}</h2>
              <span>{{ group.weekday }}</span>
            </div>
            <div class="ledger-group__totals">
              <span class="is-income">收 {{ formatMoney(group.totalIncome) }}</span>
              <span class="is-expense">支 {{ formatMoney(group.totalExpense) }}</span>
            </div>
          </header>

          <div class="ledger-group__list">
            <article v-for="item in group.items" :key="item.id" class="ledger-item">
              <div class="ledger-item__icon">{{ item.badge }}</div>
              <div class="ledger-item__content">
                <div class="ledger-item__headline">
                  <strong>{{ item.category }}</strong>
                  <span>{{ item.time }}</span>
                </div>
                <p>{{ item.note }}</p>
                <small v-if="item.imageName" class="ledger-item__media">图片：{{ item.imageName }}</small>
              </div>
              <strong :class="['ledger-item__amount', item.type === 'income' ? 'is-income' : 'is-expense']">
                {{ signedMoney(item.amount, item.type) }}
              </strong>
            </article>
          </div>
        </article>
      </div>

      <aside class="ledger-side">
        <section class="ledger-side__section page-card">
          <header class="ledger-side__header">
            <div>
              <h3>常用分类</h3>
              <p>按本月使用频次展示</p>
            </div>
            <button @click="openCategoryManager">管理</button>
          </header>

          <div class="category-grid">
            <article v-for="item in currentMonth.categories" :key="item.name" class="category-card">
              <span class="category-card__badge">{{ item.badge }}</span>
              <strong>{{ item.name }}</strong>
              <small>{{ item.count }} 笔</small>
            </article>
          </div>
        </section>

        <section class="ledger-side__section page-card">
          <header class="ledger-side__header">
            <div>
              <h3>本月概况</h3>
              <p>保留 Make 原型里的信息节奏</p>
            </div>
          </header>

          <article v-for="item in currentMonth.overview" :key="item.label" class="overview-card">
            <div class="overview-card__meta">
              <span>{{ item.label }}</span>
              <strong>{{ item.value }}</strong>
            </div>
            <div class="overview-card__track">
              <span :style="{ width: item.progress + '%' }"></span>
            </div>
          </article>
        </section>
      </aside>
    </section>
  </div>

  <el-dialog v-model="entryDialogVisible" :title="entryDialogTitle" width="560px" destroy-on-close>
    <form class="ledger-entry-form" @submit.prevent="submitEntry">
      <label class="ledger-entry-form__field">
        <span>分类</span>
        <el-select v-model="entryForm.categoryId" placeholder="请选择分类">
          <el-option v-for="item in entryCategoryOptions" :key="item.id" :label="item.name" :value="item.id">
            <div class="ledger-option">
              <strong>{{ item.badge }}</strong>
              <span>{{ item.name }}</span>
            </div>
          </el-option>
        </el-select>
      </label>

      <label class="ledger-entry-form__field">
        <span>金额</span>
        <input
          v-model.number="entryForm.amount"
          class="ledger-entry-form__number"
          type="number"
          min="0.01"
          step="0.01"
          placeholder="请输入金额"
        />
      </label>

      <label class="ledger-entry-form__field">
        <span>备注</span>
        <el-input v-model="entryForm.note" type="textarea" :rows="3" maxlength="60" show-word-limit />
      </label>

      <label class="ledger-entry-form__field">
        <span>日期</span>
        <input v-model="entryForm.date" class="ledger-entry-form__date" type="date" />
      </label>

      <label class="ledger-entry-form__field">
        <span>图片（可选，最多 1 张）</span>
        <input class="ledger-entry-form__file" type="file" accept="image/*" @change="onImageChange" />
      </label>

      <p v-if="entryForm.imageName" class="ledger-entry-form__image-name">已选择：{{ entryForm.imageName }}</p>

      <footer class="ledger-entry-form__footer">
        <button class="ledger-entry-form__button ledger-entry-form__button--ghost" type="button" @click="entryDialogVisible = false">
          取消
        </button>
        <button class="ledger-entry-form__button ledger-entry-form__button--primary" type="submit">保存</button>
      </footer>
    </form>
  </el-dialog>

  <el-dialog v-model="categoryDialogVisible" title="分类管理" width="640px" destroy-on-close>
    <div class="ledger-category-manager">
      <el-tabs v-model="categoryTab">
        <el-tab-pane label="支出分类" name="expense"></el-tab-pane>
        <el-tab-pane label="收入分类" name="income"></el-tab-pane>
      </el-tabs>

      <div class="ledger-category-manager__create">
        <label class="ledger-category-manager__field">
          <span>分类名称</span>
          <el-input v-model="categoryDraft[categoryTab].name" maxlength="12" placeholder="例如：旅行" />
        </label>
        <label class="ledger-category-manager__field ledger-category-manager__field--short">
          <span>徽章</span>
          <el-input v-model="categoryDraft[categoryTab].badge" maxlength="2" placeholder="旅" />
        </label>
        <button class="ledger-category-manager__add" type="button" @click="addCategory">新增分类</button>
      </div>

      <div class="ledger-category-manager__list">
        <article v-for="item in currentCategoryList" :key="item.id" class="ledger-category-card">
          <div class="ledger-category-card__meta">
            <span class="ledger-category-card__badge">{{ item.badge }}</span>
            <div>
              <strong>{{ item.name }}</strong>
              <small>{{ item.isDefault ? "默认分类" : "自定义分类" }}</small>
            </div>
          </div>
          <button
            class="ledger-category-card__remove"
            type="button"
            :disabled="item.isDefault"
            @click="removeCategory(item)"
          >
            删除
          </button>
        </article>
      </div>
    </div>
  </el-dialog>
</template>

<script>
import {
  createUserCategory,
  createUserLedger,
  deleteUserCategory,
  getUserLedger,
  listUserCategories
} from "@/api/userLedger";

function pad(value) {
  return String(value).padStart(2, "0");
}

function getTodayDate() {
  var now = new Date();
  return now.getFullYear() + "-" + pad(now.getMonth() + 1) + "-" + pad(now.getDate());
}

function getCurrentMonthKey() {
  return getTodayDate().slice(0, 7);
}

function monthLabel(monthKey) {
  var parts = String(monthKey || "").split("-");
  if (parts.length !== 2) {
    return monthKey || "暂无数据";
  }

  return parts[0] + "年" + Number(parts[1]) + "月";
}

function formatLedgerCurrency(value) {
  return "¥" + Number(value || 0).toLocaleString("zh-CN");
}

function buildMonthPlaceholder(monthKey) {
  return {
    key: monthKey,
    label: monthLabel(monthKey),
    summary: {
      income: 0,
      expense: 0,
      balance: 0
    },
    groups: [],
    categories: [],
    overview: []
  };
}

function shiftMonthKey(monthKey, offset) {
  var parts = String(monthKey || "").split("-");
  var year = Number(parts[0]);
  var month = Number(parts[1]);
  var date = new Date(year, month - 1 + offset, 1);

  return date.getFullYear() + "-" + pad(date.getMonth() + 1);
}

function buildInitialMonths() {
  var list = [];
  var currentMonthKey = getCurrentMonthKey();
  var index = 0;

  for (index = 0; index < 12; index += 1) {
    list.push(buildMonthPlaceholder(shiftMonthKey(currentMonthKey, -index)));
  }

  return list;
}

function normalizeAmount(value) {
  var amount = Number(value || 0);
  return Number.isFinite(amount) ? amount : 0;
}

function extractPayload(result) {
  if (result && result.data) {
    return result.data;
  }

  return result || {};
}

function toBadge(name, fallback) {
  var text = String(name || fallback || "").trim();
  return text ? text.slice(0, 1) : "-";
}

function normalizeLedgerMonth(result, fallbackMonthKey) {
  var payload = extractPayload(result);
  var summary = payload.summary || {};
  var groups = Array.isArray(payload.groups) ? payload.groups : [];
  var categories = Array.isArray(payload.categories) ? payload.categories : [];
  var overview = Array.isArray(payload.overview) ? payload.overview : [];
  var key = payload.key || fallbackMonthKey;

  return {
    key: key,
    label: payload.label || monthLabel(key),
    summary: {
      income: normalizeAmount(summary.income),
      expense: normalizeAmount(summary.expense),
      balance: normalizeAmount(summary.balance)
    },
    groups: groups.map(function(group) {
      var items = Array.isArray(group.items) ? group.items : [];

      return {
        date: group.date || "",
        weekday: group.weekday || "",
        totalIncome: normalizeAmount(group.totalIncome),
        totalExpense: normalizeAmount(group.totalExpense),
        items: items.map(function(item) {
          return {
            id: item.id,
            badge: item.badge || toBadge(item.category),
            category: item.category || "未分类",
            time: item.time || "00:00",
            note: item.note || "",
            amount: normalizeAmount(item.amount),
            type: item.type === "income" ? "income" : "expense",
            imageName: item.imageName || ""
          };
        })
      };
    }),
    categories: categories.map(function(item) {
      return {
        name: item.name || "未分类",
        badge: item.badge || toBadge(item.name),
        count: Number(item.count || 0)
      };
    }),
    overview: overview.map(function(item) {
      return {
        label: item.label || "-",
        value: item.value || "-",
        progress: Number(item.progress || 0)
      };
    })
  };
}

function normalizeCategoryList(result, type) {
  var payload = extractPayload(result);
  var list = Array.isArray(payload.list) ? payload.list : [];

  return list.map(function(item) {
    return {
      id: item.id,
      type: item.type || type,
      name: item.name || "",
      badge: item.badge || toBadge(item.name),
      isDefault: Boolean(item.isDefault)
    };
  });
}

function buildErrorMessage(error, fallback) {
  if (error && error.response && error.response.data && error.response.data.message) {
    return error.response.data.message;
  }

  return fallback;
}

export default {
  name: "UserHome",
  data() {
    return {
      monthIndex: 0,
      months: buildInitialMonths(),
      entryDialogVisible: false,
      categoryDialogVisible: false,
      categoryTab: "expense",
      categoryGroups: {
        expense: [],
        income: []
      },
      ledgerLoading: false,
      entryForm: {
        type: "expense",
        categoryId: "",
        amount: null,
        note: "",
        date: getTodayDate(),
        imageName: ""
      },
      categoryDraft: {
        expense: {
          name: "",
          badge: ""
        },
        income: {
          name: "",
          badge: ""
        }
      }
    };
  },
  computed: {
    currentMonth() {
      if (this.months.length) {
        return this.months[this.monthIndex];
      }

      return {
        key: "",
        label: "暂无数据",
        summary: { income: 0, expense: 0, balance: 0 },
        groups: [],
        categories: [],
        overview: []
      };
    },
    entryDialogTitle() {
      return this.entryForm.type === "income" ? "新增收入" : "新增支出";
    },
    entryCategoryOptions() {
      return this.categoryGroups[this.entryForm.type] || [];
    },
    currentCategoryList() {
      return this.categoryGroups[this.categoryTab] || [];
    }
  },
  created() {
    this.refreshCategories();
    this.refreshLedger(this.currentMonth.key);
  },
  methods: {
    formatMoney: formatLedgerCurrency,
    signedMoney(value, type) {
      return (type === "income" ? "+" : "-") + this.formatMoney(value);
    },
    prevMonth() {
      if (this.monthIndex > 0) {
        this.monthIndex -= 1;
        this.loadCurrentMonth();
      }
    },
    nextMonth() {
      if (this.monthIndex < this.months.length - 1) {
        this.monthIndex += 1;
        this.loadCurrentMonth();
      }
    },
    ensureMonthSlot(monthKey) {
      var exists = this.months.some(function(item) {
        return item.key === monthKey;
      });

      if (!exists) {
        this.months = this.months.concat([buildMonthPlaceholder(monthKey)]).sort(function(a, b) {
          return b.key.localeCompare(a.key);
        });
      }
    },
    selectMonth(monthKey) {
      var nextIndex = this.months.findIndex(function(item) {
        return item.key === monthKey;
      });

      this.monthIndex = nextIndex >= 0 ? nextIndex : 0;
    },
    replaceMonth(monthData) {
      this.ensureMonthSlot(monthData.key);
      this.months = this.months
        .map(function(item) {
          return item.key === monthData.key ? monthData : item;
        })
        .sort(function(a, b) {
          return b.key.localeCompare(a.key);
        });
      this.selectMonth(monthData.key);
    },
    loadCurrentMonth() {
      if (!this.currentMonth || !this.currentMonth.key) {
        return;
      }

      this.refreshLedger(this.currentMonth.key);
    },
    refreshLedger(targetMonthKey) {
      var monthKey = targetMonthKey || (this.currentMonth ? this.currentMonth.key : getCurrentMonthKey());

      this.ensureMonthSlot(monthKey);
      this.ledgerLoading = true;

      return getUserLedger(monthKey)
        .then(
          function(result) {
            this.replaceMonth(normalizeLedgerMonth(result, monthKey));
          }.bind(this)
        )
        .catch(
          function(error) {
            this.$message.error(buildErrorMessage(error, "账本数据加载失败，请稍后重试。"));
          }.bind(this)
        )
        .finally(
          function() {
            this.ledgerLoading = false;
          }.bind(this)
        );
    },
    refreshCategories() {
      return Promise.all([listUserCategories("expense"), listUserCategories("income")])
        .then(
          function(results) {
            this.categoryGroups = {
              expense: normalizeCategoryList(results[0], "expense"),
              income: normalizeCategoryList(results[1], "income")
            };
            this.ensureCategorySelection();
          }.bind(this)
        )
        .catch(
          function(error) {
            this.$message.error(buildErrorMessage(error, "分类加载失败，请稍后重试。"));
          }.bind(this)
        );
    },
    ensureCategorySelection() {
      var options = this.entryCategoryOptions;
      var matched = options.some(
        function(item) {
          return String(item.id) === String(this.entryForm.categoryId);
        }.bind(this)
      );

      if (!matched) {
        this.entryForm.categoryId = options.length ? options[0].id : "";
      }
    },
    openAction(label) {
      if (label === "支出") {
        this.openEntryDialog("expense");
        return;
      }

      if (label === "收入") {
        this.openEntryDialog("income");
        return;
      }

      this.openCategoryManager();
    },
    openEntryDialog(type) {
      this.entryForm.type = type;
      this.entryForm.amount = null;
      this.entryForm.note = "";
      this.entryForm.date = getTodayDate();
      this.entryForm.imageName = "";
      this.ensureCategorySelection();
      this.entryDialogVisible = true;
    },
    openCategoryManager() {
      this.categoryDialogVisible = true;
    },
    onImageChange(event) {
      var file = event.target && event.target.files ? event.target.files[0] : null;
      this.entryForm.imageName = file ? file.name : "";
    },
    submitEntry() {
      if (!this.entryForm.categoryId) {
        this.$message.warning("请选择分类");
        return;
      }

      if (!(Number(this.entryForm.amount) > 0)) {
        this.$message.warning("金额必须大于 0");
        return;
      }

      createUserLedger({
        type: this.entryForm.type,
        category_id: Number(this.entryForm.categoryId),
        amount: Number(this.entryForm.amount),
        remark: this.entryForm.note,
        record_date: this.entryForm.date,
        image_url: this.entryForm.imageName
      })
        .then(
          function(result) {
            var payload = extractPayload(result);
            var monthKey = payload.month || String(this.entryForm.date || "").slice(0, 7) || getCurrentMonthKey();

            this.entryDialogVisible = false;
            this.refreshLedger(monthKey);
            this.$message.success(this.entryDialogTitle + "已保存并回显到列表。");
          }.bind(this)
        )
        .catch(
          function(error) {
            this.$message.error(buildErrorMessage(error, "保存失败，请稍后重试。"));
          }.bind(this)
        );
    },
    addCategory() {
      var draft = this.categoryDraft[this.categoryTab];
      if (!String(draft.name || "").trim()) {
        this.$message.warning("分类名称不能为空");
        return;
      }

      createUserCategory({
        type: this.categoryTab,
        name: draft.name,
        badge: draft.badge
      })
        .then(
          function() {
            draft.name = "";
            draft.badge = "";
            return Promise.all([this.refreshCategories(), this.refreshLedger(this.currentMonth.key)]);
          }.bind(this)
        )
        .then(
          function() {
            this.$message.success("分类已新增。");
          }.bind(this)
        )
        .catch(
          function(error) {
            this.$message.error(buildErrorMessage(error, "分类新增失败。"));
          }.bind(this)
        );
    },
    removeCategory(item) {
      deleteUserCategory(item.id)
        .then(
          function() {
            return Promise.all([this.refreshCategories(), this.refreshLedger(this.currentMonth.key)]);
          }.bind(this)
        )
        .then(
          function() {
            this.$message.success("分类已删除。");
          }.bind(this)
        )
        .catch(
          function(error) {
            this.$message.warning(buildErrorMessage(error, "分类删除失败。"));
          }.bind(this)
        );
    }
  }
};
</script>

<style scoped>
.ledger-page {
  display: flex;
  flex-direction: column;
  gap: 22px;
}

.ledger-page__toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.month-switcher {
  display: inline-flex;
  align-items: center;
  gap: 14px;
}

.month-switcher button {
  width: 40px;
  height: 40px;
  border-radius: 14px;
  border: 1px solid var(--border-color);
  background: rgba(255, 255, 255, 0.9);
  color: var(--text-main);
  font-size: 22px;
}

.month-switcher button:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.month-switcher strong {
  font-size: 22px;
}

.ledger-page__actions {
  display: flex;
  gap: 12px;
}

.ledger-page__action {
  min-height: 44px;
  padding: 0 18px;
  border: none;
  border-radius: 999px;
  color: #ffffff;
  font-weight: 700;
  box-shadow: var(--shadow-sm);
}

.ledger-page__action--expense {
  background: #ef4444;
}

.ledger-page__action--income {
  background: #22c55e;
}

.ledger-page__action--category {
  background: #f59e0b;
}

.ledger-page__summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.summary-card {
  padding: 22px 24px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.88);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
}

.summary-card span {
  display: block;
  margin-bottom: 12px;
  color: var(--text-muted);
  font-size: 14px;
}

.summary-card strong {
  font-size: 24px;
}

.summary-card strong.is-income {
  color: var(--success-color);
}

.summary-card strong.is-expense {
  color: var(--danger-color);
}

.ledger-page__content {
  display: grid;
  grid-template-columns: minmax(0, 1.55fr) 320px;
  gap: 18px;
}

.ledger-feed {
  padding: 24px;
}

.ledger-empty {
  padding: 28px;
  border-radius: 20px;
  background: var(--card-muted);
  border: 1px dashed var(--border-color);
}

.ledger-empty strong {
  display: block;
  font-size: 18px;
}

.ledger-empty p {
  margin: 10px 0 0;
  color: var(--text-muted);
}

.ledger-group + .ledger-group {
  margin-top: 26px;
}

.ledger-group__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
}

.ledger-group__header h2 {
  margin: 0;
  font-size: 28px;
}

.ledger-group__header span {
  margin-left: 10px;
  color: var(--text-muted);
  font-size: 14px;
}

.ledger-group__totals {
  display: inline-flex;
  gap: 12px;
  font-weight: 700;
}

.ledger-group__totals .is-income {
  color: var(--success-color);
}

.ledger-group__totals .is-expense {
  color: var(--danger-color);
}

.ledger-group__list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.ledger-item {
  display: grid;
  grid-template-columns: 54px minmax(0, 1fr) auto;
  gap: 16px;
  align-items: center;
  padding: 18px 16px;
  border-radius: 22px;
  background: var(--card-muted);
  border: 1px solid rgba(229, 231, 235, 0.92);
}

.ledger-item__icon {
  width: 54px;
  height: 54px;
  border-radius: 18px;
  background: #ffffff;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 700;
  color: var(--text-subtle);
}

.ledger-item__headline {
  display: flex;
  align-items: center;
  gap: 12px;
}

.ledger-item__headline strong {
  font-size: 16px;
}

.ledger-item__headline span,
.ledger-item__content p {
  color: var(--text-muted);
}

.ledger-item__content p {
  margin: 6px 0 0;
  font-size: 14px;
}

.ledger-item__media {
  display: inline-block;
  margin-top: 8px;
  color: #a16207;
  font-size: 12px;
}

.ledger-item__amount {
  font-size: 18px;
}

.ledger-item__amount.is-income {
  color: var(--success-color);
}

.ledger-item__amount.is-expense {
  color: var(--danger-color);
}

.ledger-side {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.ledger-side__section {
  padding: 22px;
}

.ledger-side__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 14px;
  margin-bottom: 18px;
}

.ledger-side__header h3 {
  margin: 0 0 4px;
  font-size: 22px;
}

.ledger-side__header p {
  margin: 0;
  color: var(--text-muted);
  font-size: 13px;
}

.ledger-side__header button {
  border: none;
  background: transparent;
  color: #d19b00;
  font-weight: 700;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.category-card {
  padding: 16px;
  border-radius: 20px;
  background: var(--card-muted);
  border: 1px solid rgba(229, 231, 235, 0.92);
}

.category-card__badge {
  width: 36px;
  height: 36px;
  border-radius: 14px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: #ffffff;
  color: var(--text-subtle);
  font-weight: 700;
}

.category-card strong {
  display: block;
  margin-top: 12px;
}

.category-card small {
  display: block;
  margin-top: 6px;
  color: var(--text-muted);
}

.overview-card + .overview-card {
  margin-top: 14px;
}

.overview-card__meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 10px;
}

.overview-card__meta span {
  color: var(--text-muted);
  font-size: 14px;
}

.overview-card__meta strong {
  font-size: 18px;
}

.overview-card__track {
  height: 8px;
  border-radius: 999px;
  background: rgba(229, 231, 235, 0.9);
  overflow: hidden;
}

.overview-card__track span {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(135deg, var(--brand-color) 0%, var(--brand-hover) 100%);
}

.ledger-entry-form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.ledger-entry-form__field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.ledger-entry-form__field > span {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-subtle);
}

.ledger-entry-form__date,
.ledger-entry-form__number,
.ledger-entry-form__file {
  min-height: 40px;
  border: 1px solid var(--border-color);
  border-radius: 10px;
  padding: 0 12px;
  background: #fff;
}

.ledger-entry-form__image-name {
  margin: -4px 0 0;
  color: #a16207;
  font-size: 12px;
}

.ledger-entry-form__footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 8px;
}

.ledger-entry-form__button {
  min-height: 40px;
  border-radius: 10px;
  padding: 0 16px;
  border: 1px solid var(--border-color);
  font-weight: 600;
}

.ledger-entry-form__button--ghost {
  background: #fff;
}

.ledger-entry-form__button--primary {
  border-color: transparent;
  background: var(--brand-color);
}

.ledger-option {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.ledger-option strong {
  width: 24px;
  height: 24px;
  border-radius: 8px;
  background: rgba(246, 211, 74, 0.35);
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.ledger-category-manager__create {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 130px auto;
  gap: 12px;
  align-items: end;
  margin: 8px 0 16px;
}

.ledger-category-manager__field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.ledger-category-manager__field span {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-subtle);
}

.ledger-category-manager__add {
  min-height: 40px;
  border-radius: 10px;
  border: none;
  background: var(--brand-color);
  font-weight: 700;
}

.ledger-category-manager__list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.ledger-category-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 16px;
  border-radius: 14px;
  border: 1px solid var(--border-color);
  background: rgba(255, 255, 255, 0.92);
}

.ledger-category-card__meta {
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.ledger-category-card__badge {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  background: rgba(246, 211, 74, 0.4);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
}

.ledger-category-card__meta strong {
  display: block;
}

.ledger-category-card__meta small {
  color: var(--text-muted);
}

.ledger-category-card__remove {
  min-height: 32px;
  padding: 0 12px;
  border-radius: 8px;
  border: 1px solid rgba(239, 68, 68, 0.4);
  background: rgba(239, 68, 68, 0.08);
  color: #dc2626;
}

.ledger-category-card__remove:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

@media (max-width: 1200px) {
  .ledger-page__content {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 920px) {
  .ledger-page__toolbar,
  .ledger-page__actions {
    flex-direction: column;
    align-items: stretch;
  }

  .ledger-page__summary {
    grid-template-columns: 1fr;
  }

  .ledger-group__header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .ledger-category-manager__create {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 820px) {
  .category-grid {
    grid-template-columns: 1fr;
  }
}
</style>
