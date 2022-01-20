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
        <tr v-for="preset in presets" :key="preset.url">
          <td>
            <preset-edit-button :preset="preset" />
          </td>
          <td>
            <a :href="preset.url">{{ preset.url }}</a>
          </td>
          <td :class="{ 'is-info': radio.url == preset.url }">{{ preset.newName }}</td>
          <td>
            <a :href="preset.newUrl">{{ preset.newUrl }}</a>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
  <progress v-else class="progress is-success" max="100" />
  <preset-form-modal />
</template>

<script>
import { mapState } from "vuex";

import PresetFormModal from "../components/PresetFormModal.vue";
import PresetEditButton from "../components/PresetEditButton.vue";

export default {
  components: {
    PresetFormModal,
    PresetEditButton,
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
