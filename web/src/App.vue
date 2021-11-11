<script>
import { mapGetters, mapState } from "vuex";

import Message from "primevue/message";
import Toast from "primevue/toast";

import PresetEditor from "./components/preset/PresetEditor.vue";
import PresetPlayer from "./components/radio/PresetList.vue";
import RadioPlayer from "./components/radio/RadioPlayer.vue";

export default {
  components: {
    Message,
    PresetEditor,
    PresetPlayer,
    RadioPlayer,
    Toast,
  },
  created() {
    this.$store.dispatch("loadRadios");
  },
  computed: {
    ...mapGetters(["radioReady"]),
    ...mapState(["page", "message"]),
  },
};
</script>

<template>
  <div>
    <radio-player class="mb-3" />
    <Message v-if="message" :severity="message.severity" :closable="false">{{
      message.content
    }}</Message>
    <preset-player v-if="page == 'player' && radioReady" />
    <preset-editor v-else-if="page == 'edit'" />
    <Toast
      :breakpoints="{ '480px': { width: '100%', right: '0', left: '0' } }"
      position="bottom-right"
    />
  </div>
</template>

<style>
</style>
