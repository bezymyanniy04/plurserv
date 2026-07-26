


window.addEventListener('pageshow', function(event){
    if (event.persisted){
        window.location.reload();
        this.sessionStorage.setItem("scroll", this.window.scrollY)

    }

});
    

const env_link = "https://plurserv.fly.dev"
const api_link = "https://plurserv.fly.dev/api"

const rootStyles = getComputedStyle(document.querySelector(':root'))

sessionStorage.removeItem("diary_id");
sessionStorage.removeItem("entry_id");
sessionStorage.removeItem("alter_id");


redirect_to_login();
get_user_myself();
        get_fronting_alters();
let data_alters;

function hov(){
    const headmates_list = document.querySelectorAll('#frontContainer');
headmates_list.forEach((element, index) =>{
    element.addEventListener('mouseenter', () =>{
        
        element.style.backgroundColor = element.style.backgroundColor.slice(0, -4)+"1)"
    });
    element.addEventListener('mouseleave', () =>{

        element.style.backgroundColor = element.style.backgroundColor.slice(0, 3)+"a"+element.style.backgroundColor.slice(3, -1)+", 0.8)"
    })


}
) 
}



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
        get_user_myself();
        get_fronting_alters();
    }
    catch(error){
        console.error(error);
    }

}

async function get_user_myself() {
    
    try{
        var token = sessionStorage.getItem("token")
        var userId = localStorage.getItem("userId")
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
        try1 = 0;
       document.getElementById("systemName").textContent = data.system_name + "!"
    }
    catch(error){
        console.error(error);
    }

}


async function get_fronting_alters() {
    
    try{
        var token = sessionStorage.getItem("token")
        var userId = localStorage.getItem("userId")
        var response = await fetch(`${api_link}/alters/fronting/${userId}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
  
    });

        if (!response.ok){

            throw new Error();
        } 
        var data = await response.json();
        try2 = 0
        data_alters = data
        if (data.length == 0){
            var noone = document.getElementById("noFrontMessage")
            noone.style.display = "flex"
            return
        }else{
            var noone = document.getElementById("noFrontMessage")
            noone.style.display = "none"
        }
        
        
        var fronting = document.getElementById("frontSection")

        for (var i = 0; i < data.length && i < 4; i++){
            
            let color = data[i].colour.split(" ");
            let text_color;
            if (color[0]>130 & color[1]>130 & color[2]>130){
                text_color = rootStyles.getPropertyValue('--dark_theme_bg_colour') 
            }else{
                text_color = rootStyles.getPropertyValue('--bg_colour')
            }
            var alter = `
                    <button  id="frontContainer" onclick="hm_page(${i})" class="card" style="background-color: rgba(${parseInt(color[0])}, ${parseInt(color[1])}, ${parseInt(color[2])}, 0.8)">
                        <div id="frontContainerInner" class="card" >
                        <img src="${data[i].avatar}" id="frontAv">
                        </div>
                        
                        <p id="frontName" style="color: ${text_color}">${data[i].name}</p>
                    </button>
            `
            fronting.innerHTML += alter
        }
        if (data.length > 4){
            fronting.innerHTML += `
                    <button id="frontButtonContainer" onclick="location.href='fronting'">
                        <p>>>></p>
                    </button>`
        }
        if (sessionStorage.getItem("scroll") != null){
            window.scroll(0, sessionStorage.getItem("scroll"))
            sessionStorage.removeItem("scroll")
        }
        hov()
    }
    catch(error){
        console.error(error);
    }

}



function hm_page(alter_ident) {
    let alter_id = data_alters[alter_ident].id
    sessionStorage.setItem("alter_id", alter_id);
    window.location.href=`headmates/headmate`

}



