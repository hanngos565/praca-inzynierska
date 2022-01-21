const address = "http://localhost:8081"

export const getImages = async() => {
    const response = await fetch(`${address}/v1/images`, {
        method: 'GET',
        mode: 'cors',
    })
    return await response.json()
}

export const uploadImage = async(data) => {
    return await fetch(`${address}/v1/images`, {
        method: 'PUT',
        mode: 'cors',
        body: JSON.stringify(data),
        headers: {
            'Content-type': 'application/json'
        }
    })
}

export const getModels = async(alg) => {
    const response = await fetch(`${address}/v1/models/${alg}`, {
        method: 'GET',
        mode: 'cors',
    })
    return await response.json()
}

export const addModel = async(data) => {
    return await fetch(`${address}/v1/models`, {
        method: 'PUT',
        mode: 'cors',
        body: data
    });
}

export const runSimulation = async(opType, data) => {
    return await fetch(`${address}/v1/simulation-results/${opType}`, {
        method: 'POST',
        mode: 'cors',
        body: data
    });
}

export const getResults = async(opType, alg) => {
    const response = await fetch(`${address}/v1/simulation-results/${opType}/${alg}`, {
        method: 'GET',
        mode: 'cors',
    })
    return response.json()
}