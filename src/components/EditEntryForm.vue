<template>
  <Teleport to="body">
    <Transition name="form">
      <div v-if="visible" class="form-overlay" @click.self="$emit('cancel')">
        <div class="form-sheet">
          <div class="form-navbar">
            <button class="form-nav-left" @click="$emit('cancel')">取消</button>
            <span class="form-nav-title">修改记录</span>
            <button class="form-nav-right" @click="handleSave">保存</button>
          </div>
          <div class="form-body">
            <div class="form-section">
              <label class="form-label">记录时间</label>
              <div class="form-time-picker" @click="tp.openTimePicker()">
                <span class="form-time-text">{{ tp.displayTime.value }}</span>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="form-time-arrow"><polyline points="9 18 15 12 9 6"/></svg>
              </div>
            </div>

            <!-- Asset -->
            <div class="form-section" v-if="entryType === 'asset'">
              <label class="form-label">资产余额</label>
              <input v-model="f1" type="number" class="form-input-wide" />
            </div>

            <!-- Uric -->
            <div class="form-section" v-if="entryType === 'uric'">
              <label class="form-label">尿酸值</label>
              <input v-model="f1" type="number" class="form-input-wide" />
            </div>

            <!-- Exercise -->
            <template v-if="entryType === 'exercise'">
              <div class="form-section">
                <label class="form-label">运动类型</label>
                <div class="capsules">
                  <button v-for="s in sports" :key="s" class="capsule" :class="{ active: f2 === s }" @click="f2 = s">{{ s }}</button>
                </div>
              </div>
              <div class="form-section">
                <label class="form-label">{{ sportLabel }}</label>
                <input v-model="f1" type="number" class="form-input-wide" />
              </div>
            </template>

            <!-- Discipline / Nosugar -->
            <div class="form-section" v-if="entryType === 'discipline' || entryType === 'nosugar'">
              <label class="form-label">连续天数</label>
              <input v-model="f1" type="number" class="form-input-wide" />
            </div>

            <!-- Thought / 随记 -->
            <template v-if="entryType === 'thought'">
              <div class="form-section">
                <label class="form-label">类型</label>
                <div class="capsules">
                  <button class="capsule" :class="{ active: f2 === 'positive' }" @click="f2 = 'positive'">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="4"/><path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M6.34 17.66l-1.41 1.41M19.07 4.93l-1.41 1.41"/></svg>
                  </button>
                  <button class="capsule" :class="{ active: f2 === 'negative' }" @click="f2 = 'negative'">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17.5 19a2.5 2.5 0 1 0 0-5 3.5 3.5 0 0 0-6.5-1.5 4 4 0 0 0-5.5 4 2.5 2.5 0 0 0 0 5"/><line x1="10" y1="19" x2="10" y2="22"/><line x1="13" y1="19" x2="13" y2="22"/></svg>
                  </button>
                </div>
              </div>
              <div class="form-section">
                <label class="form-label">备注</label>
                <textarea v-model="f1" class="form-textarea" rows="4"></textarea>
              </div>
            </template>

            <!-- Movie / Reading: 独立字段 -->
            <template v-if="entryType === 'movie'">
              <div class="form-section" v-for="f in movieFields.filter(f => f.key !== 'tags')" :key="f.key">
                <label class="form-label">{{ f.label }}</label>
                <input v-model="f.value" class="form-input-wide" :placeholder="f.placeholder" />
              </div>
              <!-- 标签：胶囊展示 + 添加 -->
              <div class="form-section">
                <label class="form-label">标签</label>
                <div class="tag-capsules">
                  <span class="tag-capsule" v-for="(t, i) in editMovieTags" :key="i">
                    {{ t }}
                    <button class="tag-close" @click="removeMovieTag(i)">×</button>
                  </span>
                  <button class="tag-add" @click="showTagInput = true">+</button>
                </div>
              </div>
              <div class="form-section">
                <label class="form-label">{{ ratingLabel }}</label>
                <div class="rating-row">
                  <button v-for="n in 10" :key="n" class="star-btn" @click="editRating = n">
                    <svg width="24" height="24" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" :fill="editRating >= n ? '#D4A574' : 'var(--color-pencil)'" stroke="none"/></svg>
                  </button>
                  <span class="rating-text">{{ editRating }}分</span>
                </div>
              </div>
            </template>

            <template v-if="entryType === 'reading'">
              <div class="form-section" v-for="f in readingFields.filter(f => f.key !== 'tags')" :key="f.key">
                <label class="form-label">{{ f.label }}</label>
                <input v-model="f.value" class="form-input-wide" :placeholder="f.placeholder" />
              </div>
              <!-- 标签：胶囊展示 + 添加 -->
              <div class="form-section">
                <label class="form-label">标签</label>
                <div class="tag-capsules">
                  <span class="tag-capsule" v-for="(t, i) in editReadingTags" :key="i">
                    {{ t }}
                    <button class="tag-close" @click="removeReadingTag(i)">×</button>
                  </span>
                  <button class="tag-add" @click="showTagInput = true">+</button>
                </div>
              </div>
              <div class="form-section">
                <label class="form-label">{{ readingCatLabel }}</label>
                <select v-model="editReadingCat" class="form-select">
                  <option v-for="c in readingCats" :key="c" :value="c">{{ c }}</option>
                </select>
              </div>
              <div class="form-section">
                <label class="form-label">{{ ratingLabel }}</label>
                <div class="rating-row">
                  <button v-for="n in 10" :key="n" class="star-btn" @click="editRating = n">
                    <svg width="24" height="24" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" :fill="editRating >= n ? '#D4A574' : 'var(--color-pencil)'" stroke="none"/></svg>
                  </button>
                  <span class="rating-text">{{ editRating }}分</span>
                </div>
              </div>
            </template>

            <!-- Todo: basic edit -->
            <div class="form-section" v-if="entryType === 'todo'">
              <label class="form-label">内容描述</label>
              <textarea v-model="f1" class="form-textarea" rows="6"></textarea>
            </div>
          </div>

          <!-- Tag input popup -->
          <Teleport to="body">
            <div v-if="showTagInput" class="tag-popup-overlay" @click.self="showTagInput = false">
              <div class="tag-popup-card">
                <div class="tag-popup-title">新增标签</div>
                <input v-model="newTagName" class="tag-popup-input" placeholder="最多4个汉字" maxlength="4" @keyup.enter="addTag" />
                <div class="tag-popup-btns">
                  <button class="tag-popup-cancel" @click="showTagInput = false">取消</button>
                  <button class="tag-popup-confirm" @click="addTag" :disabled="!newTagName.trim()">添加</button>
                </div>
              </div>
            </div>
          </Teleport>
        </div>
        <TimePickerModal
          :visible="tp.showTimeModal.value"
          :pickYear="tp.pickYear.value" :pickMonth="tp.pickMonth.value" :pickDay="tp.pickDay.value"
          :pickHour="tp.pickHour.value" :pickMinute="tp.pickMinute.value"
          @close="tp.closeTimePicker()" @confirm="tp.confirmTime()"
          @adjustYear="(d) => tp.pickYear.value += d"
          @adjustMonth="tp.adjustMonth" @adjustDay="tp.adjustDay"
          @adjustHour="tp.adjustHour" @adjustMinute="tp.adjustMinute"
        />
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import TimePickerModal from './TimePickerModal.vue'
import { useTimePicker } from '../composables/useTimePicker.js'

