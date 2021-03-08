<template>
    <Panel header="Commands" class="commands-panel-custom"></Panel>
    <Splitter style="height: 300px">
      <SplitterPanel :size="20">
        <Panel header=" " class="commands-panel-custom">
          <template #icons>
            <Button
              class="p-panel-header-icon p-link p-mr-2 p-mt-0 p-mb-0 p-pt-0 p-pb-0"
              @click="toggle"
              icon="pi pi-plus"
            />
            <Button
              class="p-panel-header-icon p-link p-mr-2 p-mt-0 p-mb-0 p-pt-0 p-pb-0"
              @click="toggle"
              icon="pi pi-arrow-up"
            />
            <Button
              class="p-panel-header-icon p-link p-mr-2 p-mt-0 p-mb-0 p-pt-0 p-pb-0"
              @click="toggle"
              icon="pi pi-arrow-down"
            />
            <Button
              class="p-panel-header-icon p-link p-mr-2 p-mt-0 p-mb-0 p-pt-0 p-pb-0"
              @click="toggle"
              icon="pi pi-trash"
            />
          </template>
          <Listbox
            v-model="activeCommand"
            :options="action.commands"
            optionLabel="name"
            listStyle="max-height:240px"
          >
          </Listbox>
        </Panel>
      </SplitterPanel>
      <SplitterPanel :size="80">
        <Command :command="activeCommand" v-on:change="changeCommand"/>
      </SplitterPanel>
    </Splitter>
</template>

<script>
import Command from "./Command.vue";
export default {
  name: "Commands",
  components: {
    Command,
  },
  props: {
    action: {},
  },
  data() {
    return {
      activeCommand: {},
    };
  },
 watch: {
    action(action) {
      if (action.commands && (action.commands.length > 0)) {
        this.activeCommand = action.commands[0]
      }
    },
  },
  methods: {
    changeCommand(command) {
      console.log("command changed:" + command.name );
      console.log(JSON.stringify(this.action))
    }
  }
};
</script>

<style>
.commands-panel-custom .p-panel-header {
  margin: 0px;
  padding: 2px !important;
}
</style>