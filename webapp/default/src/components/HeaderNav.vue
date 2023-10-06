<script setup lang="ts">
import { useAppStore } from "../stores/app";
import  HeaderNavItem  from "./HeaderNavItem.vue";
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
<p>test</p>
    <HeaderNavItem 
        :pages="appStore.navi.header"
        :level="0"
        :base="base"
        @navi="navi"
        >
    </HeaderNavItem>

    <!--
    <div v-if="appStore.navi.title" v-for="item in appStore.navi.pages">
        <router-link :to="'/public/' + item.link">
            {{  item.title }}
        </router-link>
    </div>
-->
</template>


