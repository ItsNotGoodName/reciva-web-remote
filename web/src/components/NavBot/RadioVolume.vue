<template>
  <b-button
    :title="'Volume ' + (this.radio.isMuted ? 'Muted' : this.radio.volume)"
    :loading="radioVolumeRefreshing || !radioLoaded"
    @click="refreshRadioVolume"
  >
    <b-icon v-if="radio.isMuted" icon="fa-volume-mute" />
    <span v-else>
      {{ radio.volume }}
    </span>
  </b-button>
</template>

<script>
import { mapState, mapActions, mapGetters } from "vuex";

import BButton from "../Bulma/BButton.vue";
import BIcon from "../Bulma/BIcon.vue";

export default {
  components: {
    BButton,
    BIcon,
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
