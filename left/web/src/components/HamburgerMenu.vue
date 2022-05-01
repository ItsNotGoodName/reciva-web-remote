<script setup lang="ts">
import { computed } from "vue"

import { PAGE_HOME, PAGE_EDIT, GITHUB_URL } from "../constants"
import { useBuildQuery } from "../hooks";

import DDropdownButton from './DaisyUI/DDropdownButton.vue';

const emit = defineEmits<{ (e: 'update:page', page: string): void }>()

defineProps({
  page: {
    type: String,
    required: true,
  },
});

const { data, isLoading } = useBuildQuery();
const versionUrl = computed(() => {
  if (data.value && data.value.version != 'dev') {
    return `${GITHUB_URL}/releases/tag/v${data.value.version}`;
  }
  return '#'
})
</script>

<template>
  <div class="dropdown dropdown-top">
    <d-dropdown-button :class="{ 'btn-success': page != PAGE_HOME }">
      <v-icon name="fa-bars" />
    </d-dropdown-button>
    <ul tabindex="0" class="menu menu-compact dropdown-content mb-2 p-2 shadow bg-base-200 rounded-box w-52">
      <li>
        <a :class="{ 'active': page == PAGE_HOME }" @click="emit('update:page', PAGE_HOME)">
          <v-icon name="fa-home" />Home Page
        </a>
      </li>
      <li>
        <a :class="{ 'active': page == PAGE_EDIT }" @click="emit('update:page', PAGE_EDIT)">
          <v-icon name="fa-edit" />Edit Presets
        </a>
      </li>
      <li>
        <a :href="GITHUB_URL">
          <v-icon name="fa-github" />Source Code
        </a>
      </li>
      <li v-if="isLoading">
        <a>
          <v-icon name="fa-spinner" animation="spin" />Version
        </a>
      </li>
      <li v-else-if="data">
        <a :href="versionUrl">
          <v-icon name="fa-tag" />
          {{ data.version }}
        </a>
      </li>
    </ul>
  </div>
</template>

<style>
</style>
