<script setup lang="ts">
import DButton from "./DaisyUI/DButton.vue";

defineProps({
  radio: {
    type: Object as () => Radio,
    required: true,
  },
});
</script>

<template>
  <div v-if="!radio.is_muted" class="btn-group">
    <d-button class="btn-info" aria-label="Volume Down">
      <v-icon name="fa-volume-down" />
    </d-button>
    <d-button class="btn-info px-0 w-10">{{ radio.volume }}%</d-button>
    <d-button class="btn-info" aria-label="Volume Up">
      <v-icon name="fa-volume-up" />
    </d-button>
  </div>
  <d-button v-else class="btn-error" aria-label="Volume Muted">
    <v-icon name="fa-volume-mute" />
  </d-button>
</template>

<style>
</style>
