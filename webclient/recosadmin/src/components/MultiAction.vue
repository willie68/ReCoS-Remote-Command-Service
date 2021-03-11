<template>
  <Panel header="Multiaction" class="multi-panel-custom">
    <ActionSelectionList
      :sourceValue="actionNames"
      v-model="activeAction.actions"
    >
      <template #sourceHeader> Availble Action </template>
      <template #targetHeader> Stages </template>
      <template #item="slotProps">
        <div class="p-caritem">
          <span class="p-caritem-vin">{{ slotProps.item }}</span>
        </div>
      </template>
    </ActionSelectionList>
  </Panel>
</template>

<script>
import ActionSelectionList from "./ActionSelectionList.vue";
export default {
  name: "MultiAction",
  components: {
    ActionSelectionList,
  },
  props: {
    action: {},
    profile: {},
  },
  data() {
    return {
      activeAction: { actions: [] },
      actionNames: null,
      selection: "",
    };
  },
  watch: {
    profile(profile) {
//      console.log("multi profile changed: " + profile.name);
      this.actionNames = [];
      profile.actions.forEach((element) => {
        this.actionNames.push(element.name);
      });
      this.actionlist = [this.actionNames, []];
    },
    action(action) {
//      console.log("multi action changed:" + action.name);
      this.activeAction = action;
    },
  },
  methods: {
    changeCommand(command) {
//      console.log("multi command changed:" + command.name);
//      console.log(JSON.stringify(this.action));
    },
  },
};
</script>

<style>
.multi-panel-custom .p-panel-header {
  margin: 0px;
  padding: 2px !important;
}
</style>