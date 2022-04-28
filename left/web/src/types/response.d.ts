declare interface APIResponse<T> {
  ok: boolean
  code: number
  data: T
  error: {
    message: string
  }
}