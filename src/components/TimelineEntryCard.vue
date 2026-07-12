<template>
  <div class="entry-card-row">
    <div class="time-col">
      <span class="time-text">{{ entry.time }}</span>
    </div>
    <div class="dot-col">
      <div class="hollow-dot"></div>
    </div>
    <div class="card-col">
      <div class="card-wrapper" :class="{ swiped: swipeOffset < -30 }">
        <div
          class="content-card" :class="'cat-' + (entry.category || 'yellow')"
          :style="{ transform: 'translateX(' + swipeOffset + 'px)' }"
          @pointerdown.prevent="onDown"
          @pointermove.prevent="onMove"
          @pointerup="onUp"
          @pointercancel="onUp"
        >
          <div class="card-header">
            <span class="cat-line" :class="'cat-' + (entry.category || 'yellow')"></span>
            <h3 class="card-title">
              {{ entry.title }}
              <span v-if="entry.isThought && entry.valence === 'positive'" class="thought-icon thought-icon-sun">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="4"/><path d="M12 2v2"/><path d="M12 20v2"/><path d="M4.93 4.93l1.41 1.41"/><path d="M17.66 17.66l1.41 1.41"/><path d="M2 12h2"/><path d="M20 12h2"/><path d="M6.34 17.66l-1.41 1.41"/><path d="M19.07 4.93l-1.41 1.41"/></svg>
              </span>
              <span v-if="entry.isThought && entry.valence === 'negative'" class="thought-icon thought-icon-rain">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.5 19a2.5 2.5 0 1 0 0-5 3.5 3.5 0 0 0-6.5-1.5 4 4 0 0 0-5.5 4 2.5 2.5 0 0 0 0 5"/><line x1="10" y1="19" x2="10" y2="22"/><line x1="13" y1="19" x2="13" y2="22"/></svg>
              </span>
            </h3>
          </div>
          <p v-if="entry.description" class="card-body" :class="{ asset: entry.isAsset, reading: entry.isReading, movie: entry.isMovie }">
          <span v-if="entry.isAsset || entry.isReading || entry.isMovie" v-html="entry.description"></span>
          <template v-else>{{ entry.description }}</template>
        </p>
        </div>
        <div class="card-actions">
          <button v-if="!(entry.isDiscipline || entry.isNosugar)" class="action-btn action-edit" @click.stop="confirmEdit">修改</button>
          <button class="action-btn action-delete" @click.stop="confirmDelete">删除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({ entry: Object })
const emit = defineEmits(['delete-entry', 'edit-entry'])

const swipeOffset = ref(0)
let swipeStartX = 0

function onDown(e) {
  swipeOffset.value = 0
  swipeStartX = e.clientX
}
function onMove(e) {
  const dx = e.clientX - swipeStartX
  if (dx < 0) swipeOffset.value = Math.max(dx, -150)
}
function onUp() {
  if (swipeOffset.value < -50) {
    swipeOffset.value = -130
  } else {
    swipeOffset.value = 0
  }
}
function confirmDelete() {
  swipeOffset.value = 0
  emit('delete-entry', { id: props.entry.id, isCheckIn: props.entry.isDiscipline || props.entry.isNosugar, checkType: props.entry.isDiscipline ? 'discipline' : props.entry.isNosugar ? 'nosugar' : '' })
}
function confirmEdit() {
  swipeOffset.value = 0
  emit('edit-entry', { id: props.entry.id, title: props.entry.title, description: props.entry.description, category: props.entry.category, time: props.entry.time, valence: props.entry.valence })
}
</script>

<style scoped>
.entry-card-row { display: contents }

.time-col {
  grid-column: 1; display: flex; align-items: flex-start;
  justify-content: flex-end; padding-right: 0; padding-top: 16px;
}
.time-text {
  font-size: 12px; color: var(--color-graphite); font-weight: 400;
  line-height: 20px; white-space: nowrap; font-variant-numeric: tabular-nums;
}
@media (min-width: 600px) { .time-text { font-size: 13px } }

.dot-col {
  grid-column: 3; display: flex; align-items: flex-start;
  justify-content: center; padding-top: 20px;
}
.hollow-dot {
  width: 8px; height: 8px; border-radius: 50%;
  border: 1.5px solid var(--color-pencil); background: var(--color-card);
  flex-shrink: 0; z-index: 2;
}

.card-col { grid-column: 5; min-width: 0; overflow: visible }
.card-wrapper { position: relative; overflow: hidden; border-radius: 0 var(--radius-md) var(--radius-md) 0; margin-bottom: 12px }
.card-wrapper.swiped { z-index: 5 }

.content-card {
  background: var(--color-card); border: 1px solid var(--color-border-light);
  border-left-width: 2px; border-radius: 0 var(--radius-md) var(--radius-md) 0;
  padding: 14px 16px;
  transition: transform .2s ease;
  position: relative; z-index: 2;
}
.content-card.cat-yellow { border-left-color: #D4A574 }
.content-card.cat-red    { border-left-color: #D4787A }
.content-card.cat-green  { border-left-color: #7BA88A }

.card-header { display:flex;align-items:center;gap:8px;margin-bottom:4px }
.cat-line { width:6px;height:6px;border-radius:50%;flex-shrink:0 }
.cat-line.cat-yellow { background:#D4A574 }
.cat-line.cat-red    { background:#D4787A }
.cat-line.cat-green  { background:#7BA88A }

.card-title {
  font-size:14px;font-weight:600;color:var(--color-ink);line-height:1.4;
  margin:0;word-break:break-word;display:flex;align-items:center;gap:4px;
}
@media (min-width: 600px) { .card-title { font-size:15px } }
.thought-icon { display:inline-flex;align-items:center;justify-content:center;flex-shrink:0 }
.thought-icon-sun { color:#5A8A5A }
.thought-icon-rain { color:#B06060 }

.card-body {
  font-size:13px;color:var(--color-graphite);line-height:1.65;
  margin:0;word-break:break-word;text-align:left;
}
.card-body.reading,.card-body.movie { line-height:1.8 }
@media (min-width: 600px) { .card-body { font-size:14px } }
.card-body :deep(.r-cat) {
  display:inline-block;padding:2px 12px;border-radius:10px;color:#fff;font-size:12px;font-weight:600;margin-top:2px;
}
.card-body :deep(.r-tag) {
  display:inline-block;padding:2px 10px;border-radius:10px;background:var(--color-surface-dim);color:var(--color-graphite);font-size:12px;font-weight:500;margin-right:4px;margin-top:2px;
}
.card-body :deep(.rainbow) {
  font-weight:700;
  background:linear-gradient(90deg,#D4787A,#D4A574,#9DB5C9,#84B8A4,#8CA4BD,#C88C8C,#D4787A);
  background-size:200% auto;-webkit-background-clip:text;background-clip:text;
  -webkit-text-fill-color:transparent;color:transparent;
  animation:rainbow-flow 3s linear infinite;
}
@keyframes rainbow-flow { to { background-position:200% center } }

/* Swipe actions */
.card-actions {
  position:absolute; top:0; right:0; bottom:0;
  display:flex; align-items:stretch;
  z-index:1;
}
.action-btn {
  width:65px;border:none;font-size:13px;font-weight:600;color:#fff;
  cursor:pointer;display:flex;align-items:center;justify-content:center;
  flex-shrink:0;
}
.action-delete { background:#D4787A }
.action-edit   { background:#7B9FC6 }
</style>
