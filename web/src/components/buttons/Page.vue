<template>
  <button
    title="Toggle Page"
    class="button"
    :class="buttonClass"
    @click="togglePage"
  >
    <span class="icon is-small">
      <i class="fas fa-bars"></i>
    </span>
  </button>
</template>
<script>
import { mapActions, mapState } from "vuex";

export default {
  data() {
    return { loading: false };
  },
  computed: {
    ...mapState(["page"]),
    buttonClass() {
      if (this.page == "play") {
        return "is-info";
      } else {
        return "";
      }
    },
  },
  methods: {
    ...mapActions(["togglePage"]),
  },
};
</script>
<style></style>
