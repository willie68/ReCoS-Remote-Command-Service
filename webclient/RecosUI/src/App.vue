<template>
  <div class="display" ref="display">
    <AppConfig
      v-on:pageChanged="saveNewPage($event)"
      v-on:profileChanged="saveNewProfile($event)"
    />
    <div
      class="row"
      v-for="(row, x) in cellactions"
      :key="x"
      :style="{ height: cellHeight + 'px' }"
    >
      <div
        class="col"
        v-for="(col, y) in cellactions[x]"
        :key="y"
        :style="{ width: cellWidth + 'px' }"
      >
        <Action
          :title="cellactions[x][y].title"
          :actionUrl="actionURL"
          :actionHeight="actionHeight"
          :actionWidth="actionWidth"
          :profile="profileName"
          :actionName="cellactions[x][y].name"
          :icon="cellactions[x][y].icon"
          :ref="cellactions[x][y].name"
          :fontsize="cellactions[x][y].fontsize"
          :fontcolor="cellactions[x][y].fontcolor"
          :outlined="cellactions[x][y].outlined"
          :baseURL="this.baseURL"
          :actionType="cellactions[x][y].type"
          v-if="
            cellactions[x][y].type == 'SINGLE' ||
            cellactions[x][y].type == 'MULTI'
          "
          :action="cellactions[x][y]"
        />
        <Display
          :title="cellactions[x][y].title"
          :actionHeight="actionHeight"
          :actionWidth="actionWidth"
          :profile="profileName"
          :actionName="cellactions[x][y].name"
          :icon="cellactions[x][y].icon"
          :ref="cellactions[x][y].name"
          :fontsize="cellactions[x][y].fontsize"
          :fontcolor="cellactions[x][y].fontcolor"
          :outlined="cellactions[x][y].outlined"
          :baseURL="this.baseURL"
          v-if="cellactions[x][y].type == 'DISPLAY'"
        />
        <None
          :actionHeight="actionHeight"
          :actionWidth="actionWidth"
          v-if="cellactions[x][y].type == 'NONE'"
        />
      </div>
    </div>
  </div>
</template>

<script>
import AppConfig from "./components/AppConfig.vue";
import Action from "./components/Action.vue";
import Display from "./components/Display.vue";
import None from "./components/None.vue";

