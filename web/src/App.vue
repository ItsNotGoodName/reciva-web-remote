<script>
import { mapActions, mapGetters, mapState } from "vuex";
import Toast from "primevue/toast";

import RadioPlayer from "./components/RadioPlayer.vue";
import PresetPlayer from "./components/radio/PresetList.vue";

export default {
  components: {
    Toast,
    RadioPlayer,
    PresetPlayer,
  },
  mounted() {
    this.$store.dispatch("init");
  },
  methods: {
    ...mapActions(["clearMessage"]),
  },
  computed: {
    ...mapState({}),
    ...mapGetters(["radioReady"]),
  },
};
</script>

<template>
  <div>
    <radio-player class="mb-3" />
    <preset-player v-if="radioReady" />
    <Toast position="bottom-right" />
  </div>
</template>

<style>
</style>
