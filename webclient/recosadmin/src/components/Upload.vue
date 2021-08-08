<template>
  <Dialog v-model:visible="dialogVisible" :modal="true">
    <FileUpload
      name="demo[]"
      url="./upload.php"
      @upload="onUpload"
      :multiple="true"
      accept="image/*"
      :maxFileSize="1000000"
    >
      <template #empty>
        <p>Drag and drop files to here to upload.</p>
      </template>
    </FileUpload>
  </Dialog>
</template>

<script>
export default {
  name: "Upload",
  components: {},
  props: {
    visible: Boolean,
  },
  data() {
    return {
      dialogVisible: false,
    };
  },
  methods: {},
  watch: {
    visible(visible) {
      this.dialogVisible = visible;
    },
  },
};
</script>