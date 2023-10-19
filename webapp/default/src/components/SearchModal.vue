<template>
    <Transition name="modal">
        <modal id="bl-search" class="modal-wrapper" @click="appStore.closeSearch" v-if="appStore.searchOpened"
            aria-modal="true" role="dialog" tabindex="-1">
            <div class="modal-mask">
                <div class="modal-wrapper">
                    <div class="modal-container" @click.stop="">
                        <div class="modal-header">
                            <input class="search-input" v-model="searchText" @input="doSearch" type="text">
                        </div>

                        <div class="modal-body">
                            <div class="bl-result">
                                <div class="bl-result-item" v-for="item in result" @click="navi(item)">
                                    <div class="bl-result-title"> {{ item.title }}</div>
                                    <div class="bl-result-description"> {{ (item.text) ? item.text : '&nbsp;' }}</div>
                                    <div class="bl-result-path"> {{ item.path }}</div>
                                </div>
                            </div>
                        </div>

                        <div class="modal-footer">
                            <!--                            &nbsp;
                            <button class="modal-default-button" @click.stop="appStore.closeSearch">
                                OK
                            </button>
                            -->

                        </div>
                    </div>
                </div>
            </div>
        </modal>
    </Transition>
</template>
  
<script setup lang="ts">
import { SearchPage } from "../models/content";
import { useAppStore } from "../stores/app";
import { onMounted, ref } from 'vue'
import { debounce, copyObject } from '../utils/helper'

const emit = defineEmits(['navi'])
const result = ref([] as Array<SearchPage>)
const searchText = ref("")
const debounceFn = ref()
const appStore = useAppStore()

function navi(p: SearchPage) {
    const v = appStore.pages.get(p.id);
    if (v && v.link) {
        appStore.searchOpened = false
        emit('navi', v)

    }
}

function isMatch(s: string, p: string): boolean {
    let sIdx = 0, pIdx = 0, lastWildcardIdx = -1, sBacktrackIdx = -1, nextToWildcardIdx = -1;

    while (sIdx < s.length) {
        if (pIdx < p.length && (p[pIdx] === '?' || p[pIdx] === s[sIdx])) {
            // Characters match
            sIdx++;
            pIdx++;
        } else if (pIdx < p.length && p[pIdx] === '*') {
            // Wildcard, so characters match - store the index.
            lastWildcardIdx = pIdx;
            nextToWildcardIdx = ++pIdx;
            sBacktrackIdx = sIdx;
        } else if (lastWildcardIdx === -1) {
            // No match, and no wildcard has been found.
            return false;
        } else {
            // Backtrack - no match, but a previous wildcard was found.
            pIdx = nextToWildcardIdx;
            sIdx = ++sBacktrackIdx;
        }
    }

    // Check if there are only wildcards left in the pattern.
    for (let i = pIdx; i < p.length; i++) {
        if (p[i] !== '*') return false;
    }

    return true;
}

function doSearch() {
    debounceFn.value()
}

onMounted(() => {
    debounceFn.value = debounce(() => {

        appStore.loadSearch().then(() => {
   
        let pattern = searchText.value.trim()



        if (pattern == "" || pattern == null || pattern.length <= 3) {
            result.value = []
            return
        }

        if (appStore.searchList.pages.length == 0) {
            result.value = []
            return
        }

        let res = [] as Array<SearchPage>
        pattern = pattern.toLowerCase()
        pattern = '*' + pattern.replace(" ", "*").trim() + '*';

        console.log("PATTEN", pattern)

        //@ts-ignore
        for (let i = 0; i < appStore.searchList.pages.length; ++i) {
            const value = appStore.searchList.pages[i]
            let mateched = false
            if (isMatch(value.title.toLowerCase(), pattern)) {
                mateched = true
            } else if (value.text && isMatch(value.text.toLowerCase(), pattern)) {
                mateched = true
            } else {
                /*
                if (appStore.pages.get(value.id)) {
                    const v = appStore.pages.get(value.id);
                    if (v) {
                        if (isMatch(v?.title.toLowerCase(), pattern)) {
                           mateched = true
                        }
                    }
                }
                */
            }
            if (mateched) {
                const obj = copyObject(value)
                obj.title = obj.title.substring(0, 200)
                if (obj.text) {
                    obj.text = obj.text.substring(0, 200)
                }
                res.push(obj)            

            }

        }
        result.value = res
    })


    }, 800)




});

</script>
  
<style scoped>
.bl-result-item {
    border: 1px;
    padding: 4px;
    margin-top: 4px;
    margin-bottom: 4px;
    border-radius: 4px;
    background: #eee;
}
.bl-result-item:hover {
    background: #ccc;
    cursor: pointer;
}
.search-input {
    width: 100%;
}

.modal-mask {
    position: fixed;
    z-index: 9998;
    top: 0;
    left: 0;
    width: 100%;
    height: 100vh;
    background-color: rgba(0, 0, 0, 0.5);
    display: table;
    transition: opacity 0.3s ease;
}

.modal-wrapper {
    display: table-cell;
    vertical-align: top;

}

.modal-container {

    width: 600px;
    margin: 100px auto 10px auto;
    padding: 20px 30px;

    background-color: #fff;
    border-radius: 2px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.33);
    transition: all 0.3s ease;
    font-family: Helvetica, Arial, sans-serif;
}

.modal-header h3 {
    margin-top: 0;
    color: #42b983;
}

.modal-body {
    margin: 20px 0;
    overflow: scroll;
    max-height: 500px;

}

.modal-default-button {
    float: right;

}


.modal-enter-from,
.modal-leave-to {
    opacity: 0;
}

.modal-enter-active .modal-container,
.modal-leave-active .modal-container {
    -webkit-transform: scale(1.1);
    transform: scale(1.1);
}
</style>