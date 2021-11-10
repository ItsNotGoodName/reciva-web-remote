<template>
  <Button
    v-if="preset.number == radio.preset"
    :loading="loading"
    :label="preset.name"
    @click="playRadioPreset(preset.number)"
    :badge="'' + preset.number"
    class="w-full h-full"
  />
  <Button
    v-else
    :loading="loading"
    :label="preset.name"
    @click="playRadioPreset(preset.number)"
    :badge="'' + preset.number"
    badgeClass="p-badge-info"
    class="w-full h-full p-button-secondary p-button-outlined"
  />
</template>

<script>
import Button from "primevue/button";

export default {
  components: {
    Button,
  },
  props: {
    preset: {
      type: Object,
      required: true,
    },
    radio: {
      type: Object,
      required: true,
    },
  },
  data() {
    return { loading: false };
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