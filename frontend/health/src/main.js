import "@babel/polyfill";
import "mutationobserver-shim";
import Vue from "vue";
import "./plugins/bootstrap-vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
// eslint-disable-next-line
import styles from "./styles/source_sans_pro.css";
import Cookies from "vue-cookies";

Vue.config.productionTip = false;

Vue.use(Cookies);

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
