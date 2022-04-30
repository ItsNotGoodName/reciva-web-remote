<script setup lang="ts">
import { computed } from "vue"

import { STATUS_CONNECTING, STATUS_PLAYING, STATUS_STOPPED } from "../constants";

import DButton from "./DaisyUI/DButton.vue";

const props = defineProps({
  radio: {
    type: Object as () => Radio,
    required: true,
  },
  loading: {
    type: Boolean,
    default: false,
  }
});

const status = computed(() => {
  return props.radio.status ? props.radio.status : "Unknown"
})

</script>

<template>
  <div class="tooltip tooltip-bottom" :data-tip="status">
    <d-button v-if="status == STATUS_CONNECTING" class="btn-circle btn-warning" :aria-label="status" :loading="loading">
      <v-icon name="fa-sync" animation="spin" />
    </d-button>
    <d-button v-else-if="status == STATUS_PLAYING" class="btn-circle btn-success pl-1" :aria-label="status"
      :loading="loading">
      <v-icon name="fa-play" />
    </d-button>
    <d-button v-else-if="status == STATUS_STOPPED" class="btn-circle btn-error" :aria-label="status" :loading="loading">
      <v-icon name="fa-stop" />
    </d-button>
    <d-button v-else class="btn-circle no-animation btn-info" :aria-label="status" :loading="loading">
      <v-icon name="fa-question" />
    </d-button>
  </div>
</template>

<style>
</style>
