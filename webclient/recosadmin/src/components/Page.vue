<template>
  <Toolbar class="p-pb-1 p-pt-1">
    <template #left>
      <b>ReCoS Admin</b>
      <p class="p-ml-6">Profiles:</p>
      <Dropdown
        class="p-ml-1 dropdownwidth"
        v-model="activeProfileName"
        :options="profiles"
        placeholder="Select a Profile"
      />
      <SplitButton
        v-tooltip="'Edit'"
        icon="pi pi-pencil"
        :model="profileMenuItems"
        class="p-button-warning"
        @click="dialogProfileVisible = true"
      ></SplitButton>

      <p class="p-ml-2">Pages:</p>
      <Dropdown
        class="p-ml-1 dropdownwidth"
        v-model="activePage"
        :options="activeProfile.pages"
        optionLabel="name"
        placeholder="Select a Page"
      />
      <SplitButton
        v-tooltip="'Edit'"
        icon="pi pi-pencil"
        :model="pageMenuItems"
        class="p-button-warning"
      ></SplitButton>

      <p class="p-ml-2">Actions:</p>
      <Dropdown
        class="p-ml-1 dropdownwidth"
        v-model="selectedAction"
        :options="activeProfile.actions"
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

    <template #right>
      <span class="p-input-icon-right">
        Password
        <InputText ref="pwd" v-model="pwd" :type="pwdType"></InputText>
        <i class="pi pi-eye-slash" @click="togglePwdView()" />
        <i v-if="!showPwd" class="pi pi-eye-slash" @click="togglePwdView()" />
        <i v-if="showPwd" class="pi pi-eye" @click="togglePwdView()" />
      </span>
      <SplitButton
        label="Save"
        icon="pi pi-check"
        :model="items"
        class="p-button-warning"
      ></SplitButton>
      <Button icon="pi pi-cog" class="p-mr-1" />
    </template>
  </Toolbar>
  <Splitter style="height: 500px">
    <SplitterPanel :size="20">
      <div class="p-pb-2 p-pt-2"><b>Descriptions</b></div>
      <Fieldset
        :legend="'Profile: ' + activeProfile.name"
        :toggleable="true"
        class="p-pt-2"
      >
        {{ activeProfile.description }}
      </Fieldset>
      <Fieldset :legend="'Page: ' + activePage.name" :toggleable="true">
        {{ activePage.description }}
      </Fieldset>
    </SplitterPanel>
    <SplitterPanel :size="80">
      <PageSettings :page="activePage" :profile="activeProfile"></PageSettings>
    </SplitterPanel>
  </Splitter>

  <EditProfile
    :visible="dialogProfileVisible"
    v-on:save="saveProfile"
    v-on:cancel="this.dialogProfileVisible = false"
  ></EditProfile>
</template>

<script>
import PageSettings from "./PageSettings.vue";
import EditProfile from "./EditProfile.vue";

export default {
  name: "Page",
  components: {
    PageSettings,
    EditProfile,
  },
  props: {
    msg: String,
  },
  data() {
    return {
      pwd: "",
      showPwd: false,
      pwdType: "password",
      profiles: [],
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
      activeProfile: {},
      activePage: { name: "" },
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
      selectedAction: {},
      profileMenuItems: [
        {
          label: "Add",
          icon: "pi pi-plus",
        },
        {
          label: "Delete",
          icon: "pi pi-trash",
          class: "p-button-warning",
        },
      ],
      pageMenuItems: [
        {
          label: "Add",
          icon: "pi pi-plus",
        },
        {
          label: "Delete",
          icon: "pi pi-trash",
        },
      ],
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
      dialogProfileVisible: false,
      editProfile: { name: "", description: "" },
    };
  },
  computed: {
    activeProfileName: {
      get: function () {
        return this.profileName;
      },
      set: function (newProfile) {
        if (newProfile) {
          this.profileName = newProfile;
          fetch(this.profileURL + "/" + this.profileName)
            .then((res) => res.json())
            .then((data) => {
              this.activeProfile = data;
              this.activePage = this.activeProfile.pages[0];
              console.log(this.activeProfile);
            })
            .catch((err) => console.log(err.message));
        }
      },
    },
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
    openProfile(e) {
      console.log("open profile:" + e);
      this.profileName = this.profiles[e.index];
      fetch(this.profileURL + "/" + this.profileName)
        .then((res) => res.json())
        .then((data) => {
          this.activeProfile = data;
          this.activePage = this.activeProfile.pages[0];
          console.log(this.activeProfile);
        })
        .catch((err) => console.log(err.message));
    },
    togglePwdView() {
      this.showPwd = !this.showPwd;
      if (this.showPwd) {
        this.pwdType = "text";
      } else {
        this.pwdType = "password";
      }
    },
    toggle(event) {
      this.$refs.menu.toggle(event);
    },
    changePage(page) {
      this.selectedPage = page;
    },
    saveProfile: function (profile) {
      console.log("Save profile:" + profile.name + "#" + profile.description)
      this.dialogProfileVisible = false
    },
  },
  mounted() {},
  created() {
    let that = this;
    this.unsubscribe = this.$store.subscribe((mutation, state) => {
      if (mutation.type === "baseURL") {
        console.log(`Updating to ${state.baseURL}`);
        that.profileURL = state.baseURL + "/profiles";
        console.log("page profiles url:" + that.profileURL);

        fetch(that.profileURL)
          .then((res) => res.json())
          .then((data) => {
            console.log(data);
            that.profiles = data;
            that.activeProfileName = that.profiles[0];
            console.log(that.profiles);
          })
          .catch((err) => console.log(err.message));
      }
    });
  },
  beforeUnmount() {
    this.unsubscribe();
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
