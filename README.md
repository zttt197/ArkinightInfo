# ArkinightInfo 明日方舟干员查询

一款纯本地运行的 Windows 桌面工具，用于查询《明日方舟》干员信息。

## 功能

- 干员列表：按名称/代号搜索，按职业、星级筛选
- 干员详情：面板属性（各精英化满级）、信赖加成、天赋、技能、基建技能
- 手动更新：点击"检查并更新"，先比对数据版本，有新版才下载
- 完全本地使用：数据下载到本地后离线可用，不联网也能查询

## 运行环境

- Windows x64
- .NET Desktop Runtime 8.0.x（若提示缺少运行时，安装
  https://aka.ms/dotnet/8.0/dotnet-runtime-win-x64.exe）

## 使用

1. 打开程序，点击右上角"检查并更新"，首次使用会自动下载游戏数据（约 30MB）。
2. 数据下载完成后即可搜索、筛选、查看干员详情。
3. 游戏更新后，再次点击"检查并更新"即可同步最新数据。

## 开发

- 技术栈：C# / WPF（net8.0-windows），UI 样式基于 HandyControl
- 数据解析：`src/ArkinightInfo/Data/GameData.cs`
- 更新逻辑：`src/ArkinightInfo/Services/DataUpdater.cs`
- 界面：`src/ArkinightInfo/MainWindow.xaml`
- 开发过程遇到的问题与后续计划见 `DEVELOPMENT_NOTES.md`

## 数据来源

- 干员数据：[Kengxxiao/ArknightsGameData](https://github.com/Kengxxiao/ArknightsGameData)
  （社区从游戏客户端解包整理，随游戏版本自动更新，仅作个人学习使用）
- 头像：暂未启用（原头像源已停更，见 DEVELOPMENT_NOTES.md）

## 许可说明

本项目代码仅用于个人学习。游戏数据版权归鹰角网络所有，请勿将本工具及数据用于商业用途。
