<template>
  <div class="grid">
    <div
      class="col-12 md:col-6 lg:col-3"
      v-for="preset in radio.presets"
      :key="preset.number"
    >
      <preset-button :radio="radio" :preset="preset" />
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";

import PresetButton from "./PresetButton.vue";

export default {
  components: {
    PresetButton,
  },
  computed: {
    ...mapState(["radio"]),
  },
};
</script>

<style lang="scss" scoped>
</style>