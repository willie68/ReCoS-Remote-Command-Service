<template>
  <Page></Page>
  <h6>{{ val }}</h6>
</template>

<script>
import Page from "./components/Page.vue";
export default {
  name: "RecosAdmin",
  components: {
    Page,
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
