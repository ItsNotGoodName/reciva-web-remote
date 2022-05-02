declare type APIResponse<T> = APIData<T> | APIError;

type APIData<T> = {
  ok: true
  data: T
}
type APIError = {
  ok: false
  error: {
    message: string
  }
}

declare type RadioMutation = {
  uuid: string,
  power?: boolean
  audio_source?: string
  preset?: number
  volume?: number
}

declare type RadiosDiscoverMutation = {
  force?: boolean
}

declare type PresetMutation = {
  title_new: string,
  url: string,
  url_new: string,
}
