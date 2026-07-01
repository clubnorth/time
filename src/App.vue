<template>
  <div class="app-shell">
    <YearMonthHeader v-model="currentMonth" :months="availableMonths" />

    <div class="timeline-container">
      <div class="timeline-line"></div>

      <TimelineEntryGroup
        v-for="group in allGroups"
        :key="group.date"
        :group="group"
      />
    </div>

    <div class="bottom-zone">
      <AddButton @add="handleAdd" />
    </div>
    <AddEntryPanel :visible="showAddPanel" @close="showAddPanel = false" @select="handleSelect" />
    <ThoughtFormPanel :visible="showThoughtForm" :kind="thoughtKind" @cancel="showThoughtForm = false" @create="handleThoughtCreate" />
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted, onUnmounted } from 'vue'
import YearMonthHeader from './components/YearMonthHeader.vue'
import TimelineEntryGroup from './components/TimelineEntryGroup.vue'
import AddButton from './components/AddButton.vue'
import ThoughtFormPanel from './components/ThoughtFormPanel.vue'
import AddEntryPanel from './components/AddEntryPanel.vue'

const currentMonth = ref('2025-06')
const skipScrollWatch = ref(false)

const availableMonths = ref([
  { label: '2025��1��', value: '2025-01' },
  { label: '2025��2��', value: '2025-02' },
  { label: '2025��3��', value: '2025-03' },
  { label: '2025��4��', value: '2025-04' },
  { label: '2025��5��', value: '2025-05' },
  { label: '2025��6��', value: '2025-06' },
])

