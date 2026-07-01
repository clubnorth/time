<template>
  <div class="date-group">
    <div class="date-row">
      <div class="date-label-col">
        <span class="date-weekday">{{ group.weekday }}</span>
      </div>
      <div class="dot-col">
        <div class="solid-dot" :class="dotClass" :data-date="group.date"></div>
      </div>
      <div class="date-num-col">
        <span class="date-num">{{ group.dateNum }}</span>
      </div>
    </div>
    <TimelineEntryCard
      v-for="entry in group.entries"
      :key="entry.id"
      :entry="entry"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import TimelineEntryCard from './TimelineEntryCard.vue'

const props = defineProps({
  group: Object
})

const dotClass = computed(() => {
  const w = props.group.weekday
  if (w === '周五') return 'dot-yellow'
  if (w === '周六' || w === '周日') return 'dot-green'
  return 'dot-default'
})
</script>

<style scoped>
.date-group {
  display: contents;
}

.date-row {
  display: contents;
}

.date-label-col {
  grid-column: 1;
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  padding-right: 0;
  padding-top: 16px;
}

.date-weekday {
  font-size: 18px;
  font-weight: 700;
  color: #2c2c2c;
  line-height: 20px;
  white-space: nowrap;
}

.dot-col {
  grid-column: 3;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 20px;
}

.solid-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
  z-index: 2;
}

.solid-dot.dot-default { background: #3a3a3a; }
.solid-dot.dot-yellow  { background: #f5a623; }
.solid-dot.dot-green   { background: #4caf50; }

.date-num-col {
  grid-column: 5;
  display: flex;
  align-items: flex-start;
  padding-top: 16px;
  padding-bottom: 20px;
}

.date-num {
  font-size: 18px;
  font-weight: 700;
  color: #2c2c2c;
  line-height: 20px;
  white-space: nowrap;
}
</style>