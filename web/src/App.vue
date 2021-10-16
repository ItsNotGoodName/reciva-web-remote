<script>
import SelectRadio from './components/SelectRadio.vue';
import RadioPanel from './components/RadioPanel.vue';
import RadioPresets from './components/RadioPresets.vue';
import RadioStatus from './components/RadioStatus.vue';

import Store from "./store"

export default {
  name: "App",
  components: {
    SelectRadio,
    RadioPanel,
    RadioPresets,
    RadioStatus
  },
  setup() {
    return {
      state: Store.getState(),
    }
  },
  methods: {
    updateRadios: () => Store.updateRadios()
  },
  mounted() {
    Store.updateRadios()
  },
  methods: {
    selectRadio: (uuid) => Store.selectRadio(uuid),
    setRadioPreset: (preset) => Store.setRadioPreset(preset),
    setRadioVolume: (volume) => Store.setRadioVolume(volume),
    discoverRadios: () => Store.discoverRadios(),
    toggleRadioPower: () => Store.toggleRadioPower()
  }
}
</script>

<template>
  <RadioStatus v-if="state.uuid" :connecting="state.connecting" :connected="state.connected" />
  <div class="container mx-auto px-2">
    <SelectRadio
      class="mt-2 pb-2 border-b-2"
      :uuid="state.uuid"
      :radios="state.radios"
      :selectRadio="selectRadio"
      :discoverRadios="discoverRadios"
    />
    <div v-if="state.uuid">
      <RadioPanel
        class="mt-2 pb-2 border-b-2"
        :radio="state.radio"
        :toggleRadioPower="toggleRadioPower"
        :setRadioVolume="setRadioVolume"
      />
      <RadioPresets class="mt-2" :radio="state.radio" :setRadioPreset="setRadioPreset" />
    </div>
  </div>
</template>

<style>
</style>
