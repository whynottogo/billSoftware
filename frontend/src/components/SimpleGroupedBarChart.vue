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
          rx="8"
          ry="8"
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
    groupWidth() {
      return this.plotWidth / this.labels.length;
    },
    barWidth() {
      return Math.min(24, this.groupWidth / (this.series.length + 1.5));
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
      var totalWidth = this.barWidth * this.series.length + 8 * (this.series.length - 1);
      var start = this.groupCenter(index) - totalWidth / 2;

      return this.series.map(function(item, seriesIndex) {
        var value = item.values[index] || 0;
        var x = start + seriesIndex * (self.barWidth + 8);
        var y = self.yFor(value);
        var height = self.padding.top + self.plotHeight - y;

        return {
          key: item.name + "-" + index,
          x: x,
          y: y,
          width: self.barWidth,
          height: height,
          color: item.color
        };
      });
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

.simple-bar-chart__x-label,
.simple-bar-chart__y-label {
  fill: #6b7280;
  font-size: 12px;
}

.simple-bar-chart__legend {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: center;
  margin-top: 8px;
}

.simple-bar-chart__legend-item {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--text-subtle);
  font-size: 13px;
  font-weight: 600;
}

.simple-bar-chart__legend-item i {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  display: inline-block;
}
</style>
