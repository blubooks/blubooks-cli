

import { useAppStore } from "../stores/app";
import api from "./api";


class AppService {


  navi() {
    const appStore = useAppStore()
    appStore.startLoading()
    return api.get("api/navi", {
      responseType: 'arraybuffer',
      // query URL without using browser cache
      headers: {
        'Cache-Control': 'no-cache',
        'Pragma': 'no-cache',
        'Expires': '0',
      },
    }).finally(() => {
      appStore.endLoading()
    });
  }


  loadBinary(filename: string) {
    const appStore = useAppStore()
    appStore.startLoading()
    return api.get("api/" + filename, {
      responseType: 'arraybuffer',
    }).finally(() => {
      appStore.endLoading()
    });
  }

  loadJson(filename: string) {

  const appStore = useAppStore()
  appStore.startLoading()
  return api.get("api/" + filename).finally(() => {
    appStore.endLoading()
  });
}

}

export default new AppService();
