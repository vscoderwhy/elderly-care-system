# API 文档

## 基础信息

- **Base URL**: `http://localhost:8080/api`
- **认证方式**: Bearer Token (JWT)
- **Content-Type**: `application/json`

## 认证相关

### 注册

```
POST /auth/register
```

**请求体**:
```json
{
  "phone": "13800138000",
  "password": "123456",
  "nickname": "张三"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "phone": "13800138000",
      "nickname": "张三",
      "avatar": ""
    }
  }
}
```

### 登录

```
POST /auth/login
```

**请求体**:
```json
{
  "phone": "13800138000",
  "password": "123456"
}
```

**响应**: 同注册

### 微信登录

```
POST /auth/wechat-login
```

**请求体**:
```json
{
  "code": "wx_code",
  "userInfo": {
    "nickName": "微信用户",
    "avatarUrl": "https://..."
  }
}
```

## 用户相关

### 获取个人信息

```
GET /user/profile
Authorization: Bearer {token}
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "phone": "13800138000",
    "nickname": "张三",
    "avatar": "",
    "roles": [
      {
        "id": 1,
        "name": "family"
      }
    ]
  }
}
```

### 更新个人信息

```
PUT /user/profile
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "nickname": "新昵称",
  "avatar": "https://..."
}
```

### 获取关联老人列表

```
GET /user/elderly-list
Authorization: Bearer {token}
```

## 老人管理

### 老人列表

```
GET /elderly?page=1&page_size=20
Authorization: Bearer {token}
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "name": "张奶奶",
        "gender": "女",
        "phone": "13900139000",
        "care_level": 2,
        "bed": {
          "id": 1,
          "name": "301-1"
        }
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 老人详情

```
GET /elderly/:id
Authorization: Bearer {token}
```

### 创建老人档案

```
POST /elderly
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "name": "李奶奶",
  "gender": "女",
  "id_card": "110101195001011234",
  "phone": "13900139000",
  "emergency_contact": "张三",
  "emergency_phone": "13800138000",
  "care_level": 2
}
```

### 更新老人档案

```
PUT /elderly/:id
Authorization: Bearer {token}
```

### 删除老人档案

```
DELETE /elderly/:id
Authorization: Bearer {token}
```

## 护理相关

### 护理记录列表

```
GET /care/records?elderly_id=1&page=1&page_size=20
Authorization: Bearer {token}
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "elderly_id": 1,
        "care_item": {
          "id": 1,
          "name": "喂饭",
          "category": "feeding"
        },
        "staff": {
          "id": 2,
          "nickname": "王护工"
        },
        "notes": "进食良好",
        "images": ["https://..."],
        "recorded_at": "2024-01-01T12:00:00Z"
      }
    ],
    "total": 50
  }
}
```

### 创建护理记录

```
POST /care/records
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "elderly_id": 1,
  "care_item_id": 1,
  "notes": "进食良好",
  "images": ["https://..."]
}
```

### 我的护理任务

```
GET /care/my-tasks
Authorization: Bearer {token}
```

### 护理项目列表

```
GET /care/items
Authorization: Bearer {token}
```

## 财务相关

### 账单列表

```
GET /bills?elderly_id=1&page=1&page_size=20
Authorization: Bearer {token}
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "bill_no": "202401001",
        "elderly_id": 1,
        "elderly": {
          "name": "张奶奶"
        },
        "total_amount": 3000.00,
        "status": "unpaid",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 12
  }
}
```

### 账单详情

```
GET /bills/:id
Authorization: Bearer {token}
```

### 支付账单

```
POST /bills/:id/pay
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "amount": 3000.00,
  "method": "wechat",
  "transaction_no": "wx123456"
}
```

## 错误码

| Code | Message | 说明 |
|------|---------|------|
| 0 | success | 成功 |
| 400 | Invalid request | 请求参数错误 |
| 401 | Unauthorized | 未授权 |
| 404 | Not found | 资源不存在 |
| 500 | Internal server error | 服务器错误 |

## 状态码

| HTTP Code | 说明 |
|-----------|------|
| 200 | 请求成功（业务错误见 code） |
| 401 | Token 无效或过期 |
| 500 | 服务器错误 |
