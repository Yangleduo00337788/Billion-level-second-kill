<template>
  <div class="max-w-[1440px] mx-auto px-6 lg:px-10 py-6 h-[calc(100vh-4rem)]">
    <div class="bg-white rounded-2xl border border-gray-100 h-full flex flex-col overflow-hidden">
      <div class="flex items-center justify-between px-4 lg:px-6 py-3 border-b border-gray-100 bg-white">
        <div class="flex items-center gap-3 flex-1 min-w-0">
          <n-input
            v-model:value="title"
            placeholder="输入文章标题..."
            size="large"
            :bordered="false"
            class="text-lg font-semibold flex-1"
          />
        </div>
        <div class="flex items-center gap-2 flex-shrink-0">
          <n-select
            v-model:value="categoryId"
            :options="categoryOptions"
            placeholder="选择分类"
            size="small"
            style="min-width: 120px; max-width: none;"
            clearable
          />
          <n-button size="small" quaternary @click="generateTitle" class="hidden sm:inline-flex">
            AI 标题
          </n-button>
          <n-button size="small" quaternary @click="generateSummary" class="hidden sm:inline-flex">
            AI 摘要
          </n-button>
          <n-button size="small" quaternary @click="saveDraft">
            草稿
          </n-button>
          <n-upload
            :action="'/api/v1/upload'"
            :headers="uploadHeaders"
            :show-file-list="false"
            @finish="handleUploadFinish"
            class="inline-flex"
          >
            <n-button size="small" quaternary>图片</n-button>
          </n-upload>
          <n-button size="small" type="primary" @click="handlePublish">
            发布
          </n-button>
        </div>
      </div>

      <div class="px-4 lg:px-6 py-2 border-b border-gray-100 flex items-center gap-3">
        <n-input v-model:value="coverUrl" placeholder="文章封面图片 URL (可选)" size="small" clearable class="max-w-md" />
        <div v-if="coverUrl" class="w-12 h-8 rounded overflow-hidden flex-shrink-0">
          <img :src="coverUrl" class="w-full h-full object-cover" @error="coverUrl = ''" />
        </div>
      </div>

      <div class="flex flex-1 overflow-hidden">
        <div class="flex-1 flex flex-col border-r border-gray-100">
          <div class="flex items-center gap-1 px-4 py-2 bg-gray-50 border-b border-gray-100 overflow-x-auto">
            <n-button size="tiny" quaternary @click="insertMarkdown('**', '**')">
              <strong>B</strong>
            </n-button>
            <n-button size="tiny" quaternary @click="insertMarkdown('*', '*')">
              <em>I</em>
            </n-button>
            <n-button size="tiny" quaternary @click="insertMarkdown('### ', '')">
              H
            </n-button>
            <n-button size="tiny" quaternary @click="insertMarkdown('[', '](url)')">
              Link
            </n-button>
            <n-button size="tiny" quaternary @click="insertCode">
              Code
            </n-button>
            <n-button size="tiny" quaternary @click="insertMarkdown('> ', '')">
              Quote
            </n-button>
            <n-button size="tiny" quaternary @click="insertMarkdown('- ', '')">
              List
            </n-button>
          </div>
          <textarea
            ref="textareaRef"
            v-model="content"
            class="flex-1 w-full p-4 lg:p-6 text-sm leading-relaxed resize-none outline-none bg-white font-mono"
            placeholder="开始写文章... (支持 Markdown)"
          ></textarea>
        </div>

        <div class="hidden lg:block flex-1 overflow-y-auto bg-white p-4 lg:p-6">
          <div class="markdown-body" v-html="renderedContent"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import type { UploadFileInfo } from 'naive-ui'
import { marked } from 'marked'
import { useArticleStore } from '@/stores/article'
import { post } from '@/api/request'
import { articleApi } from '@/api/article'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const articleStore = useArticleStore()

const title = ref('')
const content = ref('')
const categoryId = ref<number | null>(null)
const textareaRef = ref<HTMLTextAreaElement | null>(null)
const articleId = ref<number | null>(null)
const coverUrl = ref('')

const uploadHeaders = computed(() => ({
  Authorization: 'Bearer ' + localStorage.getItem('token')
}))

