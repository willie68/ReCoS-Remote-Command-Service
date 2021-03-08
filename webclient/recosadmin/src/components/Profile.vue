<template>
  <TabView class="tabview-custom">
    <TabPanel header="Actions" class="tabpanel-custom">
      <Actions :profile="profile"></Actions>
    </TabPanel>
    <TabPanel header="Pages">
      <Pages :profile="profile"></Pages>
    </TabPanel>
  </TabView>
</template>

<script>
import Actions from "./Actions.vue";
import Pages from "./Pages.vue";

export default {
  name: "Profile",
  components: {
    Actions,
    Pages,
  },
  props: {
    profile: {},
  },
};
</script>

<style>
.tabview-custom .p-tabview-panels {
  margin: 0;
  padding: 0 !important;
}
</style>