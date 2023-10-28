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

        if (profile && !profile.contains(mutation.target)) {
            console.log('send request')
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
                        console.log(data)
                    })
                }
            })
        }

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
                            document.body.querySelector('#view-post-image').src= data.Image
                            document.body.querySelector('#view-post-content').innerHTML = data.Content
                            console.log(data)

                            styleCategoryButton(post, data.Thread)
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