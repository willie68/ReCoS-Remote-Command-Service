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
            :accept="filetype"
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
  },
  emits: ["save", "cancel"],
  data() {
    return {
      dialogVisible: false,
      profileURL: String,
      uploadURL: String,
      file: File,
      hasFile: false,
      action: { name: "", description: "" },
    };
  },
  mounted() {
    this.hasFile = false;
  },
  methods: {
    importFile() {
      this.$emit("save", this.action)
    },
    choose() {
      this.$refs.fileInput.click();
    },
    cancel() {
      this.hasFile = false
      this.file = null
      this.action = { name: "", description: "" }
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