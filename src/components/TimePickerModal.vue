<template>
  <!-- 底层半透明遮罩 -->
  <Transition name="picker">
    <div v-if="visible" class="picker-overlay" @click.self="$emit('close')">
      <!-- 底部弹出白色圆角弹窗 -->
      <div class="picker-sheet">
        <!-- 顶部导航栏：左取消 / 中标题 / 右确认 -->
        <div class="picker-nav">
          <button class="nav-cancel" @click="$emit('close')">取消</button>
          <span class="nav-title">记录时间</span>
          <button class="nav-confirm" @click="onConfirm">确认</button>
        </div>

        <div class="picker-content">
          <!-- ========== 区域1：月历选择模块 ========== -->
          <div class="calendar-section">
            <!-- ① 年月栏：左右箭头 + 居中显示 -->
            <div class="month-nav">
              <button class="month-arrow" @click="prevMonth">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><polyline points="15 18 9 12 15 6"/></svg>
              </button>
              <span class="month-label">{{ year }}年{{ month }}月</span>
              <button class="month-arrow" @click="nextMonth">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><polyline points="9 18 15 12 9 6"/></svg>
              </button>
            </div>

            <!-- ② 星期表头：7等分 -->
            <div class="weekday-row">
              <span v-for="w in weekdays" :key="w" class="weekday">{{ w }}</span>
            </div>

            <!-- ③ 7列网格日期布局 -->
            <div class="date-grid">
              <div class="date-cell" v-for="(d, i) in calendarDays" :key="i">
                <button
                  v-if="d"
                  class="date-btn"
                  :class="{ selected: d === day }"
                  @click="selectDay(d)"
                >{{ d }}</button>
              </div>
            </div>
          </div>

          <!-- ========== 区域2：时分双列滚动选择器 ========== -->
          <div class="time-section">
            <div class="time-picker-row">
              <!-- 小时列 -->
              <div class="time-col" ref="hourCol" @scroll="onHourScroll">
                <div class="time-scroll-spacer"></div>
                <div
                  v-for="h in hours" :key="h"
                  class="time-item"
                  :class="{ active: h === hour }"
                  @click="selectHour(h)"
                  :style="hourItemStyle(h)"
                >{{ pad(h) }}</div>
                <div class="time-scroll-spacer"></div>
              </div>

              <!-- 分隔符 -->
              <span class="time-sep">:</span>

              <!-- 分钟列 -->
              <div class="time-col" ref="minCol" @scroll="onMinScroll">
                <div class="time-scroll-spacer"></div>
                <div
                  v-for="m in minutes" :key="m"
                  class="time-item"
                  :class="{ active: m === minute }"
                  @click="selectMinute(m)"
                  :style="minItemStyle(m)"
                >{{ pad(m) }}</div>
                <div class="time-scroll-spacer"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'

const props = defineProps({
  visible: { type: Boolean, default: false },
  pickYear: { type: Number, default: 2026 },
  pickMonth: { type: Number, default: 7 },
  pickDay: { type: Number, default: 12 },
  pickHour: { type: Number, default: 14 },
  pickMinute: { type: Number, default: 30 },
})

const emit = defineEmits(['close', 'confirm', 'adjustYear', 'adjustMonth', 'adjustDay', 'adjustHour', 'adjustMinute'])

const weekdays = ['一', '二', '三', '四', '五', '六', '日'].map(d => '周' + d)

const year = ref(props.pickYear)
const month = ref(props.pickMonth)
const day = ref(props.pickDay)
const hour = ref(props.pickHour)
const minute = ref(props.pickMinute)

// 打开时从 props 同步
import { watch } from 'vue'
watch(() => props.visible, (v) => {
  if (v) {
    year.value = props.pickYear
    month.value = props.pickMonth
    day.value = props.pickDay
    hour.value = props.pickHour
    minute.value = props.pickMinute
  }
})

const hours = Array.from({ length: 24 }, (_, i) => i)
const minutes = Array.from({ length: 60 }, (_, i) => i)

