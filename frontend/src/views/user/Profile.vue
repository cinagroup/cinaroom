<template>
  <div class="profile-container page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="header-title">个人信息</span>
          <el-tag v-if="userStore.roles.length" type="success" size="small">
            {{ userStore.roles.join(', ') }}
          </el-tag>
        </div>
      </template>

      <el-row :gutter="40">
        <!-- Avatar Section -->
        <el-col :span="6">
          <div class="avatar-section">
            <el-avatar :size="100" :src="userStore.avatar" :icon="UserFilled" />
            <el-upload
              :show-file-list="false"
              :before-upload="beforeAvatarUpload"
              :http-request="handleAvatarUpload"
              accept="image/*"
            >
              <el-button type="primary" size="small" class="avatar-btn">
                <el-icon><Camera /></el-icon>更换头像
              </el-button>
            </el-upload>
          </div>
        </el-col>

        <!-- Form Section -->
        <el-col :span="18">
          <el-form
            ref="formRef"
            :model="profileForm"
            :rules="rules"
            label-width="100px"
          >
            <el-form-item label="用户名">
              <el-input :model-value="profileForm.username" disabled>
                <template #prepend>
                  <el-icon><User /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="昵称" prop="nickname">
              <el-input v-model="profileForm.nickname" placeholder="设置一个昵称">
                <template #prepend>
                  <el-icon><EditPen /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="邮箱" prop="email">
              <el-input v-model="profileForm.email" placeholder="请输入邮箱">
                <template #prepend>
                  <el-icon><Message /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="手机号" prop="phone">
              <el-input v-model="profileForm.phone" placeholder="请输入手机号">
                <template #prepend>
                  <el-icon><Phone /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="注册时间">
              <el-text type="info">{{ profileForm.createdAt || '-' }}</el-text>
            </el-form-item>

            <el-form-item label="最后登录">
              <el-text type="info">{{ profileForm.lastLoginAt || '-' }}</el-text>
            </el-form-item>

            <el-form-item>
              <el-button type="primary" :loading="saving" @click="handleSave">
                <el-icon><Check /></el-icon>保存修改
              </el-button>
              <el-button @click="handleReset">
                <el-icon><RefreshLeft /></el-icon>重置
              </el-button>
            </el-form-item>
          </el-form>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { UserFilled, User, EditPen, Message, Phone, Camera, Check, RefreshLeft } from '@element-plus/icons-vue'
import type { FormInstance, FormRules, UploadRequestOptions } from 'element-plus'
import { uploadAvatar } from '@/api/auth'

const userStore = useUserStore()
const formRef = ref<FormInstance>()
const saving = ref(false)

const profileForm = reactive({
  username: '',
  email: '',
  nickname: '',
  phone: '',
  createdAt: '',
  lastLoginAt: ''
})

const rules = reactive<FormRules>({
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ]
})

const loadUserInfo = () => {
  const info = userStore.userInfo
  profileForm.username = info.username || ''
  profileForm.email = info.email || ''
  profileForm.nickname = info.nickname || ''
  profileForm.phone = info.phone || ''
  profileForm.createdAt = info.createdAt || '-'
  profileForm.lastLoginAt = info.lastLoginAt || '-'
}

const handleSave = async () => {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  saving.value = true
  try {
    await userStore.updateProfile({
      email: profileForm.email,
      nickname: profileForm.nickname,
      phone: profileForm.phone
    })
    ElMessage.success('个人信息已更新')
  } catch (error: any) {
    ElMessage.error(error?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

const handleReset = () => {
  loadUserInfo()
  ElMessage.info('已重置为原始数据')
}

const beforeAvatarUpload = (file: File) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2
  if (!isImage) {
    ElMessage.error('只能上传图片文件！')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB！')
    return false
  }
  return true
}

const handleAvatarUpload = async (options: UploadRequestOptions) => {
  try {
    const res = await uploadAvatar(options.file) as any
    const data = res.data || res
    userStore.setUserInfo({ ...userStore.userInfo, avatar: data.url })
    ElMessage.success('头像已更新')
  } catch (error: any) {
    ElMessage.error(error?.message || '头像上传失败')
  }
}

onMounted(() => {
  loadUserInfo()
  // Fetch latest user info from server
  userStore.fetchUserInfo().then(() => loadUserInfo())
})
</script>

<style scoped lang="scss">
.profile-container {
  max-width: 900px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;

  .header-title {
    font-size: 18px;
    font-weight: 600;
  }
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding-top: 20px;

  .avatar-btn {
    margin-top: 8px;
  }
}
</style>
