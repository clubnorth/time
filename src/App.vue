<template>
  <div class="app-shell">
    <YearMonthHeader v-if="!showStats && !showTodo" v-model="currentMonth" :months="availableMonths" />

    <div class="timeline-container" v-if="!showStats && !showTodo">
      <div class="timeline-line"></div>

      <TimelineEntryGroup
        v-for="group in allGroups"
        :key="group.date"
        :group="group"
        @delete-entry="handleDeleteEntry"
      />
    </div>

    <div class="bottom-zone">
      <button class="bottom-side-btn" @click="showStats = true; showTodo = false">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg>
        <span>统计</span>
      </button>
      <AddButton @add="handleAdd" />
      <button class="bottom-side-btn" @click="showTodo = true; showStats = false">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="5" width="6" height="6" rx="1"/><path d="M3 17h3l2-2"/><path d="M11 12h8"/><path d="M11 16h5"/></svg>
        <span>待办事项</span>
      </button>
    </div>
    <StatisticsPage v-if="showStats" @back="showStats = false" />
    <TodoPage v-if="showTodo" :key="todoKey" @back="showTodo = false" />
    <AddEntryPanel :visible="showAddPanel" @close="showAddPanel = false" @select="handleSelect" @backfill="handleBackfill" />
    <ThoughtFormPanel :visible="showThoughtForm" :kind="thoughtKind" @cancel="showThoughtForm = false" @create="handleThoughtCreate" />
    <AssetFormPanel :visible="showAssetForm" @cancel="showAssetForm = false" @create="handleAssetCreate" />
    <UricFormPanel :visible="showUricForm" @cancel="showUricForm = false" @create="handleUricCreate" />
    <ExerciseFormPanel :visible="showExerciseForm" @cancel="showExerciseForm = false" @create="handleExerciseCreate" />
    <TodoFormPanel :visible="showTodoForm" @cancel="showTodoForm = false" @create="handleTodoCreate" />
    <ReadingFormPanel :visible="showReadingForm" @cancel="showReadingForm = false" @create="handleReadingCreate" />
    <MovieFormPanel :visible="showMovieForm" @cancel="showMovieForm = false" @create="handleMovieCreate" />
    <ConfirmModal :visible="showConfirm" :message="confirmMessage" @close="handleConfirmClose" />
    <TimePickerModal
      :visible="bfShowPicker"
      :pickYear="bfPickY" :pickMonth="bfPickM" :pickDay="bfPickD"
      :pickHour="bfPickH" :pickMinute="bfPickI"
      @close="bfCloseAndHide"
      @confirm="handleBackfillConfirm"
      @adjustYear="(d) => { bfPickY += d }"
      @adjustMonth="bfAdjM"
      @adjustDay="bfAdjD"
      @adjustHour="bfAdjH"
      @adjustMinute="bfAdjI"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import YearMonthHeader from './components/YearMonthHeader.vue'
import TimelineEntryGroup from './components/TimelineEntryGroup.vue'
import AddButton from './components/AddButton.vue'
import ThoughtFormPanel from './components/ThoughtFormPanel.vue'
import AddEntryPanel from './components/AddEntryPanel.vue'
import AssetFormPanel from './components/AssetFormPanel.vue'
import UricFormPanel from './components/UricFormPanel.vue'
import ExerciseFormPanel from './components/ExerciseFormPanel.vue'
import TodoFormPanel from './components/TodoFormPanel.vue'
import ReadingFormPanel from './components/ReadingFormPanel.vue'
import MovieFormPanel from './components/MovieFormPanel.vue'
import StatisticsPage from './components/StatisticsPage.vue'
import TodoPage from './components/TodoPage.vue'
import ConfirmModal from './components/ConfirmModal.vue'
import TimePickerModal from './components/TimePickerModal.vue'

import { formatRecordTime } from './utils/date.js'
import { API_BASE } from './config.js'
import { useTimePicker } from './composables/useTimePicker.js'
const PAGE_SIZE = 30
const LOAD_MORE_SIZE = 20
const SCROLL_THRESHOLD = 300

