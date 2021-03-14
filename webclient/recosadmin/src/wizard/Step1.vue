<template>
  <div class="w-page">
    What kind of action do you like to add?<br />
    <Listbox
      v-model="commandType"
      :options="commandtypes"
      placeholder="select a type"
      optionLabel="description"
      optionValue="type"
      editable
    />
  </div>
</template>

<script>
export default {
  name: "Step1",
  components: {},
  props: {
    modelValue: {},
    profile: {},
  },
  data() {
    return {
      commandtypes: [],
      commandType: null,
      type: null,
      value: {},
    };
  },
  methods: {
    check(commandType) {
      console.log("Step1: select command:", JSON.stringify(commandType));
      if (commandType) {
        this.value.type = commandType;
      } else {
        this.value = { type: "" };
      }
      console.log("Step1: emit:", JSON.stringify(this.value));
      this.$emit("update:value", this.value);
      if (commandType) {
        this.$emit("next", true);
        return;
      }
      this.$emit("next", false);
    },
  },
  updated() {},
  mounted() {
    let url = this.$store.state.baseURL + "/config/commands";
    fetch(url)
      .then((res) => res.json())
      .then((data) => {
        this.commandtypes = [];
        data.forEach((element) => {
          if (element.wizard && element.wizard == true) {
            this.commandtypes.push(element);
          }
        });
      })
      .catch((err) => console.log(err.message));
  },
  watch: {
    commandType(commandType) {
      this.check(commandType);
    },
    modelValue(modelValue) {
      this.value = modelValue;
    },
  },
};
</script>