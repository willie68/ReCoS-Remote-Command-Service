<template>
  <Panel :header="activeAction.name" class="action-panel-custom">
    <template #icons>
      <button
        class="p-panel-header-icon p-link p-mr-2 p-mt-0 p-mb-0 p-pt-0 p-pb-0"
        @click="toggle"
      >
        <span class="pi pi-plus"></span>
      </button>
      <button
        class="p-panel-header-icon p-link p-mr-2 p-mt-0 p-mb-0 p-pt-0 p-pb-0"
        @click="toggle"
      >
        <span class="pi pi-trash"></span>
      </button>
    </template>

    <div class="p-fluid p-formgrid p-grid">
      <div class="p-field p-col">
        <label for="name">Name</label>
        <InputText id="name" type="text" v-model="activeAction.name" />
      </div>
      <div class="p-field p-col">
        <label for="rows">Type</label>
        <Dropdown
          v-model="activeAction.type"
          :options="enumActionTypes"
          placeholder="select a type"
          optionLabel="name"
          optionValue="type"
        />
      </div>
      <div class="p-field p-col">
        <label for="title">Title</label>
        <InputText id="title" type="text" v-model="activeAction.title" />
      </div>
      <div class="p-field p-col">
        <label for="icon">Icon</label>
        <InputText id="icon" type="text" v-model="activeAction.icon" />
      </div>
      <div class="p-field p-col">
        <label for="fontsize">Font size</label>
        <Dropdown
          id="fontsize"
          v-model="activeAction.fontsize"
          :options="enumFontSizes"
          placeholder="select a size"
          optionLabel="name"
          optionValue="value"
        />
      </div>
      <div class="p-field p-col">
        <label for="fontcolor">Font color</label>
        <ColorPicker
          v-model="activeAction.fontcolor"
          :inline="false"
          defaultColor="#FFFFFF"
        />
      </div>
      <div class="p-field p-col">
        <label for="outlined">Outlined</label><br />
        <Checkbox
          id="outlined"
          v-model="activeAction.outlined"
          :binary="true"
        />
      </div>
      <div class="p-field p-col">
        <label for="runone">Only run one</label><br />
        <Checkbox id="runone" v-model="activeAction.runone" :binary="true" />
      </div>
    </div>
    <div class="p-fluid p-formgrid p-grid">
      <div class="p-field p-col">
        <label for="description">Description</label>
        <InputText id="description" type="text" v-model="activeAction.description" />
      </div>
    </div>
  </Panel>
</template>

<script>
export default {
  name: "Action",
  props: {
    action: {},
  },
  data() {
    return {
      activeAction: {},
      enumActionTypes: [
        { name: "Single", type: "SINGLE" },
        { name: "Display", type: "DISPLAY" },
        { name: "Multi", type: "MULTI" },
      ],
      enumFontSizes: [
        { name: "8", value: 8 },
        { name: "10", value: 10 },
        { name: "12", value: 12 },
        { name: "14", value: 14 },
        { name: "16", value: 16 },
        { name: "18", value: 18 },
        { name: "20", value: 20 },
        { name: "22", value: 22 },
        { name: "24", value: 24 },
        { name: "26", value: 26 },
        { name: "28", value: 28 },
        { name: "32", value: 32 },
        { name: "36", value: 36 },
        { name: "40", value: 40 },
        { name: "44", value: 44 },
        { name: "48", value: 48 },
      ],
    };
  },
  watch: {
    action(action) {
      this.activeAction = action;
    },
  },
};
</script>

<style>
.action-panel-custom .p-panel-header {
  margin: 0px;
  padding: 2px !important;
}
</style>