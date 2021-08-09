<template>
  <Dialog v-model:visible="dialogVisible" :modal="true">
    <template #header>
      <h3>Select Import File</h3>
    </template>
    <Toolbar>
      <template #left>
        <span
          class="p-button p-component p-fileupload-choose"
          tabindex="0"
          @click="choose"
          @keydown.enter="choose"
        >
          <input
            ref="fileInput"
            type="file"
            multiple=""
            accept=".act"
            @change="handleFileChange"
          />
          <span
            class="p-button-icon p-button-icon-left pi pi-fw pi-plus"
          ></span>
          <span class="p-button-label">Choose</span>
          <span class="p-ink"></span>
        </span>
        <Button
          label="Import"
          icon="pi pi-upload"
          class="p-ml-2"
          @click="importFile"
          :disabled="uploadDisabled"
        />
        <Button
          label="Cancel"
          icon="pi pi-times"
          class="p-ml-2"
          @click="cancel"
        />
      </template>
    </Toolbar>
    <Card>
      <template #title>
        {{ action.name }}
      </template>
      <template #content>
        {{ action.description }}
      </template>
    </Card>
  </Dialog>
</template>

<script>
export default {
  name: "Upload",
  components: {},
  props: {
    visible: Boolean,
    filetype: String,
    profileName: String,
    actionName: String,
  },
  emits: ["import", "cancel", "update:modelValue"],
  data() {
    return {
      dialogVisible: false,
      profileURL: String,
      uploadURL: String,
      file: File,
      hasFile: false,
      actionContent: "",
      action: { name: "", description: "" },
    };
  },
  mounted() {
    this.profileURL = this.$store.state.baseURL + "profiles";
    this.uploadURL =
      this.profileURL +
      "/" +
      this.profileName +
      "/actions/" +
      this.actionName +
      "/check";
    this.hasFile = false;
  },
  methods: {
    importFile() {
      console.log("import: " + this.file.name);
    },
    choose() {
      this.$refs.fileInput.click();
    },
    cancel() {
      this.$emit("cancel");
    },
    handleFileChange(e) {
      console.log("input: " + e.target.files[0]);
      this.file = e.target.files[0];
      this.hasFile = true;
      const reader = new FileReader();
      reader.onload = (res) => {
        let content = res.target.result;
        this.action = JSON.parse(content);
      };
      reader.onerror = (err) => console.log(err);
      reader.readAsText(this.file);

      //this.$emit('input', e.target.files[0])
    },
  },
  computed: {
    uploadDisabled() {
      return this.disabled || !this.hasFile;
    },
  },
  watch: {
    visible(visible) {
      this.dialogVisible = visible;
    },
  },
};
</script>