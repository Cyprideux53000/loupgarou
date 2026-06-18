export type Role = 'Wolf' | 'Villager'
export type Trait = 'Cunning' | 'Aggressive' | 'Brave' | 'Timid' | 'Sly'
export type Phase = 'wolfAttack' | 'DayVote'
export type NextStep = 'wolves_attack' | 'village_vote' | ''
export type Winner = 'Wolves' | 'Villagers' | ''

export interface Player {
  id: string
  name: string
  role: Role
  trait: Trait
  alive: boolean
  mayor: boolean
}

export type GameMode = 'random' | 'llm'

export interface Game {
  id: string
  players: Player[]
  wolf_number: number
  night: boolean
  current_step: Phase
  mode: GameMode
}

export interface Status {
  wolves_alive: number
  villagers_alive: number
  next_step: NextStep
  is_game_over: boolean
  winner: Winner
}

export interface PlayerVote {
  voter: string
  target: string
}

export interface StepResult {
  victim: Player
  phase: Phase
  message: string
  new_mayor?: Player
  discussion?: string[]
  votes?: PlayerVote[]
}

export interface StepResponse {
  game: Game
  step: StepResult
}

export interface CreateGameRequest {
  names: string[]
  wolf_count: number
  mode: GameMode
}
