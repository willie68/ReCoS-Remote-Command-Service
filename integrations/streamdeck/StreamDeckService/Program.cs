using CommandLine;
using OpenMacroBoard.SDK;
using ReCoS;
using StreamDeckSharp;
using System;
using System.Drawing;
using System.Text.Json;
using System.Threading;


namespace TestStreamDeck
{
    class Program
    {
        public class Options
        {
            [Option('u', "url", Required = true, HelpText = "Setting the url to the ReCoS service")]
            public string ReCoSURL { get; set; }

            [Option('p', "profile", Required = false, HelpText = "the profile to be shown. Defaults are: for the normal streamdeck: streamdeck, for the streamdeck xl: streamdeck_xl, for the streamdeck mini: streamdeck_mini, all others: default")]
            public string Profile { get; set; }
        }


        static void Main(string[] args)
        {
            Parser.Default.ParseArguments<Options>(args)
                   .WithParsed<Options>(o => Connect(o));
            //This example is designed for the 5x3 (original) Stream Deck.
        }

        private const string DEFAULT_PROFILE_NAME = "default";

        private static GridKeyPositionCollection keyGrid;
        private static IStreamDeckBoard deck;
        private static RecosClient client;
        private static Profile activeProfile;
        private static Button[] buttons;
        private static readonly Mutex Btnmut = new();
        private static Page activePage;
        private static Options flags;

        private static bool IsDeckConnected;
        private static bool IsReCoSConnected;


        static void Connect(Options options)
        {
            IsDeckConnected = false;
            IsReCoSConnected = false;

            flags = options;
            Console.WriteLine($"URL to the ReCoS Service: -u {flags.ReCoSURL}");

            while (!(IsReCoSConnected && IsDeckConnected))
            {
                Connect2StreamDeck();
                Connect2ReCoS();
                if (IsReCoSConnected && IsDeckConnected)
                {
                    break;
                }
                Thread.Yield();
                Thread.Sleep(1000);
            }

            InitApplication();

            Console.ReadKey();
            deck.Dispose();
            client.Dispose();
            /*            else
                        {
                            Console.WriteLine("no streamdecks found");
                            Environment.Exit(1);
                        }
            */
        }

        static bool Connect2StreamDeck()
        {
            if (!IsDeckConnected)
            {
                var Devices = StreamDeck.EnumerateDevices();
                if (Devices.GetEnumerator().MoveNext())
                {
                    foreach (IStreamDeckRefHandle Device in Devices)
                    {
                        Console.WriteLine($"found device: {Device.DeviceName}");
                    }
                    deck = StreamDeck.OpenDevice();

                    if (deck != null)
                    {
                        IsDeckConnected = true;
                        return true;
                    }
                }
                else
                {
                    Console.Error.WriteLine("No streamdeck found.");
                }
            }
            return false;
        }


        static bool Connect2ReCoS()
        {
            if (client == null)
            {
                client = new RecosClient(flags.ReCoSURL);
                client.ImageWidth = 70;
                client.ImageHeight = 70;
            }
            if (!client.IsConnected())
            {
                client.Connect();
                IsReCoSConnected = client.IsConnected();
            }
            return IsReCoSConnected;
        }

        static void InitApplication()
        {
            var defaultProfile = "default";
            keyGrid = (deck.Keys as GridKeyPositionCollection) ?? throw new InvalidOperationException("Deck not supported");
            switch (keyGrid.Count)
            {
                case 6:
                    defaultProfile = "streamdeck_mini";
                    break;
                case 15:
                    defaultProfile = "streamdeck";
                    break;
                case 32:
                    defaultProfile = "streamdeck_xl";
                    break;
            }

            Console.WriteLine($"profile: {flags.Profile}");
            Console.WriteLine($"default: {defaultProfile}");

            activeProfile = ReadProfile(flags.Profile, defaultProfile);

            //            Console.WriteLine(JsonSerializer.Serialize(activeProfile));
            deck.SetBrightness(100);

            deck.KeyStateChanged += Deck_KeyPressed;

            client.SetProfile(activeProfile.Name);
            client.MessageReceived += MessageReceived;

            activePage = activeProfile.Pages[0];
            SwitchPage(activePage.Name);
        }

        private static void SwitchPage(string pagename)
        {
            activePage = GetPage(pagename);
            if (activePage == null)
            {
                activePage = activeProfile.Pages[0];
            }
            deck.ClearKeys();

            var kID = 0;
            if (Btnmut.WaitOne(1000))
            {
                buttons = new Button[activePage.Columns * activePage.Rows];
                foreach (string cellActionName in activePage.Cells)
                {
                    ReCoS.Action action = GetAction(cellActionName);
                    if (action != null)
                    {
                        buttons[kID] = new Button(action);
                        var bmp = GenerateKeyBitmap(action, null, null, null);
                        deck.SetKeyBitmap(kID, bmp);
                    }
                    kID++;
                }
                Btnmut.ReleaseMutex();
            }

        }

