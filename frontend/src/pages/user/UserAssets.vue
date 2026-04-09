<template>
  <div class="finance-page asset-page">
    <section class="finance-toolbar">
      <div class="finance-toolbar__meta">
        <h1>资产管家</h1>
        <p>把账户放在同一张资产地图里，先看全局，再处理单账户细节。</p>
      </div>
      <div class="finance-toolbar__actions">
        <button class="finance-button finance-button--ghost" @click="reloadFromSeed">重置演示数据</button>
        <button class="finance-button finance-button--primary" @click="openAddDialog">新增账户</button>
      </div>
    </section>

    <section class="finance-hero finance-hero--soft">
      <div class="finance-hero__eyebrow">
        <span class="finance-hero__eyebrow-mark">资</span>
        <span>资产总览</span>
      </div>
      <div class="finance-hero__headline">
        <h2>{{ formatCurrency(summary.netAsset) }}</h2>
        <p>净资产 = 资产 - 负债。最近一月净变化 {{ monthChangeLabel }}</p>
      </div>

      <div class="finance-stat-grid finance-stat-grid--triple">
        <article class="finance-stat-card">
          <span>总资产</span>
          <strong class="finance-tone-income">{{ formatCurrency(summary.totalAsset) }}</strong>
          <small>{{ summary.accountCount }} 个账户</small>
        </article>
        <article class="finance-stat-card">
          <span>总负债</span>
          <strong class="finance-tone-expense">{{ formatCurrency(summary.totalLiability) }}</strong>
          <small>与资产分层展示</small>
        </article>
        <article class="finance-stat-card">
          <span>更新时间</span>
          <strong>{{ summary.updatedAt }}</strong>
          <small>本页为原型数据演示</small>
        </article>
      </div>
    </section>

    <section class="page-card finance-panel">
      <header class="finance-panel__header">
        <div>
          <h3>账户分类列表</h3>
          <p>按账户分类展示总额和账户余额，支持新增、编辑和进入详情。</p>
        </div>
        <div class="asset-page__filter">
          <span class="finance-inline-tag">筛选分类</span>
          <el-select v-model="activeCategoryId" size="large" style="min-width: 220px;">
            <el-option v-for="item in categorySelectOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </div>
      </header>

      <div class="asset-page__category-list">
        <article v-for="category in visibleCategories" :key="category.id" class="asset-page__category-card">
          <header class="asset-page__category-header">
            <div class="asset-page__category-title">
              <span class="finance-badge">{{ categoryMark(category) }}</span>
              <div>
                <strong>{{ category.name }}</strong>
                <p>{{ category.accounts.length }} 个账户</p>
              </div>
            </div>
            <div class="asset-page__category-total">
              <span>分类总额</span>
              <strong :class="category.id === 'category-liability' ? 'finance-tone-expense' : 'finance-tone-income'">
                {{ formatCurrency(category.total) }}
              </strong>
            </div>
          </header>

          <div class="asset-page__account-list">
            <article
              v-for="account in category.accounts"
              :key="account.id"
              class="asset-page__account-row"
            >
              <div class="asset-page__account-main">
                <div>
                  <strong>{{ account.name }}</strong>
                  <p>{{ account.remark }}</p>
                </div>
                <span class="finance-inline-tag">{{ account.typeLabel }}<template v-if="account.provider"> · {{ providerLabel(account.provider) }}</template></span>
              </div>

              <div class="asset-page__account-balance">
                <small>账户余额</small>
                <strong :class="account.direction === 'liability' ? 'finance-tone-expense' : 'finance-tone-income'">
                  {{ formatCurrency(account.balance) }}
                </strong>
                <span :class="account.monthlyChange >= 0 ? 'finance-tone-income' : 'finance-tone-expense'">
                  {{ account.monthlyChange >= 0 ? "↑" : "↓" }} {{ formatCurrency(Math.abs(account.monthlyChange)) }}
                </span>
              </div>

              <div class="asset-page__actions">
                <button class="finance-button finance-button--ghost" @click="openEditDialog(category, account)">编辑</button>
                <button class="finance-button finance-button--primary" @click="openDetail(account.id)">详情</button>
              </div>
            </article>
          </div>
        </article>
      </div>
    </section>

    <el-dialog
      v-model="formDialogVisible"
      :title="formMode === 'create' ? '新增账户' : '编辑账户'"
      width="620px"
      destroy-on-close
    >
      <el-form ref="accountFormRef" :model="formModel" :rules="formRules" label-width="100px" status-icon>
        <el-form-item label="账户类型" prop="type">
          <el-select v-model="formModel.type" style="width: 100%;" :disabled="formMode === 'edit'" @change="onTypeChange">
            <el-option v-for="item in typeOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="账户名称" prop="name">
          <el-input v-model="formModel.name" maxlength="24" show-word-limit />
        </el-form-item>

        <el-form-item label="账户备注" prop="remark">
          <el-input v-model="formModel.remark" maxlength="40" show-word-limit />
        </el-form-item>

        <el-form-item label="账户余额" prop="balance">
          <el-input-number v-model="formModel.balance" :precision="2" :step="100" style="width: 100%;" />
        </el-form-item>

        <el-form-item label="所属分类" prop="categoryId">
          <el-select v-model="formModel.categoryId" style="width: 100%;">
            <el-option
              v-for="item in categoryOptionsForForm"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item v-if="requiresCardNo" label="卡号" prop="cardNo">
          <el-input v-model="formModel.cardNo" placeholder="例如：6225 **** **** 1902" />
        </el-form-item>

        <el-form-item v-if="requiresVirtualProvider" label="虚拟渠道" prop="provider">
          <el-select v-model="formModel.provider" style="width: 100%;">
            <el-option
              v-for="item in providerOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="asset-page__dialog-footer">
          <span v-if="formMode === 'edit'" class="finance-note">账户类型创建后不可修改。</span>
          <div class="asset-page__dialog-actions">
            <button class="finance-button finance-button--ghost" @click="closeFormDialog">取消</button>
            <button class="finance-button finance-button--primary" @click="submitForm">
              {{ formMode === "create" ? "创建账户" : "保存修改" }}
            </button>
          </div>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import {
  formatCurrency,
  getAssetOverview,
  getAssetTypeOptions,
  getVirtualProviderOptions
} from "@/utils/userAssetMock";

