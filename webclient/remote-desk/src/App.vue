<template>
 <h1>{{ title }}</h1>
 <p>Welcome</p>
 <div v-if="showModal">
 <Modal :header="header" :text="text" theme="sale" @close="toggleModal"/>
 </div>
 <br/>
 <button @click="toggleModal">show Modal</button>
</template>

<script>
import Modal from './components/Modal.vue'

export default {
  name: 'App',
  components: {
    Modal
  },
  data() {
    return {
      title: 'my first vue app',
      header: "Title me",
      text: 'this is a text',
      showModal: false
    }
  },
  methods: {
    toggleModal() {
      this.showModal = !this.showModal
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
