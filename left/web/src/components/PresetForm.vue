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

const { data, isFetching: dataIsFetching, isError: dataIsError, error: dataError, isSuccess: dataSuccess, refetch: dataRefetch } = usePresetQuery(computed(() => props.presetUrl))
const { mutate, isLoading: mutateIsLoading, isError: mutateIsError, error: mutateError } = usePresetMutation()
const loading = computed(() => dataIsFetching.value || mutateIsLoading.value)
const form = reactive({
  url: '',
  url_new: '',
  url_new_changed: false,
  title_new: '',
  title_new_changed: false
})

const formSync = () => {
  if (dataSuccess.value && data.value) {
    form.url = data.value.url
    form.url_new = data.value.url_new
    form.url_new_changed = false
    form.title_new = data.value.title_new
    form.title_new_changed = false
  }
}
const reset = () => dataRefetch.value().then(formSync)
const submit = () => mutate({ url: form.url, url_new: form.url_new, title_new: form.title_new })

watchEffect(formSync)
</script>

<template>
  <form class="space-y-2" @submit.prevent="submit">
    <h1 class="text-center text-2xl">Edit Preset</h1>
    <d-error-alert v-if="dataIsError">Failed to get preset.</d-error-alert>
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
    <div class="btn-group fleex">
      <d-button class="flex-1" type="button" @click="emit('close')">Close</d-button>
      <d-button class="flex-1 btn-error" :loading="loading" type="button" @click="reset">
        Reset
      </d-button>
      <d-button class="flex-1 btn-success" :loading="loading" type="submit">Save</d-button>
    </div>
    <d-error-alert v-if="mutateIsError">Failed to update preset.</d-error-alert>
  </form>
</template>

<style>
</style>
