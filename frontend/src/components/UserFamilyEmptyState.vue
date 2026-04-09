<template>
  <section class="page-card family-empty">
    <div class="family-empty__badge">家庭</div>
    <h3>还没有家庭记录</h3>
    <p>你可以先创建一个家庭，或者通过家庭 ID / 邀请链接加入已有家庭。</p>
    <div class="family-empty__actions">
      <button class="finance-button finance-button--primary" @click="$emit('create')">创建家庭</button>
      <button class="finance-button finance-button--ghost" @click="$emit('join')">加入家庭</button>
    </div>
  </section>
</template>

<script>
export default {
  name: "UserFamilyEmptyState",
  emits: ["create", "join"]
};
</script>

<style scoped>
.family-empty {
  padding: 30px;
  text-align: center;
}

.family-empty__badge {
  width: 56px;
  height: 56px;
  margin: 0 auto;
  border-radius: 18px;
  background: rgba(246, 211, 74, 0.22);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
}

.family-empty h3 {
  margin: 16px 0 0;
  font-size: 24px;
}

.family-empty p {
  margin: 10px auto 0;
  max-width: 520px;
  color: var(--text-subtle);
  line-height: 1.7;
}

.family-empty__actions {
  margin-top: 18px;
  display: inline-flex;
  gap: 12px;
}
</style>
