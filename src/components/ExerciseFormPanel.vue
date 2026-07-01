<template>
  <Teleport to="body">
    <Transition name="form">
      <div v-if="visible" class="form-overlay" @click.self="$emit('cancel')">
        <div class="form-sheet">
          <div class="form-navbar">
            <button class="form-nav-left" @click="$emit('cancel')">取消</button>
            <span class="form-nav-title">运动记录</span>
            <button class="form-nav-right" @click="handleCreate" :disabled="!String(amount).trim()">新建</button>
          </div>

          <div class="form-body">
            <div class="form-section">
              <label class="form-label">记录时间</label>
              <div class="form-time-picker" @click="openTimePicker">
                <span class="form-time-text">{{ displayTime }}</span>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="form-time-arrow">
                  <polyline points="9 18 15 12 9 6" />
                </svg>
              </div>
            </div>

            <div class="form-section">
              <label class="form-label">运动类型</label>
              <div class="sport-capsules">
                <button v-for="s in sports" :key="s" class="sport-capsule" :class="{ active: selectedSport === s }" @click="selectedSport = s">{{ s }}</button>
              </div>
            </div>

            <div class="form-section">
              <label class="form-label">{{ sportLabel }}</label>
              <input v-model="amount" type="number" class="form-input" :placeholder="'请输入' + sportUnit" />
            </div>
          </div>
        </div>

        <Transition name="picker">
          <div v-if="showTimeModal" class="picker-overlay" @click.self="closeTimePicker">
            <div class="picker-card">
              <div class="picker-header">
                <span class="picker-title">选择时间</span>
                <button class="picker-close" @click="closeTimePicker" aria-label="关闭">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
                    <line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" />
                  </svg>
                </button>
              </div>
              <div class="picker-body">
                <div class="picker-col"><button class="picker-arrow" @click="pickYear++">▲</button><div class="picker-value picker-value-sm">{{ pickYear }}</div><button class="picker-arrow" @click="pickYear--">▼</button><span class="picker-unit">年</span></div>
                <div class="picker-col"><button class="picker-arrow" @click="adjustMonth(1)">▲</button><div class="picker-value picker-value-sm">{{ pad(pickMonth) }}</div><button class="picker-arrow" @click="adjustMonth(-1)">▼</button><span class="picker-unit">月</span></div>
                <div class="picker-col"><button class="picker-arrow" @click="adjustDay(1)">▲</button><div class="picker-value picker-value-sm">{{ pad(pickDay) }}</div><button class="picker-arrow" @click="adjustDay(-1)">▼</button><span class="picker-unit">日</span></div>
                <div class="picker-col"><button class="picker-arrow" @click="adjustHour(1)">▲</button><div class="picker-value">{{ pad(pickHour) }}</div><button class="picker-arrow" @click="adjustHour(-1)">▼</button><span class="picker-unit">时</span></div>
                <div class="picker-col"><button class="picker-arrow" @click="adjustMinute(1)">▲</button><div class="picker-value">{{ pad(pickMinute) }}</div><button class="picker-arrow" @click="adjustMinute(-1)">▼</button><span class="picker-unit">分</span></div>
              </div>
              <button class="picker-confirm" @click="confirmTime">确定</button>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>
<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({ visible: { type: Boolean, default: false } })
const emit = defineEmits(['cancel', 'create'])

const sports = ['跑步', '骑行', '俯卧撑', '跳绳', '游泳', '徒步', '自由活动']
const selectedSport = ref('跑步')
const amount = ref('')
const showTimeModal = ref(false)
const pickYear = ref(2025)
const pickMonth = ref(6); const pickDay = ref(1); const pickHour = ref(0); const pickMinute = ref(0)
const pad = (n) => String(n).padStart(2, '0')
const currentTime = ref(new Date())

watch(() => props.visible, (val) => {
  if (val) { amount.value = ''; selectedSport.value = '跑步'; currentTime.value = new Date() }
})

