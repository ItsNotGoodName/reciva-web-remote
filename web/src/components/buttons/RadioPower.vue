<template>
  <button
    title="Power"
    class="button"
    :class="{
      'is-loading': loading || !radioLoaded,
      'is-success': radio.power,
      'is-danger': !radio.power,
    }"
    :disabled="loading || !radioLoaded"
    @click="togglePower"
  >
    <span class="icon is-small mr-1">
      <i class="fas fa-power-off"></i>
    </span>
    {{ radio.power ? "ON" : "OFF" }}
  </button>
</template>
<script>
import { mapGetters, mapState } from "vuex";

export default {
  data() {
    return { loading: false };
  },
  computed: {
    ...mapState({
      radio: (state) => state.r.radio,
    }),
    ...mapGetters(["radioLoaded"]),
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
