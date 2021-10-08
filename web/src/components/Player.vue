<script setup>
import {
  PlayIcon,
  VolumeUpIcon,
  VolumeOffIcon
} from '@heroicons/vue/solid'


defineProps({
  radio: Object
})
</script>

<template>
  <div class="flex bottom-0 fixed inset-x-0 p-4 bg-blue-300">
    <PlayIcon class="flex-shrink h-6 w-6 mr-2" />
    <div class="flex-grow">{{ radio.title }} | {{ radio.metadata }}</div>
    <VolumeUpIcon v-if="!radio.isMuted" class="flex-shrink h-6 w-6 mr-2" />
    <VolumeOffIcon v-else class="flex-shrink h-6 w-6 mr-2" />
    <div class="flex-shrink">{{ radio.volume }}</div>
  </div>
</template>