<template>
  <div class="year-month-header">
    <button class="header-btn" @click.stop="toggle">
      <span class="ym-text">{{ displayText }}</span>
      <svg class="ym-chevron" :class="{ open: showDropdown }" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
        <polyline points="6 9 12 15 18 9" />
      </svg>
    </button>

    <button class="filter-btn" @click.stop="openFilter">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" :stroke="activeFilters.length < filterTypes.length ? 'var(--color-ink)' : 'var(--color-graphite)'" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/>
      </svg>
    </button>

    <transition name="fade">
      <div v-if="showDropdown" class="dropdown-overlay" @click="showDropdown = false">
        <div class="dropdown-panel" @click.stop>
          <button v-for="item in months" :key="item.value" class="dropdown-item" :class="{ active: item.value === modelValue }" @click="select(item.value)">{{ item.label }}</button>
        </div>
      </div>
    </transition>

    <transition name="fade">
      <div v-if="showFilter" class="dropdown-overlay" @click="closeFilter">
        <div class="dropdown-panel filter-panel" @click.stop>
          <div class="filter-header">
            <span class="filter-title">筛选类型</span>
            <button class="filter-clear" @click="clearFilters">清除</button>
          </div>
          <div class="filter-list">
          <button v-for="t in filterTypes" :key="t.key" class="filter-item" :class="{ on: tempFilters.includes(t.key) }" @click="toggleFilter(t.key)">
            <span class="filter-dot" :style="{ background: t.color }"></span>
            <span class="filter-label">{{ t.label }}</span>
            <svg v-if="tempFilters.includes(t.key)" class="filter-check" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#2563EB" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
            </button>
          </div>
          <div class="filter-footer">
            <button class="filter-confirm" @click="closeFilter">确定</button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  modelValue: String,
  months: Array,
  activeFilters: { type: Array, default: () => [] },
})

const emit = defineEmits(['update:modelValue', 'update:activeFilters'])

const showDropdown = ref(false)
const showFilter = ref(false)
const tempFilters = ref([...props.activeFilters])

const filterTypes = [
  { key: 'thought', label: '随记', color: '#9DB5C9' },
  { key: 'asset', label: '资产记录', color: '#D4B87A' },
  { key: 'exercise', label: '运动', color: '#84B8A4' },
  { key: 'reading', label: '读书', color: '#A099C4' },
  { key: 'movie', label: '影视', color: '#CB99B0' },
  { key: 'discipline', label: '自律', color: '#8CA4BD' },
  { key: 'uric', label: '尿酸', color: '#BE999B' },
  { key: 'nosugar', label: '禁止糖分', color: '#C88C8C' },
  { key: 'todo', label: '完成待办', color: '#6C7A89' },
]

const displayText = computed(() => {
  const found = props.months.find(m => m.value === props.modelValue)
  return found ? found.label : props.modelValue
})

function select(value) {
  emit('update:modelValue', value)
  showDropdown.value = false
}
function toggle() { showDropdown.value = !showDropdown.value }

function openFilter() {
  tempFilters.value = [...props.activeFilters]
  showFilter.value = !showFilter.value
}
function toggleFilter(key) {
  const i = tempFilters.value.indexOf(key)
  if (i >= 0) tempFilters.value.splice(i, 1)
  else tempFilters.value.push(key)
}
function clearFilters() {
  tempFilters.value = []
  emit('update:activeFilters', [])
  showFilter.value = false
}
function closeFilter() {
  emit('update:activeFilters', [...tempFilters.value])
  showFilter.value = false
}
</script>

<style scoped>
.year-month-header {
  position: fixed; top: 0; z-index: 50;
  left: 50%; transform: translateX(-50%);
  width: 100%; max-width: var(--content-width);
  background: var(--color-card);
  padding: 16px 24px 12px;
  border-bottom: 1px solid var(--color-border-light);
  backdrop-filter: blur(16px); -webkit-backdrop-filter: blur(16px);
  display: flex; align-items: center; justify-content: space-between;
}
@media (min-width: 600px) { .year-month-header { padding: 18px 28px 14px; } }
@media (min-width: 768px) { .year-month-header { padding: 20px 36px 14px; } }

.header-btn {
  display: inline-flex; align-items: center; gap: 6px;
  background: none; border: none; cursor: pointer; padding: 4px 0;
}
.ym-text { font-size: 22px; font-weight: 700; color: var(--color-ink); letter-spacing: 0.3px; }
@media (min-width: 600px) { .ym-text { font-size: 24px; } }
.ym-chevron { color: var(--color-graphite); margin-top: 2px; transition: transform .2s ease; }
.ym-chevron.open { transform: rotate(180deg); }

.filter-btn {
  width: 36px; height: 36px; border-radius: 50%; border: none;
  background: var(--color-surface-dim); color: var(--color-ink);
  display: flex; align-items: center; justify-content: center;
  cursor: pointer; transition: background .12s;
}
.filter-btn:hover { background: var(--color-pencil); }

.dropdown-overlay { position: fixed; inset: 0; z-index: 40; background: rgba(0,0,0,.12); }
.dropdown-panel {
  position: absolute; top: 56px; left: 24px; right: 24px;
  max-width: calc(var(--content-width) - 48px); margin: 0 auto;
  background: var(--color-card); border-radius: var(--radius-md);
  box-shadow: 0 12px 40px rgba(0,0,0,.08); overflow: hidden;
}
.filter-panel { padding: 16px 0 0; }

.filter-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 0 20px 12px;
}
.filter-title { font-size: 16px; font-weight: 600; color: var(--color-ink); }
.filter-clear { font-size: 14px; color: #999; background: none; border: none; cursor: pointer; }
.filter-list { max-height: 320px; overflow-y: auto; }
.filter-item {
  display: flex; align-items: center; gap: 10px;
  width: 100%; padding: 12px 20px; border: none;
  background: none; text-align: left; cursor: pointer;
  font-size: 15px; color: var(--color-ink); transition: background .1s;
}
.filter-item:hover { background: var(--color-surface-dim); }
.filter-dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; }
.filter-label { flex: 1; }
.filter-check { flex-shrink: 0; }

.filter-footer {
  padding: 12px 20px; border-top: 1px solid var(--color-border-light);
  display: flex; justify-content: center;
}
.filter-confirm {
  width: 100%; height: 42px; border-radius: 10px; border: none;
  background: var(--color-ink); color: var(--color-card);
  font-size: 15px; font-weight: 600; cursor: pointer;
  transition: opacity .15s;
}
.filter-confirm:hover { opacity: .85; }

.dropdown-item {
  display: block; width: 100%; padding: 13px 20px;
  background: none; border: none; text-align: left;
  font-size: 15px; color: var(--color-ink); cursor: pointer;
  transition: background .12s;
}
.dropdown-item:hover { background: var(--color-surface-dim); }
.dropdown-item.active { font-weight: 600; background: var(--color-surface-dim); }

.fade-enter-active, .fade-leave-active { transition: opacity .2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
