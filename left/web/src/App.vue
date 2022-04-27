<script lang="ts" setup>
import { onMounted, ref } from "vue";

import RadioStatus from "./components/RadioStatus.vue";
import RadioTitle from "./components/RadioTitle.vue";
import RadioPreset from "./components/RadioPreset.vue";
import RadioPower from "./components/RadioPower.vue";
import RadioName from "./components/RadioName.vue";

const status = ref("Connecting")

onMounted(() => {
  setTimeout(() => { status.value = "Stopped" }, 3000);
})

const edit = ref(false)

</script>

<template>
  <div class="h-screen">
    <div class="navbar bg-base-300 fixed top-0 flex gap-2">
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
      </div>
      <div v-else>
        Hello World
      </div>
    </div>
    <div class="navbar bg-base-300 fixed bottom-0 flex flex-row-reverse flex-wrap gap-2">
      <div class="grow flex gap-2">
        <radio-power class="flex-grow" :power="false" />
        <div v-if="true" class="btn-group">
          <button class="btn btn-info" aria-label="Volume Down">
            <v-icon name="fa-volume-down" />
          </button>
          <button class="btn btn-info">40%</button>
          <button class="btn btn-info" aria-label="Volume Up">
            <v-icon name="fa-volume-up" />
          </button>
        </div>
        <button v-else class="btn btn-error" aria-label="Volume Muted">
          <v-icon name="fa-volume-mute" />
        </button>
        <radio-name model_name="Grace" model_number="423" name="Room" />
      </div>
      <div class="grow flex gap-2">
        <div class="tooltip" data-tip="Edit Presets">
          <button class="btn" :class="{ 'btn-success': edit }" aria-label="Edit Presets" @click="edit = !edit">
            <v-icon name="fa-edit" />
          </button>
        </div>
        <div class="flex grow">
          <div class="tooltip" data-tip="Discover">
            <button class="btn btn-primary rounded-none rounded-l-md" aria-label="Discover">
              <v-icon name="fa-search" />
            </button>
          </div>
          <select class="select select-primary rounded-none flex-grow">
            <option disabled selected>Select Radio</option>
          </select>
          <div class="tooltip" data-tip="Refresh">
            <button class="btn btn-primary rounded-none rounded-r-md" aria-label="Refresh">
              <v-icon name="fa-redo" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
</style>
