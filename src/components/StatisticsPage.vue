<template>
  <div class="stats-page">
    <!-- Header -->
    <div class="stats-header">
      <button class="back-btn" @click="$emit('back')">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="15 18 9 12 15 6"/></svg>
      </button>
      <span class="header-title">统计</span>
      <span class="header-spacer"></span>
    </div>

    <!-- Period selector -->
    <div class="stat-grid">
      <div class="stat-card" v-for="s in topStats" :key="s.label">
        <span class="stat-value">{{ s.value }}</span>
        <span class="stat-label">{{ s.label }}</span>
      </div>
    </div>

    <!-- Top stat cards 2x2 -->
    <div class="period-bar">
      <button v-for="p in periods" :key="p" class="period-btn" :class="{ active: period === p }" @click="period = p">{{ periodLabels[p] }}</button>
    </div>

    <!-- Category heatmap cards -->
    <div class="heatmaps">
      <div class="heatmap-card" v-for="cat in categories" :key="cat.type">
        <div class="heatmap-title-row">
          <div class="heatmap-dot" :style="{ background: cat.color }"></div>
          <span class="heatmap-title">{{ cat.name }}</span>
          <span class="heatmap-count">{{ cat.count }}次</span>
        </div>
        <div class="heatmap-grid">
          <div v-for="(day, i) in cat.days" :key="i"
            class="heatmap-cell"
            :class="{ filled: day.filled }"
            :style="{ background: day.filled ? cat.color : '#f0f0f0' }"
            :title="day.date"
          ></div>
        </div>
        <div class="heatmap-months">
          <span v-for="m in 12" :key="m" class="heatmap-month">{{ m }}月</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const emit = defineEmits(['back'])

const period = ref('year')
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

const topStats = computed(() => {
  const total = allEntries.value.length
  const now = new Date()
  const thisMonth = allEntries.value.filter(e => {
    const m = e.recorded_at.substring(0, 7)
    return m === `${now.getFullYear()}-${String(now.getMonth()+1).padStart(2,'0')}`
  }).length
  const today = allEntries.value.filter(e => {
    return e.recorded_at.substring(0, 10) === `${now.getFullYear()}-${String(now.getMonth()+1).padStart(2,'0')}-${String(now.getDate()).padStart(2,'0')}`
  }).length
  // Find max consecutive days from discipline/nosugar
  let maxStreak = 0
  for (const type of ['discipline', 'nosugar']) {
    const entries = allEntries.value.filter(e => e.type === type)
    if (entries.length > 0) {
      const n = parseInt(entries[0].description) || 0
      if (n > maxStreak) maxStreak = n
    }
  }
  return [
    { label: '总记录', value: total },
    { label: '本月', value: thisMonth },
    { label: '今日', value: today },
    { label: '最长连续', value: maxStreak + '天' },
  ]
})

const categoryData = computed(() => {
  const now = new Date()
  const year = now.getFullYear()
  const daysInYear = (new Date(year, 11, 31).getTime() - new Date(year, 0, 1).getTime()) / 86400000 + 1

  return categories.map(cat => {
    const entries = allEntries.value.filter(e => e.type === cat.type)
    const count = entries.length
    const filledDates = new Set(entries.map(e => e.recorded_at.substring(0, 10)))

    const days = []
    const start = new Date(year, 0, 1)
    for (let i = 0; i < daysInYear; i++) {
      const d = new Date(start.getTime() + i * 86400000)
      const dateStr = `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')}`
      days.push({ date: dateStr, filled: filledDates.has(dateStr) })
    }
    return { ...cat, count, days }
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
</script>

<style scoped>
.stats-page {
  min-height: 100dvh;
  background: #fff;
  padding: 0 16px 60px;
  width: 100%;
}

.stats-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 0 8px;
}
.back-btn { width: 32px; height: 32px; border: none; background: none; color: #333; cursor: pointer; display: flex; align-items: center; justify-content: center; }
.header-title { font-size: 17px; font-weight: 600; color: #1a1a1a; }
.header-spacer { width: 32px; }

.period-bar { display: flex; gap: 0; margin: 8px 0 16px; background: #f5f5f5; border-radius: 8px; padding: 2px; }
.period-btn { flex: 1; padding: 6px 0; text-align: center; font-size: 13px; color: #999; background: none; border: none; border-radius: 6px; cursor: pointer; }
.period-btn.active { background: #fff; color: #1a1a1a; box-shadow: 0 0 0 1px #e0e0e0; }

.stat-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-bottom: 24px; }
.stat-card { background: #fafafa; border-radius: 14px; padding: 18px 14px; display: flex; flex-direction: column; gap: 4px; }
.stat-value { font-size: 34px; font-weight: 700; color: #1a1a1a; line-height: 1; }
.stat-label { font-size: 13px; color: #999; }

.heatmaps { display: flex; flex-direction: column; gap: 12px; }
.heatmap-card { background: #fafafa; border-radius: 14px; padding: 14px 36px; }
.heatmap-title-row { display: flex; align-items: center; gap: 8px; margin-bottom: 10px; }
.heatmap-dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; }
.heatmap-title { font-size: 14px; font-weight: 600; color: #1a1a1a; }
.heatmap-count { margin-left: auto; font-size: 12px; color: #999; }

.heatmap-grid { display: flex; flex-wrap: wrap; gap: 2px; margin-bottom: 4px; }
.heatmap-cell { width: 11px; height: 11px; border-radius: 3px; }
.heatmap-cell.filled { opacity: 0.9; }

.heatmap-months { display: flex; justify-content: space-between; padding: 0 1px; }
.heatmap-month { font-size: 9px; color: #ccc; }
</style>