// handle home and profile 

const observer = new MutationObserver((mutations) => {
    mutations.forEach((mutation) => {
        const posts = document.body.querySelector('#post-scroll-area')
        const profile = document.body.querySelector('#view-profile-area')

        if (posts && !posts.contains(mutation.target) && mutation.target != document.body){
            console.log(mutation)
            console.log(posts)
            console.log('send posts request');
            fetch('http://127.0.0.1:8080/api/posts', {
                method: 'GET',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                }
            }).then(response => {
                if (response.ok) {
                    response.json().then((data) => {
                        fetch('./static/templates/post.html').then(postTemplate => postTemplate.text())
                            .then(postTemplateText => {
                                data.forEach(post => {
                                    let div = document.createElement('div');
                                    div.innerHTML = postTemplateText;
                                    div.querySelector('#post-thread').innerHTML = post.Thread;
                                    div.querySelector('#post-title').innerHTML = post.Title;
                                    div.querySelector('img').src = post.Image;
                                    posts.insertBefore(div, document.body.querySelector('#load-more'));
                                });
                            }).catch(error => {
                                console.error('Error:', error);
                            });
                        console.log(data);
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