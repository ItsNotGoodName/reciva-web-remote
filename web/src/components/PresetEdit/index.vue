<template>
  <div v-if="!presetsLoading" class="is-flex">
    <table class="table is-bordered is-striped is-narrow is-hoverable mx-auto">
      <thead>
        <tr>
          <th>Edit</th>
          <th>URL</th>
          <th>New Name</th>
          <th>New URL</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="preset in presets"
          :key="preset.url"
          :class="{ 'is-info': radio.url == preset.url }"
        >
          <td>
            <preset-edit :preset="preset" />
          </td>
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
  <preset-form />
</template>

<script>
import { mapState } from "vuex";

import PresetForm from "./PresetForm.vue";
import PresetEdit from "./PresetEdit.vue";

export default {
  components: {
    PresetForm,
    PresetEdit,
  },
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
};
</script>

<style lang="scss" scoped></style>