const weekdays = ['日', '一', '二', '三', '四', '五', '六'].map(d => '周' + d)

const timelineData = ref([])
const currentMonth = ref('')
const availableMonths = ref([])
const skipScrollWatch = ref(false)
const isLoading = ref(false)
const hasMore = ref(true)

const showAddPanel = ref(false)
const showThoughtForm = ref(false)
const thoughtKind = ref('positive')
const showAssetForm = ref(false)
const showUricForm = ref(false)
const showExerciseForm = ref(false)
const showTodoForm = ref(false)
const showReadingForm = ref(false)
const showMovieForm = ref(false)
const showStats = ref(false)
const showTodo = ref(false)
const todoKey = ref(0)
const showConfirm = ref(false)
const confirmMessage = ref('')
const pendingDeleteId = ref(null)
const backfillType = ref('')
const {
  currentTime: bfCurrentTime,
  showTimeModal: bfShowPicker,
  pickYear: bfPickY,
  pickMonth: bfPickM,
  pickDay: bfPickD,
  pickHour: bfPickH,
  pickMinute: bfPickI,
  openTimePicker: bfOpen,
  closeTimePicker: bfClose,
  confirmTime: bfConfirmTime,
  adjustMonth: bfAdjM,
  adjustDay: bfAdjD,
  adjustHour: bfAdjH,
  adjustMinute: bfAdjI,
} = useTimePicker()

async function fetchEntries(limit, before) {
  let url = `${API_BASE}/api/entries?limit=${limit}`
  if (before) url += `&before=${encodeURIComponent(before)}`
  try {
    const res = await fetch(url)
    const json = await res.json()
    if (json.code === 0 && Array.isArray(json.data)) {
      return json.data
    }
  } catch (e) {
    console.error('Failed to fetch entries:', e)
  }
  return []
}

async function fetchAllEntries() {
  let all = []
  let before = ''
  while (true) {
    const data = await fetchEntries(100, before)
    if (data.length === 0) break
    all = all.concat(data)
    if (data.length < 100) break
    before = data[data.length - 1].recorded_at
  }
  return all
}

function starHTML(rating) {
  if (!rating || rating <= 0) return ''
  const filled = rating
  let html = ''
  for (let i = 1; i <= 5; i++) {
    const star = i * 2
    if (filled >= star) {
      html += '<svg class="r-star" width="14" height="14" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" fill="#D4A574"/></svg>'
    } else if (filled >= star - 1) {
      html += '<svg class="r-star" width="14" height="14" viewBox="0 0 24 24"><path d="M12 5.5l2.1 4.2 4.6.7-3.3 3.2.8 4.6-4.2-2.2-4.2 2.2.8-4.6L5.3 10.4l4.6-.7L12 5.5zM12 2l-3.09 6.26L2 9.27l5 4.87-1.18 6.88L12 17.77l6.18 3.25L17 14.14l5-4.87-6.91-1.01L12 2z" fill="#D4A574"/></svg>'
    } else {
      html += '<svg class="r-star" width="14" height="14" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" fill="var(--color-pencil)" opacity="0.3"/></svg>'
    }
  }
  html += ' ' + rating + '分'
  return html
}

function rebuildMonths(data) {
  const months = new Set()
  for (const e of data) {
    months.add(e.recorded_at.substring(0, 7))
  }
  const sorted = [...months].sort().reverse()
  availableMonths.value = sorted.map(m => ({
    label: `${m.substring(0,4)}年${parseInt(m.substring(5,7))}月`,
    value: m,
  }))
  if (sorted.length > 0 && !currentMonth.value) {
    currentMonth.value = sorted[0]
  }
}

async function loadInitial() {
  isLoading.value = true
  const data = await fetchEntries(PAGE_SIZE, '')
  timelineData.value = data
  hasMore.value = data.length >= PAGE_SIZE
  rebuildMonths(data)
  isLoading.value = false
}

