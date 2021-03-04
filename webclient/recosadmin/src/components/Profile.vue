<template>
  <TabView>
    <TabPanel header="Actions"> 
      <Actions :profile="profile"></Actions> </TabPanel>
    <TabPanel header="Pages">
      <Pages :profile="profile"></Pages>
    </TabPanel>
  </TabView>
</template>

<script>
import Actions from "./Actions.vue";
import Pages  from "./Pages.vue";
export default {
  name: "Profile",
  components: {
    Actions,
    Pages,
  },
  props: {
    profile: {},
  },

}
</script>