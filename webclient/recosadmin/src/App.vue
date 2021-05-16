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
        v-tooltip="'Save'"
        icon="pi pi-save"
        :model="profileMenuItems"
        class="p-button-warning"
        @click="saveProfile()"
      ></SplitButton>
      <Button
        icon="pi pi-flag"
        class="p-mr-1 p-button-warning"
        @click="actionWizard()"
      />
      <div v-if="profileDirty">*</div>
    </template>

    <template #right>
      <span v-if="isPWDNeeded" class="p-input-icon-right">
        Password
        <InputText
          ref="pwd"
          v-model="password"
          :type="pwdType"
          name="password"
          :class="{ 'p-valid': isPwdOK, 'p-invalid': !isPwdOK }"
        />
        <i class="pi pi-eye-slash" @click="togglePwdView()" />
        <i v-if="!showPwd" class="pi pi-eye-slash" @click="togglePwdView()" />
        <i v-if="showPwd" class="pi pi-eye" @click="togglePwdView()" />
      </span>
      <Button icon="pi pi-bars" class="p-mr-1" @click="toggleHelpMenu" />
      <Menu
        id="overlay_menu"
        ref="helpmenu"
        :model="helpMenuItems"
        :popup="true"
      />
    </template>
  </Toolbar>

  <Profile :profile="activeProfile"></Profile>
  <AppFooter></AppFooter>
  <AddProfile
    :visible="dialogProfileVisible"
    v-on:save="saveNewProfile($event)"
    v-on:cancel="this.dialogProfileVisible = false"
    :profiles="profiles"
  ></AddProfile>
  <Toast position="top-right" />
  <ConfirmDialog></ConfirmDialog>
  <ActionWizard
    :visible="actionWizardVisible"
    v-on:save="saveWizardProfile($event)"
    v-on:cancel="this.actionWizardVisible = false"
    :profile="activeProfile"
  ></ActionWizard>
  <Settings
    :visible="settingsVisible"
    v-on:save="saveSettings($event)"
    v-on:cancel="this.settingsVisible = false"
  ></Settings>
  <Dialog header="About" v-model:visible="helpAboutVisible">
    This is the about dialog. For more info see:<br />
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
import Profile from "./components/Profile.vue";
import AppFooter from "./components/AppFooter.vue";
import AddProfile from "./components/AddProfile.vue";
import ActionWizard from "./wizard/ActionWizard.vue";
import Settings from "./settings/Settings.vue";

