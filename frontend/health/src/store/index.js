import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";
import qs from "qs";

Vue.use(Vuex);

axios.defaults.headers.post["Content-Type"] =
  "application/x-www-form-urlencoded";
// axios.defaults.baseURL = "https://example.com/health/v1";
axios.defaults.baseURL = "http://localhost:LISTEN_PORT/health/v1";
axios.defaults.withCredentials = true;

export default new Vuex.Store({
  state: {
    status: "",
    loggedIn: false
  },
  mutations: {
    auth_request(state) {
      state.status = "loading";
      state.loggedIn = false;
    },
    auth_success(state) {
      state.status = "success";
      state.loggedIn = true;
    },
    auth_error(state) {
      state.status = "error";
      state.loggedIn = false;
      Vue.$cookies.remove("auth");
    },
    logout(state) {
      state.status = "";
      state.loggedIn = false;
      Vue.$cookies.remove("auth");
    }
  },
  actions: {
    async login({ commit }, user) {
      try {
        commit("auth_request");
        let response = await axios.post("/login", qs.stringify(user));
        if (response.status === 200) {
          commit("auth_success");
        } else {
          commit("auth_error");
        }
      } catch (err) {
        console.log(`err: ${err.message}`);
        commit("auth_error");
      }
    },
    async logout({ commit }) {
      let response = await axios.post("/logout").catch(err => {
        console.log("logout error " + err);
        commit("logout");
      });
      if (response.status == 200) {
        commit("logout");
      } else {
        console.log("logout error!");
        commit("logout");
      }
    },
    async checkAuthentication({ commit }) {
      if (Vue.$cookies.get("auth")) {
        commit("auth_success");
      }
    }
  },
  getters: {
    authStatus: state => state.status,
    isLoggedIn: state => state.loggedIn
  }
});
