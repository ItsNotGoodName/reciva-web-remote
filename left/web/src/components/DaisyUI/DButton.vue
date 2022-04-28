<script setup lang="ts">
defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
});
</script>

<template>
  <button class="btn" :class="{ 'loading': loading }">
    <slot v-if="!loading"></slot>
  </button>
</template>

<style scoped>
button {
  touch-action: manipulation;
}
</style>