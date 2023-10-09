<template>
    <SearchModal v-if="appStore.navi.searchId"  @navi="navi"></SearchModal>
    <header id="bl-header" class="bl-header">
        <div class="bl-container"> 
            <div class="bl-inner">
                <HeaderNav  @navi="navi"></HeaderNav>
            </div>
            <div class="clear"></div>
        </div>

    </header>   
    <div id="bl-view" class="bl-view bl-container">
        <nav id="bl-nav" class="bl-nav">
            <div class="bl-nav-inner">
                <Navi @navi="navi"></Navi>
            </div>
        </nav>
        <div id="bl-page" class="bl-page">
            <div class="bl-page-inner">
                <div id="bl-content" class="bl-content">
                    <div class="bl-content-inner">
                        <div id="bl-content-body" class="markdown-body" v-if="appStore.content" v-html="appStore.content.html" />
                    </div>                        
                </div>            
        
                <aside id="bl-sidebar" class="bl-sidebar" >
                    <div class="bl-inner">
                        <Toc v-if="appStore.content.toc" :items="appStore.content.toc"  @scrolling="scrolling" />
                    </div>
                </aside>
            </div>
        </div>
    </div>   
    
</template>
<script lang="ts" setup>
import {onMounted } from 'vue'
import { useAppStore } from "../stores/app";
import Navi from '../components/Navi.vue'
import Toc from '../components/Toc.vue'
import SearchModal from '../components/SearchModal.vue'
import HeaderNav from '../components/HeaderNav.vue'
import { onBeforeRouteUpdate, useRoute } from 'vue-router'
import { ModelPage } from '../models/navi';
import  router  from '../router';


const appStore = useAppStore()
const route = useRoute();
function scrolling(id: string) {
    const el = document.getElementById(id);
    if (el) {
        var elementPosition = el.offsetTop;
        window.scrollTo({
        top: elementPosition - 10, //add your necessary value
        behavior: "smooth"  //Smooth transition to roll
        });

        //const y = el.getBoundingClientRect().top + 80 ;
        //window.scrollTo({top: y, behavior: 'smooth'});
        //el.scrollIntoView({behavior: "smooth", block: "start"});
        //el.scrollTo({behavior: "smooth", block: "start"});

    }
}


function navi(page: ModelPage) {

    if (page.type == "link-extern") {
        if (page.link) {
            window.open(page.link, '_blank');            
         }
    }

    if (page.type == "link") {
        if (page.link) {
            if (route.path != page.link) {
                router.push(page.link)
                return
            }       
         }
    }

    if (page.type == "chapter") {
        let chapter = appStore.pages.get(page.id)
        if (chapter) {
            if (chapter.pages && chapter.pages.length > 0) {
                
                let p = chapter.pages[0]
                if (p && p.link && route.path != p.link) {
                    router.push(p.link)
                    return
                }
            }

        }
    }

    if (page.show) {
        page.show = false
    }else {
        page.show = true
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

<style lang="scss" scoped>



</style>