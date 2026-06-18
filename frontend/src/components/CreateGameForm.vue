<script setup lang="ts">
import { ref } from 'vue'
import { useGame } from '@/composables/useGame'
import type { GameMode } from '@/types/game'

const emit = defineEmits<{ created: [] }>()
const { create, error, clearError } = useGame()

const playerNames = ref('Alice, Bob, Charlie, Diana, Eve')
const wolfCount = ref(2)
const mode = ref<GameMode>('random')
const localError = ref('')

async function onSubmit() {
  localError.value = ''
  clearError()
  const names = playerNames.value.split(',').map(n => n.trim()).filter(n => n)
  if (names.length < 2) {
    localError.value = 'Il faut au moins 2 joueurs.'
    return
  }
  if (wolfCount.value < 1 || wolfCount.value >= names.length) {
    localError.value = 'Nombre de loups invalide.'
    return
  }
  try {
    await create(names, wolfCount.value, mode.value)
    emit('created')
  } catch {
    localError.value = error.value ?? 'Erreur inconnue'
  }
}
</script>

<template>
  <div class="panel">
    <h2>Nouvelle Partie</h2>
    <form @submit.prevent="onSubmit">
      <label for="player-names">Noms des joueurs (separes par des virgules)</label>
      <input id="player-names" v-model="playerNames" type="text" placeholder="Alice, Bob, Charlie, Diana, Eve">
      <div class="row">
        <div>
          <label for="wolf-count">Nombre de loups</label>
          <input id="wolf-count" v-model.number="wolfCount" type="number" min="1">
        </div>
        <div>
          <label for="game-mode">Mode de jeu</label>
          <div class="mode-selector">
            <button
              type="button"
              class="mode-btn"
              :class="{ active: mode === 'random' }"
              @click="mode = 'random'"
            >
              Random
            </button>
            <button
              type="button"
              class="mode-btn"
              :class="{ active: mode === 'llm' }"
              @click="mode = 'llm'"
            >
              LLM (Ollama)
            </button>
          </div>
        </div>
      </div>
      <button class="btn-primary" type="submit" style="width:100%; margin-top: 14px;">Creer la partie</button>
    </form>
    <p v-if="localError" class="error-msg">{{ localError }}</p>
  </div>
</template>

<style scoped>
.mode-selector {
  display: flex;
  gap: 0;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid var(--border);
}

.mode-btn {
  flex: 1;
  padding: 10px 14px;
  border: none;
  border-radius: 0;
  background: var(--bg);
  color: var(--text-dim);
  font-size: 0.9rem;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s;
}

.mode-btn.active {
  background: var(--accent);
  color: #fff;
}

.mode-btn:hover:not(.active) {
  background: var(--border);
  color: var(--text);
}
</style>
