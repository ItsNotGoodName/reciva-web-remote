<template>
  <div class="buttons has-addons is-flex">
    <RadioVolumeChange title="Volume Down" :change="-5" icon="fa-volume-down" />
    <RadioVolume />
    <RadioVolumeChange title="Volume Up" :change="5" icon="fa-volume-up" />
    <RadioRefresh />
    <RadioPower class="is-flex-grow-1" />
  </div>
</template>

<script>
import { mapState } from "vuex";

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
    ...mapState({
      radiosLoading: (state) => state.r.radiosLoading,
    }),
  },
};
</script>

<style lang="scss" scoped></style>
