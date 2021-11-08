<script>
import { mapState } from "vuex";

import Alerts from "./components/Alerts.vue";
import Player from "./components/Player.vue";
import PresetPlayer from "./components/PresetPlayer.vue";
import PresetEditor from "./components/PresetEditor.vue";

export default {
  name: "App",
  components: {
    Alerts,
    Player,
    PresetPlayer,
    PresetEditor,
  },
  computed: {
    ...mapState(["edit"]),
  },
  mounted() {
    this.$store
      .dispatch("init")
      .catch((err) =>
        this.$store.dispatch("addNotification", { type: "error", message: err })
      );
  },
};
</script>

<template >
  <div class="container mx-auto px-1 h-screen">
    <Player />
    <PresetEditor v-if="edit" />
    <PresetPlayer v-else />
    <Alerts />
  </div>
</template>

<style>
</style>
