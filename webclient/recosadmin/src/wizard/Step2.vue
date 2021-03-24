<template>
  <div class="w-page">
    <span class="p-mb-2"
      >Set the parameters for {{ activeCommandType.name }}</span
    >
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
          class="fullwidth"
        />
      </div>
    </div>
    <div class="p-field p-grid p-mb-2 p-mt-2">
      <label :for="icon" class="p-col-12 p-mb-2 p-ml-2 p-md-2 p-mb-md-0"
        >Icon</label
      >
      <div class="p-col-12 p-md-8">
        <span class="p-input-icon-right fullwidth">
          <InputText id="icon" v-model="icon" placeholder="select a icon" class="fullwidth"/>
          <i class="pi pi-chevron-down" @click="selectIconDialog = true" />
        </span>
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
          class="fullwidth"
        />
        <DropdownParameter
          v-if="param.type == 'string' && param.list.length > 0"
          :id="param.name"
          :options="param.list"
          v-tooltip="param.description"
          v-model="localValue.parameters[param.name]"
          :placeholder="param.unit"
          scrollHeight="120px"
          class="fullwidth"
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
          class="fullwidth"
        />
        <Checkbox
          v-if="param.type == 'bool'"
          :id="param.name"
          v-tooltip="param.description"
          v-model="localValue.parameters[param.name]"
          :binary="true"
          class="fullwidth"
        />
        <ArgumentList
          v-if="param.type == '[]string'"
          :id="param.name"
          v-tooltip="param.description"
          v-model="localValue.parameters[param.name]"
          v-on:add="addArgument(param.name)"
          v-on:remove="removeArgument(param.name, $event)"
          class="fullwidth"
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
          class="fullwidth"
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
  <SelectIcon
    :visible="selectIconDialog"
    :iconlist="iconlist"
    @cancel="this.selectIconDialog = false"
    @save="this.saveIcon($event)"
    ><template #sourceHeader>Select Icon</template></SelectIcon
  >
</template>

<script>
import ArgumentList from "./ArgumentList.vue";
import AddName from "../components/AddName.vue";
import DropdownParameter from "./DropdownParameter.vue";
import SelectIcon from "./../components/SelectIcon.vue";
import { ObjectUtils } from "primevue/utils";

export default {
  name: "Step2",
  components: {
    ArgumentList,
    AddName,
    DropdownParameter,
    SelectIcon,
  },
  props: {
    value: {},
    profile: {},
    commandTypes: Array,
    iconlist: Array,
  },
  emits: ["next", "value"],
  data() {
    return {
      localValue: {},
      activeCommandType: { name: "", parameters: [] },
      addArgDialog: false,
      newArgName: null,
      activeParam: null,
      selectIconDialog: false,
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
        this.checkType();
        return this.localValue.title;
      },
      set: function (newTitle) {
        this.localValue.title = newTitle;
      },
    },
    icon: {
      get: function () {
        this.checkType();
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
    removeArgument(param, data) {
      let value = data;
      if (Array.isArray(data)) {
        value = data[0];
      }
      this.$confirm.require({
        message:
          "Deleting argument: " +
          param +
          ":" +
          value +
          ". Are you sure you want to proceed?",
        header: "Confirmation",
        icon: "pi pi-exclamation-triangle",
        accept: () => {
          this.deleteArgument(param, value);
        },
        reject: () => {
          //callback to execute when user rejects the action
        },
      });
    },
    deleteArgument(param, value) {
      let index = ObjectUtils.findIndexInList(
        value,
        this.localValue.parameters[param]
      );
      if (index >= 0) {
        this.localValue.parameters[param].splice(index, 1);
      }
    },
    checkType() {
      if (!this.localValue) {
        this.localValue = this.value;
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
    saveIcon(icon) {
      console.log("Step2: save icon: " + icon);
      this.localValue.icon = icon;
      this.selectIconDialog = false;
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
        if (!this.localValue.icon) {
          this.localValue.icon = this.activeCommandType.icon
        }
      },
    },
  },
};
</script>

<style>
.fullwidth {
  width: 100%;
}
</style>