using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.ComponentModel;
using System.IO;
using System.Linq;
using System.Runtime.CompilerServices;
using System.Threading.Tasks;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using ArkinightInfo.Data;
using ArkinightInfo.Models;
using ArkinightInfo.Services;

namespace ArkinightInfo.ViewModels;

public abstract class ObservableObject : INotifyPropertyChanged
{
    public event PropertyChangedEventHandler? PropertyChanged;

    protected void OnPropertyChanged([CallerMemberName] string? name = null)
        => PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(name));
}

public sealed class RelayCommand : ICommand
{
    private readonly Func<Task> _execute;
    private readonly Func<bool>? _canExecute;

    public RelayCommand(Func<Task> execute, Func<bool>? canExecute = null)
    {
        _execute = execute;
        _canExecute = canExecute;
    }

    public event EventHandler? CanExecuteChanged;

    public bool CanExecute(object? parameter) => _canExecute?.Invoke() ?? true;

    public async void Execute(object? parameter) => await _execute();

    public void RaiseCanExecuteChanged() => CanExecuteChanged?.Invoke(this, EventArgs.Empty);
}

public sealed class MainViewModel : ObservableObject
{
    private readonly string _dataRoot = GameData.ResolveDataRoot();
    private List<OperatorItemVm> _all = [];

    public ObservableCollection<OperatorItemVm> Operators { get; } = [];

    public IReadOnlyList<string> ProfessionOptions { get; } =
        new[] { "全部职业" }.Concat(Localization.Professions.Values).ToArray();

    public IReadOnlyList<string> RarityOptions { get; } =
        new[] { "全部星级", "1星", "2星", "3星", "4星", "5星", "6星" };

    private string _searchText = "";
    public string SearchText
    {
        get => _searchText;
        set { if (_searchText != value) { _searchText = value; OnPropertyChanged(); ApplyFilter(); } }
    }

    private string _selectedProfession = "全部职业";
    public string SelectedProfession
    {
        get => _selectedProfession;
        set { if (_selectedProfession != value) { _selectedProfession = value; OnPropertyChanged(); ApplyFilter(); } }
    }

    private string _selectedRarity = "全部星级";
    public string SelectedRarity
    {
        get => _selectedRarity;
        set { if (_selectedRarity != value) { _selectedRarity = value; OnPropertyChanged(); ApplyFilter(); } }
    }

    private OperatorItemVm? _selectedOperator;
    public OperatorItemVm? SelectedOperator
    {
        get => _selectedOperator;
        set
        {
            if (_selectedOperator != value)
            {
                _selectedOperator = value;
                OnPropertyChanged();
                SelectedDetail = value?.ToDetail();
            }
        }
    }

    private OperatorDetailVm? _selectedDetail;
    public OperatorDetailVm? SelectedDetail
    {
        get => _selectedDetail;
        private set { _selectedDetail = value; OnPropertyChanged(); }
    }

    private string _statusText = "正在加载…";
    public string StatusText
    {
        get => _statusText;
        private set { _statusText = value; OnPropertyChanged(); }
    }

    private string _dataVersion = "未知";
    public string DataVersion
    {
        get => _dataVersion;
        private set { _dataVersion = value; OnPropertyChanged(); }
    }

    private bool _isBusy;
    public bool IsBusy
    {
        get => _isBusy;
        private set
        {
            if (_isBusy != value)
            {
                _isBusy = value;
                OnPropertyChanged();
                UpdateCommand.RaiseCanExecuteChanged();
            }
        }
    }

    public RelayCommand UpdateCommand { get; }

    public MainViewModel()
    {
        UpdateCommand = new RelayCommand(UpdateAsync, () => !IsBusy);
    }

    public async Task LoadAsync()
    {
        IsBusy = true;
        try
        {
            StatusText = "正在加载本地数据…";
            var game = await Task.Run(() => GameData.Load(_dataRoot)).ConfigureAwait(true);
            DataVersion = game.DataVersion;

            var items = await Task.Run(() => game.Operators.Select(o => new OperatorItemVm(o)).ToList()).ConfigureAwait(true);
            _all = items;
            ApplyFilter();
            SelectedOperator = null;

            StatusText = $"共 {_all.Count} 名干员 · 数据版本 {DataVersion}";
        }
        catch (Exception ex)
        {
            StatusText = "数据加载失败：" + ex.Message;
        }
        finally
        {
            IsBusy = false;
        }
    }

