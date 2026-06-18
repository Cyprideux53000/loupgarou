<script setup lang="ts">
import type { LogEntry } from '@/composables/useGame'

defineProps<{ entries: LogEntry[] }>()
</script>

<template>
  <div v-if="entries.length" class="log">
    <div
      v-for="(entry, i) in entries"
      :key="i"
      class="log-entry"
      :class="{
        'log-night': entry.type === 'night',
        'log-day': entry.type === 'day',
        'log-end': entry.type === 'end',
      }"
    >
      &gt; {{ entry.message }}
    </div>
  </div>
</template>

<style scoped>
.log {
  background: #111;
  border-radius: 8px;
  padding: 14px;
  max-height: 200px;
  overflow-y: auto;
  font-family: 'Courier New', monospace;
  font-size: 0.85rem;
  margin-top: 16px;
  border: 1px solid #222;
}

.log-entry {
  padding: 4px 0;
  border-bottom: 1px solid var(--bg);
}

.log-entry:last-child {
  border-bottom: none;
}

.log-night {
  color: #a29bfe;
}

.log-day {
  color: #fdcb6e;
}

.log-end {
  color: var(--accent);
  font-weight: 700;
}
</style>
