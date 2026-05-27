<template>
  <div class="max-w-[900px] mx-auto px-6 lg:px-10 py-8">
    <button
      class="mb-6 inline-flex items-center gap-2 text-sm text-gray-500 hover:text-dark transition-colors glass-button px-3 py-1.5 rounded-xl"
      @click="router.back()"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
      返回
    </button>

    <div class="bg-white rounded-2xl border border-gray-100 p-6 lg:p-8">
      <h1 class="text-2xl font-bold text-dark mb-6">发布 Prompt</h1>

      <n-form ref="formRef" :model="form" label-placement="top">
        <n-form-item label="标题" path="title" :rule="{ required: true, message: '请输入标题', trigger: 'blur' }">
          <n-input v-model:value="form.title" placeholder="给你的 Prompt 起个名字" maxlength="200" show-count />
        </n-form-item>

        <n-form-item label="Prompt 内容" path="content" :rule="{ required: true, message: '请输入 Prompt 内容', trigger: 'blur' }">
          <n-input v-model:value="form.content" type="textarea" placeholder="输入 Prompt 内容，支持 {{变量}} 模板语法" :rows="8" />
        </n-form-item>

        <n-form-item label="描述">
          <n-input v-model:value="form.description" type="textarea" placeholder="简要描述这个 Prompt 的用途和使用场景" :rows="3" maxlength="500" show-count />
        </n-form-item>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <n-form-item label="分类">
            <n-select v-model:value="form.category" :options="categoryOptions" placeholder="选择分类" clearable class="w-full" />
          </n-form-item>

          <n-form-item label="适用模型">
            <n-select v-model:value="form.model" :options="modelOptions" placeholder="选择模型" clearable class="w-full" />
          </n-form-item>
        </div>

        <n-form-item label="标签">
          <n-input v-model:value="form.tags" placeholder="多个标签用逗号分隔，如：代码审查,编程,GPT" />
        </n-form-item>

        <div class="flex justify-end gap-3 mt-4">
          <n-button @click="router.back()">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="handleSubmit">发布</n-button>
        </div>
      </n-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { get, post } from '@/api/request'

const router = useRouter()
const message = useMessage()
const formRef = ref()
const submitting = ref(false)

const form = ref({
  title: '',
  content: '',
  description: '',
  category: null as string | null,
  model: null as string | null,
  tags: ''
})

const categoryOptions = ref<Array<{ label: string; value: string }>>([])

async function loadCategories() {
  try {
    const res = await get('/articles/categories')
    const cats = res.data || []
    categoryOptions.value = cats.map((c: any) => ({ label: c.name, value: c.name }))
  } catch {
    categoryOptions.value = [
      { label: '写作', value: '写作' },
      { label: '编程', value: '编程' },
      { label: '对话', value: '对话' },
      { label: '其他', value: '其他' }
    ]
  }
}

loadCategories()

const modelOptions = [
  { label: 'GPT-4', value: 'gpt-4' },
  { label: 'GPT-3.5', value: 'gpt-3.5' },
  { label: 'DeepSeek', value: 'deepseek' },
  { label: '通用', value: '通用' }
]

async function handleSubmit() {
  if (!form.value.title.trim()) {
    message.warning('请输入标题')
    return
  }
  if (!form.value.content.trim()) {
    message.warning('请输入 Prompt 内容')
    return
  }

  submitting.value = true
  try {
    const data: any = {
      title: form.value.title,
      content: form.value.content,
      description: form.value.description
    }
    if (form.value.category) data.category = form.value.category
    if (form.value.model) data.model = form.value.model
    if (form.value.tags.trim()) data.tags = JSON.stringify(form.value.tags.split(',').map((t: string) => t.trim()).filter(Boolean))

    await post('/prompts', data)
    message.success('发布成功')
    router.push('/prompt')
  } catch {
    message.error('发布失败')
  } finally {
    submitting.value = false
  }
}
</script>