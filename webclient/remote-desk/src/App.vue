<template>
  <b>{{ title }}</b>
  <label> Profile</label>
  <select v-model="profileName" :disabled="readonly" @change="changeProfile()">
    <option v-for="item in items" :value="item.name" :key="item.name" v-text="item.name" :title="item.description" ></option>
  </select>
  <div>
    <button v-for="page in activeProfile.pages" :value="page.name" :key="page.name" v-text="page.name" :title="page.description" @click="changePage(page.name)"></button>
  </div>
  <div>
  <table border="1">
    <tr v-for="(row, x) in cells" :key="x">
      <td v-for="(col, y) in cells[x]" :key="y">
      <Action :text="cells[x][y]" :actionurl="actionurl" ></Action>
      </td>
    </tr>
  </table>
  {{ activePage.name }}
  </div>
</template>

<script>
import Action from './components/Action.vue'

export default {
  name: 'App',
  components: {
    Action
  },
  data() {
    return {
      baseurl: "https://localhost:9543/api/v1/show",
      actionurl: "https://localhost:9543/api/v1/action",
      title: 'remote commands',
      header: "Title me",
      text: 'this is a text',
      showModal: false,
      readonly: true,
      profileName: "none",
      items: [ {
        "name" : "none",
        "description": "nothing found here",
      }],
      activeProfile: {},
      pageName: "",
      activePage: {},
      cells: [
        []
      ]
    }
  },
  mounted() {
    fetch(this.baseurl)
    .then(res => res.json())
    .then(data => {
      this.items = data.profiles
      this.profileName = data.profiles[0].name
      this.readonly = false
      this.changeProfile()
    })
    .catch(err => console.log(err.message))
  },
  methods: {
    toggleModal() {
      this.showModal = !this.showModal
    },
    changeProfile() {
      fetch(this.baseurl + "/" + this.profileName)
      .then(res => res.json())
      .then(data => {
        this.activeProfile = data
        this.activePage = this.activeProfile.pages[0]
        this.changePage(this.activePage.name)
      })
      .catch(err => console.log(err.message))
    },
    changePage(pageName) {
      this.pageName = pageName
      this.activeProfile.pages.forEach(page => {
        if (pageName == page.name) {
          this.activePage = page
        }
      })
      console.log(this.activePage)
      this.cells = new Array(this.activePage.rows)
      for (let x = 0; x < this.activePage.rows; x++) {
        this.cells[x] = new Array(this.activePage.columns)
        for (let y = 0; y < this.activePage.columns; y++) {
          let action = this.activeProfile.actions[(x*this.activePage.rows) + y]
          if (action) {
            this.cells[x][y] = action.name
          } else {
            this.cells[x][y] = " "
          }
        }
      }
    }
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
h1 {
  border-bottom:  1px solid #ddd;
  display: inline-block;
  padding-bottom: 10px;
}
</style>
