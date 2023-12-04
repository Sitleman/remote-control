#include <WiFi.h>
#include <WebSocketClient.h>
#include <ArduinoJson.h>
#include <Wire.h>
#include <Adafruit_GFX.h>
#include <Adafruit_SSD1306.h>

#define BUTTON_PIN 12

//const char* ssid = "newton_2.4G";
//const char* password = "rekashastia";

char path[] = "/ws-control";
//char host[] = "192.168.1.78";


const char* ssid = "OnePlus 7";
const char* password = "12121234";
char host[] = "192.168.191.193";

#define SCREEN_WIDTH 128 // OLED display width, in pixels
#define SCREEN_HEIGHT 32 // OLED display height, in pixels

// Declaration for an SSD1306 display connected to I2C (SDA, SCL pins)
Adafruit_SSD1306 display(SCREEN_WIDTH, SCREEN_HEIGHT, &Wire, -1);
 
WebSocketClient webSocketClient;
WiFiClient client;

void setup() {
  Serial.begin(115200);
  delay(10);

  if(!display.begin(SSD1306_SWITCHCAPVCC, 0x3C)) { // Address 0x3D for 128x64
    Serial.println(F("SSD1306 allocation failed"));
    for(;;);
  }
  // We start by connecting to a WiFi network

  Serial.println();
  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(ssid);
  
  WiFi.begin(ssid, password);
  
  while (WiFi.status() != WL_CONNECTED) {
    delay(1000);
    Serial.print(".");
  }

  Serial.println("");
  Serial.println("WiFi connected");  
  Serial.println("IP address: ");
  Serial.println(WiFi.localIP());

  delay(5000);
  

  // Connect to the websocket server
  if (client.connect(host, 8080)) {
    Serial.println("Connected");
  } else {
    Serial.println("Connection failed.");
    while(1) {
      // Hang on failure
    }
  }

  // Handshake with the server
  webSocketClient.path = path;
  webSocketClient.host = host;
  if (webSocketClient.handshake(client)) {
    Serial.println("Handshake successful");
  } else {
    Serial.println("Handshake failed.");
    while(1) {
      // Hang on failure
    }  
  }

  pinMode(LED_BUILTIN, OUTPUT);
  pinMode(BUTTON_PIN, INPUT_PULLUP);
}

long long button_delay = 0;

void loop() {
  String data;
  StaticJsonDocument<200> doc;
  
  if (client.connected()) {
    webSocketClient.getData(data);
    if (data.length() > 0) {
      Serial.print("Received data: ");
      Serial.println(data);

       DeserializationError error = deserializeJson(doc, data);
       if (error){
          Serial.print("DeserializationError: ");
          Serial.println(error.f_str());
          return;   
       }
       if (doc["elem"] == "led") {
        if (doc["data"] == "on") {
          digitalWrite(LED_BUILTIN, HIGH);
        } else {
          digitalWrite(LED_BUILTIN, LOW);
        }
      }
      
      if (doc["elem"] == "display") {
        display.clearDisplay();
        display.setCursor(0, 10);
        display.setTextSize(1);
        display.setTextColor(WHITE);
        // Display static text
        String printData = doc["data"];
        display.println(printData);
        display.display(); 
      
      }
    }
  }
  delay(50);
}

//  int button = digitalRead(BUTTON_PIN);
//  if (millis() - button_delay > 1000 and button == 0) {
//    button_delay = millis();
//    digitalWrite(LED_BUILTIN, HIGH);
//    Serial.println("Button pressed!");
//    sendMessage();
//  }
//  if (millis() - button_delay > 500) {
//    digitalWrite(LED_BUILTIN, LOW);
//  }
//  
//  delay(50);
//}

//void sendMessage() {
//
//   if (client.connected()) {
//    webSocketClient.getData(data);
//    if (data.length() > 0) {
//      Serial.print("Received data: ");
//      Serial.println(data);
//
//      JsonObject& root = jsonBuffer.parseObject(data);
//      if (root["elem"] == "led") {
//        if (root["data"] == "on") {
//          digitalWrite(LED_BUILTIN, HIGH);
//        } else {
//          digitalWrite(LED_BUILTIN, LOW);
//        }
//      }
//    }
//
//    
//    
////    pinMode(1, INPUT);
////    data = String(analogRead(1));
//   } else {
//      Serial.println("Client disconnected.");
//   }
//
////  if (client.connected()) {
////    webSocketClient.sendData("Hello, it's ESP32!");
////    Serial.println("Successful send message to websocket.");
////  } else {
////    Serial.println("Client disconnected.");
////  }
//}
//
////  String data;
////
////  if (client.connected()) {
// 
////    webSocketClient.getData(data);
////    if (data.length() > 0) {
////      Serial.print("Received data: ");
////      Serial.println(data);
////    }
//    
////    pinMode(1, INPUT);
////    data = String(analogRead(1));
////    
////    webSocketClient.sendData(data);
////    
////  } else {
////    Serial.println("Client disconnected.");
////    while (1) {
////      // Hang on disconnect.
////    }
////  }
////  
////  // wait to fully let the client disconnect
////  delay(3000);
//  
////}
