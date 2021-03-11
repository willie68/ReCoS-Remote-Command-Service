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
            @click="toggle"
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
      <Action :action="activeAction" :profile="profile"></Action>
    </SplitterPanel>
  </Splitter>
  <AddName
    :visible="addActionDialog"
    v-model="newActionName"
    :excludeList="actionNames"
    v-on:save="saveNewAction($event)"
    v-on:cancel="this.dialogNameVisible = false"
  />
</template>

<script>
import Action from "./Action.vue";
import AddName from "./AddName.vue";

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
      activeAction: {},
      splitterHeight: "600px",
      addActionDialog: false,
      actionNames: null,
      newActionName: null,
    };
  },
  methods: {
    addAction() {
      this.addActionDialog = true;
    },
    saveNewAction(value) {
      console.log("Actions: add new action: " + value)
      this.addActionDialog = false
      newAction = {
        name: value,
        type: "SINGLE",
        description: "",
      }
      this.profile.actions.push(newAction)
    }
  },
  watch: {
    profile(profile) {
      this.actionNames = [];
      if (profile.actions) {
        this.activeAction = profile.actions[0];
        profile.actions.forEach((element) => {
          this.actionNames.push(element.name);
        });
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