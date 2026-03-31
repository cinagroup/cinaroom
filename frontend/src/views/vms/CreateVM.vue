<template>
  <div class="create-vm-container">
    <Header>
      <template #left>
        <span class="page-title">创建虚拟机</span>
      </template>
    </Header>
    
    <el-card>
      <el-form
        ref="formRef"
        :model="vmForm"
        :rules="rules"
        label-width="120px"
        style="max-width: 700px"
      >
        <el-form-item label="虚拟机名称" prop="name">
          <el-input v-model="vmForm.name" placeholder="请输入虚拟机名称" />
        </el-form-item>
        
        <el-form-item label="镜像" prop="image">
          <el-select v-model="vmForm.image" placeholder="请选择镜像" style="width: 100%">
            <el-option label="Ubuntu 22.04 LTS" value="ubuntu-22.04" />
            <el-option label="Ubuntu 20.04 LTS" value="ubuntu-20.04" />
            <el-option label="CentOS 7" value="centos-7" />
            <el-option label="CentOS Stream 9" value="centos-stream-9" />
            <el-option label="Debian 12" value="debian-12" />
            <el-option label="Debian 11" value="debian-11" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="CPU 核心数" prop="cpu">
          <el-input-number v-model="vmForm.cpu" :min="1" :max="8" :step="1" />
        </el-form-item>
        
        <el-form-item label="内存大小" prop="memory">
          <el-input-number v-model="vmForm.memory" :min="1" :max="16" :step="1" />
          <span style="margin-left: 10px">GB</span>
        </el-form-item>
        
        <el-form-item label="磁盘大小" prop="disk">
          <el-input-number v-model="vmForm.disk" :min="10" :max="500" :step="10" />
          <span style="margin-left: 10px">GB</span>
        </el-form-item>
        
        <el-form-item label="网络配置" prop="network">
          <el-radio-group v-model="vmForm.network">
            <el-radio label="nat">NAT</el-radio>
            <el-radio label="bridge">桥接</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-divider content-position="left">高级设置</el-divider>
        
        <el-form-item label="SSH 密钥">
          <el-input
            v-model="vmForm.sshKey"
            type="textarea"
            :rows="4"
            placeholder="粘贴 SSH 公钥（可选）"
          />
        </el-form-item>
        
        <el-form-item label="初始化脚本">
          <el-input
            v-model="vmForm.initScript"
            type="textarea"
            :rows="4"
            placeholder="Cloud-init 脚本（可选）"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleCreate">创建虚拟机</el-button>
          <el-button @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import Header from '@/components/Header.vue'

const router = useRouter()
const formRef = ref(null)
const loading = ref(false)

const vmForm = reactive({
  name: '',
  image: '',
  cpu: 2,
  memory: 4,
  disk: 50,
  network: 'nat',
  sshKey: '',
  initScript: ''
})

const rules = {
  name: [
    { required: true, message: '请输入虚拟机名称', trigger: 'blur' },
    { min: 2, max: 32, message: '名称长度 2-32 个字符', trigger: 'blur' }
  ],
  image: [
    { required: true, message: '请选择镜像', trigger: 'change' }
  ]
}

const handleCreate = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        // TODO: 调用创建虚拟机 API
        ElMessage.success('虚拟机创建成功')
        router.push('/vms')
      } catch (error) {
        ElMessage.error('创建失败：' + error.message)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped lang="scss">
.create-vm-container {
  max-width: 900px;
  
  .page-title {
    font-size: 20px;
    font-weight: bold;
  }
}
</style>
