
let th = localStorage.getItem("theme")
if (th != null){
    document.documentElement.style.setProperty('color-scheme', th)

}



window.addEventListener('pageshow', function(event){
    if (event.persisted){
        this.sessionStorage.setItem("scroll", this.window.scrollY)
        window.location.reload();
    }

});



const env_link = "https://plurserv.fly.dev"
const api_link = "https://plurserv.fly.dev/api"
sessionStorage.removeItem("entry_id")
redirect_to_login();
const rootStyles = getComputedStyle(document.documentElement)
const body = document.getElementById("diary_body")

let alter_id = sessionStorage.getItem("alter_id")
let diary_id = sessionStorage.getItem("diary_id")
let data_entries;


const bg_colour_picker =  document.getElementById("bg_colour")
const bg_colour2_picker = document.getElementById("bg_colour2")
const block_colour_picker = document.getElementById("block_colour")
const text_colour_picker = document.getElementById("text_colour")
const font_picker = document.getElementById("font")

bg_colour_picker.addEventListener('input', () => {
    colour_hex({
        bg_colour: bg_colour_picker.value,
        bg_colour2: bg_colour2_picker.value,
        block_colour: block_colour_picker.value,
        text_colour: text_colour_picker.value,
        font: font_picker.value
})
})
bg_colour_picker.addEventListener('input', () => {
    colour_hex({
        bg_colour: bg_colour_picker.value,
        bg_colour2: bg_colour2_picker.value,
        block_colour: block_colour_picker.value,
        text_colour: text_colour_picker.value,
        font: font_picker.value
})
})
bg_colour2_picker.addEventListener('input', () => {
    colour_hex({
        bg_colour: bg_colour_picker.value,
        bg_colour2: bg_colour2_picker.value,
        block_colour: block_colour_picker.value,
        text_colour: text_colour_picker.value,
        font: font_picker.value
})
})
block_colour_picker.addEventListener('input', () => {
    colour_hex({
        bg_colour: bg_colour_picker.value,
        bg_colour2: bg_colour2_picker.value,
        block_colour: block_colour_picker.value,
        text_colour: text_colour_picker.value,
        font: font_picker.value
})
})
text_colour_picker.addEventListener('input', () => {
    colour_hex({
        bg_colour: bg_colour_picker.value,
        bg_colour2: bg_colour2_picker.value,
        block_colour: block_colour_picker.value,
        text_colour: text_colour_picker.value,
        font: font_picker.value
})
})
font_picker.addEventListener('input', () => {
    colour_hex({
        bg_colour: bg_colour_picker.value,
        bg_colour2: bg_colour2_picker.value,
        block_colour: block_colour_picker.value,
        text_colour: text_colour_picker.value,
        font: font_picker.value
})
})





get_diary();
get_entries();

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
            window.location.href=`${env_link}/app/login`
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
        console.log(data)
        sessionStorage.setItem("diary_id", data.id)
        document.getElementById("name").textContent = data.alter_name+"'s Diary";
        
        document.getElementById("font").value = data.font;
        colour(data)
    }
    
    catch(error){
        console.error(error);
    }

}



function colour_hex(data){
    let bg_colour = data.bg_colour;
    let bg_colour2 = data.bg_colour2;
    let bg_colour2_rgb = [parseInt(bg_colour2.slice(1, 3), 16), parseInt(bg_colour2.slice(3, 5), 16), parseInt(bg_colour2.slice(5, 7), 16)]
    let bg_colour_rgb = [parseInt(bg_colour.slice(1, 3), 16), parseInt(bg_colour.slice(3, 5), 16), parseInt(bg_colour.slice(5, 7), 16)]

    let block_colour = data.block_colour;
    let text_colour = data.text_colour;
    let del = document.getElementById("header_diary_div").offsetHeight/window.screen.height

    // console.log(bg_colour)
    document.documentElement.style.cssText = `--back_diary: ${bg_colour};
    --back2_diary: ${bg_colour2};
    --back3_diary: rgba(${parseInt(bg_colour2_rgb[0]*del+bg_colour_rgb[0]*(1-del))}, ${parseInt(bg_colour2_rgb[1]*del+bg_colour_rgb[1]*(1-del))}, ${parseInt(bg_colour2_rgb[2]*del+bg_colour_rgb[2]*(1-del))});

    --block_diary: ${block_colour};
    --text_diary: ${text_colour};
    --font_diary: ${data.font}`;


    
    // rootStyles.setProperty('--back_diary', bg_colours) 
}