const props = defineProps({ visible: { type: Boolean, default: false }, data: Object })
const emit = defineEmits(['cancel', 'save'])

const tp = useTimePicker()
const f1 = ref(''); const f2 = ref(''); const entryType = ref('')

const sports = ['跑步','骑行','俯卧撑','跳绳','游泳','徒步','自由活动']
const sportLabel = computed(() => {
  const m = {'跑步':'跑步距离(km)','骑行':'骑行距离(km)','俯卧撑':'俯卧撑(个)','跳绳':'跳绳(个)','游泳':'游泳距离(m)','徒步':'徒步距离(km)','自由活动':'活动时长(min)'}
  return m[f2.value] || ''
})

// Movie fields
const movieFields = ref([
  { key:'name', label:'名称', value:'', placeholder:'电影/剧集名称' },
  { key:'director', label:'导演', value:'', placeholder:'导演（/ 分隔）' },
  { key:'cast', label:'主演', value:'', placeholder:'主演（/ 分隔）' },
  { key:'country', label:'国家', value:'', placeholder:'国家' },
  { key:'year', label:'首次上映', value:'', placeholder:'上映年份' },
  { key:'duration', label:'时长', value:'', placeholder:'时长(集数)' },
  { key:'tags', label:'标签', value:'', placeholder:'标签（/ 分隔）' },
])
const editRating = ref(7)
const ratingLabel = computed(() => '评分')
const editReadingCat = ref('')
const readingCats = ['推理悬疑','明日方舟','人物历史','社会科学','小说文学','自我成长','杂项拾遗']
const readingCatLabel = computed(() => '分类')
// Tags
const editMovieTags = ref([])
const editReadingTags = ref([])
const showTagInput = ref(false)
const newTagName = ref('')

