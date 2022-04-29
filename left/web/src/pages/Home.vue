<script setup lang="ts">
import { computed } from "vue";

import { useRadioMutation } from "../hooks";

import RadioPreset from "../components/RadioPreset.vue";

defineProps({
  radio: {
    type: Object as () => Radio,
    required: true,
  },
});

const { mutate, isLoading, variables } = useRadioMutation()

const loadingNumber = computed(() => {
  if (!isLoading.value) {
    return undefined;
  }

  return variables.value?.preset
})
</script>

<template>
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
    <radio-preset :key="p.number" v-for="p of radio.presets" :selected="radio.preset_number == p.number" :preset="p"
      :loading="loadingNumber == p.number" @click="() => mutate({ uuid: radio.uuid, preset: p.number })" />
  </div>
</template>

<style>
</style>
