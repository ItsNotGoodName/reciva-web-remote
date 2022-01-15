<template>
  <div
    class="dropdown is-flex"
    :class="{ 'is-active': isActive }"
    @mouseleave="isActive = false"
  >
    <div class="dropdown-trigger is-flex-grow-1">
      <button
        :title="radio.title"
        class="button is-info has-text-left is-fullwidth"
        :class="{ 'is-loading': !radioLoaded }"
        aria-haspopup="true"
        aria-controls="title-dropdown-menu"
        @click="isActive = !isActive"
      >
        <span class="tag mr-2 is-white">{{ radio.preset }}</span>
        <span class="is-flex-grow-1">
          {{ radio.title }}
        </span>
      </button>
    </div>
    <div class="dropdown-menu" id="title-dropdown-menu" role="menu">
      <div class="dropdown-content">
        <div class="dropdown-item">
          <span class="tag mr-2 is-rounded is-warning">Metadata</span>
          {{ radio.metadata }}
        </div>
        <div class="dropdown-item">
          <span class="tag mr-2 is-rounded is-info">URL</span>
          <a :href="radio.url">
            {{ radio.url }}
          </a>
        </div>
        <div class="dropdown-item">
          <span class="tag mr-2 is-rounded is-success">New URL</span>
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

export default {
  data() {
    return {
      isActive: false,
    };
  },
  computed: {
    ...mapGetters(["radioLoaded"]),
    ...mapState({
      radio: (state) => state.r.radio,
    }),
  },
};
</script>

<style scoped>
.button.has-text-left {
  justify-content: flex-start;
}
</style>
