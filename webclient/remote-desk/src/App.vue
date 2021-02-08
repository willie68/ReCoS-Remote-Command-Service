<template>
  <b>{{ title }}</b>
  <label> Profile</label>
  <select v-model="profileName" :disabled="readonly" @change="changeProfile()">
    <option
      v-for="item in items"
      :value="item.name"
      :key="item.name"
      v-text="item.name"
      :title="item.description"
    ></option>
  </select>
  <div>
    <button
      v-for="page in activeProfile.pages"
      :value="page.name"
      :key="page.name"
      v-text="page.name"
      :title="page.description"
      @click="changePage(page.name)"
    ></button>
  </div>
  <p>{{ cellHeight }} x {{ cellWidth }}</p>
  <div class="display" ref="display">
    <div
      class="row"
      v-for="(row, x) in cells"
      :key="x"
      :style="{ height: cellHeight + 'px' }"
    >
      <div
        class="col"
        v-for="(col, y) in cells[x]"
        :key="y"
        :style="{ width: cellWidth + 'px' }"
      >
        <Action
          :text="cells[x][y]"
          :actionUrl="actionURL"
          :actionHeight="actionHeight"
          :actionWidth="actionWidth"
          :profile="profileName"
          :actionName="cells[x][y]"
          :icon="icons[x][y]"
        ></Action>
      </div>
    </div>
  </div>
</template>

<script>
import Action from "./components/Action.vue";

export default {
  name: "App",
  components: {
    Action,
  },
  data() {
    return {
      servicePort: 9280,
      baseURL:
        window.location.protocol +
        "//localhost:" +
        this.servicePort +
        "/api/v1/",
      showURL: this.baseURL + "/show",
      actionURL: this.baseURL + "/action",
      title: "remote commands",
      header: "Title me",
      text: "this is a text",
      showModal: false,
      readonly: true,
      profileName: "none",
      items: [
        {
          name: "none",
          description: "nothing found here",
        },
      ],
      activeProfile: {},
      pageName: "",
      activePage: {},
      cells: [[]],
      cellWidth: 20,
      cellHeight: 20,
      actionWidth: 16,
      actionHeight: 16,
    };
  },
  mounted() {
    console.log("service url:" + this.baseURL);
    this.baseURL =
      window.location.protocol +
      "//" +
      window.location.hostname +
      ":" +
      this.servicePort +
      "/api/v1";
    this.showURL = this.baseURL + "/show";
    this.actionURL = this.baseURL + "/action";

    console.log("service url:" + this.baseURL);
    console.log("ui url:" + this.showURL);
    console.log("action url:" + this.actionURL);

    fetch(this.showURL)
      .then((res) => res.json())
      .then((data) => {
        this.items = data.profiles;
        this.profileName = data.profiles[0].name;
        this.readonly = false;
        this.changeProfile();
      })
      .catch((err) => console.log(err.message));
  },
  methods: {
    toggleModal() {
      this.showModal = !this.showModal;
    },
    changeProfile() {
      fetch(this.showURL + "/" + this.profileName)
        .then((res) => res.json())
        .then((data) => {
          this.activeProfile = data;
          this.activePage = this.activeProfile.pages[0];
          this.changePage(this.activePage.name);
        })
        .catch((err) => console.log(err.message));
    },
    changePage(pageName) {
      this.pageName = pageName;
      this.activeProfile.pages.forEach((page) => {
        if (pageName == page.name) {
          this.activePage = page;
        }
      });
      console.log(this.activePage);
      this.cells = new Array(this.activePage.rows);
      this.icons = new Array(this.activePage.rows);
      for (let x = 0; x < this.activePage.rows; x++) {
        this.cells[x] = new Array(this.activePage.columns);
        this.icons[x] = new Array(this.activePage.columns);
        for (let y = 0; y < this.activePage.columns; y++) {
          let action = this.activeProfile.actions[x * this.activePage.rows + y];
          if (action) {
            this.cells[x][y] = action.name;
            this.icons[x][y] = action.icon;
          } else {
            this.cells[x][y] = "";
            this.icons[x][y] = "";
          }
        }

        this.cellWidth =
          this.$refs.display.clientWidth / this.activePage.columns - 4;
        this.cellHeight =
          this.$refs.display.clientHeight / this.activePage.rows - 4;
        this.actionWidth = this.cellWidth - 4;
        this.actionHeight = this.cellHeight - 4;
      }
    },
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: top;
  color: #71b8ff;
  background: black;
}
.display {
  position: absolute;
  top: 50px;
  bottom: 0;
  width: 100%;
  background: black;
}
.row {
  display: block;
}
.col {
  display: inline-block;
}
h1 {
  border-bottom: 1px solid #ddd;
  display: inline-block;
  padding-bottom: 10px;
}
</style>
