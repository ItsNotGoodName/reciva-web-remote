<script setup lang="ts">
import { computed, watchEffect, reactive } from "vue"
import { usePresetQuery, usePresetMutation } from "../hooks"

import DButton from './DaisyUI/DButton.vue';
import DErrorAlert from "./DaisyUI/DErrorAlert.vue";

const emit = defineEmits(['close'])

const props = defineProps({
  presetUrl: {
    type: String,
    required: true,
  },
});

const { data, isFetching: dataIsFetching, isError: dataIsError, error: dataError, isSuccess: dataSuccess } = usePresetQuery(computed(() => props.presetUrl))
const loading = computed(() => dataIsFetching.value || !data.value)
const { mutate, isLoading: mutateIsLoading, isError: mutateIsError, error: mutateError } = usePresetMutation()

const form = reactive({
  url: '',
  url_new: '',
  url_new_changed: false,
  title_new: '',
  title_new_changed: false
})

watchEffect(() => {
  if (dataSuccess.value && data.value) {
    form.url = data.value.url
    form.url_new = data.value.url_new
    form.url_new_changed = false
    form.title_new = data.value.title_new
    form.title_new_changed = false
  }
})

const submit = () => mutate({ url: form.url, url_new: form.url_new, title_new: form.title_new })
</script>

<template>
  <form class="space-y-2" @submit.prevent="submit">
    <d-error-alert v-if="dataIsError" :error="dataError">Failed to get preset.</d-error-alert>
    <h1 class="text-center text-2xl">Edit Preset</h1>
    <div class="form-control">
      <label class="label">
        <span class="label-text">URL</span>
      </label>
      <input type="text" placeholder="URL" class="input input-bordered" :value="form.url" readonly :disabled="loading">
    </div>
    <div class="form-control">
      <label class="label">
        <span class="label-text">New Title</span>
      </label>
      <input type="text" placeholder="New Title" class="input input-bordered"
        :class="{ 'input-warning': form.title_new_changed }" @input="form.title_new_changed = true"
        v-model="form.title_new" :disabled="loading">
    </div>
    <div class="form-control">
      <label class="label">
        <span class="label-text">New URL</span>
      </label>
      <textarea placeholder="New URL" class="textarea textarea-bordered h-24"
        :class="{ 'textarea-warning': form.url_new_changed }" @input="form.url_new_changed = true"
        v-model="form.url_new" :disabled="loading" />
    </div>
    <div class="flex space-x-2">
      <d-button class="btn-success ml-auto" :loading="mutateIsLoading || loading" type="submit">Save</d-button>
      <d-button type="button" @click="emit('close')">Close</d-button>
    </div>
    <d-error-alert v-if="mutateIsError" :error="mutateError">Failed to update preset.</d-error-alert>
  </form>
</template>

<style>
</style>
