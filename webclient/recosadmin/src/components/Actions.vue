<template>
  <Splitter :style="{ height: splitterHeight }">
    <SplitterPanel :size="20" style="height: 100%">
      <Panel header="Actions" class="actions-panel-custom">
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
          listStyle="height: 100%"
        >
        </Listbox>
      </Panel>
    </SplitterPanel>
    <SplitterPanel :size="80">
      <Action :action="activeAction" :profile="activeProfile"></Action>
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
  <ConfirmDialog></ConfirmDialog>
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
      splitterHeight: "600px",
      addActionDialog: false,
      actionNames: null,
      newActionName: null,
    };
  },
  methods: {
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
      if (profile.actions) {
        this.activeAction = profile.actions[0];
      } else {
        this.activeAction = {};
      }
    },
  },
};
</script>

<style>
.actions-panel-custom {
  height: 100%;
}

.actions-panel-custom .p-panel-header {
  margin: 0px;
  padding: 2px !important;
}

.actions-panel-custom .p-panel-content {
  margin: 0px;
  padding: 2px !important;
  height: 100%;
}
.p-toggleable-content {
  height: 100%;
}
</style>