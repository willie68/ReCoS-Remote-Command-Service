<template>
  <div class="w-page">
    What kind of action do you like to add?<br />
    <Listbox
      v-model="commandType"
      :options="commandTypes"
      placeholder="select a type"
      optionLabel="description"
      optionValue="type"
      :filter="true"
      dataKey="type"
      class="p-mt-2"
    />
  </div>
</template>

<script>
export default {
  name: "Step1",
  components: {},
  props: {
    value: {},
    profile: {},
    commandTypes: Array,
  },
  emits: ["next", "value"],
  data() {
    return {
      type: null,
      localValue: {},
    };
  },
  computed: {
    commandType: {
      get: function () {
        return this.localValue.type;
      },
      set: function (newType) {
        if (newType) {
          this.localValue.type = newType;
          console.log("Step1: commandTypes", JSON.stringify(this.commandTypes));
          this.commandTypes.forEach((element) => {
            if (element.type == newType) {
              this.localValue.actiontype = element.wizardactiontype;
            }
          });
        } else {
          this.localValue = { type: "", actiontype: "SINGLE" };
        }
        this.check();
      },
    },
  },
  methods: {
    check() {
      console.log("Step1: select command:", JSON.stringify(this.localValue));
      this.$emit("value", this.localValue);
      if (this.localValue.type) {
        this.$emit("next", true);
        return;
      }
      this.$emit("next", false);
    },
  },
  mounted() {
    this.localValue = this.value;
    console.log(
      "Step1: mounted update value:",
      JSON.stringify(this.localValue)
    );
    this.check();
  },
  beforeUpdate() {
    this.localValue = this.value;
    console.log("Step1: before update value:", JSON.stringify(this.localValue));
  },
  updated() {
    console.log("Step1: updated value:", JSON.stringify(this.localValue));
  },
  watch: {
    value: {
      deep: true,
      handler(value) {
        this.localValue = value;
        console.log("Step1: value:", JSON.stringify(this.localValue));
      },
    },
  },
};
</script>