<template>
  <Teleport to="body">
    <Transition name="form">
      <div v-if="visible" class="form-overlay" @click.self="$emit('cancel')">
        <div class="form-sheet">
          <div class="form-navbar">
            <button class="form-nav-left" @click="$emit('cancel')">取消</button>
            <span class="form-nav-title">读书记录</span>
            <button class="form-nav-right" @click="handleCreate" :disabled="!name.trim() || loading">{{ loading ? '查询中...' : '新建' }}</button>
          </div>
          <div class="form-body">
            <div class="form-section">
              <label class="form-label">记录时间</label>
              <div class="form-time-picker" @click="tp.openTimePicker()">
                <span class="form-time-text">{{ tp.displayTime.value }}</span>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="form-time-arrow"><polyline points="9 18 15 12 9 6"/></svg>
              </div>
            </div>
            <div class="form-section">
              <label class="form-label">阅读状态</label>
              <div class="capsules">
                <button class="capsule" :class="{ active: status === 'reading' }" @click="status = 'reading'">在读</button>
                <button class="capsule" :class="{ active: status === 'finished' }" @click="status = 'finished'">读完</button>
              </div>
            </div>
            <div class="form-section" v-if="status === 'finished'">
              <label class="form-label">分类</label>
              <div class="cat-grid">
                <button v-for="c in categories" :key="c.key" class="catc" :class="{ active: catKey === c.key }" @click="catKey = c.key" :style="{ borderColor: c.color, color: catKey === c.key ? '#fff' : c.color, background: catKey === c.key ? c.color : 'transparent' }">{{ c.label }}</button>
              </div>
            </div>
            <div class="form-section">
              <label class="form-label">书名</label>
              <input v-model="name" class="form-input-wide" placeholder="请输入书名" />
            </div>
            <div class="form-section" v-if="status === 'finished'">
              <label class="form-label">阅读用时（小时）</label>
              <input v-model="duration" class="form-input-wide" placeholder="例如：2" />
            </div>
            <div class="form-section" v-if="status === 'finished'">
              <label class="form-label">评分（1-10分）</label>
              <div class="rating-row">
                <button v-for="n in 10" :key="n" class="star-btn" :class="{ active: rating >= n }" @click="rating = n">
                  <svg width="28" height="28" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" :fill="rating >= n ? '#D4A574' : 'var(--color-pencil)'" stroke="none"/></svg>
                </button>
                <span class="rating-text">{{ rating }}分</span>
              </div>
            </div>
          </div>
        </div>
        <TimePickerModal
          :visible="tp.showTimeModal.value"
          :pickYear="tp.pickYear.value" :pickMonth="tp.pickMonth.value" :pickDay="tp.pickDay.value"
          :pickHour="tp.pickHour.value" :pickMinute="tp.pickMinute.value"
          @close="tp.closeTimePicker()" @confirm="tp.confirmTime()"
          @adjustYear="(d) => tp.pickYear.value += d"
          @adjustMonth="tp.adjustMonth" @adjustDay="tp.adjustDay"
          @adjustHour="tp.adjustHour" @adjustMinute="tp.adjustMinute"
        />
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, watch } from 'vue'
import { API_BASE } from '../config.js'
import TimePickerModal from './TimePickerModal.vue'
import { useTimePicker } from '../composables/useTimePicker.js'

const props = defineProps({ visible: { type: Boolean, default: false } })
const emit = defineEmits(['cancel', 'create'])

const categories = [
  { key: '推理悬疑', label: '推理悬疑', color: '#8B6B4A' },
  { key: '明日方舟', label: '明日方舟', color: '#5B7A9E' },
  { key: '人物历史', label: '人物历史', color: '#9E7A5B' },
  { key: '社会科学', label: '社会科学', color: '#6B8E7A' },
  { key: '小说文学', label: '小说文学', color: '#8A6B9E' },
  { key: '自我成长', label: '自我成长', color: '#7A9E6B' },
  { key: '杂项拾遗', label: '杂项拾遗', color: '#9E8A6B' },
]

const tp = useTimePicker()
const status = ref('reading')
const catKey = ref('小说文学')
const name = ref('')
const duration = ref('')
const rating = ref(7)
const loading = ref(false)

