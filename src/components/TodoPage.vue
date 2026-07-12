<template>
  <div class="todo-page">
    <div class="todo-header">
      <button class="back-btn" @click="$emit('back')">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="15 18 9 12 15 6"/></svg>
      </button>
      <span class="header-title">待办事项</span>
      <span class="header-spacer"></span>
    </div>

    <div class="todo-content">
      <div v-for="cat in visibleCategories" :key="cat.key" class="category-section">
        <div class="category-header" :style="{ background: cat.color }">
          <span class="category-label">{{ cat.label }}</span>
          <span class="category-count">{{ getItems(cat.key).length }}</span>
        </div>
        <div class="category-list">
          <div v-for="todo in getItems(cat.key)" :key="todo.id" class="todo-card-wrapper">
            <div
              class="todo-card"
              :class="{ swiped: swipeId === todo.id }"
              :style="{ transform: swipeId === todo.id ? 'translateX(' + swipeOffset + 'px)' : 'translateX(0)' }"
              @pointerdown.prevent="onDown($event, todo)"
              @pointermove.prevent="onMove($event, todo)"
              @pointerup="onUp()"
              @pointerleave="onUp()"
              @pointercancel="onUp()"
            >
              <button class="todo-check" @click.stop="toggleComplete(todo)" :class="{ done: todo.completed }">
                <svg v-if="todo.completed" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
              </button>
              <div class="todo-info">
                <span class="todo-time">{{ fmtTime(todo.due_date, todo.created_at) }}</span>
                <span class="todo-title">{{ todo.title }}</span>
              </div>
            </div>
            <div class="todo-actions">
              <button v-if="!todo.completed" class="action-btn action-edit" @click.stop="openEdit(todo)">修改</button>
              <button class="action-btn action-delete" @click.stop="deleteTodo(todo)">删除</button>
            </div>
          </div>
        </div>
      </div>

      <div v-if="allActive.length === 0" class="todo-empty">暂无待办事项</div>
    </div>

    <!-- Edit form -->
    <Teleport to="body">
      <Transition name="form">
        <div v-if="editingId" class="form-overlay" @click.self="editingId = null">
          <div class="form-sheet">
            <div class="form-navbar">
              <button class="form-nav-left" @click="editingId = null">取消</button>
              <span class="form-nav-title">修改待办</span>
              <button class="form-nav-right" @click="saveEdit" :disabled="!editTitle.trim() || editTitle.length > 50">保存</button>
            </div>
            <div class="form-body">
              <div class="form-section">
                <label class="form-label">截止时间</label>
                <div class="form-time-picker" @click="etp.openTimePicker()">
                  <span class="form-time-text">{{ editPickerTime }}</span>
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="form-time-arrow">
                    <polyline points="9 18 15 12 9 6" />
                  </svg>
                </div>
              </div>
              <div class="form-section">
                <label class="form-label">内容</label>
                <textarea v-model="editTitle" class="form-textarea" placeholder="待办内容" maxlength="55" rows="3"></textarea>
                <span class="form-count" :class="{ warn: editTitle.length > 50 }">{{ editTitle.length }}/50</span>
              </div>
            </div>
          </div>
          <TodoTimePicker
            :visible="etp.showTimeModal.value"
            :pickYear="etp.pickYear.value" :pickMonth="etp.pickMonth.value" :pickDay="etp.pickDay.value"
            :pickHour="etp.pickHour.value" :pickMinute="etp.pickMinute.value"
            @close="etp.closeTimePicker()"
            @confirm="etp.confirmTime()"
            @adjustYear="(d) => etp.pickYear.value += d"
            @adjustMonth="etp.adjustMonth"
            @adjustDay="etp.adjustDay"
            @adjustHour="etp.adjustHour"
            @adjustMinute="etp.adjustMinute"
          />
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { API_BASE } from '../config.js'
import TodoTimePicker from './TodoTimePicker.vue'
import { useTimePicker } from '../composables/useTimePicker.js'

defineEmits(['back'])

const categories = [
  { key: 'overdue',  label: '你忘记做了', color: '#C0392B' },
  { key: 'today',    label: '今日',       color: '#E65C4A' },
  { key: 'upcoming', label: '未来三日',   color: '#E88278' },
  { key: 'later',    label: '以后',       color: '#D9A09A' },
]

const visibleCategories = computed(() => categories.filter(c => getItems(c.key).length > 0))

const todos = ref([])
const allActive = computed(() => todos.value.filter(t => !t.completed))

