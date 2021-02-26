<template>
  <div class="layout-footer">
    <div class="layout-footer-left">
      <span
        >ReCoS Admin - Remote Command Service - (C) 2021 Wilfried Klaas</span
      >
    </div>
  </div>
</template>

<style scoped>
.layout-footer {
  position: absolute;
  bottom: 0;
  width: 100%;
  font-size: 1rem;
  padding: 2rem;
  background-color: var(--surface-a);
  display: flex;
  align-items: center;
  justify-content: space-between;

  a {
    color: var(--text-color);
    font-weight: 600;
  }

  .layout-footer-right {
    a {
      i {
        color: var(--text-secondary-color);
        font-size: 1.5rem;
      }
    }
  }
}
</style>