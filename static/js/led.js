ledSwitch = document.getElementById('led1')

ledSwitch.onclick = function () {
    if (ledSwitch.checked) {
        mysocket.send("led", "on")
    } else {
        mysocket.send("led", "off")
    }
}


displayText = document.getElementById('display-text1')
displayTextSend = document.getElementById('display-text1-send')

displayTextSend.onclick = function () {
    mysocket.send("display", displayText.value)
}

