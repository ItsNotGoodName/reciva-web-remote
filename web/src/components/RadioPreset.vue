<template>
  <div v-if="radioLoaded" class="mx-auto columns is-multiline">
    <div
      v-for="preset in radio.presets"
      :key="preset.number"
      class="column is-one-third"
    >
      <RadioPreset :preset="preset" :radio="radio" />
    </div>
  </div>
  <progress v-else-if="radioSelected" class="progress is-info" max="100" />
</template>

<script>
import { mapState, mapGetters } from "vuex";

import RadioPreset from "./buttons/RadioPreset.vue";

export default {
  computed: {
    ...mapGetters(["radioLoaded", "radioSelected"]),
    ...mapState({
      radio: (state) => state.r.radio,
    }),
  },
  components: { RadioPreset },
};
</script>

<style></style>
