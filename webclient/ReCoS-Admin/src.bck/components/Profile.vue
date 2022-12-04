<template>
  <div>
    <div class="p-fluid p-mt-2">
      <div class="p-field p-grid">
        <label for="description" class="p-col-2 p-mb-2 p-md-2 p-mb-0 p-ml-2"
          >Description</label
        >
        <div class="p-col-2 p-md-9">
          <InputText
            id="description"
            type="text"
            v-model="activeProfile.description"
          />
        </div>
      </div>
    </div>
    <hr />
    <TabView class="tabview-custom">
      <TabPanel header="Pages">
        <Pages :profile="activeProfile"></Pages>
      </TabPanel>
      <TabPanel header="Actions">
        <Actions :profile="activeProfile"></Actions>
      </TabPanel>
    </TabView>
  </div>
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
  data() {
    return {
      activeProfile: {},
      profileDirty: false,
    };
  },
  created() {
    this.emitter.on('insertAction', (action) => this.addAction(action))
  },
  methods: {
    addAction(action) {
      console.log("Profile: addAction: " + JSON.stringify(action))
      this.activeProfile.actions.push(action)
    }
  },
  watch: {
    profile(newProfile) {
      if (newProfile) {
        console.log("changing profile to " + newProfile.name);
        this.activeProfile = newProfile;
      }
    },
  },
};
</script>

<style>
.tabview-custom .p-tabview-panels {
  margin: 0 0 0 0;
  padding: 0 1em 0 0 !important;
}
</style>