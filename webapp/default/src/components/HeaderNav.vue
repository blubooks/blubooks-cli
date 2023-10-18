<script setup lang="ts">
import { useAppStore } from "../stores/app";
import  HeaderNavItem  from "./HeaderNavItem.vue";
import { Page } from '../models/content'

const appStore = useAppStore()
const emit = defineEmits(['navi'])

function navi(page: Page) {
  emit('navi', page)
}

</script>

<template>
    <div id="bl-header-nav" class="bl-header-nav">
        <button @click="appStore.openSearch()">Search</button>
        <HeaderNavItem 
            :pages="appStore.navi.header"
            :level="0"
            @navi="navi"
            >
        </HeaderNavItem>
    </div>
</template>


