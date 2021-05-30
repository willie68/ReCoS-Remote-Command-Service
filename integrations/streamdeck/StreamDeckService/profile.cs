using System.Text.Json.Serialization;

namespace ReCoS
{
    public class Profiles
    {
        [JsonPropertyName("profiles")]
        public Profile[] ProfileList { get; set; }
    }

    public class Button
    {
        public Button(Action action)
        {
            Action = action;
        }

        public Action Action { get; set; }
        public string Text { get; set; }
    }

    public class Profile
    {
        [JsonPropertyName("name")]
        public string Name { get; set; }
        [JsonPropertyName("description")]
        public string Description { get; set; }
        [JsonPropertyName("pages")]
        public Page[] Pages { get; set; }
        [JsonPropertyName("actions")]
        public Action[] Actions { get; set; }
    }

    public class Page
    {
        [JsonPropertyName("name")]
        public string Name { get; set; }
        [JsonPropertyName("description")]
        public string Description { get; set; }

        [JsonPropertyName("icon")]
        public string Icon { get; set; }

        [JsonPropertyName("columns")]
        public int Columns { get; set; }

        [JsonPropertyName("rows")]
        public int Rows { get; set; }

        [JsonPropertyName("toolbar")]
        public string Toolbar { get; set; }

        [JsonPropertyName("cells")]
        public string[] Cells { get; set; }
    }

    public class Action
    {
        [JsonPropertyName("type")]
        public string Type { get; set; }

        [JsonPropertyName("name")]
        public string Name { get; set; }

        [JsonPropertyName("title")]
        public string Title { get; set; }

        [JsonPropertyName("description")]
        public string Description { get; set; }

        [JsonPropertyName("icon")]
        public string Icon { get; set; }

        [JsonPropertyName("fontsize")]
        public int Fontsize { get; set; }

        [JsonPropertyName("fontcolor")]
        public string Fontcolor { get; set; }

        [JsonPropertyName("outlined")]
        public bool Outlined { get; set; }
    }

    public class Message
    {
        [JsonPropertyName("profile")]
        public string Profile { get; set; }

        [JsonPropertyName("action")]
        public string Action { get; set; }

        [JsonPropertyName("page")]
        public string Page { get; set; }

        [JsonPropertyName("imageurl")]
        public string ImageURL { get; set; }

        [JsonPropertyName("title")]
        public string Title { get; set; }

        [JsonPropertyName("text")]
        public string Text { get; set; }

        [JsonPropertyName("state")]
        public int State { get; set; }

        [JsonPropertyName("command")]
        public string Command { get; set; }
    }
}
