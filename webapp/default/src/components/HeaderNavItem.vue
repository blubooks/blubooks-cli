<script setup lang="ts">
import { Page } from '../models/content'
import type { PropType } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from "../stores/app";

const route = useRoute();
const appStore = useAppStore()
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
    
    <ul :class="{ 'bl-group': page && !page.link, 'bl-link-list': page && page.link }">
        <template v-for="pg of pages" :key="pg.link">
            <li :class="{ 'bl-group-item': pg && !pg.link, 'bl-link-list-item': pg && pg.link }">
                <div class="bl-inner-item" :class="{ 'link': pg.link }">
                    <a v-if="pg.link" href="#" :class="{  'active-item': pg.link == route.path, 'active-parent-item': pg.activeParent  }" :to="pg.link"
                        @click.prevent="$emit('navi', pg)">
                        {{ pg.title }}
                    </a>
                    <span class="title" v-else  :class="{ 'active-item':appStore.currentBook.id == pg.id }"   @click="$emit('navi', pg)">{{ pg.title }}</span>
                </div>
                <HeaderNavItem v-if="pg.pages" :pages="pg.pages" :page="pg" :level="level + 1" @navi="navi" />
            </li>
        </template>
    </ul>
</template>


