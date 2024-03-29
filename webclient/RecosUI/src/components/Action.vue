<template>
  <div
    ref="image"
    class="action"
    :class="{ noicon: icon === '', noaction: actionName === '' }"
    :style="{
      height: actionHeight + 'px',
      width: actionWidth + 'px',
      backgroundImage: 'url(' + imageUrl + ')',
    }"
    @click="actionClick"
    @dblclick="actionDblClick"
  >
    <div>
      <p class="title" :style="textStyle">{{ mytitle }}</p>
      <p class="text" :style="textStyle">{{ mytext }}</p>
    </div>
  </div>
</template>

<script>
export default {
  name: "Action",
  props: [
    "title",
    "text",
    "actionUrl",
    "actionHeight",
    "actionWidth",
    "profile",
    "actionName",
    "icon",
    "fontsize",
    "fontcolor",
    "outlined",
    "baseURL",
    "actionType",
    "action",
  ],
  data() {
    return {
      // imageUrl: "assets/point_red.png",
      saveImg: "",
      saveTitle: "",
      saveText: "",
      myBaseUrl: "",
      timerID: null,
    };
  },
  computed: {
    textStyle() {
      return {
        fontSize: this.fontsize ? this.fontsize + "px" : "14px",
        color: this.fontcolor ? this.fontcolor : "black",
        textShadow: this.outlined
          ? "-1px 1px 2px #fff, 1px 1px 2px #fff, 1px -1px 0 #fff, -1px -1px 0 #fff"
          : "",
      };
    },
    imageUrl() {
      if (this.actionName) {
        if (this.saveImg) {
          return this.buildImageSrc(this.saveImg);
        }
        return this.action.icon ? this.buildImageSrc(this.action.icon) : "";
      }
      return "";
    },
    mytitle() {
      if (this.actionName) {
        if (this.saveTitle) {
          return this.saveTitle;
        }
        return this.action.title;
      }
      return "";
    },
    mytext() {
      if (this.actionName) {
        if (this.saveText) {
          return this.saveText;
        }
      }
      return "";
    },
  },
  methods: {
    closeModal() {
      this.$emit("close");
    },
    buildImageSrc(data) {
      if (data.startsWith("/")) {
        let imgWidth = this.actionWidth
        let imgHeight = this.actionHeight
        return (
          this.myBaseUrl +
          data +
          "?width=" +
          Math.floor(imgWidth) +
          "&height=" +
          Math.floor(imgHeight)
        );
      }
      if (data.startsWith("data:")) {
        return data;
      }
      return "assets/" + data;
    },
    actionClick() {
      console.log(
        "action " + this.profile + ":" + this.actionName + " clicked"
      );
      if (this.actionName) {
        var actionPostUrl =
          this.actionUrl + "/" + this.profile + "/" + this.actionName;
        var options = {
          method: "POST",
          body: JSON.stringify({
            profile: this.profile,
            action: this.actionName,
            page: this.page,
            command: "click",
          }),
          headers: {
            "Content-Type": "application/json",
          },
        };
        this.saveImg = "hourglass.svg";
        if (this.actionType != "MULTI") {
          console.log("set timeout");
          if (this.timerID) {
            clearTimeout(this.timerID);
            this.timerID = null;
          }
          this.timerID = setTimeout(() => (this.saveImg = ""), 20000);
        }
        let that = this;
        fetch(actionPostUrl, options)
          .then((res) => res.json())
          .then((data) => {
            console.log(that.actionType);
          })
          .catch((err) => console.log(err.message));
      }
    },
    actionDblClick() {
      console.log(
        "action " + this.profile + ":" + this.actionName + " clicked"
      );
      if (this.actionName) {
        var actionPostUrl =
          this.actionUrl + "/" + this.profile + "/" + this.actionName;
        var options = {
          method: "POST",
          body: JSON.stringify({
            profile: this.profile,
            action: this.actionName,
            page: this.page,
            command: "dblclick",
          }),
          headers: {
            "Content-Type": "application/json",
          },
        };
        if (this.actionType != "MULTI") {
          this.saveImg = "hourglass.svg";
          if (this.timerID) {
            clearTimeout(this.timerID);
            this.timerID = null;
          }
          this.timerID = setTimeout(() => (this.saveImg = ""), 20000);
        }
        fetch(actionPostUrl, options)
          .then((res) => res.json())
          .then((data) => {})
          .catch((err) => console.log(err.message));
      }
    },
  },
  mounted() {
    console.log("Action: mounted ", this.action.name);
    this.myBaseUrl = this.baseURL.substring(0, this.baseURL.indexOf("/", 8));
    this.saveImg = "";
  },
  unmounted() {
    console.log("Action: unmounted ", this.action.name);
  },
  watch: {
    action: {
      deep: true,
      handler(newAction) {
        console.log("Action: changing action " + JSON.stringify(newAction));
        this.saveImg = "";
        this.saveTitle = "";
        this.saveText = "";
      },
    },
    baseURL(baseURL) {
      this.myBaseUrl = baseURL.substring(0, baseURL.indexOf("/", 8));
    },
    saveImg: {
      deep: false,
      handler(newImg) {
        if (this.timerID) {
          console.log("stop timer");
          clearTimeout(this.timerID);
          this.timerID = null;
        }
      },
    },
  },
};
</script>

<style>
.action {
  padding: 0px;
  margin: 0px;
  border: 10px;
  border-radius: 10px;
  text-align: center;
  color: black;
  background: darkgray;
  background-repeat: no-repeat;
  background-size: contain;
  background-position: center;
  justify-content: center;
  align-items: center;
  display: flex;
  transition: background-image 1s ease-in-out;
}

.action .title {
  color: #000;
  font-size: 16px;
  font-weight: bold;
}

.action .text {
  color: #000;
  font-size: 16px;
  font-weight: bold;
}

.action h1 {
  color: #03cfb4;
  font-style: italic;
  border: none;
  padding: 0;
}

.action p {
  font-style: normal;
}

.action.noicon {
  background: lightgray;
  background-repeat: no-repeat;
  background-size: contain;
  background-position: center;
}

.action.noaction {
  background: rgb(45, 45, 45);
  color: black;
  background-size: 100% 100%;
}
</style>