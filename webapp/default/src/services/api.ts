import axios from "axios";


const path = "/public"

const instance = axios.create({
  //baseURL: "http://dev.cidb.de:4061/api/v1",
  //baseURL: "http://192.168.80.19:4061/api/v1",
  baseURL: path,
  headers: {
    "Content-Type": "application/json",
  },
});




export default instance;


