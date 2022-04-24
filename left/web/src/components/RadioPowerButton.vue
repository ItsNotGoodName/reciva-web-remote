<template>
  <b-button
    :title="power.title"
    :loading="loading || !radioReady"
    :class="power.class"
    @click="togglePower"
  >
    <b-icon icon="fa-power-off" />
    <span>{{ power.text }}</span>
  </b-button>
</template>

<script>
import { mapGetters, mapState } from "vuex";

import BButton from "./Bulma/BButton.vue";
import BIcon from "./Bulma/BIcon.vue";

export default {
  components: {
    BButton,
    BIcon,
  },
  data() {
    return { loading: false };
  },
  computed: {
    ...mapGetters(["radioReady"]),
    ...mapState({
      radio: (state) => state.r.radio,
    }),
    power() {
      if (this.radio.power) {
        return {
          text: "ON",
          title: "Powered ON",
          class: "is-success",
        };
      } else {
        return {
          text: "OFF",
          title: "Powered OFF",
          class: "is-danger",
        };
      }
    },
  },
  methods: {
    togglePower() {
      if (this.loading) return;
      this.loading = true;
      this.$store.dispatch("toggleRadioPower").finally(() => {
        this.loading = false;
      });
    },
  },
};
</script>

<style lang="scss" scoped>
button {
  min-width: 6rem;
}
</style>
