<template>
  <div v-if="radioReady" class="columns is-multiline m-0">
    <div v-for="preset in radio.presets" :key="preset.number" class="column is-one-third">
      <radio-preset-button class="is-fullwidth" :preset="preset" :radio="radio" />
    </div>
  </div>
  <progress v-else-if="radioSelected" class="progress is-info" max="100" />
</template>

<script>
import { mapState, mapGetters } from "vuex";

import RadioPresetButton from "../components/RadioPresetButton.vue";

export default {
  computed: {
    ...mapGetters(["radioReady", "radioSelected"]),
    ...mapState({
      radio: (state) => state.r.radio,
    }),
  },
  components: { RadioPresetButton },
};
</script>

<style lang="scss" scoped></style>
