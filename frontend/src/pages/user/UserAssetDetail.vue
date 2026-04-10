<template>
  <div class="finance-page asset-detail-page">
    <div v-if="errorMessage" class="finance-inline-notice">
      <span>{{ errorMessage }}</span>
      <button type="button" @click="loadDetail">重新加载</button>
    </div>

    <section class="finance-toolbar">
      <div class="finance-toolbar__actions">
        <button class="finance-button finance-button--ghost" @click="goBack">返回资产总览</button>
        <button class="finance-button finance-button--ghost" @click="openSettingsDialog">账户设置</button>
      </div>
      <div class="finance-toolbar__actions">
        <span class="finance-pill">{{ categoryName }} · {{ account.typeLabel }}</span>
      </div>
    </section>

    <section class="finance-hero">
      <div class="finance-hero__eyebrow">
        <span class="finance-hero__eyebrow-mark">户</span>
        <span>{{ account.name }}</span>
      </div>
      <div class="finance-hero__headline">
        <h1>{{ formatCurrency(account.balance) }}</h1>
        <p>{{ account.remark }}</p>
      </div>
      <div class="finance-stat-grid finance-stat-grid--triple">
        <article class="finance-stat-card">
          <span>本月净变动</span>
          <strong :class="monthNetChange >= 0 ? 'finance-tone-income' : 'finance-tone-expense'">
            {{ monthNetChange >= 0 ? "+" : "-" }}{{ formatCurrency(Math.abs(monthNetChange)) }}
          </strong>
          <small>{{ selectedMonthLabel }}</small>
        </article>
        <article class="finance-stat-card">
          <span>记录条数</span>
          <strong>{{ filteredRecords.length }} 条</strong>
          <small>按月筛选可回溯历史</small>
        </article>
        <article class="finance-stat-card">
          <span>账户标识</span>
          <strong>{{ account.cardNo || providerText || "无卡号" }}</strong>
          <small>详情页支持调整、增加、减少余额</small>
        </article>
      </div>
    </section>

    <section class="finance-grid-2">
      <article class="page-card finance-panel">
        <header class="finance-panel__header">
          <div>
            <h3>余额操作</h3>
            <p>支持调整、增加、减少。所有操作都会生成余额变动记录。</p>
          </div>
        </header>

        <el-form ref="operationFormRef" :model="operationForm" :rules="operationRules" label-width="104px" status-icon>
          <el-form-item label="操作类型" prop="actionType">
            <el-radio-group v-model="operationForm.actionType">
              <el-radio-button label="adjust">调整余额</el-radio-button>
              <el-radio-button label="increase">增加余额</el-radio-button>
              <el-radio-button label="decrease">减少余额</el-radio-button>
            </el-radio-group>
          </el-form-item>

          <el-form-item :label="amountLabel" prop="amount">
            <el-input-number
              v-model="operationForm.amount"
              :precision="2"
              :step="100"
              :min="0"
              style="width: 100%;"
            />
          </el-form-item>

          <el-form-item label="操作备注" prop="note">
            <el-input
              v-model="operationForm.note"
              type="textarea"
              :autosize="{ minRows: 2, maxRows: 3 }"
              maxlength="40"
              show-word-limit
            />
          </el-form-item>
        </el-form>

        <div class="asset-detail-page__op-actions">
          <button class="finance-button finance-button--ghost" @click="resetOperationForm">重置</button>
          <button class="finance-button finance-button--primary" :disabled="isSubmitting" @click="submitOperation">
            {{ isSubmitting ? "写入中..." : "确认写入记录" }}
          </button>
        </div>
      </article>

      <article class="page-card finance-panel">
        <header class="finance-panel__header">
          <div>
            <h3>余额变动记录</h3>
            <p>记录只属于资产模块，不与收支流水共享。</p>
          </div>
          <div class="asset-detail-page__month-filter">
            <span class="finance-inline-tag">筛选月份</span>
            <el-select v-model="selectedMonthKey" size="large" style="min-width: 170px;">
              <el-option label="全部月份" value="all" />
              <el-option
                v-for="monthKey in monthKeys"
                :key="monthKey"
                :label="getMonthLabel(monthKey)"
                :value="monthKey"
              />
            </el-select>
          </div>
        </header>

        <div class="asset-detail-page__record-list">
          <article v-for="record in filteredRecords" :key="record.id" class="asset-detail-page__record-row">
            <div class="asset-detail-page__record-main">
              <strong>{{ record.action }}</strong>
              <p>{{ record.note }}</p>
              <small>{{ record.dateLabel }} · {{ getMonthLabel(record.monthKey) }} · {{ record.source }}</small>
            </div>
            <div class="asset-detail-page__record-change">
              <strong :class="record.change >= 0 ? 'finance-tone-income' : 'finance-tone-expense'">
                {{ record.change >= 0 ? "+" : "" }}{{ formatCurrency(record.change) }}
              </strong>
              <span>余额 {{ formatCurrency(record.balanceAfter) }}</span>
            </div>
          </article>
          <div v-if="filteredRecords.length === 0" class="asset-detail-page__empty">当前月份暂无记录</div>
        </div>
      </article>
    </section>

    <el-dialog v-model="settingsDialogVisible" title="账户设置" width="520px" destroy-on-close>
      <el-form ref="settingsFormRef" :model="settingsForm" :rules="settingsRules" label-width="110px" status-icon>
        <el-form-item label="账户备注" prop="remark">
          <el-input v-model="settingsForm.remark" maxlength="40" show-word-limit />
        </el-form-item>
        <el-form-item label="设置新余额" prop="targetBalance">
          <el-input-number v-model="settingsForm.targetBalance" :precision="2" :step="100" :min="0" style="width: 100%;" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="asset-detail-page__op-actions">
          <button class="finance-button finance-button--ghost" @click="settingsDialogVisible = false">取消</button>
          <button class="finance-button finance-button--primary" :disabled="isSubmitting" @click="submitSettings">
            {{ isSubmitting ? "保存中..." : "保存设置" }}
          </button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import {
  buildAssetOperationPayload,
  buildAssetUpdatePayload,
  buildUserAssetsError,
  createUserAssetOperation,
  formatAssetCurrency,
  getMonthLabel,
  getUserAssetDetail,
  monthKeysFromRecords,
  normalizeAssetDetailPayload,
  updateUserAsset
} from "@/api/userAssets";

