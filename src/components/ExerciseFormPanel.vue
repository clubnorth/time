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
        <TimePickerModal :visible="showTimeModal" :pickYear="pickYear" :pickMonth="pickMonth" :pickDay="pickDay" :pickHour="pickHour" :pickMinute="pickMinute" @close="closeTimePicker" @confirm="confirmTime" @adjustYear="(d) => pickYear += d" @adjustMonth="adjustMonth" @adjustDay="adjustDay" @adjustHour="adjustHour" @adjustMinute="adjustMinute" />
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import TimePickerModal from './TimePickerModal.vue'
import { useTimePicker } from '../composables/useTimePicker.js'

const props = defineProps({ visible: { type: Boolean, default: false } })
const emit = defineEmits(['cancel', 'create'])
const { currentTime, showTimeModal, pickYear, pickMonth, pickDay, pickHour, pickMinute, displayTime, openTimePicker, closeTimePicker, confirmTime, adjustMonth, adjustDay, adjustHour, adjustMinute } = useTimePicker()

const sports = ['跑步', '骑行', '俯卧撑', '跳绳', '游泳', '徒步', '自由活动']
const selectedSport = ref('跑步')
const amount = ref('')

watch(() => props.visible, (val) => { if (val) { amount.value = ''; selectedSport.value = '跑步'; currentTime.value = new Date() } })

const sportLabel = computed(() => {
  const m = { '跑步': '跑步距离', '骑行': '骑行距离', '俯卧撑': '俯卧撑个数', '跳绳': '跳绳个数', '游泳': '游泳距离', '徒步': '徒步距离', '自由活动': '活动时长' }
  return m[selectedSport.value] || ''
})

const sportUnit = computed(() => {
  const m = { '跑步': '千米', '骑行': '千米', '俯卧撑': '个', '跳绳': '个', '游泳': '米', '徒步': '公里', '自由活动': '分钟' }
  return m[selectedSport.value] || ''
})

function handleCreate() {
  if (!String(amount.value).trim()) return
  emit('create', { time: currentTime.value, sport: selectedSport.value, amount: String(amount.value).trim() })
}
</script>

<style scoped>
.form-overlay { position: fixed; inset: 0; z-index: 1100; background: rgba(0,0,0,0.25); display: flex; align-items: flex-end; justify-content: center; }
.form-sheet { width: 100%; max-width: var(--content-width); height: 78%; background: var(--color-card); border-radius: 16px 16px 0 0; display: flex; flex-direction: column; overflow: hidden; box-shadow: 0 -2px 20px rgba(0,0,0,0.08); }
.form-navbar { display: flex; align-items: center; justify-content: space-between; padding: 14px 16px; background: var(--color-card); border-bottom: 1px solid var(--color-border-light); flex-shrink: 0; }
.form-nav-left { font-size: 15px; color: var(--color-ink); background: none; border: none; cursor: pointer; padding: 6px 4px; }
.form-nav-title { font-size: 16px; font-weight: 600; color: var(--color-ink); }
.form-nav-right { font-size: 15px; color: #7BA88A; background: none; border: none; cursor: pointer; padding: 6px 4px; font-weight: 500; }
.form-nav-right:disabled { color: var(--color-pencil); cursor: default; }
.form-body { flex: 1; overflow-y: auto; padding: 20px 16px 40px; -webkit-overflow-scrolling: touch; }
.form-section { margin-bottom: 24px; }
.form-label { display: block; font-size: 13px; color: var(--color-graphite); margin-bottom: 8px; font-weight: 500; }
.form-time-picker { display: flex; align-items: center; justify-content: center; position: relative; background: var(--color-surface-dim); border-radius: var(--radius-md); padding: 14px 16px; cursor: pointer; border: 1px solid var(--color-border-light); min-height: 48px; }
.form-time-text { font-size: 15px; color: var(--color-ink); }
.form-time-arrow { position: absolute; right: 16px; color: var(--color-pencil); }
.sport-capsules { display: flex; flex-wrap: wrap; gap: 8px; }
.sport-capsule { padding: 6px 14px; border-radius: 16px; border: 1px solid var(--color-pencil); background: var(--color-surface-dim); color: var(--color-graphite); font-size: 13px; cursor: pointer; transition: all 0.15s; }
.sport-capsule:hover { border-color: #84B8A4; color: #84B8A4; }
.sport-capsule.active { background: #84B8A4; color: #fff; border-color: #84B8A4; }
.form-input { width: 100%; background: var(--color-surface-dim); border: 1px solid var(--color-border-light); border-radius: var(--radius-md); padding: 14px 16px; font-size: 20px; font-weight: 600; color: var(--color-ink); text-align: center; outline: none; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", sans-serif; -webkit-appearance: none; -moz-appearance: textfield; }
.form-input::placeholder { color: var(--color-pencil); font-weight: 400; font-size: 15px; }
</style>