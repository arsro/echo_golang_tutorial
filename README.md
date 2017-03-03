# echo_golang
1. ```clone this repository```
2. ```docker build -t echo_image .```
3. ```sh ./docker_run.sh```


## メモ

### CURD
- 参考
	<a href="http://qiita.com/CST_negi/items/5e276ddc0412cefef7e3" >Golang+Echo+dbrでMySQLのCRUDをする／JSONでDBの値を返却する話 - Qiita</a>
 
####  JSON形式でのリクエストを受けとる
以下のリクエストを送る
```$curl -v POST -H "Content-Type: application/json" "http://localhost:3000/users" -d '{"name" : "Anthony Kiedis", "age":20}'```


```go
// modelsで設定したUser構造体をインスタンス化(構造体の初期化方法は複数ある)
user := new(models.User)
c.JSON(http.StatusCreated, user)
//->{"name":"","age":0}

//この部分でリクエストをUser構造体にバインドしている
//User構造体で定義したtag(`json:"hoge"`)の部分でどのリクエストパラメータのキーをバインドするかを設定しているっぽい
if err := c.Bind(user); err != nil {
	return err
}
c.JSON(http.StatusCreated, user)
//->{"name":"Anthony Kiedis","age":20}
```

- 参考
	- <a href="https://medium.com/@kyawmyintthein/building-golang-restful-api-with-echo-framework-1-422abc78e3a7#.3g32545eo" >Building Golang Restful API with Echo framework (part - 1) – Medium</a>
	- <a href="http://omiend.hatenablog.jp/entry/2017/01/31/203314" >その２ ファイル分割してみる - golang製Web Framework 「echo」事始め - the industrial</a>
	- <a href="http://dev.classmethod.jp/server-side/language/golang-5/" >急いで学ぶGo lang#5 構造体 ｜ Developers.IO</a>

#### DBへの接続, SELECT

1. 別コンテナでmysqlを立ち上げて、echoコンテナと繋げる
- docker-composeを利用
```yml
mysql:
  image: mysql
  ports:
    - "3306:3306"
  expose:
    - 3306
  environment:
    MYSQL_ROOT_PASSWORD: root
    MYSQL_DATABASE: test_db
    MYSQL_USER: user
    MYSQL_PASSWORD: root
```

``` docker-compose up -d ```でmysqlコンテナを起動

- ```sh ./docker_run.sh```を実行。echoコンテナにアクセスする。

- echoコンテナ内に接続先のmysqlコンテナの環境変数が入っていることを確認

```
# env | grep MYSQL

MYSQL_ENV_MYSQL_DATABASE=test_db
MYSQL_ENV_MYSQL_ROOT_PASSWORD=root
MYSQL_ENV_GOSU_VERSION=1.7
MYSQL_PORT_3306_TCP_PORT=3306
MYSQL_PORT_3306_TCP=tcp://172.17.0.2:3306
MYSQL_ENV_MYSQL_USER=tomorrow
MYSQL_ENV_MYSQL_PASSWORD=root
MYSQL_ENV_MYSQL_VERSION=5.7.17-1debian8
MYSQL_ENV_no_proxy=*.local, 169.254/16
MYSQL_NAME=/echo/mysql
MYSQL_PORT_3306_TCP_PROTO=tcp
MYSQL_PORT_3306_TCP_ADDR=172.17.0.2
MYSQL_ENV_MYSQL_MAJOR=5.7
MYSQL_PORT=tcp://172.17.0.2:3306

```

- ```/etc/hosts```にも接続先のエントリが追加されている

```
# cat /etc/hosts | grep mysql
172.17.0.2	mysql bdc58137530b mysql_mysql_1

// mysql			: --linkオプションで設定したエアリアス名
// bdc58137530b 	: mysqlコンテナのコンテナID
// mysql_mysql_1 	: mysqlコンテナのコンテナ名
```

- 参考
	- <a href="http://qiita.com/Arturias/items/75828479c1f9eb8d43fa" >Docker の基本学習 ~ コンテナ間のリンク - Qiita</a>
	- <a href="http://qiita.com/astrsk_hori/items/e3d6c237d68be1a6f548" >dockerでmysqlを使う - Qiita</a>
	- <a href="http://48n.jp/blog/2016/09/16/links-container-for-docker/" >Dockerのコンテナ間を繋ぐLinksを使ってPHPとMySQLコンテナを連携させる - SHOYAN BLOG</a>
	
	
2. golangからデータベースに接続する
- 以下の```struct```を使用

``` go
package db
// usersテーブルをid, name, ageで設定している場合は以下の構成に
type (
	Users struct {
				Id   	int64  `db:"id"`
				Name	string `db:"name"`
				Age		int	   `db:"age"`
	}
)
```


