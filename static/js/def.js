/*Vue.component('account-delete', {
  template: '<div class="w-25 mx-auto"><div class="p-4  border rounded-2"></div></div>'
                                                                                                                                                                                                                                                                                                              </template>'
})*/

Vue.component('account-delete', {
  template: '<div class="alert alert-danger text-center" role="alert"> Chystáte se smazat účet, myslíte-li to opravdu vážně, pokračujte <a href="#" v-on:click.prevent = "test">zde</a>.</div>',
  methods:{
    test(){
      alert()
    }
  }

})


if(document.getElementById('accountDelete') != null) {
  new Vue({
    delimiters: ['${', '}'],
    el: '#accountDelete'
  })
}



















if(document.getElementById('accountEditContainer') != null) {
  new Vue({
    delimiters: ['${', '}'],
    el: '#accountEditContainer',
    data:{
      formular: true,
      hlaska: false,
      hlaskaText: "",
      alertDanger : false,
      alertPrimary: false,
    },
    methods: {
      accountEdit(e){
        try{
          (async () => {
            
            posts = {
              firstname: e.target.firstname.value,
              surname: e.target.surname.value,
              gender: e.target.gender.value,
              birdthyear: e.target.birdthyear.value,
              ido: e.target.ido.value
            }

            const response = await  axios.post('/account-edit',posts);
            this.formular = false;
            this.hlaska = true;
            if(response.data.status == "ok"){
              this.hlaskaText = "Účet změněn";
              this.alertPrimary = true
            }else if(response.data.status == "error") {
              this.hlaskaText = "Error";
              this.alertDanger = true
            }





          })();

        }catch(err){
          console.log(err)
        }
      }

    }
  })
}







if(document.getElementById('loginContainer') != null){

  new Vue({
    delimiters: ['${', '}'],
    el: '#loginContainer',
    data:{
      formular: true,
      hlaska: false,
      hlaskaText: "",
      alertDanger : false,
      alertPrimary: false,
      loginPassive: true,
      loginActive: false,
      posts : {
        email:"",
        passwordFromForm:""
      },
    },
    methods:{
      loginUser(){
        try{
          (async () => {
            const response = await  axios.post('/login',this.posts);
            this.formular = false;
            this.hlaska = true;
            if(response.data.status == "ok"){
              switch(response.data.code){
                case 11:
                  //this.hlaskaText = 'Je to ok'
                  //this.alertPrimary = true,
                  //this.loginPassive = false,
                  //this.loginActive = true
                  document.getElementById("tlacitko").innerHTML =  '<div  class="dropdown"><button  class="btn btn-outline-primary dropdown-toggle" type="button" id="dropdownMenuButton1" data-bs-toggle="dropdown" aria-expanded="false">'+response.data.firstname+'</button><ul class="dropdown-menu" aria-labelledby="dropdownMenuButton1"><li><a class="dropdown-item" href="/account-summary">Můj účet</a></li><li><a class="dropdown-item" href="/account-delete/challenge">Smazat účet</a></li><li><a class="dropdown-item" href="/logout">Odhlásit se</a></li></ul></div>';
                  window.location = response.data.referer;
                break;
              }
            }else if(response.data.status == "error") {
              switch(response.data.code){
                case 12:
                  this.hlaskaText = "Registrace na základě emailu "+this.posts.email+" neexistuje, zkuste to znovu."
                  this.alertDanger = true
                break;
                case 13:
                  this.hlaskaText = "Účet ještě nebyl funkční, nejprve je třeba dokončit autorizaci, která vám byla zaslána emailem na adresu "+this.posts.email+"."
                  this.alertDanger = true
                break;
                case 14:
                  this.hlaskaText = "Uvedené heslo k emailu "+this.posts.email+" není správné, zkuste to znovu."
                  this.alertDanger = true
                break;
                case 15:
                  this.hlaskaText = "Uživatel  s emailem "+this.posts.email+" už sice v databázi je, ale používá přihlašování přes sociální sítě, takže je nutné se přihlásit přes ně";
                  this.alertDanger = true
                break;



                
              }
            }
          })();
        }catch(err){

        }
      }
    }
  })
}


if(document.getElementById('registrationContainer') != null){

  new Vue({
      delimiters: ['${', '}'],
      el: '#registrationContainer',
      data : {
        email: "",
        formular: true,
        hlaska: false,
        hlaskaText: "",
        alert:"",
        userExists: false,
        posts : {
          firstname:"",
          lastname:"",
          gender:"",
          birdthyear:"",
          email:"",
          password:"",
          passwordOld:"",
          passwordNew:"",
          passwordNewConfirm:""
        },

      },
      methods: {
        checkUserExists(){
          try{
            (async () => {
                const response = await  axios.get('/checkuserexists',{params: {email:this.posts.email}});
                //alert()
                  if(response.data == true){
                    this.userExists = true
                    this.hlaskaText = "Tento email již je použit a podruhé se použít nedá."
                  }
                  else{
                    this.userExists = false
                    this.hlaskaText = ""
                  }

            })();
          }catch(err){
            console.log(err);
          }
        },
        processForm(){
            try{
            (async () => {
                if(this.posts.email == ""){
                  this.posts.email = document.getElementById('inputEmail').value; //pro formular se soc siti
                }
                if(this.posts.email == ""){ //pro hypo situaci, že by přes všechny překážky nebyl zadám e-mail
                    this.formular = false;
                    this.hlaska = true;
                    this.hlaskaText = "Je potřeba zadat nějaký e-mail";
                }else{
                  const response = await  axios.post('/registration',this.posts);
                  this.formular = false;
                  this.hlaska = true;
                  if(response.data.status == "ok"){
                    switch(response.data.code){
                      case 1:
                        this.hlaskaText = 'Na email <a href="mailto:'+this.posts.email+'">'+this.posts.email+'</a> byla zaslána zpráva s autorizačním kódem, registrace bude dokončená po úspěšné autorizaci'
                      break;
                    }
                  }else if(response.data.status == "error") {
                    switch(response.data.code){
                      case 11:
                        this.hlaskaText = 'Tento email již je použit a podruhé se použít nedá.'
                      break;
                    }
                  }
                }

            })();
          }catch(err){

          }
        },
        passwordChange(){
          try{
            (async () => {
              const response = await  axios.post('/password-change',this.posts);
              this.formular = false;
              this.hlaska = true;
              if(response.data.status == "ok"){
                this.hlaskaText = "Heslo bylo úspěšně změněno"
                this.alertPrimary = true
              }else if(response.data.status == "error") {
                switch(response.data.code){
                  case 21:
                    this.hlaskaText = "Zadali jste nesprávné heslo, <a href=\"/password-change\">zkuste to prosím znovu</a>"
                    this.alertDanger = true
                  break;
                  case 22:
                    this.hlaskaText = "Nová hesla se neshodují, <a href=\"/password-change\">zkuste to prosím znovu</a>"
                    this.alertDanger = true
                  break;
                }

              }

            })();
          }catch(err){

          }
        }
    }
  })
}