async function loadMore() {
  if (isLoading.value || !hasMore.value) return
  const oldest = timelineData.value[timelineData.value.length - 1]
  if (!oldest) return

  isLoading.value = true
  const data = await fetchEntries(LOAD_MORE_SIZE, oldest.recorded_at)
  if (data.length < LOAD_MORE_SIZE) hasMore.value = false
  timelineData.value = [...timelineData.value, ...data]
  rebuildMonths([...timelineData.value])
  isLoading.value = false
}

const allGroups = computed(() => {
  if (timelineData.value.length === 0) return []
  
  const byDate = {}
  for (const entry of timelineData.value) {
    const d = entry.recorded_at.substring(0, 10)
    if (!byDate[d]) byDate[d] = []
    byDate[d].push(entry)
  }
  
  const groups = []
  const dates = Object.keys(byDate).sort().reverse()
  for (const date of dates) {
    const entries = byDate[date]
      .sort((a, b) => b.recorded_at.localeCompare(a.recorded_at))
      .map(e => {
        const isThought = e.type === 'thought'
        const isAsset = e.type === 'asset'
        const isUric = e.type === 'uric'
        const isExercise = e.type === 'exercise'
        const isReading = e.type === 'reading'
        const isMovie = e.type === 'movie'
        let title = e.title
        let desc = e.description
        const isDiscipline = e.type === 'discipline'
        const isNosugar = e.type === 'nosugar'
        if (isDiscipline) {
          desc = '今天是连续自律第 ' + (e.description || '1') + ' 天'
        } else if (isNosugar) {
          desc = '今天是连续无糖第 ' + (e.description || '1') + ' 天'
        } else if (isThought) {
          title = '随记'
        } else if (isExercise) {
          const sportName = e.title.replace('运动·', '')
          const baseDesc = {
            '跑步': '你今天跑步跑了 ' + e.description + ' 千米',
            '骑行': '你今天骑行了 ' + e.description + ' 千米',
            '俯卧撑': '你今天俯卧撑了 ' + e.description + ' 个',
            '跳绳': '你今天跳绳了 ' + e.description + ' 个',
            '游泳': '你今天游泳了 ' + e.description + ' 米',
            '徒步': '你今天徒步了 ' + e.description + ' 公里',
            '自由活动': '你今天自由活动了 ' + e.description + ' 分钟',
          }[sportName] || (e.description + ' ' + sportName)
          desc = baseDesc + (e.valence ? ' 本次消耗 ' + e.valence + ' 卡' : '')
        } else if (isAsset) {
          desc = '你现在的余额是 <span class="rainbow">' + (e.description || '0') + '</span> 元'
        } else if (isUric) {
          desc = '你今天的尿酸值是 ' + (e.description || '0') + ' mol'
        }
        return {
          id: e.id,
          time: e.recorded_at.substring(11, 16),
          title: isExercise ? '运动·' + e.title.replace('运动·', '') : title,
          description: desc,
          category: e.category,
          isAsset,
          isUric,
          isThought,
          isReading,
          isMovie,
          valence: e.valence,
        }
      })
    
    const d = new Date(date)
    groups.push({
      date,
      weekday: weekdays[d.getDay()],
      dateNum: `${d.getMonth() + 1}月${d.getDate()}日`,
      entries,
    })
  }
  return groups
})

function handleAdd() { showAddPanel.value = true }
function handleSelect(item) {
  if (item.kind) { thoughtKind.value = item.kind; showThoughtForm.value = true }
  else if (item.id === 'asset') { showAssetForm.value = true }
  else if (item.id === 'uric') { showUricForm.value = true }
  else if (item.id === 'discipline' || item.id === 'nosugar') { handleQuickCreate(item.id) }
  else if (item.id === 'exercise') { showExerciseForm.value = true }
  else if (item.id === 'todo') { showTodoForm.value = true }
  else if (item.id === 'reading') { showReadingForm.value = true }
  else if (item.id === 'movie') { showMovieForm.value = true }
  else { showAddPanel.value = false }
}

