<template>
  <div class="footer">
    ReCoS Admin V{{ this.$store.state.packageVersion }} - Remote Command Service - (C) 2021 Wilfried Klaas
  </div>
</template>

<style>
.footer {
  position: fixed;
  left: 0;
  bottom: 0;
  width: 100%;
  color: rgba(255, 255, 255, 255);
  text-align: center;
}
</style>