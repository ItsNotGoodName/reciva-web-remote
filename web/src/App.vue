<template>
  <div>
    <NavTop />
    <NavBottom />
    <div style="margin-bottom: 6rem; margin-top: 1rem">
      <RadioPreset />
    </div>
  </div>
</template>

<script setup>
import { onBeforeMount } from "vue";
import { useStore } from "vuex";

import NavTop from "./components/NavTop.vue";
import RadioPreset from "./components/RadioPreset.vue";
import NavBottom from "./components/NavBottom.vue";

const store = useStore();

onBeforeMount(() => {
  store.dispatch("initRadio");
});
</script>

<style></style>