async function handleExerciseCreate(data) {
  const { recordedAt, monthKey } = formatRecordTime(data.time)

  // Calorie calculation（基于 69kg 体重、平均心率 160、标准 MET 公式）
  // kcal = MET × 体重(kg) × 时间(小时)
  const amount = parseFloat(data.amount) || 0
  const weight = 69
  const sport = data.sport
  const metCalcs = {
    '跑步':   Math.round(amount * 11.0 / 60 * weight * 6),       // 10km/h → 76 kcal/km
    '骑行':   Math.round(amount * 9.0  / 60 * weight * 60/22),   // 22km/h → 28 kcal/km
    '俯卧撑': Math.round(amount * 4.5  / 60 * weight / 30),      // 30个/min → 0.17 kcal/个
    '跳绳':   Math.round(amount * 11.5 / 60 * weight / 120),     // 120个/min → 0.11 kcal/个
    '游泳':   Math.round(amount * 8.0  / 60 * weight / 50),      // 50m/min → 0.18 kcal/m  
    '徒步':   Math.round(amount * 6.5  / 60 * weight * 60/5),    // 5km/h → 90 kcal/km
    '自由活动': Math.round(amount * 8.0 / 60 * weight),           // 9.2 kcal/min
  }
  const calories = metCalcs[sport] || 0

  try {
    const res = await fetch(`${API_BASE}/api/entries`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        type: 'exercise',
        title: '运动·' + sport,
        description: data.amount,
        category: 'green',
        valence: String(calories),
        recorded_at: recordedAt,
      }),
    })
    if ((await res.json()).code === 0) {
      await loadInitial()
      skipScrollWatch.value = true
      currentMonth.value = monthKey
      showStats.value = false
    }
  } catch (e) {
    console.error('Failed to create exercise entry:', e)
  }

  showExerciseForm.value = false
  showAddPanel.value = false
}

async function handleUricCreate(data) {
  const { recordedAt, monthKey } = formatRecordTime(data.time)

  try {
    const res = await fetch(`${API_BASE}/api/entries`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        type: 'uric',
        title: '尿酸记录',
        description: data.amount,
        category: 'yellow',
        recorded_at: recordedAt,
      }),
    })
    if ((await res.json()).code === 0) {
      await loadInitial()
      skipScrollWatch.value = true
      currentMonth.value = monthKey
      showStats.value = false
    }
  } catch (e) {
    console.error('Failed to create uric entry:', e)
  }

  showUricForm.value = false
  showAddPanel.value = false
}

async function handleAssetCreate(data) {
  const { recordedAt, monthKey } = formatRecordTime(data.time)

  try {
    const res = await fetch(`${API_BASE}/api/entries`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        type: 'asset',
        title: '资产记录',
        description: data.amount,
        category: 'yellow',
        recorded_at: recordedAt,
      }),
    })
    if ((await res.json()).code === 0) {
      await loadInitial()
      skipScrollWatch.value = true
      currentMonth.value = monthKey
      showStats.value = false
    }
  } catch (e) {
    console.error('Failed to create asset entry:', e)
  }

  showAssetForm.value = false
  showAddPanel.value = false
}



async function getConsecutiveDay(type) {
  try {
    const res = await fetch(`${API_BASE}/api/entries?limit=30`)
    const json = await res.json()
    if (json.code !== 0 || !Array.isArray(json.data)) return 1
    const filtered = json.data.filter(e => e.type === type)
    if (filtered.length === 0) return 1
    const latest = filtered[0]
    const latestDate = new Date(latest.recorded_at.substring(0, 10))
    const today = new Date()
    const yesterday = new Date(today)
    yesterday.setDate(yesterday.getDate() - 1)
    const lds = latestDate.toDateString()
    if (lds === today.toDateString()) return -1
    if (lds === yesterday.toDateString()) return (parseInt(latest.description) || 0) + 1
    return 1
  } catch (e) { return 1 }
}

