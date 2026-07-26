


window.addEventListener('pageshow', function(event){
    if (event.persisted){
        window.location.reload();
        this.sessionStorage.setItem("scroll", this.window.scrollY)

    }

});
    
const env_link = "https://plurserv.fly.dev"
const api_link = "https://plurserv.fly.dev/api"

sessionStorage.removeItem("diary_id");
sessionStorage.removeItem("entry_id");
sessionStorage.removeItem("alter_id");

redirect_to_login();
get_user();
get_fronting_alters();

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

async function get_user() {
    
    try{
        var token = sessionStorage.getItem("token")
        var userId = sessionStorage.getItem("friend_id")
        var response = await fetch(`${api_link}/users/${userId}`, {
        method: "GET",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
  
    });

        if (!response.ok){
            throw new Error();
        } 
        var data = await response.json();
       document.getElementById("name").textContent = data.system_name
    }
    catch(error){
        console.error(error);
    }

}

let data_alters;

async function get_fronting_alters() {
    
    try{
        let token = sessionStorage.getItem("token")
        let userId = sessionStorage.getItem("friend_id")
         
        let response = await fetch(`${api_link}/alters/fronting/${userId}`, {
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
        data_alters = data
        
        let alters = document.getElementById("frontSection")
        alters.innerHTML = ""
        if (data.length == 0){
            let noone = document.getElementById("noFrontMessage")
            noone.style.display = "block"
            return
        }else{
            let noone = document.getElementById("noFrontMessage")
            noone.style.display = "none"
        }
        for (let i = 0; i < data.length; i++){
            let alter = `
                <div id="headmate_row" style="display: flex; height: 150px; justify-content: stretch; margin: 10px">
                    <button id="headmate_id" onclick="hm_page(${i})" style="width: 80%;  border-top-left-radius: 15px; border-bottom-left-radius: 15px; background-color: lightpink; padding: 5px; display: flex;">
                        
                        <img src="${data[i].avatar}" style="width: 110px; height: 110px; margin: 5px; border-radius: 10px; border: 5px solid ${data[i].colour};">
                        <div style="padding: 5px;">
                            <h2 id="Headmate_name">${data[i].name}</h2>
                            <p>${data[i].pronouns}</p>
                        </div>
                    </button>
                </div>

            `           
            alters.innerHTML += alter
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


function hm_page(alter_ident) {
    let alter_id = data_alters[alter_ident].id
    sessionStorage.setItem("friend_alter_id", alter_id);
    window.location.href=`headmates/headmate`

}



