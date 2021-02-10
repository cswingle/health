<template>
  <div>
    <b-container>
      <b-form-row>
        <b-col>
          <b-form-group label="Variable" label-cols="3">
            <b-form-select
              v-model="healthEntry.variable"
              :options="refVariables"
              value-field="variable"
              text-field="variable"
            />
          </b-form-group>
        </b-col>
      </b-form-row>
      <b-form-row>
        <b-col>
          <b-form-group label="Value" label-cols="3">
            <b-form-input v-model.number="healthEntry.value" type="number" />
          </b-form-group>
        </b-col>
      </b-form-row>
    </b-container>
  </div>
</template>

<script>
export default {
  name: "Health",
  props: {
    healthEntry: Object,
    refVariables: Array
  },
  data: function() {
    return {};
  },
  methods: {}
};
</script>

<style scoped></style>