const readingFields = ref([
  { key:'name', label:'书名', value:'', placeholder:'书名' },
  { key:'author', label:'作者', value:'', placeholder:'作者' },
  { key:'duration', label:'用时', value:'', placeholder:'阅读用时' },
  { key:'year', label:'首次出版', value:'', placeholder:'出版年份' },
  { key:'tags', label:'标签', value:'', placeholder:'标签（/ 分隔）' },
])

function removeMovieTag(i) { editMovieTags.value.splice(i, 1) }
function removeReadingTag(i) { editReadingTags.value.splice(i, 1) }
function addTag() {
  const v = newTagName.value.trim()
  if (!v || v.length > 4) return
  if (entryType.value === 'movie') {
    if (!editMovieTags.value.includes(v)) editMovieTags.value.push(v)
  } else {
    if (!editReadingTags.value.includes(v)) editReadingTags.value.push(v)
  }
  newTagName.value = ''
  showTagInput.value = false
}

function stripHtml(s) { return (s || '').replace(/<[^>]+>/g, '').trim() }
function parseField(text, key) {
  const re = new RegExp(key + '[：:]\\s*(.+?)(?:<br>|$)', 'i')
  const m = text.match(re)
  return m ? stripHtml(m[1]) : ''
}
// 从HTML描述中提取<span class="r-tag">内文本，返回数组
function parseTags(text) {
  const re = /<span class="r-tag">([^<]+)<\/span>/g
  const tags = []
  let m
  while ((m = re.exec(text)) !== null) {
    tags.push(m[1])
  }
  // 兼容无HTML标签的纯文本标签（/ 分隔）
  if (tags.length === 0) {
    const raw = parseField(text, '标签')
    if (raw) return raw.split(/[\/,，、]/).map(t => t.trim()).filter(Boolean)
  }
  return tags
}

