<template>
    <header class="bl-header">
        <div class="bl-container"> 
            <HeaderNav></HeaderNav>
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
                            <div class="markdown-body" v-if="appStore.content" v-html="appStore.content.html" />
                        </div>                        
                    </div>            
            
                    <aside id="bl-sidebar" >
                        <div v-if="appStore.content.toc">
                            <Toc :items="appStore.content.toc"  @scrolling="scrolling" />
                        </div>
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
import HeaderNav from '../components/HeaderNav.vue'
import { onBeforeRouteUpdate, useRoute } from 'vue-router'
import Toc from '../components/Toc.vue';


const appStore = useAppStore()
const route = useRoute();
function scrolling(id: string) {
    const el = document.getElementById(id);
    if (el) {
        el.scrollIntoView({behavior: "smooth"});

    }
}

/*
onBeforeRouteLeave((to, from) => {
    if (to.path !== from.path) {
        loadData(to.path)
        }
    })
*/

onBeforeRouteUpdate( (to, from) => {
    if (to.path !== from.path) {
        console.log(to.params)
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