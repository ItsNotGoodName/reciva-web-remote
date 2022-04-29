<script setup lang="ts">
import { useRadioMutation, useRadioVolumeMutation } from "../hooks";

import DButton from "./DaisyUI/DButton.vue";

defineProps({
  radio: {
    type: Object as () => Radio,
    required: true,
  },
});

const { mutate, isLoading } = useRadioMutation()
const { mutate: refreshVolume, isLoading: refreshVolumeLoading } = useRadioVolumeMutation()
</script>

<template>
  <div v-if="!radio.is_muted" class="btn-group">
    <d-button class="btn-info w-14" aria-label="Volume Down" :loading="isLoading"
      @click="mutate({ uuid: radio.uuid, volume: radio.volume - 5 })">
      <v-icon name="fa-volume-down" />
    </d-button>
    <d-button class="btn-info px-0 w-12" :loading="refreshVolumeLoading" @click="() => refreshVolume(radio.uuid)">
      {{ radio.volume }}%
    </d-button>
    <d-button class="btn-info w-14" aria-label="Volume Up" :loading="isLoading"
      @click="mutate({ uuid: radio.uuid, volume: radio.volume + 5 })">
      <v-icon name="fa-volume-up" />
    </d-button>
  </div>
  <d-button v-else class="btn-error" aria-label="Volume Muted">
    <v-icon name="fa-volume-mute" />
  </d-button>
</template>

<style>
</style>
