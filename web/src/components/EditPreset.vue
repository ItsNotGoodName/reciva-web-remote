<template>
  <div :class="{ hidden: !presetShow }">
    <div class="p-2 bg-info rounded-t">
      <strong class="px-2">Preset Details</strong>
    </div>
    <div class="flex border-l-2 border-r-2 b">
      <form class="w-full mx-auto bg-white rounded p-2">
        <div class="mb-4">
          <label class="block text-sm font-bold mb-2" for="newName">
            Preset New Name
          </label>
          <input
            v-model="presetNewName"
            :class="{ 'border-yellow-500': preset.newNameChanged }"
            class="border rounded w-full py-2 px-3"
            :disabled="presetLoading"
            id="newName"
            type="text"
            placeholder="Preset New Name"
          />
        </div>
        <div class="mb-4">
          <label class="block text-sm font-bold mb-2" for="newUrl">
            Preset New URL
          </label>
          <input
            v-model="presetNewUrl"
            :class="{ 'border-yellow-500': preset.newUrlChanged }"
            class="border rounded w-full py-2 px-3"
            :disabled="presetLoading"
            id="newUrl"
            type="text"
            placeholder="Preset New URL"
          />
        </div>
      </form>
    </div>
    <div class="p-2 flex flex-row-reverse rounded-b bg-light">
      <button @click="closePreset()" class="btn btn-secondary">Close</button>
      <loading-button
        :on-click="updatePreset"
        class="btn text-white btn-success rounded-l w-16"
      >
        Save
      </loading-button>
    </div>
  </div>
</template>

<script>
import { mapActions, mapState } from "vuex";

import LoadingButton from "./LoadingButton.vue";

export default {
  components: {
    LoadingButton,
  },
  computed: {
    ...mapState({
      presetShow: (state) => state.p.presetShow,
      preset: (state) => state.p.preset,
      presetLoading: (state) => state.p.presetLoading,
    }),
    presetNewName: {
      get() {
        return this.$store.state.p.preset.newName;
      },
      set(newName) {
        this.$store.commit("SET_PRESET_NEW_NAME", newName);
      },
    },
    presetNewUrl: {
      get() {
        return this.$store.state.p.preset.newUrl;
      },
      set(newUrl) {
        this.$store.commit("SET_PRESET_NEW_URL", newUrl);
      },
    },
  },
  methods: {
    ...mapActions(["loadPresets", "closePreset", "updatePreset"]),
  },
};
</script>

<style scoped>
</style>