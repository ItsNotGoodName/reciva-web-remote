<template>
  <Button
    icon="pi pi-refresh"
    label="Discover"
    :loading="discovering"
    title="Discovers radios on the local network."
    @click="discoverRadios"
    class="p-button-success"
  />
</template>

<script>
import Button from "primevue/button";
import { mapState } from "vuex";

export default {
  components: {
    Button,
  },
  computed: {
    ...mapState(["discovering"]),
  },
  methods: {
    discoverRadios() {
      if (this.discovering) {
        return;
      }
      this.$store
        .dispatch("discoverRadios")
        .then(() => {
          this.$toast.add({
            severity: "success",
            summary:
              "discovered " + this.$store.state.radios.length + " radios",
            life: 3000,
          });
        })
        .catch((err) => {
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
