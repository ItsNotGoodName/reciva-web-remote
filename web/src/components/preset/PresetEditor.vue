<template>
  <div class="p-fluid">
    <div class="card" v-if="!loading">
      <DataTable
        :value="presets"
        sortField="url"
        :sortOrder="1"
        editMode="cell"
        @cell-edit-complete="onCellEditComplete"
        class="editable-cells-table"
        responsiveLayout="scroll"
      >
        <Column :sortable="true" field="url" header="URL" key="url" />
        <Column :sortable="true" field="newUrl" header="New URL" key="newUrl">
          <template #editor="{ data, field }">
            <InputText v-model="data[field]" autofocus />
          </template>
        </Column>
        <Column
          :sortable="true"
          field="newName"
          header="New Name"
          key="newName"
        >
          <template #editor="{ data, field }">
            <InputText v-model="data[field]" autofocus />
          </template>
        </Column>
      </DataTable>
    </div>
    <div class="flex" v-else>
      <ProgressSpinner class="mx-auto" />
    </div>
  </div>
</template>

<script>
import { mapActions, mapState } from "vuex";

import Column from "primevue/column";
import DataTable from "primevue/datatable";
import InputText from "primevue/inputtext";
import ProgressSpinner from "primevue/progressspinner";

export default {
  components: {
    Column,
    DataTable,
    InputText,
    ProgressSpinner,
  },
  data() {
    return {
      loading: false,
    };
  },
  created() {
    this.loading = true;
    this.readPresets()
      .catch((err) => {
        this.$toast.add({
          severity: "error",
          summary: "Could not fetch Presets",
          detail: err,
        });
      })
      .finally(() => {
        this.loading = false;
      });
  },
  methods: {
    ...mapActions(["readPresets", "updatePreset"]),
    onCellEditComplete(event) {
      let { data, newValue, field } = event;
      data[field] = newValue;
      this.updatePreset(data)
        .then(() => {
          this.$toast.add({
            severity: "success",
            summary: "Preset Updated",
            life: 1000,
          });
        })
        .catch((err) => {
          this.$toast.add({
            severity: "error",
            summary: "Preset Update Failed",
            detail: err,
            life: 3000,
          });
        });
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
::v-deep(.editable-cells-table td.p-cell-editing) {
  padding-top: 0;
  padding-bottom: 0;
}
</style>