function toNumber(value) {
  return Number(value || 0);
}

export default {
  name: "UserAssets",
  data() {
    return {
      summary: {
        netAsset: 0,
        totalAsset: 0,
        totalLiability: 0,
        monthlyChange: 0,
        accountCount: 0,
        updatedAt: ""
      },
      categories: [],
      activeCategoryId: "all",
      formDialogVisible: false,
      formMode: "create",
      editingAccountId: "",
      formModel: this.createEmptyForm(),
      typeOptions: getAssetTypeOptions(),
      providerOptions: getVirtualProviderOptions(),
      formRules: {
        type: [{ required: true, message: "请选择账户类型", trigger: "change" }],
        name: [{ required: true, message: "请输入账户名称", trigger: "blur" }],
        remark: [{ required: true, message: "请输入账户备注", trigger: "blur" }],
        balance: [{ required: true, message: "请输入账户余额", trigger: "change" }],
        categoryId: [{ required: true, message: "请选择所属分类", trigger: "change" }],
        cardNo: [{ validator: function(rule, value, callback) { callback(); }, trigger: "blur" }],
        provider: [{ validator: function(rule, value, callback) { callback(); }, trigger: "change" }]
      }
    };
  },
  computed: {
    monthChangeLabel() {
      var value = Number(this.summary.monthlyChange || 0);
      var symbol = value >= 0 ? "+" : "-";
      return symbol + formatCurrency(Math.abs(value));
    },
    visibleCategories() {
      if (this.activeCategoryId === "all") {
        return this.categories;
      }

      return this.categories.filter(
        function(category) {
          return category.id === this.activeCategoryId;
        }.bind(this)
      );
    },
    categorySelectOptions() {
      return [{ label: "全部分类", value: "all" }].concat(
        this.categories.map(function(category) {
          return {
            label: category.name,
            value: category.id
          };
        })
      );
    },
    categoryOptionsForForm() {
      return this.categories.map(function(category) {
        return {
          label: category.name,
          value: category.id
        };
      });
    },
    requiresCardNo() {
      return this.formModel.type === "bankCard" || this.formModel.type === "creditCard";
    },
    requiresVirtualProvider() {
      return this.formModel.type === "virtual";
    }
  },
  created() {
    this.reloadFromSeed();
  },
  methods: {
    formatCurrency: formatCurrency,
    createEmptyForm() {
      return {
        type: "cash",
        name: "",
        remark: "",
        balance: 0,
        categoryId: "category-liquid",
        cardNo: "",
        provider: ""
      };
    },
    categoryMark(category) {
      var marks = {
        "category-liquid": "流",
        "category-investment": "投",
        "category-liability": "负"
      };

      return marks[category.id] || "资";
    },
    providerLabel(value) {
      if (value === "wechat") return "微信";
      if (value === "alipay") return "支付宝";
      return "";
    },
    typeLabel(type) {
      var target = this.typeOptions.find(function(item) {
        return item.value === type;
      });
      return target ? target.label : "账户";
    },
    ensureCategory(categoryId) {
      var target = this.categories.find(function(item) {
        return item.id === categoryId;
      });

      if (!target) {
        return this.categories[0];
      }

      return target;
    },
    recalculateSummary() {
      var totalAsset = 0;
      var totalLiability = 0;
      var monthlyChange = 0;
      var accountCount = 0;

      this.categories = this.categories.map(function(category) {
        var categoryTotal = 0;

        category.accounts.forEach(function(account) {
          var balance = toNumber(account.balance);
          var change = toNumber(account.monthlyChange);

          categoryTotal += balance;
          monthlyChange += change;
          accountCount += 1;

          if (account.direction === "liability") {
            totalLiability += balance;
          } else {
            totalAsset += balance;
          }
        });

        return Object.assign({}, category, {
          total: categoryTotal
        });
      });

      this.summary = Object.assign({}, this.summary, {
        totalAsset: totalAsset,
        totalLiability: totalLiability,
        netAsset: totalAsset - totalLiability,
        monthlyChange: monthlyChange,
        accountCount: accountCount,
        updatedAt: new Date().toLocaleString("zh-CN", { hour12: false })
      });
    },
    reloadFromSeed() {
      var overview = getAssetOverview();
      this.summary = overview.summary;
      this.categories = overview.categories;
      this.activeCategoryId = "all";
    },
    openDetail(accountId) {
      this.$router.push("/user/assets/" + accountId);
    },
    openAddDialog() {
      this.formMode = "create";
      this.editingAccountId = "";
      this.formModel = this.createEmptyForm();

      if (this.categories.length > 0) {
        this.formModel.categoryId = this.categories[0].id;
      }

      this.formDialogVisible = true;
      this.$nextTick(
        function() {
          if (this.$refs.accountFormRef) {
            this.$refs.accountFormRef.clearValidate();
          }
        }.bind(this)
      );
    },
    openEditDialog(category, account) {
      this.formMode = "edit";
      this.editingAccountId = account.id;
      this.formModel = {
        type: account.type,
        name: account.name,
        remark: account.remark,
        balance: toNumber(account.balance),
        categoryId: category.id,
        cardNo: account.cardNo || "",
        provider: account.provider || ""
      };
      this.formDialogVisible = true;
      this.$nextTick(
        function() {
          if (this.$refs.accountFormRef) {
            this.$refs.accountFormRef.clearValidate();
          }
        }.bind(this)
      );
    },
    closeFormDialog() {
      this.formDialogVisible = false;
    },
    onTypeChange() {
      this.formModel.cardNo = "";
      this.formModel.provider = "";
    },
    validateDynamicFields() {
      if (this.requiresCardNo && !this.formModel.cardNo.trim()) {
        this.$message.error("银行卡和信用卡账户必须填写卡号");
        return false;
      }

      if (this.requiresVirtualProvider && !this.formModel.provider) {
        this.$message.error("虚拟账户必须区分微信或支付宝");
        return false;
      }

      return true;
    },
    submitForm() {
      this.$refs.accountFormRef.validate(
        function(valid) {
          if (!valid) {
            return;
          }

          if (!this.validateDynamicFields()) {
            return;
          }

          if (this.formMode === "create") {
            this.createAccount();
          } else {
            this.updateAccount();
          }
        }.bind(this)
      );
    },
    createAccount() {
      var category = this.ensureCategory(this.formModel.categoryId);
      var now = Date.now();
      var newAccount = {
        id: "acc-local-" + now,
        name: this.formModel.name.trim(),
        type: this.formModel.type,
        typeLabel: this.typeLabel(this.formModel.type),
        remark: this.formModel.remark.trim(),
        balance: toNumber(this.formModel.balance),
        cardNo: this.formModel.cardNo.trim(),
        provider: this.formModel.provider,
        direction: this.formModel.type === "liability" ? "liability" : "asset",
        monthlyChange: 0
      };

      category.accounts.unshift(newAccount);
      this.recalculateSummary();
      this.formDialogVisible = false;
      this.$message.success("账户已创建（原型本地数据）");
    },
    updateAccount() {
      var targetAccount = null;

      this.categories.forEach(function(category) {
        category.accounts.forEach(function(account) {
          if (account.id === this.editingAccountId) {
            targetAccount = account;
          }
        }, this);
      }, this);

      if (!targetAccount) {
        this.$message.error("未找到要编辑的账户");
        return;
      }

      targetAccount.name = this.formModel.name.trim();
      targetAccount.remark = this.formModel.remark.trim();
      targetAccount.balance = toNumber(this.formModel.balance);
      targetAccount.cardNo = this.formModel.cardNo.trim();
      targetAccount.provider = this.formModel.provider;
      targetAccount.typeLabel = this.typeLabel(targetAccount.type);
      targetAccount.direction = targetAccount.type === "liability" ? "liability" : "asset";

      this.recalculateSummary();
      this.formDialogVisible = false;
      this.$message.success("账户信息已更新（原型本地数据）");
    }
  }
};
</script>

