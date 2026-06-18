<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useGame } from '@/composables/useGame'
import PhaseBar from '@/components/PhaseBar.vue'
import StatusBar from '@/components/StatusBar.vue'
import PlayerCard from '@/components/PlayerCard.vue'
import GameLog from '@/components/GameLog.vue'

const route = useRoute()
const router = useRouter()
const { game, status, log, loading, isGameOver, isDayPhase, step, load, addDiscussionMessage } = useGame()

const chatInput = ref('')
const selectedSpeaker = ref('')

function sendMessage() {
  const msg = chatInput.value.trim()
  if (!msg || !selectedSpeaker.value) return
  addDiscussionMessage(`${selectedSpeaker.value}: ${msg}`)
  chatInput.value = ''
}

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

    <div v-if="isDayPhase && game.mode === 'llm'" class="chat-box">
      <p class="chat-label">Discussion du village</p>
      <form class="chat-input-row" @submit.prevent="sendMessage">
        <select v-model="selectedSpeaker" class="chat-speaker">
          <option value="" disabled>Joueur</option>
          <option
            v-for="p in game.players.filter(p => p.alive)"
            :key="p.id"
            :value="p.name"
          >
            {{ p.name }}
          </option>
        </select>
        <input
          v-model="chatInput"
          type="text"
          placeholder="Ecrivez un message..."
          class="chat-input"
        >
        <button type="submit" class="chat-send" :disabled="!selectedSpeaker || !chatInput.trim()" title="Envoyer">&#10004;</button>
      </form>
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

.chat-box {
  background: #111;
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 14px;
  margin-bottom: 14px;
}

.chat-label {
  font-size: 0.85rem;
  color: #74b9ff;
  font-weight: 600;
  margin-bottom: 10px;
}

.chat-input-row {
  display: flex;
  gap: 8px;
}

.chat-speaker {
  width: 130px;
  flex-shrink: 0;
  margin-bottom: 0;
  padding: 10px;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: var(--bg);
  color: var(--text);
  font-size: 0.9rem;
}

.chat-input {
  flex: 1;
  margin-bottom: 0;
}

.chat-send {
  background: var(--green);
  color: var(--bg);
  border: none;
  border-radius: 8px;
  padding: 10px 16px;
  font-size: 1.1rem;
  cursor: pointer;
  font-weight: 700;
  transition: opacity 0.2s;
}

.chat-send:hover {
  opacity: 0.85;
}
</style>
