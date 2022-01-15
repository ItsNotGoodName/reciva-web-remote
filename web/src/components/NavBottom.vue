<template>
  <nav
    class="navbar is-fixed-bottom p-2"
    role="navigation"
    aria-label="bottom navigation"
  >
    <div class="columns is-gapless" style="width: 100%">
      <div class="column m-2">
        <RadioSelect />
      </div>
      <div v-if="radioSelected" class="column m-2">
        <RadioControl />
      </div>
    </div>
  </nav>
</template>

<script>
import { mapGetters } from "vuex";

import RadioSelect from "./RadioSelect.vue";
import RadioControl from "./RadioControl.vue";

export default {
  components: {
    RadioSelect,
    RadioControl,
  },
  computed: {
    ...mapGetters(["radioSelected"]),
  },
};
</script>

<style scoped>
nav {
  border-top: 2px solid #e5e5e5;
}
</style>
