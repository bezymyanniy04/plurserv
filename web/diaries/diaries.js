window.addEventListener('pageshow', function(event){
    if (event.persisted){
        // this.sessionStorage.setItem("query", this.document.getElementById("query").value)
        this.sessionStorage.setItem("scroll", this.window.scrollY)
        window.location.reload();
    }

});

const env_link = "https://plurserv.fly.dev"
const api_link = "https://plurserv.fly.dev/api"


const rootStyles = getComputedStyle(document.querySelector(':root'))

sessionStorage.removeItem("diary_id");
sessionStorage.removeItem("entry_id");
sessionStorage.removeItem("alter_id");

redirect_to_login();
get_diaries();





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

let data_diaries;

async function get_diaries() {
    
    try{
        
        let token = sessionStorage.getItem("token")
        let response = await fetch(`${api_link}/diaries`, {
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
        
        data_diaries = data
        let diaries = document.getElementById("diaries_list")
        diaries.innerHTML = ""
        if (data.length == 0){
            document.getElementById('no_one').style.display = 'flex';
            return;
        }else{
            document.getElementById('no_one').style.display = 'none';

        }
        for (let i = 0; i < data.length; i++){
            let color = data[i].block_colour.split(" ");
            let text_color;
            if (color[0]>130 & color[1]>130 & color[2]>130){
                text_color = rootStyles.getPropertyValue('--dark_theme_bg_colour') 
            }else{
                text_color = rootStyles.getPropertyValue('--bg_colour')
            }
            let diary = `
                <button onclick="open_diary(${i})" class = "btn flex_cont corn" style="background-color: rgb(${parseInt(color[0])}, ${parseInt(color[1])}, ${parseInt(color[2])});; color: ${text_color};">
                <img class="img_button" src="${"../"+data[i].alter_avatar}">
                <p class="text_button">${data[i].alter_name}'s Diary</p>
                </button>

            `       
            diaries.innerHTML += diary
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



function open_diary(i){
    sessionStorage.setItem("alter_id", data_diaries[i].alter_id);
    sessionStorage.setItem("diary_id", data_diaries[i].id);

    window.location.href=`diary`;

}

let data_alters;

async function get_alters() {
    
    try{

        
        let token = sessionStorage.getItem("token")
        let response = await fetch(`${api_link}/alters_wo_diaries`, {
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
        data_alters = data;
        document.getElementById("adding_popover").style.display = 'block';
        let alters = document.getElementById("hm_list");
        alters.innerHTML = "";
        if (data.length == 0){
            document.getElementById('no_one_hm').style.display = 'flex';
            return;
        }else{
            document.getElementById('no_one_hm').style.display = 'none';

        }
        for (let i = 0; i < data.length; i++){
            let color = data[i].colour.split(" ");
            let alter = `

                    <button id="headmate_id"  onclick="post_diary(${i})" class="btn_popup flex_cont">
                        
                        <img src="${data[i].avatar}" class="img_popup" style="border: 0.3vw solid rgb(${parseInt(color[0])}, ${parseInt(color[1])}, ${parseInt(color[2])});">
                        
                        <p id="Headmate_name" style="padding-left: 5px"><b>${data[i].name}</b></p>
                        
                    </button>


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


async function post_diary(i) {
    
    try{
        let alter = data_alters[i].id;
        let token = sessionStorage.getItem("token");
        let response = await fetch(`${api_link}/diaries/${alter}`, {
        method: "POST",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
  
    });

        if (!response.ok){
            throw new Error();
        } 
        get_diaries()
        document.getElementById("hm_list").innerHTML = ""
        document.getElementById("adding_popover").style.display = 'none'
    }
    catch(error){
        console.error(error);
    }

}