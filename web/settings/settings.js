
window.addEventListener('pageshow', function(event){
    if (event.persisted){
        window.location.reload();
        this.sessionStorage.setItem("scroll", this.window.scrollY)

    }

});
    

const env_link = "https://plurserv.fly.dev"
const api_link = "https://plurserv.fly.dev/api"


redirect_to_login();



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
get_user_myself()
    }
    catch(error){
        console.error(error);
    }

}

function log_out(){
    sessionStorage.clear();
    localStorage.clear();
    redirect_to_login();
}

let userId = localStorage.getItem("userId")



async function get_user_myself() {
    
    try{
        var token = sessionStorage.getItem("token")
        
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
       document.getElementById("systemName").value = data.system_name
       document.getElementById("themeSelect").value = data.theme
       document.getElementById("fontSelect").value = data.font
       document.getElementById("avatar_img").src = "../"+data.avatar
    }
    catch(error){
        console.error(error);
    }

}
async function edit_info() {
    
    try{
        var token = sessionStorage.getItem("token");
        let sysname = document.getElementById("systemName").value;
        let theme = document.getElementById("themeSelect").value;
        let font = document.getElementById("fontSelect").value
        var response = await fetch(`${api_link}/userinfo`, {
        method: "PUT",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
            system_name: sysname,
            theme: parseInt(theme),
            font: font,
        })
  
    });

        if (!response.ok){
            throw new Error();
        } 
        var data = await response.json();
       document.getElementById("avatar_img").src = data.avatar
    }
    catch(error){
        console.error(error);
    }

}

const file_input = document.getElementById("file_add")
file_input.addEventListener('input', ()=>{
    add_avatar()
})
async function download_avatar() {
    let filename = document.getElementById("avatar_img").src.split('app/')[1]
    // Option B: Force download programmatically
    const response = await fetch(`/app/${filename}`);
    const blob = await response.blob();
    
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename; // This forces download instead of display
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
}

async function delete_avatar() {
    
    try{    
        let token = sessionStorage.getItem("token");
        let response = await fetch(`${api_link}/users/avatar`, {
        method: "DELETE",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
  
    });

        if (!response.ok){
            throw new Error();
        } 
        get_user_myself();


    }
    catch(error){
        console.error(error);
    }

}
async function add_avatar() {
    
    try{    
        let token = sessionStorage.getItem("token");
        const fileInput = document.getElementById('file_add');
        const file = fileInput.files[0];
        const formData = new FormData();
        formData.append('photo', file);
        let response = await fetch(`${api_link}/users/avatar`, {
        method: "PUT",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: formData
  
    });

        if (!response.ok){
            throw new Error();
        } 

        get_user_myself();


    }
    catch(error){
        console.error(error);
    }

}