watch(() => props.visible, (val) => {
  if (!val || !props.data) return
  const d = props.data
  entryType.value = d.type || ''

  if (d.type === 'thought') {
    f1.value = d.description || ''
    f2.value = d.valence === 'negative' ? 'negative' : 'positive'
  } else if (d.type === 'exercise') {
    f1.value = d.description || ''
    f2.value = d.title ? d.title.replace('运动·','') : '跑步'
  } else if (d.type === 'movie') {
    const desc = d.description || ''
    movieFields.value.forEach(field => {
      if (field.key === 'tags') {
        editMovieTags.value = parseTags(desc)
      } else if (field.key === 'country') {
        field.value = parseField(desc, '国家').replace(/[【】]/g, '')
      } else {
        field.value = parseField(desc, field.label)
        if (field.key === 'year') field.value = field.value.replace(/年$/, '')
      }
    })
    editRating.value = parseInt(parseField(desc, '评分')) || 7
  } else if (d.type === 'reading') {
    const desc = d.description || ''
    readingFields.value.forEach(field => {
      if (field.key === 'tags') {
        editReadingTags.value = parseTags(desc)
      } else {
        field.value = parseField(desc, field.label)
        if (field.key === 'year') field.value = field.value.replace(/年$/, '')
      }
    })
    editReadingCat.value = parseField(desc, '分类')
    editRating.value = parseInt(parseField(desc, '评分')) || 7
  } else {
    f1.value = d.description || ''
  }

  if (d.recordedAt) {
    const parts = d.recordedAt.split(/[- :]/)
    tp.currentTime.value = new Date(+parts[0], +parts[1]-1, +parts[2], +parts[3]||0, +parts[4]||0, 0)
  } else {
    tp.currentTime.value = new Date()
  }
})

function starHTML(rating) {
  if (!rating || rating <= 0) return ''
  let html = ''
  for (let i = 1; i <= 5; i++) {
    const s = i * 2
    if (rating >= s) html += '<svg class="r-star" width="14" height="14" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" fill="#D4A574"/></svg>'
    else if (rating >= s - 1) html += '<svg class="r-star" width="14" height="14" viewBox="0 0 24 24"><path d="M12 5.5l2.1 4.2 4.6.7-3.3 3.2.8 4.6-4.2-2.2-4.2 2.2.8-4.6L5.3 10.4l4.6-.7L12 5.5zM12 2l-3.09 6.26L2 9.27l5 4.87-1.18 6.88L12 17.77l6.18 3.25L17 14.14l5-4.87-6.91-1.01L12 2z" fill="#D4A574"/></svg>'
    else html += '<svg class="r-star" width="14" height="14" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" fill="var(--color-pencil)" opacity="0.3"/></svg>'
  }
  html += ' ' + rating + '分'
  return html
}

