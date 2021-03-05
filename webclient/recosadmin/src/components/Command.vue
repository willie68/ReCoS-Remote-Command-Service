<template>
  <Panel :header="activeCommand.name" class="command-panel-custom">
    <div class="p-fluid p-formgrid p-grid">
      <div class="p-field p-col">
        <label for="name">Name</label>
        <InputText id="name" type="text" v-model="activeCommand.name" />
      </div>
      <div class="p-field p-col">
        <label for="title">Titel</label>
        <InputText id="title" type="text" v-model="activeCommand.titel" />
      </div>
      <div class="p-field p-col">
        <label for="rows">Type</label>
        <Dropdown
          v-model="activeCommand.type"
          :options="enumCommandTypes"
          placeholder="select a type"
          optionLabel="name"
          optionValue="type"
        />
      </div>
      <div class="p-field p-col">
        <label for="icon">Icon</label>
        <Dropdown
          id="icon"
          v-model="activeCommand.icon"
          :options="iconlist"
          placeholder="select a icon"
        >
          <template #option="slotProps">
            <div class="icon-item">
              <img :src="'assets/' + slotProps.option" />
              <div>{{ slotProps.option }}</div>
            </div>
          </template>
        </Dropdown>
      </div>
    </div>
    <div class="p-fluid">
      <div class="p-field p-grid">
        <label for="description" class="p-col-6 p-mb-2 p-md-2 p-mb-md-0"
          >Description</label
        >
        <div class="p-col-12 p-md-10">
          <InputText
            id="description"
            type="text"
            v-model="activeCommand.description"
          />
        </div>
      </div>
    </div>
  </Panel>
</template>

<script>
export default {
  name: "Command",
  components: {},
  props: {
    command: {},
  },
  data() {
    return {
      activeCommand: {},
      iconlist: [],
    };
  },
  watch: {
    command(command) {
      this.activeCommand = command;
    },
  },
  created() {
    let that = this;
    this.unsubscribe = this.$store.subscribe((mutation, state) => {
      if (mutation.type === "baseURL") {
        that.iconurl = state.baseURL + "/config/icons";
        fetch(that.iconurl)
          .then((res) => res.json())
          .then((data) => {
            //console.log(data);
            that.iconlist = data;
          })
          .catch((err) => console.log(err.message));
      }
    });
  },
  beforeUnmount() {
    this.unsubscribe();
  },
};
</script>

<style>
.command-panel-custom .p-panel-header {
  margin: 0px;
  padding: 2px !important;
}
</style>