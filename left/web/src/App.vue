<script lang="ts" setup>
import { onMounted, ref } from "vue";

import RadioStatus from "./components/RadioStatus.vue";
import RadioTitle from "./components/RadioTitle.vue";
import RadioPreset from "./components/RadioPreset.vue";
import RadioPower from "./components/RadioPower.vue";
import RadioName from "./components/RadioName.vue";
import DButton from "./components/DaisyUI/DButton.vue";
import RadioAudiosource from "./components/RadioAudiosource.vue";

const edit = ref(false)
const status = ref("")
const audiosources = ref(["Aux", "Internet radio"])
const audiosource = ref("")

const radio = ref({
  name: "Living Room",
  model_name: "Grace",
  model_number: "412"
})

onMounted(() => {
  setTimeout(() => { status.value = "Stopped" }, 2000);
  setTimeout(() => { status.value = "Connecting" }, 4000);
  setTimeout(() => { status.value = "Playing" }, 6000);
})

const setAudiosource = async (value: string) => {
  await new Promise(resolve => setTimeout(resolve, 1000))
  audiosource.value = value
}
</script>

<template>
  <div class="h-screen">
    <div class="navbar bg-base-200 fixed top-0 flex gap-2 z-50 border-b-2 border-b-base-300">
      <radio-status :status="status" />
      <radio-title class="flex-grow w-full" title="Lorem Ipsum" url="http://www.google.com"
        url_new="http://example.com" />
    </div>
    <div class="mx-5 pt-20 pb-36">
      <div v-if="!edit" class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        <radio-preset :selected="true" :number="1" title="Preset 1" />
        <radio-preset :number="2" title="Preset 2" />
        <radio-preset :number="3" title="Preset 3" />
        <radio-preset :number="4" title="Preset 4" />
        <radio-preset :number="5" title="Preset 5" />
        <radio-preset :number="6" title="Preset 6" />
        <radio-preset :number="7" title="Preset 7" />
        <radio-preset :number="8" title="Preset 8" />
        <radio-preset :number="9" title="Preset 9" />
        <radio-preset :number="10" title="Preset 11" />
        <radio-preset :number="11" title="Preset 11" />
      </div>
      <div v-else>
        Hello World
      </div>
    </div>
    <div
      class="navbar bg-base-200 fixed bottom-0 flex flex-row-reverse flex-wrap gap-2 pb-4 px-4 z-50 border-t-2 border-t-base-300">
      <!--- Radio Toolbar -->
      <div class="flex-grow md:flex-grow-0 flex gap-2">
        <radio-power class="flex-grow" :power="false" />
        <div v-if="true" class="btn-group">
          <d-button class="btn-info" aria-label="Volume Down">
            <v-icon name="fa-volume-down" />
          </d-button>
          <d-button class="btn-info px-0 w-10">100%</d-button>
          <d-button class="btn-info" aria-label="Volume Up">
            <v-icon name="fa-volume-up" />
          </d-button>
        </div>
        <d-button v-else class="btn-error" aria-label="Volume Muted">
          <v-icon name="fa-volume-mute" />
        </d-button>
        <radio-audiosource :audiosource="audiosource" :audiosources="audiosources" :set-audiosource="setAudiosource" />
        <radio-name :model_name="radio.model_name" :model_number="radio.model_number" :name="radio.name" />
      </div>
      <!--- Radios Toolbar -->
      <div class="grow flex gap-2">
        <div class="tooltip" data-tip="Edit Presets">
          <d-button :class="{ 'btn-success': edit }" aria-label="Edit Presets" @click="edit = !edit">
            <v-icon name="fa-edit" />
          </d-button>
        </div>
        <div class="grow flex">
          <div class="tooltip" data-tip="Discover">
            <d-button class="btn-primary rounded-none rounded-l-md" aria-label="Discover">
              <v-icon name="fa-search" />
            </d-button>
          </div>
          <select class="select select-primary rounded-none flex-grow">
            <option disabled selected>Select Radio</option>
          </select>
          <div class="tooltip" data-tip="Refresh">
            <d-button class="btn-primary rounded-none rounded-r-md" aria-label="Refresh">
              <v-icon name="fa-redo" />
            </d-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
</style>
