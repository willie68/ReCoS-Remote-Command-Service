<template>
  <InputText v-model="val"></InputText>
  <h6>{{ val }}</h6>
</template>

<script>
export default {
  name: "RecosAdmin",
  components: {
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
