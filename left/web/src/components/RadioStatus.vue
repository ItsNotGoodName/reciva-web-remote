<script setup lang="ts">
const { status } = defineProps({
  status: {
    type: String,
    default: "",
  },
});
</script>

<template>
  <div class="tooltip tooltip-bottom" :data-tip="status">
    <button v-if="status == 'Connecting'" class="btn btn-circle btn-warning" :aria-label="status">
      <v-icon name="fa-sync" animation="spin" />
    </button>
    <button v-else-if="status == 'Playing'" class="btn btn-circle btn-success" :aria-label="status">
      <v-icon name="fa-play" />
    </button>
    <button v-else-if="status == 'Stopped'" class="btn btn-circle btn-error" :aria-label="status">
      <v-icon name="fa-stop" />
    </button>
    <button v-else class="btn btn-circle no-animation btn-info" :aria-label="status">
      <v-icon name="fa-question" />
    </button>
  </div>
</template>

<style>
</style>
