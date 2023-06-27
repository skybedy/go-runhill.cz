
    

             

              
       
            



$(document).ready(function(){
    //$(this).scrollTop(5000);
   // $('body,html').animate({scrollTop: 156}, 800); 
   //window.scrollTo(2000, 2000);

        $('buttonxxx').on('click', function () {
            var center = SMap.Coords.fromWGS84(14.41790, 50.12655);
            var m = new SMap(JAK.gel("m"), center, 13);
            m.addControl(new SMap.Control.Sync()); /* Aby mapa reagovala na změnu velikosti průhledu */
            m.addDefaultLayer(SMap.DEF_TURIST).enable(); /* Turistický podklad */
        
            var mouse = new SMap.Control.Mouse(SMap.MOUSE_PAN | SMap.MOUSE_WHEEL | SMap.MOUSE_ZOOM); /* Ovládání myší */
            m.addControl(mouse);
        
            var sync = new SMap.Control.Sync();
            sync.setBottomSpace(10);
            m.addControl(sync);
            
            var xhr = new JAK.Request(JAK.Request.XML);
            xhr.setCallback(window, "response");
            //xhr.send("//api.mapy.cz/xml/sample.xml");
            xhr.send("https://api.timechip.cz/uploads/test.xml");

            var response = function(xmlDoc) {
                var gpx = new SMap.Layer.GPX(xmlDoc);
                m.addLayer(gpx);
                gpx.enable();
                gpx.fit();
                
            }
            return false;   
        }); 
       








    $('butto').on('click', function () {
        const data =  new FormData($('#test_form')[0]);
        $.ajax({
            method: 'post',
            processData: false,
            contentType: false,
             cache: false,
            data: data,
            enctype: 'multipart/form-data',
            //url: 'http://localhost:8000'
            url: 'http://api.timechip.loc/virtual-race',
            dataType: "json",
            success : function(data){
  

                var center = SMap.Coords.fromWGS84(14.41790, 50.12655);
                var m = new SMap(JAK.gel("m"), center, 13);
                m.addControl(new SMap.Control.Sync()); /* Aby mapa reagovala na změnu velikosti průhledu */
                m.addDefaultLayer(SMap.DEF_TURIST).enable(); /* Turistický podklad */
            
                var mouse = new SMap.Control.Mouse(SMap.MOUSE_PAN | SMap.MOUSE_WHEEL | SMap.MOUSE_ZOOM); /* Ovládání myší */
                m.addControl(mouse);
            
                var sync = new SMap.Control.Sync();
                sync.setBottomSpace(10);
                m.addControl(sync);
                
    var xhr = new JAK.Request(JAK.Request.XML);
    xhr.setCallback(window, "response");
    xhr.send("//api.mapy.cz/xml/sample.xml");
    //xhr.send("https://api.timechip.cz/uploads/test.xml");

    var response = function(xmlDoc) {
        var gpx = new SMap.Layer.GPX(xmlDoc);
        m.addLayer(gpx);
        gpx.enable();
        gpx.fit();
    }
            }
	    }); 
        return false;
    })










})




/* pure JavaScript
var form = document.querySelector('form');
var request = new XMLHttpRequest();
form.addEventListener('submit',function(e){
    e.preventDefault();
    var formdata = new FormData(form);
    console.log(formdata);
    request.open('POST','test.php')
    request.send(formdata);
},false)*/
var center = SMap.Coords.fromWGS84(14.41790, 50.12655);
var m = new SMap(JAK.gel("m"));

var xhr = new JAK.Request(JAK.Request.XML);
xhr.setCallback(window, "response");




var form = document.querySelector('form');

form.addEventListener('submit',function(e){
    e.preventDefault();


            
            
            xhr.send("//api.mapy.cz/xml/sample.xml");
            var signals = m.getSignals();
            signals.addListener(window, "marker-click", function() { 
				    
				});




            xhr.send("https://api.timechip.cz/uploads/test.xml");
            
            var response = function(xmlDoc) {
                var gpx = new SMap.Layer.GPX(xmlDoc);
                m.addLayer(gpx);
                gpx.enable();
                gpx.fit();
                
            }

},false)

var center = SMap.Coords.fromWGS84(14.41790, 50.12655);
var m = new SMap(JAK.gel("m"), center, 13);
m.addControl(new SMap.Control.Sync()); /* Aby mapa reagovala na změnu velikosti průhledu */
m.addDefaultLayer(SMap.DEF_TURIST).enable(); /* Turistický podklad */
var signals = m.getSignals();
signals.addListener(window, "marker-click", function(e) {
   alert();
   return false;
  });





