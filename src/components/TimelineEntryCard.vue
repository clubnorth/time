<template>
  <div class="entry-card-row">
    <div class="time-col">
      <span class="time-text">{{ entry.time }}</span>
    </div>
    <div class="dot-col">
      <div class="hollow-dot"></div>
    </div>
    <div class="card-col">
      <div class="content-card">
        <div class="card-header">
          <span class="cat-line" :class="'cat-' + (entry.category || 'yellow')"></span>
          <h3 class="card-title">{{ entry.title }}</h3>
        </div>
        <p v-if="entry.description" class="card-body" :class="{ bold: entry.boldDesc }">{{ entry.description }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  entry: Object
})
</script>

<style scoped>
.entry-card-row {
  display: contents;
}

.time-col {
  grid-column: 1;
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  padding-right: 0;
  padding-top: 14px;
}

.time-text {
  font-size: 13px;
  color: #999;
  font-weight: 400;
  line-height: 20px;
  white-space: nowrap;
}

.dot-col {
  grid-column: 3;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 18px;
}

.hollow-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  border: 1.5px solid #a0a0a0;
  background: #ffffff;
  flex-shrink: 0;
  z-index: 2;
}

.card-col {
  grid-column: 5;
  min-width: 0;
  overflow: hidden;
}

.content-card {
  background: #fafafa;
  border: 1px solid #ebebeb;
  border-radius: 8px;
  padding: 14px 16px;
  margin-bottom: 14px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.cat-line {
  width: 3px;
  height: 14px;
  border-radius: 1.5px;
  flex-shrink: 0;
}

.cat-yellow { background: #f5a623; }
.cat-red    { background: #e85d75; }
.cat-green  { background: #4caf50; }

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: #2c2c2c;
  line-height: 1.4;
  margin: 0;
  word-break: break-word;
}

.card-body.bold {
  font-weight: 700;
  font-size: 15px;
  color: #1a1a1a;
}
.card-body {
  font-size: 13px;
  color: #666;
  line-height: 1.65;
  margin: 0;
  word-break: break-word;
  text-align: left;
}
</style>
