<script setup lang="ts">
defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
});
</script>

<template>
  <button class="btn touch-manipulation" :class="{ 'loading ': loading }">
    <slot v-if="!loading"></slot>
  </button>
</template>

<style>
</style>