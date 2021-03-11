<template>
  <Panel :header="activeAction.name" class="action-panel-custom">
    <div class="p-fluid p-formgrid p-grid">
      <div class="p-field p-col">
        <label for="name">Name</label>
        <InputText id="name" type="text" v-model="activeAction.name" />
      </div>
      <div class="p-field p-col">
        <label for="rows">Type</label>
        <Dropdown
          v-model="activeAction.type"
          :options="enumActionTypes"
          placeholder="select a type"
          optionLabel="name"
          optionValue="type"
        />
      </div>
      <div class="p-field p-col">
        <label for="title">Title</label>
        <InputText id="title" type="text" v-model="activeAction.title" />
      </div>
      <div class="p-field p-col">
        <label for="icon">Icon</label>
        <Dropdown
          id="icon"
          v-model="activeAction.icon"
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
      <div class="p-field p-col">
        <label for="fontsize">Font size</label>
        <Dropdown
          id="fontsize"
          v-model="activeAction.fontsize"
          :options="enumFontSizes"
          placeholder="select a size"
          optionLabel="name"
          optionValue="value"
        />
      </div>
      <div class="p-field p-col">
        <label for="fontcolor">Font color</label>
        <ColorPicker
          v-model="activeAction.fontcolor"
          :inline="false"
          defaultColor="#FFFFFF"
        />
      </div>
      <div class="p-field p-col">
        <label for="outlined">Outlined</label><br />
        <Checkbox
          id="outlined"
          v-model="activeAction.outlined"
          :binary="true"
        />
      </div>
      <div class="p-field p-col">
        <label for="runone">Only run one</label><br />
        <Checkbox id="runone" v-model="activeAction.runone" :binary="true" />
      </div>
    </div>
    <div class="p-fluid">
      <div class="p-field p-grid">
        <label for="description" class="p-col-6 p-mb-2 p-md-2 p-mb-md-0"
          >Description</label
        >
        <div class="p-col-12 p-md-10">
          <InputText
            id="description"
            type="text"
            v-model="activeAction.description"
          />
        </div>
      </div>
    </div>
    <Commands :action="activeAction" v-if="activeAction.type != `MULTI`" />
    <MultiAction
      :action="activeAction"
      :profile="activeProfile"
      v-show="activeAction.type == `MULTI`"
    />
  </Panel>
</template>

<script>
import Commands from "./Commands.vue";
import MultiAction from "./MultiAction.vue";

export default {
  name: "Action",
  components: {
    Commands,
    MultiAction,
  },
  props: {
    action: {},
    profile: {},
  },
  data() {
    return {
      activeProfile: {},
      activeAction: { name: "" },
      enumActionTypes: [
        { name: "Single", type: "SINGLE" },
        { name: "Display", type: "DISPLAY" },
        { name: "Multi", type: "MULTI" },
      ],
      enumFontSizes: [
        { name: "8", value: 8 },
        { name: "10", value: 10 },
        { name: "12", value: 12 },
        { name: "14", value: 14 },
        { name: "16", value: 16 },
        { name: "18", value: 18 },
        { name: "20", value: 20 },
        { name: "22", value: 22 },
        { name: "24", value: 24 },
        { name: "26", value: 26 },
        { name: "28", value: 28 },
        { name: "32", value: 32 },
        { name: "36", value: 36 },
        { name: "40", value: 40 },
        { name: "44", value: 44 },
        { name: "48", value: 48 },
      ],
      iconlist: [],
      activeCommand: {},
    };
  },
  watch: {
    action(action) {
      this.activeAction = action;
    },
    activeAction(activeAction) {
      console.log("Action: activeAction changed: " + JSON.stringify(activeAction));
    },
    profile(profile) {
//      console.log("Action change profile to " + profile.name);
      this.activeProfile = profile;
    },
  },
  mounted() {},
  created() {
    let that = this;
    this.unsubscribe = this.$store.subscribe((mutation, state) => {
      if (mutation.type === "baseURL") {
        that.iconurl = state.baseURL + "/config/icons";
        fetch(that.iconurl)
          .then((res) => res.json())
          .then((data) => {
            //console.log(data);
            that.iconlist = data;
          })
          .catch((err) => console.log(err.message));
      }
    });
  },
  beforeUnmount() {
    this.unsubscribe();
  },
};
</script>

<style>
.action-panel-custom .p-panel-header {
  margin: 0px;
  padding: 2px !important;
  text-align: center;
  height: 34px;
}
</style>