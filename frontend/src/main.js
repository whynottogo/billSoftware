import { createApp } from "vue";
import ElementPlus from "element-plus";
import "element-plus/dist/index.css";
import * as ElementPlusIconsVue from "@element-plus/icons-vue";

import App from "./App.vue";
import router from "./router";
import "./styles/global.css";

var app = createApp(App).use(router).use(ElementPlus);
for (var key in ElementPlusIconsVue) {
  app.component(key, ElementPlusIconsVue[key]);
}
app.mount("#app");

