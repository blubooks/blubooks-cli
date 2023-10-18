<script setup lang="ts">
//import { useAppStore } from "../stores/app";
import { Page } from '../models/content'
import type { PropType } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute();

const emit = defineEmits(['navi'])

function navi(page: Page) {
    emit('navi', page)
}

defineProps({
    pages: {
        type: Array<Page>,
        required: true
    },
    page: {
        type: Object as PropType<Page>,
        required: false,
    },
    level: {
        type: Number,
        required: false,
        default: 0
    },

});

</script>

<template>
    <ul class="bl-nav-ul" :class="{ 'bl-group': page && page.type == 'group', 'bl-link-list': page && page.link }">
        <template v-for="pg of pages" :key="pg.link">
            <li :class="{ 'bl-group-item': pg && pg.type == 'group', 'bl-link-list-item': pg && pg.link }">
                <div class="bl-inner-item" :class="{ 'link': pg.link }">                 
                    <a v-if="pg.link" href="#"
                        :class="{ 'active-item': pg.link == route.path, 'active-parent-item': pg.activeParent }"
                         @click.prevent="navi(pg)">
                        {{ pg.title }}
                    </a>
                    <span class="title" v-else :class="{ 'bl-group-text': pg && pg.type == 'group', 'active-parent-item': pg.activeParent }" @click="navi(pg)">{{
                        pg.title }}</span>
                </div>
                <NaviItem v-if="pg.pages && pg.show" :pages="pg.pages" :page="pg" :level="level + 1" @navi="navi" />
            </li>
        </template>
    </ul>
</template>


