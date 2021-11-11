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
      <Button class="m-1" icon="pi pi-bars" @click="toggle" />
      <OverlayPanel ref="op">
        <div class="flex flex-column">
          <discover-button class="m-1" />
          <refresh-button class="m-1" v-if="radioReady" />
          <page-button :toggle="toggle" class="m-1" />
        </div>
      </OverlayPanel>
      <radio-dropdown class="m-1 flex-auto" />
      <Button
        class="m-1"
        :loading="loading"
        icon="pi pi-refresh"
        @click="loadRadios"
      />
      <template v-if="radioReady" class="m-1 flex flex-1">
        <radio-volume class="m-1 mr-3 flex-grow-1" />
        <power-button class="m-1 flex-auto" />
      </template>
      <playing-text class="w-full m-1" />
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters, mapActions } from "vuex";

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
    ...mapActions(["loadRadios"]),
    toggle(event) {
      this.$refs.op.toggle(event);
    },
  },
  computed: {
    ...mapState(["radio", "loading"]),
    ...mapGetters(["radioReady"]),
  },
};
</script>

<style scoped>
</style>