function handleDeleteEntry(id) {
  pendingDeleteId.value = id
  confirmMessage.value = '确定删除这条记录？'
  showConfirm.value = true
}

async function handleConfirmClose() {
  showConfirm.value = false
  const id = pendingDeleteId.value
  pendingDeleteId.value = null
  if (id) {
    try {
      const res = await fetch(`${API_BASE}/api/entries/${id}`, { method: 'DELETE' })
      if ((await res.json()).code === 0) {
        await loadInitial()
      }
    } catch (e) { console.error('Failed to delete entry:', e) }
  }
}

function handleBackfill(item) {
  backfillType.value = item.id
  bfOpen()
}

function bfCloseAndHide() {
  bfClose()
  showAddPanel.value = false
}

async function handleBackfillConfirm() {
  bfConfirmTime()
  const d = bfCurrentTime.value
  const { recordedAt, monthKey } = formatRecordTime(d)
  const type = backfillType.value
  const dateStr = recordedAt.substring(0, 10)

  // Check for duplicate entry (same type, same date)
  let duplicate = false
  try {
    const res = await fetch(`${API_BASE}/api/entries?limit=200`)
    const json = await res.json()
    if (json.code === 0) {
      duplicate = json.data.some(e => e.type === type && e.recorded_at.substring(0, 10) === dateStr)
    }
  } catch (e) {
    console.error('Failed to check duplicate entry:', e)
    confirmMessage.value = '网络异常，请重试'
    showConfirm.value = true
    return
  }
  if (duplicate) {
    confirmMessage.value = '该天已经打过卡了'
    showConfirm.value = true
    return
  }

  const title = type === 'discipline' ? '自律' : '禁止糖分'
  try {
    const res = await fetch(`${API_BASE}/api/entries`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ type, title, description: '1', category: 'green', recorded_at: recordedAt }),
    })
    if ((await res.json()).code === 0) {
      try { await fetch('/api/entries/recalculate?type='+type, { method: 'POST' }) } catch (e) { console.error('Recalculate failed:', e) }
      await loadInitial()
      skipScrollWatch.value = true
      currentMonth.value = monthKey
      showStats.value = false
      showAddPanel.value = false
    }
  } catch (e) { console.error('Backfill create failed:', e) }
  bfClose()
}

async function handleQuickCreate(type) {
  const dayCount = await getConsecutiveDay(type)
  if (dayCount < 0) { confirmMessage.value = '今天已经打过卡了，明天再来吧'; showConfirm.value = true; showAddPanel.value = false; return }
  const { recordedAt, monthKey } = formatRecordTime(new Date())

  const title = type === 'discipline' ? '自律' : '禁止糖分'

  try {
    const res = await fetch(`${API_BASE}/api/entries`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ type, title, description: String(dayCount), category: 'green', recorded_at: recordedAt }),
    })
    if ((await res.json()).code === 0) {
      await loadInitial()
      skipScrollWatch.value = true
      currentMonth.value = monthKey
      showStats.value = false
    }
  } catch (e) { console.error('Quick create failed:', e) }

  showAddPanel.value = false
}

async function handleThoughtCreate(data) {
  const { recordedAt, monthKey } = formatRecordTime(data.time)

  try {
    const res = await fetch(`${API_BASE}/api/entries`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        type: 'thought',
        title: data.kind === 'positive' ? '随记·太阳' : '随记·阴雨',
        description: data.note,
        category: data.kind === 'positive' ? 'green' : 'red',
        valence: data.kind,
        recorded_at: recordedAt,
      }),
    })
    if ((await res.json()).code === 0) {
      await loadInitial()
      skipScrollWatch.value = true
      currentMonth.value = monthKey
      showStats.value = false
    }
  } catch (e) {
    console.error('Failed to create entry:', e)
  }

  showThoughtForm.value = false
  showAddPanel.value = false
}

