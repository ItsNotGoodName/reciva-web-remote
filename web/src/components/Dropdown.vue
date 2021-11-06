<template>
  <div class="relative">
    <button
      class="
        btn btn-secondary
        relative
        p-2
        focus:outline-none focus:bg-gray-600
      "
      @click="toggle()"
    >
      D
    </button>
    <div
      v-if="isOpen"
      class="
        absolute
        right-0
        mt-2
        w-48
        bg-white
        rounded
        overflow-hidden
        shadow-xl
        z-20
      "
      v-on:unfocus="toggle()"
    >
      <button
        @click="updatePreset({ ...p, sid: sid })"
        class="btn w-full btn-white"
        :key="p.uri"
        v-for="p in presets"
      >
        {{ p.uri }}
      </button>
    </div>
  </div>
</template>

<script>
import { mapActions, mapState } from "vuex";

export default {
  name: "Dropdown",
  props: {
    sid: {
      type: Number,
    },
  },
  data() {
    return {
      isOpen: false,
    };
  },
  methods: {
    toggle() {
      this.isOpen = !this.isOpen;
    },
    ...mapActions(["updatePreset"]),
  },
  computed: {
    ...mapState({
      presets: (state) => state.p.presets,
    }),
  },
};
</script>

<style lang="scss" scoped>
</style>