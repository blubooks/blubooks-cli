

import { useAppStore } from "../stores/app";
import api from "./api";


class AppService {


  navi() {
    const appStore = useAppStore()
    appStore.startLoading()
    return api.get("api/navi.json", {
      // query URL without using browser cache
      headers: {
        'Cache-Control': 'no-cache',
        'Pragma': 'no-cache',
        'Expires': '0',
      },
    }).finally(() => {
      appStore.endLoad()
    });
  }

  loadNaviBinary() {
    const appStore = useAppStore()
    appStore.startLoading()
    return api.get("api/data.bin", {
      // query URL without using browser cache
      responseType: 'arraybuffer',
      headers: {
        'Cache-Control': 'no-cache',
        'Pragma': 'no-cache',
        'Expires': '0',
      },
    }).finally(() => {
      appStore.endLoad()
    });
  }

  loadBinary(filename: string) {
    const appStore = useAppStore()
    appStore.startLoading()
    return api.get("api/" + filename, {
      responseType: 'arraybuffer',
    }).finally(() => {
      appStore.endLoad()
    });
  }

  loadJson(filename: string) {

  const appStore = useAppStore()
  appStore.startLoading()
  return api.get("api/" + filename).finally(() => {
    appStore.endLoad()
  });
}

}

export default new AppService();
