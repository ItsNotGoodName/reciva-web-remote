<template>
  <div
    style="background: white"
    class="fixed top-0 left-0 px-2 py-2 w-full border-bottom-2 border-300"
  >
    <div class="flex flex-wrap">
      <Button class="m-1" type="button" icon="pi pi-bars" @click="toggle" />
      <radio-dropdown class="m-1 flex-auto" />
      <overlay-panel ref="op">
        <div class="p-buttonset">
          <discover-button />
          <refresh-button />
        </div>
      </overlay-panel>
      <div v-if="radioReady" class="m-1 flex flex-1">
        <radio-volume class="mr-2" />
        <power-button class="flex-grow-1" />
      </div>
      <playing-text class="flex-auto m-1" />
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";

import Button from "primevue/button";
import OverlayPanel from "primevue/overlaypanel";
import Toolbar from "primevue/toolbar";

import DiscoverButton from "./radio/DiscoverButton.vue";
import RadioDropdown from "./radio/RadioDropdown.vue";
import RadioVolume from "./radio/RadioVolume.vue";
import RefreshButton from "./radio/RefreshButton.vue";
import PowerButton from "./radio/PowerButton.vue";
import PlayingText from "./radio/PlayingText.vue";

export default {
  components: {
    Button,
    DiscoverButton,
    OverlayPanel,
    RadioDropdown,
    RadioVolume,
    RefreshButton,
    Toolbar,
    PowerButton,
    PlayingText,
  },
  methods: {
    toggle(event) {
      this.$refs.op.toggle(event);
    },
  },
  computed: {
    ...mapState(["radio"]),
    ...mapGetters(["radioReady"]),
  },
};
</script>

<style scoped>
</style>