const allTimelineData = reactive({
  '2025-01': [
    { date: '2025-01-03', weekday: '周五', dateNum: '1月3日', entries: [
      { id: 'j1', time: '09:00', title: '新年规划', description: '制定 Q1 工作计划与目标', category: 'red' },
      { id: 'j2', time: '14:00', title: '团队 kickoff', description: '全员大会，同步年度方向', category: 'yellow' },
    ]},
    { date: '2025-01-10', weekday: '周五', dateNum: '1月10日', entries: [
      { id: 'j3', time: '10:00', title: '技术方案评审', description: '新架构方案讨论', category: 'green' },
      { id: 'j4', time: '15:30', title: '代码审查', description: '核心模块 review', category: 'red' },
      { id: 'j5', time: '19:00', title: '部门聚餐', description: '年底尾牙，吃火锅', category: 'yellow' },
    ]},
    { date: '2025-01-18', weekday: '周六', dateNum: '1月18日', entries: [
      { id: 'j6', time: '08:00', title: '晨跑', description: '10 公里 LSD，配速 6:00', category: 'green' },
      { id: 'j7', time: '14:00', title: '看书', description: '《人月神话》第二章', category: 'yellow' },
    ]},
  ],
  '2025-02': [
    { date: '2025-02-05', weekday: '周三', dateNum: '2月5日', entries: [
      { id: 'f1', time: '09:30', title: '节前冲刺', description: '发版前回归测试', category: 'red' },
      { id: 'f2', time: '12:00', title: '午饭', description: '跟同事去吃了拉面', category: 'green' },
      { id: 'f3', time: '16:00', title: '上线部署', description: '春节前最后一版上线', category: 'yellow' },
      { id: 'f4', time: '18:30', title: '放假啦', description: '收拾工位，准备过年', category: 'green' },
    ]},
    { date: '2025-02-20', weekday: '周四', dateNum: '2月20日', entries: [
      { id: 'f5', time: '10:00', title: '节后收心会', description: '对齐年后 sprint 计划', category: 'yellow' },
      { id: 'f6', time: '14:30', title: '设计评审', description: '新版 UI 交互稿讨论', category: 'red' },
    ]},
    { date: '2025-02-28', weekday: '周五', dateNum: '2月28日', entries: [
      { id: 'f7', time: '09:00', title: '月末复盘', description: '回顾 2 月进展，调整 3 月计划', category: 'red' },
      { id: 'f8', time: '11:00', title: '线上分享', description: '分享了一次关于性能优化的主题', category: 'green' },
      { id: 'f9', time: '17:00', title: '周五下午茶', description: '行政准备的蛋糕和水果', category: 'yellow' },
    ]},
  ],
  '2025-03': [
    { date: '2025-03-03', weekday: '周一', dateNum: '3月3日', entries: [
      { id: 'm1', time: '09:00', title: '周一站会', description: '同步本周重点事项', category: 'yellow' },
      { id: 'm2', time: '11:30', title: '客户沟通', description: '与甲方确认需求变更', category: 'red' },
      { id: 'm3', time: '15:00', title: '编码时间', description: '集中开发新功能模块', category: 'green' },
    ]},
    { date: '2025-03-12', weekday: '周三', dateNum: '3月12日', entries: [
      { id: 'm4', time: '08:30', title: '健身', description: '力量训练，背+二头', category: 'green' },
      { id: 'm5', time: '10:00', title: '技术分享会', description: '团队内部分享了 Rust 入门', category: 'yellow' },
      { id: 'm6', time: '19:00', title: '看电影', description: '新上映的科幻片', category: 'green' },
    ]},
  ],
  '2025-04': [
    { date: '2025-04-15', weekday: '周二', dateNum: '4月15日', entries: [
      { id: 'a1', time: '09:30', title: '晨间站会', description: '同步开发进度，讨论阻塞项与当日计划', category: 'yellow' },
      { id: 'a2', time: '11:00', title: '产品方案评审', description: '新功能 PRD 评审', category: 'red' },
      { id: 'a3', time: '15:00', title: '代码走读', description: '支付模块重构方案讨论', category: 'green' },
    ]},
    { date: '2025-04-20', weekday: '周日', dateNum: '4月20日', entries: [
      { id: 'a4', time: '10:00', title: '周末读书', description: '读完了《系统设计面试》', category: 'green' },
      { id: 'a5', time: '16:30', title: '健身房', description: '腿日训练', category: 'yellow' },
    ]},
  ],
  '2025-05': [
    { date: '2025-05-03', weekday: '周六', dateNum: '5月3日', entries: [
      { id: 'b1', time: '08:15', title: '高铁出发', description: 'G123 次，北京南 -> 上海虹桥', category: 'yellow' },
      { id: 'b2', time: '12:00', title: '抵达上海', description: '入住酒店，午餐后前往会场', category: 'green' },
      { id: 'b3', time: '14:30', title: '主题演讲', description: '技术大会分享：前端性能优化实践', category: 'red' },
    ]},
    { date: '2025-05-10', weekday: '周六', dateNum: '5月10日', entries: [
      { id: 'b4', time: '09:00', title: '团队复盘', description: 'Q2 目标回顾，识别改进方向', category: 'red' },
      { id: 'b5', time: '14:00', title: '一对一沟通', description: '与团队成员月度 1:1', category: 'yellow' },
    ]},
    { date: '2025-05-24', weekday: '周六', dateNum: '5月24日', entries: [
      { id: 'b6', time: '19:00', title: '朋友聚餐', description: '预约了炭火烤肉，四个人', category: 'green' },
      { id: 'b7', time: '21:30', title: '电影夜', description: '看了一部老电影', category: 'yellow' },
    ]},
  ],
  '2025-06': [
    { date: '2025-06-07', weekday: '周六', dateNum: '6月7日', entries: [
      { id: 'c1', time: '07:45', title: '晨跑', description: '5 公里轻松跑，配速 5:30', category: 'green' },
      { id: 'c2', time: '10:00', title: '产品迭代会', description: '对齐下个 sprint 优先级', category: 'yellow' },
      { id: 'c3', time: '13:00', title: '午饭 + 散步', description: '公司附近新开的轻食店', category: 'green' },
      { id: 'c4', time: '16:00', title: '线上技术分享', description: 'Rust 异步运行时分享', category: 'red' },
    ]},
    { date: '2025-06-15', weekday: '周日', dateNum: '6月15日', entries: [
      { id: 'c5', time: '09:30', title: '整理笔记', description: '把这周学到的知识点整理成文档', category: 'yellow' },
      { id: 'c6', time: '14:00', title: '线上课程', description: '完成了系统设计课程第三章', category: 'green' },
    ]},
    { date: '2025-06-25', weekday: '周三', dateNum: '6月25日', entries: [
      { id: 'c7', time: '11:00', title: '线上评审会议', description: '讨论新版 UI 组件库方案', category: 'red' },
      { id: 'c8', time: '17:30', title: '下班后小聚', description: '和同事喝了杯咖啡', category: 'green' },
    ]},
  ],
})

