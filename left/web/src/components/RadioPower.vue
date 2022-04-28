<script setup lang="ts">
import { useRadioMutation } from "../hooks";
import DButton from "./DaisyUI/DButton.vue";

const { radio } = defineProps({
  radio: {
    type: Object as () => Radio,
    required: true,
  },
});

const { mutate, isLoading } = useRadioMutation(radio.uuid)
</script>

<template>
  <d-button v-if="radio.power" class="btn-success" aria-label="Powered ON" :loading="isLoading"
    @click="() => mutate({ power: false })">
    <v-icon name="fa-power-off" />
  </d-button>
  <d-button v-else class="btn-error" aria-label="Powered OFF" :loading="isLoading"
    @click="() => mutate({ power: true })">
    <v-icon name="fa-power-off" />
  </d-button>
</template>

<style>
</style>
