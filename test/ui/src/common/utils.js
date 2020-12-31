export async function apiDelete(url) {
    return fetch(url, {
            credentials: 'same-origin',
            method: 'DELETE',
        })
        .then(async function(response) {
            let error = null;
            if (response.status !== 200) {
                error = await response.text();
            }
            return {
                error: error
            };
        });
}

export async function apiGet(url) {
    return fetch(url, {
            credentials: 'same-origin'
        })
        .then(async function(response) {
            let error = null;
            let data = null;
            if (response.status === 200) {
                if (response.headers.get("Content-Length") > 0) {
                    data = await response.json();
                }
            } else {
                error = await response.text();
            }
            return {
                error: error,
                data: data
            };
        });
}

export async function apiPost(url, body) {
    return fetch(url, {
            credentials: 'same-origin',
            method: 'POST',
            headers: {
            'Content-Type': 'application/json'
            },
            body: JSON.stringify(body)
        })
        .then(async function(response) {
            let error = null;
            let data = null;
            if (response.status === 200) {
                if (response.headers.get("Content-Length") > 0) {
                    data = await response.json();
                }
            } else {
                error = await response.text();
            }
            return {
                error: error,
                data: data
            };
        });
}

export async function apiPut(url, body) {
    return fetch(url, {
            credentials: 'same-origin',
            method: 'PUT',
            headers: {
            'Content-Type': 'application/json'
            },
            body: JSON.stringify(body)
        })
        .then(async function(response) {
            let error = null;
            let data = null;
            if (response.status === 200) {
                if (response.headers.get("Content-Length") > 0) {
                    data = await response.json();
                }
            } else {
                error = await response.text();
            }
            return {
                error: error,
                data: data
            };
        });
}
