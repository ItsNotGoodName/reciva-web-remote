
<script setup lang="ts">
const { power } = defineProps({
  power: {
    type: Boolean,
    default: false,
  },
});
</script>

<template>
  <button v-if="power" class="btn btn-success" aria-label="Powered ON">
    <v-icon name="fa-power-off" />
  </button>
  <button v-else class="btn btn-error" aria-label="Powered OFF">
    <v-icon name="fa-power-off" />
  </button>
</template>

<style>
</style>
