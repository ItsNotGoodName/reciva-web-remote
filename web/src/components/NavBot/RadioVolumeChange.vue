<template>
  <b-button :loading="loading" @click="setVolume">
    <b-icon :icon="icon" />
  </b-button>
</template>

<script>
import { mapState } from "vuex";

import BButton from "../Bulma/BButton.vue";
import BIcon from "../Bulma/BIcon.vue";

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
  data() {
    return { loading: false };
  },
  computed: {
    ...mapState({
      radio: (state) => state.r.radio,
    }),
  },
  methods: {
    setVolume() {
      if (this.loading) return;
      this.loading = true;
      this.$store
        .dispatch("setRadioVolume", this.radio.volume + this.change)
        .finally(() => {
          this.loading = false;
        });
    },
  },
};
</script>

<style></style>
