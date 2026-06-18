<script setup lang="ts">
import type { Status } from '@/types/game'

defineProps<{
  night: boolean
  currentStep: string
  status: Status
}>()
</script>

<template>
  <div
    class="phase-bar"
    :class="{
      'phase-over': status.is_game_over,
      'phase-night': !status.is_game_over && night,
      'phase-day': !status.is_game_over && !night,
    }"
  >
    <span class="phase-text">
      <template v-if="status.is_game_over">Partie terminee</template>
      <template v-else-if="night">&#127769; Nuit</template>
      <template v-else>&#9728;&#65039; Jour</template>
    </span>
    <span class="phase-step">
      <template v-if="status.is_game_over && status.winner === 'Wolves'">&#129418; Les Loups gagnent !</template>
      <template v-else-if="status.is_game_over && status.winner === 'Villagers'">&#127968; Les Villageois gagnent !</template>
      <template v-else-if="night">Les loups attaquent...</template>
      <template v-else>Vote du village</template>
    </span>
  </div>
</template>

<style scoped>
.phase-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 20px;
  border-radius: 10px;
  margin-bottom: 16px;
  font-weight: 600;
}

.phase-night {
  background: linear-gradient(135deg, #1a1a2e, #16213e);
  border: 1px solid var(--purple);
}

.phase-day {
  background: linear-gradient(135deg, #e9a820, #d4830a);
  color: var(--bg);
  border: 1px solid #d4830a;
}

.phase-over {
  background: linear-gradient(135deg, var(--accent), #c0392b);
  border: 1px solid var(--accent);
}
</style>
