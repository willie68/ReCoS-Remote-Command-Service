<template>
  <Panel
    :header="activeCommand.name"
    class="command-panel-custom"
    v-if="activeCommand.name != ''"
  >
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
          :options="commandtypes"
          placeholder="select a type"
          optionLabel="name"
          optionValue="type"
          editable
        />
      </div>
      <div class="p-field p-col">
        <label for="icon">Icon</label>
        <Dropdown
          id="icon"
          v-model="activeCommand.icon"
          :options="iconlist"
          placeholder="select a icon"
        >
          <template #option="slotProps">
            <div class="icon-item">
              <img :src="'assets/' + slotProps.option" />
              <div>{{ slotProps.option }}</div>
            </div>
          </template>
        </Dropdown>
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
              v-if="param.type == 'string' && param.list.length == 0"
              :id="param.name"
              type="text"
              v-tooltip="param.description"
              v-model="activeCommand.parameters[param.name]"
            />
            <Dropdown
              v-if="param.type == 'string' && param.list.length > 0"
              :id="param.name"
              :options="param.list"
              v-tooltip="param.description"
              v-model="activeCommand.parameters[param.name]"
              :placeholder="param.unit"
            />
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
          </div>
        </div>
      </div>
    </div>
  </Panel>
  <AddName
    :visible="addArgDialog"
    v-model="newArgName"
    v-on:save="saveNewArgument($event)"
    v-on:cancel="this.addArgDialog = false"
    ><template #sourceHeader>New Argument</template></AddName
  >
</template>

<script>
import ArgumentList from "./ArgumentList.vue";
import AddName from "./AddName.vue";
import { ObjectUtils } from "primevue/utils";

export default {
  name: "Command",
  components: {
    ArgumentList,
    AddName,
  },
  props: {
    command: {},
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
    };
  },
  watch: {
    command(command) {
      if (command) {
        console.log("Command: changing command to " + command.name);
        this.activeCommand = command;
      } else {
        this.activeCommand = { name: "" };
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
  },
  created() {
    this.updateIcons();
    this.upadteCommandTypes();
    let that = this;
    this.unsubscribe = this.$store.subscribe((mutation) => {
      if (mutation.type === "baseURL") {
        that.updateIcons();
        that.upadteCommandTypes();
      }
    });
  },
  methods: {
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
      console.log("Command: remove argument.", param, JSON.stringify(data));
      this.$confirm.require({
        message:
          "Deleting argument: " +
          param + ":" + data +
          ". Are you sure you want to proceed?",
        header: "Confirmation",
        icon: "pi pi-exclamation-triangle",
        accept: () => {
          this.deleteCommand(param, data);
        },
        reject: () => {
          //callback to execute when user rejects the action
        },
      });
    },
    deleteCommand(param, data) {
        console.log("Command: delete argument ", param, data);
        let index = ObjectUtils.findIndexInList(
          this.data,
          this.activeCommand.parameters[param]
        );
        this.activeCommand.parameters[param].splice(index, 1);
    },
    updateIcons() {
      let iconurl = this.$store.state.baseURL + "/config/icons";
      fetch(iconurl)
        .then((res) => res.json())
        .then((data) => {
          //console.log(data);
          this.iconlist = data;
        })
        .catch((err) => console.log(err.message));
    },
    upadteCommandTypes() {
      let url = this.$store.state.baseURL + "/config/commands";
      fetch(url)
        .then((res) => res.json())
        .then((data) => {
          //console.log(data);
          this.commandtypes = data;
        })
        .catch((err) => console.log(err.message));
    },
  },
  beforeUnmount() {
    this.unsubscribe();
  },
  mounted() {
    if (this.command) {
      console.log("Commands: action: " + JSON.stringify(this.command));
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