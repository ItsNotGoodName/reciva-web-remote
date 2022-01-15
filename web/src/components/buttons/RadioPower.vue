<template>
  <button
    class="button"
    :title="power.title"
    :class="[{ 'is-loading': loading || !radioLoaded }, power.class]"
    :disabled="loading || !radioLoaded"
    @click="togglePower"
  >
    <span class="icon is-small mr-1">
      <i class="fas fa-power-off"></i>
    </span>
    {{ power.text }}
  </button>
</template>
<script>
import { mapGetters, mapState } from "vuex";

export default {
  data() {
    return { loading: false };
  },
  computed: {
    ...mapGetters(["radioLoaded"]),
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
<style></style>
