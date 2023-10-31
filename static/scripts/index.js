window.sidepanelOffset = 0
window.sidepanelLimit = 15
window.chatOffset = 0
window.chatLimit = 10

function debounce(delay, fn) {
    let timerId;
    return function (...args) {
        if (timerId) {
            clearTimeout(timerId);
        }
        timerId = setTimeout(() => {
            fn(...args);
            timerId = null;
        }, delay);
    };
}

function createUserlistElement(id, username, status, notificationCount) {
    let chatButton = createChatButton(id, username, status, notificationCount)

    chatButton.addEventListener("click", () => handleChat(chatButton))
    chatButton.addEventListener("click", () => loadChat(chatButton))
}

async function loadChat(chatButton, insert = "bottom") {
    await fetch('http://127.0.0.1:8080/api/load-chat?chatBtnId=' + chatButton.id + '&offset=' + window.chatOffset + '&limit=' + window.chatLimit, {
        method: 'GET',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'credentials': 'include',
        }
    }).then(response => {
        if (response.ok) {
            response.json().then((data) => {
                if (Object.values(data)[0] !== null) {
                    stylizeIncomingMessages(data, chatButton.id, insert)
                    window.chatOffset += window.chatLimit
                }
            })
        }
        else { /*display error*/ }
    })
}

function stylizeIncomingMessages(data, recipientId, insert) {
    const chatFiller = document.getElementById("chat-scroll-area")
    let messages = Object.values(data)[0]
    for (let i = 0; i < messages.length; i++) {
        let [date, message, sender] = Array.from({ length: 3 }, () => document.createElement("div"));
        let senderName = document.createTextNode(Object.values(data)[1][i])
        let messageContent = document.createTextNode(messages[i]['Content'])
        let dateContent = document.createTextNode(messages[i]['Timestamp'])

        sender.appendChild(senderName)
        message.appendChild(messageContent)
        date.appendChild(dateContent)

        if (messages[i]['SenderId'] == recipientId) {
            sender.classList.add("recipient-name")
            message.classList.add("recipient")
            date.classList.add("recipient-date")
        } else {
            sender.classList.add("sender-name")
            message.classList.add("sender")
            date.classList.add("sender-date")
        }
        chatFiller.insertBefore(date, chatFiller.firstChild);
        chatFiller.insertBefore(message, chatFiller.firstChild);
        chatFiller.insertBefore(sender, chatFiller.firstChild);
    }
    if (insert == "bottom") {
        chatFiller.scroll(0, chatFiller.scrollHeight)
    } else {
        // make smooth transition on loading old messages with scroll
    }
}

function createChatButton(id, username, status, notificationCount) {
    let template = document.getElementById('user-template')
    let userChatBtn = template.cloneNode(true)
    userChatBtn.classList.remove("d-none")

    userChatBtn.querySelector(".chat-with").innerHTML = username
    userChatBtn.querySelector(".user-status").innerHTML = status
    userChatBtn.classList.add("chat-with-user")
    userChatBtn.id = id
    userChatBtn.classList.add("chat-button")
    if (notificationCount === 0) {
        userChatBtn.getElementsByClassName("badge")[0].innerHTML = ""
    } else {
        userChatBtn.getElementsByClassName("badge")[0].innerHTML = notificationCount
    }
    document.getElementById("userlist-scroll-area").appendChild(userChatBtn)
    return userChatBtn
}

function removeActiveBtn() {
    document.getElementById("chat-scroll-area").innerHTML = ""
    let elements = document.getElementsByClassName("chat-button");

    for (var i = 0; i < elements.length; i++) {
        elements[i].classList.remove("active");
    }
}

function handleChat(userChatBtn) { // use to open/close chat
    let chatArea = document.getElementById("chat-area").classList

    window.chatOffset = 0

    if (chatArea.contains("d-none")) {
        userChatBtn.classList.add("active")
        document.getElementById("chat-scroll-area").innerHTML = ""
        document.getElementById("chat-username").innerHTML = userChatBtn.getElementsByClassName("chat-with")[0].innerHTML
        document.getElementById("chat-user-status").innerHTML = userChatBtn.getElementsByClassName("user-status")[0].innerHTML
        chatArea.remove("d-none")
    } else {
        if (document.getElementById("chat-username").innerHTML === userChatBtn.getElementsByClassName("chat-with")[0].innerHTML) {
            chatArea.add("d-none")
            document.getElementById("chat-scroll-area").innerHTML = ""
            userChatBtn.classList.remove("active")
        } else {
            removeActiveBtn()
            userChatBtn.classList.add("active")
            document.getElementById("chat-username").innerHTML = userChatBtn.getElementsByClassName("chat-with")[0].innerHTML
            document.getElementById("chat-user-status").innerHTML = userChatBtn.getElementsByClassName("user-status")[0].innerHTML
        }
    }
}

function fillUserlist() {
    fetch('http://127.0.0.1:8080/api/users?offset=' + window.sidepanelOffset + '&limit=' + window.sidepanelLimit, {
        method: 'GET',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'credentials': 'include',
        }
    }).then(response => {
        if (response.status === 200) {
            response.json().then((data) => {
                for (let i = 0; i < Object.keys(data.Users).length; i++) {
                    createUserlistElement(Object.values(data.Users[i])[0], Object.values(data.Users[i])[2], Object.values(data.Users[i])[8], data.Count[i])
                }
            })
            window.sidepanelOffset += window.sidepanelLimit
        }
        else { /*display error*/ }
    })
}

