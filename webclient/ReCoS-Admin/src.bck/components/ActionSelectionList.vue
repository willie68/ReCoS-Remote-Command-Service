<template>
  <div class="p-orderlist p-component">
    <div class="p-orderlist-list-container">
      <div class="p-orderlist-header" v-if="$slots.sourceHeader">
        <slot name="sourceHeader"></slot>
      </div>
      <Listbox
        v-model="selectedSourceValue"
        :options="sourceValue"
        listStyle="height: 100%"
      >
      </Listbox>
    </div>
    <div class="p-orderlist-controls">
      <OLButton type="button" icon="pi pi-plus" @click="add"></OLButton>
      <OLButton type="button" icon="pi pi-minus" @click="remove"></OLButton>
      <OLButton
        type="button"
        icon="pi pi-angle-double-up"
        @click="moveTop"
      ></OLButton>
      <OLButton type="button" icon="pi pi-angle-up" @click="moveUp"></OLButton>
      <OLButton
        type="button"
        icon="pi pi-angle-down"
        @click="moveDown"
      ></OLButton>
      <OLButton
        type="button"
        icon="pi pi-angle-double-down"
        @click="moveBottom"
      ></OLButton>
    </div>
    <div class="p-orderlist-list-container">
      <div class="p-orderlist-header" v-if="$slots.targetHeader">
        <slot name="targetHeader"></slot>
      </div>
      <Listbox
        v-model="selectedDestination"
        :options="destinationList"
        optionLabel="value"
        dataKey="id"
        listStyle="height: 100%"
      >
      </Listbox>
    </div>
  </div>
</template>

<script>
import Button from "primevue/button";
import { ObjectUtils } from "primevue/utils";

export default {
  name: "ActionSelectionList",
  emits: ["update:modelValue"],
  props: {
    modelValue: {
      type: Array,
      default: null,
    },
    sourceValue: {
      type: Array,
      default: null,
    },
  },
  data() {
    return {
      selectedSourceValue: null,
      selectedDestination: null,
      selectedDestinationIndex: 0,
      destinationList: null,
      idcounter: 0,
    };
  },
  beforeUnmount() {},
  updated() {
    if (this.reorderDirection) {
      this.updateListScroll();
      this.reorderDirection = null;
    }
  },
  mounted() {},
  methods: {
    getItemKey(item, index) {
      return this.dataKey
        ? ObjectUtils.resolveFieldData(item, this.dataKey)
        : index;
    },
    newEntry(name) {
      let entry = {
        value: name,
        id: this.idcounter,
      };
      this.idcounter++;
      return entry;
    },
    add() {
      if (this.selectedSourceValue) {
        this.destinationList.push(this.newEntry(this.selectedSourceValue));
        this.destinationChanged();
      }
    },
    remove() {
      if (this.selectedDestination) {
        let index = ObjectUtils.findIndexInList(
          this.selectedDestination,
          this.destinationList
        );
        this.destinationList.splice(index, 1);
        this.destinationChanged();
      }
    },
    moveUp() {
      if (this.selectedDestination) {
        let index = ObjectUtils.findIndexInList(
          this.selectedDestination,
          this.destinationList
        );
        if (index != 0) {
          let a = this.destinationList;
          var b = a[index];
          a[index] = a[index - 1];
          a[index - 1] = b;
          this.destinationChanged();
        }
      }
    },
    moveTop() {
      if (this.selectedDestination) {
        let index = ObjectUtils.findIndexInList(
          this.selectedDestination,
          this.destinationList
        );
        if (index != 0) {
          let a = this.destinationList;
          var b = a[index];
          for (let x = index; x > 0; x--) {
            a[x] = a[x - 1];
          }
          a[0] = b;
          this.destinationChanged();
        }
      }
    },
    moveDown() {
      if (this.selectedDestination) {
        let a = this.destinationList;
        let index = ObjectUtils.findIndexInList(this.selectedDestination, a);
        if (index < a.length - 1) {
          var b = a[index];
          a[index] = a[index + 1];
          a[index + 1] = b;
          this.destinationChanged();
        }
      }
    },
    moveBottom() {
      if (this.selectedDestination) {
        let a = this.destinationList;
        let index = ObjectUtils.findIndexInList(this.selectedDestination, a);
        if (index != a.length - 1) {
          var b = a[index];
          for (let x = index; x < a.length - 1; x++) {
            a[x] = a[x + 1];
          }
          a[a.length - 1] = b;
          this.destinationChanged();
        }
      }
    },
    destinationChanged() {
      console.log(
        "ActionSelectionList modelValue changed" +
          JSON.stringify(this.destinationList)
      );
      let result = [];
      this.destinationList.forEach((element) => {
        result.push(element.value);
      });
      this.$emit("update:modelValue", result);
    },
  },
  computed: {},
  components: {
    OLButton: Button,
  },
  watch: {
    modelValue(modelvalue) {
      this.destinationList = [];
      if (modelvalue) {
        modelvalue.forEach((element) => {
          this.destinationList.push(this.newEntry(element));
        });
      }
    },
  },
};
</script>

<style>
</style>