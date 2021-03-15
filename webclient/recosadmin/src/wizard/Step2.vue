<template>
  <div class="w-page">
    <span class="p-mb-2"
      >Set the parameters for {{ activeCommandType.name }}</span
    >
    <Button @click="onClick" icon="pi pi-eye"/>
    <div class="p-field p-grid p-mb-2 p-mt-2">
      <label :for="title" class="p-col-12 p-mb-2 p-ml-2 p-md-2 p-mb-md-0"
        >Title</label
      >
      <div class="p-col-12 p-md-8">
        <InputText
          :id="title"
          type="text"
          v-tooltip="'The Title of the action'"
          v-model="title"
        />
      </div>
    </div>
    <div class="p-field p-grid p-mb-2 p-mt-2">
      <label :for="icon" class="p-col-12 p-mb-2 p-ml-2 p-md-2 p-mb-md-0"
        >Icon</label
      >
      <div class="p-col-12 p-md-8">
        <InputText
          :id="icon"
          type="text"
          v-tooltip="'The Icon of the action'"
          v-model="icon"
        />
      </div>
    </div>
    <div
      class="p-field p-grid p-mb-2 p-mt-2"
      v-for="(param, x) in localParameter"
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
          v-model="localValue.parameters[param.name]"
        />
        <Dropdown
          v-if="param.type == 'string' && param.list.length > 0"
          :id="param.name"
          :options="param.list"
          v-tooltip="param.description"
          v-model="localValue.parameters[param.name]"
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
          v-model="localValue.parameters[param.name]"
        />
        <Checkbox
          v-if="param.type == 'bool'"
          :id="param.name"
          v-tooltip="param.description"
          v-model="localValue.parameters[param.name]"
          :binary="true"
        />
        <ArgumentList
          v-if="param.type == '[]string'"
          :id="param.name"
          v-tooltip="param.description"
          v-model="localValue.parameters[param.name]"
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
          v-model="localValue.parameters[param.name]"
        />
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
</template>

<script>
import ArgumentList from "./ArgumentList.vue";
import AddName from "../components/AddName.vue";
import Button from 'primevue/button/Button.vue';
//import { ObjectUtils } from "primevue/utils";

export default {
  name: "Step2",
  components: {
    ArgumentList,
    AddName,
    Button,
  },
  props: {
    value: {},
    profile: {},
    commandTypes: Array,
  },
  emits: ["next", "value"],
  data() {
    return {
      localValue: {},
      activeCommandType: { name: "", parameters: [] },
      addArgDialog: false,
      newArgName: null,
      activeParam: null,
    };
  },
  computed: {
    localParameter: {
      get: function () {
        let params = [];
        if (this.activeCommandType && this.activeCommandType.parameter) {
          this.activeCommandType.parameter.forEach((element) => {
            if (element.wizard && element.wizard == true) {
              params.push(element);
            }
          });
        }
        return params;
      },
      set: function (newParam) {
        console.log("Step2: new param", newParam);
      },
    },
    title: {
      get: function () {
        this.checkType()
        return this.localValue.title;
      },
      set: function (newTitle) {
        this.localValue.title = newTitle;
      },
    },
    icon: {
      get: function () {
        this.checkType()
        return this.localValue.icon;
      },
      set: function (newIcon) {
        this.localValue.icon = newIcon;
      },
    },
  },
  methods: {
    onClick() {
      console.log(JSON.stringify(this.localValue));
    },
    addArgument(param) { 
      this.activeParam = param; 
      console.log("Step2: add argument for", param); 
      this.addArgDialog = true; 
    }, 
    saveNewArgument(data) { 
      if (this.activeParam) { 
        console.log("Step2: sav argument for", this.activeParam, data); 
        if (!this.localValue.parameters[this.activeParam]) { 
          this.localValue.parameters[this.activeParam] = []; 
        } 
        this.localValue.parameters[this.activeParam].push(data); 
      } 
      this.addArgDialog = false; 
    }, 
    checkType() {
      if (!this.localValue) {
        this.localValue = this.value
      }
      if (!this.localValue.title) {
        this.localValue.title = "";
      }
      if (!this.localValue.icon) {
        this.localValue.icon = "";
      }
      if (!this.localValue.parameters) {
        this.localValue.parameters = new Map();
      }
    },
  },
  mounted() {
    console.log("Step2: mounted value: ", JSON.stringify(this.value));
    this.localValue = this.value;
    let type = this.localValue.type;
    if (this.commandTypes) {
      this.commandTypes.forEach((element) => {
        if (element.type == type) {
          this.activeCommandType = element;
        }
      });
    }
    this.checkType();
  },
  watch: {
    value: {
      deep: true,
      handler(value) {
        this.localValue = value;
        console.log("Step2: watch value: ", JSON.stringify(this.localValue));
        this.checkType();
      },
    },
  },
};
</script>