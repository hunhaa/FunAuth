package main

import (
"bufio"
"bytes"
"encoding/json"
"fmt"
"io"
"math/rand"
"net/http"
"os"
"strings"
"time"
)

func init() {
rand.Seed(time.Now().UnixNano())
}

// LoginRequest 登录请求结构
type LoginRequest struct {
Username string `json:"username"`
Password string `json:"password"`
Cookie   string `json:"cookie"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Code int `json:"code"`
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Message string `json:"message"`
}

const ServerURL = "http://127.0.0.1:8080"

func main() {
reader := bufio.NewReader(os.Stdin)
var fbtoken string

fmt.Println("========================================")
fmt.Println("   FunAuth - 网易 MC 3.9 租赁服工具")
fmt.Println("========================================")
fmt.Println()

for {
fmt.Println("\n请选择操作:")
fmt.Println("1. 登录获取 fbtoken (使用 Cookie)")
fmt.Println("2. 搜索租赁服")
fmt.Println("3. 查看可用租赁服列表")
fmt.Println("4. 查看租赁服详情")
fmt.Println("5. 进入租赁服")
fmt.Println("6. 退出")
fmt.Print("> ")

input, _ := reader.ReadString('\n')
input = strings.TrimSpace(input)

switch input {
case "1":
fmt.Print("请输入您的网易登录 Cookie: ")
cookie, _ := reader.ReadString('\n')
cookie = strings.TrimSpace(cookie)
if cookie != "" {
// 向验证服务器发送登录请求
loginReq := LoginRequest{Cookie: cookie}
jsonData, err := json.Marshal(loginReq)
if err != nil {
fmt.Printf("JSON 编码失败：%v\n", err)
break
}

resp, err := http.Post(ServerURL+"/api/phoenix/login", "application/json", bytes.NewBuffer(jsonData))
if err != nil {
fmt.Printf("连接服务器失败：%v\n请确保 funauth-server 正在运行\n", err)
break
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
fmt.Printf("读取响应失败：%v\n", err)
break
}

var loginResp LoginResponse
if err := json.Unmarshal(body, &loginResp); err != nil {
fmt.Printf("解析响应失败：%v\n", err)
break
}

if loginResp.Code != 0 {
fmt.Printf("登录失败：%s\n", loginResp.Message)
break
}

fbtoken = loginResp.Data.Token
fmt.Println("✓ 认证成功!")
// 生成 fbtoken 文件
fileName := fmt.Sprintf("fbtoken-%s", generateRandomString(5))
err = os.WriteFile(fileName, []byte(fbtoken), 0600)
if err != nil {
fmt.Printf("警告：保存 token 文件失败：%v\n", err)
} else {
fmt.Printf("✓ Token 已保存到文件：%s\n", fileName)
}
}

case "2":
if fbtoken == "" {
fmt.Println("请先登录获取 fbtoken (选项 1)")
break
}
fmt.Print("请输入要搜索的租赁服名称：")
name, _ := reader.ReadString('\n')
name = strings.TrimSpace(name)
if name == "" {
fmt.Println("名称不能为空")
break
}
fmt.Println("正在搜索...")
// 调用服务端 API 搜索
req, _ := http.NewRequest("GET", ServerURL+"/api/open/g79/rental_search?keyword="+name, nil)
req.Header.Set("Authorization", "Bearer "+fbtoken)
resp, err := http.DefaultClient.Do(req)
if err != nil {
fmt.Printf("请求失败：%v\n", err)
break
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
fmt.Printf("搜索结果：%s\n", string(body))

case "3":
if fbtoken == "" {
fmt.Println("请先登录获取 fbtoken (选项 1)")
break
}
fmt.Println("正在获取可用租赁服列表...")
req, _ := http.NewRequest("GET", ServerURL+"/api/open/g79/rental_available", nil)
req.Header.Set("Authorization", "Bearer "+fbtoken)
resp, err := http.DefaultClient.Do(req)
if err != nil {
fmt.Printf("请求失败：%v\n", err)
break
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
fmt.Printf("可用租赁服列表：%s\n", string(body))

case "4":
if fbtoken == "" {
fmt.Println("请先登录获取 fbtoken (选项 1)")
break
}
fmt.Print("请输入租赁服 ID: ")
worldID, _ := reader.ReadString('\n')
worldID = strings.TrimSpace(worldID)
if worldID == "" {
fmt.Println("ID 不能为空")
break
}
fmt.Println("正在获取详情...")
req, _ := http.NewRequest("GET", ServerURL+"/api/open/g79/rental_details?world_id="+worldID, nil)
req.Header.Set("Authorization", "Bearer "+fbtoken)
resp, err := http.DefaultClient.Do(req)
if err != nil {
fmt.Printf("请求失败：%v\n", err)
break
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
fmt.Printf("租赁服详情：%s\n", string(body))

case "5":
if fbtoken == "" {
fmt.Println("请先登录获取 fbtoken (选项 1)")
break
}
fmt.Print("请输入租赁服 ID: ")
worldID, _ := reader.ReadString('\n')
worldID = strings.TrimSpace(worldID)
if worldID == "" {
fmt.Println("ID 不能为空")
break
}
fmt.Print("请输入服务器密码 (如无密码直接回车): ")
password, _ := reader.ReadString('\n')
password = strings.TrimSpace(password)

fmt.Println("正在获取进入凭证...")
enterReq := map[string]string{"world_id": worldID, "password": password}
jsonData, _ := json.Marshal(enterReq)
req, _ := http.NewRequest("POST", ServerURL+"/api/open/g79/rental_enter_world", bytes.NewBuffer(jsonData))
req.Header.Set("Authorization", "Bearer "+fbtoken)
req.Header.Set("Content-Type", "application/json")
resp, err := http.DefaultClient.Do(req)
if err != nil {
fmt.Printf("请求失败：%v\n", err)
break
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
fmt.Printf("\n✓ 认证成功!\n原始响应：%s\n", string(body))
fmt.Println("\n请在游戏中使用以上信息加入服务器")

case "6":
fmt.Println("再见!")
os.Exit(0)

default:
fmt.Println("无效选项，请重新选择")
}
}
}

func getServerStatus(status int) string {
switch status {
case 0:
return "离线"
case 1:
return "运行中"
case 2:
return "维护中"
default:
return fmt.Sprintf("未知 (%d)", status)
}
}

func generateRandomString(n int) string {
chars := "abcdefghijklmnopqrstuvwxyz0123456789"
result := make([]byte, n)
for i := 0; i < n; i++ {
result[i] = chars[rand.Intn(len(chars))]
}
return string(result)
}
