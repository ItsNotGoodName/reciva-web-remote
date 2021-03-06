declare type State = {
  audio_source: string
  audio_sources: string[]
  is_muted: boolean
  metadata: string
  model_name: string
  model_number: string
  name: string
  power: boolean
  preset_number: number
  presets: StatePreset[]
  status: string
  title: string
  title_new: string
  url: string
  url_new: string
  uuid: string
  volume: number
}

declare type StatePreset = {
  number: number
  title: string
  title_new: string
  url: string
  url_new: string
}