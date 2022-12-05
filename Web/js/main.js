const message = {
    data() {
        return {
            messageName:"",
            cmdNumber:"178",
            paraNumber:"140",
            searchResult:0  // 1:success 2:fail
        }
    },
    methods: {
        searchMessage() {
            console.log(this.messageName)
        }
    }
}

const serverList = {
    data() {
        return {
            servers: [],
            selectServer: "",
            timer: ""
        }
    },
    methods: {
        chooseServer(event) {
            this.selectServer = event.currentTarget.value
            vUser.getUserList()
        },
        getServerList() {
            axios
                .get('/server/showList')
                .then(response => {
                    if (response.data != null)
                        response.data.sort()
                    this.servers = response.data
                    let found = false
                    if (this.servers != null) {
                        this.servers.forEach((item) => {
                            if (item === this.selectServer)
                                found = true
                        })
                    }
                    if (!found)
                        this.selectServer = ""
                })
                .catch(function (error) {
                    console.log(error);
                });
        }
    },
    mounted() {
        this.getServerList()
        this.timer = setInterval(this.getServerList, 2000)
    }
}

const userList = {
    data() {
        return {
            users: [],
            selectUser: "",
            timer: ""
        }
    },
    methods: {
        chooseUser(event) {
            this.selectUser = event.currentTarget.value
        },
        getUserList() {
            if (vServer.selectServer === "") {
                return
            }
            axios
                .get('/user/showList', {
                    params: {"serverName": vServer.selectServer}
                })
                .then(response => {
                    if (response.data != null)
                        response.data.sort()
                    this.users = response.data
                    let found = false
                    if (this.users != null) {
                        this.users.forEach((item) => {
                            if (item === this.selectUser)
                                found = true
                        })
                    }
                    if (!found)
                        this.selectUser = ""
                })
                .catch(function (error) {
                    console.log(error);
                });
        }
    },
    mounted() {
        this.timer = setInterval(this.getUserList, 20000)
    }
}

vServer = Vue.createApp(serverList).mount('#serverList')
vUser = Vue.createApp(userList).mount('#userList')
vMessage = Vue.createApp(message).mount('#message')
