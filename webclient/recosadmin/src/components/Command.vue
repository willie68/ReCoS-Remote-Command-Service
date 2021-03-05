<template>
    <Panel header="Command" class="command-panel-custom">
      {{ activeCommand.name }}
    </Panel>
</template>

<script>

export default {
  name: "Command",
  components: {
  },
  props: {
    command: {},
  },
  data() {
    return {
      activeCommand: {},
    };
  },
 watch: {
    command(command) {
      this.activeCommand = command
    },
  },
};
</script>

<style>
.command-panel-custom .p-panel-header {
  margin: 0px;
  padding: 2px !important;
}
</style>