<template>
  <div class="app-shell">
    <YearMonthHeader v-model="currentMonth" :months="availableMonths" />

    <div class="timeline-container">
      <div class="timeline-line"></div>

      <TimelineEntryGroup
        v-for="group in allGroups"
        :key="group.date"
        :group="group"
      />
    </div>

    <div class="bottom-zone">
      <AddButton @add="handleAdd" />
    </div>
    <AddEntryPanel :visible="showAddPanel" @close="showAddPanel = false" @select="handleSelect" />
    <ThoughtFormPanel :visible="showThoughtForm" :kind="thoughtKind" @cancel="showThoughtForm = false" @create="handleThoughtCreate" />
    <AssetFormPanel :visible="showAssetForm" @cancel="showAssetForm = false" @create="handleAssetCreate" />
    <UricFormPanel :visible="showUricForm" @cancel="showUricForm = false" @create="handleUricCreate" />
    <ExerciseFormPanel :visible="showExerciseForm" @cancel="showExerciseForm = false" @create="handleExerciseCreate" />
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

const API_BASE = 'http://localhost:8080'
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
        const isAsset = e.type === 'asset'
        const isUric = e.type === 'uric'
        const isExercise = e.type === 'exercise'
        let title = e.title
        let desc = e.description
        if (isExercise) {
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
  else if (item.id === 'exercise') { showExerciseForm.value = true }
  else { showAddPanel.value = false }
}

async function handleExerciseCreate(data) {
  const d = data.time
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const monthKey = `${y}-${m}`
  const recordedAt = `${monthKey}-${day} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}:00`

  // Calorie calculation (weight: 70kg, based on MET formulas)
  const amount = parseFloat(data.amount) || 0
  const weight = 70
  const sport = data.sport
  const metCalcs = {
    '跑步': amount * 56,           // 跑步: per km
    '骑行': amount * 21,           // 骑行: per km
    '俯卧撑': amount * 0.15,   // 俯卧撑: per rep
    '跳绳': amount * 0.12,         // 跳绳: per rep
    '游泳': amount * 0.28,         // 游泳: per m
    '徒步': amount * 56,           // 徒步: per km
    '自由活动': amount * 3.5, // 自由活动: per min
  }
  const calories = Math.round(metCalcs[sport] || 0)

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
    }
  } catch (e) {
    console.error('Failed to create exercise entry:', e)
  }

  showExerciseForm.value = false
  showAddPanel.value = false
}

async function handleUricCreate(data) {
  const d = data.time
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const monthKey = `${y}-${m}`
  const recordedAt = `${monthKey}-${day} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}:00`

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
    }
  } catch (e) {
    console.error('Failed to create uric entry:', e)
  }

  showUricForm.value = false
}

async function handleAssetCreate(data) {
  const d = data.time
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const monthKey = `${y}-${m}`
  const recordedAt = `${monthKey}-${day} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}:00`

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
    }
  } catch (e) {
    console.error('Failed to create asset entry:', e)
  }

  showAssetForm.value = false
}

async function handleThoughtCreate(data) {
  const d = data.time
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const monthKey = `${y}-${m}`
  const recordedAt = `${monthKey}-${day} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}:00`

  try {
    const res = await fetch(`${API_BASE}/api/entries`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        type: 'thought',
        title: data.kind === 'positive' ? '念头·正' : '念头·负',
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
    }
  } catch (e) {
    console.error('Failed to create entry:', e)
  }

  showThoughtForm.value = false
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
  padding: 64px 24px 48px 24px;
  flex: 1;
}

.timeline-line {
  position: absolute;
  left: calc(24px + 46px + 4px + 10px);
  top: 0;
  bottom: 0;
  width: 1px;
  background: #d0d0d0;
  z-index: 0;
  pointer-events: none;
}

.bottom-zone {
  position: fixed;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 100%;
  max-width: 480px;
  height: 60px;
  background: #f8f8f8;
  border-top: 1px solid #ebebeb;
  display: flex;
  flex-direction: column;
  align-items: center;
  z-index: 50;
}

</style>
