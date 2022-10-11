## 插件开发模版(普通)

```javascript
// [name:Cookie检测（demo 插件开发演示）]

// 第一行为插件名称， 在后台显示的使用会用到

// 返回数据格式
// return {
//      // 代表是否允许通过
//     "bool": true,
//      // 处理后的变量
//     "env": env
// }

// 必须以main为函数名, env为传入变量
function main(env) {
    let result = request({
        "method": "get",
        "url": "https://plogin.m.jd.com/cgi-bin/ml/islogin",
        "headers": {
            "Cookie": env,
            "User-Agent": "jdapp;iPhone;9.4.4;14.3;network/4g;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
            "Referer": "https://h5.m.jd.com/"
        },
        "dataType": "json",
        "timeout": 5 * 1000
    })

    if (result) {
        // 判断是否过期
        if (result["islogin"] === "1"){
            // Cookie有效
            return {
                "bool": true,
                "env": env
            }
        } else {
            // Cookie无效
            return {
                "bool": false,
                "env": "Cookie已失效"
            }
        }
    } else {
        return {
            "bool": false,
            "env": "请求失败"
        }
    }
}
```

## 插件开发模版(定时)

```javascript
// [name:Cookie检测（demo 定时插件开发演示）]
// [cron:0 0/4 * * *]
// [env:JD_COOKIE]

// 目前 CRON 定时仅支持5位，请勿填写6位cron定时
// 在线CRON调试网站：http://cron.ciding.cc/
// 0 0/4 * * * 四小时运行一次
// 第一行为插件名称， 在后台显示的使用会用到

// env中的值是青龙面板里面开发者需要获取的变量名，插件运行前。系统会自动获取所有面板中符合此名的变量传入env中

// env 传入格式（举例）：[{ID:1, OldID: "2.10.13的唯一值", Name: "JD_COOKIE", Value: "pt_key=xxxx;pt_pin=xxxx;", Remarks: "nuanxinqing"}]

// 必须以main为函数名, env为传入变量
function main(env) {
    for (let i = 0; i < env.length; i++) {
        if (check(env[i]["Value"])) {
            // Cookie有效
            continue
        } else {
            // Cookie失效，通知提交者
            let b = mass_send_message([env[i]["Remarks"]], "您的Cookie：" + env[i]["Value"] + "已过期，请及时更新")
            if (!b) {
                console.error("Cookie检测插件：" + env[i]["Remarks"] + "用户消息推送失败")
            }
        }
    }
}

function check(env) {
    let result = request({
        "method": "get",
        "url": "https://plogin.m.jd.com/cgi-bin/ml/islogin",
        "headers": {
            "Cookie": env,
            "User-Agent": "jdapp;iPhone;9.4.4;14.3;network/4g;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
            "Referer": "https://h5.m.jd.com/"
        },
        "dataType": "json",
        "timeout": 5 * 1000
    })

    if (result) {
        // 判断是否过期
        if (result["islogin"] === "1"){
            // Cookie有效
            return true
        } else {
            // Cookie无效
            return false
        }
    } else {
        return false
    }
}
```



## 封装可用方法

### *request*

```javascript
let result = request({
    // 请求方式（默认get）
    "method": "get",
    // 请求地址
    "url": "https://me-api.jd.com/user_new/info/GetJDUserInfoUnion",
    // 数据类型（返回数据如果是JSON那么就需要指定为json，否则默认为location）
    "dataType": "",
    // 请求头
    "headers": {
        "Cookie": env,
        "User-Agent": "jdapp;iPhone;9.4.4;14.3;network/4g;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
        "Referer": "https://h5.m.jd.com/"
    },
    // 请求体(body为json请求体， 二选一即可)
    "body": {},
    "formdata": {},
    // 自定义超时(单位：纳秒)(5 * 1000 = 5秒)
    "timeout": 5 * 1000
})
```

### *console*

```javascript
// 参数说明：填写打印信息即可（string）
console.info()
console.debug()
console.warn()
console.error()
console.log()
```

### *refind*

```javascript
// 参数说明：ReFind(正则表达式, 待匹配数据) 
// 返回：匹配结果列表（数组）（string）
let result = refind("pt_pin=(.*?);", "pt_pin=xxx")
```

### *replace*

```javascript
// 参数说明：替换文本中的关键词
// s（string）：原始字符串
// old（string）：需要替换的内容
// ne（string）：替换后的内容
// count（int）：需要替换的数量，不填写默认为替换第一个
let result = replace(s, old, ne, count)
```

### *mass_send_message*

```javascript
// 参数说明：WxPusher消息推送
// UserName：字符串数组，可以同时对一或多个用户推送消息【传入参数：用户名】
// Msg：需要推送信息内容
// UserName（补充）：用户名会在插件入口参数env中获取，类型为JS对象。【  env[i]["Remarks"]  】
// 返回值：布尔类型，true：发送成功。false：发送失败
let result = mass_send_message(UserName, Msg)
```

### *sleep*

```javascript
// 参数说明：等待运行，延迟N秒后再继续执行
// s：单位：秒
sleep(s)
```

