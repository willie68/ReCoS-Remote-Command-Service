using System;
using System.Collections.Generic;
using System.Drawing;
using System.Linq;
using System.Net.Http;
using System.Net.Http.Headers;
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

        readonly HttpClient client = new HttpClient();
        readonly ClientWebSocket ws = new ClientWebSocket();

        public int ImageWidth { get; set; }
        public int ImageHeight { get; set; }

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
            this.ImageHeight = 100;
            this.ImageWidth = 100;
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
                Uri uri = new Uri(url);
                Uri wsUri = new Uri($"ws://{uri.Host}:{uri.Port}/ws");
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

            var task = client.GetAsync(src, HttpCompletionOption.ResponseHeadersRead);
            using (var response = task.Result)
            {
                if (response.IsSuccessStatusCode)
                {
                    var stream = response.Content.ReadAsStreamAsync().Result;
                    var trailingHeaders = response.Headers;
                    IEnumerable<string> headerValues;
                    string contentType = GetHeaderString(response.Content.Headers, "Content-Type");
                    if ("image/svg+xml".Equals(contentType))
                    {
                        Svg.SvgDocument svgImg = Svg.SvgDocument.Open<Svg.SvgDocument>(stream);
                        svgImg.Width = ImageWidth;
                        svgImg.Height = ImageHeight;

                        return svgImg.Draw();
                    }
                    else if ("image/bmp".Equals(contentType) || "image/png".Equals(contentType))
                    {
                        return Image.FromStream(stream);
                    }
                }
            }
            return null;
        }

        private string buildImageSource(string data)
        {
            if (data.StartsWith("/"))
            {
                return url + data + $"?width={ImageWidth}&height={ImageHeight}";
            }
            if (data.StartsWith("data:"))
            {
                return data;
            }
            return url + "/webclient/assets/" + data;
        }

        public void SetProfile(string profileName)
        {
            Message msg = new Message();
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
            //Console.WriteLine($"button press result: {result}");
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
        private static string GetHeaderString(HttpHeaders headers, string name)
        {
            IEnumerable<string> values;

            if (headers.TryGetValues(name, out values))
            {
                return values.FirstOrDefault();
            }

            return null;
        }
    }
}