<script setup lang="ts">
import DButton from "./DaisyUI/DButton.vue";

const { status } = defineProps({
  status: {
    type: String,
    default: "",
  },
});
</script>

<template>
  <div class="tooltip tooltip-bottom" :data-tip="status">
    <d-button v-if="status == 'Connecting'" class="btn-circle btn-warning" :aria-label="status">
      <v-icon name="fa-sync" animation="spin" />
    </d-button>
    <d-button v-else-if="status == 'Playing'" class="btn-circle btn-success pl-1" :aria-label="status">
      <v-icon name="fa-play" />
    </d-button>
    <d-button v-else-if="status == 'Stopped'" class="btn-circle btn-error" :aria-label="status">
      <v-icon name="fa-stop" />
    </d-button>
    <d-button v-else class="btn-circle no-animation btn-info" :aria-label="status">
      <v-icon name="fa-question" />
    </d-button>
  </div>
</template>

<style>
</style>
