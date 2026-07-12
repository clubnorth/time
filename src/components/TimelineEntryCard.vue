<template>
  <div class="entry-card-row">
    <div class="time-col">
      <span class="time-text">{{ entry.time }}</span>
    </div>
    <div class="dot-col">
      <div class="hollow-dot"></div>
    </div>
    <div class="card-col">
      <div
        class="card-wrapper"
        @pointerdown.prevent="onDown"
        @pointerup="onUp"
        @pointercancel="onUp"
      >
        <div class="content-card" :class="'cat-' + (entry.category || 'yellow')">
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

        <Transition name="popover">
          <div v-if="showPopover" class="delete-popover">
            <div class="popover-arrow"></div>
            <div class="popover-body">
              <button v-if="!(entry.isDiscipline || entry.isNosugar)" class="popover-edit-btn" @click.stop="confirmEdit">修改</button>
              <button class="popover-delete-btn" @click.stop="confirmDelete">删除</button>
              <button class="popover-cancel-btn" @click.stop="showPopover = false">取消</button>
            </div>
          </div>
        </Transition>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onUnmounted } from 'vue'

const props = defineProps({ entry: Object })
const emit = defineEmits(['delete-entry', 'edit-entry'])

const showPopover = ref(false)
let longPressTimer = null

function onDown() {
  longPressTimer = setTimeout(() => { showPopover.value = true }, 3000)
}
function onUp() {
  if (longPressTimer) { clearTimeout(longPressTimer); longPressTimer = null }
}
function confirmDelete() {
  showPopover.value = false
  emit('delete-entry', { id: props.entry.id, isCheckIn: props.entry.isDiscipline || props.entry.isNosugar, checkType: props.entry.isDiscipline ? 'discipline' : props.entry.isNosugar ? 'nosugar' : '' })
}
function confirmEdit() {
  showPopover.value = false
  emit('edit-entry', { id: props.entry.id, title: props.entry.title, description: props.entry.description, category: props.entry.category, time: props.entry.time, valence: props.entry.valence })
}

onUnmounted(() => { if (longPressTimer) clearTimeout(longPressTimer) })
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
  flex-shrink: 0; z-index: 2; transition: transform .15s ease, border-color .15s ease;
}

.card-col { grid-column: 5; min-width: 0; overflow: visible }
.card-wrapper { position: relative }

.content-card {
  background: var(--color-card); border: 1px solid var(--color-border-light);
  border-left-width: 2px; border-radius: 0 var(--radius-md) var(--radius-md) 0;
  padding: 14px 16px; margin-bottom: 12px;
  transition: box-shadow .15s ease, transform .12s ease;
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
.card-body.reading { line-height:1.8 }
.card-body.movie { line-height:1.8 }
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

/* Delete popover */
.delete-popover {
  position:absolute;bottom:-8px;left:50%;transform:translateX(-50%) translateY(100%);
  z-index:30;display:flex;flex-direction:column;align-items:center;
}
.popover-arrow {
  width:0;height:0;
  border-left:8px solid transparent;border-right:8px solid transparent;
  border-bottom:8px solid #fff;
  filter:drop-shadow(0 -1px 1px rgba(0,0,0,.06));
}
.popover-body {
  display:flex;align-items:center;gap:0;
  background:#fff;border-radius:10px;
  box-shadow:0 4px 16px rgba(0,0,0,.12);
  overflow:hidden;
}
.popover-delete-btn {
  padding:10px 20px;border:none;background:none;
  font-size:14px;font-weight:500;color:#D4787A;cursor:pointer;
  transition:background .12s;
  border-left:1px solid var(--color-border-light);
}
.popover-delete-btn:hover { background:#FDF0F0 }
.popover-edit-btn {
  padding:10px 20px;border:none;background:none;
  font-size:14px;font-weight:500;color:#7B9FC6;cursor:pointer;
  transition:background .12s;
}
.popover-edit-btn:hover { background:#EFF6FF }
.popover-cancel-btn {
  padding:10px 20px;border:none;background:none;
  font-size:14px;color:var(--color-graphite);cursor:pointer;
  border-left:1px solid var(--color-border-light);
  transition:background .12s;
}
.popover-cancel-btn:hover { background:var(--color-surface-dim) }

.popover-enter-active { transition:opacity .15s ease,transform .15s ease }
.popover-leave-active { transition:opacity .1s ease }
.popover-enter-from,.popover-leave-to { opacity:0;transform:translateX(-50%) translateY(100%) scale(.9) }
</style>
