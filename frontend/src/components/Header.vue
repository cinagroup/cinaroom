<template>
  <div class="header-bar">
    <div class="header-left">
      <slot name="left"></slot>
    </div>
    <div class="header-center">
      <slot name="center"></slot>
    </div>
    <div class="header-right">
      <slot name="right"></slot>
    </div>
  </div>
</template>

<script setup>
</script>

<style scoped lang="scss">
.header-bar {
  background-color: var(--card-bg);
  border-radius: var(--card-radius);
  border: 1px solid var(--border-color-light);
  padding: 15px 20px;
  margin-bottom: $content-margin;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: var(--shadow-sm);
  transition: background-color 0.3s, border-color 0.3s;
}
</style>
