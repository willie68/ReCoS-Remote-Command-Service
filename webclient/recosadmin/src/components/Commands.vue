<template>
  <Panel header="Commands" class="commands-panel-custom"></Panel>
  <Splitter style="height: 300px">
    <SplitterPanel :size="20">
      <Panel header=" " class="commands-panel-custom">
        <template #icons>
          <Button
            class="p-panel-header-icon p-link p-mr-2 p-mt-0 p-mb-0 p-pt-0 p-pb-0"
            icon="pi pi-plus"
            @click="addCommand"
          />
          <Button
            class="p-panel-header-icon p-link p-mr-2 p-mt-0 p-mb-0 p-pt-0 p-pb-0"
            icon="pi pi-arrow-up"
            @click="moveUp"
          />
          <Button
            class="p-panel-header-icon p-link p-mr-2 p-mt-0 p-mb-0 p-pt-0 p-pb-0"
            icon="pi pi-arrow-down"
            @click="moveDown"
          />
          <Button
            class="p-panel-header-icon p-link p-mr-2 p-mt-0 p-mb-0 p-pt-0 p-pb-0"
            icon="pi pi-trash"
            @click="deleteConfirm"
          />
        </template>
        <Listbox
          v-model="activeCommand"
          :options="activeAction.commands"
          optionLabel="name"
          listStyle="max-height:240px"
        >
        </Listbox>
      </Panel>
    </SplitterPanel>
    <SplitterPanel :size="80">
      <Command :command="activeCommand" v-on:change="changeCommand" />
    </SplitterPanel>
  </Splitter>
  <AddName
    :visible="addCmdDialog"
    v-model="newCmdName"
    :excludeList="cmdNames"
    v-on:save="saveNewCommand($event)"
    v-on:cancel="this.addCmdDialog = false"
    ><template #sourceHeader>New Command</template></AddName
  >
</template>

<script>
import Command from "./Command.vue";
import AddName from "./AddName.vue";
import { ObjectUtils } from "primevue/utils";

export default {
  name: "Commands",
  components: {
    Command,
    AddName,
  },
  props: {
    action: {},
  },
  data() {
    return {
      activeAction: { name: "" },
      activeCommand: {},
      addCmdDialog: false,
      newCmdName: null,
      cmdNames: [],
    };
  },
  watch: {
    action(action) {
      if (action) {
        console.log("Commands: action: " + JSON.stringify(action));
        this.activeAction = action;
        if (action.commands && action.commands.length > 0) {
          this.activeCommand = action.commands[0];
        }
      } else {
        this.activeAction = { name: "", commands: [] };
      }
    },
  },
  methods: {
    changeCommand(command) {
      console.log("Commands: command changed:" + command.name);
      //      console.log(JSON.stringify(this.action));
    },
    addCommand() {
      this.cmdNames = [];
      if (this.action.commands) {
        this.action.commands.forEach((element) => {
          this.cmdNames.push(element.name);
        });
      }
      this.addCmdDialog = true;
    },
    saveNewCommand(value) {
      console.log("Commands: add new command: " + value);
      this.addCmdDialog = false;
      let newCommand = {
        name: value,
        title: value,
        parameters: new Map(),
      };
      if (!this.activeAction.commands) {
        this.activeAction.commands = [];
      }
      this.activeAction.commands.push(newCommand);
      this.activeCommand = newCommand;
    },
    deleteConfirm() {
      if (this.activeCommand) {
        console.log("Commands: delete confirm pressed");
        this.$confirm.require({
          message:
            "Deleting command: " +
            this.activeCommand.name +
            ". Are you sure you want to proceed?",
          header: "Confirmation",
          icon: "pi pi-exclamation-triangle",
          accept: () => {
            this.deleteCommand();
          },
          reject: () => {
            //callback to execute when user rejects the action
          },
        });
      }
    },
    deleteCommand() {
      if (this.activeCommand) {
        console.log("Commands: delete command " + this.activeCommand.name);
        let index = ObjectUtils.findIndexInList(
          this.activeCommand,
          this.activeAction.commands
        );
        this.activeAction.commands.splice(index, 1);
        if (this.activeAction.commands.length > 0) {
          this.activeCommand = this.activeAction.commands[0];
        } else {
          this.activeCommand = null;
        }
      }
    },
    moveUp() {
      if (this.activeCommand) {
        let a = this.activeAction.commands;
        let index = ObjectUtils.findIndexInList(this.activeCommand, a);
        if (index != 0) {
          var b = a[index];
          a[index] = a[index - 1];
          a[index - 1] = b;
        }
      }
    },
    moveDown() {
      if (this.activeCommand) {
        let a = this.activeAction.commands;
        let index = ObjectUtils.findIndexInList(this.activeCommand, a);
        if (index < a.length - 1) {
          var b = a[index];
          a[index] = a[index + 1];
          a[index + 1] = b;
        }
      }
    },
  },
  mounted() {
    if (this.action) {
      this.activeAction = this.action;
      if (this.action.commands && this.action.commands.length > 0) {
        this.activeCommand = this.action.commands[0];
      }
    } else {
      this.activeAction = { name: "", commands: [] };
    }
  },
};
</script>

<style>
.commands-panel-custom .p-panel-header {
  margin: 0px;
  padding: 2px !important;
}
</style>