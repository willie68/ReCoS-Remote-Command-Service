<template>
  <HelloWorld> </HelloWorld>
  <InputText v-model="val"></InputText>
  <h6>{{ val }}</h6>
</template>

<script>
import HelloWorld from "./components/HelloWorld.vue";
export default {
  name: "RecosAdmin",
  components: {
    HelloWorld,
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
