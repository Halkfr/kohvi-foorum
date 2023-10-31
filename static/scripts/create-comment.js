document.addEventListener('click', function (e) {
    if (location.pathname === "/view-post") {
        if (e.target.matches('#comment-submit-btn')) {
            if (document.querySelector('form.create-comment').checkValidity()) {
                e.preventDefault();
                let formData = new FormData();
                formData.append('content', document.querySelector('form.create-comment #create-comment-content').value);
                const urlParams = new URLSearchParams(window.location.search);
                const postId = urlParams.get('id')
                formData.append('postid', postId)

                fetch('http://127.0.0.1:8080/api/add-comment', {
                    method: 'POST',
                    headers: {
                        'credentials': 'include',
                    },
                    body: formData
                }).then(response => {
                    if (response.ok) {
                        fetchComments(postId)
                        document.getElementById("create-comment").reset();
                    } else {/*display error*/ }
                })
            } else {
                if (!document.querySelector('form.create-comment #create-comment-content').checkValidity()) {
                    document.querySelector('form.create-comment #create-comment-content').reportValidity()
                }
            }
        }
    }
})

// Id        int
// Content   string
// PostId    int
// UserId    int
// Timestamp string