<template>
  <button class="add-btn" @click="$emit('add')" aria-label="新增条目">
    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
      <line x1="12" y1="5" x2="12" y2="19" />
      <line x1="5" y1="12" x2="19" y2="12" />
    </svg>
  </button>
</template>

<script setup>
defineEmits(['add'])
</script>

<style scoped>
.add-btn {
  position: absolute;
  top: -26px;
  left: 50%;
  transform: translateX(-50%);
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: var(--color-ink);
  color: var(--color-card);
  border: 3px solid #fcfcfc;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  transition: transform 0.15s, box-shadow 0.15s;
  z-index: 1;
}
.add-btn:hover {
  transform: translateX(-50%) scale(1.05);
  box-shadow: 0 6px 22px rgba(0, 0, 0, 0.18);
}
.add-btn:active {
  transform: translateX(-50%) scale(0.96);
}
</style>