<script setup lang="ts">
import { PropType, } from "vue"

import { usePresetsQuery } from "../hooks";
import { PAGE_HOME, } from "../constants"

import DButton from "../components/DaisyUI/DButton.vue";
import DErrorAlert from "../components/DaisyUI/DErrorAlert.vue";

defineProps({
  setPage: {
    type: Function as PropType<(page: string) => void>,
    required: true,
  },
});

const { data, isLoading, isError, error } = usePresetsQuery();
</script>

<template>
  <div v-if="isLoading" class="flex">
    <v-icon class="mx-auto" name="fa-spinner" animation="spin" scale="2" />
  </div>
  <d-error-alert v-else-if="isError" :error="error">Failed to list presets.</d-error-alert>
  <div v-else-if="data" class="overflow-x-auto w-full">
    <table class="table table-compact table-zebra w-full">
      <thead>
        <tr>
          <th>
            <d-button class="btn-primary btn-sm" @click="() => setPage(PAGE_HOME)" aria-label="Home">
              <v-icon name="fa-home" />
            </d-button>
          </th>
          <th>URL</th>
          <th>New Title</th>
          <th>New URL</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="p in data" :key="p.url">
          <th class="w-0">
            <d-button class="btn-success btn-sm" aria-label="Edit">
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
  </div>
</template>

<style>
</style>
