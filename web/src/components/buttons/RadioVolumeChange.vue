<template>
  <button
    class="button"
    :class="{
      'is-loading': loading || !radioLoaded,
    }"
    :disabled="loading || !radioLoaded"
    @click="setVolume"
  >
    <span class="icon">
      <i class="fas" :class="icon"></i>
    </span>
  </button>
</template>
<script>
import { mapState, mapGetters } from "vuex";

export default {
  props: {
    icon: {
      type: String,
      default: "",
    },
    change: {
      type: Number,
      default: 0,
    },
  },
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
    setVolume() {
      if (this.loading) return;
      this.loading = true;
      this.$store
        .dispatch("setRadioVolume", this.radio.volume + this.change)
        .finally(() => {
          this.loading = false;
        });
    },
  },
};
</script>
<style></style>
