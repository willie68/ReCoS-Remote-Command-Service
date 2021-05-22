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

            Console.WriteLine(JsonSerializer.Serialize(activeProfile));
            deck.SetBrightness(100);

            activePage = activeProfile.Pages[0];
            var kID = 0;
            buttons = new Button[activePage.Columns * activePage.Rows];
            foreach (string cellActionName in activePage.Cells)
            {
                ReCoS.Action action = GetAction(cellActionName);
                if (action != null)
                {
                    buttons[kID] = new Button(action);
                    var bmp = GenerateKeyBitmap(action);
                    deck.SetKeyBitmap(kID, bmp);
                }
                kID++;
            }
            deck.KeyStateChanged += Deck_KeyPressed;
            Message msg = new();
            msg.Profile = activeProfile.Name;
            msg.Command = "change";
            client.SendMessage(msg).Wait();
        }
        static KeyBitmap GenerateKeyBitmap(ReCoS.Action action)
        {
            return KeyBitmap.Create.FromGraphics(72, 72, (g) =>
            {

                if (!String.IsNullOrEmpty(action.Icon))
                {
                    Image img = client.GetImage(action.Icon);
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
                var fb = new Font("Arial", fontsize, FontStyle.Bold);
                var size = g.MeasureString(action.Title, fb);
                var xPos = 0;
                if (size.Width < 72)
                {
                    xPos = (72 - Convert.ToInt32(size.Width)) / 2;
                }
                var origin = new PointF(xPos, 0);
                g.DrawString(action.Title, fb, b, origin);

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

            Console.WriteLine($"key {e.Key} pressed. IsDown: {e.IsDown}");
            if (e.IsDown)
            {
                var button = buttons[e.Key];
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
