<template>
  <div class="field has-addons is-flex">
    <div class="control">
      <b-button
        title="Toggle Page"
        :class="[(page ? 'is-success' : 'is-info')]"
        @click="togglePage"
      >
        <b-icon icon="fa-bars" />
      </b-button>
    </div>
    <div class="control">
      <b-button
        title="Discover"
        class="is-success"
        :loading="radiosDiscovering"
        @click="discoverRadios"
      >
        <b-icon icon="fa-search" />
      </b-button>
    </div>
    <div class="control is-flex-grow-1">
      <b-select class="is-fullwidth" v-model="radioUUID" :loading="radiosLoading">
        <option selected disabled hidden value>Select Radio</option>
        <option v-for="radio in radios" :value="radio.uuid">{{ radio.name }}</option>
      </b-select>
    </div>
    <div class="control">
      <radio-refresh-button />
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters, mapActions } from "vuex";

import BButton from "../Bulma/BButton.vue";
import BIcon from "../Bulma/BIcon.vue";
import BSelect from "../Bulma/BSelect.vue";

import RadioRefreshButton from "../RadioRefreshButton.vue";

export default {
  components: {
    BButton,
    BIcon,
    BSelect,
    RadioRefreshButton
  },
  methods: {
    ...mapActions(["discoverRadios", "togglePage"]),
  },
  computed: {
    ...mapGetters(["radioSelected"]),
    ...mapState({
      page: (state) => state.page,
      radios: (state) => state.r.radios,
      radiosLoading: (state) => state.r.radiosLoading,
      radiosDiscovering: (state) => state.r.radiosDiscovering,
    }),
    radioUUID: {
      get() {
        return this.$store.state.r.radioUUID;
      },
      set(radioUUID) {
        this.$store.dispatch("setRadioUUID", radioUUID);
      },
    },
  },
};
</script>

<style lang="scss" scoped></style>
