<template>
  <Dialog v-model:visible="dialogProfileVisible">
    <template #header>
      <h3>Add profile</h3>
    </template>
    <div class="p-fluid">
      <div class="p-field p-grid">
        <label for="name" class="p-col-24 p-mb-2 p-md-2 p-mb-md-0">name</label>
        <div class="p-col-24 p-md-10">
          <InputText
            id="name"
            type="text"
            v-model="addProfile.name"
            class="p-ml-2"
            :disabled="edit"
            :class="{ nameMissing: !isNameOK }"
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
            v-model="addProfile.description"
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
      <Button label="Save" icon="pi pi-check" autofocus @click="save" :disabled="!isNameOK" />
    </template>
  </Dialog>
</template>

<script>
export default {
  name: "AddProfile",
  components: {},
  props: {
    profile: {},
    visible: Boolean,
    edit: Boolean,
    profiles: {},
  },
  data() {
    return {
      dialogProfileVisible: false,
      addProfile: { name: "", description: "" },
      isNameOK: true,
    };
  },
  methods: {
    cancel() {
      this.$emit("cancel");
    },
    save() {
      this.$emit("save", this.addProfile);
    },
    checkName(name) {
        if (name == "") {
          this.isNameOK = false
          return
        }
        this.isNameOK = !this.profiles.map(elem => elem.toLowerCase()).includes(name.toLowerCase());
    }
  },
  watch: {
    visible(visible) {
      this.dialogProfileVisible = visible;
      this.checkName(this.addProfile.name)
    },
    profile(profile) {
        this.addProfile = profile;
    },
    addProfile: {
      deep: true,
      handler(profile) {
        this.checkName(profile.name)
      },
    },
  },
};
</script>

<style>
.nameMissing {
  background: lightcoral !important;
}
</style>