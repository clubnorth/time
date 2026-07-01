<template>
  <div class="year-month-header">
    <button class="header-btn" @click.stop="toggle">
      <span class="ym-text">{{ displayText }}</span>
      <svg class="ym-chevron" :class="{ open: showDropdown }" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
        <polyline points="6 9 12 15 18 9" />
      </svg>
    </button>
    <transition name="fade">
      <div v-if="showDropdown" class="dropdown-overlay" @click="showDropdown = false">
        <div class="dropdown-panel" @click.stop>
          <button
            v-for="item in months"
            :key="item.value"
            class="dropdown-item"
            :class="{ active: item.value === modelValue }"
            @click="select(item.value)"
          >
            {{ item.label }}
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  modelValue: String,
  months: Array
})

const emit = defineEmits(['update:modelValue'])

const showDropdown = ref(false)

const displayText = computed(() => {
  const found = props.months.find(m => m.value === props.modelValue)
  return found ? found.label : props.modelValue
})

function select(value) {
  emit('update:modelValue', value)
  showDropdown.value = false
}

function toggle() {
  showDropdown.value = !showDropdown.value
}
</script>

<style scoped>
.year-month-header {
  position: fixed;
  top: 0;
  z-index: 50;
  left: 50%;
  transform: translateX(-50%);
  width: 100%;
  max-width: 480px;
  background: #ffffff;
  padding: 16px 24px 12px;
  border-bottom: 1px solid #f0f0f0;
}

.header-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px 0;
}

.ym-text {
  font-size: 22px;
  font-weight: 700;
  color: #2c2c2c;
  letter-spacing: 0.5px;
}

.ym-chevron {
  color: #999;
  margin-top: 2px;
  transition: transform 0.2s ease;
}

.ym-chevron.open {
  transform: rotate(180deg);
}

.dropdown-overlay {
  position: fixed;
  inset: 0;
  z-index: 40;
  background: rgba(0, 0, 0, 0.15);
}

.dropdown-panel {
  position: absolute;
  top: 56px;
  left: 24px;
  right: 24px;
  max-width: 432px;
  margin: 0 auto;
  background: #ffffff;
  border-radius: 10px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  overflow: hidden;
  padding: 4px 0;
}

.dropdown-item {
  display: block;
  width: 100%;
  padding: 13px 20px;
  background: none;
  border: none;
  text-align: left;
  font-size: 15px;
  color: #2c2c2c;
  cursor: pointer;
  transition: background 0.15s;
}

.dropdown-item:hover {
  background: #f5f5f5;
}

.dropdown-item.active {
  font-weight: 600;
  color: #1a1a1a;
  background: #f8f8f8;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>