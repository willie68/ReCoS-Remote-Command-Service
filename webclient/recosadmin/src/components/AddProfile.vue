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
            :class="{ 'p-invalid': !isNameOK }"
            autofocus
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
        label="Import"
        icon="pi pi-cloud-upload"
        class="p-button-text"
        @click="this.importAction();"
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
  <Upload :visible="dialogUploadVisible" filetype=".profile" @cancel="dialogUploadVisible = false" @save="doImport"/>
</template>

<script>
import Upload from "./Upload.vue";
export default {
  name: "AddProfile",
  components: {
    Upload,
  },
  props: {
    profile: {},
    visible: Boolean,
    edit: Boolean,
    profiles: {},
  },
  data() {
    return {
      dialogProfileVisible: false,
      dialogUploadVisible: false,
      addProfile: { name: "", description: "" },
      isNameOK: true,
    };
  },
  methods: {
    importAction() {
      this.dialogUploadVisible = true;
    },
    doImport(event) {
      let newProfile = event;
      console.log("new profile: " + JSON.stringify(newProfile))
      this.dialogUploadVisible = false;
      this.addProfile = newProfile;
      //this.emitter.emit("insertAction", newAction);
    },
    cancel() {
      this.$emit("cancel");
    },
    save() {
      this.$emit("save", this.addProfile);
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
  watch: {
    visible(visible) {
      this.dialogProfileVisible = visible;
      this.checkName(this.addProfile.name);
    },
    profile(profile) {
      this.addProfile = profile;
    },
    addProfile: {
      deep: true,
      handler(profile) {
        this.checkName(profile.name);
      },
    },
  },
};
</script>
