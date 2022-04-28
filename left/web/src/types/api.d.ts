declare interface APIResponse<T> {
  ok: boolean
  code: number
  data: T
  error: {
    message: string
  }
}

declare interface RadioMutation {
  uuid: string,
  power?: boolean
  audio_source?: string
  preset?: number
  volume?: number
}

declare interface RadiosDiscoverMutation {
  force?: boolean
}