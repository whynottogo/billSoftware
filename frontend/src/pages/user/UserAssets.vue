<template>
  <div class="finance-page asset-page">
    <div v-if="errorMessage" class="finance-inline-notice">
      <span>{{ errorMessage }}</span>
      <button type="button" @click="loadOverview">重新加载</button>
    </div>

    <section class="finance-toolbar">
      <div class="finance-toolbar__meta">
        <h1>资产管家</h1>
        <p>把账户放在同一张资产地图里，先看全局，再处理单账户细节。</p>
      </div>
      <div class="finance-toolbar__actions">
        <button class="finance-button finance-button--ghost" @click="loadOverview">刷新资产</button>
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
          <strong>{{ summary.updatedAt || "--" }}</strong>
          <small>本页已切到真实接口数据</small>
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
              <strong :class="isLiabilityCategory(category) ? 'finance-tone-expense' : 'finance-tone-income'">
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
            <button class="finance-button finance-button--primary" :disabled="isSubmitting" @click="submitForm">
              {{ isSubmitting ? "提交中..." : (formMode === "create" ? "创建账户" : "保存修改") }}
            </button>
          </div>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import {
  buildAssetPayload,
  buildUserAssetsError,
  createUserAsset,
  formatAssetCurrency,
  getAssetTypeOptions,
  getUserAssets,
  getVirtualProviderOptions,
  normalizeAssetOverviewPayload,
  updateUserAsset
} from "@/api/userAssets";

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
      isLoading: false,
      isSubmitting: false,
      errorMessage: "",
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
      return symbol + this.formatCurrency(Math.abs(value));
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
    this.loadOverview();
  },
  methods: {
    formatCurrency: formatAssetCurrency,
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
      var name = String(category.name || "");
      var hasLiability = (category.accounts || []).some(function(account) {
        return account.direction === "liability";
      });

      if (String(category.id || "") === "category-liability" || name.indexOf("负债") >= 0 || hasLiability) {
        return "负";
      }

      if (String(category.id || "") === "category-investment" || name.indexOf("投资") >= 0) {
        return "投";
      }

      if (String(category.id || "") === "category-liquid" || name.indexOf("流动") >= 0) {
        return "流";
      }

      return name.slice(0, 1) || "资";
    },
    providerLabel(value) {
      if (value === "wechat") return "微信";
      if (value === "alipay") return "支付宝";
      return "";
    },
    isLiabilityCategory(category) {
      return (category.accounts || []).some(function(account) {
        return account.direction === "liability";
      });
    },
    typeLabel(type) {
      var target = this.typeOptions.find(function(item) {
        return item.value === type;
      });
      return target ? target.label : "账户";
    },
    loadOverview() {
      this.isLoading = true;
      this.errorMessage = "";

      return getUserAssets()
        .then(
          function(result) {
            var overview = normalizeAssetOverviewPayload(result);
            this.summary = overview.summary;
            this.categories = overview.categories;

            var hasActiveCategory = this.categories.some(function(category) {
              return category.id === this.activeCategoryId;
            }, this);

            if (!hasActiveCategory) {
              this.activeCategoryId = "all";
            }
          }.bind(this)
        )
        .catch(
          function(error) {
            this.errorMessage = buildUserAssetsError(error, "资产数据加载失败，请稍后重试");
          }.bind(this)
        )
        .finally(
          function() {
            this.isLoading = false;
          }.bind(this)
        );
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
        balance: Number(account.balance || 0),
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
      this.isSubmitting = true;

      return createUserAsset(buildAssetPayload(this.formModel))
        .then(
          function() {
            return this.loadOverview();
          }.bind(this)
        )
        .then(
          function() {
            this.formDialogVisible = false;
            this.$message.success("账户已创建");
          }.bind(this)
        )
        .catch(
          function(error) {
            this.$message.error(buildUserAssetsError(error, "创建账户失败，请稍后重试"));
          }.bind(this)
        )
        .finally(
          function() {
            this.isSubmitting = false;
          }.bind(this)
        );
    },
    updateAccount() {
      if (!this.editingAccountId) {
        this.$message.error("未找到要编辑的账户");
        return;
      }

      this.isSubmitting = true;

      return updateUserAsset(this.editingAccountId, buildAssetPayload(this.formModel))
        .then(
          function() {
            return this.loadOverview();
          }.bind(this)
        )
        .then(
          function() {
            this.formDialogVisible = false;
            this.$message.success("账户信息已更新");
          }.bind(this)
        )
        .catch(
          function(error) {
            this.$message.error(buildUserAssetsError(error, "更新账户失败，请稍后重试"));
          }.bind(this)
        )
        .finally(
          function() {
            this.isSubmitting = false;
          }.bind(this)
        );
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
