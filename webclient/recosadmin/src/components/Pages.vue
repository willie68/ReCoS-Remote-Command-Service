<template>
  <Toolbar class="p-mt-0 p-pt-0 p-pb-0">
    <template #left>
    <p class="p-ml-2">Pages:</p>
    <Dropdown
      class="p-ml-1 dropdownwidth"
      v-model="activePage"
      :options="profile.pages"
      optionLabel="name"
      placeholder="Select a Page"
    />
    <Button icon="pi pi-plus" @click="newPage" class="p-ml-1" />
    <Button icon="pi pi-trash" class="p-button-warning" />

    <p class="p-ml-2">Actions:</p>
    <Dropdown
      class="p-ml-1 dropdownwidth"
      v-model="selectedAction"
      :options="profile.actions"
      optionLabel="name"
      placeholder="Select an action"
    />
    <SplitButton
      v-tooltip="'Edit'"
      icon="pi pi-pencil"
      :model="actionMenuItems"
      class="p-button-warning"
    ></SplitButton>
    </template>
  </Toolbar>
  <Splitter style="height: 500px">
    <SplitterPanel :size="20">
      <div class="p-pb-2 p-pt-2"><b>Descriptions</b></div>
      <Fieldset
        :legend="'Profile: ' + profile.name"
        :toggleable="true"
        class="p-pt-2"
      >
        {{ profile.description }}
      </Fieldset>
      <Fieldset :legend="'Page: ' + activePage.name" :toggleable="true">
        {{ activePage.description }}
      </Fieldset>
    </SplitterPanel>
    <SplitterPanel :size="80">
      <PageSettings :page="activePage" :profile="profile"></PageSettings>
    </SplitterPanel>
  </Splitter>
</template>

<script>
import PageSettings from "./PageSettings.vue";

export default {
  name: "Pages",
  components: {
    PageSettings,
  },
  props: {
    profile: {},
  },
  data() {
    return {
      profileName: "",
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
      selectedAction: {},
      activePage: {},
      actionMenuItems: [
        {
          label: "Add",
          icon: "pi pi-plus",
        },
        {
          label: "Delete",
          icon: "pi pi-trash",
        },
      ],
    };
  },
  computed: {
    selectedPage: {
      get: function () {
        return this.activePage;
      },
      set: function (newPage) {
        if (newPage) {
          this.activePage = newPage;
        }
      },
    },
  },
  methods: {
    newPage() {
      this.activePage = {
        name: "Your Name Here",
        description: "Your description here",
        rows: 3,
        columns: 5,
      };
      this.activeProfile.pages.push(this.activePage);
      this.selectedPage = this.activePage;
    },
    toggle(event) {
      this.$refs.menu.toggle(event);
    },
    changePage(page) {
      this.selectedPage = page;
    },
 },
  mounted() {},
  created() {
    //    this.unsubscribe = this.$store.subscribe((mutation, state) => {
    //      if (mutation.type === "baseURL") {
    //      }
    //    });
  },
  beforeUnmount() {
    // this.unsubscribe();
  },
};
</script>

<style>
.p-panel-content {
  margin: 0px;
  padding: 0px !important;
}

.dropdownwidth {
  min-width: 12em;
}
</style>
