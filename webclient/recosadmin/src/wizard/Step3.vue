<template>
  <div class="w-page">
    <span class="p-mb-2"
      >Where should i display your action?<br />
      (To create a new 5x3 page simply add a new name.)</span
    >
    <div class="p-field p-grid p-mb-2 p-mt-2">
      <label :for="page" class="p-col-12 p-mb-2 p-ml-2 p-md-2 p-mb-md-0"
        >Page</label
      >
      <div class="p-col-12 p-md-8">
        <Dropdown
          id="page"
          v-model="page"
          :options="pages"
          optionLabel="name"
          class="p-ml-1 dropdownwidth"
          :editable="true"
        />
        <Button icon="pi pi-plus" @click="addPage" v-tooltip="'add new page'" />
      </div>
    </div>
    <transition-group
      v-for="row of page.rows"
      :key="row"
      name="dynamic-box"
      tag="div"
      class="p-grid"
    >
      <div v-for="col of page.columns" :key="col" class="p-col">
        <div class="box">
          <Button
            :ref="'btn' + ((row - 1) * page.columns + (col - 1))"
            v-if="isSet((row - 1) * page.columns + (col - 1))"
            style="width=100%;"
            :label="page.cells[(row - 1) * page.columns + (col - 1)]"
          ></Button>
          <Button
            :ref="'btn' + ((row - 1) * page.columns + (col - 1))"
            @click="clickButton((row - 1) * page.columns + (col - 1))"
            v-if="isActual((row - 1) * page.columns + (col - 1))"
            class="p-button-danger"
            style="width=100%;"
            label="actual"
          ></Button>
          <Button
            :ref="'btn' + ((row - 1) * page.columns + (col - 1))"
            @click="clickButton((row - 1) * page.columns + (col - 1))"
            v-if="
              !isSet((row - 1) * page.columns + (col - 1)) &&
              !isActual((row - 1) * page.columns + (col - 1))
            "
            class="p-button-success"
            style="width=100%;"
            label="empty"
          ></Button>
        </div>
      </div>
    </transition-group>
  </div>
</template>

<script>
export default {
  name: "Step3",
  components: {},
  props: {
    value: {},
    profile: {},
  },
  emits: ["next", "value"],
  data() {
    return {
      localValue: {},
      localPage: {},
      cells: [],
      newPage: null,
    };
  },
  computed: {
    page: {
      get: function () {
        return this.localPage;
      },
      set: function (newPage) {
        this.localPage = newPage;
        this.localValue.page = newPage.name;
      },
    },
    pages: {
      get: function () {
        let pages = [];
        this.profile.pages.forEach((element) => {
          pages.push(element);
        });
        if (this.newPage) {
          pages.push(this.newPage);
        }
        return pages;
      },
    },
  },
  mounted() {
    this.localValue = this.value;
    this.page = this.profile.pages[0];
    console.log("Step3: mounted value: ", JSON.stringify(this.localValue));
  },
  methods: {
    addPage() {
      if (
        typeof this.localPage === "string" ||
        this.localPage instanceof String
      ) {
        let isPresent = this.profile.pages
          .map((elem) => elem.name)
          .map((elem) => elem.toLowerCase())
          .includes(this.localPage.toLowerCase());
        console.log("Step3: add page is present: ", isPresent);
        if (!isPresent) {
          let myname = this.localPage;
          this.newPage = { name: myname, rows: 3, columns: 5, cells: [] };
          this.localPage = this.newPage;
        }
      }
    },
    clickButton(index) {
      this.localValue.page = this.localPage.name;
      this.localValue.index = index;
    },
    isSet(index) {
      if (this.page) {
        let cell = this.page.cells[index];
        if (cell && cell != "none") {
          return true;
        }
      }
      return false;
    },
    isActual(index) {
      if (this.localValue.index >= 0) {
        if (this.localValue.index == index) {
          return true;
        }
      }
      return false;
    },
    onClick() {
      console.log(JSON.stringify(this.page));
      console.log(JSON.stringify(this.localValue));
      //      console.log(JSON.stringify(this.profile));
    },
  },
  watch: {
    value: {
      deep: true,
      handler(value) {
        this.localValue = value;
        console.log("Step3: watch value: ", JSON.stringify(this.localValue));
      },
    },
    localPage: {
      deep: true,
      handler() {
        this.localValue.index = -1;
      },
    },
  },
};
</script>