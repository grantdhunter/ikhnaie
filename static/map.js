mapboxgl.accessToken = MAPBOX_API_KEY

const ROOM_ID = getRoomId()
const CLIENT_ID = getClientId()

let markers = {}

function getRoomId() {
    let path = window.location.pathname
    if (path == '/') {
        let roomId = window.crypto.getRandomValues(new Uint32Array(1))[0].toString(16)
        window.history.pushState({ roomId: roomId }, "", "/" + roomId)
        return roomId
    } else {
        let elem = path.split("/")
        return elem.pop()
    }

}

function getClientId() {
    let id = localStorage.getItem("clientId")
    if (id) {
        return id
    }
    id = window.crypto.getRandomValues(new Uint32Array(1))[0].toString(16)
    localStorage.setItem("clientId", id)
    return id
}


function showError(error) {
    let x = document.getElementById("container")
    switch (error.code) {
        case error.PERMISSION_DENIED:
            x.innerHTML = "User denied the request for Geolocation."
            break;
        case error.POSITION_UNAVAILABLE:
            x.innerHTML = "Location information is unavailable."
            break;
        case error.TIMEOUT:
            x.innerHTML = "The request to get user location timed out."
            break;
        case error.UNKNOWN_ERROR:
            x.innerHTML = "An unknown error occurred."
            break;
    }
}
function showLoc(loc, clientId, {colour='#200ca0', fly=false} = {}) {
    const coordinates = [loc.lng, loc.lat];

    if (fly) {
        map.flyTo({
            center: coordinates,
            zoom: 14
        })
    }

    let m = markers[clientId]
    if (!m) {
        m = new mapboxgl.Marker({
            color: colour
        })
        markers[clientId] = m
    }

    m.setLngLat(coordinates)
        .addTo(map)

}

function positionToLoc(position) {
    return {
        lat: position.coords.latitude,
        lng: position.coords.longitude,
        acc: position.coords.accuracy,
        heading: position.coords.heading,
        speed: position.coords.speed,
        ts: position.timestamp
    }
}
function sendLoc(loc) {
    socket.send(JSON.stringify({
        roomId: ROOM_ID,
        clientId: CLIENT_ID,
        loc:loc
    }))
}

var map = new mapboxgl.Map({
    container: 'map',
    style: 'mapbox://styles/mapbox/streets-v11'
});
map.on('load', function() {
    if (navigator.geolocation) {
        let options = {
            enableHighAccuracy: true,
            maximumAge: 60000
        }
        navigator.geolocation.getCurrentPosition((position) => {
            document.getElementById("message").style.display = 'none'
            document.getElementById("map").style.visibility = 'visible'
            let loc = positionToLoc(position)
            showLoc(loc, CLIENT_ID, {fly:true})
            sendLoc(loc)
        }, showError, options);
        navigator.geolocation.watchPosition((position) => {
            let loc = positionToLoc(position)
            showLoc(loc, CLIENT_ID)
            sendLoc(loc)
        })
    } else {
        doc.innerHTML("Oh no, you has no location")

    }
})

map.on('click', function(e) {
    sendLoc({lat: e.lngLat.lat,
             lng: e.lngLat.lng})
})


var socket
var reconnectTimer
function connectWebSocket() {
    socket = new WebSocket("wss://" + window.location.host + "/ws/?clientId=" + CLIENT_ID)
}

connectWebSocket()
socket.onclose = function() {
    console.log("ws: connecting...")
    reconnectTimer = setTimeout(connectWebSocket)
}


socket.onopen = function(e) {
    if (reconnectTimer) {
        setTimeout(reconnectTimer, 1000)
    }
    console.log("ws: connected", e)
    socket.send(JSON.stringify({
        roomId: ROOM_ID,
        clientId: CLIENT_ID
    }))
}
socket.onmessage = function(e) {
    console.log("ws: message ", e)
    let evt = JSON.parse(e.data)
    if (evt.loc) {
        showLoc(evt.loc, evt.clientId, {colour: '#a0360c'})
    }
}



window.onbeforeunload = function() {
    socket.close()
}
