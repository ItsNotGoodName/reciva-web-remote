<template>
  <div
    class="
      sticky
      top-0
      mx-auto
      border-l-2 border-b-2 border-r-2
      rounded-b
      p-2
      mb-2
      bg-white
      w-full
      sm:flex
      flex-wrap
      space-y-2
      gap-2
      sm:space-y-0
    "
  >
    <div class="flex gap-2">
      <button
        @click="toggleEdit"
        class="w-16 flex-1 sm:flex-inital btn-sm rounded"
        :class="{ 'btn-primary': !edit, 'btn-warning': edit }"
      >
        {{ edit ? "Edit" : "Play" }}
      </button>
      <loading-button
        class="w-24 flex-1 sm:flex-initial btn-sm btn-success rounded"
        :on-click="discoverRadios"
        >Discover</loading-button
      >
    </div>
    <div class="flex gap-2 flex-grow md:flex-auto">
      <select class="flex-1 btn-sm btn-light rounded" v-model="radioUUID">
        <option :value="null" disabled>Select Radio</option>
        <option :key="uuid" v-bind:value="uuid" v-for="(name, uuid) in radios">
          {{ name }}
        </option>
      </select>
      <loading-button
        v-if="radioUUID"
        class="w-24 Preset 02 rounded btn-sm btn-light"
        className="border-blue-600"
        :on-click="refreshRadio"
        >Refresh</loading-button
      >
    </div>
    <template v-if="radioConnected">
      <div class="flex gap-2 flex-1">
        <VolumeOffIcon v-if="radioisMuted" class="h-8" />
        <VolumeUpIcon v-else class="h-8" />
        <loading-button
          className="ml-0 border-blue-600"
          class="text-left flex-grow w-10"
          :on-click="refreshRadioVolume"
          >{{ radio.volume }}%</loading-button
        >
        <loading-button
          class="w-16 btn-light"
          className="border-blue-600"
          :on-click="decreaseRadioVolume"
        >
          <ChevronDownIcon class="h-8 rounded w-full" />
        </loading-button>
        <loading-button
          class="w-16 btn-light"
          className="border-blue-600"
          :on-click="increaseRadioVolume"
        >
          <ChevronUpIcon class="h-8 rounded w-full" />
        </loading-button>
        <loading-button
          v-if="radio.power"
          class="w-16 rounded btn-sm btn-success"
          :on-click="toggleRadioPower"
          >ON</loading-button
        >
        <loading-button
          v-else
          class="w-16 rounded btn-sm btn-danger"
          :on-click="toggleRadioPower"
          >OFF</loading-button
        >
      </div>
      <div class="flex gap-2 flex-grow">
        <div class="w-8 h-8 flex-shrink-0">
          <PlayIcon v-if="radio.state == 'Playing'" />
          <StopIcon v-else-if="radio.state == 'Stopped'" />
          <RefreshIcon v-else />
        </div>
        <div class="my-auto truncate space-x-2">
          <span
            class="px-2 py-0.5 text-base rounded-full text-white bg-blue-500"
            >{{ radio.title }}</span
          >
          <span v-if="radio.metadata">{{ radio.metadata }}</span>
        </div>
      </div>
    </template>
  </div>
</template>

<script>
import { mapState, mapActions } from "vuex";

import ChevronDownIcon from "@heroicons/vue/solid/ChevronDownIcon";
import ChevronUpIcon from "@heroicons/vue/solid/ChevronUpIcon";
import PlayIcon from "@heroicons/vue/solid/PlayIcon";
import RefreshIcon from "@heroicons/vue/solid/RefreshIcon";
import StopIcon from "@heroicons/vue/solid/StopIcon";
import VolumeOffIcon from "@heroicons/vue/solid/VolumeOffIcon";
import VolumeUpIcon from "@heroicons/vue/solid/VolumeUpIcon";

import LoadingButton from "./LoadingButton.vue";

export default {
  components: {
    ChevronDownIcon,
    ChevronUpIcon,
    PlayIcon,
    RefreshIcon,
    StopIcon,
    VolumeOffIcon,
    VolumeUpIcon,
    LoadingButton,
  },
  computed: {
    ...mapState([
      "config",
      "edit",
      "radio",
      "radioConnected",
      "radioConnecting",
      "radios",
    ]),
    radioUUID: {
      get() {
        return this.$store.state.radioUUID;
      },
      set(uuid) {
        this.$store.dispatch("setRadioUUID", uuid);
      },
    },
  },
  methods: {
    ...mapActions([
      "toggleEdit",
      "decreaseRadioVolume",
      "discoverRadios",
      "increaseRadioVolume",
      "refreshRadio",
      "refreshRadioVolume",
      "toggleRadioPower",
    ]),
  },
};
</script>

<style scoped>
</style>


