using System;
using System.Drawing;
using System.Net.Http;
using System.Text.Json;
using System.Threading.Tasks;

namespace ReCoS
{

    public class RecosClient
    {
        private const string DefaultURL = "http://127.0.0.1:9280";

        readonly HttpClient client = new HttpClient();

        private string url;
        private string baseUrl;
        private bool isConnected;
        private Profile[] Profiles;


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
                isConnected = true;
            }
            catch (Exception e)
            {
                Console.Error.WriteLine($"can't connect to ReCoS {e.Message}");
            }
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
    }
}