const observer = new MutationObserver((mutations) => {
    mutations.forEach((mutation) => {
        const posts = document.body.querySelector('#post-area')
        const profile = document.body.querySelector('#view-profile-area')
        const post = document.body.querySelector('#view-post-area')

        if (posts && document.body.querySelector('#post-area').classList.contains("initial")) {
            document.getElementById("view-posts").innerHTML = "Viewall"
            window.postOffset = 0;
            window.postLimit = 5;
            loadPost()
            posts.classList.remove("initial")
        }

        // fill profile page
        if (profile && !profile.contains(mutation.target)) {
            fetch('http://127.0.0.1:8080/api/user', {
                method: 'GET',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                }
            }).then(response => {
                if (response.ok) {
                    response.json().then((data) => {
                        document.body.querySelector('#username').innerHTML = data.Username
                        document.body.querySelector('#full-name').innerHTML = data.Firstname + " " + data.Lastname
                        document.body.querySelector('#birthdate').innerHTML = data.Birthdate
                        document.body.querySelector('#gender').innerHTML = data.Gender
                        document.body.querySelector('#email').innerHTML = data.Email
                        document.body.querySelector('#joined').innerHTML = data.Timestamp
                    })
                }
            })
        }

        // fill view-post page
        if (post && !post.contains(mutation.target)) {
            const params = new URLSearchParams(window.location.search);

            if (params.has('id')) {
                const id = params.get('id')

                return fetch('http://127.0.0.1:8080/api/post?id=' + id, {
                    method: 'GET',
                    headers: {
                        'Accept': 'application/json',
                        'Content-Type': 'application/json',
                    },
                }).then(response => {
                    if (response.ok) {
                        response.json().then((data) => {
                            document.body.querySelector('#view-post-category').innerHTML = data.Thread
                            document.body.querySelector('#view-post-title').innerHTML = data.Title
                            if (data.Image) {
                                document.body.querySelector('#view-post-image').src = data.Image
                            } else {
                                document.body.querySelector('#view-post-image').remove()
                            }
                            document.body.querySelector('#view-post-content').innerHTML = data.Content
                            getUsername(data.UserId).then((username) => {
                                document.body.querySelector('.post-creator').innerHTML = username
                            })
                            document.body.querySelector('.post-creation-date').innerHTML = data.Timestamp
                            // console.log(data)

                            styleCategoryButton(post, data.Thread)

                            fetchComments(id);
                        })
                    } else {
                        throw new Error('Error fetching post');
                    }
                })

            } else { // no id in query parameters
                window.history.pushState({}, '', 'home');
                handleLocation();
            }
        }
    });
});

observer.observe(document.body, { childList: true, subtree: true });

async function fetchComments(id) {
    return fetch('http://127.0.0.1:8080/api/comments?id=' + id, {
        method: 'GET',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        },
    }).then(response => {
        if (response.ok) {
            response.json().then((comments) => {
                fetch('./static/templates/comment.html').then(commentTemplate => commentTemplate.text())
                    .then(commentTemplateText => {
                        document.body.querySelector('.comments').innerHTML = ""
                        if (comments.length != 0) {
                            comments.forEach(comment => {
                                let div = document.createElement('div');
                                div.classList.add("comment-container")
                                div.innerHTML = commentTemplateText;
                                div.querySelector('.comment-content').innerHTML = comment.Content;
                                div.querySelector('.comment-timestamp').innerHTML = comment.Timestamp;

                                getUsername(comment.UserId).then(username => {
                                    div.querySelector('.comment-username').innerHTML = username
                                });
                                document.body.querySelector('.comments').prepend(div);
                            });
                        }
                    }).catch(error => {
                        console.error('Error:', error);
                    });
            })
        } else {
            throw new Error('Error fetching comments');
        }
    })
}

// observe side panel

const observerSidepanel = new MutationObserver((mutation) => {
    const sidepanel = document.body.querySelector('#sidepanel')

    if (sidepanel) {
        document.getElementById("chat-scroll-area").addEventListener('scroll', (event) => {
            const e = event.target;
            if (e.scrollTop === 0 && e.scrollHeight > e.clientHeight) {
                debounce(1000, loadChat(document.getElementById("userlist-holder").getElementsByClassName("active")[0], "top"))
            }
        });

        document.getElementById("userlist-scroll-area").addEventListener('scroll', (event) => {
            const e = event.target;
            if (e.scrollHeight - e.scrollTop === e.clientHeight) {
                debounce(100, fillUserlist());
            }
        });

        fetch('http://127.0.0.1:8080/api/user-notifications-number', {
            method: 'GET',
            headers: {
                'credentials': 'include',
            },
        }).then(response => response.text()).then(data => {
            if (data != 0) { // change total count for current user
                document.getElementById("user-list").getElementsByClassName("badge")[0].innerHTML = data
            } else {
                document.getElementById("user-list").getElementsByClassName("badge")[0].innerHTML = ""
            }
        })

        observerSidepanel.disconnect()
    }
});

observerSidepanel.observe(document.body, { childList: true, subtree: true });