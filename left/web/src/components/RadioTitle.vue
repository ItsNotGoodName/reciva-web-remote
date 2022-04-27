<script setup lang="ts">
const { title, metadata, url, url_new, preset_number } = defineProps({
  title: {
    type: String,
    default: "",
  },
  metadata: {
    type: String,
    default: "",
  },
  url: {
    type: String,
    default: "",
  },
  url_new: {
    type: String,
    default: "",
  },
  preset_number: {
    type: Number,
    default: 0,
  },
});
</script>

<template>
  <div class="dropdown dropdown-end">
    <label tabindex="0" class="btn btn-primary w-full justify-start truncate">
      <div class="w-0">
        <span class="badge badge-info badge-lg mr-2 ">
          {{ preset_number }}
        </span>
        {{ title }}
      </div>
    </label>
    <div tabindex="0" class="dropdown-content card card-compact w-full p-2 mt-2 shadow bg-primary text-primary-content">
      <div class="card-body break-words">
        <h3 class="card-title">Stream Information</h3>
        <p>
          <span class="badge badge-info mr-2">Title</span>
          {{ title }}
        </p>
        <p>
          <span class="badge badge-info mr-2">Preset Number</span>
          {{ preset_number }}
        </p>
        <p>
          <span class="badge badge-info mr-2">Metadata</span>
          {{ metadata }}
        </p>
        <p>
          <span class="badge badge-info mr-2">URL</span>
          <a class="link link-hover" :href="url">{{ url }}</a>
        </p>
        <p>
          <span class="badge badge-info mr-2">New URL</span>
          <a class="link link-hover" :href="url_new">{{ url_new }}</a>
        </p>
      </div>
    </div>
  </div>
</template>

<style>
</style>
