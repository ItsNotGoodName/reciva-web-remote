<script setup>
import {
  VolumeUpIcon,
  VolumeOffIcon,
} from '@heroicons/vue/outline'

defineProps({
  radio: Object,
  toggleRadioPower: Function,
  setRadioVolume: Function,
  renewRadio: Function
})
</script>

<template>
  <div>
    <div class="flex">
      <div class="flex space-x-2">
        <button
          v-if="radio.power"
          @click="toggleRadioPower()"
          class="bg-green-300 hover:bg-green-500 inline-block w-16 rounded-xl p-2"
        >On</button>
        <button
          v-else
          @click="toggleRadioPower()"
          class="bg-red-300 hover:bg-red-500 inline-block w-16 rounded-xl p-2"
        >OFF</button>
        <div v-if="!radio.isMuted" class="space-x-1 flex">
          <button
            class="hover:bg-gray-300 p-1 rounded-full"
            @click="setRadioVolume(radio.volume - 5)"
          >
            <VolumeOffIcon class="inline-block h-8 w-8" />
          </button>
          <div class="m-auto w-6">{{ radio.volume }}</div>
          <button
            class="hover:bg-gray-300 p-1 rounded-full"
            @click="setRadioVolume(radio.volume + 5)"
          >
            <VolumeUpIcon class="inline-block h-8 w-8" />
          </button>
        </div>
        <div v-else class="bg-red-300 p-2">Muted</div>
      </div>
      <button
        class="ml-auto bg-yellow-300 hover:bg-yellow-500 rounded-xl p-2"
        @click="renewRadio()"
      >Renew</button>
    </div>
    <div class="mt-2">
      <div class="inline-block bg-gray-300 p-2">{{ radio.state }}</div>
      <div class="inline-block bg-yellow-300 p-2">{{ radio.title }}</div>
      <div class="inline-block p-2">{{ radio.metadata }}</div>
    </div>
  </div>
</template>