<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="visible" class="confirm-overlay" @click.self="$emit('close')">
        <div class="confirm-card">
          <div class="confirm-icon">
            <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="#4caf50" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10"/>
              <path d="M8 12l3 3 5-5"/>
            </svg>
          </div>
          <p class="confirm-msg">{{ message }}</p>
          <button class="confirm-btn" @click="$emit('close')">确定</button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
defineProps({
  visible: { type: Boolean, default: false },
  message: { type: String, default: '' },
})
defineEmits(['close'])
</script>

<style scoped>
.confirm-overlay {
  position: fixed; inset: 0; z-index: 1400;
  background: rgba(0,0,0,0.25);
  display: flex; align-items: center; justify-content: center;
}
.confirm-card {
  background: var(--color-card);
  border-radius: var(--radius-lg);
  padding: 32px 28px 24px;
  width: 260px;
  box-shadow: 0 12px 40px rgba(0,0,0,0.1);
  display: flex; flex-direction: column; align-items: center; gap: 16px;
}
.confirm-icon {
  width: 56px; height: 56px;
  background: #EBF5EE;
  border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
}
.confirm-msg {
  font-size: 15px; color: var(--color-ink); text-align: center; line-height: 1.6; margin: 0;
}
.confirm-btn {
  margin-top: 4px;
  padding: 8px 36px;
  border-radius: var(--radius-md); border: none;
  background: var(--color-ink); color: var(--color-card);
  font-size: 14px; font-weight: 500; cursor: pointer;
  transition: opacity 0.15s;
}
.confirm-btn:hover { opacity: .85; }
.confirm-btn:active { opacity: .7; }

.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-active .confirm-card, .fade-leave-active .confirm-card { transition: transform 0.22s cubic-bezier(0.32,0.72,0,1); }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.fade-enter-from .confirm-card, .fade-leave-to .confirm-card { transform: scale(0.92); }
</style>