<template>
  <b-button
    :class="state.buttonClass"
    :loading="!radioLoaded && radioSelected"
    :title="state.title"
    class="is-rounded"
    style="width: 0%"
  >
    <b-icon :icon="state.iconClass" />
  </b-button>
</template>

<script>
import { mapState, mapGetters } from "vuex";

import BButton from "../Bulma/BButton.vue";
import BIcon from "../Bulma/BIcon.vue";

export default {
  components: {
    BButton,
    BIcon,
  },
  computed: {
    ...mapGetters(["radioLoaded", "radioSelected"]),
    ...mapState({
      radio: (state) => state.r.radio,
    }),
    state() {
      if (this.radio.state == "Playing") {
        return {
          buttonClass: "is-success",
          iconClass: "ml-1 fa-play",
          title: "Playing",
        };
      } else if (this.radio.state == "Connecting") {
        return {
          buttonClass: "is-warning",
          iconClass: "fa-sync fa-pulse",
          title: "Connecting",
        };
      } else if (this.radio.state == "Stopped") {
        return {
          buttonClass: "is-danger",
          iconClass: "fa-pause",
          title: "Stopped",
        };
      } else {
        return {
          buttonClass: "is-info",
          iconClass: "fa-question",
          title: "Unknown",
        };
      }
    },
  },
};
</script>

<style lang="scss" scoped></style>
