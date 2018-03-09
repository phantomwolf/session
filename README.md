# Session
HTTP协议是无状态的。为了记录用户的登录状态，我们需要session。

## 原理
当用户登录时，为其生成一个session，存放在服务器上(如redis中)。每个session拥有一个唯一的id(采用UUID)，并包含uid、name等信息。服务器生成session后，利用Set-Cookie头，将SESSION_ID写入客户端的Cookie中。这样，以后客户端的请求都会带上这个Cookie。服务器用SESSION_ID搜索相应的Session数据，若session存在且未过期，说明用户依然处于登录状态；若session不存在或已过期，说明用户需要重新登录，将其跳转到登录页面。

## 存储
### redis
redis中存储两样东西，一样是session本身，以redis hash的形式存储：
<session_id> => {"_uid": "<uid>", "_expire": "<expire>", ...}

一种是单个用户所有的session id，以redis list的形式存储：
<user_id> => [<session_id>, <session_id>]

通过限制第二个list的长度，我们可以限制一个用户可拥有的session数量。
