<template>
  <div class="p-fluid p-formgrid p-grid">
    <div class="p-field p-col">
      <label for="name">Name</label>
      <InputText id="name" type="text" v-model="activeAction.name" />
    </div>
    <div class="p-field p-col">
      <label for="types">Type</label>
      <Dropdown
        id="types"
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
      <span class="p-input-icon-right">
        <InputText
          id="icon"
          v-model="activeAction.icon"
          placeholder="select a icon"
        />
        <i class="pi pi-chevron-down" @click="selectIconDialog = true" />
      </span>
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
        id="fontcolor"
        v-model="activeAction.fontcolor"
        :inline="false"
        defaultColor="#FFFFFF"
      />
    </div>
    <div class="p-field p-col">
      <label for="outlined">Outlined</label><br />
      <Checkbox id="outlined" v-model="activeAction.outlined" :binary="true" />
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
  <SelectIcon
    :visible="selectIconDialog"
    :iconlist="iconlist"
    @cancel="this.selectIconDialog = false"
    @save="this.saveIcon($event)"
    ><template #sourceHeader>Select Icon</template></SelectIcon
  >
</template>

<script>
import Commands from "./Commands.vue";
import MultiAction from "./MultiAction.vue";
import SelectIcon from "./SelectIcon.vue";

export default {
  name: "Action",
  components: {
    Commands,
    MultiAction,
    SelectIcon,
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
      selectIconDialog: false,
    };
  },
  methods: {
    saveIcon(icon) {
      console.log("Action: save icon: " + icon);
      this.activeAction.icon = icon;
      this.selectIconDialog = false;
    },
  },
  watch: {
    action(action) {
      if (action) {
        this.activeAction = action;
      } else {
        this.activeAction = { name: "" };
        //this.$refs.actionpanel
      }
    },
    activeAction(activeAction) {
      console.log("Action: activeAction changed: " + activeAction.name);
    },
    profile(profile) {
      //      console.log("Action change profile to " + profile.name);
      this.activeProfile = profile;
    },
    newIcon(icon) {
      console.log("Action: newIcon changed: " + icon);
      this.newIcon = icon;
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
};
</script>

<style>
.action-panel-custom {
  margin: 0px;
  text-align: center;
  height: 100%;
}

.action-panel-custom .p-panel-header {
  margin: 0px;
  padding: 2px !important;
  text-align: center;
  height: 34px;
}

.action-panel-custom .p-panel-content {
  margin: 0px;
  height: 400px;
  border-width: 0px;
}
</style>