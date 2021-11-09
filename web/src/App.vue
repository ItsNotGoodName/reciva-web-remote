<script>
import { mapActions, mapState } from "vuex";
import Toast from "primevue/toast";

import RadioPlayer from "./components/RadioPlayer.vue";

export default {
  components: {
    Toast,
    RadioPlayer,
  },
  mounted() {
    this.$store.dispatch("init");
  },
  methods: {
    ...mapActions(["clearMessage"]),
  },
  computed: {
    ...mapState({}),
  },
};
</script>

<template>
  <div>
    <RadioPlayer />
    <Toast position="bottom-right" />
  </div>
</template>

<style>
</style>
