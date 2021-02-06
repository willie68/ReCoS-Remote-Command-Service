<template>
  <div
    class="action"
    :class="{ sale: theme === 'sale' }"
    :style="{ height: actionHeight + 'px', width: actionWidth + 'px' 
      , backgroundImage: 'url('+ imageUrl + ')'
    }"

  >
    <b>{{ text }} {{ actionWidth }} {{ actionHeight }}</b>
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
  props: ["text", "actionurl", "theme", "actionHeight", "actionWidth"],
  data() {
      return {
          imageUrl: "assets/point_red.png",
      };
  },
  methods: {
    closeModal() {
      this.$emit("close");
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