<template>
  <b-button @click="setRadioVolume(this.radio.volume + this.change)">
    <b-icon :icon="icon" />
  </b-button>
</template>

<script>
import { mapActions, mapState } from "vuex";

import BButton from "./Bulma/BButton.vue";
import BIcon from "./Bulma/BIcon.vue";

export default {
  components: {
    BButton,
    BIcon,
  },
  props: {
    icon: {
      type: String,
      default: "",
    },
    change: {
      type: Number,
      default: 0,
    },
  },
  computed: {
    ...mapState({
      radio: (state) => state.r.radio,
    }),
  },
  methods: {
    ...mapActions(["setRadioVolume"]),
  },
};
</script>

<style lang="scss" scoped></style>
