const socketAPI = 'ws://localhost:8080/ws-web'

class MySocket{
    constructor(){
        this.mysocket =  null;
        this.vMsgContainer = document.getElementById("msgcontainer");
        this.vMsgIpt = document.getElementById("ipt");
    }

    showMessage(text, myself){
        var div = document.createElement("div");
        div.innerHTML = text;
        var cself = (myself)? "self" : "";
        div.className="msg " + cself;
        this.vMsgContainer.appendChild(div);
    }

    send(messageType, data){
        const message = JSON.stringify({
            "type": messageType,
            "data": data,
        });
        this.mysocket.send(message);
        console.log(message);
    }

    connectSocket(){
        console.log("memulai socket");
        var socket = new WebSocket(socketAPI);
        this.mysocket = socket;

        socket.addEventListener("message", (event) => {
            // this.showMessage(event.data,false);
            // var mess = JSON.parse(event.data)
            // if (mess.source === "process") {
            //     switch (mess.status) {
            //         case "success":
            //             processFinish(mess)
            //             break;
            //         case "info":
            //             processInfo(mess)
            //             break;
            //         case "error":
            //             processInfo(mess)
            //             showProcessPause()
            //             break;
            //     }
            // }
            // alert(event.data)
            document.getElementsByTagName("body")[0].style.backgroundColor=event.data;
            // document.getElementsByTagName('body')[0].innerHTML="";
        });

        socket.onopen = ()=> {
            console.log("socket opend")
        };
        socket.onclose = ()=>{
            console.log("socket close")
        }
    }
}

var mysocket = new MySocket()
mysocket.connectSocket();