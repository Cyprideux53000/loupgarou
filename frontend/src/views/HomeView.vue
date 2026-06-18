<script setup lang="ts">
import { useRouter } from 'vue-router'
import CreateGameForm from '@/components/CreateGameForm.vue'
import LoadGameForm from '@/components/LoadGameForm.vue'
import { useGame } from '@/composables/useGame'

const router = useRouter()
const { gameId } = useGame()

function onGameReady() {
  if (gameId.value) {
    router.push({ name: 'game', params: { id: gameId.value } })
  }
}
</script>

<template>
  <CreateGameForm @created="onGameReady" />
  <LoadGameForm @loaded="onGameReady" />
</template>
