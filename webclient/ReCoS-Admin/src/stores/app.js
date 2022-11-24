import { defineStore } from "pinia";

export const useAppStore = defineStore("app", {
  state() {
    return {
      count: 0,
      servicePort: 9280,
      baseURL: window.location.protocol + "//localhost:9280/api/v1/",
      password: "",
      authheader: {},
      inconlist: [],
      packageVersion: "0",
    };
  },
  actions: {
    increment(state) {
      state.count++;
    },
    baseURL(state, baseurl) {
      state.baseURL = baseurl;
      if (!baseurl.endsWith("/")) {
        state.baseURL = state.baseURL + "/";
      }
    },
    password(state, password) {
      state.password = password;
      state.authheader = {
        Authorization: `Basic ${btoa(`admin:${password}`)}`,
      };
    },
    iconlist(state, iconlist) {
      state.iconlist = iconlist;
    },
  },
});