function handleSidepanel() {
    if (document.getElementById("userlist-holder").classList.contains("d-none")) {
        document.getElementById("userlist-holder").classList.remove("d-none")
    } else {
        document.getElementById("chat-area").classList.add("d-none")
        document.getElementById("userlist-holder").classList.add("d-none")
        let elements = document.getElementsByClassName("chat-button");

        for (var i = 0; i < elements.length; i++) {
            elements[i].classList.remove("active");
        }
    }
}

document.addEventListener('click', function (e) {
    if (e.target.classList.contains("btn-filter", "post-category")) {
        window.postOffset = 0;
        window.postLimit = 5;

        let posts = document.getElementsByClassName("post-container")
        while (posts[0]) {
            document.getElementById("post-area").removeChild(posts[0])
        }

        document.getElementById("view-posts").innerHTML = e.target.innerHTML.replace(/\s/g, '')

        loadPost(e.target.innerHTML)

        document.getElementById("load-more-btn").classList.remove("disabled")
        document.getElementById("load-more-btn").innerHTML = "Load more"
    }

    if (e.target.id === "sign-out-nav") {
        fetch('/api/sign-out', {
            method: 'POST',
            credentials: 'include',
        }).then(response => {
            if (response.ok) { window.location.href = '/sign-in' }
        })
    }

    if (e.target.id === "send-message") {
        e.preventDefault()
        if (e.target.matches('#send-message')) {
            let x = document.querySelector('form.chat-form').elements;
            let obj = {}
            obj.content = x['message-text'].value
            obj.recipientUsername = document.getElementById("chat-username").innerHTML
            obj.senderAccount = true

            if (obj.content !== "") {
                ws.send(JSON.stringify(obj))
            }
            document.getElementById('message-text').value = ''
        }
    }

    if (e.target.id === "user-list") {
        window.sidepanelOffset = 0
        if (document.getElementById("userlist-holder").classList.contains("d-none")) {
            clearUserlist()
            fillUserlist()
            handleSidepanel()
        } else {
            clearUserlist()
            handleSidepanel()
        }
    }

    if (e.target.id === "load-more-btn") {
        loadPost(document.getElementById("view-posts").innerHTML)
    }
});

function clearUserlist() {
    let div = document.getElementById("userlist-scroll-area");
    div.scrollTo(0, 0) // prevents scroll event by clicking userlist btn
    var children = div.children;
    for (let i = children.length - 1; i >= 0; i--) {
        let child = children[i];
        if (!child.classList.contains("d-none")) {
            div.removeChild(child);
        }
    }
}

async function loadPost(thread = "Viewall") {
    thread = thread.replace(/\s/g, '')
    await fetch('http://127.0.0.1:8080/api/posts?offset=' + window.postOffset + '&limit=' + window.postLimit + '&thread=' + thread, {
        method: 'GET',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        }
    }).then(response => {
        if (response.status === 200) {
            response.json().then((data) => {
                fetch('./static/templates/post.html').then(postTemplate => postTemplate.text())
                    .then(postTemplateText => {
                        console.log(data)
                        data.forEach(post => {
                            let div = document.createElement('div');
                            div.classList.add("post-container")
                            div.innerHTML = postTemplateText;
                            div.querySelector('.post-thread').innerHTML = post.Thread;
                            div.querySelector('.post-title').innerHTML = post.Title;
                            div.querySelector('.post-title').addEventListener('click', (e) => {
                                window.history.pushState({}, '', 'view-post?id=' + post.Id);
                                handleLocation();
                            })
                            getUsername(post.UserId).then(username => {
                                div.querySelector('.post-creator').innerHTML = username;
                            });

                            getPostCreationDate(post.Id).then(date => {
                                div.querySelector('.post-creation-date').innerHTML = date
                            });

                            if (post.Image !== "") {
                                div.querySelector('img').src = post.Image;
                            } else {
                                div.querySelector('.post-img').remove()
                            }

                            styleCategoryButton(div, post.Thread)
                            document.body.querySelector('#post-area').insertBefore(div, document.body.querySelector('#load-more'));
                        });
                        window.postOffset += window.postLimit
                    }).catch(error => {
                        console.error('Error:', error);
                    });
            })
        }
        else if (response.status === 204) {
            document.getElementById("load-more-btn").classList.add("disabled")
            document.getElementById("load-more-btn").innerHTML = "No more posts available"
        }
    });
}

function styleCategoryButton(element, category) {
    switch (category) {
        case "Question":
            element.querySelector('.post-thread').classList.add("btn-primary")
            break;
        case "Buy/Sell":
            element.querySelector('.post-thread').classList.add("btn-warning")
            break;
        case "Help!":
            element.querySelector('.post-thread').classList.add("btn-danger")
            break;
        case "Discussion":
            element.querySelector('.post-thread').classList.add("btn-success")
            break;
    }
}

async function getUsername(id) {
    return fetch('http://127.0.0.1:8080/api/username?id=' + id, {
        method: 'GET',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        },
    }).then(response => {
        if (response.ok) {
            return response.json();
        } else {
            throw new Error('Error fetching username');
        }
    }).then(data => {
        console.log(String(data));
        return String(data);
    }).catch(error => {
        console.error(error);
        return ''
    });
}

async function getPostCreationDate(id) {
    return fetch('http://127.0.0.1:8080/api/post-creation-date?id=' + id, {
        method: 'GET',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        },
    }).then(response => {
        if (response.ok) {
            return response.json();
        } else {
            throw new Error('Error fetching username');
        }
    }).then(data => {
        console.log(String(data));
        return String(data);
    }).catch(error => {
        console.error(error);
        return ''
    });
}