// Utilities

import { defineStore } from 'pinia'
import appService from '../services/app.service'
import  { ModelContent, ModelNavi, ModelPage } from '../models/navi'


/*
function isObjEmpty (obj: any) {
  return Object.keys(obj).length === 0;
}
*/



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
          return Promise.resolve();
        }
      }

      return appService.loadJson(this.pagesIdx.get(path) + ".json").then(
        (response: any) => {
          this.content = <ModelContent>response.data
          this.pages.set(this.content.id, this.content)
       
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
          this.prepareNavi(navi.pages)
          this.navi = navi
       
        },
        (err: any) => {
          return Promise.reject(err);
        }
      )
    },
    prepareNavi(pages: Array<ModelPage>) {
      pages.forEach(page => {
        
        if (page.link) {
         // page.link = '/' + page.link
          this.pagesIdx.set(page.link, page.id);
        }

        if (page.pages) {
          
          this.prepareNavi(page.pages)
        }
      });
    }
  }

})
