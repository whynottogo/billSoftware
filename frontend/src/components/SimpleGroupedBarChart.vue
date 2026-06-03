<template>
  <div class="simple-bar-chart">
    <svg :viewBox="'0 0 ' + svgWidth + ' ' + height" class="simple-bar-chart__svg" preserveAspectRatio="none">
      <g v-for="tick in ticks" :key="'tick-' + tick.value">
        <line
          :x1="padding.left"
          :x2="svgWidth - padding.right"
          :y1="yFor(tick.value)"
          :y2="yFor(tick.value)"
          class="simple-bar-chart__grid"
        />
        <text
          :x="padding.left - 10"
          :y="yFor(tick.value) + 4"
          class="simple-bar-chart__y-label"
          text-anchor="end"
        >
          {{ tick.label }}
        </text>
      </g>

      <g v-for="(label, index) in labels" :key="'group-' + label + '-' + index">
        <rect
          v-for="bar in barsFor(index)"
          :key="bar.key"
          :x="bar.x"
          :y="bar.y"
          :width="bar.width"
          :height="bar.height"
          :fill="bar.color"
          class="simple-bar-chart__bar"
          rx="6"
          ry="6"
          @mouseenter="setActiveBar(bar)"
          @mousemove="setActiveBar(bar)"
          @mouseleave="clearActiveBar"
        />
        <text
          :x="groupCenter(index)"
          :y="height - 12"
          class="simple-bar-chart__x-label"
          text-anchor="middle"
        >
          {{ label }}
        </text>
      </g>

      <g v-if="activeBar" :transform="tooltipTransform()" class="simple-bar-chart__tooltip">
        <rect width="124" height="42" rx="10" ry="10" />
        <text x="12" y="17" class="simple-bar-chart__tooltip-label">{{ activeBar.label }}</text>
        <text x="12" y="32" class="simple-bar-chart__tooltip-value">{{ activeBar.name }} · {{ formatCurrency(activeBar.value) }}</text>
      </g>
    </svg>

    <div class="simple-bar-chart__legend">
      <span v-for="item in series" :key="item.name" class="simple-bar-chart__legend-item">
        <i :style="{ backgroundColor: item.color }"></i>
        {{ item.name }}
      </span>
    </div>
  </div>
</template>

