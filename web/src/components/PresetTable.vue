<template>
  <div v-if="!presetsLoading" class="table-container">
    <table class="table is-striped is-hoverable">
      <tbody>
        <tr>
          <th>Edit</th>
          <th>URL</th>
          <th>New Name</th>
          <th>New URL</th>
        </tr>
        <tr
          v-for="preset in presets"
          :key="preset.url"
          :class="{ 'is-info': radio.url == preset.url }"
        >
          <td><PresetEdit :preset="preset" /></td>
          <td>
            {{ preset.url }}
          </td>
          <td>{{ preset.newName }}</td>
          <td>{{ preset.newUrl }}</td>
        </tr>
      </tbody>
    </table>
  </div>
  <progress v-else class="progress is-success" max="100" />
</template>

<script>
import { mapState } from "vuex";

import PresetEdit from "./buttons/PresetEdit.vue";

export default {
  beforeMount() {
    this.$store.dispatch("listPresets");
  },
  computed: {
    ...mapState({
      presets: (state) => state.p.presets,
      presetsLoading: (state) => state.p.presetsLoading,
      radio: (state) => state.r.radio,
    }),
  },
  components: { PresetEdit },
};
</script>

<style></style>
