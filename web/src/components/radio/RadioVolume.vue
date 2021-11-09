<template>
  <div class="flex">
    <Button
      :icon="radio.isMuted ? 'pi pi-volume-off' : 'pi pi-volume-up'"
      :label="'' + radio.volume"
      class="p-button-secondary p-button-text mr-2 w-5rem"
      :loading="loading"
      @click="refreshRadioVolume"
    />
    <Slider
      v-model="radioVolume"
      class="my-auto w-10rem flex-grow"
      @slideend="setRadioVolume"
    />
  </div>
</template>

<script>
import Slider from "primevue/slider";
import Button from "primevue/button";
import { mapState } from "vuex";

export default {
  components: {
    Slider,
    Button,
  },
  data() {
    return {
      loading: false,
    };
  },
  methods: {
    setRadioVolume() {
      this.$store.dispatch("setRadioVolume", this.radio.volume);
    },
    refreshRadioVolume() {
      if (this.loading) {
        return;
      }
      this.loading = true;
      this.$store
        .dispatch("refreshRadioVolume")
        .then(() => {
          this.loading = false;
        })
        .catch((err) => {
          this.loading = false;
          this.$toast.add({
            severity: "error",
            summary: err || "Error",
            life: 3000,
          });
        });
    },
  },
  computed: {
    radioVolume: {
      get() {
        return this.$store.state.radio.volume;
      },
      set(value) {
        return this.$store.commit("SET_RADIO_VOLUME", value);
      },
    },
    ...mapState(["radio"]),
  },
};
</script>

<style lang="scss" scoped>
</style>