function handleSave() {
  const dt = tp.currentTime.value
  const recordedAt = dt.getFullYear()+'-'+String(dt.getMonth()+1).padStart(2,'0')+'-'+String(dt.getDate()).padStart(2,'0')+' '+String(dt.getHours()).padStart(2,'0')+':'+String(dt.getMinutes()).padStart(2,'0')+':00'

  const out = {
    id: props.data.id,
    category: props.data.category || 'green',
    recorded_at: recordedAt,
    description: f1.value,
  }

  if (entryType.value === 'thought') {
    out.title = f2.value === 'negative' ? '随记·阴雨' : '随记·太阳'
    out.valence = f2.value
  } else if (entryType.value === 'exercise') {
    out.title = '运动·' + f2.value
    out.description = f1.value
    // Recalculate calories (same formula as App.vue)
    const amount = parseFloat(f1.value) || 0
    const met = { '跑步':11,'骑行':9,'俯卧撑':4.5,'跳绳':11.5,'游泳':8,'徒步':6.5,'自由活动':8 }[f2.value] || 8
    const kcal = Math.round(met / 60 * 69 * (f2.value === '跑步' ? amount * 6 : f2.value === '骑行' ? amount * 60/22 : f2.value === '俯卧撑' ? amount / 30 : f2.value === '跳绳' ? amount / 120 : f2.value === '游泳' ? amount / 50 : f2.value === '徒步' ? amount * 12 : amount))
    out.valence = String(kcal)
    out.category = 'green'
  } else if (entryType.value === 'movie') {
    out.title = props.data.title
    const m = movieFields.value
    const lines = ['名称：' + m[0].value]
    if (m[1].value) lines.push('导演：' + m[1].value)
    if (m[2].value) lines.push('主演：' + m[2].value)
    if (m[3].value) lines.push('国家：' + m[3].value.split(/[\/,，、]/).map(c => '【' + c.trim() + '】').join(' '))
    if (m[4].value) lines.push('首次上映：' + m[4].value + '年')
    if (m[5].value) lines.push('时长：' + m[5].value)
    if (editMovieTags.value.length) {
      const tags = editMovieTags.value.map(t => '<span class="r-tag">' + t + '</span>').join(' ')
      lines.push('标签：' + tags)
    }
    if (editRating.value) lines.push('评分：' + starHTML(editRating.value))
    out.description = lines.join('<br>')
  } else if (entryType.value === 'reading') {
    out.title = props.data.title
    const r = readingFields.value
    const lines = ['书名：' + r[0].value]
    if (r[1].value) lines.push('作者：' + r[1].value)
    if (r[2].value) lines.push('用时：' + r[2].value)
    if (r[3].value) lines.push('首次出版：' + r[3].value + '年')
    if (editReadingCat.value) {
      const catColor = { '推理悬疑':'#8B6B4A','明日方舟':'#5B7A9E','人物历史':'#9E7A5B','社会科学':'#6B8E7A','小说文学':'#8A6B9E','自我成长':'#7A9E6B','杂项拾遗':'#9E8A6B' }[editReadingCat.value] || '#888'
      lines.push('分类：<span class="r-cat" style="background:'+catColor+'">' + editReadingCat.value + '</span>')
    }
    if (editReadingTags.value.length) {
      const tags = editReadingTags.value.map(t => '<span class="r-tag">' + t + '</span>').join(' ')
      lines.push('标签：' + tags)
    }
    if (editRating.value) lines.push('评分：' + starHTML(editRating.value))
    out.description = lines.join('<br>')
  } else if (entryType.value === 'todo') {
    out.title = props.data.title
    out.description = f1.value
  } else {
    out.title = props.data.title
  }

  emit('save', out)
}
</script>

