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
      <div v-if="profileDirty"> *</div>
    </template>

    <template #right>
      <span class="p-input-icon-right">
        Password
        <InputText
          ref="pwd"
          v-model="password"
          :type="pwdType"
          name="password"
          :class="{ passwordOK: isPwdOK, passwordMissing: !isPwdOK }"
        />
        <i class="pi pi-eye-slash" @click="togglePwdView()" />
        <i v-if="!showPwd" class="pi pi-eye-slash" @click="togglePwdView()" />
        <i v-if="showPwd" class="pi pi-eye" @click="togglePwdView()" />
      </span>
      <Button icon="pi pi-cog" class="p-mr-1" />
    </template>
  </Toolbar>

  <Profile :profile="activeProfile"></Profile>
  <AppFooter></AppFooter>
  <EditProfile
    :visible="dialogProfileVisible"
    v-on:save="saveProfile"
    v-on:cancel="this.dialogProfileVisible = false"
  ></EditProfile>
</template>

<script>
import Profile from "./components/Profile.vue";
import AppFooter from "./components/AppFooter.vue";
import EditProfile from "./components/EditProfile.vue";

export default {
  components: {
    Profile,
    AppFooter,
    EditProfile,
  },
  data() {
    return {
      name: "RecosAdmin",
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
        },
        {
          label: "Delete",
          icon: "pi pi-trash",
          class: "p-button-warning",
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
      dialogProfileVisible: false,
      editProfile: { name: "", description: "" },
    };
  },
  computed: {
    activeProfileName: {
      get: function () {
        return this.profileName
      },
      set: function (newProfile) {
        if (newProfile) {
          let that = this
          this.profileName = newProfile
          fetch(this.profileURL + "/" + this.profileName, {
            method: "GET",
            headers: this.$store.state.authheader
          })
            .then((res) => res.json())
            .then((data) => {
              that.activeProfile = data
              that.profileDirty = false
              console.log("profile dirty: "+ that.profileDirty)
            })
            .catch((err) => console.log(err.message))
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
          fetch(this.$store.state.baseURL + "/config/check", {
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
      "/api/v1";
    this.$store.commit("baseURL", basepath);
    console.log(`Updating to ${basepath}`);
    let that = this;

    that.profileURL = basepath + "/profiles";
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
  },
  methods: {
    togglePwdView() {
      this.showPwd = !this.showPwd;
      if (this.showPwd) {
        this.pwdType = "text";
      } else {
        this.pwdType = "password";
      }
    },
    saveProfile(profile) {
      console.log("Save profile:" + profile.name + "#" + profile.description);
      this.dialogProfileVisible = false;
    },
    exportProfile() {
      console.log("export profile: " + this.activeProfileName);
      window.open(this.profileURL + "/" + this.activeProfileName + "/export");
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
        this.profileDirty = true;
      },
    },
  },
};
</script>

<style>
html {
  font-size: 14px;
}

body {
  margin: 0px;
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

.passwordOK {
  background: lightgreen !important;
}

.passwordMissing {
  background: lightcoral !important;
}

.profiledirty {
  background: lightsalmon;
}
</style>
