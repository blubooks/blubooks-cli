<template>
    <header class="bl-header">
        <div class="bl-container"> 
            Head
        </div>

    </header>   
    <div class="bl-container"> 
        <div id="bl-view">

            <nav id="bl-nav">
                <div class="bl-nav-inner">
                    <Navi></Navi>
                </div>
            </nav>
            <div id="bl-page">
                <div class="bl-page-inner">
                    <div id="bl-content">
                        <div class="bl-content-inner">
                            <div class="markdown-body" v-if="appStore.content" v-html="appStore.content.html">
                            </div>
                        
                        </div>
                        
                    </div>            
            
                    <aside id="bl-sidebar">
                        Side
                    </aside>
                </div>

            </div>
        </div>   
    </div> 
</template>
<script lang="ts" setup>
import {onMounted } from 'vue'
import { useAppStore } from "../stores/app";
import Navi from '../components/Navi.vue'
import { onBeforeRouteUpdate, useRoute } from 'vue-router'


const appStore = useAppStore()
const route = useRoute();


/*
onBeforeRouteLeave((to, from) => {
    if (to.path !== from.path) {
        loadData(to.path)
        }
    })
*/

onBeforeRouteUpdate( (to, from) => {
    if (to.path !== from.path) {
        loadData(to.path)
        }
    })

function loadData(path: string) {

    appStore.loadContent(path).then() 

}

onMounted(() => {
    console.log("moun")
  appStore.loadNavi().then(() => {
    appStore.loadContent(route.path).then() 

  }) 

});



</script>