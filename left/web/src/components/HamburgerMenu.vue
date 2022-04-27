<script setup lang="ts">
import { PropType } from 'vue';
import { PAGE_HOME, PAGE_EDIT } from "../constants"

import DDropdownButton from './DaisyUI/DDropdownButton.vue';

const { version, versionLoading, page, setPage } = defineProps({
  version: {
    type: String,
    default: "",
  },
  versionLoading: {
    type: Boolean,
    default: false,
  },
  page: {
    type: String,
    default: "",
  },
  setPage: {
    type: Function as PropType<(page: string) => void>,
    default: () => { },
  },
});
</script>

<template>
  <div class="dropdown dropdown-top">
    <d-dropdown-button :class="{ 'btn-success': page != PAGE_HOME }">
      <v-icon name="fa-bars" />
    </d-dropdown-button>
    <ul tabindex="0" class="menu menu-compact dropdown-content mb-2 p-2 shadow bg-base-200 rounded-box w-52">
      <li>
        <a :class="{ 'active': page == PAGE_HOME }" @click="() => setPage(PAGE_HOME)">
          <v-icon name="fa-home" />Homepage
        </a>
      </li>
      <li>
        <a :class="{ 'active': page == PAGE_EDIT }" @click="() => setPage(PAGE_EDIT)">
          <v-icon name="fa-edit" />Edit Presets
        </a>
      </li>
      <li>
        <a href="https://github.com/ItsNotGoodName/reciva-web-remote">
          <v-icon name="fa-github" />Source Code
        </a>
      </li>
      <li v-if="versionLoading">
        <a>
          <v-icon name="fa-spinner" animation="spin" />Version
        </a>
      </li>
      <li v-else-if="version">
        <a :href="'https://github.com/ItsNotGoodName/reciva-web-remote/releases/tag/' + version">
          <v-icon name="fa-tag" />
          {{ version }}
        </a>
      </li>
    </ul>
  </div>
</template>

<style>
</style>