const displayTime = computed(() => {
  const d = currentTime.value
  return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日 ${pad(d.getHours())}:${pad(d.getMinutes())}`
})

const sportLabel = computed(() => {
  const m = { '跑步': '跑步距离', '骑行': '骑行距离', '俯卧撑': '俯卧撑个数', '跳绳': '跳绳个数', '游泳': '游泳距离', '徒步': '徒步距离', '自由活动': '活动时长' }
  return m[selectedSport.value] || ''
})

const sportUnit = computed(() => {
  const m = { '跑步': '千米', '骑行': '千米', '俯卧撑': '个', '跳绳': '个', '游泳': '米', '徒步': '公里', '自由活动': '分钟' }
  return m[selectedSport.value] || ''
})

const daysInMonth = computed(() => new Date(pickYear.value, pickMonth.value, 0).getDate())
function openTimePicker() { const d = currentTime.value; pickYear.value = d.getFullYear(); pickMonth.value = d.getMonth()+1; pickDay.value = d.getDate(); pickHour.value = d.getHours(); pickMinute.value = d.getMinutes(); showTimeModal.value = true }
function closeTimePicker() { showTimeModal.value = false }
function confirmTime() { currentTime.value = new Date(pickYear.value, pickMonth.value-1, pickDay.value, pickHour.value, pickMinute.value, 0); showTimeModal.value = false }
function adjustMonth(delta) { let m = pickMonth.value + delta; if (m > 12) { m=1; pickYear.value++ } if (m < 1) { m=12; pickYear.value-- } pickMonth.value = m; if (pickDay.value > daysInMonth.value) pickDay.value = daysInMonth.value }
function adjustDay(delta) { let d = pickDay.value + delta; const max = daysInMonth.value; if (d > max) { d=1; adjustMonth(1) } if (d < 1) { adjustMonth(-1); d = new Date(pickYear.value, pickMonth.value-1, 0).getDate() } pickDay.value = d }
function adjustHour(delta) { pickHour.value = (pickHour.value + delta + 24) % 24 }
function adjustMinute(delta) { pickMinute.value = (pickMinute.value + delta + 60) % 60 }
function handleCreate() { if (!String(amount.value).trim()) return; emit('create', { time: currentTime.value, sport: selectedSport.value, amount: amount.value.trim() }) }
</script>


<style scoped>
.form-overlay { position: fixed; inset: 0; z-index: 1100; background: rgba(0,0,0,0.25); display: flex; align-items: flex-end; justify-content: center; }
.form-sheet { width: 100%; max-width: 480px; height: 78%; background: #f5f5f7; border-radius: 16px 16px 0 0; display: flex; flex-direction: column; overflow: hidden; box-shadow: 0 -2px 20px rgba(0,0,0,0.08); }
.form-navbar { display: flex; align-items: center; justify-content: space-between; padding: 14px 16px; background: #fff; border-bottom: 0.5px solid #e8e8e8; flex-shrink: 0; }
.form-nav-left { font-size: 15px; color: #333; background: none; border: none; cursor: pointer; padding: 6px 4px; }
.form-nav-title { font-size: 16px; font-weight: 600; color: #1a1a1a; }
.form-nav-right { font-size: 15px; color: #3b82f6; background: none; border: none; cursor: pointer; padding: 6px 4px; font-weight: 500; }
.form-nav-right:disabled { color: #c0c0c0; cursor: default; }
.form-body { flex: 1; overflow-y: auto; padding: 20px 16px 40px; -webkit-overflow-scrolling: touch; }
.form-section { margin-bottom: 24px; }
.form-label { display: block; font-size: 13px; color: #888; margin-bottom: 8px; font-weight: 500; }
.form-time-picker { display: flex; align-items: center; justify-content: center; position: relative; background: #fff; border-radius: 10px; padding: 14px 16px; cursor: pointer; border: 1px solid #eee; min-height: 48px; }
.form-time-text { font-size: 15px; color: #1a1a1a; }
.form-time-arrow { position: absolute; right: 16px; color: #bbb; }

.sport-capsules { display: flex; flex-wrap: wrap; gap: 8px; }
.sport-capsule { padding: 6px 14px; border-radius: 16px; border: 1px solid #d0d0d0; background: #fff; color: #666; font-size: 13px; cursor: pointer; transition: all 0.15s; }
.sport-capsule:hover { border-color: #84B8A4; color: #84B8A4; }
.sport-capsule.active { background: #84B8A4; color: #fff; border-color: #84B8A4; }

.form-input { width: 100%; background: #fff; border: 1px solid #eee; border-radius: 10px; padding: 14px 16px; font-size: 20px; font-weight: 600; color: #1a1a1a; text-align: center; outline: none; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", sans-serif; -webkit-appearance: none; -moz-appearance: textfield; }
.form-input::placeholder { color: #bbb; font-weight: 400; font-size: 15px; }

.picker-overlay { position: fixed; inset: 0; z-index: 1200; background: rgba(0,0,0,0.35); display: flex; align-items: center; justify-content: center; }
.picker-card { width: 340px; background: #fff; border-radius: 16px; overflow: hidden; box-shadow: 0 8px 32px rgba(0,0,0,0.12); }
.picker-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px 10px; }
.picker-title { font-size: 16px; font-weight: 600; color: #1a1a1a; }
.picker-close { width: 28px; height: 28px; border-radius: 50%; border: none; background: #f0f0f0; color: #888; display: flex; align-items: center; justify-content: center; cursor: pointer; }
.picker-body { display: flex; align-items: center; justify-content: center; padding: 16px 12px 20px; gap: 6px; }
.picker-col { display: flex; flex-direction: column; align-items: center; gap: 3px; flex: 1; max-width: 62px; }
.picker-arrow { width: 100%; height: 30px; border: none; background: #f5f5f5; border-radius: 8px; color: #888; font-size: 11px; cursor: pointer; display: flex; align-items: center; justify-content: center; }
.picker-value { font-size: 34px; font-weight: 700; color: #1a1a1a; line-height: 1.2; text-align: center; }
.picker-value-sm { font-size: 28px; }
.picker-unit { font-size: 11px; color: #999; margin-top: -2px; }
.picker-confirm { display: block; width: auto; margin: 0 auto 16px; padding: 8px 36px; border: none; background: #2c2c2c; color: #fff; font-size: 14px; font-weight: 500; border-radius: 10px; cursor: pointer; }

.picker-enter-active, .picker-leave-active { transition: opacity 0.2s ease; }
.picker-enter-active .picker-card, .picker-leave-active .picker-card { transition: transform 0.22s cubic-bezier(0.32,0.72,0,1); }
.picker-enter-from, .picker-leave-to { opacity: 0; }
.picker-enter-from .picker-card, .picker-leave-to .picker-card { transform: scale(0.92); }
.form-enter-active, .form-leave-active { transition: opacity 0.22s ease; }
.form-enter-active .form-sheet, .form-leave-active .form-sheet { transition: transform 0.28s cubic-bezier(0.32,0.72,0,1); }
.form-enter-from, .form-leave-to { opacity: 0; }
.form-enter-from .form-sheet, .form-leave-to .form-sheet { transform: translateY(100%); }
</style>