<template>
  <div>
    <NavTop />
    <NavMargin>
      <MainBody />
    </NavMargin>
    <NavBottom />
  </div>
</template>

<script setup>
import { onBeforeMount } from "vue";
import { useStore } from "vuex";

import NavTop from "./components/NavTop.vue";
import NavBottom from "./components/NavBottom.vue";
import NavMargin from "./components/NavMargin.vue";
import MainBody from "./components/MainBody.vue";

const store = useStore();

onBeforeMount(() => {
  store.dispatch("initRadio");
});
</script>

<style></style>
