using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Net.Http;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;

namespace ArkinightInfo.Services;

/// <summary>手动更新：先比对 data_version.txt，有新版才下载数据与缺失头像。</summary>
public static class DataUpdater
{
    private const string DataRepo = "Kengxxiao/ArknightsGameData";
    private const string DataBranch = "master";
    private const string DataFolder = "zh_CN/gamedata/excel";
    private const string ImagesRepo = "Aceship/Arknight-Images";
    private const string ImagesBranch = "main";

    private static readonly string[] DataFiles =
    {
        "character_table.json",
        "skill_table.json",
        "building_data.json",
        "char_meta_table.json",
    };

    // 头像源（Aceship/Arknight-Images）已停更，新干员头像缺失，暂时关闭头像下载。
    // 等换用可用的头像源（如 PRTS 维基）后再打开。见 DEVELOPMENT_NOTES.md。
    private static readonly bool DownloadAvatars = false;

    private static readonly HttpClient Http = CreateHttp();

    public sealed record CheckResult(bool NeedsUpdate, string Message, string? LocalVersion, string? RemoteVersion);

    private static HttpClient CreateHttp()
    {
        var client = new HttpClient { Timeout = TimeSpan.FromMinutes(5) };
        client.DefaultRequestHeaders.UserAgent.ParseAdd("ArkinightInfo/1.0");
        return client;
    }

    public static async Task<CheckResult> CheckAsync(string dataRoot, CancellationToken ct = default)
    {
        var local = ReadVersionFile(dataRoot);
        var remote = await HttpGetStringAsync(DataUrl("data_version.txt"), ct).ConfigureAwait(false);
        remote = remote.Trim();

        if (local is null)
            return new CheckResult(true, "本地还没有数据，点击更新开始首次下载。", null, remote);

        return local == remote
            ? new CheckResult(false, "已是最新版本，无需更新。", local, remote)
            : new CheckResult(true, "发现新版本，开始下载数据。", local, remote);
    }

    public static async Task UpdateAsync(string dataRoot, IProgress<string>? progress = null, CancellationToken ct = default)
    {
        var gamedataDir = Path.Combine(dataRoot, "gamedata");
        var avatarDir = Path.Combine(dataRoot, "avatars");
        Directory.CreateDirectory(gamedataDir);
        Directory.CreateDirectory(avatarDir);

        // 1. 下载 JSON 数据（先写临时文件，成功后原子替换）
        for (var i = 0; i < DataFiles.Length; i++)
        {
            ct.ThrowIfCancellationRequested();
            var name = DataFiles[i];
            progress?.Report($"正在下载 {name}（{i + 1}/{DataFiles.Length}）…");
            var tmp = Path.Combine(gamedataDir, name + ".tmp");
            var target = Path.Combine(gamedataDir, name);
            await HttpGetToFileAsync(DataUrl(name), tmp, ct).ConfigureAwait(false);
            File.Move(tmp, target, overwrite: true);
        }

        // 2. 补下缺失的头像（当前已关闭，见 DownloadAvatars 注释）
        if (DownloadAvatars)
        {
            var ids = ReadOperatorIds(Path.Combine(gamedataDir, "character_table.json"));
            var missing = ids.Where(id => !File.Exists(Path.Combine(avatarDir, id + ".png"))).ToList();
            progress?.Report($"正在检查头像（共 {ids.Count} 名干员，缺 {missing.Count} 个头像）…");

            var semaphore = new SemaphoreSlim(8);
            var done = 0;
            await Parallel.ForEachAsync(missing, new ParallelOptions { MaxDegreeOfParallelism = 8, CancellationToken = ct },
                async (id, token) =>
                {
                    await semaphore.WaitAsync(token).ConfigureAwait(false);
                    try
                    {
                        await HttpGetToFileAsync(AvatarUrl(id), Path.Combine(avatarDir, id + ".png"), token).ConfigureAwait(false);
                    }
                    catch (HttpRequestException)
                    {
                        // 个别干员没有头像文件，忽略
                    }
                    finally
                    {
                        semaphore.Release();
                        var current = Interlocked.Increment(ref done);
                        if (current % 20 == 0 || current == missing.Count)
                            progress?.Report($"头像下载进度：{current}/{missing.Count}");
                    }
                });
        }

        // 3. 最后更新版本文件（作为“更新成功”的标记）
        var versionTmp = Path.Combine(gamedataDir, "data_version.txt.tmp");
        var versionTarget = Path.Combine(gamedataDir, "data_version.txt");
        await HttpGetToFileAsync(DataUrl("data_version.txt"), versionTmp, ct).ConfigureAwait(false);
        File.Move(versionTmp, versionTarget, overwrite: true);

        progress?.Report($"更新完成：{DataFiles.Length} 个数据文件。");
    }

    private static string? ReadVersionFile(string dataRoot)
    {
        var path = Path.Combine(dataRoot, "gamedata", "data_version.txt");
        return File.Exists(path) ? File.ReadAllText(path).Trim() : null;
    }

    private static List<string> ReadOperatorIds(string characterTablePath)
    {
        var result = new List<string>();
        if (!File.Exists(characterTablePath))
            return result;
        using var doc = JsonDocument.Parse(File.ReadAllText(characterTablePath));
        foreach (var prop in doc.RootElement.EnumerateObject())
            result.Add(prop.Name);
        return result;
    }

    private static string DataUrl(string file)
        => $"https://raw.githubusercontent.com/{DataRepo}/{DataBranch}/{DataFolder}/{file}";

    private static string AvatarUrl(string id)
        => $"https://raw.githubusercontent.com/{ImagesRepo}/{ImagesBranch}/avatars/{id}.png";

    private static async Task<string> HttpGetStringAsync(string url, CancellationToken ct)
    {
        using var resp = await Http.GetAsync(url, ct).ConfigureAwait(false);
        resp.EnsureSuccessStatusCode();
        return await resp.Content.ReadAsStringAsync(ct).ConfigureAwait(false);
    }

    private static async Task HttpGetToFileAsync(string url, string path, CancellationToken ct)
    {
        using var resp = await Http.GetAsync(url, HttpCompletionOption.ResponseHeadersRead, ct).ConfigureAwait(false);
        resp.EnsureSuccessStatusCode();
        await using var fs = new FileStream(path, FileMode.Create, FileAccess.Write, FileShare.None);
        await resp.Content.CopyToAsync(fs, ct).ConfigureAwait(false);
    }
}
