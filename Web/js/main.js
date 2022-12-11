const message = {
    data() {
        return {
            messageName: "stReqTest",
            cmdNumber: "0",
            paraNumber: "0",
            searchState: 0,  // 1:searching 2:success 3:fail
            searchBtnTip: "查找消息",
            params: []
        }
    },
    methods: {
        searchMessage() {
            if (this.searchState === 1)
                return
            if (this.messageName === "" || this.messageName.length < 0) {
                alert("请输入有效消息名")
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

                    this.searchBtnTip = "查找消息"
                    if (response.data.cmdFlag === undefined) {
                        this.searchState = 3
                    } else {
                        this.searchState = 2
                        this.cmdNumber = response.data.cmdNumber
                        this.paraNumber = response.data.paraNumber
                        this.params = response.data.params
                        if (this.params != null) {
                            this.params.forEach((item) => {
                                console.log(item)
                                let arrayLen = 0
                                if (item[2].includes("MAX_USERNAME")) {
                                    arrayLen = 30
                                } else {

                                }
                                if (item[2].includes("+") && item[2].includes("1")) {
                                    arrayLen += 1
                                }
                                if (arrayLen) {
                                    item[2] = arrayLen
                                }
                            })
                        }
                    }
                    //todo 解析返回值 
                    console.log(response.data)
                })
                .catch(function (error) {
                    this.searchBtnTip = "查找消息"
                    this.searchState = 3
                    console.log(error);
                })
        },
        sendMessage() {
            if (this.searchState !== 2) {
                return
            }
            if (vServer.selectServer === "") {
                alert("请选择服务器")
                return
            }
            if (vUser.selectUser === "") {
                alert("请选择在线角色")
                return
            }
            console.log(this.params)
            axios
                .post('/message/send', {
                    "serverName": vServer.selectServer,
                    "user": vUser.selectUser,
                    "message": this.messageName,
                    "params": this.params
                })
                .then(response => {
                    if (response.data.status === "NoServer") {
                        alert("请刷新网页重新发送")
                    } else {
                        alert("发送成功")
                    }
                })
                .catch(function (error) {
                    alert("请刷新网页重新发送")
                })
        }
    }
}

const serverList = {
    data() {
        return {
            servers: [],
            selectServer: "4546",
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
        this.timer = setInterval(this.getServerList, 30000)
    }
}

const userList = {
    data() {
        return {
            users: [],
            selectUser: "1001",
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
        this.timer = setInterval(this.getUserList, 30000)
    }
}

vServer = Vue.createApp(serverList).mount('#serverList')
vUser = Vue.createApp(userList).mount('#userList')
vMessage = Vue.createApp(message).mount('#message')


