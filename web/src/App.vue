<template>
  <div>
    <NavTop />
    <Notifications />
    <MainBody />
    <NavBottom />
  </div>
</template>

<script setup>
import { onBeforeMount } from "vue";
import { useStore } from "vuex";

import NavTop from "./components/NavTop/index.vue";
import MainBody from "./components/MainBody.vue";
import NavBottom from "./components/NavBottom/index.vue";
import Notifications from "./components/Notifications.vue";

const store = useStore();

onBeforeMount(() => {
  store.dispatch("initRadio");
});
</script>

<style lang="scss" scoped></style>
