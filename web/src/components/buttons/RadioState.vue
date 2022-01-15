<template>
  <button
    class="button is-rounded is-info"
    :class="{ 'is-loading': !radioLoaded }"
    style="width: 0%"
    :title="state.title"
  >
    <span class="icon">
      <i :class="state.iconClass" />
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
    state() {
      if (this.radio.state == "Playing") {
        return { iconClass: "fas ml-1 fa fa-play", title: "Playing" };
      } else if (this.radio.state == "Connecting") {
        return { iconClass: "fas fa-sync fa-pulse", title: "Connecting" };
      } else if (this.radio.state == "Stopped") {
        return { iconClass: "fas fa-pause", title: "Stopped" };
      } else {
        return { iconClass: "fas fa-question", title: "Unknown" };
      }
    },
  },
};
</script>

<style></style>
