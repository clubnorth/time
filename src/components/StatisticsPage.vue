<template>
  <div class="stats-page">
    <div class="stats-header">
      <button class="back-btn" @click="$emit('back')">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="15 18 9 12 15 6"/></svg>
      </button>
      <span class="header-title">统计</span>
      <span class="header-spacer"></span>
    </div>

    <div class="stat-grid">
      <div class="stat-card" v-for="s in topStats" :key="s.label">
        <span class="stat-value">{{ s.value }}</span>
        <span class="stat-label">{{ s.label }}</span>
      </div>
    </div>

    <div class="year-nav">
      <button class="year-arrow" @click="year--" aria-label="???">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"/></svg>
      </button>
      <span class="year-text">{{ year }}</span>
      <button class="year-arrow" @click="year++" :disabled="year >= maxYear" aria-label="???">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"/></svg>
      </button>
    </div>

    <div class="heatmaps">
      <div class="heatmap-card" v-for="cat in categoryData" :key="cat.type">
        <div class="heatmap-title-row">
          <div class="heatmap-dot" :style="{ background: cat.color }"></div>
          <span class="heatmap-title">{{ cat.name }}</span>
          <span class="heatmap-count">{{ cat.count }}次</span>
        </div>
        <div class="hm-scroll">
            <div class="hm-weeks" v-for="(week, wi) in cat.weeks" :key="wi">
              <template v-for="(day, di) in week" :key="di">
              <template v-if="day"><div
                class="hm-cell"
                :class="{ filled: day.filled, today: day.isToday }"
                :style="day.filled ? { background: cat.color } : {}"
                :title="day.date"
              >
                <span class="hm-date">{{ day.dateNum }}</span>
                <span class="hm-month">{{ day.label }}</span>
              </div>
              </template>
              <template v-else>
              <div class="hm-cell hm-empty"></div>
              </template>
              </template>
            </div>
          </div></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { API_BASE } from '../config.js'

const emit = defineEmits(['back'])
const year = ref(new Date().getFullYear())
const maxYear = ref(new Date().getFullYear())
const allEntries = ref([])

const categories = [
  { type: 'thought',  name: '随记',     color: '#9DB5C9' },
  { type: 'asset',    name: '资产记录', color: '#D4B87A' },
  { type: 'uric',     name: '尿酸记录', color: '#BE999B' },
  { type: 'exercise', name: '运动',     color: '#84B8A4' },
  { type: 'discipline', name: '自律',   color: '#8CA4BD' },
  { type: 'nosugar',    name: '禁止糖分', color: '#C88C8C' },
  { type: 'reading',   name: '读书',     color: '#A099C4' },
  { type: 'movie',     name: '影视',     color: '#CB99B0' },
]

const monthNames = ['一','二','三','四','五','六','七','八','九','十','十一','十二']

