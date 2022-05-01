<script setup lang="ts">
import { useStateMutation } from "../hooks";

import DButton from "./DaisyUI/DButton.vue";

defineProps({
  state: {
    type: Object as () => State,
    required: true,
  },
});

const { mutate, isLoading } = useStateMutation()
</script>

<template>
  <d-button v-if="state.power" class="btn-success w-14" aria-label="Powered ON" :loading="isLoading"
    @click="() => mutate({ uuid: state.uuid, power: false })">
    <v-icon name="fa-power-off" />
  </d-button>
  <d-button v-else class="btn-error w-14" aria-label="Powered OFF" :loading="isLoading"
    @click="() => mutate({ uuid: state.uuid, power: true })">
    <v-icon name="fa-power-off" />
  </d-button>
</template>

<style>
</style>