async function handleTodoCreate(data) {
  const d = data.time
  const dueDate = d.getFullYear() + '-' +
    String(d.getMonth() + 1).padStart(2, '0') + '-' +
    String(d.getDate()).padStart(2, '0') + ' ' +
    String(d.getHours()).padStart(2, '0') + ':' +
    String(d.getMinutes()).padStart(2, '0') + ':00'

  const today = new Date().toISOString().substring(0, 10)
  const dueDay = dueDate.substring(0, 10)
  let cat = 'later'
  if (dueDay < today) cat = 'overdue'
  else if (dueDay === today) cat = 'today'
  else {
    const t = new Date(today)
    t.setDate(t.getDate() + 3)
    if (dueDay <= t.toISOString().substring(0, 10)) cat = 'upcoming'
  }

  try {
    const res = await fetch(`${API_BASE}/api/todos`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: data.note, due_date: dueDate, category: cat }),
    })
    if ((await res.json()).code === 0) {
      showTodoForm.value = false
      showAddPanel.value = false
      showStats.value = false
      todoKey.value++
      showTodo.value = true
    }
  } catch (e) { console.error('Failed to create todo:', e) }
}

async function handleReadingCreate(data) {
  const { recordedAt, monthKey } = formatRecordTime(data.time)
  const d = data.time
  const dateStr = d.getFullYear() + '年' + (d.getMonth() + 1) + '月' + d.getDate() + '日'

  if (data.status === 'reading') {
    const desc = dateStr + ' 开始阅读《' + data.name + '》'
    try {
      const res = await fetch(`${API_BASE}/api/entries`, {
        method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ type: 'reading', title: '开始阅读', description: desc, category: 'green', valence: 'reading', recorded_at: recordedAt }),
      })
      if ((await res.json()).code === 0) {
        await loadInitial(); skipScrollWatch.value = true; currentMonth.value = monthKey; showStats.value = false
      }
    } catch (e) { console.error(e) }
  } else {
    // Check if there's a matching "开始阅读" record
    const bookName = data.name.trim()
    let hasReadingRecord = false
    try {
      const all = await fetchAllEntries()
      hasReadingRecord = all.some(e =>
        e.type === 'reading' && e.title === '开始阅读' &&
        e.description && e.description.includes(bookName)
      )
    } catch (e) { console.error('Failed to check reading records:', e) }

    if (!hasReadingRecord) {
      confirmMessage.value = '未找到《' + bookName + '》的在读记录，不能添加已读'
      showConfirm.value = true
      showReadingForm.value = false
      return
    }

    await createFinishedReading(data)
  }

  showReadingForm.value = false
  showAddPanel.value = false
}

async function createFinishedReading(data) {
  const { recordedAt, monthKey } = formatRecordTime(data.time)
  const lines = ['书名：' + data.name]
  if (data.bookInfo && data.bookInfo.author) {
    const bi = data.bookInfo
    const authors = bi.author.split(/[,，、/]/).map(a => a.trim()).filter(Boolean)
    if (bi.nationality) {
      lines.push('作者：' + authors.map(a => a + '【' + bi.nationality + '】').join(' / '))
    } else {
      lines.push('作者：' + authors.join(' / '))
    }
  } else {
    lines.push('作者：未知')
  }
  if (data.duration) lines.push('用时：' + data.duration + 'h')
  if (data.bookInfo && data.bookInfo.first_publish_date) lines.push('首次出版：' + data.bookInfo.first_publish_date + '年')
  if (data.category) {
    const catColor = { '推理悬疑':'#8B6B4A','明日方舟':'#5B7A9E','人物历史':'#9E7A5B','社会科学':'#6B8E7A','小说文学':'#8A6B9E','自我成长':'#7A9E6B','杂项拾遗':'#9E8A6B' }[data.category] || '#888'
    lines.push('分类：<span class="r-cat" style="background:'+catColor+'">' + data.category + '</span>')
  }
  if (data.bookInfo && data.bookInfo.tags && data.bookInfo.tags.length) {
    const tags = data.bookInfo.tags.map(t => '<span class="r-tag">' + t + '</span>').join(' ')
    lines.push('标签：' + tags)
  }
  if (data.rating) lines.push('评分：' + starHTML(data.rating))
  const desc = lines.join('<br>')

  try {
    const res = await fetch(`${API_BASE}/api/entries`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ type: 'reading', title: '读完', description: desc, category: 'green', valence: data.category || 'finished', recorded_at: recordedAt }),
    })
    if ((await res.json()).code === 0) {
      await loadInitial(); skipScrollWatch.value = true; currentMonth.value = monthKey; showStats.value = false
    }
  } catch (e) { console.error(e) }
}

