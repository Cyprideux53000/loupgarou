import { ref, computed } from 'vue'
import * as api from '@/api/client'
import type { Game, GameMode, Status } from '@/types/game'

export interface LogEntry {
  message: string
  type: 'night' | 'day' | 'end' | 'discussion' | 'vote'
}

const game = ref<Game | null>(null)
const status = ref<Status | null>(null)
const log = ref<LogEntry[]>([])
const discussion = ref<string[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

export function useGame() {
  const isGameOver = computed(() => status.value?.is_game_over ?? false)
  const gameId = computed(() => game.value?.id ?? null)
  const isDayPhase = computed(() => game.value != null && !game.value.night && !isGameOver.value)

  async function create(names: string[], wolfCount: number, mode: GameMode) {
    loading.value = true
    error.value = null
    try {
      log.value = []
      discussion.value = []
      game.value = await api.createGame({ names, wolf_count: wolfCount, mode })
      status.value = await api.getStatus(game.value.id)
      log.value.unshift({ message: 'Partie creee !', type: 'day' })
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Erreur inconnue'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function load(id: string) {
    loading.value = true
    error.value = null
    try {
      log.value = []
      discussion.value = []
      game.value = await api.getGame(id)
      status.value = await api.getStatus(id)
      log.value.unshift({ message: 'Partie chargee.', type: 'day' })
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Erreur inconnue'
      throw e
    } finally {
      loading.value = false
    }
  }

  function addDiscussionMessage(message: string) {
    discussion.value.push(message)
    log.value.unshift({ message, type: 'discussion' })
  }

  async function step() {
    if (!game.value) return
    loading.value = true
    error.value = null
    try {
      const response = await api.executeStep(game.value.id, discussion.value.length > 0 ? discussion.value : undefined)
      game.value = response.game
      status.value = await api.getStatus(game.value.id)
      discussion.value = []

      const logType: 'night' | 'day' = response.step.phase === 'wolfAttack' ? 'night' : 'day'

      if (response.step.discussion && response.step.discussion.length > 0) {
        for (const line of [...response.step.discussion].reverse()) {
          log.value.unshift({ message: line, type: 'discussion' })
        }
      }

      if (response.step.votes && response.step.votes.length > 0) {
        for (const vote of [...response.step.votes].reverse()) {
          log.value.unshift({ message: `${vote.voter} vote contre ${vote.target}`, type: 'vote' })
        }
      }

      log.value.unshift({ message: response.step.message, type: logType })

      if (response.step.new_mayor) {
        log.value.unshift({
          message: `Nouveau maire : ${response.step.new_mayor.name}`,
          type: logType,
        })
      }

      if (status.value.is_game_over) {
        const winMsg = status.value.winner === 'Wolves'
          ? 'Les Loups remportent la partie !'
          : 'Les Villageois remportent la partie !'
        log.value.unshift({ message: winMsg, type: 'end' })
      }
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Erreur inconnue'
    } finally {
      loading.value = false
    }
  }

  function clearError() {
    error.value = null
  }

  return { game, status, log, discussion, loading, error, isGameOver, isDayPhase, gameId, create, load, step, addDiscussionMessage, clearError }
}
