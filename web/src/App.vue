<script>
import { mapGetters, mapState } from "vuex";

import Toast from "primevue/toast";

import PresetEditor from "./components/preset/PresetEditor.vue";
import PresetPlayer from "./components/radio/PresetList.vue";
import RadioPlayer from "./components/radio/RadioPlayer.vue";

export default {
  components: {
    Toast,
    PresetEditor,
    PresetPlayer,
    RadioPlayer,
  },
  mounted() {
    this.$store.dispatch("init");
  },
  computed: {
    ...mapGetters(["radioReady"]),
    ...mapState(["page"]),
  },
};
</script>

<template>
  <div>
    <radio-player class="mb-3" />
    <preset-player v-if="page == 'player' && radioReady" />
    <preset-editor v-else-if="page == 'edit'" />
    <Toast position="bottom-right" />
  </div>
</template>

<style>
</style>
