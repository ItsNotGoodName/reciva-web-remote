<template >
  <div
    v-if="radioUUID && radioConnected"
    class="gap-1 grid lg:grid-cols-4 md:grid-cols-3 sm:grid-cols-2"
  >
    <div v-for="preset in radio.presets" :key="preset.number">
      <loading-button
        v-if="preset.number == radio.preset"
        class="w-full h-10 rounded btn btn-primary"
        :className="{ 'border-blue-600': preset.number != radio.preset }"
        :on-click="() => playRadioPreset(preset.number)"
        >{{ preset.name }}</loading-button
      >
      <loading-button
        v-else
        class="w-full h-10 rounded btn btn-white"
        :className="{ 'border-blue-600': preset.number != radio.preset }"
        :on-click="() => playRadioPreset(preset.number)"
      >
        {{ preset.name }}
      </loading-button>
    </div>
  </div>
</template>

<script>
import { mapState, mapActions } from "vuex";

import LoadingButton from "./LoadingButton.vue";

export default {
  components: {
    LoadingButton,
  },
  computed: {
    ...mapState(["radio", "radioConnected", "radioConnecting", "radioUUID"]),
  },
  methods: {
    ...mapActions(["playRadioPreset"]),
  },
};
</script>

<style scoped>
</style>