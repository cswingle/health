<template>
  <b-container>
    <b-toast
      title="Data Uploaded"
      variant="success"
      auto-hide-delay="2000"
      v-model="submitSuccess"
    >
      Uploaded {{ rowsInserted }} <span v-if="rowsInserted > 1">rows</span><span v-else>row</span>
    </b-toast>
    <b-toast
      title="Submit Error"
      variant="error"
      auto-hide-delay="2000"
      v-model="submitFail"
    >
      Failed to insert data
    </b-toast>
    <h1>Data Entry</h1>
    <hr />
    <b-container v-if="isLoggedIn">
      <b-row>
        <b-col>
          <Health
            v-for="h in healthEntries"
            v-bind:key="h.variable"
            v-bind:healthEntry="h"
            v-bind:refVariables="refVariables"
          />
        </b-col>
      </b-row>
      <b-row>
        <b-col cols="10" />
        <b-col cols="2">
          <b-button size="sm" variant="primary" v-on:click="addRow">
            add entry
          </b-button>
        </b-col>
      </b-row>
      <hr />
      <b-row>
        <b-col />
        <b-col cols="2">
          <b-button size="sm" variant="primary" v-on:click="submit">
            submit
          </b-button>
        </b-col>
        <b-col />
      </b-row>
    </b-container>
    <b-container v-else>
      <b-row>
        <b-col>
          <b-alert show
            ><b-link to="/login">Login</b-link> to enter data</b-alert
          >
        </b-col>
      </b-row>
    </b-container>
  </b-container>
</template>

<script>
import axios from "axios";
import Health from "./Health.vue";

export default {
  name: "Data",
  components: {
    Health
  },
  async mounted() {
    await axios
      .get(this.baseUrl + "/ref_variables")
      .then(
        response =>
          (this.refVariables = response.data.sort(
            (a, b) => a.sequence - b.sequence
          ))
      );
  },
  data() {
    return {
      // baseUrl: "https://example.com/health/v1",
      baseUrl: "http://localhost:LISTEN_PORT/health/v1",
      refVariables: null,
      submitSuccess: false,
      submitFail: false,
      rowsInserted: 0,
      healthEntries: [{ variable: null, value: null }]
    };
  },
  computed: {
    isLoggedIn: function() {
      return this.$store.getters.isLoggedIn;
    }
  },
  methods: {
    addRow: function() {
      let newHealth = { variable: null, value: null };
      this.healthEntries.push(newHealth);
    },
    submit: function() {
      let ts = new Date().toISOString();
      let payload = [];
      for (let h of this.healthEntries) {
        if ((h.variable) && (h.variable.length > 0) && (h.value)) {
          let entry = {
            ts: ts,
            username: null,
            variable: h.variable,
            value: h.value
          };
          payload.push(entry);
        }
      }
      if (payload.length) {
        axios
          .post(this.baseUrl + "/health/keys/u", payload)
          .then(response => {
            // TODO: use this for a success toast
            if (response.status === 200) {
              let rowsInserted = 0;
              for (let i in response.data.message) {
                if (
                  response.data.status[i] != "error" &&
                  response.data.message[i].search("inserted") > -1
                ) {
                  rowsInserted = rowsInserted + 1;
                }
              }
              if (rowsInserted) {
                this.rowsInserted = rowsInserted;
                this.submitSuccess = true;
              }
            } else {
              // TODO: failure toast
              this.submitFail = true;
            }
          })
          .catch(err => console.log("error posting", err));
      }
      this.healthEntries = [{ variable: null, value: null }];
      this.rowsInserted = 0;
      this.submitSuccess = false;
      this.submitFail = false;
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
