<template>
  <nav
    class="navbar is-fixed-top"
    role="navigation"
    aria-label="top navigation"
  >
    <RadioState />
    <RadioTitle />
  </nav>
</template>

<script>
import { mapGetters } from "vuex";

import RadioState from "./RadioState.vue";
import RadioTitle from "./RadioTitle.vue";

export default {
  components: {
    RadioState,
    RadioTitle,
  },
  computed: {
    ...mapGetters({
      radioSelected: "radioSelected",
    }),
  },
};
</script>

<style></style>
