<script setup lang="ts">
import { useRadioMutation } from "../hooks";

import Preset from "./Preset.vue";

defineProps({
  radio: {
    type: Object as () => Radio,
    required: true,
  },
});

const { mutate, isLoading } = useRadioMutation()
</script>

<template>
  <preset :key="p.number" v-for="p of radio.presets" :selected="radio.preset_number == p.number" :preset="p"
    :loading="isLoading" @click="() => mutate({ uuid: radio.uuid, preset: p.number })" />
</template>

<style>
</style>
