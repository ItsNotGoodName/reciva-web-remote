<template>
  <div>
    <NavTop />
    <MainBody />
    <NavBot />
  </div>
</template>

<script setup>
import { onBeforeMount } from "vue";
import { useStore } from "vuex";

import NavTop from "./components/NavTop/index.vue";
import MainBody from "./components/MainBody.vue";
import NavBot from "./components/NavBot/index.vue";

const store = useStore();

onBeforeMount(() => {
  store.dispatch("initRadio");
});
</script>

<style></style>