function buildHeatmapWeeks(entries, y) {
  const filledSet = new Set(entries.map(e => e.recorded_at.substring(0, 10)))
  const today = new Date()
  const currentYear = today.getFullYear()
  const currentDate = today.getDate()
  const currentMonth = today.getMonth()

  // Start from Jan 1 of the year
  const start = new Date(y, 0, 1)
  const startDay = start.getDay() // 0=Sun

  // Calculate end date: if viewing current year, end at today; else end at Dec 31
  let endDate
  if (y === currentYear) {
    endDate = new Date(y, currentMonth, currentDate)
  } else if (y > currentYear) {
    endDate = new Date(y, currentMonth, currentDate) // same as today for future years
  } else {
    endDate = new Date(y, 11, 31)
  }

  // Count total days
  const totalDays = Math.floor((endDate - start) / 86400000) + 1
	
  // Build weeks array: each week has 7 slots [Sun..Sat]
  const weeks = []
  let dayOffset = 0

  // First week: pad leading empty days
  let firstWeek = []
  for (let d = 0; d < 7; d++) {
    if (d < startDay) {
      firstWeek.push(null) // empty cell
    } else if (dayOffset < totalDays) {
      const date = new Date(y, 0, 1 + dayOffset)
      const ds = `${date.getFullYear()}-${String(date.getMonth()+1).padStart(2,'0')}-${String(date.getDate()).padStart(2,'0')}`
      firstWeek.push({
        date: ds,
        filled: filledSet.has(ds),
        isToday: ds === `${today.getFullYear()}-${String(today.getMonth()+1).padStart(2,'0')}-${String(today.getDate()).padStart(2,'0')}`,
        label: date.getDate() === 1 ? monthNames[date.getMonth()] : '',
        dateNum: isLastDayOfMonth(date) ? String(date.getDate()) : '',
      })
      dayOffset++
    } else {
      firstWeek.push(null)
    }
  }
  weeks.push(firstWeek)

  // Remaining weeks
  while (dayOffset < totalDays) {
    if (weeks.length > 53) break
    const week = []
    for (let d = 0; d < 7; d++) {
      if (dayOffset < totalDays) {
        const date = new Date(y, 0, 1 + dayOffset)
        const ds = `${date.getFullYear()}-${String(date.getMonth()+1).padStart(2,'0')}-${String(date.getDate()).padStart(2,'0')}`
        week.push({
          date: ds,
          filled: filledSet.has(ds),
          isToday: ds === `${today.getFullYear()}-${String(today.getMonth()+1).padStart(2,'0')}-${String(today.getDate()).padStart(2,'0')}`,
          label: date.getDate() === 1 ? monthNames[date.getMonth()] : '',
          dateNum: isLastDayOfMonth(date) ? String(date.getDate()) : '',
        })
        dayOffset++
      } else {
        week.push(null)
      }
    }
    weeks.push(week)
  }
  return weeks
}

function isLastDayOfMonth(date) {
  const next = new Date(date.getFullYear(), date.getMonth() + 1, 1)
  next.setDate(next.getDate() - 1)
  return date.getDate() === next.getDate()
}

const topStats = computed(() => {
  const total = allEntries.value.length
  const now = new Date()
  const thisMonth = allEntries.value.filter(e => e.recorded_at.substring(0,7) === `${now.getFullYear()}-${String(now.getMonth()+1).padStart(2,'0')}`).length
  const daysWithRecords = new Set(allEntries.value.map(e => e.recorded_at.substring(0, 10))).size
  let daysSinceStart = 0
  if (allEntries.value.length > 0) {
    const earliest = allEntries.value.reduce((min, e) => e.recorded_at < min ? e.recorded_at : min, allEntries.value[0].recorded_at)
    const startDate = new Date(earliest.substring(0, 10))
    const todayDate = new Date()
    todayDate.setHours(0, 0, 0, 0)
    daysSinceStart = Math.floor((todayDate - startDate) / 86400000) + 1
  }
  return [
    { label: '总记录', value: total },
    { label: '本月', value: thisMonth },
    { label: '记录天数', value: daysWithRecords },
    { label: '已使用', value: daysSinceStart + '天' },
  ]
})

const categoryData = computed(() => {
  return categories.map(cat => {
    const entries = allEntries.value.filter(e => e.type === cat.type)
    return { ...cat, count: entries.length, weeks: buildHeatmapWeeks(entries, year.value) }
  })
})


async function fetchAll() {
  try {
    let allData = []
    let before = ''
    while (true) {
      const url = `${API_BASE}/api/entries?limit=100${before ? '&before=' + encodeURIComponent(before) : ''}`
      const res = await fetch(url)
      const json = await res.json()
      if (json.code !== 0 || !Array.isArray(json.data) || json.data.length === 0) break
      allData = allData.concat(json.data)
      if (json.data.length < 100) break
      before = json.data[json.data.length - 1].recorded_at
    }
    allEntries.value = allData
  } catch (e) { console.error(e) }
}

