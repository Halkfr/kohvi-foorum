if (window.location.pathname === '/') window.location.href = '/home';

document.addEventListener('click', e => {
    const activeElements = ["title", "home-nav", "create-post-nav", "profile-nav", "sign-out-nav", "register-link"];
    if (activeElements.includes(e.target.id)) {
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
    '/sign-up': 'sign-up.html'
}

const handleLocation = async () => {
    const isAuthorizated = await checkActiveSession()

    const path = window.location.pathname;
    const html = await fetch('./static/templates/' + routers[path]).then((data) => data.text()).catch(error => {
        console.error('Error:', error);
    });

    if (isAuthorizated === true) {
        if (path === "/sign-in" || path === "/sign-up") {
            window.history.replaceState({}, '', 'http://127.0.0.1:8080/home',)
            handleLocation()
        } else {
            if (document.body.querySelector('#main-container') === null) {
                const header = await fetch('./static/templates/header.html').then((data) => data.text());
                const sidepanel = await fetch('./static/templates/sidepanel.html').then((data) => data.text());
                document.body.innerHTML = header + sidepanel;
            }
            document.querySelector('#content-container').innerHTML = html
        }
    } else {
        if (path === "/sign-in") {
            document.body.innerHTML = html;
        } else if (path === "/sign-up") {
            const header = await fetch('./static/templates/header.html').then((data) => data.text());
            document.body.innerHTML = header + html;
        } else {
            window.history.replaceState({}, '', 'http://127.0.0.1:8080/sign-in',)
            handleLocation()
        }
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

async function checkActiveSession() {
    const response = await fetch('/api/session-status', {
        method: 'POST',
        credentials: 'include',
    });
    if (response.ok) {
        return true
    }
    return false
}