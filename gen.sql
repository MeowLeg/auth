create table if not exists clicks (
	ip text,
	logtime datetime default (datetime('now', 'localtime'))
);

create table if not exists register (
	openid text primary key,
	nickname text not null
);

create table if not exists project (
	key text primary key,
	url text not null,
	weixin text default ''
);


-- alter table project add column weixin text default '';

create table if not exists weixin (
	weixin text unique,
	appid text,
	appsecret text,
	access_token text
);

-- insert into weixin values ('zstvjbnt', 'wx8e3e7383c989c94a', '6547db3da6a229a659bda118f48f8038', '');
-- insert into weixin values ('zsgd93', 'wxee6284b0c702c21e', '86707f6633745e51c4639f11c82abed2', '');
-- update project set weixin = 'zstvjbnt' where key = 'choice1';
-- insert into project values('choice2', 'http://develop.zsgd.com:11002/choices/wechat/index.html?vote_id=16', 'zstvjbnt');
