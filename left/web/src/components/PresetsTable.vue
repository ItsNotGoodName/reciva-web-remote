<script setup lang="ts">
import { usePresetsQuery } from "../hooks";

import DButton from "../components/DaisyUI/DButton.vue";
import DErrorAlert from "../components/DaisyUI/DErrorAlert.vue";

const { data, isLoading, isError, error } = usePresetsQuery();

const emit = defineEmits<{
  (e: 'update:presetUrl', url: string): void
}>()

defineProps({
  presetUrl: {
    type: String,
    default: "",
  },
});
</script>

<template>
  <div v-if="isLoading" class="flex">
    <v-icon class="mx-auto" name="fa-spinner" animation="spin" scale="2" />
  </div>
  <d-error-alert v-else-if="isError" :error="error">Failed to list presets.</d-error-alert>
  <table v-else-if="data" class="table table-compact">
    <thead>
      <tr>
        <th></th>
        <th>URL</th>
        <th>New Title</th>
        <th>New URL</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="p in data" :key="p.url" :class="{ 'active': p.url == presetUrl }">
        <th class="w-0">
          <d-button class="btn-success btn-sm" aria-label="Edit" @click="emit('update:presetUrl', p.url)">
            <v-icon name="fa-edit" />
          </d-button>
        </th>
        <td class="w-0">
          {{ p.url }}
        </td>
        <td class="w-0">
          {{ p.title_new }}
        </td>
        <td>
          {{ p.url_new }}
        </td>
      </tr>
    </tbody>
    <tfoot>
      <tr>
        <th></th>
        <th>URL</th>
        <th>New Title</th>
        <th>New URL</th>
      </tr>
    </tfoot>
  </table>
</template>

<style>
</style>
