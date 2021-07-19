//#include <ArduinoWebsockets.h>
#include <WebSockets2_Generic.h>
#include <WiFi.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>
#include <Adafruit_GFX.h> //Grafik Bibliothek
#include <Adafruit_ILI9341.h> // Display Treiber
#include <upng.h> //PNG decoder
#include <XPT2046_Touchscreen.h> //Touchscreen Treiber
#include <TouchEvent.h> //Auswertung von Touchscreen Ereignissen
#include <TFTForm.h> //Formulare am Touchscreen
#include <SPIFFS.h>  //Filesystem

//*********************Konfiguration*********************************
//Diese Einstellungen werden in einer späteren Version am Display einstellbar
#define SSID "93w8562k"
#define PKEY "akteon00"

#define SERVER "192.168.178.45"
#define PORT 9280
#define API "/api/v1"
#define ICONS "/config/icons/"
#define WEBSOCKET "/ws"
#define ASSETS "/webclient/assets"
#define PROFILEINDEX 0
//********************************************************************
#define DEBUG_WEBSOCKETS_PORT     Serial
// Debug Level from 0 to 4
#define _WEBSOCKETS_LOGLEVEL_     3

//Definitionen der verwendeten Pins
#define TFT_CS   5
#define TFT_DC   4
#define TFT_MOSI 23
#define TFT_CLK  18
#define TFT_RST  22
#define TFT_MISO 19
#define TFT_LED  15
#define TOUCH_CS 14
#define TOUCH_IRQ 27 //Für AZ-Touch alte Version #define TOUCH_IRQ 2
#define TOUCH_ROTATION 1 //Muss für 2.4 Zoll Display 1 sein

typedef
struct __attribute__((__packed__)) {
  char code[2];
  uint32_t fileSize;
  uint32_t creatorBytes;
  uint32_t imageOffset;
  uint32_t headerSize;
  uint32_t width;
  uint32_t height;
  uint16_t planes;
  uint16_t depth;
  uint32_t format;
} BITMAPHDR;

//die Instanzen der verwendeten Klassen
Adafruit_ILI9341 tft = Adafruit_ILI9341(TFT_CS, TFT_DC, TFT_RST);
XPT2046_Touchscreen touch(TOUCH_CS, TOUCH_IRQ);
TouchEvent tevent(touch);

//using namespace websockets;
using namespace websockets2_generic;

//WebsocketsClient client;
WebsocketsClient client;
//global variables
DynamicJsonDocument profileList(5000);
DynamicJsonDocument profile(10000);
JsonArray profileArray;
JsonArray pages;
JsonArray actions;
uint32_t wsCnt = 0;
boolean reconnect = false;
//Basiswerte für Kalibrierung
uint16_t xMin = 242;
uint16_t yMin = 242;
uint16_t xMax = 3655;
uint16_t yMax = 3877;
boolean active = true;

uint16_t convertColor(uint8_t r, uint8_t g, uint8_t b, uint8_t a, uint16_t bg) {
  r = r / 8;
  g = g / 4;
  b = b / 8;
  uint16_t tftColor = (a == 0) ? bg : r * 2048 + g * 32 + b;
  return tftColor;
}

void displayImage(const unsigned char * buf, uint16_t x, uint16_t y,
                  uint16_t w, uint16_t h, uint8_t elements, boolean rgb,
                  boolean topleft, uint16_t bgcolor) {
  uint16_t xo, yo, cpix;
  uint8_t e[4];
  e[3] = 255; //if no alpha channel display all pixels;
  xo = (w < 80) ? (80 - w) / 2 : 0;
  yo = (h < 80) ? (80 - h) / 2 : 0;
  if ((w == 1) && (h == 1)) {
    if (rgb) {
      cpix = convertColor(buf[0], buf[1], buf[2], 255, bgcolor);
    } else {
      cpix = convertColor(buf[2], buf[1], buf[0], 255, bgcolor);
    }
    //Serial.printf("w = %i, h = %i, c= %x\n",w,h,cpix);
    tft.fillRect(x + 4, y + 4, 72, 72, cpix);
  } else {
    tft.startWrite();
    for (uint16_t row = 0; row < h; row++) {
      for (uint16_t col = 0; col < w; col++) {
        for (uint8_t el = 0; el < elements; el++) {
          e[el] = buf[(row * w + col) * elements + el];
        }
        if (rgb) {
          cpix = convertColor(e[0], e[1], e[2], e[3], bgcolor);
        } else {
          cpix = convertColor(e[2], e[1], e[0], e[3], bgcolor);
        }
        if (((xo + col) < 80) && ((yo + row) < 80)) {
          if (topleft) {
            tft.writePixel(x + xo + col, y + yo + row, cpix);
          } else {
            tft.writePixel(x + xo + col, y + yo + h - row, cpix);
          }
        }
      }
    }
    tft.endWrite();
  }
}

