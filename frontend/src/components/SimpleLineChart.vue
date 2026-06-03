<template>
  <div class="simple-chart">
    <svg :viewBox="'0 0 ' + svgWidth + ' ' + height" class="simple-chart__svg" preserveAspectRatio="none">
      <g v-for="tick in ticks" :key="'tick-' + tick.value">
        <line
          :x1="padding.left"
          :x2="svgWidth - padding.right"
          :y1="yFor(tick.value)"
          :y2="yFor(tick.value)"
          class="simple-chart__grid"
        />
        <text
          :x="padding.left - 10"
          :y="yFor(tick.value) + 4"
          class="simple-chart__y-label"
          text-anchor="end"
        >
          {{ tick.label }}
        </text>
      </g>

      <path
        v-for="item in series"
        :key="item.name"
        :d="linePath(item.values)"
        :stroke="item.color"
        class="simple-chart__line"
      />

      <circle
        v-for="point in points"
        :key="point.key"
        :cx="xFor(point.index)"
        :cy="yFor(point.value)"
        :fill="point.color"
        class="simple-chart__point"
        r="3"
        @mouseenter="setActivePoint(point)"
        @mousemove="setActivePoint(point)"
        @mouseleave="clearActivePoint"
      >
      </circle>

      <g v-if="activePoint" :transform="tooltipTransform" class="simple-chart__tooltip">
        <rect width="118" height="42" rx="10" ry="10" />
        <text x="12" y="17" class="simple-chart__tooltip-label">{{ activePoint.label }}</text>
        <text x="12" y="32" class="simple-chart__tooltip-value">{{ activePoint.name }} · {{ formatCurrency(activePoint.value) }}</text>
      </g>

      <text
        v-for="(label, index) in labels"
        :key="'x-' + label + '-' + index"
        :x="xFor(index)"
        :y="height - 12"
        class="simple-chart__x-label"
        text-anchor="middle"
      >
        {{ label }}
      </text>
    </svg>

    <div v-if="series.length > 1" class="simple-chart__legend">
      <span v-for="item in series" :key="item.name" class="simple-chart__legend-item">
        <i :style="{ backgroundColor: item.color }"></i>
        {{ item.name }}
      </span>
    </div>
  </div>
</template>

<script>
export default {
  name: "SimpleLineChart",
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
      default: 280
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
        bottom: 40,
        left: 52
      },
      activePoint: null
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
    points() {
      var list = [];
      var labels = this.labels;

      this.series.forEach(function(item) {
        item.values.forEach(function(value, index) {
          list.push({
            key: item.name + "-" + index,
            value: value,
            index: index,
            color: item.color,
            name: item.name,
            label: labels[index] || ""
          });
        });
      });

      return list;
    },
    tooltipTransform() {
      if (!this.activePoint) {
        return "translate(0,0)";
      }

      var tooltipWidth = 118;
      var tooltipHeight = 42;
      var x = this.activePoint.x + 10;
      var y = this.activePoint.y - tooltipHeight - 8;

      if (x + tooltipWidth > this.svgWidth - this.padding.right) {
        x = this.activePoint.x - tooltipWidth - 10;
      }

      if (x < this.padding.left) {
        x = this.padding.left;
      }

      if (y < this.padding.top) {
        y = this.activePoint.y + 12;
      }

      return "translate(" + x + "," + y + ")";
    }
  },
  methods: {
    xFor(index) {
      if (this.labels.length <= 1) {
        return this.padding.left + this.plotWidth / 2;
      }

      return this.padding.left + (this.plotWidth / (this.labels.length - 1)) * index;
    },
    yFor(value) {
      return this.padding.top + this.plotHeight - (Number(value || 0) / this.maxValue) * this.plotHeight;
    },
    linePath(values) {
      var self = this;

      return values
        .map(function(value, index) {
          return (index === 0 ? "M" : "L") + self.xFor(index) + "," + self.yFor(value);
        })
        .join(" ");
    },
    setActivePoint(point) {
      this.activePoint = {
        name: point.name,
        label: point.label || "-",
        value: point.value,
        color: point.color,
        x: this.xFor(point.index),
        y: this.yFor(point.value)
      };
    },
    clearActivePoint() {
      this.activePoint = null;
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
    }
  }
};
</script>

<style scoped>
.simple-chart {
  width: 100%;
}

.simple-chart__svg {
  width: 100%;
  display: block;
}

.simple-chart__grid {
  stroke: rgba(229, 231, 235, 0.95);
  stroke-width: 1;
}

.simple-chart__line {
  fill: none;
  stroke-width: 2.4;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.simple-chart__point {
  cursor: pointer;
  stroke: transparent;
  stroke-width: 10;
}

.simple-chart__tooltip {
  pointer-events: none;
}

.simple-chart__tooltip rect {
  fill: rgba(23, 23, 23, 0.9);
  filter: drop-shadow(0 8px 18px rgba(23, 23, 23, 0.18));
}

.simple-chart__tooltip-label {
  fill: rgba(255, 255, 255, 0.72);
  font-size: 8px;
  font-weight: 600;
}

.simple-chart__tooltip-value {
  fill: #ffffff;
  font-size: 9px;
  font-weight: 700;
}

.simple-chart__x-label,
.simple-chart__y-label {
  fill: #6b7280;
  font-size: 8px;
}

.simple-chart__legend {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  justify-content: center;
  margin-top: 6px;
}

.simple-chart__legend-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--text-subtle);
  font-size: 9px;
  font-weight: 600;
}

.simple-chart__legend-item i {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  display: inline-block;
}
</style>
