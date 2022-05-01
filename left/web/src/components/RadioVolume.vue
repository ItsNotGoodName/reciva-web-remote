<script setup lang="ts">
import { useStateMutation, useRadioVolumeMutation } from "../hooks";

import DButton from "./DaisyUI/DButton.vue";

defineProps({
  state: {
    type: Object as () => State,
    required: true,
  },
});

const { mutate, isLoading } = useStateMutation()
const { mutate: refreshVolume, isLoading: refreshVolumeLoading } = useRadioVolumeMutation()
</script>

<template>
  <div v-if="!state.is_muted" class="btn-group flex-nowrap">
    <d-button class="btn-info w-14" aria-label="Volume Down" :loading="isLoading"
      @click="mutate({ uuid: state.uuid, volume: state.volume - 5 })">
      <v-icon name="fa-volume-down" />
    </d-button>
    <d-button class="btn-info px-0 w-12" :loading="refreshVolumeLoading" @click="() => refreshVolume(state.uuid)">
      {{ state.volume }}%
    </d-button>
    <d-button class="btn-info w-14" aria-label="Volume Up" :loading="isLoading"
      @click="mutate({ uuid: state.uuid, volume: state.volume + 5 })">
      <v-icon name="fa-volume-up" />
    </d-button>
  </div>
  <d-button v-else class="btn-error" aria-label="Volume Muted">
    <v-icon name="fa-volume-mute" />
  </d-button>
</template>

<style>
</style>
