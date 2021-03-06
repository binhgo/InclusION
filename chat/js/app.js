var me = {};
me.avatar = "images/me.png";

var you = {};
you.avatar = "images/you.png";

var clientID;
var channelID = "chatbot";
var name = '';

var image = "images/image.png";

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

    if (who === "me") {
        control = '<li>' +
            '<div class="d-flex">' +
            '<img src="' + me.avatar + '" alt="avatar" class="rounded-circle">' +
            '<div class="text-l">' +
            '<p class="message-text">'+ text +'</p>' +
            '<p><small>'+ date +'</small></p>' +
            '</div>' +
            '</div>' +
            '</li>'
    }else{
        control = '<li>' +
            '<div class="d-flex float-right mr-4">' +
            '<div class="text-r">' +
            '<p class="message-text">'+ text +'</p>' +
            '<p><small class="float-right">'+ date +'</small></p>' +
            '</div>' +
            '<img src="'+ you.avatar +'" alt="avatar" class="rounded-circle">' +
            '</div>' +
            '</li>';
    }
    setTimeout(
        function(){
            $("ul.message_body").append(control).scrollTop($("ul.message_body").prop('scrollHeight'));
        }, time);

}

function resetChat(){
    $("ul.message_body").empty();
}

function checkUrl(url) {
    var expression = /(https?:\/\/(?:www\.|(?!www))[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]\.[^\s]{2,}|www\.[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]\.[^\s]{2,}|https?:\/\/(?:www\.|(?!www))[a-zA-Z0-9]\.[^\s]{2,}|www\.[a-zA-Z0-9]\.[^\s]{2,})/;
    var regex = new RegExp(expression);
    if (url.match(regex))
        return true;
    return false;
}

function initItem(data) {
    var item = "<li>" +
        "<div class=\"d-flex float-right mr-4\">" +
        "<div class=\"text-r\">" +
        "<div class=\"row mt-1 mb-1 text-center\">";

    data.forEach(val => {
        var img = checkUrl(val.Image) ? val.Image : image;

        item += "<div class=\"col-lg-3 col-12 suggest-item\">" +
            "<img src=\""+ img +"\" class=\"img-fluid suggest-img\">" +
            "<a href='https://www.lazada.vn/products/apple-iphone-xs-max-i248162747-s324087916.html' target='_blank'>" +
            "<h5 class=\"suggest-title\">"+ val.Name +"</h5>" +
            "</a>" +
            "<p class=\"suggest-description\">" + val.Description + "</p>" +
            "</div>";

        $(".suggest-items").prepend(item);
    });
    item += "</div>" +
        "</div>" +
        "<img src=\"images/you.png\" alt=\"avatar\" class=\"rounded-circle\">" +
        "</div>" +
        "</li>";

    $("ul.message_body").append(item).scrollTop($("ul.message_body").prop('scrollHeight'));
}

function getResponse(data) {
    switch (data.Type) {
        case "Events":
            return data.Events;
            break;
        case "Products":
            return data.Products;
            break;
        case "Suggestions":
            return data.Suggestions;
            break;
        default:
            return data.Mess;
    }
}

$(document).ready(function(){

    // Create Centrifuge object with Websocket endpoint address set in main.go
    var centrifuge = new Centrifuge('ws://'+ window.location.hostname +':'+ window.location.port +'/connection/websocket');
	// var centrifuge = new Centrifuge('ws://localhost:8080/connection/websocket');

    centrifuge.on('connect', function(ctx) {
        console.log('connected');
        clientID = ctx.client;
    });

    centrifuge.on('disconnect', function(ctx) {
        console.log('disconnected')
    });

    if (channelID === 'chatbot') {
        $("#chat_content").load("../chat/register.html?v=2", "", function () {

            var sub = centrifuge.subscribe(channelID, function (message) {

                if (typeof message.data === 'string') {

                    if (JSON.parse(message.data).ClientId === clientID) {
                        $("#chat_content").load("../chat/chat.html?v=1", "", function () {
                            $("#your_name").html(name);

                            channelID = JSON.parse(message.data).ChannelId;

                            var sub1 = centrifuge.subscribe(channelID, function (message) {
                                if (message.info.client !== clientID) {
                                    var resData = JSON.parse(message.data);
                                    console.log(resData);
                                    if (resData.Type === "Message") {
                                        insertChat("you", resData.Mess, 0);
                                    } else {
                                        var data = getResponse(resData);
                                        setTimeout(function () {
                                            initItem(data);
                                        }, 1000);
                                    }
                                }
                            });

                            $(document).on("keydown", "#my_text", function (e) {
                                if (e.which === 13) {
                                    var text = $(this).val();
                                    if (text !== "") {
                                        insertChat("me", text, 0);
                                        $(".suggest-items").html("");
                                        sub1.publish(text);
                                        $(this).val('');
                                    }
                                }
                            });
							
							$(document).on("click", "#btn_mytext", function () {
								$("#my_text").trigger({type: 'keydown', which: 13, keyCode: 13});
							});
                        });
                    }
                }
            });

            $(document).on("submit", "#form-register", function (e) {
                e.preventDefault();
                var form_value = $(this).serializeArray();
                name = form_value[0].value;
                var info = {name: form_value[0].value, email: form_value[1].value};
                sub.publish(info);
            });

        });
    }

    centrifuge.connect();
});
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