function currentMonthKey() {
  var now = new Date();
  var yyyy = now.getFullYear();
  var MM = String(now.getMonth() + 1).padStart(2, "0");
  return yyyy + "-" + MM;
}

export default {
  name: "UserAssetDetail",
  data() {
    return {
      account: {
        id: "",
        name: "",
        type: "cash",
        typeLabel: "账户",
        remark: "",
        balance: 0,
        cardNo: "",
        provider: "",
        direction: "asset",
        categoryId: "",
        categoryName: ""
      },
      categoryName: "",
      records: [],
      isLoading: false,
      isSubmitting: false,
      errorMessage: "",
      selectedMonthKey: "all",
      operationForm: {
        actionType: "adjust",
        amount: null,
        note: ""
      },
      operationRules: {
        actionType: [{ required: true, message: "请选择操作类型", trigger: "change" }],
        amount: [{ required: true, message: "请输入金额", trigger: "change" }],
        note: [{ required: true, message: "请输入备注", trigger: "blur" }]
      },
      settingsDialogVisible: false,
      settingsForm: {
        remark: "",
        targetBalance: null
      },
      settingsRules: {
        remark: [{ required: true, message: "请输入账户备注", trigger: "blur" }],
        targetBalance: [{ required: true, message: "请输入设置余额", trigger: "change" }]
      }
    };
  },
  computed: {
    monthKeys() {
      return monthKeysFromRecords(this.records);
    },
    selectedMonthLabel() {
      if (this.selectedMonthKey === "all") {
        return "全部月份";
      }
      return getMonthLabel(this.selectedMonthKey);
    },
    filteredRecords() {
      if (this.selectedMonthKey === "all") {
        return this.records;
      }

      return this.records.filter(
        function(record) {
          return record.monthKey === this.selectedMonthKey;
        }.bind(this)
      );
    },
    monthNetChange() {
      return this.filteredRecords.reduce(function(sum, record) {
        return sum + Number(record.change || 0);
      }, 0);
    },
    amountLabel() {
      if (this.operationForm.actionType === "adjust") {
        return "调整后余额";
      }
      return "变动金额";
    },
    providerText() {
      if (this.account.provider === "wechat") return "微信";
      if (this.account.provider === "alipay") return "支付宝";
      return "";
    }
  },
  watch: {
    "$route.params.accountId": function() {
      this.loadDetail();
    }
  },
  created() {
    this.loadDetail();
  },
  methods: {
    formatCurrency: formatAssetCurrency,
    getMonthLabel: getMonthLabel,
    loadDetail(options) {
      var preserveSelectedMonth = options && options.preserveSelectedMonth;
      var nextMonthKey = options && options.nextMonthKey;

      this.isLoading = true;
      this.errorMessage = "";

      return getUserAssetDetail(this.$route.params.accountId)
        .then(
          function(result) {
            var detail = normalizeAssetDetailPayload(result);

            this.account = detail.account;
            this.categoryName = detail.category.name;
            this.records = detail.records;
            this.selectedMonthKey = preserveSelectedMonth
              ? (nextMonthKey || this.selectedMonthKey)
              : "all";
            this.resetOperationForm();
            this.settingsForm = {
              remark: this.account.remark,
              targetBalance: Number(this.account.balance || 0)
            };
          }.bind(this)
        )
        .catch(
          function(error) {
            this.errorMessage = buildUserAssetsError(error, "账户详情加载失败，请稍后重试");
          }.bind(this)
        )
        .finally(
          function() {
            this.isLoading = false;
          }.bind(this)
        );
    },
    goBack() {
      this.$router.push("/user/assets");
    },
    resetOperationForm() {
      this.operationForm = {
        actionType: "adjust",
        amount: null,
        note: ""
      };

      if (this.$refs.operationFormRef) {
        this.$refs.operationFormRef.clearValidate();
      }
    },
    submitOperation() {
      this.$refs.operationFormRef.validate(
        function(valid) {
          if (!valid) {
            return;
          }

          var amount = Number(this.operationForm.amount || 0);
          if (amount <= 0) {
            this.$message.error("金额必须大于 0");
            return;
          }

          this.isSubmitting = true;

          createUserAssetOperation(this.account.id, buildAssetOperationPayload({
            actionType: this.operationForm.actionType,
            amount: amount,
            note: this.operationForm.note.trim()
          }))
            .then(
              function() {
                return this.loadDetail({
                  preserveSelectedMonth: true,
                  nextMonthKey: currentMonthKey()
                });
              }.bind(this)
            )
            .then(
              function() {
                this.$message.success("余额操作已记录");
              }.bind(this)
            )
            .catch(
              function(error) {
                this.$message.error(buildUserAssetsError(error, "余额操作失败，请稍后重试"));
              }.bind(this)
            )
            .finally(
              function() {
                this.isSubmitting = false;
              }.bind(this)
            );
        }.bind(this)
      );
    },
    openSettingsDialog() {
      this.settingsForm = {
        remark: this.account.remark,
        targetBalance: Number(this.account.balance || 0)
      };
      this.settingsDialogVisible = true;
      this.$nextTick(
        function() {
          if (this.$refs.settingsFormRef) {
            this.$refs.settingsFormRef.clearValidate();
          }
        }.bind(this)
      );
    },
    submitSettings() {
      this.$refs.settingsFormRef.validate(
        function(valid) {
          if (!valid) {
            return;
          }

          var remark = this.settingsForm.remark.trim();
          var newBalance = Number(this.settingsForm.targetBalance || 0);
          var currentBalance = Number(this.account.balance || 0);
          var hasRemarkChange = remark !== String(this.account.remark || "");
          var hasBalanceChange = newBalance !== currentBalance;

          if (!hasRemarkChange && !hasBalanceChange) {
            this.settingsDialogVisible = false;
            this.$message.success("没有需要保存的修改");
            return;
          }

          this.isSubmitting = true;

          var chain = Promise.resolve();

          if (hasRemarkChange) {
            chain = chain.then(
              function() {
                return updateUserAsset(this.account.id, buildAssetUpdatePayload(this.account, {
                  remark: remark,
                  balance: currentBalance
                }));
              }.bind(this)
            );
          }

          if (hasBalanceChange) {
            chain = chain.then(
              function() {
                return createUserAssetOperation(this.account.id, buildAssetOperationPayload({
                  actionType: "adjust",
                  amount: newBalance,
                  note: "通过账户设置修改余额"
                }));
              }.bind(this)
            );
          }

          chain
            .then(
              function() {
                return this.loadDetail({
                  preserveSelectedMonth: true,
                  nextMonthKey: hasBalanceChange ? currentMonthKey() : this.selectedMonthKey
                });
              }.bind(this)
            )
            .then(
              function() {
                this.settingsDialogVisible = false;
                this.$message.success("账户设置已保存");
              }.bind(this)
            )
            .catch(
              function(error) {
                this.$message.error(buildUserAssetsError(error, "保存账户设置失败，请稍后重试"));
              }.bind(this)
            )
            .finally(
              function() {
                this.isSubmitting = false;
              }.bind(this)
            );
        }.bind(this)
      );
    }
  }
};
</script>