watch(() => props.visible, (val) => {
  if (val) { status.value = 'reading'; catKey.value = '小说文学'; name.value = ''; duration.value = ''; rating.value = 7; tp.currentTime.value = new Date() }
})

async function handleCreate() {
  if (!name.value.trim()) return
  let bookInfo = null
  if (status.value === 'finished') {
    loading.value = true
    try {
      const r = await fetch(`${API_BASE}/api/book-info`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ book_name: name.value.trim() }) })
      const j = await r.json()
      if (j.code === 0 && j.data) bookInfo = j.data
    } catch (e) { console.error('Book info fetch failed:', e) }
    loading.value = false
  }
  emit('create', {
    time: tp.currentTime.value,
    status: status.value,
    name: name.value.trim(),
    duration: status.value === 'finished' ? duration.value.trim() : '',
    category: status.value === 'finished' ? catKey.value : '',
    rating: status.value === 'finished' ? rating.value : 0,
    bookInfo,
  })
}
</script>

<style scoped>
.form-overlay { position:fixed;inset:0;z-index:1100;background:rgba(0,0,0,.25);display:flex;align-items:flex-end;justify-content:center }
.form-sheet { width:100%;max-width:var(--content-width);height:80%;background:var(--color-card);border-radius:16px 16px 0 0;display:flex;flex-direction:column;overflow:hidden;box-shadow:0 -2px 20px rgba(0,0,0,.08) }
.form-navbar { display:flex;align-items:center;justify-content:space-between;padding:14px 16px;background:var(--color-card);border-bottom:1px solid var(--color-border-light);flex-shrink:0 }
.form-nav-left { font-size:15px;color:var(--color-ink);background:none;border:none;cursor:pointer;padding:6px 4px }
.form-nav-title { font-size:16px;font-weight:600;color:var(--color-ink) }
.form-nav-right { font-size:15px;color:#7BA88A;background:none;border:none;cursor:pointer;padding:6px 4px;font-weight:500 }
.form-nav-right:disabled { color:var(--color-pencil);cursor:default }
.form-body { flex:1;overflow-y:auto;padding:20px 16px 40px;-webkit-overflow-scrolling:touch }
.form-section { margin-bottom:24px }
.form-label { display:block;font-size:13px;color:var(--color-graphite);margin-bottom:8px;font-weight:500 }
.form-time-picker { display:flex;align-items:center;justify-content:center;position:relative;background:var(--color-surface-dim);border-radius:var(--radius-md);padding:14px 16px;cursor:pointer;border:1px solid var(--color-border-light);min-height:48px }
.form-time-text { font-size:15px;color:var(--color-ink) }
.capsules { display:flex;gap:10px }
.capsule { flex:1;padding:10px 0;border-radius:12px;border:1.5px solid var(--color-pencil);background:var(--color-surface-dim);color:var(--color-graphite);font-size:14px;font-weight:500;cursor:pointer;transition:all .15s;text-align:center }
.capsule:hover { border-color:#A099C4;color:#A099C4 }
.capsule.active { background:#A099C4;color:#fff;border-color:#A099C4 }
.cat-grid { display:flex;flex-wrap:wrap;gap:8px }
.catc { padding:6px 14px;border-radius:14px;border:1.5px solid;font-size:13px;font-weight:500;cursor:pointer;transition:all .15s }
.catc:hover { opacity:.8 }
.form-input-wide { width:100%;height:48px;background:var(--color-surface-dim);border:1px solid var(--color-border-light);border-radius:var(--radius-md);padding:0 16px;font-size:15px;color:var(--color-ink);outline:none;font-family:inherit;transition:border-color .15s }
.form-input-wide:focus { border-color:#A099C4 }
.form-input-wide::placeholder { color:var(--color-pencil) }
.rating-row { display:flex;align-items:center;gap:2px }
.star-btn { width:32px;height:32px;display:flex;align-items:center;justify-content:center;border:none;background:none;cursor:pointer;padding:0;transition:transform .1s }
.star-btn:hover { transform:scale(1.15) }
.star-btn:active { transform:scale(.9) }
.rating-text { font-size:15px;font-weight:600;color:var(--color-ink);margin-left:8px }
</style>
