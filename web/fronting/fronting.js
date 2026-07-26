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
get_alters();
setInterval(change_time, 1000)




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
              localStorage.setItem("userId", data.user_id);

    }
    catch(error){
        console.error(error);
    }

}

let data_alters;

async function get_alters() {
    
    try{
        let token = sessionStorage.getItem("token")
        let userId = localStorage.getItem("userId")
         if (sessionStorage.getItem("query") != null){
            document.getElementById("query").value = sessionStorage.getItem("query")
            sessionStorage.removeItem("query")
        }
        let query = document.getElementById("query").value
        let response = await fetch(`${api_link}/alters/fronting/${userId}?query=${query}`, {
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
        
        let alters = document.getElementById("hm_list")
        alters.innerHTML = ""
        if (data.length == 0){
            var noone = document.getElementById("noFrontMessage")
            noone.style.display = "block"
            return
        }else{
            var noone = document.getElementById("noFrontMessage")
            noone.style.display = "none"
        }
        for (let i = 0; i < data.length; i++){
            let front = "-"
            let color_front = "var(--darker_main_colour)"
            let color_font = "var(--text_light)"
            if(!data[i].fronting){
                front = "+"
                color_front = "var(--main_colour)"
                color_font = "var(--text_dark)"
            }
            let diff_time = moment.duration(moment().diff(moment(data[i].fronting_since).utc()))

            let hours = Math.floor(diff_time.asHours()).toString().padStart(2, "0")
            let minutes = diff_time.minutes().toString().padStart(2, "0")
            let secondss = diff_time.seconds().toString().padStart(2, "0")
            let alter_colour_spl = data[i].colour.split(" ");
            let alter_colour = `rgb(${alter_colour_spl[0]}, ${alter_colour_spl[1]}, ${alter_colour_spl[2]})`
            let alter = `
            <div class="to_center_width">
                <div id="headmate_row" class="hm_row">
                    <button id="headmate_id" class="btn headmate_button" onclick="hm_page(${i})" ">
                        
                        <img src="${"../"+data[i].avatar}" class="img_alter" style="border-color: ${alter_colour}">
                        <div >
                            <h2 id="Headmate_name" class="text_hm">${data[i].name}</h2>
                            <p class="text_hm">${data[i].pronouns}</p>
                        </div>
                        <p id="time_since_${i}" class="text_hm">${hours}:${minutes}:${secondss}</p>
                    </button>
                    <button class="btn button_front_add text_hm" onclick="change_front(${i})" style="background-color: ${color_front}; color: ${color_font}"><h1>${front}</h1></button>
                </div>
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


function change_time(){
    for (let i = 0; i < data_alters.length; i++){
        let diff_time = moment.duration(moment().diff(moment(data_alters[i].fronting_since).utc()))
            let hours = Math.floor(diff_time.asHours()).toString().padStart(2, "0")
            let minutes = diff_time.minutes().toString().padStart(2, "0")
            let secondss = diff_time.seconds().toString().padStart(2, "0")
        document.getElementById(`time_since_${i}`).textContent = `${hours}:${minutes}:${secondss}` 
    }
}

function hm_page(alter_ident) {
    let alter_id = data_alters[alter_ident].id
    sessionStorage.setItem("alter_id", alter_id);
    window.location.href=`../headmates/headmate`

}



async function change_front(i) {
    
    try{
        let token = sessionStorage.getItem("token")
        let response = await fetch(`${api_link}/alters/${data_alters[i].id}`, {
        method: "PATCH",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
  
    });

        if (!response.ok){
            throw new Error();
        } 
        get_alters()
        
    }
    
    catch(error){
        console.error(error);
    }

}
