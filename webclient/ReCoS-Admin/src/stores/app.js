export const appStore = {
  servicePort: 9280,
  baseURL: window.location.protocol + "//localhost:9280/api/v1/",
  password: "",
  authheader: {},
  inconlist: [],
  packageVersion: "0",

  setPassword(password) {
    this.password = password;
    this.authheader = new Headers({
      Authorization: `Basic ${btoa(`admin:${password}`)}`,
    });
  },
};
