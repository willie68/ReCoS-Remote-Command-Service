<template>
  <Sidebar
    v-model:visible="visible"
    :baseZIndex="1000"
    position="top"
    class="p-sidebar-sm"
  >
    <div class="p-fluid">
      <div class="p-field p-grid">
        <label for="profiles" class="p-col-12 p-mb-2 p-md-2 p-mb-md-0"
          >Profiles</label
        >
        <div class="p-col-12 p-md-10">
          <Dropdown
            id="profiles"
            v-tooltip="'select the profile to edit'"
            class="p-ml-1 dropdownwidth"
            v-model="activeProfileName"
            :options="profiles"
            placeholder="Select a Profile"
          />
        </div>
      </div>
    </div>
    <label class="p-col-12 p-mb-2 p-md-2 p-mb-md-0">Pages</label>
    <Button
      v-for="page in toolbarPages"
      :value="page.name"
      :key="page.name"
      :title="page.description"
      :label="page.name"
      @click="changePage(page.name)"
      class="p-ml-1"
    >
      <img v-if="page.icon" :src="buildImageSrc(page.icon)" height="20" />
    </Button>
  </Sidebar>
  <Button
    style="position: absolute; top: 4px; left: 4px; z-index: 10000"
    icon="pi pi-bars"
    class="p-mr-1 p-button-sm"
    @click="toggleHelpMenu"
  />
  <Menu id="overlay_menu" ref="helpmenu" :model="helpMenuItems" :popup="true" />
  <Button
    style="position: absolute; top: 4px; right: 4px; z-index: 10000"
    icon="pi pi-cog"
    @click="visible = true"
    v-if="!visible"
    class="p-button-sm"
  />
  <Dialog header="About" v-model:visible="helpAboutVisible">
    This is ReCoS V{{ this.$store.state.packageVersion }} <br />For more info
    see:<br />
    <a
      href="https://github.com/willie68/ReCoS-Remote-Command-Service"
      target="_blank"
      >ReCoS on github</a
    >
    <template #footer>
      <Button label="OK" icon="pi pi-check" @click="closeHelpAbout" autofocus />
    </template>
  </Dialog>
  <Dialog header="Credits" v-model:visible="helpCreditsVisible">
    <div v-html="credits"></div>
    <template #footer>
      <Button
        label="OK"
        icon="pi pi-check"
        @click="closeHelpCredits"
        autofocus
      />
    </template>
  </Dialog>
</template>

