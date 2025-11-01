# 微信access_token接口使用说明

## 接口概述

新增的`/api/access_token`接口用于获取微信公众平台的access_token，这是调用微信API的基础凭证。

## 接口地址

```
GET /api/access_token
```

## 请求参数

无需请求参数，接口会自动从环境变量获取配置信息。

## 环境变量配置

在使用接口前，需要设置以下环境变量：

```bash
# 微信公众平台的AppID
WECHAT_APPID=your_appid_here

# 微信公众平台的AppSecret
WECHAT_SECRET=your_secret_here
```

## 响应格式

### 成功响应
```json
{
  "code": 0,
  "data": {
    "access_token": "ACCESS_TOKEN_VALUE",
    "expires_in": 7200,
    "timestamp": 1635768900,
    "success": true,
    "message": "获取access_token成功"
  }
}
```

### 错误响应
```json
{
  "code": -1,
  "errorMsg": "错误信息",
  "data": null
}
```

## 使用示例

### 1. 使用curl测试
```bash
curl http://localhost:80/api/access_token
```

### 2. 在代码中使用
```javascript
// JavaScript示例
fetch('/api/access_token')
  .then(response => response.json())
  .then(data => {
    if (data.code === 0) {
      const accessToken = data.data.access_token;
      console.log('获取到的access_token:', accessToken);
    } else {
      console.error('获取失败:', data.errorMsg);
    }
  });
```

```python
# Python示例
import requests

response = requests.get('http://localhost:80/api/access_token')
data = response.json()

if data['code'] == 0:
    access_token = data['data']['access_token']
    print(f'获取到的access_token: {access_token}')
else:
    print(f'获取失败: {data["errorMsg"]}')
```

## 注意事项

1. **access_token有效期**：微信access_token有效期为2小时（7200秒），需要定时刷新
2. **调用频率限制**：微信API对access_token的获取有频率限制，建议缓存使用
3. **安全性**：AppSecret是敏感信息，请妥善保管，不要泄露
4. **环境变量**：在生产环境中，请通过安全的方式设置环境变量

## 错误处理

接口可能返回以下类型的错误：

- `40001`: 获取access_token时AppSecret错误，或者access_token无效
- `40002`: 不合法的凭证类型
- `40013`: 不合法的AppID
- `40125`: 错误的appsecret
- `-1`: 系统内部错误

## 相关文档

- [微信公众平台API文档](https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html)
- [接口频率限制说明](https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html)