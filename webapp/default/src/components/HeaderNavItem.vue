<script setup lang="ts">
//import { useAppStore } from "../stores/app";
import { ModelPage } from '../models/navi'
import type { PropType } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute();

const emit = defineEmits(['navi'])

function navi(page: ModelPage) {
  emit('navi', page)
}
//const appStore = useAppStore()



defineProps({
    pages: {
        type: Array<ModelPage>,
        required: true
    },
    page: {
        type:   Object as PropType<ModelPage>,
        required: false,
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
    <ul class="bl-nav-ul"  :class="{'bl-group': page && !page.link, 'bl-link-list': page && page.link  }">
    <template v-for="pg of pages" :key="pg.link">
        <li :class="{'bl-group-item': pg && !pg.link, 'bl-link-list-item': pg && pg.link  }">
            <div class="bl-inner-item"  :class="{'link': pg.link }">
                
                <a v-if="pg.link" href="#" :class="{'active-item': pg.link == route.path }"  :to="pg.link" @click.prevent="$emit('navi', pg)">
            {{  pg.title }}
                </a>
            <span class="title" v-else  @click="$emit('navi', pg)">{{ pg.title }}</span>
            </div>


                <HeaderNavItem
                    v-if="pg.pages && pg.show "
                    :pages="pg.pages"
                    :page="pg"
                    :level="level + 1"
                    :base="base"
                    @navi="navi"
                 
                />
        </li>

    </template>

    </ul>
</template>


