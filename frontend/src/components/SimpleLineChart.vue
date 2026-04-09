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
        r="5"
      />

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
      }
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

      return Math.ceil(max / 4 / 1000) * 1000 * 4;
    },
    plotWidth() {
      return this.svgWidth - this.padding.left - this.padding.right;
    },
    plotHeight() {
      return this.height - this.padding.top - this.padding.bottom;
    },
    ticks() {
      var step = this.maxValue / 4;
      var list = [];
      var index;

      for (index = 0; index <= 4; index += 1) {
        list.push({
          value: step * index,
          label: "¥" + Math.round((step * index) / 1000) + "k"
        });
      }

      return list;
    },
    points() {
      var list = [];

      this.series.forEach(function(item) {
        item.values.forEach(function(value, index) {
          list.push({
            key: item.name + "-" + index,
            value: value,
            index: index,
            color: item.color
          });
        });
      });

      return list;
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
  stroke-width: 4;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.simple-chart__x-label,
.simple-chart__y-label {
  fill: #6b7280;
  font-size: 12px;
}

.simple-chart__legend {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: center;
  margin-top: 8px;
}

.simple-chart__legend-item {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--text-subtle);
  font-size: 13px;
  font-weight: 600;
}

.simple-chart__legend-item i {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  display: inline-block;
}
</style>