export default {
  name: "App",
  components: {
    AppConfig,
    Action,
    Display,
    None,
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
      title: "ReCoS",
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
      cellactions: [[]],
      cellWidth: 20,
      cellHeight: 20,
      actionWidth: 16,
      actionHeight: 16,
    };
  },
  computed: {
    newPageName: {
      get: function () {
        return this.pageName;
      },
      set: function (newPageName) {
        this.changePage(newPageName);
      },
    },
    toolbarPages: {
      get: function () {
        var a = [];
        if (this.activeProfile) {
          if (this.activeProfile.pages) {
            this.activeProfile.pages.forEach((page) => {
              if (page.toolbar != "hide") {
                a.push(page);
              }
            });
          }
        }
        return a;
      },
    },
  },
  mounted() {
    document.title = "ReCoS Client";  
    this.servicePort = 9280;
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
    this.connectWS();
  },
  methods: {
    helpHelp() {
      window
        .open(
          "https://raw.githubusercontent.com/willie68/ReCoS-Remote-Command-Service/master/documentation/README.pdf",
          "_blank"
        )
        .focus();
    },
    connectWS() {
      console.log("Starting connection to WebSocket Server");
      let that = this;
      if (this.connection) {
        this.connection.close(1000, "Work complete");
        this.connection = undefined;
      }
      this.connection = new WebSocket(
        "ws://" + window.location.hostname + ":" + this.servicePort + "/ws"
      );

      this.connection.onmessage = function (event) {
        // create a JSON object
        var jsonObject = JSON.parse(event.data);
        if (jsonObject.profile == that.profileName) {
          if (jsonObject.action) {
            if (that.$refs[jsonObject.action]) {
              console.log("Event:", event.data);
              // console.log("found action");
              if (jsonObject.command == "sendmessage") {
                alert("Message from ReCoS Service: \r\n" + jsonObject.text);
              }
              that.$refs[jsonObject.action].saveImg = jsonObject.imageurl;
              that.$refs[jsonObject.action].saveTitle = jsonObject.title;
              that.$refs[jsonObject.action].saveText = jsonObject.text;
              return;
            }
          }
          if (jsonObject.page) {
            // console.log("change page ", jsonObject.page);
            that.newPageName = jsonObject.page;
            return;
          }
          // console.log("action: ", jsonObject.action);
        }
      };

      this.connection.onopen = function (event) {
        // console.log(event);
        // console.log("Successfully connected to the websocket server...");
        var message = { profile: that.profileName, command: "change" };
        that.connection.send(JSON.stringify(message));
      };

      this.connection.onclose = function (event) {
        // console.log(event);
        // console.log("Connection closed to the websocket server...");
        if (that.connection) {
          that.connection.close(1000, "Work complete");
          that.connection = undefined;
        }
        setTimeout(() => that.connectWS(), 2000);
      };
    },
    toggleModal() {
      this.showModal = !this.showModal;
    },
    saveNewProfile(profile) {
      console.log("profile changed: " + JSON.stringify(profile));
      this.activeProfile = profile;
      this.profileName = profile.name;
      var message = { profile: this.profileName, command: "change" };
      this.connection.send(JSON.stringify(message));
      this.activePage = this.activeProfile.pages[0];
      this.changePage(this.activePage.name);
    },
    changeProfile() {
      fetch(this.showURL + "/" + this.profileName)
        .then((res) => res.json())
        .then((data) => {
          var message = { profile: this.profileName, command: "change" };
          this.connection.send(JSON.stringify(message));
          this.activeProfile = data;
          this.activePage = this.activeProfile.pages[0];
          this.changePage(this.activePage.name);
        })
        .catch((err) => console.log(err.message));
    },
    saveNewPage(page) {
      console.log("page changed: " + JSON.stringify(page));
      this.changePage(page.name);
    },
    changePage(pageName) {
      // console.log("change page to: ", pageName);
      this.pageName = pageName;
      this.activeProfile.pages.forEach((page) => {
        if (pageName == page.name) {
          this.activePage = page;
        }
      });
      // console.log("adding actions to page: " + this.activePage.name);
      // console.log("cell count: " + this.activePage.cells.length);
      this.cellactions = new Array(this.activePage.rows);
      for (let x = 0; x < this.activePage.rows; x++) {
        this.cellactions[x] = new Array(this.activePage.columns);
        for (let y = 0; y < this.activePage.columns; y++) {
          var action = undefined;
          let index = x * this.activePage.columns + y;
          if (index < this.activePage.cells.length) {
            let actionName = this.activePage.cells[index];
            if (actionName) {
              this.activeProfile.actions.forEach((profileAction, z) => {
                if (profileAction.name == actionName) {
                  action = profileAction;
                }
              });
              if (action) {
                // console.log("adding action (" + x + "," + y + ") " + index + ":" + action.name );
                this.cellactions[x][y] = action;
              } else if (actionName == "none") {
                this.cellactions[x][y] = {
                  type: "DISPLAY",
                  title: "",
                  name: "none",
                };
              } else {
                // console.log("missing action (" + x + "," + y + ") " + index);
                this.cellactions[x][y] = {
                  type: "DISPLAY",
                  title: "Action not defined",
                  name: "none",
                  fontcolor: "#FF0000",
                };
              }
              continue;
            }
          }
          // console.log("adding none(" + x + "," + y + ") " + index);
          this.cellactions[x][y] = {
            type: "NONE",
          };
        }

        this.cellWidth =
          this.$refs.display.clientWidth / this.activePage.columns - 4;
        this.cellHeight =
          this.$refs.display.clientHeight / this.activePage.rows - 4;
        this.actionWidth = this.cellWidth - 4;
        this.actionHeight = this.cellHeight - 4;
      }
    },
    buildImageSrc(data) {
      if (!data) {
        console.log("no data");
        return "data:image/bmp;base64,Qk1CAAAAAAAAAD4AAAAoAAAAAQAAAAEAAAABAAEAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP///wCAAAAA";
      }
      if (data.startsWith("/")) {
        return this.baseURL + data;
      }
      if (data.startsWith("data:")) {
        return data;
      }
      return "assets/" + data;
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
  top: 0;
  bottom: 0;
  width: 100%;
  background: black;
}

.pagebuttons {
  height: 32px;
}

.row {
  display: flex;
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
