<template>
  <div class="icon-picker" ref="pickerRef">
    <button type="button" class="icon-picker__trigger" @click="toggleDropdown">
      <CategoryIcon v-if="selected" :icon="selected" :size="40" />
      <span v-else class="icon-picker__placeholder">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="icon-picker__placeholder-icon">
          <circle cx="12" cy="12" r="10" />
          <line x1="12" y1="8" x2="12" y2="16" />
          <line x1="8" y1="12" x2="16" y2="12" />
        </svg>
      </span>
      <svg class="icon-picker__arrow" :class="{ 'is-open': visible }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="6 9 12 15 18 9" />
      </svg>
    </button>
    <transition name="icon-picker-dropdown">
      <div v-if="visible" class="icon-picker__dropdown">
        <div class="icon-picker__grid">
          <button
            v-for="item in icons"
            :key="item.key"
            type="button"
            class="icon-picker__item"
            :class="{ 'is-selected': item.key === selected }"
            :style="{ background: item.color }"
            @click="onSelect(item.key)"
          >
            <component :is="item.icon" class="icon-picker__svg" />
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script>
import { getAllCategoryIcons } from "@/icons/categoryIcons";
import CategoryIcon from "@/components/CategoryIcon.vue";

export default {
  name: "IconPicker",
  components: { CategoryIcon },
  props: {
    selected: { type: String, default: "" }
  },
  emits: ["select"],
  data() {
    return {
      visible: false
    };
  },
  computed: {
    icons() {
      return getAllCategoryIcons();
    }
  },
  mounted() {
    document.addEventListener("mousedown", this.onOutsideClick);
  },
  beforeUnmount() {
    document.removeEventListener("mousedown", this.onOutsideClick);
  },
  methods: {
    toggleDropdown() {
      this.visible = !this.visible;
    },
    onSelect(key) {
      this.$emit("select", key);
      this.visible = false;
    },
    onOutsideClick(e) {
      if (this.$refs.pickerRef && !this.$refs.pickerRef.contains(e.target)) {
        this.visible = false;
      }
    }
  }
};
</script>

<style scoped>
.icon-picker {
  position: relative;
}

.icon-picker__trigger {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  height: 40px;
  padding: 0 8px 0 0;
  border: 1px solid var(--border-color, #e5e7eb);
  border-radius: 10px;
  background: #ffffff;
  cursor: pointer;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.icon-picker__trigger:hover {
  border-color: var(--brand-color, #6366f1);
}

.icon-picker__placeholder {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
}

.icon-picker__placeholder-icon {
  width: 20px;
  height: 20px;
  color: #c0c4cc;
}

.icon-picker__arrow {
  width: 14px;
  height: 14px;
  color: #c0c4cc;
  transition: transform 0.2s;
}

.icon-picker__arrow.is-open {
  transform: rotate(180deg);
}

.icon-picker__dropdown {
  position: absolute;
  top: calc(100% + 6px);
  left: 50%;
  transform: translateX(-50%);
  z-index: 2000;
  background: #ffffff;
  border: 1px solid var(--border-color, #e5e7eb);
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  padding: 12px;
}

.icon-picker__grid {
  display: grid;
  grid-template-columns: repeat(6, 44px);
  gap: 8px;
}

.icon-picker__item {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  border: 2px solid transparent;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: border-color 0.15s, transform 0.15s;
  color: #ffffff;
}

.icon-picker__item:hover {
  transform: scale(1.15);
}

.icon-picker__item.is-selected {
  border-color: var(--text-main, #1f2937);
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.3);
}

.icon-picker__svg {
  width: 22px;
  height: 22px;
}

.icon-picker-dropdown-enter-active,
.icon-picker-dropdown-leave-active {
  transition: opacity 0.15s, transform 0.15s;
}

.icon-picker-dropdown-enter-from,
.icon-picker-dropdown-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-4px);
}
</style>
