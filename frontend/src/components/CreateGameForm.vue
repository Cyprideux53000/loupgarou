<script setup lang="ts">
import { ref } from 'vue'
import { useGame } from '@/composables/useGame'

const emit = defineEmits<{ created: [] }>()
const { create, error, clearError } = useGame()

const playerNames = ref('Alice, Bob, Charlie, Diana, Eve')
const wolfCount = ref(2)
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
    await create(names, wolfCount.value)
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
        <div style="display:flex;align-items:flex-end;">
          <button class="btn-primary" type="submit" style="width:100%">Creer la partie</button>
        </div>
      </div>
    </form>
    <p v-if="localError" class="error-msg">{{ localError }}</p>
  </div>
</template>