const hourCol = ref(null)
const minCol = ref(null)

// 日历：当月天数 + 前置空位
const totalDays = computed(() => new Date(year.value, month.value, 0).getDate())
const firstDayOfWeek = computed(() => {
  const d = new Date(year.value, month.value - 1, 1).getDay()
  return d === 0 ? 6 : d - 1
})

const calendarDays = computed(() => {
  const days = []
  for (let i = 0; i < firstDayOfWeek.value; i++) days.push(null)
  for (let d = 1; d <= totalDays.value; d++) days.push(d)
  return days
})

function prevMonth() {
  const oldMonth = month.value
  if (month.value === 1) { month.value = 12; year.value--; emit('adjustYear', -1) }
  else month.value--
  emit('adjustMonth', month.value - oldMonth)
}
function nextMonth() {
  const oldMonth = month.value
  if (month.value === 12) { month.value = 1; year.value++; emit('adjustYear', 1) }
  else month.value++
  emit('adjustMonth', month.value - oldMonth)
}

function selectDay(d) {
  const oldDay = day.value
  day.value = d
  emit('adjustDay', d - oldDay)
}

function selectHour(h) {
  const old = hour.value
  hour.value = h
  emit('adjustHour', h - old)
  nextTick(() => { if (hourCol.value) hourCol.value.scrollTo({ top: h * 40, behavior: 'smooth' }) })
}
function selectMinute(m) {
  const old = minute.value
  minute.value = m
  emit('adjustMinute', m - old)
  nextTick(() => { if (minCol.value) minCol.value.scrollTo({ top: m * 40, behavior: 'smooth' }) })
}

// 时分滚动视觉：中心高亮
function hourItemStyle(h) {
  const diff = Math.abs(h - hour.value)
  const opacity = diff === 0 ? 1 : diff === 1 ? 0.6 : diff === 2 ? 0.3 : 0.15
  const scale = diff === 0 ? 1.1 : 1
  return { opacity, transform: `scale(${scale})` }
}
function minItemStyle(m) {
  const diff = Math.abs(m - minute.value)
  const opacity = diff === 0 ? 1 : diff === 1 ? 0.6 : diff === 2 ? 0.3 : 0.15
  const scale = diff === 0 ? 1.1 : 1
  return { opacity, transform: `scale(${scale})` }
}

let hourScrollTimer = null
function onHourScroll() {
  if (hourScrollTimer) clearTimeout(hourScrollTimer)
  hourScrollTimer = setTimeout(() => {
    if (!hourCol.value) return
    const idx = Math.round(hourCol.value.scrollTop / 40)
    selectHour(Math.max(0, Math.min(23, idx)))
  }, 80)
}
let minScrollTimer = null
function onMinScroll() {
  if (minScrollTimer) clearTimeout(minScrollTimer)
  minScrollTimer = setTimeout(() => {
    if (!minCol.value) return
    const idx = Math.round(minCol.value.scrollTop / 40)
    selectMinute(Math.max(0, Math.min(59, idx)))
  }, 80)
}

function scrollToSelected() {
  nextTick(() => {
    if (hourCol.value) hourCol.value.scrollTop = hour.value * 40
    if (minCol.value) minCol.value.scrollTop = minute.value * 40
  })
}

function pad(n) { return String(n).padStart(2, '0') }

function onConfirm() {
  emit('confirm')
}
</script>

<style scoped>
/* ===== 遮罩 ===== */
.picker-overlay {
  position: fixed; inset: 0; z-index: 1200;
  background: rgba(0,0,0,0.35);
  display: flex; align-items: flex-end; justify-content: center;
}

/* ===== 底部弹窗 ===== */
.picker-sheet {
  width: 100%; max-width: var(--content-width);
  background: #fff;
  border-radius: 20px 20px 0 0;
  max-height: 85vh;
  display: flex; flex-direction: column;
  overflow: hidden;
  box-shadow: 0 -4px 30px rgba(0,0,0,0.1);
  padding-bottom: env(safe-area-inset-bottom, 12px);
}

