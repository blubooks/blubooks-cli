<script setup lang="ts">
//import { useAppStore } from "../stores/app";
import { ModelPage } from '../models/navi'
import type { PropType } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute();

//const appStore = useAppStore()



defineProps({
    pages: {
        type: Array<ModelPage>,
        required: true
    },
    page1: {
        type:   Object as PropType<ModelPage>,
        required: false,
        default: []
    },
    level: {
        type: Number,
        required: false,
        default: 0
    },
    base: {
        type: String,
        required: true,
    }    
});



</script>

<template>
    <ul class="bl-nav-ul">
    <template v-for="page of pages" :key="page.link">
        <li>
            <div class="name">
                
                <router-link v-if="page.link" :to="page.link">
            {{  page.title }} (<span v-if="page.link == route.path">{{ route.path }}</span>)
            </router-link>
            <span v-else>{{ page.title }}</span>
            </div>


                <NaviItem
                    v-if="page.pages"
                    :pages="page.pages"
                    :page1="page"
                    :level="level + 1"
                    :base="base"
                 
                />
        </li>

    </template>

    </ul>
</template>


