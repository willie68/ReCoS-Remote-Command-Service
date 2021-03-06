<template>
  <Panel header="Commands" class="commands-panel-custom">
    <Splitter style="height: 400px">
      <SplitterPanel :size="20">
        <Panel header=" " class="commands-panel-custom no-border">
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
            class="no-border"
            v-on:change="changeCommand($event)"
          >
          </Listbox>
        </Panel>
      </SplitterPanel>
      <SplitterPanel :size="80">
        <Panel
          v-show="activeCommand != null"
          :header="'Command: ' + activeCommandName"
          class="commands-panel-custom no-border"
        >
          <Command :command="activeCommand" />
        </Panel>
      </SplitterPanel>
    </Splitter>
  </Panel>
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
      activeCommandName: "",
      addCmdDialog: false,
      newCmdName: null,
      cmdNames: [],
    };
  },
  methods: {
    changeCommand(event) {
      let command = event.value;
      if (command) {
        console.log("Commands: command changed:" + command.name);
        this.activeCommandName = command.name;
      } else {
        this.activeCommandName = "";
      }
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
      this.$nextTick(function () {
        console.log(
          "Commands: set active commnad to: ",
          this.activeAction.commands[this.activeAction.commands.length - 1]
        );
        this.activeCommand = this.activeAction.commands[
          this.activeAction.commands.length - 1
        ];
      });
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
          try {
            this.activeCommand = null;
          } catch (err) {
            console.log(err)
          }
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
        this.activeCommandName = this.activeCommand.name;
      }
    } else {
      this.activeAction = { name: "", commands: [] };
    }
  },
  watch: {
    activeProfile: {
      deep: true,
      handler(newProfile) {
        console.log("app: changing profile " + newProfile.name);
        //console.log(JSON.stringify(newProfile));
        this.profileDirty = true;
      },
    },
    action: {
      deep: false,
      handler(action) {
        if (action) {
          if (this.activeAction != action) {
            console.log("Commands: action: changed: ", action);
            this.activeAction = action;
            if (action.commands && action.commands.length > 0) {
              this.activeCommand = action.commands[0];
            }
          } else {
            this.activeAction = { name: "", commands: [] };
          }
        }
      },
    },
  },
};
</script>

<style>
.commands-list {
  height: 350px;

}
.commands-panel-custom {
  height: 435px;
}

.commands-panel-custom .p-panel-content {
  border-width: 0px !important;
}

.commands-panel-custom .p-panel-header {
  margin: 0px;
  padding: 2px !important;
}

.commands-panel-custom .p-panel-content {
  margin: 0px;
  padding: 2px !important;
}

</style>