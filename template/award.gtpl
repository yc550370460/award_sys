<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>抽奖</title>
    <link rel="stylesheet" href="/static/css/bootstrap.min.css">
    <link rel="stylesheet" href="/static/css/bootstrap-theme.min.css">
    <link rel="stylesheet" href="/static/css/style.css">
    <link rel="stylesheet" href="/static/css/barrager.css">

    <script src="/static/js/jquery-2.2.1.min.js"></script>
    <script src="/static/js/header.js"></script>
    <script src="/static/js/jquery.barrager.js"></script>
</head>
<body>
    <div class='luck-back'>
        <p id="conn_status">connecting...</p>
        <div class="luck-content ce-pack-end">
            <div id="luckuser" class="slotMachine">
                <div class="slot">
                    <span class="name">姓名</span>
                </div>
            </div>
            <div class="luck-content-btn">
                <a id="start" class="start" onclick="start()">开始</a>
            </div>
            <div class="luck-user">
                <div class="luck-user-title">
                    <span>中奖名单</span>
                </div>
                <ul class="luck-user-list"></ul>
                <div class="luck-user-btn">
                    <a href="#">中奖人</a>
                </div>
            </div>
        </div>
    </div>


    <script>
        var xinm = new Array();
        xinm[0] = "/static/img/employee_ico/1.jpg"
        xinm[1] = "/static/img/employee_ico/2.jpg"
        xinm[2] = "/static/img/employee_ico/3.jpg"
        xinm[3] = "/static/img/employee_ico/4.jpg"
        xinm[4] = "/static/img/employee_ico/5.jpg"
        xinm[5] = "/static/img/employee_ico/6.jpg"
        xinm[6] = "/static/img/employee_ico/1.jpg"
        xinm[7] = "/static/img/employee_ico/2.jpg"
        xinm[8] = "/static/img/employee_ico/3.jpg"
        xinm[9] = "/static/img/employee_ico/4.jpg"
        xinm[10] = "/static/img/employee_ico/5.jpg"
        xinm[11] = "/static/img/employee_ico/6.jpg"

        var phone = new Array();
        phone[0] = "王尼玛"
        phone[1] = "张全蛋"
        phone[2] = "纸巾老撕"
        phone[3] = "王铁柱"
        phone[4] = "田二妞"
        phone[5] = "唐马儒"
        phone[6] = "张三"
        phone[7] = "李四"
        phone[8] = "王二麻子"
        phone[9] = "咯咯咯"
        phone[10] = "一二三"
        phone[11] = "四五六"

        let sock = null;
        const wsuri = "ws://127.0.0.1:1234/bubble";

        console.log("onload");
        sock = new WebSocket(wsuri);
        sock.onopen = function() {
            console.log("connected to " + wsuri);
        }
        sock.onclose = function(e) {
            console.log("connection closed (" + e.code + ")");
        }
        var msg_list = new Array();
        sock.onmessage = function(e) {
            // console.log("message received: " + e.data);
            // document.getElementById('recv_msg').innerHTML=e.data;
            msg_list.push(e.data);
            while (msg_list.length > 0)
                {
                    msg = msg_list.pop();
                    console.log(msg);
                    ico_list = ["/static/img/icon/blueE.ico", "/static/img/icon/blueN.ico", "/static/img/icon/bubble.ico", "/static/img/icon/cheers.ico",
                    "/static/img/icon/colonA.ico", "/static/img/icon/fiveX.ico", "/static/img/icon/keyboardE.ico", "/static/img/icon/party.ico",
                    "/static/img/icon/redA.ico", "/static/img/icon/redE.ico", "/static/img/icon/redN.ico", "/static/img/icon/X storage.ico", "/static/img/icon/blueX.ico"]
                    icon = ico_list[Math.floor(Math.random()*12)]
                    console.log(icon)
                    if (msg.indexOf("2018") > -1 || msg.indexOf("新年") > -1)
                        {
                            console.log("match")
                            $(".barrage").css("width", "500px");
                            $(".barrage_box").css("height", "100px");
                            $(".barrage_box div.p a").css("font-size", "25px");
                            $(".barrage_box div.p").css("margin", "9% 0px");
                            console.log(e.data.length);
                            if (e.data.length > 16)
                            {
                                console.log("too long")

                                danmu(e.data.slice(0, 16), icon);
                                danmu("程序员：太长被我截尾了", icon);
                            }
                            else{
                                danmu(e.data, icon);
                            }

                        }
                        else
                        {
                            console.log("not match")
                            $(".barrage").css("width", "500px");
                            $(".barrage_box").css("height", "40px");
                            $(".barrage_box div.p a").css("font-size", "14px");
                            $(".barrage_box div.p").css("margin", "0px 0px 0px 0px");
                            if (e.data.length > 16)
                            {
                                console.log("too long")
                                danmu(e.data.slice(0, 16), icon);
                                danmu("程序员：太长被我截尾了", icon);
                            }
                            else{
                                danmu(e.data, icon);
                            }
                        }
                }

        }


        // 弹幕
        function danmu(msg, icon) {
            var item={
                img: icon, //图片
                info:msg, //文字
                href:'http://www.baidu.com', //链接
                close:true, //显示关闭按钮
                speed:10, //延迟,单位秒,默认6
                color:'#ffffff', //颜色,默认白色
                bottom: Math.random() * 800,
                old_ie_color:'#000000', //ie低版兼容色,不能与网页背景相同,默认黑色
            }
            $('body').barrager(item);

        }
    </script>
    <script src="/static/js/Luckdraw.js" type="text/javascript"></script>
</body>
</html>
