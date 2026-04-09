<template>
  <div ref="canvas" class="finance-chart" :style="{ height: height }"></div>
</template>

<script>
import * as echarts from "echarts/dist/echarts.js";

export default {
  name: "FinanceChart",
  props: {
    option: {
      type: Object,
      required: true
    },
    height: {
      type: String,
      default: "280px"
    }
  },
  data() {
    return {
      chart: null,
      renderFrame: 0
    };
  },
  watch: {
    option() {
      this.renderChart();
    }
  },
  mounted() {
    this.renderChart();
    window.addEventListener("resize", this.handleResize);
  },
  beforeUnmount() {
    window.removeEventListener("resize", this.handleResize);

    if (this.renderFrame) {
      window.cancelAnimationFrame(this.renderFrame);
      this.renderFrame = 0;
    }

    if (this.chart) {
      this.chart.dispose();
      this.chart = null;
    }
  },
  methods: {
    renderChart() {
      if (this.renderFrame) {
        window.cancelAnimationFrame(this.renderFrame);
      }

      this.renderFrame = window.requestAnimationFrame(() => {
        var canvas = this.$refs.canvas;
        var rect;

        this.renderFrame = 0;

        if (!canvas) {
          return;
        }

        rect = canvas.getBoundingClientRect();

        if (!rect.width || !rect.height) {
          this.renderChart();
          return;
        }

        if (this.chart) {
          this.chart.dispose();
          this.chart = null;
        }

        this.chart = echarts.init(canvas);
        this.chart.setOption(this.option, true);
        this.chart.resize();
      });
    },
    handleResize() {
      if (this.chart) {
        this.chart.resize();
      }
    }
  }
};
</script>

<style scoped>
.finance-chart {
  width: 100%;
}
</style>
