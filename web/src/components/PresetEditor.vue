<template>
  <div class="lg:flex gap-2">
    <Presets class="flex-grow mb-2" />
    <div class="flex-auto">
      <edit-preset class="mb-2" />
    </div>
  </div>
</template>

<script>
import EditPreset from "./EditPreset.vue";
import Presets from "./Presets.vue";

export default {
  components: {
    EditPreset,
    Presets,
  },
};
</script>

<style>
</style>