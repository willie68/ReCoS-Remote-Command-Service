<template>
  <div
    class="action"
    :class="{ sale: theme === 'sale' }"
    :style="{
      height: actionHeight + 'px',
      width: actionWidth + 'px',
      backgroundImage: 'url(' + imageUrl + ')',
    }"
    @click="actionClick"
  >
    <b>{{ text }} {{ profile }}:{{ actionName }}</b>
    <!--
        <img
      class="stateimage"
      src="@/assets/point_red.png"
      :width="actionWidth - 10"
      :height="actionHeight - 10"
    />
    -->
  </div>
</template>

<script>
export default {
  name: "Action",
  props: [
    "text",
    "actionUrl",
    "theme",
    "actionHeight",
    "actionWidth",
    "profile",
    "actionName",
  ],
  data() {
    return {
      // imageUrl: "assets/point_red.png",
    };
  },
  computed: {
    imageUrl() {
      console.log("actionName:" + this.actionName)
      return this.actionName ? "assets/point_gray.png" : "assets/point_red.png"
    }
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
        fetch(actionPostUrl, options)
          .then((res) => res.json())
          .then((data) => {
            this.items = data.profiles;
            this.profileName = data.profiles[0].name;
            this.readonly = false;
            this.changeProfile();
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
  background: rgb(21, 32, 160);
  border-radius: 10px;
  text-align: center;
  background-repeat: no-repeat;
  background-attachment: fixed;
  background-size: 100% 100%;
}
.stateimage {
  position: relative;
  top: 0px;
  left: 0px;
}
.action h1 {
  color: #03cfb4;
  border: none;
  padding: 0;
}
.action p {
  font-style: normal;
}
.action.sale {
  background: crimson;
  color: white;
}
.action.sale h1 {
  background: crimson;
  color: black;
}
</style>