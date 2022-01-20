<template>
  <b-modal title="Edit Preset" v-model="presetVisible">
    <template v-slot:default>
      <div class="field">
        <label class="label">URL</label>
        <div class="control">
          <input :value="preset.url" class="input" type="text" placeholder="URL" disabled />
        </div>
      </div>

      <div class="field">
        <label class="label">New Name</label>
        <div class="control">
          <input v-model="newName" class="input" type="text" placeholder="New Name" />
        </div>
      </div>

      <div class="field">
        <label class="label">New URL</label>
        <div class="control">
          <textarea v-model="newUrl" class="textarea" placeholder="New URL"></textarea>
        </div>
      </div>
    </template>

    <template v-slot:footer>
      <div class="field is-grouped">
        <div class="control">
          <b-button @click="submitPreset" :loading="presetLoading" class="is-link">Submit</b-button>
        </div>
        <div class="control">
          <b-button class="is-link is-light" @click="hidePreset">Cancel</b-button>
        </div>
      </div>
    </template>
  </b-modal>
</template>

<script>
import { mapActions, mapState } from "vuex";

import BModal from "./Bulma/BModal.vue";
import BButton from "./Bulma/BButton.vue";

export default {
  computed: {
    presetVisible: {
      get() {
        return this.$store.state.p.presetVisible;
      },
      set() {
        this.$store.dispatch("hidePreset");
      },
    },
    newName: {
      get() {
        return this.$store.state.p.preset.newName;
      },
      set(value) {
        this.$store.commit("SET_PRESET_NEW_NAME", value);
      },
    },
    newUrl: {
      get() {
        return this.$store.state.p.preset.newUrl;
      },
      set(value) {
        this.$store.commit("SET_PRESET_NEW_URL", value);
      },
    },
    ...mapState({
      preset: (state) => state.p.preset,
      presetLoading: (state) => state.p.presetLoading,
    }),
  },
  methods: {
    ...mapActions(["hidePreset", "submitPreset"]),
  },
  components: {
    BModal,
    BButton,
  },
};
</script>

<style lang="scss" scoped></style>
