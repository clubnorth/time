<template>
  <Teleport to="body">
    <Transition name="panel">
      <div v-if="visible" class="panel-overlay" @click.self="$emit('close')">
        <div class="panel-sheet">
          <div class="panel-header">
<span class="panel-title" v-if="!sortMode">添加记录</span>
<span class="panel-title" v-else>拖拽排序</span>
<div class="panel-header-actions">
  <template v-if="!sortMode">
  <button class="panel-sort-btn" @click.stop="toggleSort()" aria-label="排序">
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="4" y1="6" x2="20" y2="6"/><line x1="4" y1="12" x2="14" y2="12"/><line x1="4" y1="18" x2="8" y2="18"/></svg>
  </button>
  <button class="panel-close" @click="$emit('close')" aria-label="关闭">
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
    <line x1="18" y1="6" x2="6" y2="18" />
    <line x1="6" y1="6" x2="18" y2="18" />
  </svg>
</button>
  </template>
  <button v-else class="panel-done-btn" @click.stop="toggleSort()">完成</button>
</div>
          </div>
          <div class="panel-body">
            <div
              v-for="(item, index) in items"
              :key="item.id"
              class="entry-card"
              @click="$emit('select', item)"
            >
              <div class="entry-icon" :style="{ backgroundColor: item.color }">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                  <!-- 念头: lightbulb (always show) -->
                  <template v-if="item.id === 'thought'">
                    <path d="M10 17h4" />
                    <path d="M10.5 19h3" />
                    <path d="M12 5V3" />
                    <circle cx="12" cy="11" r="5" />
                    <path d="M12 6v2" />
                  </template>
                  <!-- 资产记录: wallet -->
                  <template v-if="item.id === 'asset'">
                    <circle cx="12" cy="12" r="9" />
                    <path d="M15 9.5c-1.5-1.5-4.5-1.5-6 0" />
                    <path d="M9 14.5c1.5 1.5 4.5 1.5 6 0" />
                    <line x1="12" y1="7" x2="12" y2="9" />
                    <line x1="12" y1="15" x2="12" y2="17" />
                  </template>
                  <!-- 读书: book -->
                  <template v-if="item.id === 'reading'">
                    <path d="M12 7v10" />
                    <path d="M5 7h7v8l-3.5-2L5 15V7Z" />
                    <path d="M19 7h-7v8l3.5-2L19 15V7Z" />
                  </template>
                  <!-- 影视: film -->
                  <template v-if="item.id === 'movie'">
                    <rect x="3" y="5" width="18" height="14" rx="2" />
                    <path d="M3 10h18" />
                    <path d="M10 5v5" />
                    <path d="M14 5v5" />
                    <path d="M10 15v4" />
                    <path d="M14 15v4" />
                    <polygon points="10,12 16,15 10,18" fill="white" stroke="none" />
                  </template>
                  <!-- 运动: runner -->
                  <template v-if="item.id === 'exercise'">
                    <circle cx="12" cy="5" r="2" />
                    <path d="M10 8l-3 5l2 1l2-3l2 3v4" />
                    <path d="M12 10l3-3l3 1v3" />
                  </template>
                  <!-- 自律: target -->
                  <template v-if="item.id === 'discipline'">
                    <circle cx="12" cy="12" r="9" />
                    <circle cx="12" cy="12" r="5.5" />
                    <circle cx="12" cy="12" r="2" />
                    <path d="M12 3v2" />
                    <path d="M12 19v2" />
                    <path d="M3 12h2" />
                    <path d="M19 12h2" />
                  </template>
                  <!-- 尿酸: droplet + chart -->
                  <template v-if="item.id === 'uric'">
                    <path d="M12 3C12 3 6 9 6 13a6 6 0 0 0 12 0c0-4-6-10-6-10Z" />
                    <path d="M9 16h6" />
                    <polyline points="9,16 10,14 12,15 15,11 15,16" />
                  </template>
                  <!-- 禁止糖分: circle-slash -->
                  <template v-if="item.id === 'nosugar'">
                    <rect x="6" y="6" width="12" height="12" rx="2" />
                    <circle cx="12" cy="12" r="6" />
                    <path d="M8 8l8 8" />
                  </template>
                  <!-- 待办事项: checklist -->
                  <template v-if="item.id === 'todo'">
                    <rect x="3" y="5" width="18" height="14" rx="2" />
                    <line x1="8" y1="10" x2="16" y2="10" />
                    <line x1="8" y1="14" x2="14" y2="14" />
                    <polyline points="6,10 7,11 9,9" />
                  </template>
                </svg>
              </div>
              <div class="entry-text">
                <span class="entry-title">{{ item.title }}</span>
                <span class="entry-sub">{{ item.subtitle }}</span>
              </div>
              <div v-if="sortMode" class="sort-handle">
                <button class="sort-arrow" @click.stop="moveUp(index)" :disabled="index===0">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="18 15 12 9 6 15"/></svg>
                </button>
                <button class="sort-arrow" @click.stop="moveDown(index)" :disabled="index===items.length-1">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"/></svg>
                </button>
              </div>

              <div v-if="(item.id === 'discipline' || item.id === 'nosugar') && !sortMode" class="action-group">
                <button class="backfill-btn" @click.stop="$emit('backfill', item)" aria-label="补卡">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                </button>
                <button class="check-btn" @click.stop="$emit('select', item)" aria-label="打卡">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#7BA88A" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="20 6 9 17 4 12" />
                  </svg>
                </button>
              </div>
              <div v-if="item.id !== 'thought' && item.id !== 'discipline' && item.id !== 'nosugar' && !sortMode" style="display:contents">
              <button class="entry-plus" @click.stop="$emit('select', item)" aria-label="添加">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
                  <line x1="12" y1="5" x2="12" y2="19" />
                  <line x1="5" y1="12" x2="19" y2="12" />
                </svg>
              </button>
              </div>
              <div v-if="item.id === 'thought' && !sortMode" class="thought-actions">
                <button class="thought-btn thought-btn-sun" @click.stop="$emit('select', { ...item, kind: 'positive' })" aria-label="太阳">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="4"/><path d="M12 2v2"/><path d="M12 20v2"/><path d="M4.93 4.93l1.41 1.41"/><path d="M17.66 17.66l1.41 1.41"/><path d="M2 12h2"/><path d="M20 12h2"/><path d="M6.34 17.66l-1.41 1.41"/><path d="M19.07 4.93l-1.41 1.41"/></svg>
                </button>
                <button class="thought-btn thought-btn-rain" @click.stop="$emit('select', { ...item, kind: 'negative' })" aria-label="阴雨">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.5 19a2.5 2.5 0 1 0 0-5 3.5 3.5 0 0 0-6.5-1.5 4 4 0 0 0-5.5 4 2.5 2.5 0 0 0 0 5"/><line x1="10" y1="19" x2="10" y2="22"/><line x1="13" y1="19" x2="13" y2="22"/></svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template><script setup>
