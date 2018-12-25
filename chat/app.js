var me = {};
me.avatar = "https://en.gravatar.com/userimage/142709731/db51ed5f00030f48c8f15cae951b46ac.png";

var you = {};
you.avatar = "https://en.gravatar.com/userimage/142709731/8aa2297829b3ecd61f14a7ed89f39acc.png";

var clientID;
var channelID = "chatbot";

function formatAMPM(date) {
    var hours = date.getHours();
    var minutes = date.getMinutes();
    var ampm = hours >= 12 ? 'PM' : 'AM';
    hours = hours % 12;
    hours = hours ? hours : 12; // the hour '0' should be '12'
    minutes = minutes < 10 ? '0'+minutes : minutes;
    var strTime = hours + ':' + minutes + ' ' + ampm;
    return strTime;
}

//-- No use time. It is a javaScript effect.
function insertChat(who, text, time){
    if (time === undefined){
        time = 0;
    }
    var control = "";
    var date = formatAMPM(new Date());

    if (who == "me"){
        control = '<li style="width:100%">' +
            '<div class="msj macro">' +
            '<div class="avatar"><img class="img-circle" style="width:100%;" src="'+ me.avatar +'" /></div>' +
            '<div class="text text-l">' +
            '<p>'+ text +'</p>' +
            '<p><small>'+date+'</small></p>' +
            '</div>' +
            '</div>' +
            '</li>';
    }else{
        control = '<li style="width:100%;">' +
            '<div class="msj-rta macro">' +
            '<div class="text text-r">' +
            '<p>'+text+'</p>' +
            '<p><small>'+date+'</small></p>' +
            '</div>' +
            '<div class="avatar" style="padding:0px 0px 0px 10px !important"><img class="img-circle" style="width:100%;" src="'+you.avatar+'" /></div>' +
            '</li>';
    }
    setTimeout(
        function(){
            $("ul").append(control).scrollTop($("ul").prop('scrollHeight'));
        }, time);

}

function resetChat(){
    $("ul").empty();
}


$(document).ready(function(){

    // Create Centrifuge object with Websocket endpoint address set in main.go
    var centrifuge = new Centrifuge('ws://localhost:8080/connection/websocket');

    centrifuge.on('connect', function(ctx) {
        //drawText('Connected over ' + ctx.transport + '<br>');
        console.log('connected')
        clientID = ctx.client;
    });

    centrifuge.on('disconnect', function(ctx) {
        //drawText('Disconnected: ' + ctx.reason + '<br>');
        console.log('disconnected')
    });

    var sub = centrifuge.subscribe(channelID, function(message) {

        console.log(message.info.client)
        console.log(clientID)

        if(message.info.client == clientID) {
            insertChat("me", message.data, 0)
        } else {
            insertChat("you", message.data, 0)
        }

    });

    centrifuge.connect();



    $(".mytext").on("keydown", function(e){
        if (e.which == 13){
            var text = $(this).val();
            if (text !== ""){
                sub.publish(text);
                $(this).val('');
            }
        }
    });
});

$('body > div > div > div:nth-child(2) > span').click(function(){
    $(".mytext").trigger({type: 'keydown', which: 13, keyCode: 13});
})

//-- Clear Chat
//resetChat();

//-- Print Messages
// insertChat("me", "Hello Tom...", 0);
// insertChat("you", "Hi, Pablo", 1500);
// insertChat("me", "What would you like to talk about today?", 3500);
// insertChat("you", "Tell me a joke",7000);
// insertChat("me", "Spaceman: Computer! Computer! Do we bring battery?!", 9500);
// insertChat("you", "LOL", 12000);


//-- NOTE: No use time on insertChat.


// var input = document.getElementById("input");
// input.addEventListener('keyup', function(e) {
//     if (e.keyCode == 13) { // ENTER key pressed
//         sub.publish(this.value);
//         input.value = '';
//     }
// });

// After setting event handlers – initiate actual connection with server.
//