<template>
  <div>
    <b-form class="login" @submit.prevent="login">
      <h1>Sign in</h1>
      <b-form-row>
        <b-col>
          <b-form-group label="Username" label-for="name" label-cols="3">
            <b-form-input
              id="name"
              required
              v-model="username"
              placeholder="Username"
              autocapitalize="off"
            />
          </b-form-group>
        </b-col>
      </b-form-row>
      <b-form-row>
        <b-col>
          <b-form-group label="Password" label-for="password" label-cols="3">
            <b-form-input
              id="password"
              type="password"
              required
              v-model="password"
              placeholder="Password"
              cols="1"
            />
          </b-form-group>
        </b-col>
      </b-form-row>
      <hr />
      <b-button variant="success" type="submit">Login</b-button>
    </b-form>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: "",
      password: ""
    };
  },

  methods: {
    login() {
      let username = this.username;
      let password = this.password;
      this.$store
        .dispatch("login", {
          username: username,
          password: password
        })
        .then(() => this.$router.push("/"))
        .catch(err => console.log(err));
    }
  }
};
</script>