const showAddPanel = ref(false)
const showThoughtForm = ref(false)
const thoughtKind = ref('positive')

let thoughtIdCounter = 1000
const weekdays = ['����', '��һ', '�ܶ�', '����', '����', '����', '����']

const allGroups = computed(() => {
  const months = Object.keys(allTimelineData).sort().reverse()
  const groups = []
  for (const month of months) {
    for (const g of [...allTimelineData[month]].reverse()) {
      groups.push({
        ...g,
        entries: [...g.entries].sort((a, b) => b.time.localeCompare(a.time)),
      })
    }
  }
  return groups
})

function handleAdd() {
  showAddPanel.value = true
}

function handleSelect(item) {
  console.log('Selected:', item)
  if (item.kind) {
    thoughtKind.value = item.kind
    showThoughtForm.value = true
  } else {
    showAddPanel.value = false
  }
}

function handleThoughtCreate(data) {
  const d = data.time
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const monthKey = `${y}-${m}`
  const dateKey = `${monthKey}-${day}`
  const timeStr = `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
  const weekday = weekdays[d.getDay()]
  const dateNum = `${d.getMonth() + 1}��${d.getDate()}��`

  const newEntry = {
    id: `thought_${thoughtIdCounter++}`,
    time: timeStr,
    title: data.kind === 'positive' ? '��ͷ����' : '��ͷ����',
    description: data.note,
    category: data.kind === 'positive' ? 'green' : 'red',
  }

  // Ensure month key exists
  if (!allTimelineData[monthKey]) {
    allTimelineData[monthKey] = []
    availableMonths.value.push({
      label: `${y}��${d.getMonth() + 1}��`,
      value: monthKey,
    })
    availableMonths.value.sort((a, b) => b.value.localeCompare(a.value))
  }

  // Find existing date group or create new one
  const existingGroup = allTimelineData[monthKey].find(g => g.date === dateKey)
  if (existingGroup) {
    existingGroup.entries.push(newEntry)
  } else {
    allTimelineData[monthKey].push({
      date: dateKey,
      weekday,
      dateNum,
      entries: [newEntry],
    })
  }

  showThoughtForm.value = false
  showAddPanel.value = false
}

function handleScroll() {
  const dots = document.querySelectorAll('[data-date]')
  for (const dot of dots) {
    const rect = dot.getBoundingClientRect()
    if (rect.top > 55) {
      const date = dot.dataset.date
      const month = date.substring(0, 7)
      if (month !== currentMonth.value) {
        skipScrollWatch.value = true
        currentMonth.value = month
      }
      break
    }
  }
}

watch(currentMonth, (newMonth) => {
  if (skipScrollWatch.value) {
    skipScrollWatch.value = false
    return
  }
  const dots = document.querySelectorAll('[data-date]')
  for (const dot of dots) {
    if (dot.dataset.date.startsWith(newMonth)) {
      dot.scrollIntoView({ behavior: 'smooth', block: 'start' })
      break
    }
  }
})

onMounted(() => {
  window.addEventListener('scroll', handleScroll, { passive: true })
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
.app-shell {
  position: relative;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
}

.timeline-container {
  position: relative;
  display: grid;
  grid-template-columns: 46px 4px 20px 9px 1fr;
  padding: 64px 24px 48px 24px;
  flex: 1;
}

.timeline-line {
  position: absolute;
  left: calc(24px + 46px + 4px + 10px);
  top: 0;
  bottom: 0;
  width: 1px;
  background: #d0d0d0;
  z-index: 0;
  pointer-events: none;
}

.bottom-zone {
  position: fixed;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 100%;
  max-width: 480px;
  height: 60px;
  background: #f8f8f8;
  border-top: 1px solid #ebebeb;
  display: flex;
  flex-direction: column;
  align-items: center;
  z-index: 50;
}

</style>
