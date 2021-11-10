<template>
  <Button
    icon="pi pi-refresh"
    label="Refresh"
    :loading="loading"
    @click="refreshRadio"
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
    refreshRadio() {
      if (this.loading) {
        return;
      }
      this.loading = true;
      this.$store
        .dispatch("refreshRadio")
        .then(() => {
          this.loading = false;
          this.$toast.add({
            severity: "success",
            summary: "refreshed radio",
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