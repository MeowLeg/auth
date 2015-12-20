![新蓝广科](http://develop.zsgd.com:8081/markdown/img/zsgd.jpg)

##auth简介

auth是用于微信公众号网页认证的系统，采用go-template微服务框架编写。

##使用配置

* 配置middle.db数据库的weixin表。数据库中的weixin表存储公众号信息，四个字段分别是公众号、appid、appsecret、access_token。其中access_token需要定时获取，这里已经写了一个timerTask.go，用于每小时更新一次。

    ```sql
    create table if not exists weixin (
        weixin text unique,
    	appid text,
    	appsecret text,
    	access_token text
    );
    ```
* 配置project表。该表字段分别是网页验证的state字段（用户自定义），验证后的重定向url，所属的微信公众号（用户验证是否是公众号用户）。

    ```sql
    create table if not exists project (
        key text primary key,
        url text not null,
    	weixin text default ''
    );
    ```
* switch/switch.go中ifSubscribe函数是验证公众号用户之用，可以在其中添加无须认证的公众号，例子中是zsgd93。

##开发承建
[新蓝广科](http://www.xinlantech.com)
