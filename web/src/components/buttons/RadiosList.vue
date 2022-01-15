<template>
  <button
    title="List"
    class="button"
    :class="{ 'is-loading': radiosLoading }"
    :disabled="radiosLoading"
    @click="listRadios"
  >
    <span class="icon is-small">
      <i class="fas fa-arrow-down"></i>
    </span>
  </button>
</template>
<script>
import { mapState, mapActions } from "vuex";

export default {
  computed: {
    ...mapState({
      radiosLoading: (state) => state.r.radiosLoading,
    }),
  },
  methods: {
    ...mapActions(["listRadios"]),
  },
};
</script>
<style></style>
