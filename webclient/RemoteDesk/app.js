const app = Vue.createApp({
   data() {
       return {
           showBooks: true,
           title: 'Buch',
           author: 'Brandon',
           age: 25,
           x: 0,
           y : 0
       }
   },
   methods: {
       changeTitle(title) {
           console.log('div click')
           this.title = title
       },
       toggleShowBooks() {
           this.showBooks = !this.showBooks
       },
       handleEvent(e, data) {
           console.log(e.type, e)
           if (data) {
               console.log(data)
           }
       },
       handleMouseMove(e) {
        this.x = e.offsetX
        this.y = e.offsetY
        }
   }
})

app.mount('#app')