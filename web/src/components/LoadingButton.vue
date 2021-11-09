<template>
  <Button :loading="loading" v-bind="$props" @click="handleClick" />
</template>

<script>
import Button from "primevue/button";
export default {
  components: {
    Button,
  },
  props: {
    lClick: {
      type: Function,
      default: () => Promise.resolve(),
    },
  },
  data() {
    return {
      loading: false,
    };
  },
  methods: {
    handleClick() {
      if (this.loading) {
        return;
      }
      this.loading = true;
      this.lClick()
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

<style scoped>
</style>
