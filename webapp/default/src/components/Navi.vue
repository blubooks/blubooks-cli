<script setup lang="ts">
import { useAppStore } from "../stores/app";
import  NaviItem  from "./NaviItem.vue";
import { ModelPage } from '../models/navi'
import  router  from '../router'
import { useRoute } from 'vue-router'

const appStore = useAppStore()
const route = useRoute()


function navi(page: ModelPage) {
    if (page.link) {
        if (route.path != page.link) {
            router.push(page.link)
            return
        }
    }
    if (page.show) {
        page.show = false
    }else {
        page.show = true
    }
}

var base = ""

</script>

<template>

    <NaviItem 
        :pages="appStore.navi.pages"
        :level="0"
        :base="base"
        @navi="navi"
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


