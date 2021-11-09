<template>
  <div class="flex">
    <loading-button
      :icon="radio.isMuted ? 'pi pi-volume-off' : 'pi pi-volume-up'"
      :label="'' + radio.volume"
      class="p-button-secondary p-button-text mr-4 w-5rem"
      :lClick="refreshRadioVolume"
    />
    <Slider
      v-model="volume"
      class="my-auto w-10rem"
      @slideend="setRadioVolume"
    />
  </div>
</template>

<script>
import Slider from "primevue/slider";
import { mapActions, mapState } from "vuex";

import LoadingButton from "./buttons/LoadingButton.vue";

export default {
  components: {
    Slider,
    LoadingButton,
  },
  data() {
    return {
      volume: 0,
    };
  },
  methods: {
    ...mapActions(["refreshRadioVolume"]),
    setRadioVolume() {
      this.$store.dispatch("setRadioVolume", this.volume);
    },
  },
  computed: {
    ...mapState(["radio"]),
  },
};
</script>

<style lang="scss" scoped>
</style>