<script setup lang="ts">
import { ref } from 'vue'
import { useGame } from '@/composables/useGame'

const emit = defineEmits<{ loaded: [] }>()
const { load, error, clearError } = useGame()

const gameId = ref('')
const localError = ref('')

async function onSubmit() {
  localError.value = ''
  clearError()
  const id = gameId.value.trim()
  if (!id) {
    localError.value = 'Entrez un ID de partie.'
    return
  }
  try {
    await load(id)
    emit('loaded')
  } catch {
    localError.value = error.value ?? 'Partie introuvable'
  }
}
</script>

<template>
  <div class="panel">
    <h2>Charger une partie</h2>
    <form @submit.prevent="onSubmit">
      <div class="row">
        <div>
          <label for="load-id">ID de la partie</label>
          <input id="load-id" v-model="gameId" type="text" placeholder="uuid...">
        </div>
        <div style="display:flex;align-items:flex-end;">
          <button class="btn-secondary" type="submit" style="width:100%">Charger</button>
        </div>
      </div>
    </form>
    <p v-if="localError" class="error-msg">{{ localError }}</p>
  </div>
</template>
