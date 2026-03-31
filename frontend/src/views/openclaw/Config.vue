<template>
  <div class="config-container">
    <Header>
      <template #left>
        <span class="page-title">OpenClaw 管控配置</span>
      </template>
    </Header>
    
    <el-tabs v-model="activeTab" type="border-card">
      <el-tab-pane label="模型配置" name="model">
        <el-form label-width="150px">
          <el-form-item label="默认模型">
            <el-select v-model="config.model" style="width: 100%">
              <el-option label="qwencode/qwen3.5-plus" value="qwencode/qwen3.5-plus" />
              <el-option label="claude-3-sonnet" value="claude-3-sonnet" />
              <el-option label="gpt-4-turbo" value="gpt-4-turbo" />
            </el-select>
          </el-form-item>
          
          <el-form-item label="API 密钥">
            <el-input v-model="config.apiKey" type="password" show-password placeholder="请输入 API 密钥" />
          </el-form-item>
          
          <el-form-item label="模型切换策略">
            <el-radio-group v-model="config.switchStrategy">
              <el-radio label="auto">自动</el-radio>
              <el-radio label="manual">手动</el-radio>
            </el-radio-group>
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" @click="handleSave('model')">保存配置</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
      
      <el-tab-pane label="工具配置" name="tools">
        <el-table :data="toolList" style="width: 100%">
          <el-table-column prop="name" label="工具名称" min-width="150" />
          <el-table-column prop="description" label="描述" min-width="300" />
          <el-table-column prop="enabled" label="状态" width="100">
            <template #default="{ row }">
              <el-switch v-model="row.enabled" />
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120">
            <template #default="{ row }">
              <el-button link @click="handleEditTool(row)">配置</el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <div style="margin-top: 20px">
          <el-button type="primary" @click="handleSave('tools')">保存配置</el-button>
        </div>
      </el-tab-pane>
      
      <el-tab-pane label="技能配置" name="skills">
        <el-table :data="skillList" style="width: 100%">
          <el-table-column prop="name" label="技能名称" min-width="150" />
          <el-table-column prop="version" label="版本" width="100" />
          <el-table-column prop="enabled" label="状态" width="100">
            <template #default="{ row }">
              <el-switch v-model="row.enabled" />
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180">
            <template #default="{ row }">
              <el-button link @click="handleUpdateSkill(row)">更新</el-button>
              <el-button type="danger" link @click="handleUninstallSkill(row)">卸载</el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <div style="margin-top: 20px">
          <el-button type="primary" @click="handleInstallSkill">安装新技能</el-button>
          <el-button type="primary" @click="handleSave('skills')">保存配置</el-button>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import Header from '@/components/Header.vue'

const activeTab = ref('model')

const config = reactive({
  model: 'qwencode/qwen3.5-plus',
  apiKey: '',
  switchStrategy: 'auto'
})

const toolList = ref([
  { name: 'browser', description: '浏览器自动化控制', enabled: true },
  { name: 'exec', description: 'Shell 命令执行', enabled: true },
  { name: 'message', description: '消息发送与管理', enabled: true },
  { name: 'wecom_mcp', description: '企业微信 MCP 工具', enabled: true }
])

const skillList = ref([
  { name: 'qqbot-remind', version: '1.0.0', enabled: true },
  { name: 'wecom-doc', version: '1.2.0', enabled: true },
  { name: 'wecom-schedule', version: '1.1.0', enabled: true },
  { name: 'frontend-design', version: '2.0.0', enabled: true }
])

const handleSave = (tab) => {
  ElMessage.success(`${tab}配置已保存`)
}

const handleEditTool = (tool) => {
  ElMessage.info(`配置工具：${tool.name}`)
}

const handleUpdateSkill = (skill) => {
  ElMessage.success(`技能 ${skill.name} 已更新到最新版本`)
}

const handleUninstallSkill = (skill) => {
  skill.enabled = false
  ElMessage.success(`已卸载技能：${skill.name}`)
}

const handleInstallSkill = () => {
  ElMessage.info('打开技能安装对话框')
}
</script>

<style scoped lang="scss">
.config-container {
  .page-title {
    font-size: 20px;
    font-weight: bold;
  }
}
</style>
