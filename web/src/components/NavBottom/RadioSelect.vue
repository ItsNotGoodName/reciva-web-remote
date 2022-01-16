<template>
  <div class="field has-addons is-flex">
    <div class="control">
      <b-button
        title="Toggle Page"
        :class="{ 'is-info': page === 'play' }"
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
      <b-select
        class="is-fullwidth"
        v-model="radioUUID"
        :loading="radiosLoading"
      >
        <option selected disabled hidden value="">Select Radio</option>
        <option v-for="radio in radios" :value="radio.uuid">
          {{ radio.name }}
        </option>
      </b-select>
    </div>
    <div class="control">
      <b-button title="List" :loading="radiosLoading" @click="listRadios">
        <b-icon icon="fa-arrow-down"></b-icon>
      </b-button>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters, mapActions } from "vuex";

import BButton from "../Bulma/BButton.vue";
import BIcon from "../Bulma/BIcon.vue";
import BSelect from "../Bulma/BSelect.vue";

export default {
  components: {
    BButton,
    BIcon,
    BSelect,
  },
  methods: {
    ...mapActions(["togglePage", "discoverRadios", "listRadios"]),
  },
  computed: {
    ...mapGetters(["radioSelected"]),
    ...mapState({
      radios: (state) => state.r.radios,
      radiosLoading: (state) => state.r.radiosLoading,
      radiosDiscovering: (state) => state.r.radiosDiscovering,
      page: (state) => state.page,
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