import { ref, onMounted } from "vue"
import { API_BASE } from '../config.js'
defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['close', 'select', 'backfill'])

const sortMode = ref(false)
function toggleSort(){ sortMode.value = !sortMode.value }
async function saveOrder(){ try{ await fetch(API_BASE+'/api/settings/panel-order',{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify({value:JSON.stringify(items.value.map(i=>i.id))})}) }catch(e){ console.error('Failed to save panel order:', e) } }
async function loadOrder(){ try{ const r=await fetch(API_BASE+'/api/settings/panel-order'); const j=await r.json(); if(j.code===0&&j.data&&j.data.value){ const ids=JSON.parse(j.data.value); return ids.map(id=>items.value.find(i=>i.id===id)).filter(Boolean).concat(items.value.filter(i=>!ids.includes(i.id))) } }catch(e){ console.error('Failed to load panel order:', e) } return null }
function moveUp(idx){ if(idx<=0)return; const a=[...items.value]; [a[idx],a[idx-1]]=[a[idx-1],a[idx]]; items.value=a; saveOrder() }
function moveDown(idx){ if(idx>=items.value.length-1)return; const a=[...items.value]; [a[idx],a[idx+1]]=[a[idx+1],a[idx]]; items.value=a; saveOrder() }


const items = ref([
  { id: 'thought',    title: '随记',     subtitle: '记录生活中的阳光与阴雨',   color: '#9DB5C9' },
  { id: 'asset',      title: '资产记录', subtitle: '记录资产变动',        color: '#D4B87A' },
  { id: 'exercise',   title: '运动',     subtitle: '记录运动与锻炼',       color: '#84B8A4' },
  { id: 'reading',    title: '读书',     subtitle: '记录阅读的书目与笔记', color: '#A099C4' },
  { id: 'movie',      title: '影视',     subtitle: '记录电影·剧集·动漫',     color: '#CB99B0' },
  { id: 'discipline', title: '自律',     subtitle: '培养好习惯',           color: '#8CA4BD' },
  { id: 'uric',       title: '尿酸',     subtitle: '记录尿酸值变化',       color: '#BE999B' },
  { id: 'nosugar',    title: '禁止糖分', subtitle: '控糖饮食记录',         color: '#C88C8C' },
  { id: 'todo',       title: '待办事项', subtitle: '添加一个待办事项',       color: '#6C7A89' },
])
onMounted(async()=>{const o=await loadOrder();if(o)items.value=o})
</script>
<style scoped>
.panel-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  background: rgba(0, 0, 0, 0.25);
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.panel-sheet {
  width: 100%;
  max-width: var(--content-width);
  height: 80%;
  background: var(--color-card);
  border-radius: 20px 20px 0 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px 8px;
  flex-shrink: 0;
}

