<script setup lang="ts">
import DDropdownButton from './DaisyUI/DDropdownButton.vue';

defineProps({
  radio: {
    type: Object as () => Radio,
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
          {{ radio.preset_number }}
        </span>
        {{ radio.title_new ? radio.title_new : radio.title }}
      </div>
    </d-dropdown-button>
    <div tabindex="0" class="dropdown-content card card-compact w-full p-2 mt-1 shadow bg-primary text-primary-content">
      <div class="card-body break-words">
        <h3 class="card-title">Stream Information</h3>
        <p>
          <span class="badge badge-info mr-2">Metadata</span>
          {{ radio.metadata }}
        </p>
        <p>
          <span class="badge badge-info mr-2">Title</span>
          {{ radio.title }}
        </p>
        <p>
          <span class="badge badge-info mr-2">New Title</span>
          {{ radio.title_new }}
        </p>
        <p>
          <span class="badge badge-info mr-2">Preset Number</span>
          {{ radio.preset_number }}
        </p>
        <p>
          <span class="badge badge-info mr-2">URL</span>
          <a class="link link-hover" :href="radio.url">{{ radio.url }}</a>
        </p>
        <p>
          <span class="badge badge-info mr-2">New URL</span>
          <a class="link link-hover" :href="radio.url_new">{{ radio.url_new }}</a>
        </p>
      </div>
    </div>
  </div>
</template>

<style>
</style>
