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
      state: Store.getState()
    }
  },
  mounted() {
    Store.refreshRadio().then((radios) => {
      for (let k in radios) {
        Store.selectRadio(k)
        return
      }
      console.log("no radios found")
    })
  }
}
</script>

<template>
  The radios on the network are {{ Object.keys(this.state.radios) }}.
  <br />
  <div v-if="this.state.radio">
    The current radio is {{ this.state.radio.uuid }}.
    <br />
    Power: {{ this.state.radio.power ? "ON" : "OFF" }}.
    <br />
    URL: {{ this.state.radio.url }}.
    <br />
    State: {{ this.state.radio.state }}.
    <Player :radio="this.state.radio" />
  </div>
</template>

<style>
</style>
