<template>
  <Dialog v-model:visible="dialogVisible">
    <template #header>
      <h3>
        <div class="p-orderlist-header" v-if="$slots.sourceHeader">
          <slot name="sourceHeader"></slot>
        </div>
      </h3>
    </template>
    <div v-for="(param, x) in iconlist" :key="x">
      {{ param }}
      <div class="icon-item">
        <img :src="'assets/' + param" />
        <div>{{ param }}</div>
      </div>
    </div>
    <template #footer>
      <Button
        label="Cancel"
        icon="pi pi-times"
        class="p-button-text"
        @click="cancel"
      />
      <Button label="Save" icon="pi pi-check" autofocus @click="save" />
    </template>
  </Dialog>
</template>

<script>
export default {
  name: "SelectIcon",
  components: {},
  props: {
    modelValue: {
      type: String,
      default: "",
    },
    iconlist: {
      type: Array,
    },
    visible: Boolean,
  },
  emits: ["save", "cancel", "update:modelValue"],
  data() {
    return {
      dialogVisible: false,
      icon: "",
    };
  },
  methods: {
    cancel() {
      this.$emit("cancel");
    },
    save() {
      this.$emit("update:modelValue", this.icon);
    },
  },
  watch: {
    visible(visible) {
      this.dialogVisible = visible;
    },
    modelValue(value) {
      this.icon = value;
    },
  },
};
</script>
