using System;
using System.Drawing;
using System.Net.Http;
using System.Net.WebSockets;
using System.Text;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;

namespace ReCoS
{

    public class MessageReceived : EventArgs
    {
        public Message Message { get; }
        public MessageReceived(Message message)
        {
            this.Message = message;
        }
    }
    public class RecosClient
    {
        public event EventHandler<MessageReceived> MessageReceived;

        private const string DefaultURL = "http://127.0.0.1:9280";

        readonly HttpClient client = new();
        readonly ClientWebSocket ws = new();

        private string url;
        private string baseUrl;
        private bool isConnected;
        private Profile[] Profiles;
        private CancellationTokenSource Token;


        public RecosClient()
        {
            this.url = DefaultURL;
            Init();
        }
        public RecosClient(string url)
        {
            if (String.IsNullOrEmpty(url))
            {
                this.url = DefaultURL;
            }
            else
            {
                this.url = url;
            }
            Init();
        }

        private void Init()
        {
            this.baseUrl = $"{url}/api/v1/";
            isConnected = false;
        }

        public string GetUrl()
        {
            return url;
        }
        public string GetBaseUrl()
        {
            return baseUrl;
        }

        public bool IsConnected()
        {
            return isConnected;
        }

        public void Connect()
        {
            try
            {
                GetProfileInfo().Wait();
                Uri uri = new(url);
                Uri wsUri = new($"ws://{uri.Host}:{uri.Port}/ws");
                Task wsConnect = ws.ConnectAsync(wsUri, CancellationToken.None);
                wsConnect.Wait();

                Token = new CancellationTokenSource();
                var Ct = Token.Token;

                Task receiverTask = Receive(Ct);
                //receiverTask.Run(); // Pass same token to Task.Run.

                isConnected = true;
            }
            catch (Exception e)
            {
                Console.Error.WriteLine($"can't connect to ReCoS {e.Message}");
            }
        }

        public void Dispose()
        {
            Token.Cancel();
        }

        private async Task GetProfileInfo()
        {
            var streamTask = client.GetStreamAsync($"{baseUrl}show");
            var profiles = await JsonSerializer.DeserializeAsync<Profiles>(await streamTask);

            this.Profiles = profiles.ProfileList;
        }

        public Profile GetProfile(string profileName)
        {
            var streamTask = client.GetStreamAsync($"{baseUrl}show/{profileName}");
            var profile = JsonSerializer.DeserializeAsync<Profile>(streamTask.Result).Result;

            return profile;
        }

        public String[] ProfileNames()
        {
            if (!isConnected)
            {
                Connect();
            }

            var ProfileNames = new string[Profiles.Length];
            for (int x = 0; x < Profiles.Length; x++)
            {
                ProfileNames[x] = Profiles[x].Name;
            }
            return ProfileNames;
        }

        public Image GetImage(string data)
        {
            string src = buildImageSource(data);
            var streamTask = client.GetStreamAsync(src);
            var stream = streamTask.Result;

            Svg.SvgDocument svgImg = Svg.SvgDocument.Open<Svg.SvgDocument>(stream);
            svgImg.Width = 72;
            svgImg.Height = 72;

            return svgImg.Draw();
        }

        private string buildImageSource(string data)
        {
            if (data.StartsWith("/"))
            {
                return baseUrl + data + "?width=72&height=72";
            }
            if (data.StartsWith("data:"))
            {
                return data;
            }
            return url + "/webclient/assets/" + data;
        }

        public void SetProfile(string profileName)
        {
            Message msg = new();
            msg.Profile = profileName;
            msg.Command = "change";
            SendMessage(msg).Wait();
        }
        internal void SendClick(string profileName, string pageName, string actionName)
        {
            var message = new Message();
            message.Profile = profileName;
            message.Action = actionName;
            message.Page = pageName;
            message.Command = "click";

            var json = JsonSerializer.Serialize(message);
            StringContent httpContent = new StringContent(json, System.Text.Encoding.UTF8, "application/json");


            var actionUrl = baseUrl + "action/" + profileName + "/" + actionName;
            var stringTask = client.PostAsync(actionUrl, httpContent);
            var result = stringTask.Result;
            Console.WriteLine($"button press result: {result}");
        }

        static UTF8Encoding encoder = new UTF8Encoding();

        public async Task SendMessage(Message msg)
        {
            string jsonStr = JsonSerializer.Serialize(msg);
            byte[] buffer = encoder.GetBytes(jsonStr);
            await ws.SendAsync(new ArraySegment<byte>(buffer), WebSocketMessageType.Text, true, CancellationToken.None);
        }

        private async Task Receive(CancellationToken token)
        {
            byte[] buffer = new byte[2048];
            while ((ws.State == WebSocketState.Open) && !token.IsCancellationRequested)
            {
                var result = await ws.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);
                if (result.MessageType == WebSocketMessageType.Close)
                {
                    await ws.CloseAsync(WebSocketCloseStatus.NormalClosure, string.Empty, CancellationToken.None);
                }
                else
                {

                    var str = encoder.GetString(buffer, 0, result.Count).Replace("\0", string.Empty);
                    try
                    {
                        var message = JsonSerializer.Deserialize<Message>(str);
                        MessageReceived(this, new MessageReceived(message));
                    }
                    catch (Exception e)
                    {
                        Console.WriteLine($"{str} \r\n {e.Message}");
                    }
                }
            }
        }
    }
}