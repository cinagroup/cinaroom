<template>
  <div class="page-container">
    <Header>
      <template #left>
        <span class="page-title">系统设置</span>
      </template>
    </Header>
    <div class="cina-card">
      <div class="cina-card__header">
        <span class="cina-card__title">基本设置</span>
      </div>
      <el-form label-width="120px" style="max-width: 600px">
        <el-form-item label="系统名称">
          <el-input model-value="CinaSeek 云端开发工作室" disabled />
        </el-form-item>
        <el-form-item label="允许注册">
          <el-switch model-value />
        </el-form-item>
        <el-form-item label="默认角色">
          <el-select model-value="user" disabled>
            <el-option label="普通用户" value="user" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary">保存设置</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import Header from '@/components/Header.vue'
</script>
