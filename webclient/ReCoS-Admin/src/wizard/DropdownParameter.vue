<template>
  <Dropdown
    :options="options"
    :v-tooltip="tooltip"
    v-model="value"
    :placeholder="placeholder"
    :scrollheight="scrollHeight"
  />
</template>

<script>
export default {
  name: "DropdownParameter",
  components: {},
  props: {
    modelValue: String,
    options: Array,
    tooltip: String,
    placeholder: String,
    scrollheight: String,
  },
  emits: ["update:modelValue"],
  data() {
    return {
      localValue: {},
    };
  },
  computed: {
    value: {
      get: function () {
        if (!this.localValue) {
          return this.modelValue;
        }
        return this.localValue;
      },
      set: function (newValue) {
        console.log("DropdownParameter: new value", newValue);
        this.localValue = newValue;
        this.$emit("update:modelValue", newValue);
      },
    },
  },
  methods: {},
  mounted() {},
  watch: {
    value: {
      deep: true,
      handler(value) {
        this.localValue = value;
        console.log(
          "DorpdownParameter: watch value: ",
          JSON.stringify(this.localValue)
        );
      },
    },
  },
};
</script>