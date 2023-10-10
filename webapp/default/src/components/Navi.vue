<script setup lang="ts">
import { useAppStore } from "../stores/app";
import NaviItem from "./NaviItem.vue";
import { ModelPage } from '../models/navi'

const emit = defineEmits(['navi'])

const appStore = useAppStore()

function navi(page: ModelPage) {
    emit('navi', page)
}

</script>

<template>
    <NaviItem :pages="appStore.currentBook.pages" :level="0" @navi="navi">
    </NaviItem>
</template>


