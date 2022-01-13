<template>
  <div class="mx-auto columns is-multiline">
    <div
      class="column is-one-third"
      v-for="preset in radio.presets"
      :key="preset.number"
    >
      <RadioPreset class="is-fullwidth" v-if="radioLoaded" :preset="preset" />
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";

import RadioPreset from "./buttons/RadioPreset.vue";

export default {
  data() {
    return { loading: false };
  },
  computed: {
    ...mapGetters(["radioLoaded"]),
    ...mapState({
      radio: (state) => state.r.radio,
    }),
  },
  methods: {
    playRadioPreset(preset) {
      if (this.loading) return;
      this.loading = true;
      this.$store
        .dispatch("playRadioPreset", preset)
        .then(() => {
          this.loading = false;
        })
        .catch((err) => {
          this.loading = false;
          this.$toast.add({
            severity: "error",
            summary: err || "Error",
            life: 3000,
          });
        });
    },
  },
  components: { RadioPreset },
};
</script>

<style></style>