<style scoped>
.form-overlay { position:fixed;inset:0;z-index:1100;background:rgba(0,0,0,.25);display:flex;align-items:flex-end;justify-content:center }
.form-sheet { width:100%;max-width:var(--content-width);height:80%;background:var(--color-card);border-radius:16px 16px 0 0;display:flex;flex-direction:column;overflow:hidden;box-shadow:0 -2px 20px rgba(0,0,0,.08) }
.form-navbar { display:flex;align-items:center;justify-content:space-between;padding:14px 16px;background:var(--color-card);border-bottom:1px solid var(--color-border-light);flex-shrink:0 }
.form-nav-left { font-size:15px;color:var(--color-ink);background:none;border:none;cursor:pointer }
.form-nav-title { font-size:16px;font-weight:600;color:var(--color-ink) }
.form-nav-right { font-size:15px;color:#7BA88A;background:none;border:none;cursor:pointer;font-weight:500 }
.form-body { flex:1;overflow-y:auto;padding:20px 16px 40px;-webkit-overflow-scrolling:touch }
.form-section { margin-bottom:20px }
.form-label { display:block;font-size:13px;color:var(--color-graphite);margin-bottom:6px;font-weight:500 }
.form-time-picker { display:flex;align-items:center;justify-content:center;position:relative;background:var(--color-surface-dim);border-radius:var(--radius-md);padding:14px 16px;cursor:pointer;border:1px solid var(--color-border-light);min-height:48px }
.form-time-text { font-size:15px;color:var(--color-ink) }
.form-input-wide { width:100%;height:44px;background:var(--color-surface-dim);border:1px solid var(--color-border-light);border-radius:var(--radius-md);padding:0 14px;font-size:14px;color:var(--color-ink);outline:none;font-family:inherit }
.form-select { width:100%;height:44px;background:var(--color-surface-dim);border:1px solid var(--color-border-light);border-radius:var(--radius-md);padding:0 14px;font-size:14px;color:var(--color-ink);outline:none;font-family:inherit;appearance:none;-webkit-appearance:none;background-image:url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%23666' stroke-width='2.5' stroke-linecap='round'%3E%3Cpolyline points='6 9 12 15 18 9'/%3E%3C/svg%3E");background-repeat:no-repeat;background-position:right 14px center;cursor:pointer }
.form-select:focus { border-color:var(--color-ink) }
.form-input-wide::placeholder { color:var(--color-pencil) }
.form-textarea { width:100%;background:var(--color-surface-dim);border:1px solid var(--color-border-light);border-radius:var(--radius-md);padding:14px 16px;font-size:15px;color:var(--color-ink);line-height:1.6;resize:none;outline:none;font-family:-apple-system,BlinkMacSystemFont,sans-serif;min-height:100px }
.capsules { display:flex;gap:8px;flex-wrap:wrap }
.capsule { flex:1;padding:10px 0;border-radius:12px;border:1.5px solid var(--color-pencil);background:var(--color-surface-dim);color:var(--color-graphite);font-size:14px;font-weight:500;cursor:pointer;transition:all .15s;text-align:center;display:flex;align-items:center;justify-content:center;gap:4px }
.capsule:hover { border-color:var(--color-ink);color:var(--color-ink) }
.capsule.active { background:var(--color-ink);color:var(--color-card);border-color:var(--color-ink) }
.rating-row { display:flex;align-items:center;gap:1px }
.star-btn { width:30px;height:30px;display:flex;align-items:center;justify-content:center;border:none;background:none;cursor:pointer;padding:0;transition:transform .1s }
.star-btn:hover { transform:scale(1.15) }
.star-btn:active { transform:scale(.9) }
.rating-text { font-size:15px;font-weight:600;color:var(--color-ink);margin-left:8px }

.tag-capsules { display:flex;flex-wrap:wrap;gap:8px;align-items:center }
.tag-capsule { display:inline-flex;align-items:center;gap:2px;padding:4px 4px 4px 12px;border-radius:12px;background:var(--color-surface-dim);color:var(--color-graphite);font-size:13px;font-weight:500 }
.tag-close { width:18px;height:18px;border-radius:50%;border:none;background:var(--color-pencil);color:#666;font-size:12px;display:flex;align-items:center;justify-content:center;cursor:pointer;transition:all .12s }
.tag-close:hover { background:#D4787A;color:#fff }
.tag-add { width:30px;height:30px;border-radius:50%;border:1.5px dashed var(--color-pencil);background:transparent;color:var(--color-graphite);font-size:18px;display:flex;align-items:center;justify-content:center;cursor:pointer;transition:all .12s }
.tag-add:hover { border-color:var(--color-ink);color:var(--color-ink) }

.tag-popup-overlay { position:fixed;inset:0;z-index:1300;background:rgba(0,0,0,.3);display:flex;align-items:center;justify-content:center }
.tag-popup-card { width:280px;background:#fff;border-radius:16px;padding:24px 20px 20px;box-shadow:0 8px 30px rgba(0,0,0,.12) }
.tag-popup-title { font-size:16px;font-weight:600;color:var(--color-ink);text-align:center;margin-bottom:16px }
.tag-popup-input { width:100%;height:44px;border:1px solid var(--color-pencil);border-radius:10px;padding:0 14px;font-size:15px;color:var(--color-ink);outline:none;font-family:inherit;text-align:center }
.tag-popup-input:focus { border-color:var(--color-ink) }
.tag-popup-btns { display:flex;gap:12px;margin-top:16px }
.tag-popup-btns button { flex:1;height:40px;border-radius:10px;border:none;font-size:14px;font-weight:500;cursor:pointer }
.tag-popup-cancel { background:var(--color-surface-dim);color:var(--color-graphite) }
.tag-popup-confirm { background:var(--color-ink);color:var(--color-card);transition:opacity .15s }
.tag-popup-confirm:hover { opacity:.85 }
.tag-popup-confirm:disabled { opacity:.4;cursor:default }
</style>
