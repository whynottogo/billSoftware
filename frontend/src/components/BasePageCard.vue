<template>
  <section class="page-card base-page-card">
    <div class="base-page-card__header">
      <div>
        <h2 class="page-title">{{ title }}</h2>
        <p v-if="description" class="page-description">{{ description }}</p>
      </div>
      <slot name="extra" />
    </div>
    <div class="base-page-card__body">
      <slot />
    </div>
  </section>
</template>

<script>
export default {
  name: "BasePageCard",
  props: {
    title: {
      type: String,
      required: true
    },
    description: {
      type: String,
      default: ""
    }
  }
};
</script>

<style scoped>
.base-page-card {
  padding: 28px;
}

.base-page-card__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 24px;
}

.base-page-card__body {
  display: flex;
  flex-direction: column;
  gap: 20px;
}
</style>