onMounted(async () => { await fetchAll(); nextTick(() => scrollHmRight()) })

watch(year, () => { nextTick(() => scrollHmRight()) })

function scrollHmRight() {
  document.querySelectorAll('.hm-scroll').forEach(el => { el.scrollLeft = el.scrollWidth })
}


</script>

<style scoped>
.stats-page { min-height: 100dvh; background: var(--color-card); padding: 0 24px 60px; width: 100%; }
@media (min-width: 600px) { .stats-page { padding: 0 28px 68px; } }
@media (min-width: 768px) { .stats-page { padding: 0 36px 72px; } }

.stats-header { display: flex; align-items: center; justify-content: space-between; padding: 14px 0 10px; }
.back-btn { width: 32px; height: 32px; border: none; background: none; color: var(--color-ink); cursor: pointer; display: flex; align-items: center; justify-content: center; border-radius: 50%; transition: background .15s }
.back-btn:hover { background: var(--color-surface-dim) }
.header-title { font-size: 17px; font-weight: 600; color: var(--color-ink); }
.header-spacer { width: 32px; }

.stat-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-bottom: 24px; }
@media (min-width: 768px) { .stat-grid { grid-template-columns: 1fr 1fr 1fr 1fr; gap: 14px; } }
.stat-card { background: var(--color-surface-dim); border-radius: var(--radius-lg); padding: 18px 14px; display: flex; flex-direction: column; gap: 3px; }
.stat-value { font-size: 34px; font-weight: 700; color: var(--color-ink); line-height: 1; }
.stat-label { font-size: 13px; color: var(--color-graphite); }

.year-nav { display: flex; align-items: center; justify-content: center; gap: 2px; margin-bottom: 16px; background: var(--color-surface-dim); border-radius: var(--radius-md); padding: 4px; width: fit-content; margin-left: auto; margin-right: auto; }
.year-arrow { width: 36px; height: 36px; border: none; border-radius: 10px; background: transparent; color: var(--color-graphite); cursor: pointer; display: flex; align-items: center; justify-content: center; transition: all 0.2s ease; }
.year-arrow:hover { background: var(--color-pencil); color: var(--color-ink); }
.year-arrow:active { background: var(--color-pencil); }
.year-arrow:disabled { color: var(--color-pencil); cursor: default; pointer-events: none; }
.year-arrow:disabled:hover { background: transparent; }
.year-text { font-size: 15px; font-weight: 600; color: var(--color-ink); padding: 0 14px; min-width: 48px; text-align: center; font-variant-numeric: tabular-nums; }

.heatmaps { display: flex; flex-direction: column; gap: 12px; }
.heatmap-card { background: var(--color-surface-dim); border-radius: var(--radius-lg); padding: 14px 10px; overflow: hidden; }
.heatmap-title-row { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; padding: 0 2px; }
.heatmap-dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; }
.heatmap-title { font-size: 14px; font-weight: 600; color: var(--color-ink); }
.heatmap-count { margin-left: auto; font-size: 12px; color: var(--color-graphite); }

.hm-scroll { display: flex; gap: 3px; overflow-x: auto; scrollbar-width: none; -webkit-overflow-scrolling: touch; scroll-behavior: smooth }
.hm-scroll::-webkit-scrollbar { display: none; }
.hm-weeks { display: flex; flex-direction: column; gap: 3px; flex-shrink: 0; }
.hm-cell { width: 22px; height: 22px; border-radius: 3px; background: #EDE8E2; position: relative; flex-shrink: 0; display: flex; align-items: center; justify-content: center; overflow: hidden; }

.hm-cell.hm-empty { background: transparent; }
.hm-date { font-size: 9px; color: var(--color-graphite); line-height: 1; font-weight: 600; text-align: center; }
.hm-month { font-size: 9px; color: var(--color-graphite); line-height: 1; font-weight: 600; text-align: center; }
</style>