```go
package controllers

import (
	"github.com/labstack/echo"
	_ "github.com/go-sql-driver/mysql" 	//←これ!
	"github.com/gocraft/dbr" 			//←これ!

	"net/http"
	
	"../db"
)



func ConnectDB(c echo.Context ) error{
	// mysqlへの接続設定
	// dbr.Open("mysql", "[mysql_username]:[mysql_passwd]@tcp( [接続先のmysqlコンテナのホストネーム] :3306)/[接続DB名]", nil)
	conn, err := dbr.Open("mysql", "tomorrow:root@tcp(mysql:3306)/test_db", nil)
	if err != nil{
		return err
	}
	sess := conn.NewSession(nil)

	var user []db.Users
	// select * from users
	sess.Select("*").From("users").Load(&user)
	
	// [{"Id":1,"Name":"tomorrow","Age":24}]
	return c.JSON(http.StatusCreated, user)
}
```

```sql
-- usersテーブル
CREATE TABLE `users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(16) DEFAULT NULL,
  `age` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;
```

- 参考
	- <a href="https://eurie.co.jp/blog/engineering/2015/12/go-lang-ormapper-dbr" >dbr – Go 言語 O/R Mapper の紹介</a>


#### INSERT
- テーブル構造は上記のものと同じ

- 以下のPOSTリクエストを送信する
```
curl -v POST -H "Content-Type: application/json" "http://localhost:3000/users" -d '{"id" : 1, "name" : "Taro Yamada", "age":20}'
```

- DBに接続(上記の方法と同様)

```go
func PostUser(c echo.Context) error {
	conn, err := dbr.Open("mysql", "tomorrow:root@tcp(mysql:3306)/test_db", nil)
	if err != nil{
		return err
	}
	sess := conn.NewSession(nil)
	
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return err
	}
	c.JSON(http.StatusCreated, user)
	
	result, err := sess.InsertInto("users").
					Columns("id", "name", "age").
					Record(user).Exec()

	if err != nil {
		return err
	} else {
		return c.JSON(http.StatusCreated, "SUCCESS!")
	}
}
```

- ```struct``` を初期化して ```INSERT``` 文の ```Values``` とする場合は、```Record()``` を用いる

	- その他方法としては```Record()```ではなく、```Values()```を用いる方法も
	```go
	result, err := sess.InsertInto("users").
						Columns("id", "name", "agee").
						Values(1, "Yamada", 20).
						Exec()
	```
		
#### GETクエリを取得する
- 以下のリクエストを送信する
```
curl -v GET -H "Content-Type: application/json" "http://localhost:3000/users?id=3&name=hogehoge&age=25"
```

```go
import(
	"strconv"
)

func GetUser(c echo.Context) error {
	
	// string -> int
	/**
	 *  Package strconv
	 *  see https://golang.org/pkg/strconv/
	 */
	id,_ := strconv.Atoi(c.QueryParam("id"))
	var name string = c.QueryParam("name")
	age,_ := strconv.Atoi(c.QueryParam("age"))

	var user *models.User = newUser(id, name, age)
	
	models.Users[user.ID] = user
	return c.JSON(http.StatusCreated, models.Users)
}

//構造体の初期化
func newUser(id int, name string, age int) *models.User{
	user := new(models.User)
	user.ID = id
	user.Name = name
	user.Age = age

	return user
}

```

- ```int -> string```へのキャストに関してはstrconvパッケージを使う。
- 参考
	- <a href="http://matope.hatenablog.com/entry/2014/04/22/101127" >Golangでの文字列・数値変換 - 小野マトペの納豆ペペロンチーノ日記</a>

#### UPDATE
以下のリクエストを送信する

```
curl -v -X PUT -H "Content-Type: application/json" "http://localhost:3000/users" -d '{"id": 2, "name" :"Yamada Hanako"}'
```

```go
result, err := sess.Update("users").
				Set("name", user.Name).
				Where("id = ?", user.ID).
				Exec()
```

#### DELETE
以下のリクエストを送信する

```
curl -v -X DELETE -H "Content-Type: application/json" "http://localhost:3000/users" -d '{"id": 2}'
```

```go
result, err := sess.DeleteFrom("users").
				Where("id = ?", user.ID).
				Exec()
```



-------------

------------

## 参考
- 公式
	- <a href="https://echo.labstack.com/cookbook/hello-world" >Hello World | Echo - High performance, minimalist Go web framework</a>
	- <a href="https://github.com/labstack/echo/" >labstack/echo: High performance, minimalist Go web framework</a>
	

	
