<script setup lang="ts">
import {  PageTocItem } from '../models/content'


const emit = defineEmits(['scrolling'])
import { useRoute } from 'vue-router'

const route = useRoute();



function scrolling(id: string) {
  emit('scrolling', id)
}

defineProps({
    items: {
        type: Array<PageTocItem>,
        required: true
    },  
});



</script>

<template>
    <ul class="bl-toc-ul">
    <template v-for="item of items" :key="item.id + '_' + index">
        <li>
           <a :href="route.path  +'#' + item.id" @click.prevent="$emit('scrolling', item.id)">{{ item.title }}</a>
            <Toc
                @scrolling="scrolling"
                v-if="item.items"
                :items="item.items"
            />
        </li>

    </template>

    </ul>
</template>


