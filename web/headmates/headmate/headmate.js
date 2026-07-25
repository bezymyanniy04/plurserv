window.addEventListener('pageshow', function(event){
    if (event.persisted){
        this.sessionStorage.setItem("scroll", this.window.scrollY)
        window.location.reload();
    }

});



const env_link = "http://localhost:8080"
const api_link = "http://localhost:8080/api"

sessionStorage.removeItem("diary_id");
sessionStorage.removeItem("entry_id");
redirect_to_login();
const file_input = document.getElementById("file_add")
file_input.addEventListener('input', ()=>{
    add_avatar()
})

let alter_id = sessionStorage.getItem("alter_id")
// sessionStorage.removeItem("alter_id")

function existing_headmate(){
    get_alter();
    get_diary_by_alter()
    
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



async function get_alter() {
    
    try{
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
        let color = data.colour.split(" ");
        let colour = "#"+parseInt(color[0]).toString(16).padStart(2, '0') + parseInt(color[1]).toString(16).padStart(2, '0') + parseInt(color[2]).toString(16).padStart(2, '0');

        document.getElementById("name").value = data.name;
        document.getElementById("avatar").src = "../../"+data.avatar;
        document.getElementById("pronouns").value = data.pronouns;
        document.getElementById("age").value = data.age;
        document.getElementById("role").value = data.role;
        document.getElementById("hm_colour").value = colour;
        document.getElementById("description").value = data.description;
        // document.getElementById("av_change").value = data.avatar
        if (data.fronting){
            document.getElementById("fronting").style.backgroundColor = "magenta"
            document.getElementById("fronting").textContent = "remove from the front"
        }else{
            document.getElementById("fronting").style.backgroundColor = "pink"
            document.getElementById("fronting").textContent = "add to the front"

        }
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
async function download_avatar() {
    let filename = document.getElementById("avatar").src.split('app/')[1]
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
        let response = await fetch(`${api_link}/alters/avatar/${alter_id}`, {
        method: "DELETE",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
  
    });

        if (!response.ok){
            throw new Error();
        } 
        get_alter();


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
        let response = await fetch(`${api_link}/alters/avatar/${alter_id}`, {
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

        get_alter();


    }
    catch(error){
        console.error(error);
    }

}

async function get_diary_by_alter() {
    
    try{
        let token = sessionStorage.getItem("token")
        let response = await fetch(`${api_link}/diary/${alter_id}`, {
        method: "GET",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
  
    });
        if (!response.ok){
            document.getElementById("btn_add_diary").style.display = "block";
            document.getElementById("btn_diary").style.display = "none";


        }else{
            let data = await response.json();
            sessionStorage.setItem("diary_id", data.id)  
        } 
        

        
        
    }
    
    catch(error){
        console.error(error);
    }

}

async function add_diary() {
    
    try{
        let token = sessionStorage.getItem("token")
        let response = await fetch(`${api_link}/diaries/${alter_id}`, {
        method: "POST",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
  
    });
       if (!response.ok){
            throw new Error();
        } 
        let data = await response.json();
        sessionStorage.setItem("diary_id", data.id)
        window.location.href=`../../diaries/diary`


        
        
    }
    
    catch(error){
        console.error(error);
    }

}






async function change_front() {
    
    try{

        let token = sessionStorage.getItem("token")
        let response = await fetch(`${api_link}/alters/${alter_id}`, {
        method: "PATCH",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
  
    });

        if (!response.ok){
            throw new Error();
        } 
        get_alter()
        
    }
    
    catch(error){
        console.error(error);
    }

}

async function edit_alter(av) {
    
    try{

        let token = sessionStorage.getItem("token")
        let name = document.getElementById("name").value
        let pronouns = document.getElementById("pronouns").value
        let age = document.getElementById("age").value
        let role = document.getElementById("role").value
        let description = document.getElementById("description").value
        let colour = document.getElementById("hm_colour").value
        
        let response = await fetch(`${api_link}/alters/${alter_id}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
            name: name,
            pronouns: pronouns,
            age: age,
            role: role,
            description: description,
            colour: colour
        })
  
    });

        if (!response.ok){
            throw new Error();
        } 
        get_alter()
        
    }
    
    catch(error){
        console.error(error);
    }

}

async function delete_alter() {
    
    try{

        let token = sessionStorage.getItem("token")        
        let response = await fetch(`${api_link}/alters/${alter_id}`, {
        method: "DELETE",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
  
    });

        if (!response.ok){
            throw new Error();
        } 
        window.location.href=`..`
        
    }
    
    catch(error){
        console.error(error);
    }

}

async function post_alter() {
    
    try{

        let token = sessionStorage.getItem("token")
        let name = "StandardName"
        let avatar = ""
        let pronouns = document.getElementById("pronouns").value
        let age = document.getElementById("age").value
        let role = document.getElementById("role").value
        let description = document.getElementById("description").value
        let colour = document.getElementById("hm_colour").value

        let response = await fetch(`${api_link}/alters`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
            name: name,
            avatar: avatar,
            pronouns: pronouns,
            age: age,
            role: role,
            description: description,
            colour: colour
        })
  
    });

        if (!response.ok){
            throw new Error();
        } 
        let data = await response.json();
        sessionStorage.setItem("alter_id", data.id)
        alter_id = sessionStorage.getItem("alter_id")
        // window.location.href=`headmate`
    }
    
    catch(error){
        console.error(error);
    }

}