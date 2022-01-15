<template>
  <button
    class="button is-rounded"
    :class="[{ 'is-loading': !radioLoaded }, state.buttonClass]"
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
        return {
          buttonClass: "is-success",
          iconClass: "fas ml-1 fa fa-play",
          title: "Playing",
        };
      } else if (this.radio.state == "Connecting") {
        return {
          buttonClass: "is-warning",
          iconClass: "fas fa-sync fa-pulse",
          title: "Connecting",
        };
      } else if (this.radio.state == "Stopped") {
        return {
          buttonClass: "is-danger",
          iconClass: "fas fa-pause",
          title: "Stopped",
        };
      } else {
        return {
          buttonClass: "is-info",
          iconClass: "fas fa-question",
          title: "Unknown",
        };
      }
    },
  },
};
</script>

<style></style>
