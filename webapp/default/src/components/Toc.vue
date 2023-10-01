<script setup lang="ts">
import {  ModelToc } from '../models/navi'


const emit = defineEmits(['scrolling'])

function scrolling(id: string) {
  emit('scrolling', id)
}

defineProps({
    items: {
        type: Array<ModelToc>,
        required: true
    },  
});



</script>

<template>
    <ul class="bl-nav-ul">
    <template v-for="item of items" :key="item.id + '_' + index">
        <li>
           <a @click.prevent="$emit('scrolling', item.id)">{{ item.title }}</a>
            <Toc
                @scrolling="scrolling"
                v-if="item.items"
                :items="item.items"
            />
        </li>

    </template>

    </ul>
</template>


