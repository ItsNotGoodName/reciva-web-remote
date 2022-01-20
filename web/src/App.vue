<template>
  <NavTop />
  <div class="container nav-margin px-2">
    <Edit v-if="store.state.page == 'edit'" />
    <Home v-else />
  </div>
  <NavBottom />
</template>

<script setup>
import { useStore } from "vuex";
import { onMounted } from "vue";

import NavTop from "./components/NavTop/index.vue";
import NavBottom from "./components/NavBottom/index.vue";

import Home from "./views/Home.vue"
import Edit from "./views/Edit.vue"

const store = useStore();

onMounted(() => {
  store.dispatch("initRadio");
});
</script>

<style lang="scss" scoped>
.nav-margin {
  margin-top: 1rem;
  margin-bottom: 6rem;
}
</style>
