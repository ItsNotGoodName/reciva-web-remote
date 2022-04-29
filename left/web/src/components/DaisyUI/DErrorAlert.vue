<script setup lang="ts">
import { PropType, } from "vue"

defineProps({
  error: {
    type: Error as PropType<any>,
  },
});
</script>

<template>
  <div class="alert alert-error shadow-lg">
    <div>
      <v-icon name="fa-times-circle" />
      <slot />
      <span>{{ error }}</span>
    </div>
  </div>
</template>

<style>
</style>