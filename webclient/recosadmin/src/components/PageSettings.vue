<template>
  <Panel class="page-panel-custom" v-if="activePage && activePage.name != ''">
    <template #header>
      <b>{{ profile.name }} # {{ activePage.name }}</b>
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
        <InputText id="name" type="text" v-model="activePage.name" />
      </div>
      <div class="p-field p-col">
        <label for="icon">Icon</label>
        <span class="p-input-icon-right">
          <InputText
            id="icon"
            v-model="activePage.icon"
            placeholder="select a icon"
          />
          <i class="pi pi-chevron-down" @click="selectIconDialog = true" />
        </span>
      </div>
      <div class="p-field p-col">
        <label for="rows">rows</label>
        <InputNumber
          id="rows"
          showButtons
          v-model="activePage.rows"
          :min="1"
          :max="10"
        />
      </div>
      <div class="p-field p-col">
        <label for="columns">columns</label>
        <InputNumber
          id="columns"
          showButtons
          v-model="activePage.columns"
          :min="1"
          :max="10"
        />
      </div>
      <div class="p-field p-col">
        <label for="rows">Type</label>
        <Dropdown
          v-model="activePage.toolbar"
          :options="enumPageTypes"
          placeholder="select a toolbar type"
          optionLabel="name"
          optionValue="type"
        />
      </div>
    </div>
  </Panel>
  <ButtonPanel
    v-if="activePage"
    :actions="profile.actions"
    :page="activePage"
    :profile="profile"
  ></ButtonPanel>
  <SelectIcon
    :visible="selectIconDialog"
    :iconlist="iconlist"
    @cancel="this.selectIconDialog = false"
    @save="this.saveIcon($event)"
    ><template #sourceHeader>Select Icon</template></SelectIcon
  >
</template>

<script>
import ButtonPanel from "./ButtonPanel.vue";
import SelectIcon from "./SelectIcon.vue";

export default {
  name: "PageSettings",
  components: {
    ButtonPanel,
    SelectIcon,
  },
  props: {
    page: {},
    profile: { name: "" },
  },
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
      activePage: {},
      enumPageTypes: [
        { name: "Show", type: "show" },
        { name: "Hide", type: "hide" },
      ],
      iconlist: [],
      selectIconDialog: false,
    };
  },
  methods: {
    toggle(event) {
      this.$refs.menu.toggle(event);
    },
    changePage() {},
    saveIcon(icon) {
      console.log("Action: save icon: " + icon);
      this.activePage.icon = icon;
      this.selectIconDialog = false;
    },
  },
  mounted() {
    this.iconlist = this.$store.state.iconlist;
    let that = this;
    this.unsubscribe = this.$store.subscribe((mutation, state) => {
      if (mutation.type === "iconlist") {
        that.iconlist = state.iconlist;
      }
    });
  },
  beforeUnmount() {
    this.unsubscribe();
  },
  created() {},
  watch: {
    page(page) {
      if (page) {
        this.activePage = page;
      }
    },
  },
};
</script>


<style>
.page-panel-custom .p-panel-header {
  margin: 0px;
  padding: 2px !important;
}
</style>
