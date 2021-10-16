<script setup>
defineProps({
  setRadioPreset: Function,
  radio: Object
})
</script>

<template>
  <div class="grid gap-2 lg:grid-cols-4 md:grid-cols-3 sm:grid-cols-2">
    <button
      @click="setRadioPreset(i + 1)"
      class="rounded-xl bg-gray-100 hover:bg-gray-300 p-2"
      :value="i + 1"
      v-for="i in Array(radio.presets).keys() "
      :key="i + 1"
    >Preset {{ i + 1 }}</button>
  </div>
</template>