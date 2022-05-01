<script setup lang="ts">
import { computed } from "vue"

import { useStateMutation } from "../hooks";

import DButton from './DaisyUI/DButton.vue';
import DDropdownButton from './DaisyUI/DDropdownButton.vue';

defineProps({
  state: {
    type: Object as () => State,
    required: true,
  }
});

const { mutate, isLoading, variables } = useStateMutation()
const loadingAudioSource = computed(() => {
  if (!isLoading.value) {
    return null;
  }

  return variables.value?.audio_source
})
</script>

<template>
  <div class="dropdown dropdown-top dropdown-end">
    <d-dropdown-button :class="{ 'btn-secondary': state.audio_source }" aria-label="Audio Source">
      <v-icon name="fa-itunes-note" />
    </d-dropdown-button>
    <ul tabindex="0" class="menu menu-compact dropdown-content mb-2 p-2 shadow bg-base-200 rounded-box w-52 space-y-2">
      <span class="mx-auto">Audio Source</span>
      <d-button :loading="loadingAudioSource == a" :key="a" v-for="a in state.audio_sources"
        :class="{ 'btn-secondary': a == state.audio_source }"
        @click="() => { mutate({ uuid: state.uuid, audio_source: a }) }">
        {{ a }}
      </d-button>
    </ul>
  </div>
</template>

<style>
</style>

