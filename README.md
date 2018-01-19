# 工作原理
HTTP协议是无状态的。为了记录用户的登录状态，我们需要session。

当用户登录时，为其生成一个session，存放在服务器上(如redis中)。每个session拥有一个唯一的id(采用UUID)，并包含uid、name等信息。服务器生成session后，利用Set-Cookie头，将SESSION_ID写入客户端的Cookie中。这样，以后客户端的请求都会带上这个Cookie。服务器用SESSION_ID搜索相应的Session数据，若session存在且未过期，说明用户依然处于登录状态；若session不存在或已过期，说明用户需要重新登录，将其跳转到登录页面。
