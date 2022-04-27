<script setup lang="ts">
const { name, model_name, model_number } = defineProps({
  name: {
    type: String,
    default: "",
  },
  model_name: {
    type: String,
    default: "",
  },
  model_number: {
    type: String,
    default: "",
  },
});
</script>

<template>
  <div class="dropdown dropdown-top dropdown-end">
    <label tabindex="0" class="btn btn-primary w-full justify-start truncate">
      <v-icon name="md-radio" />
    </label>
    <div tabindex="0" class="dropdown-content card card-compact w-80 p-2 mb-2 shadow bg-primary text-primary-content">
      <div class="card-body break-words">
        <h3 class="card-title">Radio Information</h3>
        <p>
          <span class="badge badge-info mr-2">Name</span>
          {{ name }}
        </p>
        <p>
          <span class="badge badge-info mr-2">Model Name</span>
          {{ model_name }}
        </p>
        <p>
          <span class="badge badge-info mr-2">Model Number</span>
          {{ model_number }}
        </p>
      </div>
    </div>
  </div>
</template>

<style>
</style>
