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

    <div class="period-bar">
      <button v-for="p in periods" :key="p" class="period-btn" :class="{ active: period === p }" @click="period = p">{{ periodLabels[p] }}</button>
    </div>

    <div class="year-nav">
      <button class="year-arrow" @click="year--">&lt;</button>
      <span class="year-text">{{ year }}</span>
      <button class="year-arrow" @click="year++" :disabled="year >= maxYear">&gt;</button>
    </div>

    <div class="heatmaps">
      <div class="heatmap-card" v-for="cat in categoryData" :key="cat.type">
        <div class="heatmap-title-row">
          <div class="heatmap-dot" :style="{ background: cat.color }"></div>
          <span class="heatmap-title">{{ cat.name }}</span>
          <span class="heatmap-count">{{ cat.count }}次</span>
        </div>
        <div class="heatmap-body">
          <div class="heatmap-scroll" ref="scrollRef">
            <div class="heatmap-grid" v-for="(week, wi) in cat.weeks" :key="wi">
              <div v-for="(day, di) in week" :key="di" class="heatmap-cell" :class="day ? { filled: day.filled, today: day.isToday } : {}" :style="day && day.filled ? { background: cat.color } : {}" :title="day ? day.date : ''">
                <template v-if="day">
                  <span v-if="day.label" class="cell-label">{{ day.label }}</span>
                  <span v-if="day.dateNum" class="cell-datenum">{{ day.dateNum }}</span>
                </template>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const emit = defineEmits(['back'])
const period = ref('year')
const year = ref(new Date().getFullYear())
const maxYear = ref(new Date().getFullYear())
const periods = ['week', 'month', 'year']
const periodLabels = { week: '周', month: '月', year: '年' }
const API_BASE = 'http://localhost:8080'
const allEntries = ref([])

const categories = [
  { type: 'thought',  name: '念头',     color: '#9DB5C9' },
  { type: 'asset',    name: '资产记录', color: '#D4B87A' },
  { type: 'uric',     name: '尿酸记录', color: '#BE999B' },
  { type: 'exercise', name: '运动',     color: '#84B8A4' },
  { type: 'discipline', name: '自律',   color: '#8CA4BD' },
  { type: 'nosugar',  name: '禁止糖分', color: '#C88C8C' },
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
  const today = allEntries.value.filter(e => e.recorded_at.substring(0,10) === `${now.getFullYear()}-${String(now.getMonth()+1).padStart(2,'0')}-${String(now.getDate()).padStart(2,'0')}`).length
  let maxStreak = 0
  for (const t of ['discipline','nosugar']) {
    const es = allEntries.value.filter(e => e.type === t)
    if (es.length) maxStreak = Math.max(maxStreak, parseInt(es[0].description)||0)
  }
  return [
    { label: '总记录', value: total },
    { label: '本月', value: thisMonth },
    { label: '今日', value: today },
    { label: '最长连续', value: maxStreak + '天' },
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
    const res = await fetch(`${API_BASE}/api/entries?limit=500`)
    const json = await res.json()
    if (json.code === 0) allEntries.value = json.data
  } catch (e) { console.error(e) }
}

onMounted(fetchAll)

function scrollToEnd(el) { if (el) { el.scrollLeft = el.scrollWidth; } }
</script>

<style scoped>
.stats-page { min-height: 100dvh; background: #fff; padding: 0 32px 60px; width: 100%; }
.stats-header { display: flex; align-items: center; justify-content: space-between; padding: 14px 0 10px; }
.back-btn { width: 32px; height: 32px; border: none; background: none; color: #333; cursor: pointer; display: flex; align-items: center; justify-content: center; }
.header-title { font-size: 17px; font-weight: 600; color: #1a1a1a; }
.header-spacer { width: 32px; }

.stat-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-bottom: 24px; }
.stat-card { background: #fafafa; border-radius: 14px; padding: 18px 14px; display: flex; flex-direction: column; gap: 4px; }
.stat-value { font-size: 34px; font-weight: 700; color: #1a1a1a; line-height: 1; }
.stat-label { font-size: 13px; color: #999; }

.period-bar { display: flex; margin: 8px 0 12px; background: #f5f5f5; border-radius: 8px; padding: 2px; }
.period-btn { flex: 1; padding: 6px 0; text-align: center; font-size: 13px; color: #999; background: none; border: none; border-radius: 6px; cursor: pointer; }
.period-btn.active { background: #fff; color: #1a1a1a; box-shadow: 0 0 0 1px #e0e0e0; }

.year-nav { display: flex; align-items: center; justify-content: center; gap: 16px; margin-bottom: 16px; }
.year-arrow { width: 28px; height: 28px; border: 1px solid #e0e0e0; border-radius: 50%; background: #fff; color: #666; font-size: 14px; cursor: pointer; display: flex; align-items: center; justify-content: center; }
.year-text { font-size: 15px; font-weight: 600; color: #1a1a1a; }

.heatmaps { display: flex; flex-direction: column; gap: 12px; }
.heatmap-card { background: #fafafa; border-radius: 14px; padding: 12px 8px; }
.heatmap-title-row { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; }
.heatmap-dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; }
.heatmap-title { font-size: 14px; font-weight: 600; color: #1a1a1a; }
.heatmap-count { margin-left: auto; font-size: 12px; color: #999; }

.heatmap-body { display: flex; gap: 4px; }
.heatmap-weekdays { display: flex; flex-direction: column; gap: 2px; padding-top: 1px; flex-shrink: 0; }
.heatmap-weekdays span { font-size: 9px; color: #ccc; height: 14px; line-height: 14px; text-align: right; width: 16px; }

.heatmap-scroll::-webkit-scrollbar { display: none; }
.heatmap-scroll { display: flex; gap: 4px; overflow-x: auto; cursor: grab; user-select: none; direction: rtl; scrollbar-width: none; }
.heatmap-grid { display: flex; flex-direction: column; gap: 4px; direction: ltr; }
.heatmap-cell { width: 18px; height: 18px; border-radius: 3px; background: #f0f0f0; position: relative; flex-shrink: 0; display: flex; align-items: center; justify-content: center; }
.heatmap-cell.filled { background: #ccc; }
.heatmap-cell.today { box-shadow: inset 0 0 0 1.5px #333; }
.cell-label { font-size: 6px; color: #888; line-height: 1; pointer-events: none; }
.cell-datenum { position: absolute; bottom: -1px; right: 1px; font-size: 5px; color: #bbb; line-height: 1; pointer-events: none; }
</style>