<script>
export default {
  name: "AppConfig",
  components: {},
  emits: ["profileChanged", "pageChanged"],
  props: [  ],
  data() {
    return {
      profileURL: "",
      visible: false,
      helpMenuItems: [
        {
          label: "Help",
          icon: "pi pi-question-circle",
          command: () => {
            this.helpHelp();
          },
        },
        {
          separator: true,
        },
        {
          label: "Credits",
          icon: "pi pi-star",
          command: () => {
            this.helpCredits();
          },
        },
        {
          label: "About",
          icon: "pi pi-info-circle",
          command: () => {
            this.helpAbout();
          },
        },
      ],
      helpAboutVisible: false,
      helpCreditsVisible: false,
      profileName: "",
      activeProfile: {},
    };
  },
  computed: {
    activeProfileName: {
      get: function () {
        return this.profileName;
      },
      set: function (newProfile) {
        if (newProfile) {
          let that = this;
          this.profileName = newProfile;
          fetch(this.profileURL + "/" + this.profileName, {
            method: "GET",
          })
            .then((res) => res.json())
            .then((data) => {
              that.activeProfile = data;
              //console.log(JSON.stringify(that.activeProfile))
              this.$emit("profileChanged", that.activeProfile);
              if (that.activeProfile.pages && (that.activeProfile.pages.length > 0)) {
                  that.changePage(that.activeProfile.pages[0].name)
              }
            })
            .catch((err) => console.log(err.message));
        }
      },
    },
    toolbarPages: {
      get: function () {
        console.log("new pages");
        var a = [];
        if (this.activeProfile) {
          if (this.activeProfile.pages) {
            this.activeProfile.pages.forEach((page) => {
              if (page.toolbar != "hide") {
                a.push(page);
              }
            });
          }
        }
        return a;
      },
    },
  },
  mounted() {
    let servicePort = this.$store.state.servicePort;
    let basepath =
      window.location.protocol +
      "//" +
      window.location.hostname +
      ":" +
      servicePort +
      "/api/v1/";
    this.$store.commit("baseURL", basepath);
    console.log(`Updating to ${basepath}`);
    let that = this;
    that.profileURL = basepath + "profiles";
    console.log("page profiles url:" + that.profileURL);

    fetch(that.profileURL)
      .then((res) => res.json())
      .then((data) => {
        that.profiles = data;
        that.activeProfileName = that.profiles[0];
      })
      .catch((err) => {
        console.log(err.message);
        this.$toast.add({
          severity: "error",
          summary: "Error loading profile",
          detail: err.message,
          life: 3000,
        });
      });
    fetch(basepath + "config/credits")
      .then((response) => {
        return response.text();
      })
      .then((data) => {
        this.credits = data;
      });
  },
  methods: {
    helpHelp() {
      window
        .open(
          "https://raw.githubusercontent.com/willie68/ReCoS-Remote-Command-Service/master/documentation/README.pdf",
          "_blank"
        )
        .focus();
    },
    helpAbout() {
      this.helpAboutVisible = true;
    },
    closeHelpAbout() {
      this.helpAboutVisible = false;
    },
    helpCredits() {
      this.helpCreditsVisible = true;
    },
    closeHelpCredits() {
      this.helpCreditsVisible = false;
    },
    toggleHelpMenu(event) {
      this.$refs.helpmenu.toggle(event);
    },
    changePage(pageName) {
      // console.log("change page to: ", pageName);
      this.pageName = pageName;
      this.activeProfile.pages.forEach((page) => {
        if (pageName == page.name) {
          this.activePage = page;
        }
      });
      this.$emit("pageChanged", this.activePage);
      this.visible = false;

      // console.log("adding actions to page: " + this.activePage.name);
      // console.log("cell count: " + this.activePage.cells.length);
      this.cellactions = new Array(this.activePage.rows);
      for (let x = 0; x < this.activePage.rows; x++) {
        this.cellactions[x] = new Array(this.activePage.columns);
        for (let y = 0; y < this.activePage.columns; y++) {
          var action = undefined;
          let index = x * this.activePage.columns + y;
          if (index < this.activePage.cells.length) {
            let actionName = this.activePage.cells[index];
            if (actionName) {
              this.activeProfile.actions.forEach((profileAction) => {
                if (profileAction.name == actionName) {
                  action = profileAction;
                }
              });
              if (action) {
                // console.log("adding action (" + x + "," + y + ") " + index + ":" + action.name );
                this.cellactions[x][y] = action;
              } else if (actionName == "none") {
                this.cellactions[x][y] = {
                  type: "DISPLAY",
                  title: "",
                  name: "none",
                };
              } else {
                // console.log("missing action (" + x + "," + y + ") " + index);
                this.cellactions[x][y] = {
                  type: "DISPLAY",
                  title: "Action not defined",
                  name: "none",
                  fontcolor: "#FF0000",
                };
              }
              continue;
            }
          }
          // console.log("adding none(" + x + "," + y + ") " + index);
          this.cellactions[x][y] = {
            type: "NONE",
          };
        }
      }
    },
    buildImageSrc(data) {
      if (!data) {
        console.log("no data");
        return "data:image/bmp;base64,Qk1CAAAAAAAAAD4AAAAoAAAAAQAAAAEAAAABAAEAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP///wCAAAAA";
      }
      if (data.startsWith("/")) {
        return this.baseURL + data;
      }
      if (data.startsWith("data:")) {
        return data;
      }
      return "assets/" + data;
    },
  },
  watch: {},
  beforeUnmount() {},
};
</script>

<style>
</style>