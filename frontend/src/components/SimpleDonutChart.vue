<template>
  <div class="simple-donut">
    <svg :viewBox="'0 0 ' + size + ' ' + size" :width="size" :height="size" class="simple-donut__svg">
      <g :transform="'rotate(-90 ' + half + ' ' + half + ')'">
        <circle
          class="simple-donut__track"
          :cx="half"
          :cy="half"
          :r="radius"
          :stroke-width="strokeWidth"
          fill="none"
        />
        <circle
          v-for="item in chartSegments"
          :key="item.name"
          class="simple-donut__segment"
          :cx="half"
          :cy="half"
          :r="radius"
          :stroke="item.color"
          :stroke-width="strokeWidth"
          :stroke-dasharray="item.dasharray"
          :stroke-dashoffset="item.dashoffset"
          stroke-linecap="round"
          fill="none"
        />
      </g>
    </svg>

    <div class="simple-donut__center">
      <span>{{ centerTitle }}</span>
      <strong>{{ centerValue || totalLabel }}</strong>
    </div>
  </div>
</template>

<script>
export default {
  name: "SimpleDonutChart",
  props: {
    segments: {
      type: Array,
      required: true
    },
    centerTitle: {
      type: String,
      default: ""
    },
    centerValue: {
      type: String,
      default: ""
    },
    size: {
      type: Number,
      default: 220
    },
    strokeWidth: {
      type: Number,
      default: 28
    }
  },
  computed: {
    total() {
      return this.segments.reduce(function(sum, item) {
        return sum + Number(item.value || 0);
      }, 0);
    },
    totalLabel() {
      return "¥" + this.total.toLocaleString("zh-CN");
    },
    half() {
      return this.size / 2;
    },
    radius() {
      return this.half - this.strokeWidth / 2;
    },
    circumference() {
      return 2 * Math.PI * this.radius;
    },
    chartSegments() {
      var offset = 0;
      var circumference = this.circumference;
      var total = this.total || 1;

      return this.segments.map(function(item) {
        var length = (Number(item.value || 0) / total) * circumference;
        var segment = {
          name: item.name,
          color: item.color,
          dasharray: length + " " + circumference,
          dashoffset: -offset
        };

        offset += length;
        return segment;
      });
    }
  }
};
</script>

<style scoped>
.simple-donut {
  position: relative;
  width: fit-content;
  margin: 0 auto;
}

.simple-donut__svg {
  display: block;
}

.simple-donut__track {
  stroke: rgba(229, 231, 235, 0.9);
}

.simple-donut__center {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
}

.simple-donut__center span {
  color: var(--text-muted);
  font-size: 12px;
}

.simple-donut__center strong {
  margin-top: 8px;
  font-size: 28px;
  line-height: 1.1;
}
</style>
