<template>
  <div
    class="acdisplay"
    :class="{ noicon: icon === '', noaction: actionName === '' }"
    :style="{
      height: actionHeight + 'px',
      width: actionWidth + 'px',
      backgroundImage: 'url(' + imageUrl + ')',
    }"
  >
    <div>
      <p class="title" :style="textStyle">{{ mytitle }}</p>
      <p class="text" :style="textStyle">{{ mytext }}</p>
    </div>
  </div>
</template>

<script>
export default {
  name: "Display",
  props: [
    "title",
    "text",
    "actionHeight",
    "actionWidth",
    "profile",
    "actionName",
    "icon",
    "fontsize",
    "fontcolor",
    "outlined",
  ],
  data() {
    return {
      // imageUrl: "assets/point_red.png",
      saveImg: "",
      saveTitle: "",
      saveText: "",
    };
  },
  computed: {
    textStyle() {
      return {
        fontSize: this.fontsize ? this.fontsize + "px" : "14px",
        color: this.fontcolor ? this.fontcolor : "black",
        textShadow: this.outlined ? "-1px 1px 2px #fff, 1px 1px 2px #fff, 1px -1px 0 #fff, -1px -1px 0 #fff" : ""
      };
    },
    imageUrl() {
      console.log("actionName:" + this.actionName);
      if (this.actionName) {
        if (this.saveImg) {
          return this.buildImageSrc(this.saveImg);
        }
        return this.icon ? this.buildImageSrc(this.icon) : "";
      }
      return "";
    },
    mytitle() {
      console.log("actionName:" + this.actionName);
      if (this.actionName) {
        if (this.saveTitle) {
          return this.saveTitle;
        }
        return this.title;
      }
      return "";
    },
    mytext() {
      console.log("actionName:" + this.actionName);
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
        return data;
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
.acdisplay {
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
}

.acdisplay .title {
  color: #000;
  font-size: 16px;
  font-weight: bold;
}

.acdisplay .text {
  color: #000;
  font-size: 16px;
  font-weight: bold;
}
.acdisplay h1 {
  color: #03cfb4;
  font-style: italic;
  border: none;
  padding: 0;
}

.acdisplay p {
  font-style: normal;
}

.acdisplay.noicon {
  background: lightgray;
  background-repeat: no-repeat;
  background-size: contain;
  background-position: center;
}

.acdisplay.noaction {
  background: rgb(45, 45, 45);
  color: black;
  background-size: 100% 100%;
}
</style>