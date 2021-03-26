<template>
  <Splitter class="actions-splitter no-border">
    <SplitterPanel :size="20" style="height: 100%">
      <Panel header="Actions" class="actions-panel-custom no-border">
        <template #icons>
          <button
            class="p-panel-header-icon p-link p-mr-2 p-mt-0 p-mb-0 p-pt-0 p-pb-0"
            @click="addAction"
          >
            <span class="pi pi-plus"></span>
          </button>
          <button
            class="p-panel-header-icon p-link p-mr-2 p-mt-0 p-mb-0 p-pt-0 p-pb-0"
            @click="deleteConfirm"
          >
            <span class="pi pi-trash"></span>
          </button>
        </template>
        <Listbox
          v-model="activeAction"
          :options="profile.actions"
          optionLabel="name"
          listStyle="max-height:500px"
          class="no-border"
          v-on:change="checkChange($event)"
        >
        </Listbox>
      </Panel>
    </SplitterPanel>
    <SplitterPanel class="no-border" :size="80">
      <Panel
        v-show="activeAction"
        :header="'Action: ' + activeActionName"
        class="actions-panel-custom no-border"
      >
        <Action :action="activeAction" :profile="activeProfile"></Action>
      </Panel>
    </SplitterPanel>
  </Splitter>
  <AddName
    :visible="addActionDialog"
    v-model="newActionName"
    :excludeList="actionNames"
    v-on:save="saveNewAction($event)"
    v-on:cancel="this.addActionDialog = false"
    ><template #sourceHeader>New Action</template></AddName
  >
</template>

<script>
import Action from "./Action.vue";
import AddName from "./AddName.vue";
import { ObjectUtils } from "primevue/utils";

export default {
  name: "Actions",
  components: {
    Action,
    AddName,
  },
  props: {
    profile: {},
  },
  data() {
    return {
      activeProfile: null,
      activeAction: {},
      activeActionName: "",
      addActionDialog: false,
      actionNames: null,
      newActionName: null,
    };
  },
  methods: {
    checkChange(event) {
      let action = event.value;
      console.log("Actions: changed", JSON.stringify(action));
      if (action) {
        this.activeActionName = action.name;
      } else {
        this.activeActionName = "";
      }
    },
    addAction() {
      this.activeProfile.actions.forEach((element) => {
        this.actionNames.push(element.name);
      });
      this.addActionDialog = true;
    },
    saveNewAction(value) {
      console.log("Actions: add new action: " + value);
      this.addActionDialog = false;
      let newAction = {
        name: value,
        type: "SINGLE",
        description: "",
      };
      this.activeProfile.actions.push(newAction);
      this.activeAction = newAction;
      this.actionNames.push(newAction.name);
      this.activeActionName = this.activeAction.name;
    },
    deleteConfirm() {
      if (this.activeAction) {
        this.$confirm.require({
          message:
            "Deleting action: " +
            this.activeAction.name +
            ". Are you sure you want to proceed?",
          header: "Confirmation",
          icon: "pi pi-exclamation-triangle",
          accept: () => {
            this.deleteAction();
          },
          reject: () => {
            //callback to execute when user rejects the action
          },
        });
      }
    },
    deleteAction() {
      if (this.activeAction) {
        console.log("Actions: delete action " + this.activeAction.name);
        let index = ObjectUtils.findIndexInList(
          this.activeAction,
          this.activeProfile.actions
        );
        this.activeProfile.actions.splice(index, 1);
        if (this.activeProfile.actions.length > 0) {
          this.activeAction = this.activeProfile.actions[0];
        } else {
          this.activeAction = null;
        }
      }
    },
  },
  watch: {
    profile(profile) {
      this.activeProfile = profile;
      this.actionNames = [];
      if (profile.actions && profile.actions.length > 0) {
        this.activeAction = profile.actions[0];
        this.activeActionName = this.activeAction.name;
      } else {
        this.activeAction = {};
      }
    },
  },
};
</script>

<style>
.actions-splitter {
  border-width: 0px;
  height: 500px;
}

.actions-panel-custom {
  height: 100%;
}

.actions-panel-custom .p-panel-header {
  margin: 0px;
  padding: 2px !important;
  height: 34px;
}

.actions-panel-custom .p-panel-content {
  margin: 0px;
  padding: 2px !important;
  height: 100%;
}
</style>