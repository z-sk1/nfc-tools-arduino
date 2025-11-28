#include <SPI.h>
#include <MFRC522.h>

const int ssPin = 9;
const int rstPin = 8;

MFRC522 mfrc522(SS_PIN, RST_PIN);

void setup() {
  Serial.begin(9600);
  SPI.begin();
  mfrc522.PCD_Init();
  delay(50);

  Serial.println("Arduino ready!");
}

void loop() {
  if (Serial.available()) {
    String cmd = Serial.readStringUntil('\n');
    cmd.trim();

    String arg;

    int spaceIndex = cmd.indexOf(' ');
    if (spaceIndex > 0) {
      arg = cmd.substring(spaceIndex + 1);
      arg.trim();
      cmd = cmd.substring(0, spaceIndex);
    }
  
    if (cmd == "nfcRead") {
      readTag();
      return;
      
    } else if (cmd == "nfcWrite") {
      writeText(arg);
      return;

    } else {
      Serial.println("unknown command: " + cmd);

    }
  }
}

void readTag() {
  if (!mfrc522.PICC_IsNewCardPresent() || !mfrc522.PICC_ReadCardSerial()) {
    Serial.println("ERR NoTag");
    return;
  }

  // ---- Simple Read (reads block 1) ----
  MFRC522::StatusCode status;
  byte buffer[18];
  byte size = sizeof(buffer);

  status = mfrc522.MIFARE_Read(1, buffer, &size);

  if (status != MFRC522::STATUS_OK) {
    Serial.println("ERR ReadFail");
    return;
  }

  // Convert to string
  String data = "";
  for (int i = 0; i < 16; i++) {
    if (buffer[i] == 0x00) break;
    data += (char)buffer[i];
  }

  Serial.println("DATA " + data);
}

void writeText(String txt) {
  if (!mfrc522.PICC_IsNewCardPresent() || !mfrc522.PICC_ReadCardSerial()) {
    Serial.println("ERR NoTag");
    return;
  }

  // Ensure max 16 chars (block size)
  txt = txt.substring(0, 16);

  byte buffer[16];
  for (int i = 0; i < 16; i++) {
    if (i < txt.length())
      buffer[i] = txt[i];
    else
      buffer[i] = 0x00;
  }

  MFRC522::StatusCode status;

  status = mfrc522.MIFARE_Write(1, buffer, 16);
  if (status != MFRC522::STATUS_OK) {
    Serial.println("ERR WriteFail");
    return;
  }

  Serial.println("Succesfully written: " + txt);
}