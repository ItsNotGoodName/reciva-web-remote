<template>
  <div class="field has-addons is-flex">
    <div class="control">
      <RadiosDiscover />
    </div>
    <div class="control is-flex-grow-1">
      <div class="select is-fullwidth" :class="{ 'is-loading': radiosLoading }">
        <select v-model="radioUUID">
          <option selected disabled hidden value="">Select Radio</option>
          <option v-for="radio in radios" :value="radio.uuid">
            {{ radio.name }}
          </option>
        </select>
      </div>
    </div>
    <div class="control">
      <RadiosList />
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";

import RadiosDiscover from "./buttons/RadiosDiscover.vue";
import RadiosList from "./buttons/RadiosList.vue";

export default {
  components: {
    RadiosDiscover,
    RadiosList,
  },
  computed: {
    ...mapState({
      radios: (state) => state.r.radios,
      radiosLoading: (state) => state.r.radiosLoading,
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

<style></style>