async function handleMovieCreate(data) {
  const { recordedAt, monthKey } = formatRecordTime(data.time)
  const typeLabel = { movie: '电影', series: '剧集', anime: '动漫' }[data.mediaType] || '影视'
  const title = '影视·' + typeLabel
  const mediaName = data.name.trim()

  // For series/anime "看完": check for matching "在看" record
  if (data.watchStatus === 'finished' && (data.mediaType === 'series' || data.mediaType === 'anime')) {
    let hasWatchingRecord = false
    try {
      const all = await fetchAllEntries()
      hasWatchingRecord = all.some(e =>
        e.type === 'movie' && e.title === '开始观看' &&
        e.description && e.description.includes(mediaName)
      )
    } catch (e) { console.error('Failed to check watching records:', e) }
    if (!hasWatchingRecord) {
      confirmMessage.value = '未找到《' + mediaName + '》的在看记录，不能添加看完'
      showConfirm.value = true
      showMovieForm.value = false
      return
    }
  }

  // For "在看": simple format like "开始阅读"
  if (data.watchStatus === 'watching') {
    const d = data.time
    const dateStr = d.getFullYear() + '年' + (d.getMonth() + 1) + '月' + d.getDate() + '日'
    const desc = dateStr + ' 开始观看《' + mediaName + '》'
    try {
      const res = await fetch(`${API_BASE}/api/entries`, {
        method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ type: 'movie', title: '开始观看', description: desc, category: 'green', valence: data.mediaType, recorded_at: recordedAt }),
      })
      if ((await res.json()).code === 0) {
        await loadInitial(); skipScrollWatch.value = true; currentMonth.value = monthKey; showStats.value = false
      }
    } catch (e) { console.error('Failed to create movie entry:', e) }
    showMovieForm.value = false; showAddPanel.value = false
    return
  }

  // Full info card (movie or series/anime "看完")
  const lines = ['名称：' + mediaName]
  const mi = data.mediaInfo
  if (mi && mi.director) {
    const directors = mi.director.split(/[,，、/]/).map(a => a.trim()).filter(Boolean)
    lines.push('导演：' + directors.join(' / '))
  }
  if (mi && mi.cast) {
    const cast = mi.cast.split(/[,，、/]/).map(a => a.trim()).filter(Boolean)
    lines.push('主演：' + cast.join(' / '))
  }
  if (mi && mi.country) {
    const countries = mi.country.split(/[,，、/]/).map(a => a.trim()).filter(Boolean)
    lines.push('国家：' + countries.map(c => '【' + c + '】').join(' '))
  }
  if (mi && mi.first_release_date) lines.push('首次上映：' + mi.first_release_date + '年')
  if (data.mediaType === 'movie' && mi && mi.duration) {
    lines.push('时长：' + mi.duration + 'min')
  }
  if ((data.mediaType === 'series' || data.mediaType === 'anime') && mi && mi.episodes) {
    lines.push('集数：' + mi.episodes)
  }
  if (mi && mi.tags && mi.tags.length) {
    const tags = mi.tags.map(t => '<span class="r-tag">' + t + '</span>').join(' ')
    lines.push('标签：' + tags)
  }
  if (data.rating) lines.push('评分：' + starHTML(data.rating))
  const desc = lines.join('<br>')

  try {
    const res = await fetch(`${API_BASE}/api/entries`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ type: 'movie', title: data.mediaType === 'movie' ? title : '看完', description: desc, category: 'green', valence: data.mediaType, recorded_at: recordedAt }),
    })
    if ((await res.json()).code === 0) {
      await loadInitial()
      skipScrollWatch.value = true
      currentMonth.value = monthKey
      showStats.value = false
    }
  } catch (e) { console.error('Failed to create movie entry:', e) }
  showMovieForm.value = false
  showAddPanel.value = false
}

