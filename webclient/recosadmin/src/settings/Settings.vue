<template>
  <Dialog v-model:visible="dialogVisible" :modal="true" :closable="false">
    <template #header>
      <h3>Settings</h3>
    </template>
    This is the settings dialog. Take a look to all tabs with different
    settings.
    <TabView>
      <TabPanel v-for="info of settings.infos" :key="info.name">
        <template #header>
          <img :src="'assets/' + info.image" />
          <i :src="info.image"></i>
          <span> {{ info.name }} </span>
        </template>
        <div
          class="p-field p-grid p-mb-2 p-mt-2"
          v-for="param of info.parameter"
          :key="param.name"
        >
          <label
            :for="param.name"
            class="p-col-12 p-mb-2 p-ml-2 p-md-2 p-mb-md-0"
            >{{ param.name }}</label
          >
          <div class="p-col-12 p-md-8">
            <InputText
              v-if="
                param.type == 'string' &&
                (!param.list || param.list.length == 0)
              "
              :id="param.name"
              type="text"
              v-tooltip="param.description"
              v-model="settings.settings[info.name][param.name]"
              class="fullwidth"
              @change="saveEnable()"
            />
            <DropdownParameter
              v-if="
                param.type == 'string' &&
                param.list &&
                param.list.length > 0 &&
                !param.groupedlist
              "
              :id="param.name"
              :options="paramList(param)"
              optionLabel="label"
              optionValue="value"
              v-tooltip="param.description"
              v-model="settings.settings[info.name][param.name]"
              :placeholder="param.unit"
              scrollHeight="120px"
              class="fullwidth"
              @change="update()"
            />
            <Dropdown
              v-if="
                param.type == 'string' &&
                param.list &&
                param.list.length > 0 &&
                param.groupedlist
              "
              :id="param.name"
              :options="paramList(param)"
              optionLabel="label"
              optionValue="value"
              optionGroupLabel="label"
              optionGroupChildren="items"
              v-tooltip="param.description"
              v-model="settings.settings[info.name][param.name]"
              :placeholder="param.unit"
              @change="update()"
              ><template #optiongroup="slotProps">
                <div class="p-d-flex p-ai-center country-item">
                  <img
                    :src="
                      this.$store.state.baseURL +
                      'config/icons/category/' +
                      slotProps.option.label
                    "
                    width="18"
                    class="p-mr-2"
                  />
                  <div>
                    <b>{{ slotProps.option.label }}</b>
                  </div>
                </div>
              </template>
            </Dropdown>
            <InputNumber
              v-if="param.type == 'int'"
              :id="param.name"
              type="text"
              mode="decimal"
              showButtons
              :suffix="param.unit"
              v-tooltip="param.description"
              v-model="settings.settings[info.name][param.name]"
              class="fullwidth"
            />
            <Checkbox
              v-if="param.type == 'bool'"
              :id="param.name"
              v-tooltip="param.description"
              v-model="settings.settings[info.name][param.name]"
              :binary="true"
              class="fullwidth"
              @change="update()"
            />
            <ArgumentList
              v-if="param.type == '[]string'"
              :id="param.name"
              v-tooltip="param.description"
              v-model="settings.settings[info.name][param.name]"
              v-on:add="addArgument(param.name)"
              v-on:remove="removeArgument(param.name, $event)"
              class="fullwidth"
            >
              <template #item="slotProps">
                <div>
                  {{ slotProps.item }}
                </div>
              </template></ArgumentList
            >
            <ColorPicker
              v-if="param.type == 'color'"
              :id="param.name"
              :inline="false"
              defaultColor="#00FF00"
              v-tooltip="param.description"
              v-model="settings.settings[info.name][param.name]"
              class="fullwidth"
              @change="update()"
            />
            <Calendar
              v-if="param.type == 'date'"
              :id="param.name"
              :inline="false"
              dateFormat="yy-mm-dd"
              :showIcon="true"
              v-model="settings.settings[info.name][param.name]"
              @change="update()"
            />
            <span
              v-if="param.type == 'icon'"
              class="p-input-icon-right fullwidth"
            >
              <InputText
                :id="param.name"
                v-model="settings.settings[info.name][param.name]"
                v-tooltip="param.description"
                placeholder="select a icon"
                class="fullwidth"
              />
              <i
                class="pi pi-chevron-down"
                @click="
                  iconDestination = param.name;
                  selectIconDialog = true;
                "
              />
            </span>
          </div>
        </div>
      </TabPanel>
    </TabView>
    <template #footer>
      <div class="p-pt-4">
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
import ArgumentList from "./../wizard/ArgumentList.vue";
import DropdownParameter from "./../wizard/DropdownParameter.vue";
import SelectIcon from "./../components/SelectIcon.vue";
import { ObjectUtils } from "primevue/utils";

export default {
  name: "Settings",
  components: {
    ArgumentList,
    DropdownParameter,
    SelectIcon,
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
      let url = this.$store.state.baseURL + "config/integrations";
      this.settings.infos.forEach((info) => {
        fetch(url + "/" + info.name, {
          method: "POST",
          body: JSON.stringify(this.settings.settings[info.name]),
          headers: new Headers({
            "Content-Type": "application/json",
            Authorization: `Basic ${btoa(
              `admin:${this.$store.state.password}`
            )}`,
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
      let url = this.$store.state.baseURL + "config/integrations";
      fetch(url, {
        method: "GET",
        mode: "cors",
        headers: new Headers({
          "Content-Type": "application/json",
          Authorization: `Basic ${btoa(`admin:${this.$store.state.password}`)}`,
        }),
      })
        .then((res) => res.json())
        .then((data) => {
          this.settings = data;
          console.log(JSON.stringify(data));
        })
        .catch((err) => console.log(err.message));
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

<style>
.w-page {
  width: 600px;
  height: 300px;
  text-align: left;
}
</style>