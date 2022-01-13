<template>
  <button
    title="State"
    class="button is-rounded is-info"
    :class="{ 'is-loading': !radioLoaded }"
  >
    <span class="icon">
      <i v-if="radio.state == 'Playing'" class="fas ml-1 fa fa-play" />
      <i v-else-if="radio.state == 'Connecting'" class="fas fa-sync fa-pulse" />
      <i v-else-if="radio.state == 'Stopped'" class="fas fa-pause" />
      <i v-else class="fas fa-question" />
    </span>
  </button>
</template>

<script>
import { mapState, mapGetters } from "vuex";

export default {
  computed: {
    ...mapGetters(["radioLoaded"]),
    ...mapState({
      radio: (state) => state.r.radio,
    }),
  },
};
</script>

<style></style>
