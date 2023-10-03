// Utilities

import { defineStore } from 'pinia'
import appService from '../services/app.service'
import  { ModelContent, ModelNavi, ModelPage } from '../models/navi'


/*
function isObjEmpty (obj: any) {
  return Object.keys(obj).length === 0;
}
*/

function find() {
  console.log("TEST")
}

//type Timer = ReturnType<typeof setTimeout>

export const useAppStore = defineStore('app', {



  state: () => ({
    isRequesting: false,
    isLoading: false,
    navi: {} as ModelNavi,
    content: {} as ModelContent,
    pagesIdx: new Map<string, string>(),
    pages: new Map<string, ModelContent>()
  }),
  getters: {

  }, 
  actions: {
    
    startLoading() {
      this.isRequesting = true  
      this.isLoading = true;
    },    
    endLoad() {
      this.isRequesting = false;
      this.isLoading = false;
    },  
    openNavi(path: string) {
      function find(pages: any) {
        pages.forEach((page: any) => {
          if (page.link && page.link == path)  {
             page.show = true
             return
          }
          if (page.pages) {
             find(page.pages)
          }
        });
      }
      find(this.navi.pages)
    },
    closeNavi(path: string) {
      function find(pages: any) {
        pages.forEach((page: any) => {
          if (page.link && page.link == path)  {
             page.show = false
             return
          }
          if (page.pages) {
             find(page.pages)
          }
        });
      }
      find(this.navi.pages)
    },    
    loadContent(path: string) {
      if (!this.pagesIdx.has(path)) {
        return Promise.reject();
      }
      const id = this.pagesIdx.get(path)

      if (!id) {
        return Promise.reject();
      }

      if (this.pages.has(id)) {
        const c = this.pages.get(id);
        if (c) {
          this.content = c;
          this.openNavi(path)
          return Promise.resolve();
        }
      }

      return appService.loadJson(this.pagesIdx.get(path) + ".json").then(
        (response: any) => {
          this.content = <ModelContent>response.data
          this.pages.set(this.content.id, this.content)
          this.openNavi(path)

        },
        (err: any) => {
          return Promise.reject(err);
        }
      )
    }, 
    loadNavi(){

      return appService.navi().then(
        (response: any) => {
          
          let navi = <ModelNavi>response.data;
          this.pagesIdx.set("/", navi.id)
          this.prepareNavi(navi.pages, 0, 1)
          this.navi = navi
       
        },
        (err: any) => {
          return Promise.reject(err);
        }
      )
    },
    prepareNavi(pages: Array<ModelPage>, level: number, showLevel: number) {
      level = level +1;
      pages.forEach(page => {
        
        if (page.link) {
         // page.link = '/' + page.link
      if (!page.level || page.level < 1) {
        page.show = true
       }
          this.pagesIdx.set(page.link, page.id);
        }else {
          page.show = true
        }

        if (page.pages) {
          
          this.prepareNavi(page.pages, level, showLevel)
        }
      });
    }
  }

})
