<template>
  <div class="profile-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>个人信息</span>
        </div>
      </template>
      
      <el-form :model="profileForm" label-width="120px">
        <el-form-item label="头像">
          <el-avatar :size="80" :icon="User" />
          <el-button type="primary" link style="margin-left: 20px">
            更换头像
          </el-button>
        </el-form-item>
        
        <el-form-item label="用户名">
          <el-input v-model="profileForm.username" disabled />
        </el-form-item>
        
        <el-form-item label="邮箱">
          <el-input v-model="profileForm.email" />
        </el-form-item>
        
        <el-form-item label="昵称">
          <el-input v-model="profileForm.nickname" placeholder="请输入昵称" />
        </el-form-item>
        
        <el-form-item label="手机号">
          <el-input v-model="profileForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        
        <el-form-item label="创建时间">
          <span>{{ profileForm.createdAt }}</span>
        </el-form-item>
        
        <el-form-item label="最后登录">
          <span>{{ profileForm.lastLoginAt }}</span>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSave">保存修改</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()

const profileForm = reactive({
  username: '',
  email: '',
  nickname: '',
  phone: '',
  createdAt: '',
  lastLoginAt: ''
})

onMounted(() => {
  // 加载用户信息
  profileForm.username = userStore.username
  profileForm.email = userStore.email
  profileForm.nickname = userStore.userInfo.nickname || ''
  profileForm.phone = userStore.userInfo.phone || ''
  profileForm.createdAt = '2024-01-01 00:00:00'
  profileForm.lastLoginAt = new Date().toLocaleString()
})

const handleSave = () => {
  // TODO: 调用保存 API
  ElMessage.success('保存成功')
}

const handleReset = () => {
  profileForm.email = userStore.email
  profileForm.nickname = userStore.userInfo.nickname || ''
  profileForm.phone = userStore.userInfo.phone || ''
}
</script>

<style scoped lang="scss">
.profile-container {
  max-width: 800px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: bold;
}
</style>
