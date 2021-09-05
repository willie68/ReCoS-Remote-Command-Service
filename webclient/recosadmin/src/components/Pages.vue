<template>
  <Splitter style="height: 500px">
    <SplitterPanel :size="20">
      <Panel header="Pages" class="pages-panel-custom">
        <template #icons>
          <button
            class="p-panel-header-icon p-link p-mr-2 p-mt-0 p-mb-0 p-pt-0 p-pb-0"
            @click="addPage"
            v-tooltip="'create a new page'"
          >
            <span class="pi pi-plus"></span>
          </button>
          <button
            class="p-panel-header-icon p-link p-mr-2 p-mt-0 p-mb-0 p-pt-0 p-pb-0"
            @click="importPage"
            v-tooltip="'import a page'"
          >
            <span class="pi pi-cloud-upload"></span>
          </button>
          <button
            class="p-panel-header-icon p-link p-mr-2 p-mt-0 p-mb-0 p-pt-0 p-pb-0"
            @click="deleteConfirm"
            v-tooltip="'delete a page'"
          >
            <span class="pi pi-trash"></span>
          </button>
        </template>
        <Listbox
          v-model="activePage"
          :options="activeProfile.pages"
          optionLabel="name"
          listStyle="max-height:440px"
        >
        </Listbox>
      </Panel>
    </SplitterPanel>
    <SplitterPanel :size="80">
      <PageSettings :page="activePage" :profile="activeProfile"></PageSettings>
    </SplitterPanel>
  </Splitter>
  <AddName
    :visible="addPageDialog"
    v-model="newPageName"
    :excludeList="pageNames"
    v-on:save="saveNewPage($event)"
    v-on:cancel="this.addPageDialog = false"
    ><template #sourceHeader>New Page</template></AddName
  >
</template>

<script>
import PageSettings from "./PageSettings.vue";
import AddName from "./AddName.vue";
import { ObjectUtils } from "primevue/utils";

export default {
  name: "Pages",
  components: {
    PageSettings,
    AddName,
  },
  props: {
    profile: {},
  },
  data() {
    return {
      activeProfile: { name: "" },
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
      addPageDialog: false,
      pageNames: [],
      newPageName: "",
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
        toolbar: "show",
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
    addPage() {
      this.activeProfile.pages.forEach((element) => {
        this.pageNames.push(element.name);
      });
      this.addPageDialog = true;
    },
    saveNewPage(value) {
      console.log("Pages: add new page: " + value);
      this.addPageDialog = false;
      let newPage = {
        name: value,
        description: "Your description here",
        rows: 3,
        columns: 5,
        toolbar: "show",
      };
      this.activeProfile.pages.push(newPage);
      this.activePage = newPage;
    },
    deleteConfirm() {
      if (this.activePage) {
        console.log("Pages: delete confirm pressed");
        this.$confirm.require({
          message:
            "Deleting page: " +
            this.activePage.name +
            ". Are you sure you want to proceed?",
          header: "Confirmation",
          icon: "pi pi-exclamation-triangle",
          accept: () => {
            this.deletePage();
          },
          reject: () => {
            //callback to execute when user rejects the action
          },
        });
      }
    },
    deletePage() {
      if (this.activePage) {
        console.log("Commands: delete command " + this.activePage.name);
        let index = ObjectUtils.findIndexInList(
          this.activePage,
          this.activeProfile.pages
        );
        this.activeProfile.pages.splice(index, 1);
        if (this.activeProfile.pages.length > 0) {
          this.activePage = this.activeProfile.pages[0];
        } else {
          this.activePage = null;
        }
      }
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
  watch: {
    profile(profile) {
      if (profile && profile.name != "") {
        console.log("Pages: profile changed: ", profile.name);
        this.activeProfile = profile;
        if (this.activeProfile.pages) {
          if (this.activeProfile.pages.length > 0) {
            console.log("Pages: active page");
            this.activePage = this.activeProfile.pages[0];
            return;
          }
        }
      }
      this.activePage = { name: "" };
      this.activeProfile = { name: "" };
    },
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

.pages-panel-custom .p-panel-header {
  margin: 0px;
  padding: 2px !important;
}
</style>
