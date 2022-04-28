<script setup lang="ts">
import { useRadioMutation } from "../hooks";
import DButton from './DaisyUI/DButton.vue';
import DDropdownButton from './DaisyUI/DDropdownButton.vue';

const { radio } = defineProps({
  radio: {
    type: Object as () => Radio,
    required: true,
  }
});

const { mutate, isLoading } = useRadioMutation(radio.uuid)
</script>

<template>
  <div class="dropdown dropdown-top dropdown-end">
    <d-dropdown-button :class="{ 'btn-secondary': radio.audio_source }" aria-label="Audiosource">
      <v-icon name="fa-itunes-note" />
    </d-dropdown-button>
    <ul tabindex="0" class="menu menu-compact dropdown-content mb-2 p-2 shadow bg-base-200 rounded-box w-52 space-y-2">
      <span class="mx-auto">Audio Source</span>
      <d-button :loading="isLoading" :key="a" v-for="a in radio.audio_sources"
        :class="{ 'btn-secondary': a == radio.audio_source }" @click="() => { mutate({ audio_source: a }) }">
        {{ a }}
      </d-button>
    </ul>
  </div>
</template>

<style>
</style>

