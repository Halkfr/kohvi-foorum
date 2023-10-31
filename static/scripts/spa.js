if (window.location.hostname === "localhost") {
    window.location.href = "http://127.0.0.1:8080" + window.location.pathname + window.location.search;
}
if (window.location.pathname === '/') window.location.href = '/home';

document.addEventListener('click', e => {
    console.log(e.target)
    const activeElements = ["title", "home-nav", "create-post-nav", "profile-nav", "sign-out-nav", "register-link", "post-submit-btn"];
    if (activeElements.includes(e.target.id)) {
        console.log('hey')
        e.preventDefault();
        route(e);
    }
});

const route = (e) => {
    window.history.pushState({}, '', e.target.href);
    handleLocation();
}

const routers = {
    '/home': 'home.html',
    '/create-post': 'create-post.html',
    '/view-profile': 'profile.html',
    '/sign-in': 'sign-in.html',
    '/sign-up': 'sign-up.html',
    '/view-post': 'view-post.html'
}


let activeSession = null;

async function checkActiveSession() { // checks active session every 10 sec
    if (activeSession === null) {
        const response = await fetch('/api/session-status', {
            method: 'POST',
            credentials: 'include',
        });
        if (response.ok) {
            activeSession = true;
        } else {
            activeSession = false;
        }   
        setTimeout(() => {
            activeSession = null;
        }, 10000);
    }
    
    return activeSession;
}

const handleLocation = async () => {
    const path = window.location.pathname
    const html = await fetch('./static/templates/' + routers[path]).then((data) => data.text()).catch(error => {
        console.error('Error:', error);
    });

    if (activeSession === true) {
        if (path === "/sign-in" || path === "/sign-up") {
            window.history.replaceState({}, '', 'http://127.0.0.1:8080/home',)
            handleLocation()
        } else {
            if (typeof ws === 'undefined' || ws.readyState === WebSocket.CLOSED) {
                startWS()
            };
            if (document.body.querySelector('#main-container') === null) {
                const header = await fetch('./static/templates/header.html').then((data) => data.text());
                const sidepanel = await fetch('./static/templates/sidepanel.html').then((data) => data.text());
                document.body.innerHTML = header + sidepanel;
            }
            let div = document.body.querySelector('#content-container')
            if (div.innerHTML === "" || !div.classList.contains(path)) {
                div.classList.remove(...div.classList)
                div.classList.add(path)
                div.innerHTML = html
            }
        }
    } else if (activeSession === false){
        if (path === "/sign-in") {
            document.body.innerHTML = html;
        } else if (path === "/sign-up") {
            const header = await fetch('./static/templates/header.html').then((data) => data.text());
            document.body.innerHTML = header + html;
        } else {
            window.history.replaceState({}, '', 'http://127.0.0.1:8080/sign-in',)
            handleLocation()
        }
    } else {
        await checkActiveSession()
        handleLocation()
    }
}

window.addEventListener(onpopstate, async function () {
    let path = window.location.pathname
    if (path === "/sign-in" || path === "/sign-in") {
        await fetch('/api/sign-out', {
            method: 'POST',
            credentials: 'include',
        }).then(response => {
            if (response.ok) { console.log("sign-out successfully") }
        })
    }
})

window.onpopstate = handleLocation;
window.route = route;
handleLocation();

function startWS() {
    ws = new WebSocket('ws://127.0.0.1:8080/ws')
    ws.onmessage = (event) => {
        const chatName = document.getElementById("chat-username").innerHTML
        const chatFiller = document.getElementById("chat-scroll-area")

        console.log("Got message!", event.data)
        data = JSON.parse(event.data)

        // handles messages delivery to chat
        if (chatName == data["SenderName"] || chatName == data["RecipientName"]) { // load message to correct chat
            let [date, message, sender] = Array.from({ length: 3 }, () => document.createElement("div"));

            let senderName = document.createTextNode(data["SenderName"])
            let messageContent = document.createTextNode(data["Messages"]["Content"])
            let dateContent = document.createTextNode(data["Messages"]["Timestamp"])

            if (data["Sender"]) {
                sender.classList.add("sender-name")
                message.classList.add("sender")
                date.classList.add("sender-date")
            } else {
                sender.classList.add("recipient-name")
                message.classList.add("recipient")
                date.classList.add("recipient-date")
            }
            sender.appendChild(senderName)
            message.appendChild(messageContent)
            date.appendChild(dateContent)

            chatFiller.appendChild(sender)
            chatFiller.appendChild(message)
            chatFiller.appendChild(date)

            chatFiller.scrollTo(0, chatFiller.scrollHeight)
            window.chatOffset += 1
        }

        // handles notifications

        const senderId = data.Messages.SenderId
        const recipientId = data.Messages.RecipientId

        if (data["Sender"]) { // clears notifications for sender 
            document.getElementById(recipientId).getElementsByClassName("badge")[0].innerHTML = ""
            moveBtnTop(recipientId, senderId)
        } else if (document.getElementById(senderId)) { // adds notifications to recipient
            document.getElementById(senderId).getElementsByClassName("badge")[0].innerHTML = data.CurrentNotificationCount
            moveBtnTop(recipientId, senderId)
        }

        if (data["TotalNotificationCount"] != 0) { // change total count for current user
            document.getElementById("user-list").getElementsByClassName("badge")[0].innerHTML = data["TotalNotificationCount"]
        } else {
            document.getElementById("user-list").getElementsByClassName("badge")[0].innerHTML = ""
        }
    }
}

function moveBtnTop(recipientId, senderId) {
    let parent = document.getElementById("userlist-scroll-area")
    if (document.getElementById(recipientId)) { // for user sending messages
        parent.insertBefore(document.getElementById(recipientId), parent.firstElementChild)
        parent.scrollTo(0, 0)
    } else { // for user receiving messages
        parent.insertBefore(document.getElementById(senderId), parent.firstElementChild)
    }
}