function colour(data){
    let bg_colour = data.bg_colour.split(" ");
    let bg_colour2 = data.bg_colour2.split(" ");
    let block_colour = data.block_colour.split(" ");
    let text_colour = data.text_colour.split(" ");
    let del = document.getElementById("header_diary_div").offsetHeight/window.screen.height

    // console.log(bg_colour)


    
    document.documentElement.style.cssText = `--back_diary: rgba(${parseInt(bg_colour[0])}, ${parseInt(bg_colour[1])}, ${parseInt(bg_colour[2])});
    --back2_diary: rgba(${parseInt(bg_colour2[0])}, ${parseInt(bg_colour2[1])}, ${parseInt(bg_colour2[2])});
    --back3_diary: rgba(${parseInt(bg_colour2[0]*del+bg_colour[0]*(1-del))}, ${parseInt(bg_colour2[1]*del+bg_colour[1]*(1-del))}, ${parseInt(bg_colour2[2]*del+bg_colour[2]*(1-del))});

    --block_diary: rgb(${parseInt(block_colour[0])}, ${parseInt(block_colour[1])}, ${parseInt(block_colour[2])});
    --text_diary: rgb(${parseInt(text_colour[0])}, ${parseInt(text_colour[1])}, ${parseInt(text_colour[2])});
    --font_diary: ${data.font}`;

    document.getElementById("bg_colour").value = rootStyles.getPropertyValue('--back_diary');
    document.getElementById("bg_colour2").value = rootStyles.getPropertyValue('--back2_diary');
    document.getElementById("block_colour").value = rootStyles.getPropertyValue('--block_diary');
    document.getElementById("text_colour").value = rootStyles.getPropertyValue('--text_diary');
    
    sessionStorage.setItem("bg_colour", rootStyles.getPropertyValue('--back_diary'));
    sessionStorage.setItem("bg_colour2", rootStyles.getPropertyValue('--back2_diary'));
    sessionStorage.setItem("block_colour", rootStyles.getPropertyValue('--block_diary'));
    sessionStorage.setItem("text_colour", rootStyles.getPropertyValue('--text_diary'));
    sessionStorage.setItem("font", data.font);

    
    // rootStyles.setProperty('--back_diary', bg_colours) 
}

async function edit_diary() {
    
    try{
        let token = sessionStorage.getItem("token")
        let bg_colour = document.getElementById("bg_colour").value;
        let bg_colour2 = document.getElementById("bg_colour2").value;
        let block_colour = document.getElementById("block_colour").value;
        let text_colour = document.getElementById("text_colour").value;
        let font = document.getElementById("font").value;
        let response = await fetch(`${api_link}/diaries/${diary_id}`, {
        method: "PUT",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
            bg_colour: bg_colour,
            bg_colour2: bg_colour2,
            block_colour: block_colour,
            text_colour: text_colour,
            font: font
        })
  
    });

        if (!response.ok){
            throw new Error();
        } 
        document.getElementById("adding_popover_settings").style.display = "none"
        get_diary()
    }
    
    catch(error){
        console.error(error);
    }

}

async function delete_diary() {
    
    try{
        let token = sessionStorage.getItem("token")
        let response = await fetch(`${api_link}/diaries/${diary_id}`, {
        method: "DELETE",
        headers: {
            // "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
  
    });

        if (!response.ok){
            throw new Error();
        } 
        window.location.href="diaries"
    }
    
    catch(error){
        console.error(error);
    }

}


async function get_entries() {
    
    try{    
        let token = sessionStorage.getItem("token")
        let response = await fetch(`${api_link}/diaries/entries/${diary_id}`, {
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
        data_entries = data
        let entries = document.getElementById("entries")

        entries.innerHTML = ""
        if (data.length == 0){
            document.getElementById('no_one').style.display = 'flex';
            return;
        }else{
            document.getElementById('no_one').style.display = 'none';

        }
        for (let i = 0; i < data.length; i++){
            let entry = `
                    <div class="diary_entry" onclick="open_entry(${i})" id="entry" style=" margin-top: 10px;">
                        <div id="header_entry" class="entry_header">
                            <p id="name_of_entry">${data[i].name}</p>
                            <p id="date_of_entry">${moment(data[i].date).utc().local().format('LLL') }</p>
                        </div>
                        <div id="body_entry" class="entry_body_text">
                            <div id="description_entry" class="description_entry" >${data[i].text}</div>
                            <a id="entry_link" href="javascript:open_entry(${i})">Read more and see photos</a>
                        </div>
                    </div>
            `     
            entries.innerHTML += entry
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

function open_entry(ind){
    sessionStorage.setItem("entry_id", data_entries[ind].id)
    window.location.href="entry"
}

