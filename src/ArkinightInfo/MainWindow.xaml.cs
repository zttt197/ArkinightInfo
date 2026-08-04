using System.Windows;
using ArkinightInfo.ViewModels;

namespace ArkinightInfo;

public partial class MainWindow : Window
{
    public MainWindow()
    {
        InitializeComponent();
        Loaded += OnLoaded;
    }

    private async void OnLoaded(object sender, RoutedEventArgs e)
    {
        var vm = new MainViewModel();
        DataContext = vm;
        await vm.LoadAsync();
    }
}
