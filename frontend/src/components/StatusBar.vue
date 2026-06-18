<script setup lang="ts">
import type { Status } from '@/types/game'

defineProps<{ status: Status }>()
</script>

<template>
  <div class="status-bar">
    <div class="stat stat-wolves">
      <div class="stat-value">{{ status.wolves_alive }}</div>
      <div class="stat-label">Loups vivants</div>
    </div>
    <div class="stat stat-villagers">
      <div class="stat-value">{{ status.villagers_alive }}</div>
      <div class="stat-label">Villageois vivants</div>
    </div>
  </div>
</template>

<style scoped>
.status-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.stat {
  flex: 1;
  text-align: center;
  padding: 12px;
  border-radius: 8px;
  background: var(--bg);
  border: 1px solid var(--border);
}

.stat-value {
  font-size: 1.8rem;
  font-weight: 700;
}

.stat-label {
  font-size: 0.8rem;
  color: var(--text-dim);
  margin-top: 4px;
}

.stat-wolves .stat-value {
  color: var(--accent);
}

.stat-villagers .stat-value {
  color: var(--green);
}
</style>
