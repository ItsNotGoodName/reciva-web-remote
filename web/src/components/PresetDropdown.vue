<template>
  <div class="relative">
    <button
      class="btn btn-white relative h-full"
      :class="{ 'btn-white-on': isOpen }"
      @click="toggle()"
    >
      <DotsHorizontalIcon class="w-5" />
    </button>
    <div
      v-if="isOpen"
      class="
        absolute
        top-0
        right-0
        mr-14
        w-48
        bg-white
        rounded
        overflow-hidden
        shadow-xl
        z-20
      "
    >
      <div :key="p.uri" v-for="p in presets">
        <button
          v-if="p.sid == 0"
          @click="updatePreset({ ...p, sid: sid })"
          class="btn w-full btn-white"
        >
          {{ p.uri }}
        </button>
        <button
          v-else-if="p.sid != sid"
          @click="updatePreset({ ...p, sid: sid })"
          class="btn w-full btn-light"
        >
          {{ p.uri }}
        </button>
        <button
          v-else
          @click="clearPreset(p.uri)"
          class="btn w-full btn-primary"
        >
          {{ p.uri }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { mapActions, mapState } from "vuex";
import DotsHorizontalIcon from "@heroicons/vue/solid/DotsHorizontalIcon";

export default {
  name: "Dropdown",
  components: {
    DotsHorizontalIcon,
  },
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
      this.isOpen &&
        setTimeout(() => {
          document.addEventListener("click", this.close);
        }, 250);
    },
    close() {
      this.isOpen = false;
      document.removeEventListener("click", this.close);
    },
    ...mapActions(["updatePreset", "clearPreset"]),
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