<template>
  <div class="w-page">Step 2
    {{ value.type }}
  </div>
</template>

<script>
export default {
  name: "Step2",
  components: {},
  props: {
    modelValue: {},
    profile: {},
  },
  data() {
    return {
      value: {},
    };
  },
  watch: {
    modelValue(modelValue) {
      console.log("Step 2: modelValue: " , JSON.stringify(modelValue))
      this.value = modelValue;
    },
    value(value) {
      console.log("Step 2: value: " , JSON.stringify(value))
    },
  },
};
</script>