<template>
  <Dialog v-model:visible="dialogVisible" style="width: 600px; height: 400px">
    <template #header>
      <h3>
        <div class="p-orderlist-header" v-if="$slots.sourceHeader">
          <slot name="sourceHeader"></slot>
        </div>
      </h3>
    </template>
    <span class="p-input-icon-right">
      <InputText
        id="icon"
        v-model="searchTerm"
        placeholder="select a icon"
        @input="onFilterChange"
      />
      <i class="pi pi-search" />
    </span>
    <br />
    <span
      v-for="(param, x) in filteredIconlist"
      :key="x"
      class="icon-item"
      @click="select(param)"
    >
      <a href="#" @keydown="onOptionKeyDown($event, option)">
        <img
          class="p-mb-2 p-mr-2"
          :class="{ 'icon-selected': isSelected(param) }"
          :src="'assets/' + param"
          @click="select(param)"
          width="36"
          :title="param"
        />
      </a>
    </span>
    <template #footer>
      <Button
        label="Cancel"
        icon="pi pi-times"
        class="p-button-text"
        @click="cancel"
      />
      <Button label="Save" icon="pi pi-check" autofocus @click="save" />
    </template>
  </Dialog>
</template>

<script>
export default {
  name: "SelectIcon",
  components: {},
  props: {
    iconlist: {
      type: Array,
    },
    visible: Boolean,
  },
  emits: ["save", "cancel"],
  computed: {
    filteredIconlist: {
      get: function () {
        let ficons = Array();
        if (this.searchTerm == "") {
          return this.iconlist;
        }
        this.iconlist.forEach((icon) => {
          if (icon.toLowerCase().includes(this.searchTerm.toLowerCase())) {
            ficons.push(icon);
          }
        });
        return ficons;
      },
    },
  },
  data() {
    return {
      dialogVisible: false,
      icon: "",
      searchTerm: "",
    };
  },
  methods: {
    onFilterChange() {
      console.log("SelectIcon: search:", this.searchTerm);
    },
    cancel() {
      this.$emit("cancel");
    },
    save() {
      console.log("SelectIcon saved: " + this.icon);
      this.$emit("save", this.icon);
    },
    select(item) {
      console.log("SelectIcon: select new icon: ", item);
      this.icon = item;
    },
    isSelected(item) {
      return this.icon == item;
    },
    onOptionKeyDown(event, option) {
      let item = event.currentTarget;
      console.log("onOptionKeyDown", event, option);
      switch (event.which) {
        //down
        case 40:
          var nextItem = this.findNextItem(item);
          if (nextItem) {
            nextItem.focus();
          }

          event.preventDefault();
          break;

        //up
        case 38:
          var prevItem = this.findPrevItem(item);
          if (prevItem) {
            prevItem.focus();
          }

          event.preventDefault();
          break;

        //enter
        case 13:
          this.onOptionSelect(event, option);
          event.preventDefault();
          break;
      }
    },
  },
  watch: {
    visible(visible) {
      this.dialogVisible = visible;
    },
    modelValue(value) {
      this.icon = value;
    },
  },
};
</script>

<style>
.icon-item {
}

.icon-selected {
  border: 1px;
  background-color: #2a435e;
}
</style>