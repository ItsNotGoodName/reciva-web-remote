<script setup lang="ts">
import { useRadiosDiscoverMutation } from "../hooks";

import DButton from "./DaisyUI/DButton.vue";

const { mutate, isLoading } = useRadiosDiscoverMutation()
</script>

<template>
  <div class="tooltip" data-tip="Discover">
    <d-button v-bind="$attrs" aria-label="Discover" :loading="isLoading" @click="() => mutate({})">
      <v-icon name="fa-search" />
    </d-button>
  </div>
</template>

<style>
</style>
