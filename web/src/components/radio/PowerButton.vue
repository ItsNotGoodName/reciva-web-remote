<template>
  <Button
    v-if="radio.power"
    label="ON"
    :loading="loading"
    @click="toggleRadioPower"
    class="p-button-success w-6rem"
    icon="pi pi-circle-on"
    title="Turn off radio."
  />
  <Button
    v-else
    label="OFF"
    :loading="loading"
    @click="toggleRadioPower"
    class="p-button-danger w-6rem"
    icon="pi pi-circle-off"
    title="Turn on radio."
  />
</template>

<script>
import Button from "primevue/button";
import { mapState } from "vuex";

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
    toggleRadioPower() {
      if (this.loading) {
        return;
      }
      this.loading = true;
      this.$store
        .dispatch("toggleRadioPower")
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
    ...mapState(["radio"]),
  },
};
</script>

<style lang="scss" scoped>
</style>