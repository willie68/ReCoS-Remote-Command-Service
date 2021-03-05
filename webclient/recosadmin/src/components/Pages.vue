<template>
  <Splitter style="height: 500px">
    <SplitterPanel :size="20">
      <Panel header="Pages" class="pages-panel-custom">
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
          v-model="activePage"
          :options="profile.pages"
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
import PageSettings from "./PageSettings.vue"

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

.pages-panel-custom .p-panel-header {
    margin: 0px;
    padding: 2px !important;
}

</style>
