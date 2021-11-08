<template >
  <div
    v-if="radioUUID && radioConnected"
    class="gap-2 grid lg:grid-cols-4 md:grid-cols-3 sm:grid-cols-2"
  >
    <div v-for="preset in radio.presets" :key="preset.number">
      <loading-button
        v-if="preset.number == radio.preset"
        class="
          w-full
          h-10
          rounded
          btn btn-primary
          border-blue-500
          hover:border-blue-600
          border-2
        "
        :on-click="() => playRadioPreset(preset.number)"
        :title="preset.name"
      >
        <div class="flex">
          <div class="rounded-full bg-white px-2 text-black mr-2">
            {{ preset.number }}
          </div>
          <div class="flex-grow truncate">
            {{ preset.name }}
          </div>
        </div>
      </loading-button>
      <loading-button
        v-else
        class="w-full h-10 rounded btn btn-white border-2"
        className="border-blue-600"
        :on-click="() => playRadioPreset(preset.number)"
        :title="preset.name"
      >
        <div class="flex">
          <div class="rounded-full bg-black text-white px-2 mr-2">
            {{ preset.number }}
          </div>
          <div class="flex-grow truncate">
            {{ preset.name }}
          </div>
        </div>
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