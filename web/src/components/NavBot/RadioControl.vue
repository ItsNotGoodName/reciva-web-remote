<template>
  <div v-if="radioSelected" class="buttons has-addons is-flex">
    <RadioVolumeChange title="Volume Down" :change="-5" icon="fa-volume-down" />
    <RadioVolume />
    <RadioVolumeChange title="Volume Up" :change="5" icon="fa-volume-up" />
    <RadioRefresh />
    <RadioPower class="is-flex-grow-1" />
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";

import BButton from "../Bulma/BButton.vue";
import BIcon from "../Bulma/BIcon.vue";
import RadioPower from "./RadioPower.vue";
import RadioRefresh from "./RadioRefresh.vue";
import RadioVolume from "./RadioVolume.vue";
import RadioVolumeChange from "./RadioVolumeChange.vue";

export default {
  components: {
    BButton,
    BIcon,
    RadioPower,
    RadioRefresh,
    RadioVolume,
    RadioVolumeChange,
  },
  computed: {
    ...mapGetters(["radioSelected"]),
    ...mapState({
      radiosLoading: (state) => state.r.radiosLoading,
    }),
  },
};
</script>

<style></style>
