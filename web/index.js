


window.addEventListener('pageshow', function(event){
    if (event.persisted){
        window.location.reload();
        this.sessionStorage.setItem("scroll", this.window.scrollY)

    }

});
    

const env_link = "http://localhost:8080"
const api_link = "http://localhost:8080/api"

const rootStyles = getComputedStyle(document.querySelector(':root'))

sessionStorage.removeItem("diary_id");
sessionStorage.removeItem("entry_id");
sessionStorage.removeItem("alter_id");

redirect_to_login();
get_user_myself();
get_fronting_alters();


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
    if (localStorage.getItem("refresh_token")=== null){
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
    }
    catch(error){
        console.error(error);
    }

}

async function get_user_myself() {
    
    try{
        var token = sessionStorage.getItem("token")
        var userId = sessionStorage.getItem("userId")
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
       document.getElementById("systemName").textContent = data.system_name + "!"
    }
    catch(error){
        console.error(error);
    }

}

let data_alters;
async function get_fronting_alters() {
    
    try{
        var token = sessionStorage.getItem("token")
        var userId = sessionStorage.getItem("userId")
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
            console.log(sessionStorage.getItem("scroll"))
            window.scroll(0, sessionStorage.getItem("scroll"))
            sessionStorage.removeItem("scroll")
        }
        hov()
    }
    catch(error){
        console.error(error);
    }

}


// var alter = `
//                     var alter = `
    //                 <button  id="frontContainer" onclick="hm_page(${i})" class="card" style="
    // background-image: linear-gradient(rgba(${parseInt(color[0])}, ${parseInt(color[1])}, ${parseInt(color[2])}, 0.8), rgba(${parseInt(color[0])/1.5}, ${parseInt(color[1])/1.5}, ${parseInt(color[2])/1.5}, 0.8))">
    //                     <div id="frontContainerInner" class="card" >
    //                     <img src="${data[i].avatar}" id="frontAv">
    //                     </div>
                        
    //                     <p id="frontName">${data[i].name}</p>
    //                 </button>
    //         `
//             `

function hm_page(alter_ident) {
    let alter_id = data_alters[alter_ident].id
    sessionStorage.setItem("alter_id", alter_id);
    window.location.href=`headmates/headmate`

}



