<template>
  <Button
    icon="pi pi-refresh"
    label="Discover"
    :loading="loading"
    @click="discoverRadios"
    class="p-button-success"
  />
</template>

<script>
import Button from "primevue/button";

export default {
  components: {
    Button,
  },
  data() {
    return {
      loading: false,
    };
  },
  methods: {
    discoverRadios() {
      if (this.loading) {
        return;
      }
      this.loading = true;
      this.$store
        .dispatch("discoverRadios")
        .then(() => {
          this.loading = false;
          this.$toast.add({
            severity: "success",
            summary:
              "Discovered " + this.$store.state.radios.length + " radios",
            life: 3000,
          });
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