<script>
export default {
  name: "SimpleGroupedBarChart",
  props: {
    labels: {
      type: Array,
      required: true
    },
    series: {
      type: Array,
      required: true
    },
    height: {
      type: Number,
      default: 300
    },
    valueUnit: {
      type: String,
      default: "currency"
    }
  },
  data() {
    return {
      svgWidth: 760,
      padding: {
        top: 20,
        right: 18,
        bottom: 42,
        left: 52
      },
      activeBar: null
    };
  },
  computed: {
    maxValue() {
      var max = 0;

      this.series.forEach(function(item) {
        item.values.forEach(function(value) {
          if (value > max) {
            max = value;
          }
        });
      });

      if (!max) {
        return 1;
      }

      var step = this.tickStep;
      return Math.ceil(max / step) * step;
    },
    tickStep() {
      var max = 0;

      this.series.forEach(function(item) {
        item.values.forEach(function(value) {
          if (value > max) {
            max = value;
          }
        });
      });

      if (!max) {
        return 1;
      }

      return this.niceStep(max / 5);
    },
    plotWidth() {
      return this.svgWidth - this.padding.left - this.padding.right;
    },
    plotHeight() {
      return this.height - this.padding.top - this.padding.bottom;
    },
    ticks() {
      var step = this.tickStep;
      var list = [];
      var value;

      for (value = 0; value <= this.maxValue + step * 0.5; value += step) {
        list.push({
          value: value,
          label: this.formatCurrency(value)
        });
      }

      return list;
    },
    groupWidth() {
      return this.plotWidth / this.labels.length;
    },
    barWidth() {
      return Math.min(12, this.groupWidth / (this.series.length + 2));
    }
  },
  methods: {
    groupCenter(index) {
      return this.padding.left + this.groupWidth * index + this.groupWidth / 2;
    },
    yFor(value) {
      return this.padding.top + this.plotHeight - (Number(value || 0) / this.maxValue) * this.plotHeight;
    },
    barsFor(index) {
      var self = this;
      var totalWidth = this.barWidth * this.series.length + 5 * (this.series.length - 1);
      var start = this.groupCenter(index) - totalWidth / 2;
      var baseline = this.padding.top + this.plotHeight;
      var label = this.labels[index] || "";

      return this.series.map(function(item, seriesIndex) {
        var value = item.values[index] || 0;
        var x = start + seriesIndex * (self.barWidth + 5);
        var y = self.yFor(value);
        var height = self.padding.top + self.plotHeight - y;
        var displayHeight = Math.max(height, value > 0 ? 4 : 2);

        return {
          key: item.name + "-" + index,
          x: x,
          y: baseline - displayHeight,
          width: self.barWidth,
          height: displayHeight,
          color: item.color,
          name: item.name,
          value: value,
          label: label
        };
      });
    },
    setActiveBar(bar) {
      this.activeBar = {
        name: bar.name,
        label: bar.label || "-",
        value: bar.value,
        color: bar.color,
        x: bar.x + bar.width / 2,
        y: bar.y
      };
    },
    clearActiveBar() {
      this.activeBar = null;
    },
    formatCurrency(value) {
      var number = Number(value || 0);

      if (this.valueUnit === "count") {
        return Math.round(number).toLocaleString("zh-CN");
      }

      if (!number) {
        return "¥0";
      }

      if (number >= 10000) {
        return "¥" + this.trimDecimal(number / 1000, 1) + "k";
      }

      if (number >= 1000) {
        return "¥" + this.trimDecimal(number / 1000, number % 1000 === 0 ? 0 : 1) + "k";
      }

      return "¥" + Math.round(number);
    },
    trimDecimal(value, digits) {
      return Number(value).toFixed(digits).replace(/\.0$/, "");
    },
    niceStep(rawStep) {
      var step = Number(rawStep || 0);

      if (step <= 0) {
        return 1;
      }

      var magnitude = Math.pow(10, Math.floor(Math.log10(step)));
      var residual = step / magnitude;
      var nice = 1;

      if (residual <= 1) {
        nice = 1;
      } else if (residual <= 2) {
        nice = 2;
      } else if (residual <= 5) {
        nice = 5;
      } else {
        nice = 10;
      }

      return nice * magnitude;
    },
    tooltipTransform() {
      if (!this.activeBar) {
        return "translate(0,0)";
      }

      var tooltipWidth = 124;
      var tooltipHeight = 42;
      var x = this.activeBar.x + 10;
      var y = this.activeBar.y - tooltipHeight - 8;

      if (x + tooltipWidth > this.svgWidth - this.padding.right) {
        x = this.activeBar.x - tooltipWidth - 10;
      }

      if (x < this.padding.left) {
        x = this.padding.left;
      }

      if (y < this.padding.top) {
        y = this.activeBar.y + 12;
      }

      return "translate(" + x + "," + y + ")";
    }
  }
};
</script>

<style scoped>
.simple-bar-chart {
  width: 100%;
}

.simple-bar-chart__svg {
  width: 100%;
  display: block;
}

.simple-bar-chart__grid {
  stroke: rgba(229, 231, 235, 0.95);
  stroke-width: 1;
}

.simple-bar-chart__bar {
  cursor: pointer;
}

.simple-bar-chart__tooltip {
  pointer-events: none;
}

.simple-bar-chart__tooltip rect {
  fill: rgba(23, 23, 23, 0.9);
  filter: drop-shadow(0 8px 18px rgba(23, 23, 23, 0.18));
}

.simple-bar-chart__tooltip-label {
  fill: rgba(255, 255, 255, 0.72);
  font-size: 8px;
  font-weight: 600;
}

.simple-bar-chart__tooltip-value {
  fill: #ffffff;
  font-size: 9px;
  font-weight: 700;
}

.simple-bar-chart__x-label,
.simple-bar-chart__y-label {
  fill: #6b7280;
  font-size: 8px;
}

.simple-bar-chart__legend {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  justify-content: center;
  margin-top: 6px;
}

.simple-bar-chart__legend-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--text-subtle);
  font-size: 9px;
  font-weight: 600;
}

.simple-bar-chart__legend-item i {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  display: inline-block;
}
</style>
