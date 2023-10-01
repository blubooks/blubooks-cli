

import { useAppStore } from "../stores/app";
import api from "./api";


class AppService {


  navi() {

    const appStore = useAppStore()
    appStore.startLoading()
    return api.get("api/navi.json").finally(() => {
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
