<script setup lang="ts">
import { useAppStore } from "../stores/app";
import  NaviItem  from "./NaviItem.vue";

const appStore = useAppStore()

var base = ""

</script>

<template>

    <NaviItem 
        :pages="appStore.navi.pages"
        :level="0"
        :base="base"
        >
    </NaviItem>

    <!--
    <div v-if="appStore.navi.title" v-for="item in appStore.navi.pages">
        <router-link :to="'/public/' + item.link">
            {{  item.title }}
        </router-link>
    </div>
-->
</template>


