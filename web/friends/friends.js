
let th = localStorage.getItem("theme")
if (th != null){
    document.documentElement.style.setProperty('color-scheme', th)

}



window.addEventListener('pageshow', function(event){
    if (event.persisted){
        this.sessionStorage.setItem("query", this.document.getElementById("query").value)
        this.sessionStorage.setItem("scroll", this.window.scrollY)
        window.location.reload();
    }

});

const env_link = "https://plurserv.fly.dev"
const api_link = "https://plurserv.fly.dev/api"

sessionStorage.removeItem("diary_id");
sessionStorage.removeItem("entry_id");
sessionStorage.removeItem("alter_id");


redirect_to_login();
document.getElementById("user_id").textContent = "id: "+ localStorage.getItem("userId");
friends();





async function redirect_to_login() {
    if (localStorage.getItem("refresh_token")=== null || localStorage.getItem("userId")=== null){
        window.location.href=`${env_link}/app/login`
    }else{
        refresh()
    }
}

async function refresh() {
    
    try{
        var refresh = localStorage.getItem("refresh_token")
        var response = await fetch(`${api_link}/refresh`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${refresh}`
        }
  
    });

        if (!response.ok){
            window.location.href=`${env_link}/app/login`
            throw new Error();
        } 
        var data = await response.json();
       sessionStorage.setItem("token", data.token);
              localStorage.setItem("userId", data.user_id);

    }
    catch(error){
        console.error(error);
    }

}

let data_friends;

async function friends() {
    
    try{
        document.getElementById("friends_list").style.display = "block"
        document.getElementById("sent").style.display = "none"
        document.getElementById("pending_list").style.display = "none"
        
        let token = sessionStorage.getItem("token")
        let response = await fetch(`${api_link}/friends`, {
        method: "GET",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
  
    });

        if (!response.ok){
            throw new Error();
        } 
        let data = await response.json();
        data_friends = data
        let friends = document.getElementById("friends_list")
        friends.innerHTML = ""
        for (let i = 0; i < data.length; i++){
            

        let friend = `
                <button id="friend" onclick="open_friend(${i})" style="width: 80%;  border-top-left-radius: 15px; border-bottom-left-radius: 15px; background-color: lightgray; padding: 5px; display: flex;">            
                    <img src="${data[i].avatar}" style="width: 110px; height: 110px; margin: 5px; border-radius: 10px; border: 5px solid black;">
                    <h2 id="friend_name" style="padding: 5px;" >${data[i].system_name}</h2>                
                </button>
            `
            friends.innerHTML += friend
        }

        if (sessionStorage.getItem("scroll") != null){
            window.scroll(0, sessionStorage.getItem("scroll"))
            sessionStorage.removeItem("scroll")
        }
    }
    catch(error){
        console.error(error);
    }

}

async function requests() {
    
    try{
        document.getElementById("friends_list").style.display = "none"
        document.getElementById("sent").style.display = "block"
        document.getElementById("pending_list").style.display = "none"
        
        let token = sessionStorage.getItem("token")
        let response = await fetch(`${api_link}/friends/requests/sender`, {
        method: "GET",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
  
    });

        if (!response.ok){
            throw new Error();
        } 
        let data = await response.json();
        data_friends = data
        let friends = document.getElementById("sent_list")
        friends.innerHTML = ""
        for (let i = 0; i < data.length; i++){
            

        let friend = `
                <div id="friend" style="width: 80%;  border-top-left-radius: 15px; border-bottom-left-radius: 15px; background-color: lightgray; padding: 5px; display: flex;">            
                    <img src="${data[i].avatar}" style="width: 110px; height: 110px; margin: 5px; border-radius: 10px; border: 5px solid black;">
                    <h2 id="friend_name" style="padding: 5px;" >${data[i].system_name}</h2>                
                </div>
            `
            friends.innerHTML += friend
        }

        if (sessionStorage.getItem("scroll") != null){
            window.scroll(0, sessionStorage.getItem("scroll"))
            sessionStorage.removeItem("scroll")
        }
    }
    catch(error){
        console.error(error);
    }

}

let data_requests;

async function pending() {
    
    try{
        document.getElementById("friends_list").style.display = "none"
        document.getElementById("sent").style.display = "none"
        document.getElementById("pending_list").style.display = "block"
        
        let token = sessionStorage.getItem("token")
        let response = await fetch(`${api_link}/friends/requests/reciever`, {
        method: "GET",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
  
    });

        if (!response.ok){
            throw new Error();
        } 
        let data = await response.json();
        data_requests = data
        let friends = document.getElementById("pending_list")
        friends.innerHTML = ""
        for (let i = 0; i < data.length; i++){
            

        let friend = `
                <div id="friend" style="width: 80%;  border-top-left-radius: 15px; border-bottom-left-radius: 15px; background-color: lightgray; padding: 5px; display: flex;">            
                    <img src="${data[i].avatar}" style="width: 110px; height: 110px; margin: 5px; border-radius: 10px; border: 5px solid black;">
                    <h2 id="friend_name" style="padding: 5px;" >${data[i].system_name}</h2>                
                    <button onclick="answer(${i}, '1')">Yes</button>
                    <button onclick="answer(${i}, '2')">No</button>
                </div>
            `
            friends.innerHTML += friend
        }

        if (sessionStorage.getItem("scroll") != null){
            window.scroll(0, sessionStorage.getItem("scroll"))
            sessionStorage.removeItem("scroll")
        }
    }
    catch(error){
        console.error(error);
    }

}

async function answer(i, answer) {
    
    try{
        let requestId = data_requests[i].request_id
        let token = sessionStorage.getItem("token")
        let response = await fetch(`${api_link}/friends/requests`, {
        method: "PUT",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
            id: requestId,
            answer: parseInt(answer)
        })
  
    });

        if (!response.ok){
            throw new Error();
        } 
        pending();

    }
    catch(error){
        console.error(error);
    }

}

async function add_friend() {
    
    try{
        let friend_id = document.getElementById("friend_id").value
        let token = sessionStorage.getItem("token")
        let response = await fetch(`${api_link}/friends/requests/${friend_id}`, {
        method: "POST",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
  
    });

        if (!response.ok){
            throw new Error();
        } 
        requests();
        document.getElementById("friend_id").value = ""
        document.getElementById('add_friend_div').style.display = 'none'

    }
    catch(error){
        console.error(error);
    }

}

function copy_id(){
    console.log("hey")
    navigator.clipboard.writeText(localStorage.getItem("userId"))
}


function open_friend(i) {
    let friend_id = data_friends[i].id
    sessionStorage.setItem("friend_id", friend_id);
    window.location.href=`friend`

}










