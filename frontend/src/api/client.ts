import type { Game, Status, StepResponse, CreateGameRequest } from '@/types/game'

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(url, options)
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text.trim() || `Erreur ${res.status}`)
  }
  return res.json()
}

export function createGame(req: CreateGameRequest): Promise<Game> {
  return request<Game>('/game', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req),
  })
}

export function getGame(id: string): Promise<Game> {
  return request<Game>(`/game?id=${encodeURIComponent(id)}`)
}

export function getStatus(id: string): Promise<Status> {
  return request<Status>(`/status?id=${encodeURIComponent(id)}`)
}

export function executeStep(id: string): Promise<StepResponse> {
  return request<StepResponse>(`/step?id=${encodeURIComponent(id)}`, {
    method: 'POST',
  })
}
