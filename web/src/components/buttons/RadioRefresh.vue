<template>
  <button
    title="Refresh"
    class="button"
    :class="{ 'is-loading': radioRefreshing, 'is-danger': !radioWSConnected }"
    :disabled="radioRefreshing"
    @click="refreshRadio"
  >
    <span class="icon is-small">
      <i class="fas fa-redo"></i>
    </span>
  </button>
</template>
<script>
import { mapState, mapActions, mapGetters } from "vuex";

export default {
  computed: {
    ...mapState({
      radioRefreshing: (state) => state.r.radioRefreshing,
      radioWSConnected: (state) => state.r.radioWSConnected,
    }),
    ...mapGetters(["radioLoaded"]),
  },
  methods: {
    ...mapActions(["refreshRadio"]),
  },
};
</script>
<style></style>
