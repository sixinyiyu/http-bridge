## http-bridge

### 介绍
golang 实现的远程数据接口代理服务

### 特点

* golang实现，轻量、高效率
* 让不支持跨域的远程数据接口支持跨域
* 方便未备案(小程序等)、不支持http是接口支持
* 支持自定义header

### 使用

```html
http://127.0.0.1:4455/?url=https://www.mxnzp.com/api/image/girl/list?page=1&headers={"token":"1284034","deviceId":"104drldu34","appId":"1453"}

```

* 支持方法 ```GET``` ```POST```    ```DELETE``` ```PUT``` 
* ```url``` 即为需要代理的请求地址


#### 自定义headers

1. 将自定义headers写在请求头里

```html
customHeaderKey1: customHeaderVal1
customHeaderKey2: customHeaderVal2
```

2. 将自定义headers作为参数发送

```headers={"token":"1284034","deviceId":"104drldu34","appId":"1453"}```

请求参数名为```headers`` 其中kv 必须为字符串

![Xshot-0023.png](https://i.loli.net/2019/10/22/oMy2H3jg8FKB4bi.png)

#### 说明

响应头中```x-request-time```表示响应时长 毫秒

#### *License*

http-bridge is licensed under the [Apache License]((https://github.com/sixinyiyu/http-bridge/blob/master/LICENSE)), Version 2.0