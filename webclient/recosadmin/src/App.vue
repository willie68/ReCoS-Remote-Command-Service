<template>
  <Page></Page>
  <h6>{{ val }}</h6>
  <AppFooter></AppFooter>
</template>

<script>
import Page from "./components/Page.vue";
import AppFooter from "./components/AppFooter.vue";
export default {
  data() {
    return {
      name: "RecosAdmin",
    };
  },
  components: {
    Page,
    AppFooter,
  },
  mounted() {
    let servicePort = this.$store.state.servicePort;
    let basepath =
      window.location.protocol +
      "//" +
      window.location.hostname +
      ":" +
      servicePort +
      "/api/v1";
    this.$store.commit("baseURL", basepath);
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
}
</style>
