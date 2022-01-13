<template>
  <button
    title="Volume"
    class="button"
    :class="{
      'is-loading': radioVolumeRefreshing || !radioLoaded,
    }"
    :disabled="radioVolumeRefreshing || !radioLoaded"
    @click="refreshRadioVolume"
  >
    <span v-if="radio.isMuted" class="icon">
      <i class="fas fa-volume-off"></i>
    </span>
    <span v-else>
      {{ radio.volume }}
    </span>
  </button>
</template>
<script>
import { mapState, mapActions, mapGetters } from "vuex";

export default {
  data() {
    return { loading: false };
  },
  computed: {
    ...mapState({
      radio: (state) => state.r.radio,
      radioVolumeRefreshing: (state) => state.r.radioVolumeRefreshing,
    }),
    ...mapGetters(["radioLoaded"]),
  },
  methods: {
    ...mapActions(["refreshRadioVolume"]),
  },
};
</script>
<style></style>