    private async Task UpdateAsync()
    {
        if (IsBusy)
            return;

        IsBusy = true;
        try
        {
            StatusText = "正在检查更新…";
            var check = await Task.Run(() => DataUpdater.CheckAsync(_dataRoot)).ConfigureAwait(true);
            StatusText = check.Message;

            if (!check.NeedsUpdate)
                return;

            var progress = new Progress<string>(s => StatusText = s);
            await Task.Run(() => DataUpdater.UpdateAsync(_dataRoot, progress)).ConfigureAwait(true);
            StatusText = "数据已更新，正在重新加载…";
            await LoadAsync().ConfigureAwait(true);
        }
        catch (Exception ex)
        {
            StatusText = "更新失败：" + ex.Message;
        }
        finally
        {
            IsBusy = false;
        }
    }

    private void ApplyFilter()
    {
        var keyword = SearchText.Trim();
        var prof = SelectedProfession;
        var rarity = SelectedRarity;

        var view = _all.Where(o =>
            (string.IsNullOrEmpty(keyword) ||
             o.Name.Contains(keyword, StringComparison.OrdinalIgnoreCase) ||
             o.Appellation.Contains(keyword, StringComparison.OrdinalIgnoreCase)) &&
            (prof == "全部职业" || o.ClassLabel == prof) &&
            (rarity == "全部星级" || o.RarityText == rarity)).ToList();

        Operators.Clear();
        foreach (var item in view)
            Operators.Add(item);

        if (_selectedOperator is not null && !view.Contains(_selectedOperator))
            SelectedOperator = null;
    }
}

public sealed class OperatorItemVm
{
    private readonly Operator _op;

    public OperatorItemVm(Operator op)
    {
        _op = op;
        Name = op.Name;
        Appellation = op.Appellation;
        RarityText = $"{op.Rarity + 1}星";
        Stars = new string('★', op.Rarity + 1);
        ClassLabel = Localization.Map(Localization.Professions, op.Profession);
        PositionLabel = Localization.Map(Localization.Positions, op.Position);
        NationLabel = Localization.Map(Localization.Nations, op.NationId);
        TagsText = op.Tags.Count > 0 ? string.Join(" / ", op.Tags) : "—";
        Initial = op.Name.Length > 0 ? op.Name[..1] : "?";
        // 头像源（Aceship/Arknight-Images）停更于 2024-05，新干员头像缺失，
        // 按用户要求暂时不显示头像，改用首字占位。换源后恢复。见 DEVELOPMENT_NOTES.md。
        Avatar = null;
    }

    public string Name { get; }
    public string Appellation { get; }
    public string RarityText { get; }
    public string Stars { get; }
    public string ClassLabel { get; }
    public string PositionLabel { get; }
    public string NationLabel { get; }
    public string TagsText { get; }
    public string Initial { get; }
    public ImageSource? Avatar { get; }

    public OperatorDetailVm ToDetail() => new(_op, Avatar);

    internal static ImageSource? LoadImage(string path)
    {
        try
        {
            var frame = BitmapFrame.Create(new Uri(path),
                BitmapCreateOptions.None, BitmapCacheOption.OnLoad);
            frame.Freeze();
            return frame;
        }
        catch
        {
            return null;
        }
    }
}

