<template>
  <div class="buttons has-addons is-flex">
    <radio-volume-change-button title="Volume Down" :change="-5" icon="fa-volume-down" />
    <radio-volume-button />
    <radio-volume-change-button title="Volume Up" :change="5" icon="fa-volume-up" />
    <radio-power-button class="is-flex-grow-1" />
  </div>
</template>

<script>
import { mapState } from "vuex";

import BButton from "../Bulma/BButton.vue";
import BIcon from "../Bulma/BIcon.vue";

import RadioPowerButton from "../RadioPowerButton.vue";
import RadioVolumeButton from "../RadioVolumeButton.vue";
import RadioVolumeChangeButton from "../RadioVolumeChangeButton.vue";

export default {
  components: {
    BButton,
    BIcon,
    RadioPowerButton,
    RadioVolumeButton,
    RadioVolumeChangeButton,
  },
  computed: {
    ...mapState({
      radiosLoading: (state) => state.r.radiosLoading,
    }),
  },
};
</script>

<style lang="scss" scoped></style>
