<template>
  <div class="p-fluid p-formgrid p-grid">
    <div class="p-field p-col">
      <label for="name">Name</label>
      <InputText id="name" type="text" v-model="activeCommand.name" />
    </div>
    <div class="p-field p-col">
      <label for="title">Title</label>
      <InputText id="title" type="text" v-model="activeCommand.title" />
    </div>
    <div class="p-field p-col">
      <label for="rows">Type</label>
      <Dropdown
        v-model="activeCommand.type"
        :options="cmdTypes"
        placeholder="select a type"
        optionLabel="name"
        optionValue="type"
        optionGroupLabel="label"
        optionGroupChildren="items"
        editable
        :filter="true"
        filterPlaceholder="Find a command"
        ><template #optiongroup="slotProps">
          <div class="p-d-flex p-ai-center country-item">
            <img
              :src="
                this.$baseURL +
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
    </div>
    <div class="p-field p-col">
      <label for="icon">Icon</label>
      <span class="p-input-icon-right">
        <InputText
          id="icon"
          v-model="activeCommand.icon"
          placeholder="select a icon"
        />
        <i
          class="pi pi-chevron-down"
          @click="
            iconDestination = '';
            selectIconDialog = true;
          "
        />
      </span>
    </div>
  </div>
  <div class="p-fluid">
    <div class="p-field p-grid">
      <label for="description" class="p-col-6 p-mb-2 p-ml-2 p-md-2 p-mb-md-0"
        >Description</label
      >
      <div class="p-col-12 p-md-9">
        <InputText
          id="description"
          type="text"
          v-model="activeCommand.description"
        />
      </div>
    </div>
  </div>
  <div v-show="activeCommandType.parameter">
    <div class="p-pb-3">
      Command parameter for type {{ activeCommandType.name }}
    </div>
    <div class="p-fluid">
      <div
        class="p-field p-grid"
        v-for="(param, x) in activeCommandType.parameter"
        :key="x"
      >
        <label
          :for="param.name"
          class="p-col-12 p-mb-2 p-ml-2 p-md-2 p-mb-md-0"
          >{{ param.name }}</label
        >
        <div class="p-col-12 p-md-8">
          <InputText
            v-if="
              param.type == 'string' && (!param.list || param.list.length == 0)
            "
            :id="param.name"
            type="text"
            v-tooltip="param.description"
            v-model="activeCommand.parameters[param.name]"
          />

          <Dropdown
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
            v-model="activeCommand.parameters[param.name]"
            :placeholder="param.unit"
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
            v-model="activeCommand.parameters[param.name]"
            :placeholder="param.unit"
            ><template #optiongroup="slotProps">
              <div class="p-d-flex p-ai-center country-item">
                <img
                  :src="
                    this.$baseURL +
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
            v-model="activeCommand.parameters[param.name]"
          />
          <Checkbox
            v-if="param.type == 'bool'"
            :id="param.name"
            v-tooltip="param.description"
            v-model="activeCommand.parameters[param.name]"
            :binary="true"
          />
          <ArgumentList
            v-if="param.type == '[]string'"
            :id="param.name"
            v-tooltip="param.description"
            v-model="activeCommand.parameters[param.name]"
            v-on:add="addArgument(param.name)"
            v-on:remove="removeArgument(param.name, $event)"
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
            v-model="activeCommand.parameters[param.name]"
          />
          <Calendar
            v-if="param.type == 'date'"
            :id="param.name"
            :inline="false"
            dateFormat="yy-mm-dd"
            :showIcon="true"
            v-model="activeCommand.parameters[param.name]"
          />
          <span v-if="param.type == 'icon'" class="p-input-icon-right">
            <InputText
              :id="param.name"
              v-model="activeCommand.parameters[param.name]"
              v-tooltip="param.description"
              placeholder="select a icon"
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
    </div>
  </div>
  <AddName
    :visible="addArgDialog"
    v-model="newArgName"
    v-on:save="saveNewArgument($event)"
    v-on:cancel="this.addArgDialog = false"
    ><template #sourceHeader>New Argument</template></AddName
  >
  <SelectIcon
    :visible="selectIconDialog"
    :iconlist="iconlist"
    @cancel="this.selectIconDialog = false"
    @save="this.saveIcon($event)"
    ><template #sourceHeader>Select Icon</template></SelectIcon
  >
</template>

<script>
import ArgumentList from "./ArgumentList.vue";
import AddName from "./AddName.vue";
import SelectIcon from "./SelectIcon.vue";
import { ObjectUtils } from "primevue/utils";

