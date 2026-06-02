<template>
  <span class="category-icon" :style="iconStyle">
    <component v-if="iconData" :is="iconData.icon" class="category-icon__svg" />
    <span v-else class="category-icon__fallback">{{ fallbackText }}</span>
  </span>
</template>

<script>
import { getCategoryIcon } from "@/icons/categoryIcons";

export default {
  name: "CategoryIcon",
  props: {
    icon: { type: String, default: "" },
    badge: { type: String, default: "" },
    size: { type: Number, default: 36 }
  },
  computed: {
    iconData() {
      return getCategoryIcon(this.icon);
    },
    iconStyle() {
      var bgColor = this.iconData ? this.iconData.color : "#E5E7EB";
      return {
        width: this.size + "px",
        height: this.size + "px",
        minWidth: this.size + "px",
        borderRadius: Math.round(this.size * 0.3) + "px",
        background: bgColor
      };
    },
    fallbackText() {
      return this.badge || "-";
    }
  }
};
</script>

<style scoped>
.category-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #ffffff;
}

.category-icon__svg {
  width: 55%;
  height: 55%;
}

.category-icon__fallback {
  font-size: 14px;
  font-weight: 700;
  color: #ffffff;
}
</style>
