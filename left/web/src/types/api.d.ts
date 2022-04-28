declare interface APIResponse<T> {
  ok: boolean
  code: number
  data: T
  error: {
    message: string
  }
}

declare interface RadioMutation {
  power?: boolean
  audio_source?: string
  preset?: number
  volume?: number
}