function handleUploadFinish({ file }: { file: UploadFileInfo }) {
  const resp = (file as any).response
  if (resp) {
    const url = resp.url || resp.data?.url
    if (url) {
      const imgTag = '![' + (file.name || 'image') + '](' + url + ')\n'
      content.value += imgTag
      message.success('图片已插入')
    }
  }
}

const categoryOptions = ref<Array<{ label: string; value: number }>>([])

async function loadCategories() {
  try {
    const res = await articleApi.getCategories()
    categoryOptions.value = (res.data || []).map((c: any) => ({ label: c.name, value: c.id }))
  } catch {}
}

const renderedContent = computed(() => {
  if (!content.value) return '<p style="color: #9ca3af; text-align: center; padding-top: 4rem;">预览区域</p>'
  try {
    return marked(content.value) as string
  } catch {
    return content.value
  }
})

function insertMarkdown(before: string, after: string) {
  const textarea = textareaRef.value
  if (!textarea) return
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selected = content.value.substring(start, end)
  content.value = content.value.substring(0, start) + before + selected + after + content.value.substring(end)
  textarea.focus()
  const newPos = start + before.length
  textarea.setSelectionRange(newPos, newPos + selected.length)
}

function insertCode() {
  insertMarkdown('```\n', '\n```')
}

async function handlePublish() {
  if (!title.value.trim()) {
    message.warning('请输入文章标题')
    return
  }
  if (!content.value.trim()) {
    message.warning('请输入文章内容')
    return
  }
  try {
    const data = {
      title: title.value,
      content: content.value,
      summary: content.value.slice(0, 200),
      category_id: categoryId.value || undefined,
      cover: coverUrl.value || undefined
    }
    if (articleId.value) {
      await articleStore.update(articleId.value, data)
      message.success('更新成功')
    } else {
      await articleStore.create(data)
      message.success('发布成功')
    }
    localStorage.removeItem('editor-draft')
    title.value = ''
    content.value = ''
    categoryId.value = null
    coverUrl.value = ''
    articleId.value = null
    router.push('/')
  } catch {
    // handled by interceptor
  }
}

function saveDraft() {
  const draft = { title: title.value, content: content.value, categoryId: categoryId.value, coverUrl: coverUrl.value }
  localStorage.setItem('editor-draft', JSON.stringify(draft))
  message.success('草稿已保存')
}

async function generateTitle() {
  if (!content.value.trim()) {
    message.warning('请先写一些内容')
    return
  }
  try {
    const res = await post<{ title: string }>('/ai/generate-title', { content: content.value.slice(0, 2000) })
    if (res.data.title) {
      title.value = res.data.title
      message.success('标题已生成')
    }
  } catch {
    message.error('生成失败，请检查 AI 配置')
  }
}

async function generateSummary() {
  if (!content.value.trim()) {
    message.warning('请先写一些内容')
    return
  }
  try {
    const res = await post<{ summary: string }>('/ai/generate-summary', { content: content.value.slice(0, 3000) })
    if (res.data.summary) {
      content.value = '> ' + res.data.summary + '\n\n' + content.value
      message.success('摘要已生成')
    }
  } catch {
    message.error('生成失败，请检查 AI 配置')
  }
}

onMounted(async () => {
  loadCategories()
  const draft = localStorage.getItem('editor-draft')
  if (draft && !route.params.id) {
    try {
      const parsed = JSON.parse(draft)
      title.value = parsed.title || ''
      content.value = parsed.content || ''
      categoryId.value = parsed.categoryId || null
      coverUrl.value = parsed.coverUrl || ''
    } catch {}
  }

  if (route.params.id) {
    articleId.value = Number(route.params.id)
    try {
      const article = await articleStore.fetchById(articleId.value)
      title.value = article.title
      content.value = article.content
      categoryId.value = article.category_id || null
    } catch {
      message.error('文章不存在')
      router.push('/')
    }
  }
})

watch(
  () => ({ title: title.value, content: content.value, categoryId: categoryId.value }),
  () => {
    const draft = { title: title.value, content: content.value, categoryId: categoryId.value, coverUrl: coverUrl.value }
    localStorage.setItem('editor-draft', JSON.stringify(draft))
  },
  { deep: true }
)
</script>
