<template>
  <nav class="navbar is-fixed-bottom py-2" role="navigation" aria-label="bottom navigation">
    <div class="is-flex is-flex-wrap-wrap full-width">
      <div class="is-flex-grow-2 p-2">
        <radio-select />
      </div>
      <div v-if="radioSelected" class="is-flex-grow-1 p-2">
        <radio-control />
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

<style lang="scss" scoped>
nav {
  border-top: 1px solid #ddd;
}
.full-width {
  width: 100%;
}
</style>
