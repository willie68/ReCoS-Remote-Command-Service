<template>
  <Dialog v-model:visible="dialogVisible" :modal="true">
    <template #header>
      <h3>Select Action</h3>
    </template>
    <Listbox
      v-model="selectedAction"
      :options="sourceValue"
      optionLabel="name"
      :filter="true"
      listStyle="max-height:150px"
    />
    <template #footer>
      <Button
        label="Wizard"
        icon="pi pi-flag"
        class="p-button-text"
        @click="wizard"
      />
      <Button
        label="Remove"
        icon="pi pi-trash"
        class="p-button-text"
        @click="remove"
      />
      <Button
        label="Cancel"
        icon="pi pi-times"
        class="p-button-text"
        @click="cancel"
      />
      <Button
        label="Save"
        icon="pi pi-check"
        autofocus
        @click="save"
        :disabled="!isNameOK"
      />
    </template>
  </Dialog>
</template>

<script>
export default {
  name: "SelectAction",
  components: {},
  props: {
    sourceValue: {
      type: Array,
      default: null,
    },
    visible: Boolean,
    selectByName: String,
  },
  emits: ["cancel", "save" , "remove"],
  data() {
    return {
      selectedAction: {},
      dialogVisible: false,
      isNameOK: true,
    };
  },
  methods: {
    wizard() {
      this.$emit("wizard");
    },
    cancel() {
      this.$emit("cancel");
    },
    save() {
      this.$emit("save", this.selectedAction);
    },
    remove() {
      this.$emit("remove");
    },
    checkName(name) {
      if (name == "") {
        this.isNameOK = false;
        return;
      }
      this.isNameOK = !this.profiles
        .map((elem) => elem.toLowerCase())
        .includes(name.toLowerCase());
    },
  },
  beforeUpdate() {
    console.log("SelectAction: BeforeUpdate");
    this.selectedAction = null;
    if (this.selectByName) {
      if (this.sourceValue) {
        this.sourceValue.forEach((element) => {
          if (element.name == this.selectByName) {
            console.log("SelectAction: Element found");
            this.selectedAction = element;
          }
        });
      }
    }
  },
  watch: {
    visible(visible) {
      this.dialogVisible = visible;
      if (visible && this.selectByName) {
        if (this.sourceValue) {
          this.sourceValue.forEach((element) => {
            if (element.name == this.selectByName) {
              this.selectedAction = element;
            }
          });
        }
      }
    },
  },
};
</script>
