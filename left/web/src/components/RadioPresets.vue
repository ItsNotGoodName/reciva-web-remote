<script setup lang="ts">
import { computed } from "vue";

import { useRadioMutation } from "../hooks";

import Preset from "./Preset.vue";

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
  <preset :key="p.number" v-for="p of radio.presets" :selected="radio.preset_number == p.number" :preset="p"
    :loading="loadingNumber == p.number" @click="() => mutate({ uuid: radio.uuid, preset: p.number })" />
</template>

<style>
</style>