/* ===== 导航栏 ===== */
.picker-nav {
  display: flex; align-items: center; justify-content: space-between;
  padding: 16px 20px 8px;
  flex-shrink: 0;
}
.nav-cancel {
  font-size: 16px; color: #999; background: none; border: none; cursor: pointer; padding: 4px 0;
}
.nav-title {
  font-size: 17px; font-weight: 700; color: #1a1a1a;
}
.nav-confirm {
  font-size: 16px; color: #2563EB; font-weight: 600; background: none; border: none; cursor: pointer; padding: 4px 0;
}

.picker-content {
  flex: 1; overflow-y: auto; padding: 0 16px 20px;
}

/* ===== 日历区域 ===== */
.calendar-section {
  margin-bottom: 16px;
}

/* 年月栏 */
.month-nav {
  display: flex; align-items: center; justify-content: center; gap: 16px;
  padding: 8px 0 12px;
}
.month-label {
  font-size: 17px; font-weight: 600; color: #1a1a1a; min-width: 120px; text-align: center;
}
.month-arrow {
  width: 36px; height: 36px; border-radius: 50%; border: none;
  background: #f3f4f6; color: #555;
  display: flex; align-items: center; justify-content: center;
  cursor: pointer; transition: background .15s;
}
.month-arrow:hover { background: #e5e7eb; }

/* 星期表头 */
.weekday-row {
  display: grid; grid-template-columns: repeat(7, 1fr);
  padding: 4px 0 8px;
}
.weekday {
  font-size: 13px; color: #b0b0b0; text-align: center; font-weight: 500;
}

/* 日期网格 */
.date-grid {
  display: grid; grid-template-columns: repeat(7, 1fr);
  gap: 2px;
}
.date-cell {
  display: flex; align-items: center; justify-content: center;
  height: 40px;
}
.date-btn {
  width: 36px; height: 36px; border-radius: 50%; border: none;
  background: transparent; color: #1a1a1a;
  font-size: 15px; font-weight: 500; cursor: pointer;
  transition: all .15s;
}
.date-btn:hover { background: #eef2ff; }
.date-btn.selected {
  background: #2563EB; color: #fff; font-weight: 700;
}

/* ===== 时分选择器 ===== */
.time-section {
  border-top: 1px solid #f0f0f0;
  padding-top: 16px;
}
.time-picker-row {
  display: flex; align-items: center; justify-content: center;
  gap: 4px; position: relative;
}
.time-sep {
  font-size: 28px; font-weight: 700; color: #1a1a1a;
  margin: 0 -2px; z-index: 2;
}

/* 滚轮列 */
.time-col {
  flex: 1; max-width: 120px;
  height: 200px; overflow-y: auto;
  scroll-snap-type: y mandatory;
  -ms-overflow-style: none; scrollbar-width: none;
  mask-image: linear-gradient(to bottom, transparent 0%, black 30%, black 70%, transparent 100%);
  -webkit-mask-image: linear-gradient(to bottom, transparent 0%, black 30%, black 70%, transparent 100%);
}
.time-col::-webkit-scrollbar { display: none; }
.time-scroll-spacer { height: 80px; }

.time-item {
  height: 40px; display: flex; align-items: center; justify-content: center;
  font-size: 22px; font-weight: 600; color: #1a1a1a;
  scroll-snap-align: center; cursor: pointer;
  transition: opacity .2s, transform .2s;
  font-variant-numeric: tabular-nums;
}
.time-item.active {
  font-size: 26px; font-weight: 700; color: #2563EB;
}

/* ===== 动画 ===== */
.picker-enter-active, .picker-leave-active { transition: opacity 0.25s ease; }
.picker-enter-active .picker-sheet, .picker-leave-active .picker-sheet { transition: transform 0.3s cubic-bezier(0.32,0.72,0,1); }
.picker-enter-from, .picker-leave-to { opacity: 0; }
.picker-enter-from .picker-sheet, .picker-leave-to .picker-sheet { transform: translateY(100%); }
</style>
