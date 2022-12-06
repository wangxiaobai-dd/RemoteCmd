const message = {
    data() {
        return {
            messageName: "",
            cmdNumber: "178",
            paraNumber: "140",
            searchState: 0,  // 1:searching 2:success 3:fail
            searchBtnTip: "查找消息"
        }
    },
    methods: {
        searchMessage() {
            if (this.searchState === 1)
                return
            if (this.messageName === "" || this.messageName.length < 2) {
                console.log(this.messageName.length)
                alert("请输入消息名")
                return
            }
            if (vServer.selectServer === "") {
                alert("请选择服务器")
                return
            }
            console.log(this.messageName)
            this.searchState = 1
            this.searchBtnTip = "查找中..."

            axios
                .get('/message/search', {
                    params: {"serverName": vServer.selectServer, "message": this.messageName}
                })
                .then(response => {
                    if(response.data == null){
                        this.searchState = 3
                        return
                    }
                    //todo 解析返回值 
                    console.log(response.data)
                })
                .catch(function (error) {
                    this.searchState = 0
                    console.log(error);
                });
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
        this.timer = setInterval(this.getServerList, 20000)
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
