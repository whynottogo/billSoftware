<template>
  <div class="finance-page asset-detail-page">
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
          <button class="finance-button finance-button--primary" @click="submitOperation">确认写入记录</button>
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
          <button class="finance-button finance-button--primary" @click="submitSettings">保存设置</button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import {
  formatCurrency,
  getAssetAccountDetail,
  getMonthLabel,
  monthKeysFromRecords
} from "@/utils/userAssetMock";

function nowLabel() {
  var now = new Date();
  var MM = String(now.getMonth() + 1).padStart(2, "0");
  var DD = String(now.getDate()).padStart(2, "0");
  var HH = String(now.getHours()).padStart(2, "0");
  var mm = String(now.getMinutes()).padStart(2, "0");

  return MM + "-" + DD + " " + HH + ":" + mm;
}

function currentMonthKey() {
  var now = new Date();
  var yyyy = now.getFullYear();
  var MM = String(now.getMonth() + 1).padStart(2, "0");
  return yyyy + "-" + MM;
}

function nextRecordId() {
  return "record-local-" + Date.now() + "-" + Math.round(Math.random() * 1000);
}

export default {
  name: "UserAssetDetail",
  data() {
    return {
      account: {},
      categoryName: "",
      records: [],
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
    formatCurrency: formatCurrency,
    getMonthLabel: getMonthLabel,
    loadDetail() {
      var detail = getAssetAccountDetail(this.$route.params.accountId);

      this.account = detail.account;
      this.categoryName = detail.category.name;
      this.records = detail.records;
      this.selectedMonthKey = "all";
      this.resetOperationForm();
      this.settingsForm = {
        remark: this.account.remark,
        targetBalance: Number(this.account.balance || 0)
      };
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
    prependRecord(action, change, balanceAfter, note, source) {
      var record = {
        id: nextRecordId(),
        monthKey: currentMonthKey(),
        dateLabel: nowLabel(),
        action: action,
        change: change,
        balanceAfter: balanceAfter,
        note: note,
        source: source
      };

      this.records.unshift(record);
      this.account.balance = balanceAfter;
      return record;
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

          var currentBalance = Number(this.account.balance || 0);
          var actionType = this.operationForm.actionType;
          var change = 0;
          var balanceAfter = currentBalance;
          var actionLabel = "调整";

          if (actionType === "adjust") {
            balanceAfter = amount;
            change = amount - currentBalance;
            actionLabel = "调整";
          } else if (actionType === "increase") {
            change = amount;
            balanceAfter = currentBalance + amount;
            actionLabel = "增加";
          } else {
            change = -amount;
            balanceAfter = currentBalance - amount;
            actionLabel = "减少";
          }

          this.prependRecord(actionLabel, change, balanceAfter, this.operationForm.note.trim(), "手动操作");
          this.selectedMonthKey = currentMonthKey();
          this.$message.success("余额操作已记录");
          this.resetOperationForm();
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

          var newBalance = Number(this.settingsForm.targetBalance || 0);
          var currentBalance = Number(this.account.balance || 0);
          var change = newBalance - currentBalance;

          this.account.remark = this.settingsForm.remark.trim();
          this.prependRecord("调整", change, newBalance, "通过账户设置修改余额", "账户设置");
          this.selectedMonthKey = currentMonthKey();
          this.settingsDialogVisible = false;
          this.$message.success("账户设置已保存，并生成调整记录");
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
