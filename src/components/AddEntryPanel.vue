<template>
  <Teleport to="body">
    <Transition name="panel">
      <div v-if="visible" class="panel-overlay" @click.self="$emit('close')">
        <div class="panel-sheet">
          <div class="panel-header">
            <span class="panel-title">添加记录</span>
            <button class="panel-close" @click="$emit('close')" aria-label="关闭">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
                <line x1="18" y1="6" x2="6" y2="18" />
                <line x1="6" y1="6" x2="18" y2="18" />
              </svg>
            </button>
          </div>
          <div class="panel-body">
            <div
              v-for="item in items"
              :key="item.id"
              class="entry-card"
              @click="$emit('select', item)"
            >
              <div class="entry-icon" :style="{ backgroundColor: item.color }">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                  <!-- 念头: lightbulb -->
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
                  <!-- 日记: book -->
                  <template v-if="item.id === 'diary'">
                    <path d="M12 7v12" />
                    <path d="M4 7h8v12l-4-2.5L4 19V7Z" />
                    <path d="M20 7h-8v12l4-2.5L20 19V7Z" />
                  </template>
                  <!-- 读书: book -->
                  <template v-if="item.id === 'reading'">
                    <path d="M12 7v10" />
                    <path d="M5 7h7v8l-3.5-2L5 15V7Z" />
                    <path d="M19 7h-7v8l3.5-2L19 15V7Z" />
                  </template>
                  <!-- 电影: film -->
                  <template v-if="item.id === 'movie'">
                    <rect x="3" y="5" width="18" height="14" rx="2" />
                    <path d="M3 10h18" />
                    <path d="M10 5v5" />
                    <path d="M14 5v5" />
                    <path d="M10 15v4" />
                    <path d="M14 15v4" />
                    <polygon points="10,12 16,15 10,18" fill="white" stroke="none" />
                  </template>
                  <!-- 剧集: TV -->
                  <template v-if="item.id === 'tv'">
                    <rect x="2" y="3" width="20" height="15" rx="2" />
                    <path d="M8 21l4-3" />
                    <path d="M16 21l-4-3" />
                    <polygon points="9,7 16,11 9,15" fill="white" stroke="none" />
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
                  <!-- 户外时间: sun -->
                  <template v-if="item.id === 'outdoor'">
                    <circle cx="12" cy="12" r="4.5" />
                    <path d="M12 2v3" />
                    <path d="M12 19v3" />
                    <path d="M4.93 4.93l2.12 2.12" />
                    <path d="M16.95 16.95l2.12 2.12" />
                    <path d="M2 12h3" />
                    <path d="M19 12h3" />
                    <path d="M4.93 19.07l2.12-2.12" />
                    <path d="M16.95 7.05l2.12-2.12" />
                  </template>
                </svg>
              </div>
              <div class="entry-text">
                <span class="entry-title">{{ item.title }}</span>
                <span class="entry-sub">{{ item.subtitle }}</span>
              </div>
              <div v-if="item.id !== 'thought' && item.id !== 'discipline' && item.id !== 'nosugar'" style="display:contents">
              <button class="entry-plus" @click.stop="$emit('select', item)" aria-label="添加">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
                  <line x1="12" y1="5" x2="12" y2="19" />
                  <line x1="5" y1="12" x2="19" y2="12" />
                </svg>
              </button>
              </div>
              <div v-if="item.id === 'discipline' || item.id === 'nosugar'" style="display:contents">
              <button class="check-btn" @click.stop="$emit('select', item)" aria-label="打卡">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#4caf50" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="20 6 9 17 4 12" />
                </svg>
              </button>
            </div>
            <div v-if="item.id === 'thought'" class="thought-actions">
                <button class="thought-btn thought-btn-pos" @click.stop="$emit('select', { ...item, kind: 'positive' })" aria-label="正向">正</button>
                <button class="thought-btn thought-btn-neg" @click.stop="$emit('select', { ...item, kind: 'negative' })" aria-label="负向">负</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template><script setup>
defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['close', 'select'])

const items = [
  { id: 'thought',    title: '念头',     subtitle: '记录一闪而过的想法',   color: '#9DB5C9' },
  { id: 'asset',      title: '资产记录', subtitle: '记录资产变动',        color: '#D4B87A' },
  { id: 'exercise',   title: '运动',     subtitle: '记录运动与锻炼',       color: '#84B8A4' },
  { id: 'reading',    title: '读书',     subtitle: '记录阅读的书目与笔记', color: '#A099C4' },
  { id: 'movie',      title: '电影',     subtitle: '记录看过的电影',       color: '#CB99B0' },
  { id: 'tv',         title: '剧集',     subtitle: '追剧进度与评价',       color: '#D4A882' },
  { id: 'diary',      title: '日记',     subtitle: '记录今天的点点滴滴',   color: '#8CAD8C' },
  { id: 'discipline', title: '自律',     subtitle: '培养好习惯',           color: '#8CA4BD' },
  { id: 'uric',       title: '尿酸',     subtitle: '记录尿酸值变化',       color: '#BE999B' },
  { id: 'nosugar',    title: '禁止糖分', subtitle: '控糖饮食记录',         color: '#C88C8C' },
  { id: 'outdoor',    title: '户外时间', subtitle: '记录户外活动时长',     color: '#92B4C7' },
]
</script>
<style scoped>
.panel-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.panel-sheet {
  width: 100%;
  max-width: 480px;
  height: 80%;
  background: #fafafa;
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
  color: #1a1a1a;
}

.panel-close {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: #e8e8e8;
  color: #666;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background 0.15s;
}

.panel-close:hover {
  background: #ddd;
}

.panel-body {
  flex: 1;
  overflow-y: auto;
  padding: 8px 16px 32px;
  -webkit-overflow-scrolling: touch;
}

.entry-card {
  display: flex;
  align-items: center;
  background: #ffffff;
  border-radius: 12px;
  padding: 14px 12px;
  margin-bottom: 8px;
  border: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background 0.12s;
}

.entry-card:active {
  background: #f5f5f5;
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
  color: #1a1a1a;
  line-height: 1.3;
}

.entry-sub {
  font-size: 12px;
  color: #999;
  line-height: 1.3;
  margin-top: 2px;
}

.entry-plus {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: #f0f0f0;
  color: #555;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  flex-shrink: 0;
  transition: background 0.15s, color 0.15s;
}

.entry-plus:hover {
  background: #e0e0e0;
  color: #222;
}

/* Vue transition */
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
  height: 30px;
  padding: 0 14px;
  border-radius: 15px;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 13px;
  font-weight: 600;
  transition: background 0.15s, color 0.15s;
}

.thought-btn-pos {
  background: #e6f0e6;
  color: #4a7c4a;
}

.thought-btn-pos:hover {
  background: #c8e0c8;
}

.thought-btn-neg {
  background: #f0e6e6;
  color: #7c4a4a;
}

.check-btn {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  border: 1.5px solid #c8e6c9;
  background: #e8f5e9;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s;
  flex-shrink: 0;
}
.check-btn:hover { background: #c8e6c9; border-color: #4caf50; }
.thought-btn-neg:hover {
  background: #e0c8c8;
}
</style>
