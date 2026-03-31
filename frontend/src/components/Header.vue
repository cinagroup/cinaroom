<template>
  <div class="header">
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
.header {
  background-color: #fff;
  border-radius: $border-radius;
  padding: 15px 20px;
  margin-bottom: $content-margin;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}
</style>
