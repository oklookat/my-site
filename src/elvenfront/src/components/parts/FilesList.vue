<template>
  <div class="files__list" v-if="normalized && normalized.length > 0">
    <div class="file" v-for="file in normalized" :key="file.id" v-on:click="onSelected(file)">
      <div class="file__meta">
        <div class="file__item file__uploaded-date">{{ convertDate(file.created_at) }}</div>
      </div>
      <div class="file__main">
        <div class="file__item file__preview" v-on:click.stop v-if="file.extension === 'IMAGE'">
          <img :src="file.path" v-if="file.extension === 'IMAGE'" />
        </div>
        <div class="file__item file__preview" v-on:click.stop v-if="file.extension === 'VIDEO'">
          <video controls :src="file.path"></video>
        </div>

        <div class="file__item file__name">{{ file.original_name }}</div>
        <div class="file__item file__size">{{ convertSize(file.size) }}</div>
      </div>
    </div>

    <Overlay
      v-bind:active="selectionOverlay && selected != null"
      v-on:deactivated="selectionOverlay = false; selected = null"
    >
      <div class="overlay__selected">
        <div
          class="overlay__item file__play"
          v-if="selected.extension === 'AUDIO'"
          v-on:click="playAudio(selected.path)"
        >play</div>
        <div class="overlay__item file__copy-link" v-on:click="copyLink(selected)">copy link</div>
        <div class="overlay__item file__delete" v-on:click="onDelete(selected)">delete</div>
      </div>
    </Overlay>
  </div>
</template>

<script setup lang="ts">
import { IFile } from '@/types/file'
import Dates from '@/common/tools/Dates'
import Sizes from "@/common/tools/Sizes.js"
import Extensions from "@/common/tools/Extensions"
import Overlay from "@/components/ui/Overlay.vue"
import { computed, ComputedRef, ref, Ref, toRef, toRefs } from '@vue/reactivity'

const selected: Ref<IFile | null> = ref(null)
const selectionOverlay = ref(false)

const props = defineProps<{
  files: Array<IFile>
}>()

const normalized: ComputedRef<Array<IFile>> = computed(() => {
  return props.files.map(file => {
    file = converter(file)
    return file
  })
})

const emit = defineEmits<{
  (e: 'selected', file: IFile): void
  (e: 'delete', file: IFile): void
}>()

function onSelected(file: IFile) {
  selectionOverlay.value = true
  selected.value = file
  emit('selected', file)
}

function onDelete(file: IFile) {
  selectionOverlay.value = false
  emit('delete', file)
}

async function copyLink(file: IFile) {
  navigator.clipboard.writeText(file.path)
    .then(() => {
      selectionOverlay.value = false
      window.$elvenNotify.add('Link copied to clipboard.')
    })
    .catch(() => {
      window.$elvenNotify.add('Copy to clipboard: not have permission.')
    })
}

// convert file path, extension etc
function converter(file: IFile): IFile {
  file.extension = convertExtension(file.extension)
  file.path = convertPreviewPath(file.path)
  return file
}

function playAudio(path: string) {
  window.$elvenPlayer.play(path)
}

function convertDate(date) {
  return Dates.convert(date)
}

function convertSize(size) {
  return Sizes.convert(size)
}

function convertExtension(extension) {
  return Extensions.getReadable(extension)
}

function convertPreviewPath(path) {
  return `${import.meta.env.VITE_UPLOADS_URL}/${path}`
}
</script>

<style scoped>
.files__list {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  min-height: 42px;
  gap: 12px;
}

.file {
  box-shadow: 0 0 41px 0 rgba(0, 0, 0, 0.05);
  min-height: 42px;
  border-radius: var(--border-radius);
  background-color: var(--color-level-1);
  cursor: pointer;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  padding-bottom: 12px;
  gap: 8px;
}

.file__item {
  font-size: 0.9rem;
  line-height: 1.5rem;
  margin-top: 8px;
  margin-left: 12px;
  margin-right: 12px;
}

.file__meta {
  display: flex;
  flex-direction: row;
  color: var(--color-text-inactive);
}

.file__main {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.file__name {
  font-size: 1.1rem;
  line-height: 2rem;
  letter-spacing: 0.0099rem;
}

.file__size {
  font-size: 1rem;
}

.overlay__selected {
  width: 100%;
}

.overlay__item {
  height: 64px;
  width: 100%;
  font-size: 1rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.overlay__item:hover {
  background-color: var(--color-hover);
}
</style>