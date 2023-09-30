<template>
    <div class="bl-wrapper">
        <Navi></Navi>
        <div>
            <p>hallo</p>
        </div>
    </div>
</template>
<script lang="ts" setup>
import { useRoute } from 'vue-router'
useRoute();
import {onMounted } from 'vue'
import { useAppStore } from "../stores/app";
import Navi from '../components/Navi.vue'

const appStore = useAppStore()

onMounted(() => {
  appStore.loadNavi().then() 
});



</script>