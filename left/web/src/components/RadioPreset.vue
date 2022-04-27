<script setup lang="ts">
const { title, selected, number } = defineProps({
  selected: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: "",
  },
  number: {
    type: Number,
    default: 0,
  },
});
</script>

<template>
  <button class="btn flex w-full" :class="{ 'btn-primary': selected }">
    <span class="badge badge-info badge-lg mr-2">{{ number }} </span>
    <span class="flex-grow w-0 truncate ">{{ title }}</span>
  </button>
</template>

<style>
</style>
