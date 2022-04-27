<script setup lang="ts">
import DButton from "./DaisyUI/DButton.vue";

const { power } = defineProps({
  power: {
    type: Boolean,
    default: false,
  },
});
</script>

<template>
  <d-button v-if="power" class="btn-success" aria-label="Powered ON">
    <v-icon name="fa-power-off" />
  </d-button>
  <d-button v-else class="btn-error" aria-label="Powered OFF">
    <v-icon name="fa-power-off" />
  </d-button>
</template>

<style>
</style>
