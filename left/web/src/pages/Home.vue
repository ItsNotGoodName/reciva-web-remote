<script setup lang="ts">
import { computed } from "vue";

import { useStateMutation } from "../hooks";

import RadioPreset from "../components/RadioPreset.vue";

defineProps({
  state: {
    type: Object as () => State,
    required: true,
  },
});

const { mutate, isLoading, variables } = useStateMutation()

const loadingNumber = computed(() => {
  if (!isLoading.value) {
    return null;
  }

  return variables.value?.preset
})
</script>

<template>
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
    <radio-preset :key="p.number" v-for="p of state.presets" :selected="state.preset_number == p.number" :preset="p"
      :loading="loadingNumber == p.number" @click="() => mutate({ uuid: state.uuid, preset: p.number })" />
  </div>
</template>

<style>
</style>
