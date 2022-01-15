<template>
  <div class="buttons has-addons is-flex">
    <RadioVolumeChange icon="fa-volume-off" title="Volume Down" :change="-5" />
    <RadioVolume />
    <RadioVolumeChange icon="fa-volume-up" title="Volume Up" :change="5" />
    <RadioRefresh />
    <RadioPower class="is-flex-grow-1" />
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";

import RadioPower from "./buttons/RadioPower.vue";
import RadioVolume from "./buttons/RadioVolume.vue";
import RadioVolumeChange from "./buttons/RadioVolumeChange.vue";
import RadioRefresh from "./buttons/RadioRefresh.vue";

export default {
  components: {
    RadioPower,
    RadioVolume,
    RadioVolumeChange,
    RadioRefresh,
  },
  computed: {
    ...mapGetters(["radioSelected"]),
    ...mapState({
      radios: (state) => state.r.radios,
      radiosLoading: (state) => state.r.radiosLoading,
    }),
    radioUUID: {
      get() {
        return this.$store.state.r.radioUUID;
      },
      set(radioUUID) {
        this.$store.dispatch("setRadioUUID", radioUUID);
      },
    },
  },
};
</script>

<style></style>