function getItems(key) {
  const today = new Date().toISOString().substring(0, 10)
  return allActive.value.filter(t => {
    if (t.due_date) {
      const d = t.due_date.substring(0, 10)
      if (d < today) return key === 'overdue'
      if (d === today) return key === 'today'
      const max = new Date(today); max.setDate(max.getDate() + 3)
      if (d <= max.toISOString().substring(0, 10)) return key === 'upcoming'
      return key === 'later'
    }
    // No due_date: use stored category for old todos
    return (t.category || 'today') === key
  })
}

function fmtTime(d, fallback) {
  const v = d || fallback || ''
  if (!v) return ''
  return v.substring(0, 10).replace(/-/g, '/') + ' ' + v.substring(11, 16)
}

// Swipe
const swipeId = ref(null); const swipeOffset = ref(0); let swipeStartX = 0
function onDown(e, todo) {
  if (swipeId.value && swipeId.value !== todo.id) { swipeId.value = null; swipeOffset.value = 0; return }
  swipeId.value = todo.id; swipeStartX = e.clientX; swipeOffset.value = 0
}
function onMove(e, todo) {
  if (swipeId.value !== todo.id) return
  const dx = e.clientX - swipeStartX
  if (dx < 0) swipeOffset.value = Math.max(dx, -140)
}
function onUp() {
  if (swipeOffset.value < -50) swipeOffset.value = -140
  else { swipeId.value = null; swipeOffset.value = 0 }
}

async function toggleComplete(todo) {
  try {
    await fetch(API_BASE + '/api/todos/' + todo.id, {
      method: 'PUT', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ completed: !todo.completed }),
    })
    todo.completed = !todo.completed
    if (todo.completed) {
      const now = new Date()
      const at = now.getFullYear()+'-'+String(now.getMonth()+1).padStart(2,'0')+'-'+String(now.getDate()).padStart(2,'0')+' '+String(now.getHours()).padStart(2,'0')+':'+String(now.getMinutes()).padStart(2,'0')+':00'
      await fetch(API_BASE + '/api/entries', {
        method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ type: 'todo', title: '完成待办', description: todo.title, category: 'green', recorded_at: at }),
      })
    }
  } catch (e) { console.error(e) }
}

const etp = useTimePicker()

// 编辑表单显示时间：直接从 picker 状态实时读取
const editPickerTime = computed(() => {
  const y = etp.pickYear.value
  const m = String(etp.pickMonth.value).padStart(2, '0')
  const d = String(etp.pickDay.value).padStart(2, '0')
  const h = String(etp.pickHour.value).padStart(2, '0')
  const mi = String(Math.round(etp.pickMinute.value / 10) * 10).padStart(2, '0')
  return `${y}年${m}月${d}日 ${h}:${mi}`
})

const editingId = ref(null); const editTitle = ref('')
function openEdit(todo) {
  editingId.value = todo.id; editTitle.value = todo.title
  swipeId.value = null; swipeOffset.value = 0
  if (todo.due_date) {
    const parts = todo.due_date.split(/[- :]/)
    const dt = new Date(+parts[0], +parts[1]-1, +parts[2], +parts[3], +parts[4], 0)
    etp.currentTime.value = dt
    // 同步 picker 状态，保证 display 与选择器一致
    etp.pickYear.value = dt.getFullYear()
    etp.pickMonth.value = dt.getMonth() + 1
    etp.pickDay.value = dt.getDate()
    etp.pickHour.value = dt.getHours()
    etp.pickMinute.value = dt.getMinutes()
  } else {
    const dt = new Date()
    etp.currentTime.value = dt
    etp.pickYear.value = dt.getFullYear()
    etp.pickMonth.value = dt.getMonth() + 1
    etp.pickDay.value = dt.getDate()
    etp.pickHour.value = dt.getHours()
    etp.pickMinute.value = dt.getMinutes()
  }
  nextTick(() => document.querySelector('.form-textarea')?.focus())
}
async function saveEdit() {
  if (!editTitle.value.trim() || editTitle.value.length > 50) return
  const todo = todos.value.find(t => t.id === editingId.value)
  // 直接从 picker 状态读值，避免 confirm 未触发导致 currentTime 过期
  const y = etp.pickYear.value, mo = etp.pickMonth.value, da = etp.pickDay.value
  const h = etp.pickHour.value; let mi = Math.round(etp.pickMinute.value / 10) * 10
  const dueDate = `${y}-${String(mo).padStart(2,'0')}-${String(da).padStart(2,'0')} ${String(h).padStart(2,'0')}:${String(mi).padStart(2,'0')}:00`
  try {
    await fetch(API_BASE + '/api/todos/' + editingId.value, {
      method: 'PUT', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: editTitle.value.trim(), due_date: dueDate }),
    })
    if (todo) { todo.title = editTitle.value.trim(); todo.due_date = dueDate }
  } catch (e) { console.error(e) }
  editingId.value = null
}

