import { ref, computed } from 'vue'

export function useTimePicker() {
  const currentTime = ref(new Date())
  const showTimeModal = ref(false)
  const pickYear = ref(2025)
  const pickMonth = ref(6)
  const pickDay = ref(1)
  const pickHour = ref(0)
  const pickMinute = ref(0)

  const pad = (n) => String(n).padStart(2, '0')

  const displayTime = computed(() => {
    const d = currentTime.value
    return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日 ${pad(d.getHours())}:${pad(d.getMinutes())}`
  })

  const daysInMonth = computed(() => {
    return new Date(pickYear.value, pickMonth.value, 0).getDate()
  })

  function openTimePicker() {
    const d = currentTime.value
    pickYear.value = d.getFullYear()
    pickMonth.value = d.getMonth() + 1
    pickDay.value = d.getDate()
    pickHour.value = d.getHours()
    pickMinute.value = d.getMinutes()
    showTimeModal.value = true
  }

  function closeTimePicker() {
    showTimeModal.value = false
  }

  function confirmTime() {
    const d = new Date(pickYear.value, pickMonth.value - 1, pickDay.value, pickHour.value, pickMinute.value, 0)
    currentTime.value = d
    showTimeModal.value = false
  }

  function adjustMonth(delta) {
    let m = pickMonth.value + delta
    if (m > 12) { m = 1; pickYear.value++ }
    if (m < 1) { m = 12; pickYear.value-- }
    pickMonth.value = m
    if (pickDay.value > daysInMonth.value) pickDay.value = daysInMonth.value
  }

  function adjustDay(delta) {
    let d = pickDay.value + delta
    const max = daysInMonth.value
    if (d > max) { d = 1; adjustMonth(1) }
    if (d < 1) { adjustMonth(-1); d = new Date(pickYear.value, pickMonth.value - 1, 0).getDate() }
    pickDay.value = d
  }

  function adjustHour(delta) {
    pickHour.value = (pickHour.value + delta + 24) % 24
  }

  function adjustMinute(delta) {
    pickMinute.value = (pickMinute.value + delta + 60) % 60
  }

  return {
    currentTime, showTimeModal, pickYear, pickMonth, pickDay, pickHour, pickMinute,
    pad, displayTime, daysInMonth,
    openTimePicker, closeTimePicker, confirmTime,
    adjustMonth, adjustDay, adjustHour, adjustMinute,
  }
}