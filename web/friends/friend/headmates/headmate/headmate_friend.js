window.addEventListener('pageshow', function(event){
    if (event.persisted){
        this.sessionStorage.setItem("scroll", this.window.scrollY)
        window.location.reload();
    }

});



const env_link = "https://plurserv.fly.dev"
const api_link = "https://plurserv.fly.dev/api"

sessionStorage.removeItem("diary_id");
sessionStorage.removeItem("entry_id");
redirect_to_login();




// sessionStorage.removeItem("alter_id")


// get_alter();
    
    



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

get_alter();

async function get_alter() {
    
    try{
        let alter_id = sessionStorage.getItem("friend_alter_id")
        let token = sessionStorage.getItem("token")
        let response = await fetch(`${api_link}/alters/${alter_id}`, {
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
        document.getElementById("name").textContent = data.name;
        document.getElementById("avatar").src = "../../../../"+data.avatar;
        document.getElementById("pronouns").textContent = data.pronouns;
        document.getElementById("age").textContent = data.age;
        document.getElementById("role").textContent = data.role;
        document.getElementById("hm_colour").style.backgroundColor = data.colour;
        document.getElementById("description").textContent = data.description;
       
        
        if (sessionStorage.getItem("scroll") != null){
            console.log(sessionStorage.getItem("scroll"))
            window.scroll(0, sessionStorage.getItem("scroll"))
            sessionStorage.removeItem("scroll")
        }



    }
    
    catch(error){
        console.error(error);
    }

}