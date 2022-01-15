<template>
  <div :style="style">
    <slot />
  </div>
</template>

<script>
import { mapGetters } from "vuex";

export default {
  computed: {
    ...mapGetters(["radioSelected"]),
    style() {
      return {
        "margin-bottom": "6rem",
        "margin-top": this.radioSelected ? "5rem" : "1rem",
      };
    },
  },
};
</script>

<style></style>
