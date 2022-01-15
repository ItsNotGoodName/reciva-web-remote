<template>
  <div class="mx-auto columns is-multiline">
    <div
      v-if="radioLoaded"
      v-for="preset in radio.presets"
      :key="preset.number"
      class="column is-one-third"
    >
      <RadioPreset :preset="preset" :radio="radio" />
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";

import RadioPreset from "./buttons/RadioPreset.vue";

export default {
  computed: {
    ...mapGetters(["radioLoaded"]),
    ...mapState({
      radio: (state) => state.r.radio,
    }),
  },
  components: { RadioPreset },
};
</script>

<style></style>
