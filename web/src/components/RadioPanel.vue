<script setup>
import {
  PlayIcon,
  VolumeUpIcon,
  VolumeOffIcon
} from '@heroicons/vue/solid'

defineProps({
  radio: Object,
  toggleRadioPower: Function,
  setRadioVolume: Function
})
</script>

<template>
  <div>
    <div class="space-x-2">
      <button
        @click="toggleRadioPower()"
        v-bind:class="{ 'bg-green-300': radio.power, 'bg-red-300': !radio.power }"
        class="inline-block w-16 rounded-xl p-2"
      >{{ radio.power ? "ON" : "OFF" }}</button>
      <button @click="setRadioVolume(radio.volume - 5)">
        <VolumeOffIcon class="inline-block h-8 w-8" />
      </button>
      <span class>{{ radio.volume }}</span>
      <button @click="setRadioVolume(radio.volume + 5)">
        <VolumeUpIcon class="inline-block h-8 w-8" />
      </button>
    </div>
    <div class="mt-2">
      <div class="inline-block bg-gray-300 p-2">{{ radio.state }}</div>
      <div class="inline-block bg-yellow-300 p-2">{{ radio.title }}</div>
      <div class="inline-block p-2">{{ radio.metadata }}</div>
    </div>
  </div>
</template>