<style scoped>
.asset-detail-page__op-actions {
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.asset-detail-page__month-filter {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.asset-detail-page__record-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.asset-detail-page__record-row {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 16px;
  align-items: center;
  border: 1px solid var(--border-color);
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.96);
  padding: 14px 16px;
}

.asset-detail-page__record-main strong {
  display: block;
  font-size: 16px;
}

.asset-detail-page__record-main p {
  margin: 6px 0 0;
  color: var(--text-subtle);
  line-height: 1.6;
}

.asset-detail-page__record-main small {
  display: block;
  margin-top: 8px;
  color: var(--text-muted);
  font-size: 12px;
}

.asset-detail-page__record-change {
  text-align: right;
}

.asset-detail-page__record-change strong {
  display: block;
  font-size: 20px;
}

.asset-detail-page__record-change span {
  display: block;
  margin-top: 6px;
  color: var(--text-subtle);
  font-size: 12px;
}

.asset-detail-page__empty {
  min-height: 120px;
  border-radius: 16px;
  border: 1px dashed var(--border-color);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
  background: rgba(255, 255, 255, 0.7);
}

@media (max-width: 980px) {
  .asset-detail-page__record-row {
    grid-template-columns: 1fr;
  }

  .asset-detail-page__record-change {
    text-align: left;
  }
}
</style>