String currentProfile = "";
String currentPage = "";
boolean wsactive = false;

void wsSend(String profile, String command) {
  char buf[500];
  if (wsactive)
  {
    DynamicJsonDocument doc(500);
    doc["profile"] = profile;
    doc["command"] = command;
    uint16_t n = serializeJson(doc, buf);
    client.send(buf, n);
  }
  else
  {
    Serial.println("Not Connected!");
  }
}
void onEventsCallback(WebsocketsEvent event, String data)
{
  if (event == WebsocketsEvent::ConnectionOpened)
  {
    Serial.println("Connnection Opened");
  }
  else if (event == WebsocketsEvent::ConnectionClosed)
  {
    Serial.println("Connnection Closed");
    if (wsactive) {
      //client.close(CloseReason_NormalClosure);
      Serial.println("Reconnect");
      reconnect = true;
    }
  }
  else if (event == WebsocketsEvent::GotPing)
  {
    Serial.println("Got a Ping!");
  }
  else if (event == WebsocketsEvent::GotPong)
  {
    Serial.println("Got a Pong!");
  }
}

String httpGet(String path, String * mimeType) {
  String result;
  const char * keys[] = {"Content-Type"};
  *mimeType = "unknown";
  if (WiFi.status() == WL_CONNECTED) {
    HTTPClient http;
    http.begin(SERVER, PORT, path);
    http.collectHeaders(keys, 1);
    int httpCode = http.GET();
    if (httpCode > 0) {
      if (httpCode == 200) {
        result = http.getString();
        *mimeType = http.header("Content-Type");
      } else {
        result = F("ERR:Got code ");
        result += String(httpCode);
      }
    } else {
      result = F("ERR:Get Failed");
    }
    http.end();
  } else {
    result = F("ERR:Not connected");
  }
  return result;
}

String httpPost(String profile, String action) {
  Serial.println("POST");
  String result;
  if (WiFi.status() == WL_CONNECTED) {
    char msg[500];
    DynamicJsonDocument doc(500);
    doc["profile"] = profile;
    doc["action"] = action;
    doc["comman"] = "click";
    serializeJson(doc, msg);
    Serial.println(msg);
    String path = String(API) + "/action/" + profile + "/" + action;
    HTTPClient http;
    http.begin(SERVER, PORT, path);
    int httpCode = http.POST(msg);
    if (httpCode > 0) {
      if (httpCode == 200) {
        result = http.getString();
      } else {
        result = F("ERR:Got code ");
        result += String(httpCode);
      }
    } else {
      result = F("ERR:Get Failed");
    }
    http.end();
  } else {
    result = F("ERR:Not connected");
  }
  return result;
}

boolean getProfiles() {
  boolean result = false;
  String typ;
  String res = httpGet(String(API) + String("/show"), &typ);
  if (res.startsWith("ERR:")) {
    Serial.println(res);
  } else {
    DeserializationError   error = deserializeJson(profileList, res);
    if (error ) {
      Serial.println("JSON get profiles: ");
      Serial.println(error.c_str());
    } else {
      profileArray = profileList["profiles"].as<JsonArray>();
      result = true;
    }
  }
  return result;
}

boolean getProfile() {
  boolean result = false;
  String typ;
  String res = httpGet(String(API) + String("/show/") + currentProfile, &typ);
  if (res.startsWith("ERR:")) {
    Serial.println(res);
  } else {
    DeserializationError   error = deserializeJson(profile, res);
    if (error ) {
      Serial.println("JSON get profiles: ");
      Serial.println(error.c_str());
    } else {
      Serial.println(res);
      pages = profile["pages"].as<JsonArray>();
      actions = profile["actions"].as<JsonArray>();
      result = true;
    }
  }

}

void showPng(String img, uint16_t x, uint16_t y, uint16_t bg = ILI9341_BLACK) {
  boolean result = false;
  upng_t * upng;
  uint16_t w, h;
  const unsigned char * pngbmp;
  upng = upng_new_from_bytes((const unsigned char *)img.c_str(), img.length());
  upng_decode(upng);
  w = upng_get_width(upng);
  h = upng_get_height(upng);
  pngbmp = upng_get_buffer(upng);
  displayImage(pngbmp, x, y, w, h, upng_get_components(upng), true, true, bg);
  upng_free(upng);
}

