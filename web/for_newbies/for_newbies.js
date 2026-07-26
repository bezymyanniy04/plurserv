
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

get_text()


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
              sessionStorage.setItem("userId", data.user_id);

    }
    catch(error){
        console.error(error);
    }

}

async function get_text() {
    
    try{
        var token = sessionStorage.getItem("token")
        var response = await fetch(`${api_link}/for_newbies`, {
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
        document.getElementById("for_newbies").value = data.text
    }
    catch(error){
        console.error(error);
    }

}


async function edit_text() {
    
    try{
        let text = document.getElementById("for_newbies").value
        var token = sessionStorage.getItem("token")
        var response = await fetch(`${api_link}/for_newbies`, {
        method: "PUT",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
            text: text
        })
  
    });

        if (!response.ok){
            throw new Error();
        } 

    }
    catch(error){
        console.error(error);
    }

}