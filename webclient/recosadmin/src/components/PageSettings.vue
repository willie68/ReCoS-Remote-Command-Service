<template>
  <Panel>
    <template #header>
        <b>{{ page.name }}</b>
    </template>
    <template #icons>
      <Button class="p-panel-header-icon p-link p-mr-2" @click="toggle">
        <span class="pi pi-cog"></span>
      </Button>
      <Menu id="config_menu" ref="menu" :model="profileitems" :popup="true" />
    </template>

    <div class="p-fluid p-formgrid p-grid">
      <div class="p-field p-col">
        <label for="name">name</label>
        <InputText id="name" type="text" v-model="name" />
      </div>
      <div class="p-field p-col">
        <label for="rows">rows</label>
        <InputNumber id="rows" showButtons v-model="rows" :min="1" :max="10" />
      </div>
      <div class="p-field p-col">
        <label for="columns">columns</label>
        <InputNumber
          id="columns"
          showButtons
          v-model="columns"
          :min="1"
          :max="10"
        />
      </div>
    </div>
  </Panel>
  <ButtonPanel :rows="rows" :columns="columns"></ButtonPanel>
</template>

<script>
import ButtonPanel from "./ButtonPanel.vue";
export default {
  name: "PageSettings",
  components: {
    ButtonPanel,
  },
  props: {
    page: {},
  },
  data() {
    return {
      rows: 3,
      columns: 5,
      name: "",
      profileitems: [
        {
          label: "Options",
          items: [
            {
              label: "Update",
              icon: "pi pi-refresh",
              command: () => {
                this.$toast.add({
                  severity: "success",
                  summary: "Updated",
                  detail: "Data Updated",
                  life: 3000,
                });
              },
            },
            {
              label: "Delete",
              icon: "pi pi-times",
              command: () => {
                this.$toast.add({
                  severity: "warn",
                  summary: "Delete",
                  detail: "Data Deleted",
                  life: 3000,
                });
              },
            },
          ],
        },
        {
          label: "Navigate",
          items: [
            {
              label: "Vue Website",
              icon: "pi pi-external-link",
              url: "https://vuejs.org/",
            },
            {
              label: "Router",
              icon: "pi pi-upload",
              url: "https://vuejs.org/",
            },
          ],
        },
      ],
    };
  },
  methods: {
    toggle(event) {
      this.$refs.menu.toggle(event);
    },
    changePage() {},
  },
  mounted() {},
  created() {},
  watch: {
    page(page) {
      console.log("page changed");
      this.name = page.name;
      this.rows = page.rows;
      this.columns = page.columns;
    },
  },
};
</script>