# 开发过程问题记录

> 记录本项目开发中遇到的问题、原因与当前状态，供后续维护参考。

## 1. 环境问题

- **本机没有 .NET SDK**：`dotnet --list-sdks` 为空（只有运行时）。
  - 处理：用官方 `dotnet-install.ps1` 安装 SDK 8.0.423 到 `C:\Users\ztt\.dotnet`。
- **exe 双击/直接启动失败，退出码 -2147450730**：
  - 原因：机器上只有 .NET Desktop Runtime 6.0.11 和 10.0.3，WPF 程序（net8.0-windows）需要 8.0.x，找不到匹配运行时。
  - 处理：静默安装 Desktop Runtime 8.0.29（官方 aka.ms 直链）；并设置**用户环境变量** `DOTNET_ROOT=C:\Users\ztt\.dotnet`，否则双击 exe 仍找不到运行时。
  - 注意：换电脑/换用户后，需要重新安装运行时或设置 DOTNET_ROOT。
- **git 2.20.1 不支持 `git init -b main`**：
  - 处理：用 `git init` 初始化，首次提交后执行 `git branch -m main` 改名。

## 2. 游戏数据结构问题

- **新版 character_table.json 字段位置有变化**：
  - 基础属性不在顶层 `attributes`，而在 `phases[].attributesKeyFrames[]`（每档精英化的关键帧）。
  - 信赖加成在 `favorKeyFrames`（最后一档即满信赖）。
  - 画师/CV 已不在 character_table（基础版未实现该字段）。
  - 基建技能数据在独立的 `building_data.json`（`chars` + `buffs` 两张表）。
- **可玩干员过滤**：用 `favorKeyFrames` 至少 2 档判断，得到 492 条。其中包含约 40 个 `token_` 前缀的召唤物，后续可考虑再过滤。
- **技能描述是富文本**：带 `<@ba.vup>...</>`、`<$ba.stun>...</>`、`{变量:格式}` 等标签，已实现 `Localization.CleanText` 清洗。

## 3. 头像源问题（未解决，已临时关闭头像）

- 原头像源 `Aceship/Arknight-Images` 最后更新于 **2024-05-01**。
- 对比结果：492 条干员中 330 条有头像，缺 162 条（其中约 40 条是 token 召唤物，真正缺头像的新干员约 120 个）。
- **决策**：按用户要求，暂时不显示头像：
  - `OperatorItemVm` 中 `Avatar` 置空，界面显示干员名首字占位。
  - `DataUpdater` 中头像下载开关 `DownloadAvatars = false`。
- **后续换源候选（未验证）**：
  - PRTS 维基头像：`https://prts.wiki/images/...`（文件名形如 `头像_干员名.png`，需按干员名映射）。
  - 明日方舟工具箱：`https://toolbox.arknights.xyz/avatar/{干员ID}.png`（待验证）。
  - 直接解包游戏客户端资源（最准确，工作量最大）。

## 4. 代码问题（已修复）

- `GameData.Deserialize<T>` 泛型写错，多包了一层 `Dictionary`，导致编译错误 → 已修复。
- XAML 中 `Border.Style` 属性元素放错位置（写成 Grid 的子元素）→ MC3015 编译错误 → 已修复。
- 模板生成的 App.xaml 空白字符无法精确匹配补丁 → 直接删除重建。
- 增加全局异常捕获，异常写入 exe 同目录 `crash.log`，便于排查运行期问题。

## 5. 待办与注意事项

- Git 首次推送需要 SSH 凭据：配好后执行 `git push -u origin main`。
- `data/` 目录已加入 `.gitignore`（约 30MB 游戏数据不提交进仓库），首次使用点界面上的"检查并更新"下载。
- 当前数据版本：76.0.0（2026/07/30），更新机制用 `data_version.txt` 比对，有新版才下载。
