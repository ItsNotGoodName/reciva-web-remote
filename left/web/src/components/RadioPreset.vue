<script setup lang="ts">
import DButton from "./DaisyUI/DButton.vue";

defineProps({
  preset: {
    type: Object as () => RadioPreset,
    required: true,
  },
  selected: {
    type: Boolean,
    required: true,
  },
});
</script>

<template>
  <d-button class="flex w-full" :class="{ 'btn-primary': selected }">
    <span class="badge badge-info badge-lg rounded-md mr-2">{{ preset.number }}</span>
    <span class="flex-grow w-0 truncate ">{{ preset.title_new ? preset.title_new : preset.title }}</span>
  </d-button>
</template>

<style>
</style>
