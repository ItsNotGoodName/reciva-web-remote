<template>
  <div
    style="background: white"
    class="flex fixed top-0 left-0 px-2 py-2 w-full"
  >
    <Button type="button" icon="pi pi-bars" @click="toggle" class="m-2" />
    <overlay-panel ref="op">
      <div class="flex flex-column">
        <loading-button
          icon="pi pi-refresh"
          label="Discover"
          :lClick="discoverRadios"
          class="p-button-success mb-2"
        />
        <loading-button
          icon="pi pi-refresh"
          label="Refresh"
          :lClick="refreshRadio"
        />
      </div>
    </overlay-panel>
    <radio-dropdown class="flex-grow-1 m-2" />
    <radio-volume class="m-2" />
  </div>
</template>

<script>
import Toolbar from "primevue/toolbar";
import Button from "primevue/button";
import OverlayPanel from "primevue/overlaypanel";
import { mapActions, mapState } from "vuex";

import LoadingButton from "./buttons/LoadingButton.vue";
import RadioDropdown from "./RadioDropdown.vue";
import RadioVolume from "./RadioVolume.vue";

export default {
  components: {
    Toolbar,
    RadioVolume,
    Button,
    OverlayPanel,
    LoadingButton,
    RadioDropdown,
  },
  methods: {
    ...mapActions(["discoverRadios", "refreshRadio"]),
    toggle(event) {
      this.$refs.op.toggle(event);
    },
  },
  computed: {
    ...mapState(["radio"]),
  },
};
</script>

<style scoped>
</style>