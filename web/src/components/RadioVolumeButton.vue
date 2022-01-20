<template>
  <b-button
    :title="'Volume ' + (this.radio.isMuted ? 'Muted' : this.radio.volume)"
    :loading="radioVolumeChanging != 0 || radioVolumeRefreshing || !radioReady"
    @click="refreshRadioVolume"
  >
    <b-icon v-if="radio.isMuted" icon="fa-volume-mute" />
    <span v-else>{{ radio.volume }}</span>
  </b-button>
</template>

<script>
import { mapState, mapActions, mapGetters } from "vuex";

import BButton from "./Bulma/BButton.vue";
import BIcon from "./Bulma/BIcon.vue";

export default {
  components: {
    BButton,
    BIcon,
  },
  computed: {
    ...mapState({
      radio: (state) => state.r.radio,
      radioVolumeRefreshing: (state) => state.r.radioVolumeRefreshing,
      radioVolumeChanging: (state) => state.r.radioVolumeChanging,
    }),
    ...mapGetters(["radioReady"]),
  },
  methods: {
    ...mapActions(["refreshRadioVolume"]),
  },
};
</script>

<style lang="scss" scoped>
button {
  width: 3rem;
}
</style>
