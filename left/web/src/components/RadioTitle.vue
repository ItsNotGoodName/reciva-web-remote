<script setup lang="ts">
import DDropdownButton from './DaisyUI/DDropdownButton.vue';

defineProps({
  state: {
    type: Object as () => State,
    required: true,
  },
  loading: {
    type: Boolean,
    default: false
  }
});
</script>

<template>
  <div class="dropdown dropdown-end">
    <d-dropdown-button :loading="loading" class="btn-primary w-full justify-start truncate">
      <div class="w-0">
        <span class="badge badge-info badge-lg rounded-md mr-2 ">
          {{ state.preset_number }}
        </span>
        {{ state.title_new ? state.title_new : state.title }}
      </div>
    </d-dropdown-button>
    <div tabindex="0" class="dropdown-content card card-compact w-full p-2 mt-2 shadow bg-primary text-primary-content">
      <div class="card-body break-words">
        <h3 class="card-title">Stream Information</h3>
        <p>
          <span class="badge badge-info mr-2">Metadata</span>
          {{ state.metadata }}
        </p>
        <p>
          <span class="badge badge-info mr-2">Title</span>
          {{ state.title }}
        </p>
        <p>
          <span class="badge badge-info mr-2">New Title</span>
          {{ state.title_new }}
        </p>
        <p>
          <span class="badge badge-info mr-2">Preset Number</span>
          {{ state.preset_number }}
        </p>
        <p>
          <span class="badge badge-info mr-2">URL</span>
          <a class="link link-hover" :href="state.url">{{ state.url }}</a>
        </p>
        <p>
          <span class="badge badge-info mr-2">New URL</span>
          <a class="link link-hover" :href="state.url_new">{{ state.url_new }}</a>
        </p>
      </div>
    </div>
  </div>
</template>

<style>
</style>
