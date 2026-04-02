<template>
  <div class="page-container">
    <Header>
      <template #left>
        <span class="page-title">系统日志</span>
      </template>
    </Header>
    <div class="cina-card">
      <div class="cina-card__header">
        <span class="cina-card__title">操作日志</span>
        <el-button type="primary" size="small">导出</el-button>
      </div>
      <el-table :data="[]" style="width: 100%">
        <el-table-column prop="time" label="时间" />
        <el-table-column prop="user" label="操作人" />
        <el-table-column prop="action" label="操作" />
        <el-table-column prop="detail" label="详情" />
        <template #empty>
          <el-empty description="暂无日志数据" />
        </template>
      </el-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import Header from '@/components/Header.vue'
</script>
