window.sidepanelOffset = 0
window.sidepanelLimit = 15

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

function createUserlistElement(id, username, status) {
    let chatButton = createChatButton(id, username, status)
    chatButton.addEventListener("click", async function loadChat() {
        await fetch('http://127.0.0.1:8080/api/load-chat?senderId=' + id, {
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
                'credentials': 'include',
            }
        }).then(response => {
            if (response.ok) {
                response.json().then((data) => {

                    document.getElementById("chat-username").innerHTML = username
                    document.getElementById("chat-user-status").innerHTML = status
                    if (Object.values(data)[0] !== null) {
                        stylizeChat(data, id)
                    }
                    handleChat(chatButton)

                })
            }
            else { /*display error*/ }
        })
    })
}

function createChatButton(id, username, status) {
    let template = document.getElementById('user-template')
    let userChatBtn = template.cloneNode(true)
    userChatBtn.classList.remove("d-none")

    userChatBtn.querySelector(".chat-with").innerHTML = username
    userChatBtn.querySelector(".user-status").innerHTML = status
    userChatBtn.classList.add("chat-with-user")
    userChatBtn.id = id
    userChatBtn.classList.add("chat-button")

    document.getElementById("userlist-scroll-area").appendChild(userChatBtn)
    return userChatBtn
}

function handleChat(userChatBtn) {
    console.log(userChatBtn)
    let chatArea = document.getElementById("chat-area").classList

    if (document.getElementById(userChatBtn.id).classList.contains("active")) {
        chatArea.add("d-none")
        userChatBtn.classList.remove("active")
        document.getElementById("chat-scroll-area").innerHTML = ""

    } else {
        let elements = document.getElementsByClassName("chat-button");

        for (var i = 0; i < elements.length; i++) {
            elements[i].classList.remove("active");
        }

        if (chatArea.contains("d-none")) {
            userChatBtn.classList.add("active")
            document.getElementById("chat-area").classList.remove("d-none")
        } else if (!chatArea.contains("d-none")) {
            userChatBtn.classList.add("active")
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
                for (let i = 0; i < Object.keys(data).length; i++) {
                    createUserlistElement(Object.values(data[i])[0], Object.values(data[i])[2], Object.values(data[i])[8])
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
        }
    }

    if (e.target.id === "user-list") {
        if (document.getElementById("userlist-scroll-area").classList.contains("initial")) {
            document.getElementById("userlist-scroll-area").addEventListener('scroll', handleScrollEvent);
            document.getElementById("userlist-scroll-area").classList.remove("initial")
        }
        window.sidepanelOffset = 0
        console.log(window.sidepanelOffset)
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
        const posts = document.body.querySelector('#post-area')
        let thread = document.getElementById("view-posts").innerHTML
        console.log(thread)
        fetch('http://127.0.0.1:8080/api/posts?offset=' + window.postOffset + '&limit=' + window.postLimit + '&thread=' + thread, {
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
                            data.forEach(post => {
                                let div = document.createElement('div');
                                div.innerHTML = postTemplateText;
                                div.classList.add("post-container")
                                div.querySelector('.post-thread').innerHTML = post.Thread;
                                div.querySelector('.post-title').innerHTML = post.Title;
                                if (post.Image !== "") {
                                    div.querySelector('img').src = post.Image;
                                } else {
                                    div.querySelector('.post-img').remove()
                                }

                                switch (post.Thread) {
                                    case "Question":
                                        div.querySelector('.post-thread').classList.add("btn-primary")
                                        break;
                                    case "Buy/Sell":
                                        div.querySelector('.post-thread').classList.add("btn-warning")
                                        break;
                                    case "Help!":
                                        div.querySelector('.post-thread').classList.add("btn-danger")
                                        break;
                                    case "Discussion":
                                        div.querySelector('.post-thread').classList.add("btn-success")
                                        break;
                                }

                                posts.insertBefore(div, document.body.querySelector('#load-more'));
                            });
                            window.postOffset += window.postLimit
                        }).catch(error => {
                            console.error('Error:', error);
                        });
                })
            } else if (response.status === 204) {
                document.getElementById("load-more-btn").classList.add("disabled")
                document.getElementById("load-more-btn").innerHTML = "No more posts available"
            }
        });
    }
});

function handleScrollEvent(event) {
    const e = event.target;
    if (e.scrollHeight - e.scrollTop === e.clientHeight) {
        debounce(100, fillUserlist());
    }
}

function clearUserlist() {
    let div = document.getElementById("userlist-scroll-area");
    div.scrollTo(0, 0)
    var children = div.children;
    for (let i = children.length - 1; i >= 0; i--) {
        let child = children[i];
        if (!child.classList.contains("d-none")) {
            div.removeChild(child);
        }
    }
}

function stylizeChat(data, senderId) {
    let chatFiller = document.getElementById("chat-scroll-area")
    let messages = Object.values(data)[0]
    for (let i = 0; i < messages.length; i++) {
        let [date, message, sender] = Array.from({ length: 3 }, () => document.createElement("div"));
        let senderName = document.createTextNode(Object.values(data)[1][i])
        let messageContent = document.createTextNode(messages[i]['Content'])
        let dateContent = document.createTextNode(messages[i]['Timestamp'])

        sender.appendChild(senderName)
        message.appendChild(messageContent)
        date.appendChild(dateContent)

        if (messages[i]['SenderId'] == senderId) {
            sender.classList.add("recipient-name")
            message.classList.add("recipient")
            date.classList.add("recipient-date")
        } else {
            sender.classList.add("sender-name")
            message.classList.add("sender")
            date.classList.add("sender-date")
        }
        chatFiller.appendChild(sender)
        chatFiller.appendChild(message)
        chatFiller.appendChild(date)
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
                        data.forEach(post => {
                            let div = document.createElement('div');
                            div.classList.add("post-container")
                            div.innerHTML = postTemplateText;
                            div.querySelector('.post-thread').innerHTML = post.Thread;
                            div.querySelector('.post-title').innerHTML = post.Title;
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

                            switch (post.Thread) {
                                case "Question":
                                    div.querySelector('.post-thread').classList.add("btn-primary")
                                    break;
                                case "Buy/Sell":
                                    div.querySelector('.post-thread').classList.add("btn-warning")
                                    break;
                                case "Help!":
                                    div.querySelector('.post-thread').classList.add("btn-danger")
                                    break;
                                case "Discussion":
                                    div.querySelector('.post-thread').classList.add("btn-success")
                                    break;
                            }
                            document.body.querySelector('#post-area').insertBefore(div, document.body.querySelector('#load-more'));
                        });
                        window.postOffset += window.postLimit
                    }).catch(error => {
                        console.error('Error:', error);
                    });
            })
        }
    });
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