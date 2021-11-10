<template>
  <div
    class="
      sticky
      z-3
      bg-white
      top-0
      left-0
      pb-2
      w-full
      border-bottom-2 border-500
    "
  >
    <div class="flex flex-wrap">
      <Button class="m-1" type="button" icon="pi pi-bars" @click="toggle" />
      <radio-dropdown class="m-1 flex-grow-1" />
      <overlay-panel ref="op">
        <div class="flex flex-column">
          <discover-button class="m-1" />
          <refresh-button class="m-1" v-if="radioReady" />
          <page-button class="m-1" />
        </div>
      </overlay-panel>
      <div v-if="radioReady" class="m-1 flex flex-1">
        <radio-volume class="mr-3 flex-grow-1" />
        <power-button />
      </div>
      <playing-text class="w-full m-1" />
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";

import Button from "primevue/button";
import OverlayPanel from "primevue/overlaypanel";
import Toolbar from "primevue/toolbar";

import DiscoverButton from "./DiscoverButton.vue";
import RadioDropdown from "./RadioDropdown.vue";
import RadioVolume from "./RadioVolume.vue";
import RefreshButton from "./RefreshButton.vue";
import PowerButton from "./PowerButton.vue";
import PlayingText from "./PlayingText.vue";
import PageButton from "../PageButton.vue";

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
    PageButton,
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