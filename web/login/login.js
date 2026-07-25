const env_link = "http://localhost:8080"
const api_link = "http://localhost:8080/api"



let Email = document.getElementById('UserEmail').value.toLowerCase();
let Password = document.getElementById('UserPassword').value;
let sysname_el = document.getElementById('usn');


let register = false;
// function login_pressed(){
//     register = false;
//     sysname_el.style.display = "none";
// }
function reg_pressed(){
    register = true;
    sysname_el.style.display = "block";

    document.getElementById('register').style.display = "none"
    document.getElementById('btn_submit').textContent = "SIGN IN"
}


function submit(){
    if (register){
        Post_user();
        return;
    }
    Login();
}



async function Login() {
    

    let Email = document.getElementById('UserEmail').value.toLowerCase();
    let Password = document.getElementById('UserPassword').value;

    try{

    let response = await fetch(`${api_link}/login`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            email: Email,
            password: Password
        })
    });

        if (!response.ok){
            throw new Error();
        } 
        let data = await response.json();
        sessionStorage.setItem('token', data.token);
        localStorage.setItem('refresh_token', data.refresh_token);
        sessionStorage.setItem('userId', data.id)
        window.location.href=`${env_link}/app/`

    }
    catch(error){
        console.error(error);
    }

}
async function Post_user() {
    
    try{


    let Email = document.getElementById("UserEmail").value.toLowerCase();
    let Password = document.getElementById("UserPassword").value;
    let Sysname = document.getElementById("UserSystemName").value;    
    
    let response = await fetch(`${api_link}/users`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            email: Email,
            password: Password,
            system_name: Sysname
        })
    });

        if (!response.ok){
            throw new Error();
        } 


        Login()
    }
    catch(error){
        console.error(error);
    }

}


// document.getElementById('btn_login').addEventListener('click', login_pressed);
// document.getElementById('register').addEventListener('click', reg_pressed);
// document.getElementById('btn_submit').addEventListener('click', submit);