async function deleteTodo(todo) {
  if (!confirm('确定删除"' + todo.title + '"吗？')) return
  try { await fetch(API_BASE + '/api/todos/' + todo.id, { method: 'DELETE' }); todos.value = todos.value.filter(t => t.id !== todo.id); swipeId.value = null; swipeOffset.value = 0 } catch (e) { console.error(e) }
}

onMounted(async () => {
  try { const r = await fetch(API_BASE + '/api/todos'); const j = await r.json(); if (j.code===0) todos.value = j.data } catch (e) { console.error(e) }
})
</script>

<style scoped>
.todo-page { min-height:100dvh;background:var(--color-card);padding:0 24px 80px;width:100% }
@media (min-width:600px) { .todo-page { padding:0 28px 88px } }
@media (min-width:768px) { .todo-page { padding:0 36px 92px } }

.todo-header { display:flex;align-items:center;justify-content:space-between;padding:14px 0 10px;position:sticky;top:0;background:var(--color-card);z-index:5 }
.back-btn { width:32px;height:32px;border:none;background:none;color:var(--color-ink);cursor:pointer;display:flex;align-items:center;justify-content:center;border-radius:50%;transition:background .15s }
.back-btn:hover { background:var(--color-surface-dim) }
.header-title { font-size:17px;font-weight:600;color:var(--color-ink) }
.header-spacer { width:32px }

.todo-content { margin-top:8px }
.category-section { margin-bottom:20px }
.category-header { display:inline-flex;align-items:center;justify-content:space-between;gap:6px;padding:6px 16px;border-radius:20px;margin-bottom:8px;min-width:100px }
.category-label { font-size:13px;font-weight:600;color:#fff }
.category-count { font-size:12px;color:rgba(255,255,255,.7);font-weight:600 }
.category-list { display:flex;flex-direction:column;gap:4px }

.todo-card-wrapper { position:relative;overflow:hidden;border-radius:var(--radius-md) }
.todo-card { position:relative;z-index:2;display:flex;align-items:center;gap:12px;padding:14px 16px;background:var(--color-surface-dim);transition:transform .2s ease }
.todo-card.swiped { box-shadow:-2px 0 8px rgba(0,0,0,.04) }
.todo-actions { position:absolute;top:0;right:0;bottom:0;display:flex;align-items:stretch;z-index:1 }
.action-btn { width:70px;border:none;font-size:13px;font-weight:600;color:#fff;cursor:pointer;display:flex;align-items:center;justify-content:center }
.action-delete { background:#D4787A }
.action-edit { background:#7B9FC6 }

.todo-check { width:22px;height:22px;border-radius:50%;border:2px solid var(--color-pencil);background:none;cursor:pointer;display:flex;align-items:center;justify-content:center;flex-shrink:0;transition:all .15s }
.todo-check:hover { border-color:#7BA88A }
.todo-check.done { border-color:#7BA88A;background:#EBF5EE;color:#7BA88A }
.todo-info { flex:1;min-width:0 }
.todo-time { font-size:12px;color:var(--color-graphite);display:block }
.todo-title { font-size:15px;color:var(--color-ink);line-height:1.4;word-break:break-word;display:block }
.todo-empty { text-align:center;padding:60px 0;color:var(--color-pencil);font-size:14px }

.form-overlay { position:fixed;inset:0;z-index:1100;background:rgba(0,0,0,.25);display:flex;align-items:flex-end;justify-content:center }
.form-sheet { width:100%;max-width:var(--content-width);height:60%;background:var(--color-card);border-radius:16px 16px 0 0;display:flex;flex-direction:column;overflow:hidden;box-shadow:0 -2px 20px rgba(0,0,0,.08) }
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
.form-time-arrow { position:absolute;right:16px;color:var(--color-pencil);flex-shrink:0 }
.form-textarea { width:100%;background:var(--color-surface-dim);border:1px solid var(--color-border-light);border-radius:var(--radius-md);padding:14px 16px;font-size:15px;color:var(--color-ink);line-height:1.6;resize:none;outline:none;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI","PingFang SC",sans-serif;-webkit-appearance:none;min-height:100px }
.form-textarea::placeholder { color:var(--color-pencil) }
.form-count { display:block;text-align:right;font-size:12px;color:var(--color-pencil);margin-top:6px }
.form-count.warn { color:#D4787A;font-weight:600 }
</style>