export default {
  components: {
    Profile,
    AppFooter,
    AddProfile,
    ActionWizard,
    Settings,
  },
  data() {
    return {
      name: "RecosAdmin",
      isPWDNeeded: true,
      pwd: "",
      showPwd: false,
      pwdType: "password",
      isPwdOK: false,
      profiles: [],
      activeProfile: {},
      profileDirty: false,
      profileName: "",
      profileMenuItems: [
        {
          label: "Add",
          icon: "pi pi-plus",
          command: () => {
            this.createProfile();
          },
        },
        {
          label: "Delete",
          icon: "pi pi-trash",
          class: "p-button-warning",
          command: () => {
            this.deleteProfile();
          },
        },
        {
          label: "Export",
          icon: "pi pi-cloud-download",
          class: "p-button-warning",
          command: () => {
            this.exportProfile();
          },
        },
      ],
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
          label: "Settings",
          icon: "pi pi-cog",
          command: () => {
            this.helpSettings();
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
      dialogProfileVisible: false,
      actionWizardVisible: false,
      settingsVisible: false,
      helpAboutVisible: false,
      helpCreditsVisible: false,
      credits: "",
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
            headers: this.$store.state.authheader,
          })
            .then((res) => res.json())
            .then((data) => {
              that.activeProfile = data;
              that.profileDirty = false;
              console.log("profile dirty: " + that.profileDirty);
            })
            .catch((err) => console.log(err.message));
        }
      },
    },
    password: {
      get: function () {
        return this.$store.state.password;
      },
      set: function (newPassword) {
        if (newPassword) {
          let that = this;
          that.isPwdOK = false;
          fetch(this.$store.state.baseURL + "config/check", {
            method: "GET",
            headers: new Headers({
              Authorization: `Basic ${btoa(`admin:${newPassword}`)}`,
            }),
          })
            .then((response) => {
              if (response.ok) {
                console.log("authentication check ok");
                that.isPwdOK = true;
                that.$store.commit("password", newPassword);
              }
            })
            .catch((err) => {
              console.log(err.message);
            });
        }
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
        this.needPWd();
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
    let iconurl = basepath + "config/icons";
    fetch(iconurl)
      .then((res) => res.json())
      .then((data) => {
        //console.log(data);
        that.$store.commit("iconlist", data);
      })
      .catch((err) => console.log(err.message));
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
    helpSettings() {
      this.settingsVisible = true;
    },
    toggleHelpMenu(event) {
      this.$refs.helpmenu.toggle(event);
    },
    needPWd() {
      let that = this;
      fetch(this.$store.state.baseURL + "config/check")
        .then((response) => {
          if (response.ok) {
            console.log("APP: no authentication needed");
            that.isPWDNeeded = false;
            that.isPwdOK = true;
            that.$store.commit("password", "");
          } else {
            console.log("APP: authentication needed");
          }
        })
        .catch((err) => {
          console.log("APP: authentication needed");
          console.log(err.message);
        });
    },
    actionWizard() {
      if (!this.isPwdOK) {
        this.$toast.add({
          severity: "warn",
          summary: "Please enter the password",
          detail: "To use the wizard please enter the right password.",
          life: 5000,
        });
      } else {
        if (this.activeProfile) {
          this.actionWizardVisible = true;
        }
      }
    },
    saveWizardProfile(profile) {
      console.log("App: ", JSON.stringify(profile));
      this.actionWizardVisible = false;
      this.activeProfile = profile;
      this.saveProfile();
    },
    togglePwdView() {
      this.showPwd = !this.showPwd;
      if (this.showPwd) {
        this.pwdType = "text";
      } else {
        this.pwdType = "password";
      }
    },
    createProfile() {
      this.dialogProfileVisible = true;
    },
    saveProfile() {
      console.log("App: Save profile:" + this.activeProfile.name);
      //      console.log(JSON.stringify(this.activeProfile));
      this.dialogProfileVisible = false;
      fetch(this.profileURL + "/" + this.activeProfile.name, {
        method: "PUT",
        body: JSON.stringify(this.activeProfile),
        headers: new Headers({
          "Content-Type": "application/json",
          Authorization: `Basic ${btoa(`admin:${this.$store.state.password}`)}`,
        }),
      })
        .then((response) => {
          if (!response.ok) {
            response.json().then((err) => {
              console.log(err);
              this.$toast.add({
                severity: "error",
                summary: "Error on save",
                detail: err.message,
                life: 3000,
              });
            });
          }
        })
        .catch((err) => {
          console.log(err.message);
          this.$toast.add({
            severity: "error",
            summary: "Error on save",
            detail: err.message,
            life: 3000,
          });
        });
    },
    saveNewProfile(profile) {
      this.dialogProfileVisible = false;
      console.log(JSON.stringify(profile));
      console.log("Create profile:" + profile.name);
      fetch(this.profileURL + "/", {
        method: "POST",
        body: JSON.stringify(profile),
        headers: new Headers({
          "Content-Type": "application/json",
          Authorization: `Basic ${btoa(`admin:${this.$store.state.password}`)}`,
        }),
      })
        .then((response) => {
          if (!response.ok) {
            response.json().then((err) => {
              console.log(err);
              this.$toast.add({
                severity: "error",
                summary: "Create",
                detail: err.message,
                life: 3000,
              });
            });
          } else {
            this.profiles.push(profile.name);
            this.activeProfileName = profile.name;
          }
        })
        .catch((err) => {
          console.log(err.message);
          this.$toast.add({
            severity: "error",
            summary: "Create",
            detail: err.message,
            life: 3000,
          });
        });
    },
    deleteProfile() {
      console.log("Delete profile:" + this.activeProfile.name);
      fetch(this.profileURL + "/" + this.activeProfile.name, {
        method: "DELETE",
        headers: new Headers({
          Authorization: `Basic ${btoa(`admin:${this.$store.state.password}`)}`,
        }),
      })
        .then((response) => {
          if (!response.ok) {
            response.json().then((err) => {
              console.log(err);
              this.$toast.add({
                severity: "error",
                summary: "Delete",
                detail: err.message,
                life: 3000,
              });
            });
          } else {
            this.profiles.splice(this.profiles.indexOf(this.activeProfile), 1);
            if (this.profiles.length > 0) {
              this.activeProfileName = this.profiles[0];
            }
            this.activeProfileName = "";
          }
        })
        .catch((err) => {
          console.log(err.message);
          this.$toast.add({
            severity: "error",
            summary: "Delete",
            detail: err.message,
            life: 3000,
          });
        });
    },
    exportProfile() {
      console.log("export profile: " + this.activeProfileName);
      window.open(this.profileURL + "/" + this.activeProfileName + "/export");
    },
    saveSettings() {
      this.settingsVisible = false;
    },
  },
  watch: {
    profile(newProfile) {
      if (newProfile) {
        console.log("changing profile to " + newProfile.name);
        this.activeProfile = newProfile;
      }
    },
    activeProfile: {
      deep: true,
      handler(newProfile) {
        console.log("app: changing profile " + newProfile.name);
        //console.log(JSON.stringify(newProfile));
        this.profileDirty = true;
      },
    },
  },
};
</script>

<style>
html {
  font-size: 14px;
  background-color: #1f2d40;
  color: #ffffff;
}

body {
  margin: 0px;
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #ffffff;
}

.profiledirty {
  background: lightsalmon;
}

.p-inputtext.p-valid.p-component {
  border-color: #36f43f;
}
</style>
