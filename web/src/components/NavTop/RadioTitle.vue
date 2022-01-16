<template>
  <div
    class="dropdown"
    :class="{ 'is-active': isActive }"
    @mouseleave="isActive = false"
  >
    <div class="dropdown-trigger is-flex-grow-1">
      <b-button
        :title="radio.title"
        :loading="!radioLoaded && radioSelected"
        class="is-info is-justify-content-flex-start is-fullwidth"
        aria-haspopup="true"
        aria-controls="title-dropdown-menu"
        @click="isActive = !isActive"
      >
        <b-tag class="mr-2">{{ radio.preset || 0 }}</b-tag>
        <span>
          {{ radio.title }}
        </span>
      </b-button>
    </div>
    <div class="dropdown-menu" id="title-dropdown-menu" role="menu">
      <div class="dropdown-content">
        <div class="dropdown-item">
          <b-tag class="is-warning mr-2">Metadata</b-tag>
          <span>
            {{ radio.metadata }}
          </span>
        </div>
        <div class="dropdown-item">
          <b-tag class="is-info mr-2">URL</b-tag>
          <a :href="radio.url">
            {{ radio.url }}
          </a>
        </div>
        <div class="dropdown-item">
          <b-tag class="is-success mr-2">New URL</b-tag>
          <a :href="radio.newURL">
            {{ radio.newURL }}
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";

import BButton from "../Bulma/BButton.vue";
import BTag from "../Bulma/BTag.vue";

export default {
  components: {
    BButton,
    BTag,
  },
  data() {
    return {
      isActive: false,
    };
  },
  computed: {
    ...mapGetters(["radioLoaded", "radioSelected"]),
    ...mapState({
      radio: (state) => state.r.radio,
    }),
  },
};
</script>

<style lang="scss" scoped></style>