        private static void MessageReceived(object sender, MessageReceived e)
        {
            if (activeProfile.Name.Equals(e.Message.Profile))
            {
                // the message is for the actual profile
                if (Array.IndexOf(activePage.Cells, e.Message.Action) >= 0)
                {
                    // this message is for the actual page
                    // getting the button to display
                    int kID = 0;
                    Button button = null;
                    if (Btnmut.WaitOne(1000))
                    {
                        foreach (ReCoS.Button Button in buttons)
                        {
                            if (Button != null)
                            {
                                if (Button.Action.Name.Equals(e.Message.Action))
                                {
                                    button = Button;
                                    break;
                                }
                            }
                            kID++;
                        }
                        Btnmut.ReleaseMutex();
                    }
                    if (button != null)
                    {
                        // generating the bitmap
                        var bmp = GenerateKeyBitmap(button.Action, e.Message.Title, e.Message.Text, e.Message.ImageURL);
                        // sending the bitmap to the streamdeck
                        deck.SetKeyBitmap(kID, bmp);
                    }
                }
                else
                {
                    if (!String.IsNullOrEmpty(e.Message.Page))
                    {
                        var jsonStr = JsonSerializer.Serialize(e.Message);
                        Console.WriteLine($"Message received: \r\n{jsonStr}");
                        // { "profile":"streamdeck","action":"","page":"clocks","imageurl":"check_mark.svg","title":"","text":"","state":0,"command":""}
                        var page = GetPage(e.Message.Page);
                        if (page != null)
                        {
                            SwitchPage(page.Name);
                        }
                    }
                }
            }
        }

        private static Page GetPage(string pageName)
        {
            if (!String.IsNullOrEmpty(pageName))
            {
                foreach (Page page in activeProfile.Pages)
                {
                    if (pageName.Equals(page.Name))
                    {
                        return page;
                    }
                }
            }
            return null;
        }
        static KeyBitmap GenerateKeyBitmap(ReCoS.Action action, string title, string text, string image)
        {
            return KeyBitmap.Create.FromGraphics(72, 72, (g) =>
            {
                Image img = null;
                if (String.IsNullOrEmpty(image))
                {
                    if (!String.IsNullOrEmpty(action.Icon))
                    {
                        img = client.GetImage(action.Icon);
                    }
                }
                else
                {
                    img = client.GetImage(image);
                }
                if (img != null)
                {
                    g.DrawImage(img, new Point(0, 0));
                }


                var b = Brushes.White;
                if (!String.IsNullOrEmpty(action.Fontcolor))
                {
                    Color p = ColorTranslator.FromHtml(action.Fontcolor);
                    b = new SolidBrush(p);
                }
                var fontsize = 10;
                if (action.Fontsize > 0)
                {
                    fontsize = action.Fontsize;
                }
                var myTitle = action.Title;
                if (!String.IsNullOrEmpty(title))
                {
                    myTitle = title;
                }
                var fb = new Font("Arial", fontsize, FontStyle.Bold);
                var size = g.MeasureString(myTitle, fb);
                var xPos = 0;
                if (size.Width < 72)
                {
                    xPos = (72 - Convert.ToInt32(size.Width)) / 2;
                }
                var origin = new PointF(xPos, 0);
                g.DrawString(myTitle, fb, b, origin);

                if (!String.IsNullOrEmpty(text))
                {
                    fb = new Font("Arial", fontsize);
                    size = g.MeasureString(text, fb);
                    xPos = 0;
                    if (size.Width < 72)
                    {
                        xPos = (72 - Convert.ToInt32(size.Width)) / 2;
                    }
                    origin = new PointF(xPos, 36 - fontsize);
                    g.DrawString(text, fb, b, origin);
                }
            });
        }

        static ReCoS.Action GetAction(string name)
        {
            foreach (ReCoS.Action action in activeProfile.Actions)
            {
                if (String.Equals(action.Name, name))
                {
                    return action;
                }
            }
            return null;
        }

        static Profile ReadProfile(string profileName, string defaultProfileName)
        {
            string[] profilenames = client.ProfileNames();
            bool isNamedProfile = false;
            bool isDefaultStreamdeckProfile = false;
            bool isDefaultProfile = false;

            string myProfileName = profileName;
            foreach (var name in profilenames)
            {
                if (String.Equals(name, profileName))
                {
                    isNamedProfile = true;
                    break;
                }
                if (String.Equals(name, defaultProfileName))
                {
                    isDefaultStreamdeckProfile = true;
                }
                if (String.Equals(name, DEFAULT_PROFILE_NAME))
                {
                    isDefaultProfile = true;
                }
            }

            if (!isNamedProfile)
            {
                myProfileName = defaultProfileName;
                if (!isDefaultStreamdeckProfile)
                {
                    myProfileName = DEFAULT_PROFILE_NAME;
                    if (!isDefaultProfile)
                    {
                        myProfileName = profilenames[0];
                    }
                }
            }

            return client.GetProfile(myProfileName);
        }

        static StreamDeckSharp.IStreamDeckBoard ConnectToFirstStreamdeck()
        {
            try
            {
                var deck = StreamDeck.OpenDevice();
                return deck;
            }
            catch (StreamDeckSharp.Exceptions.StreamDeckNotFoundException e)
            {
                Console.WriteLine("no streamdeck found.");
            }
            return null;
        }

        private static void Deck_KeyPressed(object sender, KeyEventArgs e)
        {
            if (!(sender is IMacroBoard d))
            {
                return;
            }

            //            Console.WriteLine($"key {e.Key} pressed. IsDown: {e.IsDown}");
            if (e.IsDown)
            {
                Button button = null;
                if (Btnmut.WaitOne(1000))
                {
                    button = buttons[e.Key];
                    Btnmut.ReleaseMutex();
                }
                if (button != null)
                {
                    if (button.Action != null && (String.Equals(button.Action.Type, "SINGLE") || String.Equals(button.Action.Type, "TOGGLE") || String.Equals(button.Action.Type, "MULTI")))
                    {
                        client.SendClick(activeProfile.Name, activePage.Name, button.Action.Name);
                    }
                }
            }

        }
    }
}
