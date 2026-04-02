<template>
  <div class="page-header">
    <div class="header-left">
      <slot name="left">
        <span class="page-title">{{ pageTitle }}</span>
      </slot>
    </div>
    <div class="header-right">
      <slot name="right"></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const pageTitle = computed(() => {
  return (route.meta?.title as string) || ''
})
</script>

<style scoped lang="scss">
.page-header {
  background-color: var(--card-bg);
  border-radius: var(--card-radius);
  border: 1px solid var(--border-color-light);
  padding: 15px 20px;
  margin-bottom: $content-margin;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: var(--card-shadow);
  transition: background-color 0.3s, border-color 0.3s;

  .header-left {
    display: flex;
    align-items: center;
    gap: 12px;

    .page-title {
      font-size: 16px;
      font-weight: 600;
      color: var(--text-primary);
    }
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}
</style>
