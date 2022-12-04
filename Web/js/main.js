const serverList = {
    data() {
        return {
            servers: [],
            selectServer:""
        }
    },
    methods: {
        chooseServer(event){
            console.log(event.currentTarget.value)
            this.selectServer = event.currentTarget.value
        },
        getServerList(){
            axios
                .get('/message/search/王宇博/stReqTest')
                .then(function (response) {
                    console.log(response);
                    console.log(response.data)
                    console.log(response.data.cmd)
                })
                .catch(function (error) {
                    console.log(error);
                });
        }
    },
    mounted(){
       this.getServerList()
    }
}

const userList = {
    data() {
        return {
            users: [
                'Google',
                'Runoob',
                'Taobao'
            ]
        }
    }
}

vServerList = Vue.createApp(serverList).mount('#serverList')

Vue.createApp(userList).mount('#userList')
