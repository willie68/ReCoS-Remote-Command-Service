<template>
  <Dialog
    v-model:visible="dialogVisible"
    :modal="true"
    :closable="false"
    :style="{ width: '50vw' }"
  >
    <template #header>
      <h3>Settings</h3>
    </template>
    This is the settings dialog. Take a look to all tabs with different
    settings.
    <TabView>
      <TabPanel>
        <template #header>
          <img :src="'assets/recos.svg'" />
          <span class="ml-2">ReCoS Password</span>
        </template>
        In this Tab you can change the password on the ReCoS Admin Pages.
        <div class="field grid mt-2">
          <label for="password" class="col-12 mb-2 md:col-2 md:mb-0"
            >Password</label
          >
          <div class="col-12 md:col-10">
            <Password
              name="password"
              v-model="pwd"
              id="password"
              :feedback="false"
              autofocus
              @input="checkpwd"
            />
          </div>
        </div>
        <div class="field grid">
          <label for="newpassword" class="col-12 mb-2 md:col-2 md:mb-0"
            >new Password</label
          >
          <div class="col-12 md:col-10">
            <Password
              name="newpassword"
              v-model="newpwd"
              id="newpassword"
              toggleMask
              @input="checkpwd"
              class=""
            />
          </div>
        </div>
        <div class="field grid mb-2 mt-2">
          <label for="repeatpassword" class="col-12 mb-2 md:col-2 md:mb-0"
            >Repeat</label
          >
          <div class="col-12 md:col-10">
            <Password
              name="repeatpassword"
              v-model="rptpwd"
              id="repeatpassword"
              :feedback="false"
              @input="checkpwd"
              :class="{ 'p-valid': isPwdOK, 'p-invalid': !isPwdOK }"
            />
          </div>
        </div>
        <div class="field grid mb-2 mt-2">
          <label class="col-12 mb-2 ml-2 md:col-2 md:mb-0"></label>
          <div class="col-12 md:col-8">
            <Button
              label="Change password"
              icon="pi pi-pencil"
              @click="changepwd"
              :disabled="!isPwdOK"
            />
          </div>
        </div>
      </TabPanel>
    </TabView>
    <template #footer>
      <div class="pt-4">
        <Button label="Cancel" icon="pi pi-times" @click="cancel" />
        <Button
          label="Save"
          icon="pi pi-check"
          @click="save"
          :disabled="!isSaveOK"
        />
      </div>
    </template>
  </Dialog>
  <SelectIcon
    :visible="selectIconDialog"
    :iconlist="iconlist"
    @cancel="this.selectIconDialog = false"
    @save="this.saveIcon($event)"
    ><template #sourceHeader>Select Icon</template></SelectIcon
  >
</template>

<script>
import ArgumentList from "./../components/ArgumentList.vue";
import DropdownParameter from "./../components/DropdownParameter.vue";
import SelectIcon from "./../components/SelectIcon.vue";
import { ObjectUtils } from "primevue/utils";
import Password from "primevue/password";

