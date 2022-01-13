<template>
  <nav
    class="navbar is-fixed-top is-flex p-2"
    role="navigation"
    aria-label="main navigation"
    v-if="radioSelected"
  >
    <RadioState class="mr-2" />
    <RadioTitle class="is-flex-grow-1" />
  </nav>
</template>

<script>
import { mapGetters } from "vuex";

import RadioState from "./buttons/RadioState.vue";
import RadioTitle from "./buttons/RadioTitle.vue";

export default {
  computed: {
    ...mapGetters({
      radioSelected: "radioSelected",
    }),
  },
  components: {
    RadioState,
    RadioTitle,
  },
};
</script>

<style scoped>
nav {
  border-bottom: 2px solid #e5e5e5;
}
</style>