function handleScroll() {
  // Scroll-triggered month detection
  const dots = document.querySelectorAll('[data-date]')
  for (const dot of dots) {
    const rect = dot.getBoundingClientRect()
    if (rect.top > 55) {
      const date = dot.dataset.date
      const month = date.substring(0, 7)
      if (month !== currentMonth.value) {
        skipScrollWatch.value = true
        currentMonth.value = month
      }
      break
    }
  }
  // Load more when near bottom
  const scrollBottom = window.innerHeight + window.scrollY
  const docHeight = document.documentElement.scrollHeight
  if (docHeight - scrollBottom < SCROLL_THRESHOLD) {
    loadMore()
  }
}

watch(currentMonth, (newMonth) => {
  if (skipScrollWatch.value) { skipScrollWatch.value = false; return }
  if (!newMonth) return
  const dots = document.querySelectorAll('[data-date]')
  for (const dot of dots) {
    if (dot.dataset.date.startsWith(newMonth)) {
      dot.scrollIntoView({ behavior: 'smooth', block: 'start' })
      break
    }
  }
})

onMounted(() => {
  loadInitial()
  window.addEventListener('scroll', handleScroll, { passive: true })
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>
<style scoped>
.app-shell {
  position: relative;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
}

.timeline-container {
  position: relative;
  display: grid;
  grid-template-columns: 46px 4px 20px 9px 1fr;
  padding: 64px 24px 60px 24px;
  flex: 1;
}

@media (min-width: 600px) {
  .timeline-container {
    grid-template-columns: 52px 6px 24px 12px 1fr;
    padding: 72px 28px 64px 28px;
  }
}

@media (min-width: 768px) {
  .timeline-container {
    grid-template-columns: 56px 6px 28px 12px 1fr;
    padding: 80px 36px 68px 36px;
  }
}

.timeline-line {
  position: absolute;
  left: calc(24px + 46px + 4px + 10px);
  top: 0;
  bottom: 0;
  width: 1px;
  background: var(--color-pencil);
  z-index: 0;
  pointer-events: none;
}

@media (min-width: 600px) {
  .timeline-line {
    left: calc(28px + 52px + 6px + 12px);
  }
}

@media (min-width: 768px) {
  .timeline-line {
    left: calc(36px + 56px + 6px + 14px);
  }
}

.bottom-zone {
  position: fixed;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 100%;
  max-width: var(--content-width);
  height: 60px;
  background: #fcfcfc;
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-top: 1px solid var(--color-border-light);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  z-index: 50;
}

@media (min-width: 600px) {
  .bottom-zone {
    padding: 0 28px;
    height: 64px;
  }
}

@media (min-width: 768px) {
  .bottom-zone {
    padding: 0 36px;
  }
}

.bottom-side-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  font-size: 11px;
  font-weight: 500;
  color: var(--color-ink);
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px 12px;
  flex: 1;
  min-width: 0;
  transition: color 0.15s;
}
.bottom-side-btn:hover {
  color: var(--color-ink);
}
.bottom-side-btn svg {
  color: var(--color-graphite);
  transition: color 0.15s;
}
.bottom-side-btn:hover svg {
  color: var(--color-ink);
}

</style>
