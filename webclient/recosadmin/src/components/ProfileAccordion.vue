<template>
  <Panel header="Profiles">
    <template #icons>
      <Button class="p-panel-header-icon p-link p-mr-2" @click="toggle">
        <span class="pi pi-cog"></span>
      </Button>
      <Menu id="config_menu" ref="menu" :model="profileitems" :popup="true" />
    </template>
  </Panel>
</template>

<script>
export default {
  name: "ProfileAccordion",
  props: {},
  data() {
    return {
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
  },
};
</script>

<style scoped>
</style>