<template>
    <div class="action" :class="{ sale: theme === 'sale' }">
        <b>{{ text }}</b>
        {{ actionurl }}
    </div>
</template>

<script>
export default {
    name: 'Action',
    props: [
        'text', 'actionurl', 'theme'
    ],
    methods: {
        closeModal() {
            this.$emit('close')
        }
    }
}
</script>

<style>
.action {
    width: 100%;
    height: 100%;
    padding: 20px;
    margin: 10px auto;
    background: white;
    border-radius: 10px;
}
.action h1 {
    color: #03cfb4;
    border: none;
    padding: 0;
}
.action p {
    font-style: normal;
}
.action.sale {
    background: crimson;
    color: white;
}
.action.sale h1 {
    background: crimson;
    color: black;
}

</style>