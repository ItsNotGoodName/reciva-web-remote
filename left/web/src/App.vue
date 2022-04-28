<script setup lang="ts">
import { watch, ref } from "vue";

import { PAGE_HOME, PAGE_EDIT } from "./constants";
import { useWS, useRadiosQuery, useRadioSubscriptionMutation, useRadioUUID } from "./hooks"

import RadioStatus from "./components/RadioStatus.vue";
import RadioTitle from "./components/RadioTitle.vue";
import RadioPower from "./components/RadioPower.vue";
import RadioName from "./components/RadioName.vue";
import DButton from "./components/DaisyUI/DButton.vue";
import RadioAudioSource from "./components/RadioAudioSource.vue";
import HamburgerMenu from "./components/HamburgerMenu.vue";
import RadioVolume from "./components/RadioVolume.vue"
import RadiosDiscover from "./components/RadiosDiscover.vue";
import RadioPresets from "./components/RadioPresets.vue";

const page = ref(PAGE_HOME);
const setPage = (value: string) => {
  page.value = value
}

const radioUUID = useRadioUUID();
const { data: radios, isLoading: radiosLoading, refetch: radiosRefetch } = useRadiosQuery();
const { radio, connecting: wsConnecting, disconnected: wsDisconnected, reconnect: wsReconnect } = useWS(radioUUID);
const { mutate: radioSubscriptionMutate, isLoading: radioSubscriptionLoading } = useRadioSubscriptionMutation();

// Make sure a valid radio is selected
watch(radios, (newRadios) => {
  if (newRadios) {
    for (const r of newRadios) {
      if (r.uuid == radioUUID.value) {
        return
      }
    }
    radioUUID.value = ""
  }
});

// Refetch radios and refresh currently selected radio
const onRefreshClick = () => {
  if (radioUUID.value) {
    radioSubscriptionMutate(radioUUID.value)
  }

  radiosRefetch.value()
}
</script>

<template>
  <div class="h-screen">
    <!-- Top Navbar -->
    <div class="navbar bg-base-200 fixed top-0 flex gap-2 z-50 border-b-2 border-b-base-300">
      <radio-status :radio="radio" :loading="wsConnecting" />
      <radio-title class="flex-grow w-full" :radio="radio" :loading="wsConnecting" />
    </div>
    <div class="mx-5 pt-20 pb-36">
      <!-- Homepage -->
      <div v-if="page == PAGE_HOME && radio" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <radio-presets :radio="radio" />
      </div>
      <!-- Edit Presets -->
      <div v-else-if="page == PAGE_EDIT">
        Hello World
      </div>
    </div>
    <!-- Bottom -->
    <div class="fixed bottom-0 w-full space-y-2">
      <!--- Radio Websocket Disconnect Alert -->
      <div v-if="wsDisconnected" class="ml-auto px-2 max-w-screen-sm">
        <div class="alert shadow-lg">
          <div>
            <v-icon class="text-info" name="fa-info-circle" />
            <span>Disconnected from server.</span>
          </div>
          <div class="flex-none">
            <d-button class="btn-sm btn-primary" :loading="wsConnecting" @click="wsReconnect">Reconnect</d-button>
          </div>
        </div>
      </div>
      <!--- Navbar -->
      <div class="navbar bg-base-200 flex flex-row-reverse flex-wrap gap-2 pb-4 px-4 z-50 border-t-2 border-t-base-300">
        <!--- Radio Toolbar -->
        <div v-if="radio" class="flex-grow md:flex-grow-0 flex gap-2">
          <radio-power class="flex-grow" :radio="radio" />
          <radio-volume :radio="radio" />
          <radio-audio-source :radio="radio" />
          <radio-name :radio="radio" />
        </div>
        <!--- Radios Toolbar -->
        <div class="grow flex gap-2">
          <hamburger-menu :page="page" :set-page="setPage" />
          <div class="grow flex">
            <radios-discover class="btn-primary w-12 rounded-none rounded-l-md" />
            <select v-model="radioUUID" :disabled="radiosLoading" class="select select-primary rounded-none flex-grow">
              <option disabled selected value="">Select Radio</option>
              <template v-if="radios">
                <option :key="r.uuid" v-for="r in radios" :value="r.uuid">{{ r.name }}</option>
              </template>
            </select>
            <div class="tooltip" data-tip="Refresh">
              <d-button class="btn-primary w-12 rounded-none rounded-r-md" aria-label="Refresh"
                :loading="radiosLoading || radioSubscriptionLoading" @click="onRefreshClick">
                <v-icon name="fa-redo" />
              </d-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
</style>
