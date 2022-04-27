<script setup lang="ts">
const { loading } = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
});
</script>

<template>
  <label tabindex="0" class="btn" :class="{ 'loading': loading }">
    <slot v-if="!loading" />
  </label>
</template>

<style scoped>
button {
  touch-action: manipulation;
}
</style>