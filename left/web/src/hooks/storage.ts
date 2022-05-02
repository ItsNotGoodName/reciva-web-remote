import { watch, ref } from "vue"

export function useRadioUUIDStorage() {
  const radioUUID = ref(localStorage.getItem("lastRadioUUID") || "")

  watch(radioUUID, (newRadioUUID) => {
    localStorage.setItem("lastRadioUUID", newRadioUUID)
  })

  return radioUUID
}
