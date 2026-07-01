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
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted, onUnmounted } from 'vue'
import YearMonthHeader from './components/YearMonthHeader.vue'
import TimelineEntryGroup from './components/TimelineEntryGroup.vue'
import AddButton from './components/AddButton.vue'
import ThoughtFormPanel from './components/ThoughtFormPanel.vue'
import AddEntryPanel from './components/AddEntryPanel.vue'

const API_BASE = 'http://localhost:8080'
const weekdays = ['日', '一', '二', '三', '四', '五', '六'].map(d => '周' + d)

// Data
const timelineData = ref([])
const currentMonth = ref('')
const availableMonths = ref([])
const skipScrollWatch = ref(false)

const showAddPanel = ref(false)
const showThoughtForm = ref(false)
const thoughtKind = ref('positive')

async function fetchAllEntries() {
  try {
    const res = await fetch(`${API_BASE}/api/entries`)
    const json = await res.json()
    if (json.code === 0 && Array.isArray(json.data)) {
      timelineData.value = json.data
      // Derive months from data
      const months = new Set()
      for (const e of json.data) {
        months.add(e.recorded_at.substring(0, 7))
      }
      const sorted = [...months].sort().reverse()
      availableMonths.value = sorted.map(m => ({
        label: `${m.substring(0,4)}年${parseInt(m.substring(5,7))}月`,
        value: m,
      }))
      // Set current month to the first available
      if (sorted.length > 0 && !currentMonth.value) {
        currentMonth.value = sorted[0]
      }
    }
  } catch (e) {
    console.error('Failed to fetch entries:', e)
  }
}

// Convert entries to timeline groups (all data, sorted by time)
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
      .map(e => ({
        id: e.id,
        time: e.recorded_at.substring(11, 16),
        title: e.title,
        description: e.description,
        category: e.category,
      }))
    
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
  else { showAddPanel.value = false }
}

async function handleThoughtCreate(data) {
  const d = data.time
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const monthKey = `${y}-${m}`
  const recordedAt = `${monthKey}-${day} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}:00`

  const entry = {
    type: 'thought',
    title: data.kind === 'positive' ? '念头·正' : '念头·负',
    description: data.note,
    category: data.kind === 'positive' ? 'green' : 'red',
    valence: data.kind,
    recorded_at: recordedAt,
  }

  try {
    const res = await fetch(`${API_BASE}/api/entries`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(entry),
    })
    const json = await res.json()
    if (json.code === 0) {
      await fetchAllEntries()
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
  fetchAllEntries()
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
