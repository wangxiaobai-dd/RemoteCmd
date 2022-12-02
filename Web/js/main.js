const serverList = {
    data() {
        return {
            servers: [
                { text: 'Google' },
                { text: 'Runoob' },
                { text: 'Taobao' }
            ]
        }
    }
}

Vue.createApp(serverList).mount('#serverList')
