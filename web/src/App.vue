<template></template>

<script setup>
import { onBeforeMount } from "vue";
import { useStore } from "vuex";

const store = useStore();

onBeforeMount(() => {
  store.dispatch("listRadios");
});
</script>

<style></style>
