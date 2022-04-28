<script setup lang="ts">
import { useRadioMutation } from "../hooks";

import DButton from "./DaisyUI/DButton.vue";

defineProps({
  radio: {
    type: Object as () => Radio,
    required: true,
  },
});

const { mutate, isLoading } = useRadioMutation()
</script>

<template>
  <d-button v-if="radio.power" class="btn-success w-12" aria-label="Powered ON" :loading="isLoading"
    @click="() => mutate({ uuid: radio.uuid, power: false })">
    <v-icon name="fa-power-off" />
  </d-button>
  <d-button v-else class="btn-error w-12" aria-label="Powered OFF" :loading="isLoading"
    @click="() => mutate({ uuid: radio.uuid, power: true })">
    <v-icon name="fa-power-off" />
  </d-button>
</template>

<style>
</style>