export default {
  name: "Settings",
  components: {
    ArgumentList,
    DropdownParameter,
    SelectIcon,
    Password,
  },
  props: {
    profile: {},
    visible: Boolean,
  },
  emits: ["cancel", "save"],
  data() {
    return {
      dialogVisible: false,
      isFinishOK: false,
      ohmActive: false,
      apActive: false,
      phActive: false,
      hmActive: false,
      apSampletrates: [
        {
          name: "44,1kHz",
          value: 44100,
        },
        {
          name: "48kHz",
          value: 48000,
        },
      ],
      apSampleRate: 44100,
      settings: {},
      isSaveOK: false,
      activeParam: "",
      selectIconDialog: false,
      iconDestination: "",
      iconlist: [],
      pwd: "",
      newpwd: "",
      rptpwd: "",
      isPwdOK: false,
    };
  },
  methods: {
    cancel() {
      this.$emit("cancel");
    },
    saveEnable() {
      this.isSaveOK = true;
    },
    save() {
      let url = this.$baseURL + "config/integrations";
      this.settings.infos.forEach((info) => {
        fetch(url + "/" + info.name, {
          method: "POST",
          body: JSON.stringify(this.settings.settings[info.name]),
          headers: new Headers({
            "Content-Type": "application/json",
            Authorization: `Basic ${btoa(`admin:${this.$appStore.password}`)}`,
          }),
          mode: "cors",
        })
          .then((res) => res.json())
          .then((data) => {
            this.settings = data;
            console.log(JSON.stringify(data));
          })
          .catch((err) => console.log(err.message));
      });
      this.$toast.add({
        severity: "warn",
        summary: "Settings changed, please restart the service.",
        detail:
          "To use the new settings you have to restart the service. Please use the context menu on the service icon.",
      });
      this.$emit("save");
    },
    paramList(param) {
      if (!param.groupedlist) {
        let list = Array();
        var fieldName, filterValue, key, label, value;
        let filtered = false;
        if (param.filteredlist) {
          fieldName = param.filteredlist;
          filterValue = this.settings.settings[fieldName];
          filtered = true;
        }
        param.list.forEach((entry) => {
          label = entry;
          value = entry;
          if (filtered) {
            if (entry.indexOf(":") > 0) {
              key = entry.substring(0, entry.indexOf(":"));
              label = entry.substring(entry.indexOf(":") + 1);
              value = entry;
              console.log("filter", key, label, value);
              if (key != filterValue) {
                return;
              }
            }
          }
          let myValue = {
            label: label,
            value: value,
          };
          list.push(myValue);
        });
        return list;
      }
      if (param.groupedlist) {
        let list = Array();
        param.list.forEach((entry) => {
          var key, value;
          if (entry.indexOf(":") > 0) {
            key = entry.substring(0, entry.indexOf(":"));
            value = entry.substring(entry.indexOf(":") + 1);
          } else {
            key = "unknown";
            value = entry;
          }
          let found = false;
          for (let x = 0; x < list.length; x++) {
            if (list[x].label == key) {
              let myValue = {
                label: value,
                value: entry,
              };
              list[x].items.push(myValue);
              found = true;
            }
          }
          if (!found) {
            let myType = {
              label: key,
              items: Array(),
            };
            let myValue = {
              label: value,
              value: entry,
            };
            myType.items.push(myValue);
            list.push(myType);
          }
        });
        return list;
      }
    },
    addArgument(param) {
      this.activeParam = param;
      console.log("Step2: add argument for", param);
      this.addArgDialog = true;
    },
    saveNewArgument(data) {
      if (this.activeParam) {
        console.log("Step2: sav argument for", this.activeParam, data);
        if (!this.settings.settings[this.activeParam]) {
          this.settings.settings[this.activeParam] = [];
        }
        this.settings.settings[this.activeParam].push(data);
      }
      this.addArgDialog = false;
    },
    removeArgument(param, data) {
      let value = data;
      if (Array.isArray(data)) {
        value = data[0];
      }
      this.$confirm.require({
        message:
          "Deleting argument: " +
          param +
          ":" +
          value +
          ". Are you sure you want to proceed?",
        header: "Confirmation",
        icon: "pi pi-exclamation-triangle",
        accept: () => {
          this.deleteArgument(param, value);
        },
        reject: () => {
          //callback to execute when user rejects the action
        },
      });
    },
    deleteArgument(param, value) {
      let index = ObjectUtils.findIndexInList(
        value,
        this.settings.settings[param]
      );
      if (index >= 0) {
        this.settings.settings[param].splice(index, 1);
      }
    },
    saveIcon(icon) {
      console.log("Step2: save icon: ", icon, this.iconDestination);
      if (this.iconDestination == "") {
        this.localValue.icon = icon;
      } else {
        this.settings.settings[this.iconDestination] = icon;
      }
      this.selectIconDialog = false;
    },
    update() {
      this.$forceUpdate();
      this.saveEnable();
    },
    getSettings() {
      console.log("settings: activated");
      let url = this.$baseURL + "config/integrations";
      fetch(url, {
        method: "GET",
        mode: "cors",
        headers: new Headers({
          "Content-Type": "application/json",
          Authorization: `Basic ${btoa(`admin:${this.$appStore.password}`)}`,
        }),
      })
        .then((res) => res.json())
        .then((data) => {
          this.settings = data;
          console.log(JSON.stringify(data));
        })
        .catch((err) => console.log(err.message));
    },
    checkpwd() {
      if (this.newpwd == this.rptpwd) {
        this.isPwdOK = true;
      } else {
        this.isPwdOK = false;
      }
    },
    changepwd() {
      let url = this.$baseURL + "config/password";
      let settings = new Map();
      settings["password"] = this.pwd;
      settings["newpassword"] = this.newpwd;
      settings["repeatpassword"] = this.rptpwd;

      fetch(url, {
        method: "POST",
        body: JSON.stringify(settings),
        headers: new Headers({
          "Content-Type": "application/json",
          Authorization: `Basic ${btoa(`admin:${this.$appStore.password}`)}`,
        }),
        mode: "cors",
      }).then((response) => {
        if (!response.ok) {
          response.json().then((err) => {
            console.log(err);
            this.$toast.add({
              severity: "error",
              summary: "Error on Change Password",
              detail: err.message,
            });
          });
        } else {
          console.log(JSON.stringify(response));
          this.$toast.add({
            severity: "warn",
            summary: "Password changed",
            detail:
              "Please restart the service and the client for the change to take effect. Please use the context menu on the service icon.",
            life: 10000,
          });
        }
      });
    },
  },
  mounted() {
    this.iconlist = this.$iconlist;
  },
  beforeUnmount() {},
  watch: {
    visible(visible) {
      this.dialogVisible = visible;
      if (visible == true) {
        this.getSettings();
      }
    },
  },
};
</script>

<style scoped>
::v-deep(.p-password input) {
  width: 15rem;
}
</style>
