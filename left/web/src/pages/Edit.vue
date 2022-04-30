<script setup lang="ts">
import { ref } from "vue"

import PresetForm from "../components/PresetForm.vue";
import PresetsTable from "../components/PresetsTable.vue";

const presetFormOpen = ref(false)
const presetUrl = ref("")

const updatePresetUrl = (url: string) => {
  presetUrl.value = url
  presetFormOpen.value = true
  window.scrollTo(0, 0)
}
</script>

<template>
  <div class="flex gap-4 flex-wrap-reverse">
    <div class="lg:flex-1 flex overflow-x-auto">
      <presets-table class="w-full" :preset-url="presetUrl" @update:preset-url="updatePresetUrl" />
    </div>
    <div v-if="presetFormOpen" class="flex-1 lg:flex-initial w-96">
      <preset-form class="p-2 rounded-lg bg-base-200 border-2 border-base-200" :preset-url="presetUrl"
        @close="presetFormOpen = false" />
    </div>
  </div>
</template>

<style>
</style>
