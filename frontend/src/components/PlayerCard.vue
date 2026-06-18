<script setup lang="ts">
import type { Player } from '@/types/game'

defineProps<{ player: Player }>()
</script>

<template>
  <div
    class="player-card"
    :class="{
      dead: !player.alive,
      wolf: player.role === 'Wolf',
      villager: player.role === 'Villager',
    }"
  >
    <div class="player-icon">{{ player.role === 'Wolf' ? '&#129418;' : '&#127968;' }}</div>
    <div class="player-name">{{ player.name }}</div>
    <div class="player-role">{{ player.role }}</div>
    <div class="player-trait">{{ player.trait }}</div>
    <span v-if="player.mayor" class="badge badge-mayor">Maire</span>
    <span v-if="!player.alive" class="badge badge-dead">Mort</span>
  </div>
</template>

<style scoped>
.player-card {
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 14px;
  text-align: center;
  transition: all 0.3s;
}

.player-card.dead {
  opacity: 0.35;
  border-color: #333;
}

.player-card.wolf {
  border-color: var(--accent);
}

.player-card.villager {
  border-color: var(--green);
}

.player-icon {
  font-size: 1.6rem;
}

.player-name {
  font-weight: 700;
  font-size: 1.1rem;
  margin-bottom: 4px;
}

.player-role {
  font-size: 0.85rem;
  margin-bottom: 2px;
}

.player-trait {
  font-size: 0.8rem;
  color: var(--text-dim);
}

.badge {
  display: inline-block;
  font-size: 0.7rem;
  padding: 2px 8px;
  border-radius: 10px;
  margin-top: 6px;
  font-weight: 600;
}

.badge-mayor {
  background: var(--gold);
  color: var(--bg);
}

.badge-dead {
  background: #555;
  color: #ccc;
}
</style>
