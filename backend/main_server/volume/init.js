// Подключаемся к базе данных admin (где у root-пользователя административные права)
var admin = db.getSiblingDB("admin");

// Аутентификация root-пользователя
var success = admin.auth("root", "root");

if (success) {
  print("Authentication successful for root user.");
  // Подключаемся к создаваемой базе данных
  var test = db.getSiblingDB("tokens");
  // Создаем базу данных
  test.createCollection("tokens");

  // Создаем пользователя с правами на чтение и запись для базы данных mydb
  test.createUser({
    user: "wallets",
    pwd: "wallets",
    roles: [{ role: "readWrite", db: "tokens" }]
  });
    var result=test.auth("wallets","wallets")
    if (result){
      print("Authentication successful for wallets");
    }else{
      print("Authentication failed for wallets");
    }
} else {
  print("Authentication failed for root user.");
}


