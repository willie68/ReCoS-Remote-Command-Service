<template>
  <div
    class="nonedisplay"
    :style="{
      height: actionHeight + 'px',
      width: actionWidth + 'px',
    }">
  </div>
</template>

<script>
export default {
  name: "None",
  props: [
  ],
  data() {
    return {
    };
  },
};
</script>

<style>
.nonedisplay {
  padding: 0px;
  margin: 0px;
  border: 10px;
  border-radius: 10px;
  text-align: center;
  color: black;
  background: black;
  background-repeat: no-repeat;
  background-size: contain;
  background-position: center;
  justify-content: center;
  align-items: center;
  display: flex;
  transition: background-image 1s linear;
}
</style>