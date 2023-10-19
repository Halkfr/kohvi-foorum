// handle home and profile

const observer = new MutationObserver((mutations) => {
    mutations.forEach((mutation) => {
        const posts = document.body.querySelector('#post-scroll-area')
        const profile = document.body.querySelector('#view-profile-area')

        if (posts && !posts.contains(mutation.target) && mutation.target != document.body && posts.classList.contains("initial")) {
            console.log(mutation)
            console.log(posts)
            console.log('send posts request');

            window.postOffset = 0;
            window.postLimit = 5;
            window.sidepanelOffset = 0
            window.sidepanelLimit = 15

            fetch('http://127.0.0.1:8080/api/posts?offset=' + window.postOffset + '&limit=' + window.postLimit, {
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
                                    div.querySelector('.post-thread').innerHTML = post.Thread;
                                    div.querySelector('.post-title').innerHTML = post.Title;
                                    div.querySelector('img').src = post.Image;

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
                        console.log(data);
                        posts.classList.remove("initial")
                    })
                }
            });
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
    });
});

observer.observe(document.body, { childList: true, subtree: true });