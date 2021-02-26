<template>
  <ScrollPanel style="width: 98%; height: 200px" class="custom">
    <div>rows {{ rows }} cols {{ columns }}</div>
    <transition-group
      v-for="row of rows"
      :key="row"
      name="dynamic-box"
      tag="div"
      class="p-grid"
    >
      <div v-for="col of columns" :key="col" class="p-col">
        <div class="box" v-badge="(row - 1) * columns + (col - 1)">
          <Button></Button>
        </div>
      </div>
    </transition-group>
  </ScrollPanel>
</template>

<script>
export default {
  name: "ButtonPanel",
  components: {},
  props: {
    rows: {},
    columns: {},
  },
};
</script>

<style>
.custom .p-scrollpanel-wrapper {
    border-right: 9px solid #f4f4f4;
}

.custom .p-scrollpanel-bar {
    background-color: #1976d2;
    opacity: 1;
    transition: background-color .3s;
}

.custom .p-scrollpanel-bar:hover {
    background-color: #02386e;
}
</style>