<script setup lang="ts">
import { useAppStore } from "../stores/app";
import  HeaderNavItem  from "./HeaderNavItem.vue";
import { ModelPage } from '../models/navi'

const appStore = useAppStore()
const emit = defineEmits(['navi'])

function navi(page: ModelPage) {
  emit('navi', page)
}

</script>

<template>
    <div id="bl-header-nav">
        <HeaderNavItem 
            :pages="appStore.navi.header"
            :level="0"
            @navi="navi"
            >
        </HeaderNavItem>
    </div>
</template>