public sealed class OperatorDetailVm
{
    public OperatorDetailVm(Operator op, ImageSource? avatar)
    {
        Name = op.Name;
        Appellation = op.Appellation;
        Stars = new string('★', op.Rarity + 1);
        RarityText = $"{op.Rarity + 1}星";
        ClassLabel = Localization.Map(Localization.Professions, op.Profession);
        PositionLabel = Localization.Map(Localization.Positions, op.Position);
        NationLabel = Localization.Map(Localization.Nations, op.NationId);
        TagsText = op.Tags.Count > 0 ? string.Join(" / ", op.Tags) : "—";
        Avatar = avatar;
        Initial = op.Name.Length > 0 ? op.Name[..1] : "?";

        var last = op.Phases.Count > 0 ? op.Phases[^1] : null;
        DeployCostText = last is null ? "—" : FormatNum(last.MaxAttributes.Cost);
        RedeployText = last is null ? "—" : FormatNum(last.MaxAttributes.RespawnTime) + "s";

        PhaseRows = op.Phases.Select(p => new PhaseRowVm(p)).ToList();

        TrustText = op.TrustBonus is { } trust && (trust.Hp > 0 || trust.Atk > 0)
            ? $"生命 +{trust.Hp} / 攻击 +{trust.Atk}"
            : "—";

        Talents = op.Talents
            .Select((t, i) => new TalentVm { Title = $"天赋{i + 1} · {t.Name}", Meta = t.UnlockText, Desc = t.Description })
            .ToList();

        Skills = op.Skills
            .Select((s, i) =>
            {
                var meta = string.Join(" · ", new[]
                {
                    s.UnlockText,
                    s.SpTypeLabel,
                    s.SkillTypeLabel,
                    s.SpCost > 0 ? $"技力消耗 {FormatNum(s.SpCost)}" : "",
                    s.InitSp > 0 ? $"初始技力 {FormatNum(s.InitSp)}" : "",
                    s.Duration > 0 ? $"持续 {FormatNum(s.Duration)}s" : "",
                }.Where(x => x.Length > 0));
                return new SkillVm { Title = $"技能{i + 1} · {s.Name}", Meta = meta, Desc = s.Description };
            })
            .ToList();

        BaseSkills = op.BaseSkills
            .Select((b, i) => new BaseSkillVm { Title = $"基建技能{i + 1} · {b.Name}", Meta = b.RoomLabel, Desc = b.Description })
            .ToList();
    }

    public string Name { get; }
    public string Appellation { get; }
    public string Stars { get; }
    public string RarityText { get; }
    public string ClassLabel { get; }
    public string PositionLabel { get; }
    public string NationLabel { get; }
    public string TagsText { get; }
    public string DeployCostText { get; }
    public string RedeployText { get; }
    public string TrustText { get; }
    public string Initial { get; }
    public ImageSource? Avatar { get; }
    public IReadOnlyList<PhaseRowVm> PhaseRows { get; }
    public IReadOnlyList<TalentVm> Talents { get; }
    public IReadOnlyList<SkillVm> Skills { get; }
    public IReadOnlyList<BaseSkillVm> BaseSkills { get; }

    private static string FormatNum(double v)
        => v == Math.Floor(v) ? ((long)v).ToString() : v.ToString("0.##");
}

public sealed class PhaseRowVm
{
    public PhaseRowVm(PhaseInfo p)
    {
        Label = $"精英{p.Index} 满级 Lv.{p.MaxLevel}";
        Hp = Format(p.MaxAttributes.MaxHp);
        Atk = Format(p.MaxAttributes.Atk);
        Def = Format(p.MaxAttributes.Def);
        Res = Format(p.MaxAttributes.MagicResistance);
        Cost = Format(p.MaxAttributes.Cost);
        Redeploy = Format(p.MaxAttributes.RespawnTime) + "s";
        Block = Format(p.MaxAttributes.BlockCnt);
    }

    public string Label { get; }
    public string Hp { get; }
    public string Atk { get; }
    public string Def { get; }
    public string Res { get; }
    public string Cost { get; }
    public string Redeploy { get; }
    public string Block { get; }

    private static string Format(double v)
        => v == Math.Floor(v) ? ((long)v).ToString() : v.ToString("0.##");
}

public sealed class TalentVm
{
    public string Title { get; init; } = "";
    public string Meta { get; init; } = "";
    public string Desc { get; init; } = "";
}

public sealed class SkillVm
{
    public string Title { get; init; } = "";
    public string Meta { get; init; } = "";
    public string Desc { get; init; } = "";
}

public sealed class BaseSkillVm
{
    public string Title { get; init; } = "";
    public string Meta { get; init; } = "";
    public string Desc { get; init; } = "";
}
