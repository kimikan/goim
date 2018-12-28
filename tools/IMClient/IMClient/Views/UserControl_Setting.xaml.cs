using IMClient.Helpers;
using System;
using System.IO;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;

namespace IMClient.Views
{
    /// <summary>
    /// UserControl_Setting.xaml 的交互逻辑
    /// </summary>
    public partial class UserControl_Setting : UserControl
    {
        public UserControl_Setting()
        {
            InitializeComponent();
            DataContext = new ViewModel();
        }
        
    }

    class ViewModel : PropertyChangedBase
    {
        public ViewModel()
        {
            Load();
        }

        #region methods

        private void Load()
        {
            try
            {
                PublicKey = File.ReadAllText("PUB");
                PrivateKey = File.ReadAllText("PRI");
                Helpers.Configuration.Instance.UserInfo.PublicKey = PublicKey;
                Helpers.Configuration.Instance.UserInfo.PrivateKey = PrivateKey;
            }
            catch (Exception) { }
        }

        private void Save()
        {
            try
            {
                File.WriteAllText("PUB", PublicKey);
                File.WriteAllText("PRI", PrivateKey);
                Helpers.Configuration.Instance.UserInfo.PublicKey = PublicKey;
                Helpers.Configuration.Instance.UserInfo.PrivateKey = PrivateKey;
                MessageBox.Show("Success!");
            }
            catch (Exception) { }
        }

        #endregion

        #region commands

        public ICommand OkCommand
        {
            get
            {
                return new GenericCommand()
                {
                    CanExecuteCallback = arg =>
                    {
                        return true;
                    },
                    ExecuteCallback = arg =>
                    {
                        Save();
                    }
                };
            }
        }

        public ICommand ResetCommand
        {
            get
            {
                return new GenericCommand()
                {
                    CanExecuteCallback = arg =>
                    {
                        return true;
                    },
                    ExecuteCallback = arg =>
                    {
                        Load();
                    }
                };
            }
        }

        public ICommand CancelCommand
        {
            get
            {
                return new GenericCommand()
                {
                    CanExecuteCallback = arg =>
                    {
                        return true;
                    },
                    ExecuteCallback = arg =>
                    {
                        Load();
                    }
                };
            }
        }
        #endregion

        #region properties

        public string PublicKey
        {
            get
            {
                return _pub;
            }
            set
            {
                if (_pub != value)
                {
                    _pub = value;
                    NotifyPropertyChanged(() => PublicKey);
                }
            }
        }
        private string _pub = string.Empty;

        public string PrivateKey
        {
            get
            {
                return _pri;
            }
            set
            {
                if (_pri != value)
                {
                    _pri = value;
                    NotifyPropertyChanged(() => PrivateKey);
                }
            }
        }
        private string _pri = string.Empty;

        #endregion
    }
}
