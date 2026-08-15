package main

import (
"bufio"
"crypto/rand"
"encoding/hex"
"fmt"
"net/http"
"net/url"
"os"
"strconv"
"strings"
"syscall"

"golang.org/x/term"

"github.com/Yeah114/g79client"
)

func main() {
reader := bufio.NewReader(os.Stdin)
var client *g79client.Client
var cookie string
var fbtoken string

fmt.Println("========================================")
fmt.Println("   FunAuth - 网易 MC 3.9 租赁服工具")
fmt.Println("========================================")
fmt.Println()

for {
fmt.Println("\n请选择操作:")
fmt.Println("1. 账号密码登录")
fmt.Println("2. 邮箱登录")
fmt.Println("3. Cookie 登录")
fmt.Println("4. 搜索租赁服")
fmt.Println("5. 查看可用租赁服列表")
fmt.Println("6. 查看租赁服详情")
fmt.Println("7. 进入租赁服")
fmt.Println("8. 退出")
fmt.Print("> ")

input, _ := reader.ReadString('\n')
input = strings.TrimSpace(input)

switch input {
case "1":
fmt.Print("请输入网易账号：")
username, _ := reader.ReadString('\n')
username = strings.TrimSpace(username)
fmt.Print("请输入密码：")
passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
fmt.Println()
if err != nil {
fmt.Printf("读取密码失败：%v\n", err)
break
}
password := strings.TrimSpace(string(passwordBytes))

if username == "" || password == "" {
fmt.Println("账号或密码不能为空")
break
}

fmt.Println("正在登录...")
token, err := loginWithPassword(username, password)
if err != nil {
fmt.Printf("登录失败：%v\n", err)
break
}
fbtoken = token
fmt.Printf("✓ 登录成功！Token: %s...\n", token[:20])

err = saveTokenToFile(fbtoken)
if err != nil {
fmt.Printf("保存 Token 失败：%v\n", err)
} else {
fmt.Println("✓ Token 已保存到文件")
}

client, err = g79client.NewClient()
if err != nil {
fmt.Printf("创建客户端失败：%v\n", err)
break
}
cookie = "NTES_SESS=" + token
err = client.G79AuthenticateWithCookie(cookie)
if err != nil {
fmt.Printf("认证失败：%v\n", err)
} else {
fmt.Println("✓ 客户端认证成功!")
}

case "2":
fmt.Print("请输入网易邮箱地址：")
email, _ := reader.ReadString('\n')
email = strings.TrimSpace(email)
fmt.Print("请输入邮箱密码：")
passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
fmt.Println()
if err != nil {
fmt.Printf("读取密码失败：%v\n", err)
break
}
password := strings.TrimSpace(string(passwordBytes))

if email == "" || password == "" {
fmt.Println("邮箱或密码不能为空")
break
}

fmt.Println("正在登录...")
token, err := loginWithEmail(email, password)
if err != nil {
fmt.Printf("登录失败：%v\n", err)
break
}
fbtoken = token
fmt.Printf("✓ 登录成功！Token: %s...\n", token[:20])

err = saveTokenToFile(fbtoken)
if err != nil {
fmt.Printf("保存 Token 失败：%v\n", err)
} else {
fmt.Println("✓ Token 已保存到文件")
}

client, err = g79client.NewClient()
if err != nil {
fmt.Printf("创建客户端失败：%v\n", err)
break
}
cookie = "MAIL_SESS=" + token
err = client.G79AuthenticateWithCookie(cookie)
if err != nil {
fmt.Printf("认证失败：%v\n", err)
} else {
fmt.Println("✓ 客户端认证成功!")
}

case "3":
fmt.Print("请输入您的网易登录 Cookie (NTES_SESS 或 MAIL_SESS): ")
cookieInput, _ := reader.ReadString('\n')
cookieInput = strings.TrimSpace(cookieInput)
if cookieInput == "" {
fmt.Println("Cookie 不能为空")
break
}

if strings.HasPrefix(cookieInput, "NTES_SESS=") {
fbtoken = strings.TrimPrefix(cookieInput, "NTES_SESS=")
} else if strings.HasPrefix(cookieInput, "MAIL_SESS=") {
fbtoken = strings.TrimPrefix(cookieInput, "MAIL_SESS=")
} else {
fbtoken = cookieInput
}

cookie = cookieInput
if !strings.HasPrefix(cookie, "NTES_SESS=") && !strings.HasPrefix(cookie, "MAIL_SESS=") {
cookie = "NTES_SESS=" + cookie
}

var err error
client, err = g79client.NewClient()
if err != nil {
fmt.Printf("创建客户端失败：%v\n", err)
break
}
err = client.G79AuthenticateWithCookie(cookie)
if err != nil {
fmt.Printf("认证失败：%v\n", err)
} else {
fmt.Println("✓ 认证成功!")
fmt.Printf("Token: %s...\n", fbtoken[:20])

err = saveTokenToFile(fbtoken)
if err != nil {
fmt.Printf("保存 Token 失败：%v\n", err)
} else {
fmt.Println("✓ Token 已保存到文件")
}
}

case "4":
if client == nil {
fmt.Println("请先登录 (选项 1/2/3)")
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
result, err := client.SearchRentalServerByName(name)
if err != nil {
fmt.Printf("搜索失败：%v\n", err)
} else {
fmt.Printf("找到 %d 个结果:\n", len(result.Entities))
for i, server := range result.Entities {
fmt.Printf("%d. [ID: %s] %s - 在线：%s/%s\n", 
i+1, server.WorldID, server.ServerName, 
server.PlayerCount.String(), server.Capacity.String())
}
}

case "5":
if client == nil {
fmt.Println("请先登录 (选项 1/2/3)")
break
}
fmt.Println("正在获取可用租赁服列表...")
result, err := client.GetAvailableRentalServers(0, 0, 0)
if err != nil {
fmt.Printf("获取失败：%v\n", err)
} else {
fmt.Printf("找到 %d 个可用租赁服:\n", len(result.Entities))
for i, server := range result.Entities {
statusStr := "未知"
if s := server.Status.Int64(); s != 0 || server.Status.String() != "" {
statusVal, _ := strconv.ParseInt(server.Status.String(), 10, 64)
if statusVal != 0 || server.Status.String() == "0" {
statusStr = getServerStatus(int(statusVal))
}
}
fmt.Printf("%d. [ID: %s] %s - 状态：%s\n", 
i+1, server.WorldID, server.ServerName, statusStr)
}
}

case "6":
if client == nil {
fmt.Println("请先登录 (选项 1/2/3)")
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
result, err := client.GetRentalServerDetails(worldID)
if err != nil {
fmt.Printf("获取详情失败：%v\n", err)
} else {
fmt.Println("\n=== 租赁服详情 ===")
fmt.Printf("名称：%s\n", result.Entity.Name)
fmt.Printf("ID: %s\n", result.Entity.WorldID)
if result.Entity.BriefSummary != "" {
fmt.Printf("描述：%s\n", result.Entity.BriefSummary)
}
playerCount := result.Entity.PlayerCount.Int64()
capacity := result.Entity.Capacity.Int64()
if playerCount != 0 || capacity != 0 {
fmt.Printf("在线人数：%d/%d\n", playerCount, capacity)
}
status := result.Entity.Status.Int64()
if status != 0 || result.Entity.Status.String() == "0" {
fmt.Printf("状态：%s\n", getServerStatus(int(status)))
}
fmt.Printf("版本：%s\n", result.Entity.McVersion)
fmt.Printf("类型：%s\n", result.Entity.ServerType)
}

case "7":
if client == nil {
fmt.Println("请先登录 (选项 1/2/3)")
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
authResult, err := client.GenerateRentalGameAuthV2(worldID, password)
if err != nil {
fmt.Printf("获取认证失败：%v\n", err)
break
}
fmt.Printf("\n✓ 认证成功!\n原始响应：%s\n", string(authResult))
fmt.Println("\n请在游戏中使用以上信息加入服务器")

case "8":
fmt.Println("再见!")
os.Exit(0)

default:
fmt.Println("无效选项，请重新选择")
}
}
}

func loginWithPassword(username, password string) (string, error) {
loginURL := "https://id.feijie.cn/password/login"

data := url.Values{}
data.Set("username", username)
data.Set("password", password)

req, err := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
if err != nil {
return "", err
}

req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
return "", err
}
defer resp.Body.Close()

for _, cookie := range resp.Cookies() {
if cookie.Name == "NTES_SESS" {
return cookie.Value, nil
}
}

return "", fmt.Errorf("未找到 NTES_SESS Cookie")
}

func loginWithEmail(email, password string) (string, error) {
loginURL := "https://mail.163.com/auth/login"

data := url.Values{}
data.Set("userId", email)
data.Set("password", password)

req, err := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
if err != nil {
return "", err
}

req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
return "", err
}
defer resp.Body.Close()

for _, cookie := range resp.Cookies() {
if cookie.Name == "MAIL_SESS" {
return cookie.Value, nil
}
}

return "", fmt.Errorf("未找到 MAIL_SESS Cookie")
}

func saveTokenToFile(token string) error {
randomBytes := make([]byte, 3)
if _, err := rand.Read(randomBytes); err != nil {
return err
}
randomStr := hex.EncodeToString(randomBytes)[:5]

filename := "fbtoken-" + randomStr

file, err := os.Create(filename)
if err != nil {
return err
}
defer file.Close()

_, err = file.WriteString(token)
return err
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
