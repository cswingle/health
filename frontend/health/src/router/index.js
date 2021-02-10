import Vue from "vue";
import VueRouter from "vue-router";
import Data from "@/components/Data.vue";
import Login from "@/components/Login.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "Data",
    component: Data
  },
  {
    path: "/login",
    name: "login",
    component: Login
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

export default router;
