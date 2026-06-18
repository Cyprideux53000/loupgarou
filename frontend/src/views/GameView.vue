<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useGame } from '@/composables/useGame'
import PhaseBar from '@/components/PhaseBar.vue'
import StatusBar from '@/components/StatusBar.vue'
import PlayerCard from '@/components/PlayerCard.vue'
import GameLog from '@/components/GameLog.vue'

const route = useRoute()
const router = useRouter()
const { game, status, log, loading, isGameOver, step, load } = useGame()

onMounted(async () => {
  const id = route.params.id as string
  if (!game.value || game.value.id !== id) {
    try {
      await load(id)
    } catch {
      router.push({ name: 'home' })
    }
  }
})
</script>

<template>
  <div v-if="game && status" class="panel">
    <p class="game-id-display">
      ID: {{ game.id }}
      <span class="mode-badge" :class="game.mode">{{ game.mode === 'llm' ? 'LLM (Ollama)' : 'Random' }}</span>
    </p>

    <PhaseBar :night="game.night" :current-step="game.current_step" :status="status" />
    <StatusBar :status="status" />

    <div class="players-grid">
      <PlayerCard v-for="player in game.players" :key="player.id" :player="player" />
    </div>

    <button class="btn-step" :disabled="isGameOver || loading" @click="step">
      {{ loading ? 'En cours...' : 'Jouer le prochain tour' }}
    </button>

    <GameLog :entries="log" />
  </div>

  <div v-else class="panel">
    <p>Chargement...</p>
  </div>
</template>

<style scoped>
.game-id-display {
  font-size: 0.8rem;
  color: #666;
  word-break: break-all;
  margin-bottom: 16px;
}

.mode-badge {
  display: inline-block;
  font-size: 0.75rem;
  padding: 2px 10px;
  border-radius: 10px;
  margin-left: 8px;
  font-weight: 600;
  vertical-align: middle;
}

.mode-badge.llm {
  background: var(--purple);
  color: #fff;
}

.mode-badge.random {
  background: var(--border);
  color: var(--text-dim);
}

.players-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 12px;
  margin-bottom: 20px;
}
</style>
