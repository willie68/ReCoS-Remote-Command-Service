<template>
  <Dialog v-model:visible="dialogProfileVisible">
    <template #header>
      <h3>Edit profile</h3>
    </template>
    <div class="p-fluid">
      <div class="p-field p-grid">
        <label for="name" class="p-col-24 p-mb-2 p-md-2 p-mb-md-0">name</label>
        <div class="p-col-24 p-md-10">
          <InputText
            id="name"
            type="text"
            v-model="editProfile.name"
            class="p-ml-2"
          />
        </div>
      </div>
      <div class="p-field p-grid">
        <label for="description" class="p-col-24 p-mb-2 p-md-2 p-mb-md-0"
          >description</label
        >
        <div class="p-col-24 p-md-10">
          <InputText
            id="description"
            type="text"
            v-model="editProfile.description"
            class="p-ml-2"
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
      <Button label="Save" icon="pi pi-check" autofocus @click="save" />
    </template>
  </Dialog>
</template>

<script>
export default {
  name: "EditProfile",
  components: {},
  props: {
    profile: String,
    visible: Boolean,
  },
  data() {
    return {
      dialogProfileVisible: false,
      editProfile: { name: "", description: "" },
    };
  },
  methods: {
    cancel() {
      this.$emit("cancel");
    },
    save() {
      this.$emit("save", this.editProfile);
    },
  },
  watch: {
    visible(visible) {
      console.log("EditProfile visible changed");
      this.dialogProfileVisible = visible;
    },
  },
};
</script>