<script>
import Player from './components/Player.vue';
import Store from "./store"

export default {
  name: "App",
  components: {
    Player,
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
    selectRadio: (event) => Store.selectRadio(event.target.value),
    setRadioPreset: (event) => Store.setRadioPreset(parseInt(event.target.value)),
    setRadioVolume: (event) => Store.setRadioVolume(parseInt(event.target.value)),
    discoverRadios: () => Store.discoverRadios(),
    toggleRadioPower: () => Store.toggleRadioPower()
  }
}
</script>

<template>
  <div class="container mx-auto">
    <select v-on:change="selectRadio($event)" multiple>
      <option :value="k" v-for="r,k in this.state.radios" :key="r">{{ r }}</option>
    </select>
    <br />
    <button
      @click="discoverRadios()"
      class="bg-blue-300 text-white font-bold py-2 px-4 rounded"
    >Discover</button>

    <div v-if="this.state.uuid">
      URL: {{ this.state.radio.url }}.
      <br />
      State: {{ this.state.radio.state }}.
      <br />
      <button
        v-bind:class="{ 'hover:bg-green-500': this.state.radio.power, 'bg-green-300': this.state.radio.power, 'hover:bg-red-500': !this.state.radio.power, 'bg-red-300': !this.state.radio.power }"
        class="text-white font-bold py-2 px-4 rounded"
        @click="toggleRadioPower()"
      >{{ this.state.radio.power ? "ON" : "OFF" }}</button>
      <br />
      <select v-on:change="setRadioPreset($event)" multiple>
        <option
          :value="i + 1"
          v-for="i in Array(this.state.radio.presets).keys() "
          :key="i + 1"
        >Preset {{ i + 1 }}</option>
      </select>
      <br />
      <select v-bind="this.state.radio.volume" v-on:change="setRadioVolume($event)" multiple>
        <option :value="i + 1" v-for="i in Array(100).keys() " :key="i + 1">Volume {{ i + 1 }}</option>
      </select>
      <Player :radio="this.state.radio" />
    </div>
  </div>
</template>

<style>
</style>