void showBmp(const char * dbuf, uint16_t x, uint16_t y, uint16_t bg = ILI9341_BLACK) {
  BITMAPHDR bmh;
  uint16_t w, h;
  memcpy((char *)&bmh, dbuf, sizeof(bmh));
  w = bmh.width;
  h = bmh.height;
  //Serial.printf("Width %i, Height %i\n",w,h);
  dbuf += bmh.imageOffset;
  displayImage((const unsigned char *)dbuf, x, y, w, h, bmh.depth / 8, false, false, bg);
}


void drawCenter(uint16_t x, uint16_t y, String txt) {
  uint16_t l = txt.length();
  uint8_t o = (l * 6 < 80) ? (80 - l * 6) / 2 : 0;
  tft.setCursor(x + o, y);
  tft.print(txt.substring(0, 79));
}

void displayCell(uint8_t index, JsonObject action, JsonObject message, boolean updt) {
  //Display the image
  uint16_t x = (index % 3) * 80;
  uint16_t y = (index / 3) * 80;
  String img = "";
  if (action.containsKey("icon")) img = action["icon"].as<String>();
  String img1 = "";
  String typ = "";
  const char * buf;
  if (updt && message.containsKey("imageurl")) img1 = message["imageurl"].as<String>();
  if (img1 != "") img = img1;
  if (img != "") {
    if (img.endsWith(".svg")) {
      img = String(API) + String(ICONS) + img;
    }
    if (!img.startsWith("/")) img = String("/webclient/assets/") + img;
    img += String("?width=72&height=72");
    //Serial.println(img);
    String data = httpGet(img, &typ);
    //Serial.println(data.length());
    buf = data.c_str();
    if (data.startsWith("ERR:")) {
      Serial.println(data);
      Serial.println(img);
    } else {
      //Serial.println(typ);
      if (typ.endsWith("/png")) showPng(data, x, y, ILI9341_BLUE);
      if (typ.endsWith("/bmp")) showBmp(buf, x, y, ILI9341_BLUE);
    }
  }
  //Display Texte
  String title = updt ? message["title"] : action["title"];
  String text = updt ? message["text"].as<String>() : "";
  String scol = action["fontcolor"];
  uint16_t col = ILI9341_BLACK;
  if (scol == "white") col = ILI9341_WHITE;
  tft.setTextColor(col);
  if (title != "") drawCenter(x, y + 20, title);
  if (text != "") drawCenter(x, y + 40, text);
}


JsonObject getObject(JsonArray arr, String name)   {
  int16_t i = arr.size();
  while ((i-- > 0) && (arr[i]["name"].as<String>() != name));
  if (i >= 0) return arr[i];
}

int16_t getIndex(JsonArray arr, String name)   {
  int16_t i = arr.size();
  while ((i-- > 0) && (arr[i]["name"].as<String>() != name));
  return i;
}

void showPage() {
  tft.fillScreen(ILI9341_WHITE);
  uint8_t row = 0;
  uint8_t col = 0;
  String actionName, icon;
  JsonObject action;
  JsonObject pag = getObject(pages, currentPage);
  JsonArray cells = pag["cells"].as<JsonArray>();
  uint16_t cnt = cells.size();
  for (uint16_t i = 0; i < cnt; i++) {
    displayCell(i, getObject(actions, cells[i].as<String>()), actions[0], false);
  }
  active = true;
}

void wsMessage(String msg) {
  if (!active) return;
  DynamicJsonDocument doc(2000);
  uint16_t x, y;

  DeserializationError   error = deserializeJson(doc, msg);
  if (error ) {
    Serial.println("JSON WS message: ");
    Serial.println(error.c_str());
  } else {
    JsonObject pag = getObject(pages, currentPage);
    String npag = doc["page"];
    if (npag != "") {
      //page switch
      currentPage = npag;
      showPage();
    } else {
      JsonArray cells = pag["cells"].as<JsonArray>();
      uint16_t cnt = cells.size();
      for (uint16_t i = 0; i < cnt; i++) {
        if (cells[i].as<String>() == doc["action"].as<String>()) {
          displayCell(i, getObject(actions, doc["action"].as<String>()), doc.as<JsonObject>(), true);
        }
      }
    }
  }
}

