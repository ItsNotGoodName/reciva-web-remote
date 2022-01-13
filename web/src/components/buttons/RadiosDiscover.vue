<template>
  <button
    title="Discover"
    class="button is-success"
    :class="{ 'is-loading': radiosDiscovering }"
    :disabled="radiosDiscovering"
    @click="discoverRadios"
  >
    <span class="icon is-small">
      <i class="fas fa-search"></i>
    </span>
  </button>
</template>
<script>
import { mapState, mapActions } from "vuex";

export default {
  computed: {
    ...mapState({
      radiosDiscovering: (state) => state.r.radiosDiscovering,
    }),
  },
  methods: {
    ...mapActions(["discoverRadios"]),
  },
};
</script>
<style></style>
