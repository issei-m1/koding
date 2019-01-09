var aws = require('aws-sdk');
var request = require('request');
 
/* kintone用のパラメータ*/
var DOMAIN = 'ainmmay-kts.cybozu.com'; //kintone環境のドメイン
var APP_ID = 14;   //シンプルToDoアプリのアプリID
var BASE_URL = "https://" + DOMAIN + '/k/v1/';
var APITOKEN =  "rnU2sI32rNpKaQaszsooOw59Wcp85K5rdDI3OrCB";
var headers = {'X-Cybozu-API-Token': APITOKEN};
 
exports.handler = function(event, context) {
    var returnMessage = JSON.stringify(event);
    var temp = returnMessage.split('&');
    temp = temp[9].split('=')[1].replace(/Todo\+/g,'');
    temp = decodeURIComponent(temp);
    var title = temp.split('+')[0];
    var message = "";
    if(temp.split('+')){
        message = temp.split('+')[1];
    }
    var body_post = {
        app: APP_ID,
        record: {
            title: {
                value: title
            },
            detail: {
                value: message
            }
        }
    };
    var options_getsalesamount = {
        url: BASE_URL + 'record.json',
        method: 'POST',
        headers: headers,
        'Content-Type': 'application/json',
        json: body_post
    }
    //レコードを取得
    request(options_getsalesamount, function (error, response, body) {
        context.done(null, {text: "kintone POST success!"});
    })
};
