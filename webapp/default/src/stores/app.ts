// Utilities

import { defineStore } from 'pinia'
import appService from '../services/app.service'
import  { ModelContent, ModelNavi, ModelPage } from '../models/navi'


/*
function isObjEmpty (obj: any) {
  return Object.keys(obj).length === 0;
}
*/

function findParents(node: any, id: string): any {
  // If current node name matches the search name, return
  // empty array which is the beginning of our parent result
  if(node.id && node.id === id) {
    return []
  }
  // Otherwise, if this node has a tree field/value, recursively
  // process the nodes in this tree array
  if(Array.isArray(node.pages)) {
    for(var treeNode of node.pages) {
      // Recursively process treeNode. If an array result is
      // returned, then add the treeNode.name to that result
      // and return recursively
      const childResult = findParents(treeNode, id)
      if(Array.isArray(childResult)) {
        return [ treeNode ].concat( childResult );
      }
    }
  }
}



function findCurrent(pPages: any, path: string): any {
  let pageOut = {} as any
  function find(pages: any) {
    for (let page of pages) {
      if (page.link && page.link == path)  {
        pageOut = page      
      }else {
        if (page.pages) {
          find(page.pages)
       } 
      }
    }
    
  }
  find(pPages);

  return pageOut
}

//type Timer = ReturnType<typeof setTimeout>

export const useAppStore = defineStore('app', {



  state: () => ({
    isRequesting: false,
    isLoading: false,
    navi: {} as ModelNavi,
    content: {} as ModelContent,
    pagesIdx: new Map<string, string>(),
    pages: new Map<string, ModelPage>(),
    contents: new Map<string, ModelContent>()
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
      const pg = findCurrent(this.navi.pages, path);
      let parents = findParents(this.navi, pg.id)
      if (parents && parents.length && parents.length > 0) {
        parents.forEach((p: ModelPage) => {
            p.show = true
        });
      }
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

      if (this.contents.has(id)) {
        const c = this.contents.get(id);
        if (c) {
          this.content = c;
          this.openNavi(path)
          return Promise.resolve();
        }
      }

      return appService.loadJson(this.pagesIdx.get(path) + ".json").then(
        (response: any) => {
          this.content = <ModelContent>response.data
          this.contents.set(this.content.id, this.content)
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
        if (page.id) {
          this.pages.set(page.id, page);
        }
        if (page.link) {
          this.pagesIdx.set(page.link, page.id);
          if (!page.level || page.level < 1) {
            page.show = true
          }
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
