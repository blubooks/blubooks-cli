<script setup lang="ts">
import { useAppStore } from "../stores/app";

const appStore = useAppStore()

</script>

<template>
    <div v-if="appStore.navi.title" v-for="item in appStore.navi.pages">
        <router-link :to="'/public/' + item.link">
            {{  item.title }}
        </router-link>
    </div>
</template>


