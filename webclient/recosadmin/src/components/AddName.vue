<template>
  <Dialog v-model:visible="dialogVisible">
    <template #header>
      <h3><div class="p-orderlist-header" v-if="$slots.sourceHeader">
        <slot name="sourceHeader"></slot>
      </div></h3>
    </template>
    <div class="p-fluid">
      <div class="p-field p-grid">
        <label for="name" class="p-col-24 p-mb-2 p-md-2 p-mb-md-0">name</label>
        <div class="p-col-24 p-md-10">
          <InputText
            id="name"
            type="text"
            v-model="name"
            class="p-ml-2"
            :class="{ 'p-invalid': !isNameOK }"
            @keyup.enter="save"
          />
        </div>
      </div>
    </div>
    <template #footer>
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
  name: "AddName",
  components: {},
  props: {
    modelValue: {
      type: String,
      default: "",
    },
    excludeList: {
      type: Array,
      default: null,
    },
    visible: Boolean,
  },
  data() {
    return {
      dialogVisible: false,
      name: "", 
      isNameOK: true,
    };
  },
  methods: {
    cancel() {
      this.$emit("cancel");
    },
    save() {
      this.$emit("updatemodelValue", this.name);
      this.$emit("save", this.name);
    },
    checkName(name) {
      if (name == "") {
        this.isNameOK = false;
        return;
      }
      this.isNameOK = !this.excludeList
        .map((elem) => elem.toLowerCase())
        .includes(name.toLowerCase());
    },
  },
  watch: {
    visible(visible) {
      this.dialogVisible = visible;
      this.checkName(this.name);
    },
    modelValue(value) {
      this.name = value;
    },
    name(name) {
      this.checkName(name)
    }
  },
};
</script>
