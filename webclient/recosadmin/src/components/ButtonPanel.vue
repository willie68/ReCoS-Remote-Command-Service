<template>
  <ScrollPanel style="width: 98%; height: 400px" class="custom">
    <transition-group
      v-for="row of activePage.rows"
      :key="row"
      name="dynamic-box"
      tag="div"
      class="p-grid"
    >
      <div v-for="col of activePage.columns" :key="col" class="p-col">
        <div class="box">
          <Button
            :ref="'btn' + ((row - 1) * activePage.columns + (col - 1))"
            @click="clickButton((row - 1) * activePage.columns + (col - 1))"
            v-if="cellActions[(row - 1) * activePage.columns + (col - 1)]"
            style="width=100%"
            :label="
              cellActions[(row - 1) * activePage.columns + (col - 1)].name
            "
          ></Button>
          <Button
            :ref="'btn' + ((row - 1) * activePage.columns + (col - 1))"
            @click="clickButton((row - 1) * activePage.columns + (col - 1))"
            v-if="!cellActions[(row - 1) * activePage.columns + (col - 1)]"
            class="p-button-success"
            label="empty"
          ></Button>
        </div>
      </div>
    </transition-group>
  </ScrollPanel>
  <SelectAction
    :visible="dialogActionVisible"
    v-on:save="assignAction($event)"
    v-on:cancel="this.dialogActionVisible = false"
    :sourceValue="profile.actions"
    :selectByName="buttonActionSelected"
  ></SelectAction>
</template>

<script>
import SelectAction from "./SelectAction.vue";

export default {
  name: "ButtonPanel",
  components: {
    SelectAction,
  },
  props: {
    profile: {},
    actions: {},
    page: {},
  },
  data() {
    return {
      activePage: {},
      cellActions: [{ name: "name" }],
      dialogActionVisible: false,
      buttonActionSelected: "",
    };
  },
  methods: {
    clickButton(index) {
      console.log("button clicked: ", index, this.activePage.cells[index]);
      this.buttonActionSelected = this.activePage.cells[index];
      this.saveIndex = index;
      this.dialogActionVisible = true;
    },
    displayAllRefs() {
      console.log("refs");
      console.log(this.$refs);
    },
    assignAction(action) {
      this.activePage.cells[this.saveIndex] = action.name;
      this.updateCellActions();
      this.dialogActionVisible = false;
    },
    updateCellActions() {
      this.cellActions = [];
      if (this.activePage.cells) {
        if (this.activePage.cells.length > 0) {
          let cellcount = this.activePage.cells.length;
          for (var i = 0; i < cellcount; i++) {
            this.cellActions[i] = this.getAction(i);
          }
        }
      }
    },
    getAction(index) {
      if (index < this.activePage.cells.length) {
        let actionName = this.activePage.cells[index];
        var action = null;
        this.actions.forEach((element) => {
          if (actionName == element.name) {
            action = element;
          }
        });
      }
      return action;
    },
  },
  watch: {
    page(page) {
      if (page) {
        this.activePage = page;
        this.updateCellActions();
      }
    },
  },
};
</script>

<style>
.custom .p-scrollpanel-wrapper {
  border-right: 9px solid #f4f4f4;
}

.custom .p-scrollpanel-bar {
  background-color: #1976d2;
  opacity: 1;
  transition: background-color 0.3s;
}

.custom .p-scrollpanel-bar:hover {
  background-color: #02386e;
}

.custom .p-button {
  width: 100px;
  height: 100px;
}
</style>