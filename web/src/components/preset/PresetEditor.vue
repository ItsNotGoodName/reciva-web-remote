<template>
  <div class="p-fluid">
    <div class="card" v-if="!loading">
      <DataTable
        :value="presets"
        sortField="url"
        dataKey="url"
        :sortOrder="1"
        responsiveLayout="scroll"
      >
        <Column :sortable="true" field="url" header="URL" />
        <Column :sortable="true" field="newName" header="New Name" />
        <Column :sortable="true" field="newUrl" header="New URL">
          <template #body="slotProps">
            <pre
              style="max-height: 80px"
              class="my-0 overflow-x-hidden overflow-y-auto"
              >{{ slotProps.data.newUrl }}</pre
            >
          </template>
        </Column>
        <Column :exportable="false" class="w-1rem">
          <template #body="slotProps">
            <Button
              icon="pi pi-pencil"
              class="p-button-rounded p-button-success p-mr-2"
              @click="editPreset(slotProps.data)"
            />
          </template>
        </Column>
      </DataTable>
      <!-- eslint-disable vue/no-v-model-argument -->
      <Dialog
        v-model:visible="presetDialog"
        :style="{ width: '450px' }"
        header="Preset Details"
        :modal="true"
        class="p-fluid"
      >
        <!--eslint-enable-->
        <div class="p-field mb-2">
          <label for="url">URL</label>
          <InputText
            id="url"
            v-model.trim="preset.url"
            disabled
            required="true"
            autofocus
          />
        </div>

        <div class="p-field mb-2">
          <label for="newName">New Name</label>
          <InputText
            id="newName"
            v-model.trim="preset.newName"
            required="true"
            autofocus
          />
        </div>

        <div class="p-field mb-2">
          <label for="newUrl">New URL</label>
          <TextArea
            id="newUrl"
            v-model="preset.newUrl"
            required="true"
            rows="3"
            cols="20"
          />
        </div>

        <template #footer>
          <Button
            label="Cancel"
            icon="pi pi-times"
            class="p-button-text"
            @click="hideDialog"
          />
          <Button
            label="Save"
            icon="pi pi-check"
            class="p-button-text"
            @click="savePreset"
            :loading="submitting"
          />
        </template>
      </Dialog>
    </div>
    <div class="flex" v-else>
      <ProgressSpinner class="mx-auto" />
    </div>
  </div>
</template>

<script>
import { mapActions, mapState } from "vuex";

import Button from "primevue/button";
import Column from "primevue/column";
import DataTable from "primevue/datatable";
import Dialog from "primevue/dialog";
import InputText from "primevue/inputtext";
import ProgressSpinner from "primevue/progressspinner";
import TextArea from "primevue/textarea";

export default {
  components: {
    Button,
    Column,
    DataTable,
    Dialog,
    InputText,
    ProgressSpinner,
    TextArea,
  },
  data() {
    return {
      loading: false,
      preset: {},
      presetDialog: false,
      submitting: false,
    };
  },
  created() {
    this.loading = true;
    this.readPresets()
      .catch((err) => {
        this.$toast.add({
          severity: "error",
          summary: "could not fetch presets",
          detail: err,
        });
      })
      .finally(() => {
        this.loading = false;
      });
  },
  methods: {
    ...mapActions(["readPresets", "updatePreset"]),
    savePreset() {
      this.submitting = true;

      console.log(this.preset);
      this.updatePreset(this.preset)
        .then(() => {
          this.$toast.add({
            severity: "success",
            summary: "preset updated",
            life: 1000,
          });
        })
        .catch((err) => {
          this.$toast.add({
            severity: "error",
            summary: "preset update failed",
            detail: err,
            life: 3000,
          });
        })
        .finally(() => {
          this.submitting = false;
          this.presetDialog = false;
          this.preset = {};
        });
    },
    editPreset(preset) {
      this.preset = { ...preset };
      this.presetDialog = true;
    },
    openNew() {
      this.preset = {};
      this.presetDialog = true;
    },
    hideDialog() {
      this.presetDialog = false;
    },
  },
  computed: {
    ...mapState({
      presets: (state) => state.p.presets,
    }),
  },
};
</script>

<style scoped>
</style>