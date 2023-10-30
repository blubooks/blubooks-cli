<template>
    <SearchModal v-if="appStore.navi.searchId" @navi="navi"></SearchModal>
    <header id="bl-header" class="bl-header">
        <div class="bl-container">
            <div class="bl-inner">
                <HeaderNav @navi="navi"></HeaderNav>
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
                        <div id="bl-content-body" class="markdown-body" v-if="appStore.currentContent"
                            v-html="appStore.currentContent.html" />
                    </div>
                </div>

                <aside id="bl-sidebar" class="bl-sidebar">
                    <div class="bl-inner">

                        <a  @click="pdf">
                            PDF
                        </a>
                        <Toc v-if="appStore.currentContent.toc" :items="appStore.currentContent.toc" @scrolling="scrolling" />
                    </div>
                </aside>
            </div>
        </div>
    </div>
</template>
<script lang="ts" setup>
import { onMounted} from 'vue'
import { useAppStore } from "../stores/app";
import Navi from '../components/Navi.vue'
import Toc from '../components/Toc.vue'
import SearchModal from '../components/SearchModal.vue'
import HeaderNav from '../components/HeaderNav.vue'

import { useRoute } from 'vue-router'
import { Page } from '../models/content';
import router from '../router';
import { onBeforeRouteUpdate} from 'vue-router'


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

function pdf() {
    console.log(appStore.currentContent)
    console.log(appStore.currentPage)

    window.open("#/_doc" + appStore.currentPage.link, '_blank');

}

function navi(page: Page) {

    if (page.type == "link-extern") {
        if (page.link) {
            window.open(page.link, '_blank');
        }
    }

    if (page.type == "link") {
        if (page.link) {
            if (route.path != page.link) {
                router.push(page.link)
                console.log("push")

                return
            }
        }
    }

    if (page.type == "book") {
        let book = appStore.pages.get(page.id)
        if (book) {
            if (book.pages && book.pages.length > 0) {

                let p = book.pages[0]
                if (p && p.link && route.path != p.link) {
                    console.log("push")
                    router.push(p.link)
                    return
                }
            }

        }
    }

    if (page.show) {
        page.show = false
    } else {
        page.show = true
    }
}



function loadData(path: string) {

appStore.loadContent(path).then()

}
/*
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

onBeforeRouteUpdate((to, from) => {
    if (to.path !== from.path) {
        console.log(to.params)
        loadData(to.path)
    }
})



onMounted(() => {
    console.log("moun home")

    appStore.loadBasicNavi().then(() => {
        appStore.loadContent(route.path).then()
    })
   
});

</script>