.panel-title {
  font-size: 17px;
  font-weight: 600;
  color: var(--color-ink);
}

.panel-close {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: var(--color-surface-dim);
  color: var(--color-graphite);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background 0.15s;
}

.panel-close:hover {
  background: var(--color-pencil);
}

.panel-body {
  flex: 1;
  overflow-y: auto;
  padding: 8px 16px 32px;
  -webkit-overflow-scrolling: touch;
}

@media (min-width: 600px) {
  .panel-body { padding: 10px 20px 36px; }
}

.entry-card {
  display: flex;
  align-items: center;
  background: var(--color-surface-dim);
  border-radius: var(--radius-md);
  padding: 14px 12px;
  margin-bottom: 8px;
  border: 1px solid transparent;
  cursor: pointer;
  transition: background 0.12s, border-color 0.12s;
}

.entry-card:active {
  background: var(--color-card);
  border-color: var(--color-pencil);
}

.entry-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.entry-text {
  flex: 1;
  display: flex;
  flex-direction: column;
  margin-left: 12px;
  min-width: 0;
}

.entry-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-ink);
  line-height: 1.3;
}

.entry-sub {
  font-size: 12px;
  color: var(--color-graphite);
  line-height: 1.3;
  margin-top: 2px;
}

.entry-plus {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: var(--color-card);
  color: var(--color-graphite);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  flex-shrink: 0;
  transition: background 0.15s, color 0.15s;
}

.entry-plus:hover {
  background: var(--color-pencil);
  color: var(--color-ink);
}

.panel-enter-active,
.panel-leave-active {
  transition: opacity 0.25s ease;
}

.panel-enter-active .panel-sheet,
.panel-leave-active .panel-sheet {
  transition: transform 0.3s cubic-bezier(0.32, 0.72, 0, 1);
}

.panel-enter-from,
.panel-leave-to {
  opacity: 0;
}

.panel-enter-from .panel-sheet,
.panel-leave-to .panel-sheet {
  transform: translateY(100%);
}
.thought-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.thought-btn {
  height: 36px; width: 36px;
  border-radius: 50%;
  border: none;
  display: flex; align-items: center; justify-content: center;
  cursor: pointer;
  transition: background 0.15s, transform 0.12s;
}
.thought-btn:active { transform: scale(0.92); }

.thought-btn-sun {
  background: #E8F5E9; color: #5A8A5A;
}
.thought-btn-sun:hover {
  background: #C8E6C9;
}

.thought-btn-rain {
  background: #FDEAEA; color: #B06060;
}
.thought-btn-rain:hover {
  background: #F5CECE;
}

.check-btn {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  border: 1.5px solid #C5DCC8;
  background: #EBF5EE;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s;
  flex-shrink: 0;
}
.check-btn:hover { background: #D4E8D6; border-color: #7BA88A; }
.panel-header-actions{display:flex;align-items:center;gap:8px}
.panel-sort-btn{width:32px;height:32px;border-radius:50%;border:none;background:var(--color-surface-dim);color:var(--color-graphite);display:flex;align-items:center;justify-content:center;cursor:pointer;transition:all .15s}
.panel-sort-btn:hover{background:var(--color-pencil)}
.panel-done-btn{height:32px;padding:0 16px;border-radius:16px;border:none;background:var(--color-ink);color:var(--color-card);font-size:14px;font-weight:500;cursor:pointer}
.panel-done-btn:hover{opacity:.85}
.panel-done-btn:active{opacity:.7}
.sort-handle{display:flex;align-items:center;gap:6px;flex-shrink:0}
.sort-arrow{width:30px;height:30px;border-radius:50%;border:none;background:var(--color-surface-dim);color:var(--color-graphite);display:flex;align-items:center;justify-content:center;cursor:pointer;transition:all .12s}
.sort-arrow:hover{background:var(--color-pencil);color:var(--color-ink)}
.sort-arrow:disabled{opacity:.25;cursor:default}
.sort-arrow:disabled:hover{background:var(--color-surface-dim);color:var(--color-graphite)}
.backfill-btn{width:34px;height:34px;border-radius:50%;border:1.5px solid var(--color-pencil);background:var(--color-surface-dim);color:var(--color-graphite);display:flex;align-items:center;justify-content:center;cursor:pointer;transition:all .15s;flex-shrink:0}
.backfill-btn:hover{background:var(--color-pencil);color:var(--color-ink);border-color:var(--color-graphite)}
.action-group{display:flex;align-items:center;gap:6px;flex-shrink:0}
</style>
