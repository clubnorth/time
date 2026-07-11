<template>
  <Transition name="picker">
    <div v-if="visible" class="picker-overlay" @click.self="$emit('close')">
      <div class="picker-card">
        <div class="picker-header">
          <button class="picker-cancel" @click="$emit('close')">取消</button>
          <span class="picker-title">选择时间</span>
          <button class="picker-confirm" @click="$emit('confirm')">确定</button>
        </div>
        <div class="picker-body">
          <div class="picker-col">
            <button class="picker-spin-btn" @click="$emit('adjustYear', 1)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><polyline points="18 15 12 9 6 15"/></svg>
            </button>
            <div class="picker-value picker-value-wide" @click.prevent>{{ pickYear }}</div>
            <button class="picker-spin-btn" @click="$emit('adjustYear', -1)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><polyline points="6 9 12 15 18 9"/></svg>
            </button>
            <span class="picker-unit">年</span>
          </div>
          <div class="picker-col">
            <button class="picker-spin-btn" @click="$emit('adjustMonth', 1)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><polyline points="18 15 12 9 6 15"/></svg>
            </button>
            <div class="picker-value picker-value-wide">{{ pad(pickMonth) }}</div>
            <button class="picker-spin-btn" @click="$emit('adjustMonth', -1)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><polyline points="6 9 12 15 18 9"/></svg>
            </button>
            <span class="picker-unit">月</span>
          </div>
          <div class="picker-col">
            <button class="picker-spin-btn" @click="$emit('adjustDay', 1)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><polyline points="18 15 12 9 6 15"/></svg>
            </button>
            <div class="picker-value picker-value-wide">{{ pad(pickDay) }}</div>
            <button class="picker-spin-btn" @click="$emit('adjustDay', -1)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><polyline points="6 9 12 15 18 9"/></svg>
            </button>
            <span class="picker-unit">日</span>
          </div>
          <div class="picker-col">
            <button class="picker-spin-btn" @click="$emit('adjustHour', 1)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><polyline points="18 15 12 9 6 15"/></svg>
            </button>
            <div class="picker-value">{{ pad(pickHour) }}</div>
            <button class="picker-spin-btn" @click="$emit('adjustHour', -1)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><polyline points="6 9 12 15 18 9"/></svg>
            </button>
            <span class="picker-unit">时</span>
          </div>
          <div class="picker-col">
            <button class="picker-spin-btn" @click="$emit('adjustMinute', 1)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><polyline points="18 15 12 9 6 15"/></svg>
            </button>
            <div class="picker-value">{{ pad(pickMinute) }}</div>
            <button class="picker-spin-btn" @click="$emit('adjustMinute', -1)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><polyline points="6 9 12 15 18 9"/></svg>
            </button>
            <span class="picker-unit">分</span>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
defineProps({
  visible: { type: Boolean, default: false },
  pickYear: { type: Number, default: 2025 },
  pickMonth: { type: Number, default: 6 },
  pickDay: { type: Number, default: 1 },
  pickHour: { type: Number, default: 0 },
  pickMinute: { type: Number, default: 0 },
})

defineEmits(['close', 'confirm', 'adjustYear', 'adjustMonth', 'adjustDay', 'adjustHour', 'adjustMinute'])

const pad = (n) => String(n).padStart(2, '0')
</script>

<style scoped>
.picker-overlay {
  position: fixed; inset: 0; z-index: 1200;
  background: rgba(0,0,0,0.25);
  display: flex; align-items: flex-end; justify-content: center;
}
.picker-card {
  width: 100%; max-width: var(--content-width);
  background: var(--color-card);
  border-radius: 20px 20px 0 0;
  box-shadow: 0 -4px 30px rgba(0,0,0,0.08);
  padding-bottom: env(safe-area-inset-bottom, 16px);
}
.picker-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 18px 24px 12px;
}
.picker-cancel {
  font-size: 15px; color: var(--color-graphite);
  background: none; border: none; cursor: pointer; padding: 4px 0;
}
.picker-title {
  font-size: 16px; font-weight: 600; color: var(--color-ink);
}
.picker-confirm {
  font-size: 15px; color: var(--color-ink); font-weight: 600;
  background: none; border: none; cursor: pointer; padding: 4px 0;
}

.picker-body {
  display: flex; align-items: center; justify-content: center;
  padding: 8px 20px 24px; gap: 4px;
}
.picker-col {
  display: flex; flex-direction: column; align-items: center; gap: 4px;
  flex: 1; max-width: 68px;
}
.picker-spin-btn {
  width: 44px; height: 36px;
  border: none; border-radius: 12px;
  background: var(--color-surface-dim);
  color: var(--color-graphite);
  display: flex; align-items: center; justify-content: center;
  cursor: pointer; transition: all 0.12s;
  -webkit-tap-highlight-color: transparent;
}
.picker-spin-btn:hover { background: var(--color-pencil); color: var(--color-ink); }
.picker-spin-btn:active { transform: scale(0.94); }

.picker-value {
  font-size: 38px; font-weight: 700; color: var(--color-ink);
  line-height: 1.1; text-align: center;
  font-variant-numeric: tabular-nums;
  width: 100%; padding: 4px 0;
  user-select: none;
}
.picker-value-wide { font-size: 32px; }
.picker-unit {
  font-size: 12px; color: var(--color-graphite);
  font-weight: 500; margin-top: -2px;
}

.picker-enter-active, .picker-leave-active { transition: opacity 0.2s ease; }
.picker-enter-active .picker-card, .picker-leave-active .picker-card { transition: transform 0.24s cubic-bezier(0.32,0.72,0,1); }
.picker-enter-from, .picker-leave-to { opacity: 0; }
.picker-enter-from .picker-card, .picker-leave-to .picker-card { transform: translateY(100%); }
</style>
