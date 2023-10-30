// Utilities

import { defineStore } from 'pinia'
import appService from '../services/app.service'
import { PageContent, Navi, Page, SearchPage, SearchList } from '../models/content'

/*
function isObjEmpty (obj: any) {
  return Object.keys(obj).length === 0;
}
*/

function findParents(node: any, id: string): any {
  // If current node name matches the search name, return
  // empty array which is the beginning of our parent result
  if (node.id && node.id === id) {
    return []
  }
  // Otherwise, if this node has a tree field/value, recursively
  // process the nodes in this tree array
  if (Array.isArray(node.pages)) {
    for (var treeNode of node.pages) {
      // Recursively process treeNode. If an array result is
      // returned, then add the treeNode.name to that result
      // and return recursively
      const childResult = findParents(treeNode, id)
      if (Array.isArray(childResult)) {
        return [treeNode].concat(childResult);
      }
    }
  }
}



function findCurrent(pPages: Page[], path: string): Page {
  let pageOut = {} as Page
  function find(pages: Page[]) {
    for (let page of pages) {
      if (page.link && page.link == path) {
        pageOut = page
      } else {
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
    navi: {} as Navi,
    currentContent: {} as PageContent,
    currentPage: {} as Page,
    currentParents: [] as Array<Page>,
    currentBook: {} as Page,
    pagesPathId: new Map<string, string>(),
    pages: new Map<string, Page>(),
    books: new Map<string, Page>(),
    contents: new Map<string, PageContent>(),
    searchOpened: false,
    searchResult: [] as Array<SearchPage>,
    searchList: {} as SearchList
  }),
  getters: {

  },
  actions: {
    startLoading() {
      this.isRequesting = true
      this.isLoading = true;
    },
    endLoading() {
      this.isRequesting = false;
      this.isLoading = false;
    },    
    closeSearch() {
      this.searchOpened = false
    },
    openSearch() {
      this.searchOpened = true
    },


    loadContent(path: string) {


      var openNavi = () => {
        const pg = findCurrent(this.navi.pages, path);
        this.currentPage = pg;
        
        if (this.currentParents && this.currentParents.length && this.currentParents.length > 0) {
          this.currentParents.forEach((p: Page) => {
            if (p.link) {
              if (this.navi.options && this.navi.options.accordion) {
                p.show = false
              }
            }
            p.activeParent = false
          });
        }

        this.currentParents = findParents(this.navi, pg.id)

        if (this.currentParents && this.currentParents.length && this.currentParents.length > 0) {
          this.currentParents.forEach((p: Page) => {
            p.show = true
            p.activeParent = true
            if (p.type == "book") {
              this.currentBook = p;
            }
          });
        }
      };


      if (!this.pagesPathId.has(path)) {
        return Promise.reject();
      }
      const id = this.pagesPathId.get(path)

      if (!id) {
        return Promise.reject();
      }

      if (this.contents.has(id)) {
        const c = this.contents.get(id);
        if (c) {
          this.currentContent = c;
          
          openNavi()
          return Promise.resolve();
        }
      }

      return appService.loadBinary(this.pagesPathId.get(path) + "").then(
        (response: any) => {
          let data = new Uint8Array(response.data);
          this.currentContent = PageContent.fromBinary(data)
          this.contents.set(this.currentContent.id, this.currentContent)
          openNavi()

        },
        (err: any) => {
          return Promise.reject(err);
        }
      )
    },
    loadBasicNavi() {
      return appService.navi().then(
        (response: any) => {


          const prepareNavi = (pages: Array<Page>, level: number, showLevel: number) => {
            level = level + 1;
            pages.forEach(page => {
              if (page.id) {
                this.pages.set(page.id, page);
              }
              if (page.link) {
                this.pagesPathId.set(page.link, page.id);
                if (!page.level || page.level < 1) {
                  page.show = true
                }
              } else {
                page.show = true
              }
      
              if (page.pages) {
      
                prepareNavi(page.pages, level, showLevel)
              }
            });
      
          }

          const data = new Uint8Array(response.data);
          const navi = Navi.fromBinary(data)

          this.pagesPathId.set("/", navi.id)

          prepareNavi(navi.pages, 0, 1)

          navi.pages.forEach(book => {
            this.books.set(book.id, book)
          });

          if (navi.footer) {
            prepareNavi(navi.footer, 0, 0)
          }

          this.navi = navi
        },
        (err: any) => {
          return Promise.reject(err);
        }
      )

    },
    loadSearch() {
      if (this.searchList && this.searchList.pages && this.searchList.pages.length && this.searchList.pages.length > 0) {
        return Promise.resolve();
      }
      if (this.navi.searchId) {
        return appService.loadBinary(this.navi.searchId).then(
          (response: any) => {
            let data = new Uint8Array(response.data);
            this.searchList = SearchList.fromBinary(data)

          },
          (err: any) => {
            return Promise.reject(err);
          }
        )
      }
      return Promise.reject();
    },


  }

})
