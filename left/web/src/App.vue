<script setup lang="ts">
import { computed, watch, ref } from "vue";

import { PAGE_HOME, PAGE_EDIT } from "./constants";
import { useWS, useRadiosQuery, useRadioSubscriptionMutation, useRadioUUIDStorage } from "./hooks"

import RadioStatus from "./components/RadioStatus.vue";
import RadioTitle from "./components/RadioTitle.vue";
import RadioPower from "./components/RadioPower.vue";
import RadioName from "./components/RadioName.vue";
import DButton from "./components/DaisyUI/DButton.vue";
import DErrorAlert from "./components/DaisyUI/DErrorAlert.vue";
import RadioAudioSource from "./components/RadioAudioSource.vue";
import HamburgerMenu from "./components/HamburgerMenu.vue";
import RadioVolume from "./components/RadioVolume.vue"
import RadiosDiscover from "./components/RadiosDiscover.vue";
import HomePage from "./pages/Home.vue";
import EditPage from "./pages/Edit.vue"

const page = ref(PAGE_HOME);
const updatePage = (value: string) => {
  page.value = value
}

const radioUUID = useRadioUUIDStorage();
const { data: radios, isLoading: radiosLoading, isError: radiosIsError, refetch: radiosRefetch, isFetching: radiosFetching } = useRadiosQuery();
const { state, stateLoading, stateSelected, connecting, disconnected, reconnect } = useWS(radioUUID);
const { mutate: radioSubscriptionMutate, isLoading: radioSubscriptionLoading } = useRadioSubscriptionMutation();
const refreshing = computed(() => radiosFetching.value || radioSubscriptionLoading.value);

// Make sure websocket is connected when fetching radios
watch(radiosFetching, () => {
  reconnect()
})

// Make sure radioUUID is a valid radio
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

// Refresh current radio if selected and refetch radios
const refresh = () => {
  if (stateSelected.value) {
    radioSubscriptionMutate(radioUUID.value)
  }

  radiosRefetch.value()
}
</script>

<template>
  <div class="h-screen">
    <!-- Top Navbar -->
    <div class="navbar bg-base-200 fixed top-0 flex gap-2 z-50 border-b-2 border-b-base-300">
      <radio-status :state="state" :loading="stateLoading" />
      <radio-title class="flex-grow w-full" :state="state" :loading="stateLoading" />
    </div>
    <div class="container mx-auto px-4 pt-20 pb-36">
      <!-- Home Page -->
      <home-page v-if="page == PAGE_HOME" :state="state" />
      <!-- Edit Page -->
      <edit-page v-else-if="page == PAGE_EDIT" />
    </div>
    <!-- Bottom -->
    <div class="fixed bottom-0 w-full space-y-2 z-50">
      <!--- Alerts -->
      <div class="ml-auto px-2 max-w-screen-sm space-y-2">
        <d-error-alert v-if="radiosIsError">
          Failed to list radios.
        </d-error-alert>
        <div v-if="disconnected" class="alert shadow-lg">
          <div>
            <v-icon class="text-info" name="fa-info-circle" />
            <span>Disconnected from server.</span>
          </div>
          <div class="flex-none">
            <d-button class="btn-sm btn-primary" :loading="connecting" @click="reconnect">Reconnect</d-button>
          </div>
        </div>
      </div>
      <!--- Navbar -->
      <div class="navbar bg-base-200 flex flex-wrap-reverse gap-2 border-t-2 border-t-base-300">
        <!--- Radios Toolbar -->
        <div class="flex-auto flex gap-2 z-10">
          <hamburger-menu :page="page" @update:page="updatePage" />
          <div class="flex-auto flex">
            <radios-discover class="btn-primary w-14 rounded-none rounded-l-md" />
            <select v-model="radioUUID" :disabled="radiosLoading"
              class="flex-auto w-36 select select-primary rounded-none">
              <option disabled selected value="">Select Radio</option>
              <template v-if="radios">
                <option :key="r.uuid" v-for="r in radios" :value="r.uuid">{{ r.name }}</option>
              </template>
            </select>
            <div class="tooltip" data-tip="Refresh">
              <d-button class="btn-primary w-14 rounded-none rounded-r-md" aria-label="Refresh" :loading="refreshing"
                @click="refresh">
                <v-icon name="fa-redo" />
              </d-button>
            </div>
          </div>
        </div>
        <!--- Radio Toolbar -->
        <div v-if="stateSelected" class="flex-auto flex gap-2 ">
          <radio-power class="flex-auto" :state="state" />
          <radio-volume :state="state" />
          <radio-audio-source :state="state" />
          <radio-name :state="state" />
        </div>
      </div>
    </div>
  </div>
</template>

<style>
</style>
