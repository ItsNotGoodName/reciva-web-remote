<template>
  <div
    style="background: white"
    class="flex fixed top-0 left-0 px-4 py-3 w-full"
  >
    <Button type="button" icon="pi pi-bars" @click="toggle" />
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
    <radio-dropdown class="flex-grow-1" />
    <loading-button
      :icon="radio.isMuted ? 'pi pi-volume-off' : 'pi pi-volume-up'"
      :label="radio.volume"
      class="p-button-secondary p-button-text"
      :lClick="refreshRadioVolume"
    />
  </div>
</template>

<script>
import Toolbar from "primevue/toolbar";
import Button from "primevue/button";
import OverlayPanel from "primevue/overlaypanel";
import { mapActions, mapState } from "vuex";

import LoadingButton from "./LoadingButton.vue";
import RadioDropdown from "./RadioDropdown.vue";

export default {
  components: {
    Toolbar,
    Button,
    OverlayPanel,
    LoadingButton,
    RadioDropdown,
  },
  methods: {
    ...mapActions(["discoverRadios", "refreshRadio", "refreshRadioVolume"]),
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