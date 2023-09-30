// Utilities

import { defineStore } from 'pinia'
import appService from '../services/app.service'


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
    navi: {} as any,

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
    loadNavi(){
      return appService.navi().then(
        (response: any) => {
          
          this.navi = response.data
       
        },
        (err: any) => {
          return Promise.reject(err);
        }
      )
    },
  }

})
