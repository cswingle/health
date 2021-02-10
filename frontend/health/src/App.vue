<template>
  <b-container id="app">
    <b-navbar>
      <b-navbar-nav>
        <b-button
          to="/"
          class="m-2"
          variant="outline-primary"
          exact
          active-class="active btn-primary"
        >
          Home
        </b-button>
        <b-button
          to="/login"
          class="m-2"
          variant="outline-primary"
          active-class="active btn-primary"
          v-if="!isLoggedIn"
        >
          Login
        </b-button>
        <b-button
          @click="logout"
          class="m-2"
          variant="outline-primary"
          active-class="active btn-primary"
          v-if="isLoggedIn"
        >
          Logout
        </b-button>
      </b-navbar-nav>
    </b-navbar>
    <router-view />
    <b-toast
      id="update-ready"
      variant="info"
      toaster="b-toaster-bottom-center"
      title="App Update Available"
      solid
      no-auto-hide
    >
      <template v-slot:default>
        <div class="text-center">
          <b-button variant="primary" @click="updateApp">Update</b-button>
        </div>
      </template>
    </b-toast>
  </b-container>
</template>

<script>
export default {
  name: "App",
  async created() {
    await this.$store.dispatch("checkAuthentication");

    document.addEventListener("swUpdated", this.updateAvailable, {
      once: true
    });
    navigator.serviceWorker.addEventListener("controllerchange", () => {
      if (this.refreshing) return;
      this.refreshing = true;
      window.location.reload();
    });
  },
  data() {
    return {
      cookie: "",
      registration: null
    };
  },
  computed: {
    isLoggedIn() {
      return this.$store.getters.isLoggedIn;
    }
  },
  methods: {
    logout() {
      this.$store.dispatch("logout").then(() => {
        this.$router.push("/login");
      });
    },
    async getCookie() {
      this.cookie = await this.$cookies.get("auth");
      console.log(this.cookie);
    },
    updateAvailable(e) {
      this.registration = e.detail;
      this.$bvToast.show("update-ready");
    },
    updateApp() {
      if (!this.registration || !this.registration.waiting) {
        return;
      }
      this.registration.waiting.postMessage("skipWaiting");
    }
  }
};
</script>

<style lang="scss">
#app {
  font-family: source_sans_pro, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
</style>
