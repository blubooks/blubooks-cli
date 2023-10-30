<template>
  <div>  
    <router-view />
  </div>
</template>

<script lang="ts" setup>
/*
import {  watch} from 'vue'
import { useAppStore } from "./stores/app";
import { useRoute } from 'vue-router'
import router from '../router';


const appStore = useAppStore()
const route = useRoute();


function loadData(path: string) {

appStore.loadContent(path).then()

}

watch(
      () => route.params,
      () => {
        loadData(route.path)
      },
      // fetch the data when the view is created and the data is
      // already being observed
      { immediate: true }
    )

*/
/*
import { onBeforeRouteUpdate, onBeforeRouteLeave } from 'vue-router'


const appStore = useAppStore()

function loadData(path: string) {

appStore.loadContent(path).then()

}

onBeforeRouteLeave((to, from) => {
    if (to.path !== from.path) {
        console.log(to.params)
        loadData(to.path)
    }
})

onBeforeRouteUpdate((to, from) => {
    if (to.path !== from.path) {
        console.log(to.params)
        console.log("test")
        loadData(to.path)
    }
})
*/

</script>