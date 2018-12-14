<!DOCTYPE html>
<html>
<head>
    <title>login page</title>
    <meta charset="UTF-8">
    <link rel="stylesheet" href="/css/bootstrap.min.css">
    <link rel="stylesheet" href="/css/bootstrap-theme.min.css">
    <link rel="stylesheet" href="/css/style.css">
    <link rel="stylesheet" href="/css/barrager.css">
    <script src="/js/jquery-2.2.1.min.js"></script>
    <script src="/js/bootstrap.min.js"></script>
    <script src="/js/jquery.barrager.js"></script>
</head>
<body>
    <form action="/login" method="post">
        用户名:<input type="text" name="username">
        密码:<input type="password" name="password">
        <input type="submit" value="登录">
    </form>
    <button type="button" value="test" id="test">test</button>
    <script>
    $("#test").click(function(){
        danmu();
    })

    function danmu() {
        var item={
            img:'/img/heisenberg.png',
            info:'弹幕文字信息',
            href:'http://www.baidu.com',
            close:false,
            speed:6,
            color:'#000',
            old_ie_color:'#ffffff',
        }
        $('body').barrager(item);

    }
    </script>
</body>
</html>