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
          <InputText
            id="icon"
            v-model="icon"
            placeholder="select a icon"
            class="fullwidth"
          />
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
          v-if="
            param.type == 'string' &&
            param.list.length > 0 &&
            !param.groupedlist
          "
          :id="param.name"
          :options="paramList(param)"
          optionLabel="label"
          optionValue="value"
          v-tooltip="param.description"
          v-model="localValue.parameters[param.name]"
          :placeholder="param.unit"
          scrollHeight="120px"
          class="fullwidth"
          @change="update()"
        />
        <Dropdown
          v-if="
            param.type == 'string' && param.list.length > 0 && param.groupedlist
          "
          :id="param.name"
          :options="paramList(param)"
          optionLabel="label"
          optionValue="value"
          optionGroupLabel="label"
          optionGroupChildren="items"
          v-tooltip="param.description"
          v-model="localValue.parameters[param.name]"
          :placeholder="param.unit"
          ><template #optiongroup="slotProps">
            <div class="p-d-flex p-ai-center country-item">
              <img
                :src="
                  this.$store.state.baseURL +
                  'config/icons/category/' +
                  slotProps.option.label
                "
                width="18"
                class="p-mr-2"
              />
              <div>
                <b>{{ slotProps.option.label }}</b>
              </div>
            </div>
          </template>
        </Dropdown>
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
          @change="update()"
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
          @change="update()"
        />
        <Calendar
          v-if="param.type == 'date'"
          :id="param.name"
          :inline="false"
          dateFormat="yy-mm-dd"
          :showIcon="true"
          v-model="localValue.parameters[param.name]"
          @change="update()"
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
    paramList(param) {
      if (!param.groupedlist) {
        let list = Array();
        var fieldName, filterValue, key, label, value;
        let filtered = false;
        if (param.filteredlist) {
          fieldName = param.filteredlist;
          filterValue = this.localValue.parameters[fieldName];
          filtered = true;
        }
        param.list.forEach((entry) => {
          label = entry;
          value = entry;
          if (filtered) {
            if (entry.indexOf(":") > 0) {
              key = entry.substring(0, entry.indexOf(":"));
              label = entry.substring(entry.indexOf(":") + 1);
              value = entry;
              console.log("filter", key, label, value);
              if (key != filterValue) {
                return;
              }
            }
          }
          let myValue = {
            label: label,
            value: value,
          };
          list.push(myValue);
        });
        return list;
      }
      if (param.groupedlist) {
        let list = Array();
        param.list.forEach((entry) => {
          var key, value;
          if (entry.indexOf(":") > 0) {
            key = entry.substring(0, entry.indexOf(":"));
            value = entry.substring(entry.indexOf(":") + 1);
          } else {
            key = "unknown";
            value = entry;
          }
          let found = false;
          for (let x = 0; x < list.length; x++) {
            if (list[x].label == key) {
              let myValue = {
                label: value,
                value: entry,
              };
              list[x].items.push(myValue);
              found = true;
            }
          }
          if (!found) {
            let myType = {
              label: key,
              items: Array(),
            };
            let myValue = {
              label: value,
              value: entry,
            };
            myType.items.push(myValue);
            list.push(myType);
          }
        });
        return list;
      }
    },
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
    update() {
      this.$forceUpdate();
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
          this.localValue.icon = this.activeCommandType.icon;
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