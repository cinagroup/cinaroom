<template>
  <div class="page-container">
    <Header>
      <template #left>
        <span class="page-title">用户管理</span>
      </template>
    </Header>
    <div class="cina-card">
      <div class="cina-card__header">
        <span class="cina-card__title">用户列表</span>
        <el-button type="primary" size="small">刷新</el-button>
      </div>
      <el-table :data="[]" style="width: 100%">
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="role" label="角色" />
        <el-table-column prop="status" label="状态" />
        <el-table-column label="操作" width="200" />
        <template #empty>
          <el-empty description="暂无用户数据" />
        </template>
      </el-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import Header from '@/components/Header.vue'
</script>
