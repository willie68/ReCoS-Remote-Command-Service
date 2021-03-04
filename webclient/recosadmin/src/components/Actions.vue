<template>
  <Splitter style="height: 500px">
    <SplitterPanel :size="20">
      <Panel header="Actions" class="actions-panel-custom">
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
        <Listbox
          v-model="activeAction"
          :options="profile.actions"
          optionLabel="name"
          listStyle="max-height:440px"
        >
        </Listbox>
      </Panel>
    </SplitterPanel>
    <SplitterPanel :size="80">
      <PageSettings :page="activePage" :profile="profile"></PageSettings>
    </SplitterPanel>
  </Splitter>
</template>

<script>
export default {
  name: "Actions",
  props: {
    profile: {},
  },
  data() {
    return {
      activeAction: {},
    };
  },
};
</script>

<style>
.actions-panel-custom .p-panel-header {
    margin: 0px;
    padding: 2px !important;
}
</style>