export default {
  name: "Command",
  components: {
    ArgumentList,
    AddName,
    SelectIcon,
  },
  props: {
    command: {},
  },
  emits: ["change"],
  computed: {
    cmdTypes: {
      get: function () {
        let cmdTypes = Array();
        for (let i = 0; i < this.commandtypes.length; i++) {
          let commandType = this.commandtypes[i];
          let cat = commandType.category;
          if (!cat || cat == "") {
            cat = "unknown";
          }

          let found = false;
          for (let x = 0; x < cmdTypes.length; x++) {
            if (cmdTypes[x].label == cat) {
              cmdTypes[x].items.push(commandType);
              found = true;
            }
          }
          if (!found) {
            let myType = {
              label: cat,
              items: Array(),
            };
            myType.items.push(commandType);
            cmdTypes.push(myType);
          }
        }
        return cmdTypes;
      },
    },
  },
  data() {
    return {
      activeCommand: {},
      iconlist: [],
      commandtypes: [],
      activeCommandType: {},
      addArgDialog: false,
      newArgName: "",
      activeParam: null,
      selectIconDialog: false,
      iconDestination: "",
    };
  },
  watch: {
    command(command) {
      if (command) {
        console.log("Command: changing command to " + command.name + ", " + JSON.stringify(command));
        this.activeCommand = command;
      } else {
        this.activeCommand = { name: "", parameters: [] };
      }
    },
    activeCommand: {
      deep: true,
      handler(newCommand) {
        if (newCommand && newCommand.name != "") {
          if (newCommand.type) {
            this.commandtypes.forEach((commandType) => {
              if (commandType.type === newCommand.type) {
                this.activeCommandType = commandType;
              }
            });
          }
          this.$emit("change", this.activeCommand);
        }
      },
    },
    activeCommandType: {
      deep: false,
      handler(newCommandType) {
        if (newCommandType) {
          newCommandType.parameter.forEach((parameter) => {
            if (!this.activeCommand.parameters[parameter.name]) {
              if (parameter.type == "string") {
                this.activeCommand.parameters[parameter.name] = "";
              }
              if (parameter.type == "[]string") {
                this.activeCommand.parameters[parameter.name] = [""];
              }
              if (parameter.type == "int") {
                this.activeCommand.parameters[parameter.name] = 0;
              }
              if (parameter.type == "bool") {
                this.activeCommand.parameters[parameter.name] = false;
              }
              if (parameter.type == "color") {
                this.activeCommand.parameters[parameter.name] = "";
              }
            }
          });
        }
      },
    },
  },
  created() {
    this.updateCommandTypes();
  },
  methods: {
    paramList(param) {
      if (!param.groupedlist) {
        let list = Array();
        var fieldName, filterValue, key, label, value;
        let filtered = false;
        if (param.filteredlist) {
          fieldName = param.filteredlist;
          filterValue = this.activeCommand.parameters[fieldName];
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
      console.log("Command: add argument.");
      this.addArgDialog = true;
    },
    saveNewArgument(data) {
      if (this.activeParam) {
        let paramArray = this.activeCommand.parameters[this.activeParam];
        if (!paramArray) {
          paramArray = [];
        }
        paramArray.push(data);
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
        this.activeCommand.parameters[param]
      );
      if (index >= 0) {
        this.activeCommand.parameters[param].splice(index, 1);
      }
    },
    updateCommandTypes() {
      let url = this.$baseURL + "config/commands";
      fetch(url)
        .then((res) => res.json())
        .then((data) => {
          //console.log(data);
          this.commandtypes = data;
        })
        .catch((err) => console.log(err.message));
    },
    saveIcon(icon) {
      console.log("Action: save icon: ", icon, this.iconDestination);
      if (this.iconDestination == "") {
        this.activeCommand.icon = icon;
      } else {
        this.activeCommand.parameters[this.iconDestination] = icon;
      }
      this.selectIconDialog = false;
    },
  },
  beforeUnmount() {},
  mounted() {
    this.iconlist = this.$iconlist;
    if (this.command) {
      console.log("Command: monuted: " + JSON.stringify(this.command));
      this.activeCommand = this.command;
      this.commandtypes.forEach((commandType) => {
        if (commandType.type === this.activeCommand.type) {
          this.activeCommandType = commandType;
        }
      });
    } else {
      this.activeCommand = { name: "" };
    }
  },
};
</script>

<style>
.command-panel-custom .p-panel-header {
  margin: 0px;
  padding: 2px !important;
}
</style>