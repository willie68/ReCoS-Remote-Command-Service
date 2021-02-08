<template>
  <div
    class="action"
    :class="{ noicon: icon === '', noaction: actionName === '' }"
    :style="{
      height: actionHeight + 'px',
      width: actionWidth + 'px',
      backgroundImage: 'url(' + imageUrl + ')',
    }"
    @click="actionClick"
  >
    <span @click="actionClick"
      ><b>{{ text }}</b
      ><br />
      {{ actionName }}<br />
    </span>
  </div>
</template>

<script>
export default {
  name: "Action",
  props: [
    "text",
    "actionUrl",
    "actionHeight",
    "actionWidth",
    "profile",
    "actionName",
    "icon",
  ],
  data() {
    return {
      // imageUrl: "assets/point_red.png",
      saveImg: "",
    };
  },
  computed: {
    imageUrl() {
      console.log("actionName:" + this.actionName);
      if (this.actionName) {
        if (this.saveImg) {
          return "assets/" + this.saveImg;
        }
        return this.icon ? "assets/" + this.icon : "";
      }
      return "";
    },
  },
  methods: {
    closeModal() {
      this.$emit("close");
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
          body: JSON.stringify(""),
          headers: {
            "Content-Type": "application/json",
          },
        };
        this.saveImg = "check_mark.png";
        fetch(actionPostUrl, options)
          .then((res) => res.json())
          .then((data) => {
            setTimeout(()  => this.saveImg = "", 3000);
          })
          .catch((err) => console.log(err.message));
      }
    },
  },
};
</script>

<style>
.action {
  padding: 0px;
  margin: 0px;
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