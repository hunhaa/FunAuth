package main

import (
"bufio"
"fmt"
"math/rand"
"os"
"strconv"
"strings"
"time"

"github.com/Yeah114/g79client"
)

func init() {
rand.Seed(time.Now().UnixNano())
}

func main() {
reader := bufio.NewReader(os.Stdin)
var client *g79client.Client
var cookie string

fmt.Println("========================================")
fmt.Println("   FunAuth - 网易 MC 3.9 租赁服工具")
fmt.Println("========================================")
fmt.Println()

for {
fmt.Println("\n请选择操作:")
fmt.Println("1. 设置 Cookie")
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
cookie, _ = reader.ReadString('\n')
cookie = strings.TrimSpace(cookie)
if cookie != "" {
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
// 生成 fbtoken 文件
fbtoken := cookie
fileName := fmt.Sprintf("fbtoken-%s", generateRandomString(5))
err = os.WriteFile(fileName, []byte(fbtoken), 0600)
if err != nil {
fmt.Printf("警告：保存 token 文件失败：%v\n", err)
} else {
fmt.Printf("✓ Token 已保存到文件：%s\n", fileName)
}
}
}

case "2":
if client == nil {
fmt.Println("请先设置 Cookie (选项 1)")
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

case "3":
if client == nil {
fmt.Println("请先设置 Cookie (选项 1)")
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

case "4":
if client == nil {
fmt.Println("请先设置 Cookie (选项 1)")
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

case "5":
if client == nil {
fmt.Println("请先设置 Cookie (选项 1)")
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
