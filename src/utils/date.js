export function formatRecordTime(d) {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const monthKey = `${y}-${m}`
  const recordedAt = `${monthKey}-${day} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}:00`
  return { recordedAt, monthKey }
}