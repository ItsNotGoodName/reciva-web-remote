declare type APIResponse<T> = APIData<T> | APIError;

interface APIData<T> {
  ok: true
  code: number
  data: T
}
interface APIError {
  ok: false
  code: number
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

declare interface PresetMutation {
  title_new: string,
  url: string,
  url_new: string,
}
