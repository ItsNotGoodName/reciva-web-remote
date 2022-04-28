<script setup lang="ts">
import { STATUS_CONNECTING, STATUS_PLAYING, STATUS_STOPPED } from "../constants";

import DButton from "./DaisyUI/DButton.vue";

defineProps({
  radio: {
    type: Object as () => Radio,
    default: { status: "Unknown" }
  },
  loading: {
    type: Boolean,
    default: false,
  }
});
</script>

<template>
  <div class="tooltip tooltip-bottom" :data-tip="radio.status">
    <d-button v-if="radio.status == STATUS_CONNECTING" class="btn-circle btn-warning" :aria-label="radio.status"
      :loading="loading">
      <v-icon name="fa-sync" animation="spin" />
    </d-button>
    <d-button v-else-if="radio.status == STATUS_PLAYING" class="btn-circle btn-success pl-1" :aria-label="radio.status"
      :loading="loading">
      <v-icon name="fa-play" />
    </d-button>
    <d-button v-else-if="radio.status == STATUS_STOPPED" class="btn-circle btn-error" :aria-label="radio.status"
      :loading="loading">
      <v-icon name="fa-stop" />
    </d-button>
    <d-button v-else class="btn-circle no-animation btn-info" :aria-label="radio.status" :loading="loading">
      <v-icon name="fa-question" />
    </d-button>
  </div>
</template>

<style>
</style>