//wird immer aufgerufen wenn ein Touchscreen-Ereignis auftritt
void onTouchEvent(int16_t x, int16_t y, EV event) {
  Serial.println("Touch");
  if (event == EV::EVT_CLICK) {
    if (active) {
      Serial.printf("Click %i,%i\n", x, y);
      JsonObject pag = getObject(pages, currentPage);
      JsonArray cells = pag["cells"].as<JsonArray>();
      uint8_t col = x / 80;
      uint8_t row = (y / 80);
      Serial.printf("row = %i, col = %i \n", row, col);
      Serial.println(cells.size());
      String name = cells[row * 3 + col];
      Serial.println(name);
      JsonObject action = getObject(actions, name);
      if (action["type"].as<String>() == "SINGLE") {
        String res = httpPost(currentProfile, action["name"].as<String>());
        Serial.print("POST result = ");
        Serial.println(res);
      }
    }
  }
}

void switchPage(boolean forward) {
  uint8_t cnt = pages.size();
  uint16_t cur = getIndex(pages, currentPage);
  if (cur < 0) {
    cur = forward ? 0 : cnt;
  } else {
    cur = forward ? cur + 1 : cur - 1;
    if (cur >= cnt) cur = 0;
    if (cur < 0) cur = cnt;
  }
  currentPage = pages[cur]["name"].as<String>();
  showPage();
}

void showProfiles() {
  active = false;
  tft.fillScreen(ILI9341_WHITE);
  tft.setTextColor(ILI9341_BLACK);
  tft.setCursor(40, 10);
  tft.print("Profile");
}

void showConfig() {
  active = false;
  tft.fillScreen(ILI9341_YELLOW);
  tft.setTextColor(ILI9341_BLACK);
  tft.setCursor(40, 10);
  tft.print("Config");
}

//wird immer aufgerufen wenn über den Touchscreen gewischt wurde
//direction: 0=links, 1=rechts, 2=auf, 3=ab
void onTouchSwipe(uint8_t direction) {
  switch (direction) {
    case 0: switchPage(true); break;
    case 1: switchPage(false); break;
    case 2: showProfiles(); break;
    case 3: showConfig(); break;
  }
}
void setup()
{
  Serial.begin(115200);
  while (!Serial);

  Serial.println("\nStarting ESP32-Client on ESP32");
  pinMode(TFT_LED, OUTPUT);
  digitalWrite(TFT_LED, 0);
  tft.begin();
  tft.fillScreen(ILI9341_WHITE);
  //Touchscreen vorbereiten
  touch.begin();
  touch.setRotation(TOUCH_ROTATION);
  tevent.setResolution(tft.width(), tft.height());
  tevent.setDrawMode(false);
  //Callback Funktionen registrieren
  tevent.registerOnAllEvents(onTouchEvent);
  tevent.setSwipe(300, 400);
  tevent.registerOnTouchSwipe(onTouchSwipe);
  tevent.calibrate(xMin, yMin, xMax, yMax);
  // Connect to wifi
  WiFi.begin(SSID, PKEY);

  // Wait some time to connect to wifi
  for (int i = 0; i < 10 && WiFi.status() != WL_CONNECTED; i++)
  {
    Serial.print(".");
    delay(1000);
  }

  // Check if connected to wifi
  if (WiFi.status() != WL_CONNECTED)
  {
    Serial.println("No Wifi!");
    return;
  }

  Serial.print("Connected to Wifi, Connecting to WebSockets Server @");
  Serial.println(SERVER);

  // run callback when messages are received
  client.onMessage([&](WebsocketsMessage message)
  {
    tevent.pollTouchScreen();
    //Serial.print(".");
    wsCnt++;
    if ((wsCnt % 100L) == 0) Serial.println(wsCnt);
    wsMessage(message.data());
  });

  // run callback when events are occuring
  client.onEvent(onEventsCallback);
  // try to connect to Websockets server
  wsactive = client.connect(SERVER, PORT, WEBSOCKET);
  wsSend("none", "change");
  if (getProfiles()) {
    currentProfile = profileArray[PROFILEINDEX]["name"].as<String>();
    if (getProfile()) {
      currentPage = pages[0]["name"].as<String>();
      showPage();
    }
  }
  wsSend(currentProfile, "change");
}

void loop()
{
  tevent.pollTouchScreen();
  // let the websockets client check for incoming messages
  if (client.available())
  {
    client.poll();
  }
  if (reconnect) {
    client = {};
    Serial.println("Reset Client");
    //wsactive = client.connect(SERVER, PORT, WEBSOCKET);
    ESP.restart();
    reconnect = false;
  }

  //delay(500);
}
