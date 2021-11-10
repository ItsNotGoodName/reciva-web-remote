<template>
  <div class="grid">
    <div
      class="col-12 md:col-6 lg:col-3"
      v-for="preset in radio.presets"
      :key="preset.number"
    >
      <Button
        v-if="preset.number == radio.preset"
        :loading="loading"
        :label="preset.name"
        @click="playRadioPreset(preset.number)"
        :badge="'' + preset.number"
        class="w-full"
      />
      <Button
        v-else
        :loading="loading"
        :label="preset.name"
        @click="playRadioPreset(preset.number)"
        :badge="'' + preset.number"
        badgeClass="p-badge-info"
        class="w-full p-button-secondary p-button-outlined"
      />
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";

import Button from "primevue/button";

export default {
  components: {
    Button,
  },
  data: () => ({
    loading: false,
  }),
  computed: {
    ...mapState(["radio"]),
  },
  methods: {
    playRadioPreset(preset) {
      if (this.loading) return;

      this.loading = true;
      this.$store
        .dispatch("playRadioPreset", preset)
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
};
</script>

<style lang="scss" scoped>
</style>