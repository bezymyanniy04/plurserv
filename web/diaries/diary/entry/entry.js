window.addEventListener('pageshow', function(event){
    if (event.persisted){
        this.sessionStorage.setItem("scroll", this.window.scrollY)
        window.location.reload();
    }

});



const env_link = "https://plurserv.fly.dev"
const api_link = "https://plurserv.fly.dev/api"

redirect_to_login();


let alter_id = sessionStorage.getItem("alter_id");
let diary_id = sessionStorage.getItem("diary_id");
let entry_id = sessionStorage.getItem("entry_id");



const file_input = document.getElementById("file_add")
file_input.addEventListener('input', ()=>{
    add_file()
})

document.documentElement.style.cssText = `--back_diary: ${sessionStorage.getItem("bg_colour")} ;
    --back2_diary: ${sessionStorage.getItem("bg_colour2")} ;
    --block_diary: ${sessionStorage.getItem("block_colour")};
    --text_diary: ${sessionStorage.getItem("text_colour")};
    --font_diary: ${sessionStorage.getItem("font")}`;

let data_files;
 get_diary();
function for_existing(){
   
    get_entry();
    get_entry_files();
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
async function get_diary() {
    
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
            throw new Error();
        } 
        let data = await response.json();
        
        document.getElementById("name").textContent = data.alter_name+"'s Diary";

    }
    
    catch(error){
        console.error(error);
    }

}



async function get_entry() {
    
    try{    
        let token = sessionStorage.getItem("token")
        let response = await fetch(`${api_link}/diaries/entries/${diary_id}/${entry_id}`, {
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
    
        document.getElementById("date_of_entry").textContent = moment(data.date).utc().local().format('LLL')
        document.getElementById("name_of_entry").value = data.name
        quill.clipboard.dangerouslyPasteHTML(data.text, `silent`)
        document.getElementById("p_name_of_entry").textContent = data.name
        document.getElementById("p_description_entry").innerHTML = data.text

        // document.getElementById("editor").innerHTML = 
        if (sessionStorage.getItem("scroll") != null){
            window.scroll(0, sessionStorage.getItem("scroll"))
            sessionStorage.removeItem("scroll")
        }
    }
    catch(error){
        console.error(error);
    }

}

async function get_entry_files() {
    
    try{    
        let token = sessionStorage.getItem("token")
        
        let response = await fetch(`${api_link}/diaries/files/${entry_id}`, {
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
        data_files = data
        let files = document.getElementById("photos")
        files.innerHTML = ""
        for (let i = 0; i < data.length; i++){
            
            let image = `<div class="file"><img class="file_img" src="${"../../../"+data[i].file}" >
            <button class="close_file_btn" onclick="download_file(${i})">download</button>
                        <button class="close_file_btn" onclick="delete_file(${i})">&#10005;</button></div>`
            files.innerHTML += image
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

function open_edit(){
    document.getElementById("p_name_of_entry").style.display = "none";
    document.getElementById("name_of_entry").style.display = "block";
    document.getElementById("p_description_entry").style.display = "none";
    document.getElementById("editor").style.display = "block";
    document.getElementById("edit_entry").style.display = "block";
    document.getElementById("open_edit").style.display = "none";
    document.getElementsByClassName("ql-toolbar")[0].style.display = "block"
}


async function edit_entry() {
    
    try{
 
        let token = sessionStorage.getItem("token");
        let name = document.getElementById("name_of_entry").value;
        // let text = document.getElementById("editor").innerHTML;
        let text = quill.getSemanticHTML(0);
        let response = await fetch(`${api_link}/diaries/entries/${diary_id}/${entry_id}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
            name: name,
            text: text
        })
  
    });

        if (!response.ok){
            throw new Error();
        } 
        

            
        document.getElementById("p_name_of_entry").style.display = "block";
        document.getElementById("name_of_entry").style.display = "none";
        document.getElementById("p_description_entry").style.display = "block";
        document.getElementById("editor").style.display = "none";
        document.getElementById("edit_entry").style.display = "none";
        document.getElementById("open_edit").style.display = "block";
        document.getElementsByClassName("ql-toolbar")[0].style.display = "none"

        // console.log("hey")
        get_entry();
        


    }
    catch(error){
        console.error(error);
    }

}

async function edit_entry_adding() {
    
    try{
 
        let token = sessionStorage.getItem("token");
        let name = document.getElementById("name_of_entry").value;
        // let text = document.getElementById("editor").innerHTML;
        let text = quill.getSemanticHTML(0);
        if (token == "" || name==""){
            throw new Error("empty");
        }
        let response = await fetch(`${api_link}/diaries/entries/${diary_id}/${entry_id}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
            name: name,
            text: text
        })
  
    });

        if (!response.ok){
            throw new Error();
        } 
        

       document.location.href="../entry"
        


    }
    catch(error){
        console.error(error);
    }

}

async function download_file(i) {
    let filename = data_files[i].file 
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

async function delete_file(i) {
    
    try{    
        let token = sessionStorage.getItem("token");
        let file_id = data_files[i].id
        let response = await fetch(`${api_link}/diaries/files/${file_id}`, {
        method: "DELETE",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
  
    });

        if (!response.ok){
            throw new Error();
        } 
        get_entry_files();


    }
    catch(error){
        console.error(error);
    }

}
async function add_file() {
    
    try{    
        let token = sessionStorage.getItem("token");
        const fileInput = document.getElementById('file_add');
        const file = fileInput.files[0];
        const formData = new FormData();
        formData.append('photo', file);
        let response = await fetch(`${api_link}/diaries/files/${entry_id}`, {
        method: "POST",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: formData
  
    });

        if (!response.ok){
            throw new Error();
        } 
        get_entry_files();


    }
    catch(error){
        console.error(error);
    }

}


// Create FormData and append the file



async function add_entry() {
    
    try{    
        let token = sessionStorage.getItem("token");
        
        let response = await fetch(`${api_link}/diaries/entries/${diary_id}`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
            name: "StandardInvisibleEntryofkw",
            text: "StandardInvisibleEntryofkw"
        })
  
    });

        if (!response.ok){
            throw new Error();
        } 
        let data = await response.json();
        sessionStorage.setItem("entry_id", data.id);
        entry_id = sessionStorage.getItem("entry_id");
        quill.clipboard.dangerouslyPasteHTML("", `silent`);
        document.getElementById("editor").style.display = "block";
        document.getElementsByClassName("ql-toolbar")[0].style.display = "block"
        // get_entry_files();
        


    }
    catch(error){
        console.error(error);
    }

}