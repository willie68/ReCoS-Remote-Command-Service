<template>
  <div class="layout-content">
    <div class="content-section implementation">
      <Toolbar>
        <template #left>
          <Button label="New" icon="pi pi-plus" class="p-mr-2"></Button>
          <Button label="Upload" icon="pi pi-upload" class="p-button-success" />
          <i class="pi pi-bars p-toolbar-separator p-mr-2" />
          <i class="pi pi-cog" />
          <SplitButton
            label="Save"
            icon="pi pi-check"
            :model="items"
            class="p-button-warning"
          ></SplitButton>
        </template>

        <template #right>
          <span class="p-input-icon-right">
            <InputText v-model="pwd" type="password"></InputText>
            <i class="pi pi-spin pi-spinner" />
          </span>
          <Button icon="pi pi-cog" class="p-mr-1" />
        </template>
      </Toolbar>
    </div>
  </div>
</template>

<script>
export default {
  name: "Page",
  props: {
    msg: String,
  },
  data() {
    return {
      items: [
        {
          label: "Update",
          icon: "pi pi-refresh",
        },
        {
          label: "Delete",
          icon: "pi pi-times",
        },
        {
          label: "Vue Website",
          icon: "pi pi-external-link",
          command: () => {
            window.location.href = "https://vuejs.org/";
          },
        },
        {
          label: "Upload",
          icon: "pi pi-upload",
          command: () => {
            window.location.hash = "/fileupload";
          },
        },
      ],
    };
  },
};
</script>

<style>
.p-button {
  margin-bottom: 0.5rem;
}
</style>