<style scoped>
.asset-page__filter {
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.asset-page__category-list {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.asset-page__category-card {
  border: 1px solid var(--border-color);
  border-radius: 24px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98) 0%, rgba(250, 251, 253, 0.96) 100%);
  padding: 18px;
}

.asset-page__category-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 8px;
}

.asset-page__category-title {
  display: inline-flex;
  align-items: center;
  gap: 12px;
}

.asset-page__category-title strong {
  display: block;
  font-size: 18px;
}

.asset-page__category-title p {
  margin: 6px 0 0;
  color: var(--text-subtle);
  font-size: 13px;
}

.asset-page__category-total span {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
  text-align: right;
}

.asset-page__category-total strong {
  display: block;
  margin-top: 6px;
  font-size: 24px;
}

.asset-page__account-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 10px;
}

.asset-page__account-row {
  display: grid;
  grid-template-columns: 1fr auto auto;
  gap: 14px;
  align-items: center;
  padding: 18px;
  border-radius: 18px;
  border: 1px solid rgba(229, 231, 235, 0.95);
  background: rgba(255, 255, 255, 0.95);
}

.asset-page__account-main {
  display: inline-flex;
  align-items: center;
  gap: 14px;
  flex-wrap: wrap;
}

.asset-page__account-main strong {
  display: block;
  font-size: 16px;
}

.asset-page__account-main p {
  margin: 6px 0 0;
  color: var(--text-subtle);
  font-size: 13px;
}

.asset-page__account-balance small {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
}

.asset-page__account-balance strong {
  display: block;
  margin-top: 6px;
  font-size: 20px;
}

.asset-page__account-balance span {
  display: block;
  margin-top: 6px;
  font-size: 12px;
}

.asset-page__actions {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.asset-page__dialog-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.asset-page__dialog-actions {
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

@media (max-width: 1100px) {
  .asset-page__account-row {
    grid-template-columns: 1fr;
  }

  .asset-page__actions {
    justify-content: flex-start;
  }
}
</style>
