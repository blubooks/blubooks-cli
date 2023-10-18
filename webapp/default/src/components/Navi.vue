<script setup lang="ts">
import { useAppStore } from "../stores/app";
import NaviItem from "./NaviItem.vue";
import { Page } from '../models/content'

const emit = defineEmits(['navi'])

const appStore = useAppStore()

function navi(page: Page) {
    emit('navi', page)
}

</script>

<template>
    <NaviItem :pages="appStore.currentBook.pages" :level="0" @navi="navi">
    </NaviItem>
</template>


